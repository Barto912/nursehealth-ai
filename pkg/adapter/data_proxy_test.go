package adapter

import "testing"

func TestRedactPII(t *testing.T) {
	dp := NewDataProxy()

	tests := []struct {
		name      string
		input     string
		want      string
		wantCount int
	}{
		{
			name:      "Empty input",
			input:     "",
			want:      "",
			wantCount: 0,
		},
		{
			name:      "No PII",
			input:     "Paciente presenta fiebre de 38 grados",
			want:      "Paciente presenta fiebre de 38 grados",
			wantCount: 0,
		},
		{
			name:      "DNI without dots",
			input:     "Paciente DNI 12345678",
			want:      "Paciente DNI [REDACTED_DNI]",
			wantCount: 1,
		},
		{
			name:      "DNI with dots",
			input:     "Paciente DNI 12.345.678",
			want:      "Paciente DNI [REDACTED_DNI]",
			wantCount: 1,
		},
		{
			name:      "Email",
			input:     "Contacto: juan.perez@hospital.com",
			want:      "Contacto: [REDACTED_EMAIL]",
			wantCount: 1,
		},
		{
			name:      "Phone local continuous digits",
			input:     "Teléfono: 1145678901",
			want:      "Teléfono: [REDACTED_PHONE]",
			wantCount: 1,
		},
		{
			// Limitación conocida (hardening en PR #35): el "+" inicial
			// permanece, pero los dígitos (el PII real) quedan redactados.
			name:      "Phone with country code",
			input:     "Celular: +541145678901",
			want:      "Celular: +[REDACTED_PHONE]",
			wantCount: 1,
		},
		{
			name:      "Address",
			input:     "Domicilio en calle corrientes 1234",
			want:      "Domicilio en [REDACTED_ADDRESS]",
			wantCount: 1,
		},
		{
			name:      "Name with title",
			input:     "Atendido por Dra. García",
			want:      "Atendido por [REDACTED_NAME]",
			wantCount: 1,
		},
		{
			name:      "Mixed PII",
			input:     "DNI 12.345.678, email jp@mail.com",
			want:      "DNI [REDACTED_DNI], email [REDACTED_EMAIL]",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, count := dp.RedactPII(tt.input)
			if got != tt.want {
				t.Errorf("RedactPII() got = %q, want %q", got, tt.want)
			}
			if count != tt.wantCount {
				t.Errorf("RedactPII() count = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestIsTechnicalField(t *testing.T) {
	dp := NewDataProxy()

	technical := []string{"snomed_ct", "fhir_code", "timestamp", "vital_signs", "SNOMED_CT"}
	for _, key := range technical {
		if !dp.IsTechnicalField(key) {
			t.Errorf("IsTechnicalField(%q) = false, want true", key)
		}
	}

	nonTechnical := []string{"patient_name", "address", "notes"}
	for _, key := range nonTechnical {
		if dp.IsTechnicalField(key) {
			t.Errorf("IsTechnicalField(%q) = true, want false", key)
		}
	}
}
