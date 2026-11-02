package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ProcedureReversibility clasifica procedimientos clínicos según su reversibilidad (AX-012).
type ProcedureReversibility string

// AgentType identifica quién autoriza una acción clínica (humano vs IA).
type AgentType string

// EthicsReviewStatus representa el estado de aprobación ética de una regla clínica (AX-009).
type EthicsReviewStatus string

const (
	ProcedureReversible   ProcedureReversibility = "REVERSIBLE"
	ProcedureIrreversible ProcedureReversibility = "IRREVERSIBLE"

	AgentPractitioner AgentType = "Practitioner" // Humano consciente y matriculado
	AgentDevice       AgentType = "Device"       // MedGemma 2B INT8 (prohibido como autorizador)

	StatusApprovedByCommittee EthicsReviewStatus = "APPROVED_BY_ETHICS_COMMITTEE"
)

// Errores deterministas del Motor Fail-Close.
var (
	ErrRuleNotApprovedByCommittee     = errors.New("AX-009 violation: clinical rule lacks APPROVED_BY_ETHICS_COMMITTEE status signed by CMIO")
	ErrMissingResponsibilityChain     = errors.New("AX-008 violation: DeveloperVersion, OperatorID, and valid SNOMED/FHIR schema are mandatory")
	ErrAIAuthorizationForbidden       = errors.New("AX-012 violation: AI Disarmament enforced; Device agent cannot authorize clinical execution")
	ErrIrreversibleRequiresSupervisor = errors.New("AX-012 violation: IRREVERSIBLE procedure requires human Practitioner supervisor ID and Titan M3 signature")
	ErrInvalidFHIRorSNOMEDTaxonomy    = errors.New("Fail-Close veto: binary deterministic check failed on HL7 FHIR R4 schema or SNOMED CT taxonomy")
)

// VerificationReceipt es la prueba criptográfica almacenada en el registro forense inmutable (AX-011).
type VerificationReceipt struct {
	EventID           string    `json:"event_id"`
	ModelVersion      string    `json:"model_version"`
	DatasetHash       string    `json:"dataset_hash"`
	SNOMEDCode        string    `json:"snomed_code"`
	FHIRResourceValid bool      `json:"fhir_resource_valid"`
	OperatorID        string    `json:"operator_id"`
	SupervisorID      string    `json:"supervisor_id,omitempty"`
	TitanM3Signature  string    `json:"titan_m3_signature"`
	TimestampUTC      time.Time `json:"timestamp_utc"`
	ReceiptSHA256     string    `json:"receipt_sha256"`
}

// ComputeHash genera el digest SHA-256 determinista para el registro forense inmutable.
func (r *VerificationReceipt) ComputeHash() string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%t|%s|%s|%s|%s",
		r.EventID,
		r.ModelVersion,
		r.DatasetHash,
		r.SNOMEDCode,
		r.FHIRResourceValid,
		r.OperatorID,
		r.SupervisorID,
		r.TitanM3Signature,
		r.TimestampUTC.UTC().Format(time.RFC3339Nano),
	)
	sum := sha256.Sum256([]byte(payload))
	r.ReceiptSHA256 = hex.EncodeToString(sum[:])
	return r.ReceiptSHA256
}

// ForensicSinkRepository abstrae el único registro forense inmutable (002_extend_cognitive_audit_log_atdp.sql).
type ForensicSinkRepository interface {
	AppendForensicRecord(receipt VerificationReceipt, status string, atdpDisagreementPayload *string) error
}

// TaxonomyValidator verifica códigos SNOMED CT y conformidad HL7 FHIR R4 aprobados por el CMIO.
type TaxonomyValidator interface {
	IsValidFHIRAndSNOMED(snomedCode string, fhirPayload []byte) bool
}

// IrreversibleDecisionGuard impone AX-008, AX-009, AX-011 y AX-012 (Desarme de la IA).
type IrreversibleDecisionGuard struct {
	forensicSink ForensicSinkRepository
	taxonomy     TaxonomyValidator
}

// NewIrreversibleDecisionGuard construye el guard con sus dependencias.
func NewIrreversibleDecisionGuard(sink ForensicSinkRepository, tax TaxonomyValidator) *IrreversibleDecisionGuard {
	return &IrreversibleDecisionGuard{
		forensicSink: sink,
		taxonomy:     tax,
	}
}

// ValidateAndSeal ejecuta la evaluación binaria determinista Fail-Close.
func (g *IrreversibleDecisionGuard) ValidateAndSeal(
	eventID string,
	ruleStatus EthicsReviewStatus,
	reversibility ProcedureReversibility,
	authorizingAgent AgentType,
	modelVersion string,
	datasetHash string,
	snomedCode string,
	fhirPayload []byte,
	operatorID string,
	supervisorID *string,
	titanM3Signature *string,
) (*VerificationReceipt, error) {

	// 1. AX-009: Verificar estado de aprobación ética homologado por CMIO
	if ruleStatus != StatusApprovedByCommittee {
		return nil, ErrRuleNotApprovedByCommittee
	}

	// 2. AX-008: Verificar identificadores explícitos de cadena de responsabilidad
	if modelVersion == "" || operatorID == "" || snomedCode == "" {
		return nil, ErrMissingResponsibilityChain
	}

	// 3. Check binario determinista Fail-Close: integridad HL7 FHIR R4 + SNOMED CT
	if !g.taxonomy.IsValidFHIRAndSNOMED(snomedCode, fhirPayload) {
		return nil, ErrInvalidFHIRorSNOMEDTaxonomy
	}

	// 4. AX-012: Desarme de la IA — MedGemma 2B INT8 (Device) NUNCA puede ser agente autorizador
	if authorizingAgent == AgentDevice {
		return nil, ErrAIAuthorizationForbidden
	}

	// 5. AX-012: Procedimientos irreversibles requieren estrictamente Supervisor humano + firma Titan M3
	supID := ""
	sig := ""
	if reversibility == ProcedureIrreversible {
		if authorizingAgent != AgentPractitioner || supervisorID == nil || *supervisorID == "" || titanM3Signature == nil || *titanM3Signature == "" {
			return nil, ErrIrreversibleRequiresSupervisor
		}
		supID = *supervisorID
		sig = *titanM3Signature
	} else if titanM3Signature != nil {
		sig = *titanM3Signature
	}

	receipt := VerificationReceipt{
		EventID:           eventID,
		ModelVersion:      modelVersion,
		DatasetHash:       datasetHash,
		SNOMEDCode:        snomedCode,
		FHIRResourceValid: true,
		OperatorID:        operatorID,
		SupervisorID:      supID,
		TitanM3Signature:  sig,
		TimestampUTC:      time.Now().UTC(),
	}
	receipt.ComputeHash()

	if err := g.forensicSink.AppendForensicRecord(receipt, "APPROVED_EXECUTION", nil); err != nil {
		return nil, err
	}

	return &receipt, nil
}
