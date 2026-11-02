package governance

import (
	"errors"
	"testing"
)

// mockForensicSink implementa ForensicSinkRepository para tests.
type mockForensicSink struct {
	appendErr error
	lastCall  *VerificationReceipt
}

func (m *mockForensicSink) AppendForensicRecord(receipt VerificationReceipt, status string, atdpDisagreementPayload *string) error {
	if m.appendErr != nil {
		return m.appendErr
	}
	r := receipt
	m.lastCall = &r
	return nil
}

// mockTaxonomy implementa TaxonomyValidator para tests.
type mockTaxonomy struct {
	valid bool
}

func (m *mockTaxonomy) IsValidFHIRAndSNOMED(snomedCode string, fhirPayload []byte) bool {
	return m.valid
}

func strPtr(s string) *string { return &s }

func TestIrreversibleDecisionGuard(t *testing.T) {
	tests := []struct {
		name             string
		ruleStatus       EthicsReviewStatus
		reversibility    ProcedureReversibility
		authorizingAgent AgentType
		modelVersion     string
		snomedCode       string
		operatorID       string
		supervisorID     *string
		titanSignature   *string
		taxonomyValid    bool
		sinkErr          error
		wantErr          error
		wantReceipt      bool
	}{
		{
			name:             "AX-009: regla sin aprobacion etica es bloqueada",
			ruleStatus:       "DRAFT",
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    true,
			wantErr:          ErrRuleNotApprovedByCommittee,
		},
		{
			name:             "AX-008: falta modelVersion es bloqueada",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    true,
			wantErr:          ErrMissingResponsibilityChain,
		},
		{
			name:             "AX-008: falta operatorID es bloqueada",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "",
			taxonomyValid:    true,
			wantErr:          ErrMissingResponsibilityChain,
		},
		{
			name:             "Fail-Close: taxonomia invalida es bloqueada",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    false,
			wantErr:          ErrInvalidFHIRorSNOMEDTaxonomy,
		},
		{
			name:             "AX-012: Device (IA) como autorizador es bloqueado",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentDevice,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    true,
			wantErr:          ErrAIAuthorizationForbidden,
		},
		{
			name:             "AX-012: irreversible sin supervisor es bloqueado",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureIrreversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			supervisorID:     nil,
			titanSignature:   strPtr("SIG-123"),
			taxonomyValid:    true,
			wantErr:          ErrIrreversibleRequiresSupervisor,
		},
		{
			name:             "AX-012: irreversible sin firma Titan M3 es bloqueado",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureIrreversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			supervisorID:     strPtr("DR-001"),
			titanSignature:   nil,
			taxonomyValid:    true,
			wantErr:          ErrIrreversibleRequiresSupervisor,
		},
		{
			name:             "Happy path: reversible aprobado genera receipt valido",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    true,
			wantErr:          nil,
			wantReceipt:      true,
		},
		{
			name:             "Happy path: irreversible con supervisor y firma genera receipt valido",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureIrreversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			supervisorID:     strPtr("DR-001"),
			titanSignature:   strPtr("SIG-TITAN-ABC"),
			taxonomyValid:    true,
			wantErr:          nil,
			wantReceipt:      true,
		},
		{
			name:             "Error del sink forense se propaga",
			ruleStatus:       StatusApprovedByCommittee,
			reversibility:    ProcedureReversible,
			authorizingAgent: AgentPractitioner,
			modelVersion:     "MedGemma-2B-INT8",
			snomedCode:       "123456",
			operatorID:       "NURSE-001",
			taxonomyValid:    true,
			sinkErr:          errors.New("db connection lost"),
			wantErr:          errors.New("db connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := &mockForensicSink{appendErr: tt.sinkErr}
			tax := &mockTaxonomy{valid: tt.taxonomyValid}
			guard := NewIrreversibleDecisionGuard(sink, tax)

			receipt, err := guard.ValidateAndSeal(
				"EVT-001",
				tt.ruleStatus,
				tt.reversibility,
				tt.authorizingAgent,
				tt.modelVersion,
				"HASH-DATASET-001",
				tt.snomedCode,
				[]byte(`{"resourceType":"CareEvent"}`),
				tt.operatorID,
				tt.supervisorID,
				tt.titanSignature,
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("esperaba error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("error inesperado: got %v, want %v", err, tt.wantErr)
				}
				if receipt != nil {
					t.Fatalf("esperaba receipt nil ante error, got %+v", receipt)
				}
				return
			}

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !tt.wantReceipt {
				t.Fatalf("esperaba sin receipt")
			}
			if receipt == nil {
				t.Fatalf("esperaba receipt valido, got nil")
			}
			if receipt.ReceiptSHA256 == "" {
				t.Fatalf("esperaba ReceiptSHA256 no vacio")
			}
			if receipt.EventID != "EVT-001" {
				t.Fatalf("EventID incorrecto: got %s", receipt.EventID)
			}
			if sink.lastCall == nil {
				t.Fatalf("el sink forense no fue invocado")
			}
		})
	}
}

func TestVerificationReceiptComputeHashDeterministic(t *testing.T) {
	r := VerificationReceipt{
		EventID:           "EVT-001",
		ModelVersion:      "MedGemma-2B-INT8",
		DatasetHash:       "HASH-001",
		SNOMEDCode:        "123456",
		FHIRResourceValid: true,
		OperatorID:        "NURSE-001",
		TitanM3Signature:  "SIG-001",
	}
	h1 := r.ComputeHash()
	h2 := r.ComputeHash()
	if h1 != h2 {
		t.Fatalf("hash no determinista: %s != %s", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("hash SHA-256 debe tener 64 caracteres hex, got %d", len(h1))
	}
}
