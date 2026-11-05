package adapter

import "testing"

func TestPseudonymization(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // exactamente 32 bytes
	p, err := NewDataProxy(key, AuthorizedJurisdiction)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	pseudo := p.PseudonymizeIdentifier("20-30456789-4")
	if pseudo == "" || pseudo == "20-30456789-4" {
		t.Fatalf("pseudonymization failed: got %s", pseudo)
	}
	pseudo2 := p.PseudonymizeIdentifier("20-30456789-4")
	if pseudo != pseudo2 {
		t.Fatalf("not deterministic: %s != %s", pseudo, pseudo2)
	}
}

func TestSanitizeClinicalText(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // exactamente 32 bytes
	p, _ := NewDataProxy(key, AuthorizedJurisdiction)
	input := "Paciente CUIL 20-30456789-4 DNI 30.456.789"
	san := p.SanitizeClinicalText(input)
	if san == input {
		t.Fatalf("sanitization failed: input unchanged")
	}
}

func TestJurisdictionFence(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // exactamente 32 bytes
	_, err := NewDataProxy(key, "US-EAST1")
	if err != ErrJurisdictionFence {
		t.Fatalf("fence failed: %v", err)
	}
}

func TestInvalidHMACKey(t *testing.T) {
	_, err := NewDataProxy([]byte("short"), AuthorizedJurisdiction)
	if err != ErrInvalidHMACSecret {
		t.Fatalf("short key not rejected")
	}
}
