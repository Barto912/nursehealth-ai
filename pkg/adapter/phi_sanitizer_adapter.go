package adapter

// PHISanitizerAdapter adapta DataProxy para satisfacer la interfaz PHISanitizer
// del paquete governance sin crear dependencia cíclica (structural typing de Go).
type PHISanitizerAdapter struct {
	proxy *DataProxy
}

// NewPHISanitizerAdapter construye el adapter con el DataProxy real.
func NewPHISanitizerAdapter(proxy *DataProxy) *PHISanitizerAdapter {
	return &PHISanitizerAdapter{proxy: proxy}
}

// SanitizePHI implementa la interfaz implícita requerida por governance.PHISanitizer.
// Retorna el texto sanitizado y true si la sanitización fue exitosa.
func (a *PHISanitizerAdapter) SanitizePHI(payload string) (string, bool) {
	if payload == "" {
		return "", false
	}
	sanitized, fieldsRedacted := a.proxy.RedactPII(payload)
	_ = fieldsRedacted // reservado para telemetría futura
	return sanitized, true
}
