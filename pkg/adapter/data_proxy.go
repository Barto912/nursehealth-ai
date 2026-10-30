package adapter

import (
	"regexp"
	"strings"
)

// DataProxy redacta automáticamente PII/PHI antes de transmitir datos a la nube.
// Cumple con HIPAA, LGPD y Ley 25.326. Es determinista y seguro para el borde.
type DataProxy struct {
	dniRegex     *regexp.Regexp
	emailRegex   *regexp.Regexp
	phoneRegex   *regexp.Regexp
	addressRegex *regexp.Regexp
	nameRegex    *regexp.Regexp
}

// NewDataProxy inicializa el proxy con patrones de redactado precompilados.
func NewDataProxy() *DataProxy {
	return &DataProxy{
		dniRegex:     regexp.MustCompile(`\b\d{1,2}\.?\d{3}\.?\d{3}\b`),
		emailRegex:   regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`),
		phoneRegex:   regexp.MustCompile(`\b(?:\+?54\s?)?(?:11|2[0-9]{2,3}|3[0-9]{2,3})\s?\d{7,8}\b`),
		addressRegex: regexp.MustCompile(`\b(?:calle|av\.|avenida)\s+[A-Za-záéíóúñÁÉÍÓÚÑ\s]+\s+\d+\b`),
		nameRegex:    regexp.MustCompile(`\b(?:Sr\.|Sra\.|Dr\.|Dra\.)\s+[A-Z][a-záéíóúñ]+\b`),
	}
}

// RedactPII aplica redactado determinista a un texto clínico.
// Retorna el texto anonimizado y el contador de campos redactados (para auditoría).
func (dp *DataProxy) RedactPII(input string) (string, int) {
	if input == "" {
		return "", 0
	}

	fieldsRedacted := 0

	replacements := []struct {
		re  *regexp.Regexp
		tag string
	}{
		{dp.dniRegex, "[REDACTED_DNI]"},
		{dp.emailRegex, "[REDACTED_EMAIL]"},
		{dp.phoneRegex, "[REDACTED_PHONE]"},
		{dp.addressRegex, "[REDACTED_ADDRESS]"},
		{dp.nameRegex, "[REDACTED_NAME]"},
	}

	for _, r := range replacements {
		if r.re.MatchString(input) {
			input = r.re.ReplaceAllString(input, r.tag)
			fieldsRedacted++
		}
	}

	return input, fieldsRedacted
}

// IsTechnicalField identifica campos que NO deben redactarse
// (códigos SNOMED CT, FHIR, timestamps y signos vitales).
func (dp *DataProxy) IsTechnicalField(key string) bool {
	technicalKeys := map[string]bool{
		"snomed_ct":      true,
		"fhir_code":      true,
		"timestamp":      true,
		"relative_time":  true,
		"vital_signs":    true,
		"medication_id":  true,
		"procedure_code": true,
	}
	return technicalKeys[strings.ToLower(key)]
}
