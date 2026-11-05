package governance

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
	"testing"
	"time"
)

func TestAX012_ECDSAVerification(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	now := time.Now().UTC()
	guard := NewIrreversibleDecisionGuard(5 * time.Minute)

	cmd := SignedClinicalCommand{
		EventID:       "550e8400-e29b-41d4-a716-446655440000",
		AxiomCode:     AxiomDisarmAI,
		ClinicianID:   "MAT-AR-99812",
		PublicKeyID:   "titan-m3-key-01",
		NonceHex:      "abcdef0123456789",
		SignedAtUTC:   now,
		PrevHashHex:   "0000000000000000000000000000000000000000000000000000000000000000",
		PayloadCipher: []byte("test-payload-secure"),
	}

	digest := sha256.Sum256(cmd.CanonicalSigningBytes())
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// Normalización Low-S requerida por el estándar DELTA-02
	curveOrder := elliptic.P256().Params().N
	halfOrder := new(big.Int).Rsh(curveOrder, 1)
	if s.Cmp(halfOrder) > 0 {
		s.Sub(curveOrder, s)
	}

	sigDER, err := asn1.Marshal(ecdsaDERSignature{R: r, S: s})
	if err != nil {
		t.Fatalf("failed to marshal sig: %v", err)
	}
	cmd.SignatureDER = sigDER

	hash, err := guard.VerifyCommand(cmd, &priv.PublicKey, now)
	if err != nil {
		t.Fatalf("AX-012 failed: err=%v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty record hash")
	}
}

func TestAX008_FailClose(t *testing.T) {
	v := VitalSignsInput{
		RespiratoryRateBPM: 30, // Fuera de rango (8-24)
		SpO2Percent:        95,
		SystolicBP:         120,
		HeartRateBPM:       80,
		TemperatureCelsius: 37.0,
	}
	err := EvaluateDeterministicTriage(v)
	if err != ErrFailCloseTriggered {
		t.Fatalf("expected ErrFailCloseTriggered, got %v", err)
	}
}

func TestAX001_AES256GCM(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	plaintext := []byte("datos clinicos confidenciales")
	aad := []byte("contexto")

	sealed, err := SealPatientPayload(key, plaintext, aad)
	if err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	opened, err := OpenPatientPayload(key, sealed, aad)
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	if string(opened) != string(plaintext) {
		t.Fatalf("decrypted text mismatch")
	}
}
