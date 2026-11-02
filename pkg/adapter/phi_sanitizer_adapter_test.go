package adapter

import (
	"strings"
	"testing"
)

func TestPHISanitizerAdapter(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		wantOK         bool
		wantContains   string
		wantNotContain string
	}{
		{
			name:         "payload vacio retorna false",
			payload:      "",
			wantOK:       false,
			wantContains: "",
		},
		{
			name:           "payload con DNI es sanitizado",
			payload:        "Paciente con DNI 12345678 presenta fiebre",
			wantOK:         true,
			wantContains:   "[REDACTED_DNI]",
			wantNotContain: "12345678",
		},
		{
			name:           "payload con email es sanitizado",
			payload:        "Contacto: paciente@email.com",
			wantOK:         true,
			wantContains:   "[REDACTED_EMAIL]",
			wantNotContain: "paciente@email.com",
		},
		{
			name:         "payload sin PHI retorna texto original",
			payload:      "Presion arterial 120/80 mmHg",
			wantOK:       true,
			wantContains: "120/80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy := NewDataProxy()
			adapter := NewPHISanitizerAdapter(proxy)

			sanitized, ok := adapter.SanitizePHI(tt.payload)

			if ok != tt.wantOK {
				t.Fatalf("esperaba ok=%v, got %v", tt.wantOK, ok)
			}

			if tt.wantContains != "" && !strings.Contains(sanitized, tt.wantContains) {
				t.Fatalf("esperaba que contenga %q, got %q", tt.wantContains, sanitized)
			}

			if tt.wantNotContain != "" && strings.Contains(sanitized, tt.wantNotContain) {
				t.Fatalf("no deberia contener %q, got %q", tt.wantNotContain, sanitized)
			}
		})
	}
}

// TestPHISanitizerAdapterSatisfiesInterface verifica que el adapter satisface
// la interfaz implicita requerida por governance.PHISanitizer.
func TestPHISanitizerAdapterSatisfiesInterface(t *testing.T) {
	proxy := NewDataProxy()
	adapter := NewPHISanitizerAdapter(proxy)

	type expectedInterface interface {
		SanitizePHI(string) (string, bool)
	}

	var _ expectedInterface = adapter
}
