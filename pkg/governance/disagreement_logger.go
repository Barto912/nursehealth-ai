package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Errores deterministas del DisagreementLogger (AX-016).
var (
	ErrEmptyInference  = errors.New("AX-016 violation: MedGemma 2B INT8 inference payload cannot be empty")
	ErrEmptyRefutation = errors.New("AX-016 violation: human refutation must be substantiated (non-empty)")
	ErrPHINotSanitized = errors.New("AX-016 violation: PHI sanitization by Data Proxy not confirmed; ATDP tuple rejected")
)

// ATDPTuple representa la tupla de discrepancia anonimizada para RLHF clínico (AX-016):
// [Contexto Clínico Anonimizado SNOMED CT + Inferencia MedGemma 2B INT8 + Refutación Fundamentada del Operador].
type ATDPTuple struct {
	EventID              string    `json:"event_id"`
	SNOMEDCode           string    `json:"snomed_code"`
	MedGemmaInference    string    `json:"medgemma_inference"`
	HumanRefutation      string    `json:"human_refutation"`
	OperatorID           string    `json:"operator_id"`
	SanitizedByDataProxy bool      `json:"sanitized_by_data_proxy"`
	TimestampUTC         time.Time `json:"timestamp_utc"`
	TupleSHA256          string    `json:"tuple_sha256"`
}

// ComputeHash genera el digest SHA-256 determinista de la tupla ATDP.
func (t *ATDPTuple) ComputeHash() string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%t|%s",
		t.EventID,
		t.SNOMEDCode,
		t.MedGemmaInference,
		t.HumanRefutation,
		t.OperatorID,
		t.SanitizedByDataProxy,
		t.TimestampUTC.UTC().Format(time.RFC3339Nano),
	)
	sum := sha256.Sum256([]byte(payload))
	t.TupleSHA256 = hex.EncodeToString(sum[:])
	return t.TupleSHA256
}

// PHISanitizer abstrae la sanitización local de PHI del Data Proxy (pkg/adapter/data_proxy.go).
type PHISanitizer interface {
	SanitizePHI(payload string) (sanitized string, ok bool)
}

// DisagreementLogger implementa el mecanismo obligatorio de refutación humana (AX-016).
type DisagreementLogger struct {
	forensicSink ForensicSinkRepository
	sanitizer    PHISanitizer
}

// NewDisagreementLogger construye el logger con sus dependencias.
func NewDisagreementLogger(sink ForensicSinkRepository, san PHISanitizer) *DisagreementLogger {
	return &DisagreementLogger{
		forensicSink: sink,
		sanitizer:    san,
	}
}

// LogRefutation captura una refutación humana a una inferencia de MedGemma 2B INT8,
// sanitiza PHI localmente y empaqueta la tupla ATDP en el registro forense inmutable.
func (l *DisagreementLogger) LogRefutation(
	eventID string,
	snomedCode string,
	medGemmaInference string,
	humanRefutation string,
	operatorID string,
) (*ATDPTuple, error) {

	// 1. AX-016: la inferencia de MedGemma no puede estar vacía
	if medGemmaInference == "" {
		return nil, ErrEmptyInference
	}

	// 2. AX-016: la refutación humana debe estar fundamentada (no vacía)
	if humanRefutation == "" {
		return nil, ErrEmptyRefutation
	}

	// 3. Sanitización PHI local vía Data Proxy (soberanía sanitaria)
	sanitizedInference, okInf := l.sanitizer.SanitizePHI(medGemmaInference)
	sanitizedRefutation, okRef := l.sanitizer.SanitizePHI(humanRefutation)
	if !okInf || !okRef {
		return nil, ErrPHINotSanitized
	}

	tuple := ATDPTuple{
		EventID:              eventID,
		SNOMEDCode:           snomedCode,
		MedGemmaInference:    sanitizedInference,
		HumanRefutation:      sanitizedRefutation,
		OperatorID:           operatorID,
		SanitizedByDataProxy: true,
		TimestampUTC:         time.Now().UTC(),
	}
	tuple.ComputeHash()

	// 4. Convergencia en el único registro forense inmutable (AX-011)
	payloadJSON, err := json.Marshal(tuple)
	if err != nil {
		return nil, err
	}
	payloadStr := string(payloadJSON)

	receipt := VerificationReceipt{
		EventID:           eventID,
		ModelVersion:      "MedGemma-2B-INT8",
		DatasetHash:       tuple.TupleSHA256,
		SNOMEDCode:        snomedCode,
		FHIRResourceValid: true,
		OperatorID:        operatorID,
		TimestampUTC:      tuple.TimestampUTC,
	}
	receipt.ComputeHash()

	if err := l.forensicSink.AppendForensicRecord(receipt, "HUMAN_REFUTATION", &payloadStr); err != nil {
		return nil, err
	}

	return &tuple, nil
}
