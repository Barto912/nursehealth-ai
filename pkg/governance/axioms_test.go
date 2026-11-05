package governance

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"testing"
	"time"
)

func signHelper(t *testing.T, priv *ecdsa.PrivateKey, cmd *SignedClinicalCommand) {
	t.Helper()
	digest := sha256.Sum256(cmd.CanonicalSigningBytes())
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}
	der, err := asn1.Marshal(ecdsaDERSignature{R: r, S: s})
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	cmd.SignatureDER = der
}

func TestAX012_ECDSAVerification(t *testing.T) {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	now := time.Now().UTC()
	guard := NewIrreversibleDecisionGuard(5 * time.Minute)

	cmd := SignedClinicalCommand{
		EventID:       "550e8400-e29b-41d4-a716-446655440000",
		AxiomCode:     AxiomDisarmAI,
		ClinicianID:   "MAT-AR-99812",
		PublicKeyID:   "TITAN-M3-01",
		NonceHex:      "aabbccddeeff00112233445566778899",
		SignedAtUTC:   now,
		PayloadCipher: []byte("test"),
		PrevHashHex:   "0000000000000000000000000000000000000000000000000000000000000000",
	}
	signHelper(t, priv, &cmd)

	hash, err := guard.VerifyCommand(cmd, &priv.PublicKey, now)
	if err != nil || len(hash) != 64 {
		t.Fatalf("AX-012 failed: err=%v hash=%s", err, hash)
	}

	if _, err := guard.VerifyCommand(cmd, &priv.PublicKey, now); err != ErrReplayOrClockSkew {
		t.Fatalf("replay not blocked: %v", err)
	}
}

func TestAX008_FailClose(t *testing.T) {
	normal := VitalSignsInput{16, 98, 120, 72, 36.6}
	if err := EvaluateDeterministicTriage(normal); err != nil {
		t.Fatalf("normal vitals rejected: %v", err)
	}
	hypoxia := VitalSignsInput{28, 89, 110, 115, 38.5}
	if err := EvaluateDeterministicTriage(hypoxia); err != ErrFailCloseTriggered {
		t.Fatalf("hypoxia not caught: %v", err)
	}
}

func TestAX001_AES256GCM(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	plain := []byte("DATOS_SENSIBLES_PACIENTE")
	cipher, err := SealPatientPayload(key, plain, []byte("AX-001"))
	if err != nil {
		t.Fatalf("seal failed: %v", err)
	}
	dec, err := OpenPatientPayload(key, cipher, []byte("AX-001"))
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	if string(dec) != string(plain) {
		t.Fatalf("decryption mismatch")
	}
}
