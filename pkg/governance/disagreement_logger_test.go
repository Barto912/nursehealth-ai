package governance

import (
	"errors"
	"testing"
)

// mockLoggerSink captura receipt, status y payload ATDP para tests del logger.
type mockLoggerSink struct {
	appendErr   error
	lastStatus  string
	lastPayload *string
	lastReceipt *VerificationReceipt
}

func (m *mockLoggerSink) AppendForensicRecord(receipt VerificationReceipt, status string, atdpDisagreementPayload *string) error {
	if m.appendErr != nil {
		return m.appendErr
	}
	r := receipt
	m.lastReceipt = &r
	m.lastStatus = status
	m.lastPayload = atdpDisagreementPayload
	return nil
}

// mockSanitizer simula la sanitización PHI del Data Proxy.
type mockSanitizer struct {
	ok bool
}

func (m *mockSanitizer) SanitizePHI(payload string) (string, bool) {
	if !m.ok {
		return "", false
	}
	return "[SANITIZED] " + payload, true
}

func TestDisagreementLogger(t *testing.T) {
	tests := []struct {
		name        string
		inference   string
		refutation  string
		sanitizerOK bool
		sinkErr     error
		wantErr     error
		wantTuple   bool
	}{
		{
			name:        "AX-016: inferencia vacia es bloqueada",
			inference:   "",
			refutation:  "El paciente no presenta fiebre",
			sanitizerOK: true,
			wantErr:     ErrEmptyInference,
		},
		{
			name:        "AX-016: refutacion vacia es bloqueada",
			inference:   "Sugiero administrar antibiotico",
			refutation:  "",
			sanitizerOK: true,
			wantErr:     ErrEmptyRefutation,
		},
		{
			name:        "AX-016: sanitizacion PHI no confirmada es bloqueada",
			inference:   "Sugiero administrar antibiotico",
			refutation:  "El paciente no presenta fiebre",
			sanitizerOK: false,
			wantErr:     ErrPHINotSanitized,
		},
		{
			name:        "Happy path: refutacion valida genera tupla ATDP y asiento forense",
			inference:   "Sugiero administrar antibiotico",
			refutation:  "El paciente no presenta fiebre; descartar sepsis",
			sanitizerOK: true,
			wantErr:     nil,
			wantTuple:   true,
		},
		{
			name:        "Error del sink forense se propaga",
			inference:   "Sugiero administrar antibiotico",
			refutation:  "El paciente no presenta fiebre",
			sanitizerOK: true,
			sinkErr:     errors.New("db connection lost"),
			wantErr:     errors.New("db connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := &mockLoggerSink{appendErr: tt.sinkErr}
			san := &mockSanitizer{ok: tt.sanitizerOK}
			logger := NewDisagreementLogger(sink, san)

			tuple, err := logger.LogRefutation(
				"EVT-100",
				"123456",
				tt.inference,
				tt.refutation,
				"NURSE-001",
			)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("esperaba error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Fatalf("error inesperado: got %v, want %v", err, tt.wantErr)
				}
				if tuple != nil {
					t.Fatalf("esperaba tuple nil ante error, got %+v", tuple)
				}
				return
			}

			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !tt.wantTuple {
				t.Fatalf("esperaba sin tuple")
			}
			if tuple == nil {
				t.Fatalf("esperaba tuple valido, got nil")
			}
			if tuple.TupleSHA256 == "" {
				t.Fatalf("esperaba TupleSHA256 no vacio")
			}
			if !tuple.SanitizedByDataProxy {
				t.Fatalf("esperaba SanitizedByDataProxy true")
			}
			if sink.lastReceipt == nil {
				t.Fatalf("el sink forense no fue invocado")
			}
			if sink.lastStatus != "HUMAN_REFUTATION" {
				t.Fatalf("status incorrecto: got %s", sink.lastStatus)
			}
			if sink.lastPayload == nil || *sink.lastPayload == "" {
				t.Fatalf("esperaba payload ATDP no vacio")
			}
		})
	}
}

func TestATDPTupleComputeHashDeterministic(t *testing.T) {
	tuple := ATDPTuple{
		EventID:              "EVT-100",
		SNOMEDCode:           "123456",
		MedGemmaInference:    "inf",
		HumanRefutation:      "ref",
		OperatorID:           "NURSE-001",
		SanitizedByDataProxy: true,
	}
	h1 := tuple.ComputeHash()
	h2 := tuple.ComputeHash()
	if h1 != h2 {
		t.Fatalf("hash no determinista: %s != %s", h1, h2)
	}
	if len(h1) != 64 {
		t.Fatalf("hash SHA-256 debe tener 64 caracteres hex, got %d", len(h1))
	}
}
