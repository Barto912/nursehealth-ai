package governance

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"
)

const (
	AxiomPatientDignity       = "AX-001"
	AxiomNoAlgorithmicHarm    = "AX-008"
	AxiomForensicTransparency = "AX-009"
	AxiomDisarmAI             = "AX-012"
	AxiomHumanRefutation      = "AX-016"
)

var (
	ErrUnauthorizedAxiom        = errors.New("GOVERNANCE_FAIL_CLOSE: axiom not part of the 5 cryptographically verifiable invariants")
	ErrMissingHardwareSignature = errors.New("AX-012_VIOLATION: irreversible action requires hardware ECDSA P-256 signature")
	ErrInvalidHardwareSignature = errors.New("AX-012_VIOLATION: ECDSA P-256 cryptographic signature verification failed")
	ErrReplayOrClockSkew        = errors.New("AX-012_VIOLATION: signature timestamp outside allowed window or nonce replayed")
	ErrEmptyDisagreementNote    = errors.New("AX-016_VIOLATION: human refutation requires non-empty clinical justification")
	ErrFailCloseTriggered       = errors.New("AX-008_VIOLATION: deterministic clinical boundary exceeded — mechanical fail-close engaged")
	ErrEncryptionIntegrity      = errors.New("AX-001_VIOLATION: invalid AES-256 key length or ciphertext authentication failure")
)

type ecdsaDERSignature struct { R, S *big.Int }

type SignedClinicalCommand struct {
	EventID, AxiomCode, ClinicianID, PublicKeyID, NonceHex, PrevHashHex string
	SignedAtUTC      time.Time
	PayloadCipher    []byte
	DisagreementNote string
	SignatureDER     []byte
}

func (c SignedClinicalCommand) CanonicalSigningBytes() []byte {
	payloadHex := hex.EncodeToString(c.PayloadCipher)
	return []byte(fmt.Sprintf("NH_DC_V1|%d:%s|%d:%s|%d:%s|%d:%s|%d:%s|%d:%d|%d:%s|%d:%s|%d:%s",
		len(c.EventID), c.EventID, len(c.AxiomCode), c.AxiomCode, len(c.ClinicianID), c.ClinicianID,
		len(c.PublicKeyID), c.PublicKeyID, len(c.NonceHex), c.NonceHex,
		len(fmt.Sprintf("%d", c.SignedAtUTC.UTC().Unix())), c.SignedAtUTC.UTC().Unix(),
		len(payloadHex), payloadHex, len(c.DisagreementNote), c.DisagreementNote, len(c.PrevHashHex), c.PrevHashHex))
}

func (c SignedClinicalCommand) ComputeRecordHash() string {
	sigHex := hex.EncodeToString(c.SignatureDER)
	return hex.EncodeToString(sha256.Sum256([]byte(fmt.Sprintf("%s|%d:%s", string(c.CanonicalSigningBytes()), len(sigHex), sigHex))))
}

type IrreversibleDecisionGuard struct {
	mu           sync.Mutex
	maxClockSkew time.Duration
	seenNonces   map[string]time.Time
}

func NewIrreversibleDecisionGuard(maxClockSkew time.Duration) *IrreversibleDecisionGuard {
	return &IrreversibleDecisionGuard{maxClockSkew: maxClockSkew, seenNonces: make(map[string]time.Time)}
}

func (g *IrreversibleDecisionGuard) VerifyCommand(cmd SignedClinicalCommand, pubKey *ecdsa.PublicKey, now time.Time) (string, error) {
	switch cmd.AxiomCode {
	case AxiomPatientDignity, AxiomNoAlgorithmicHarm, AxiomForensicTransparency, AxiomDisarmAI, AxiomHumanRefutation:
	default:
		return "", ErrUnauthorizedAxiom
	}
	if pubKey == nil || len(cmd.SignatureDER) == 0 || cmd.ClinicianID == "" || cmd.NonceHex == "" {
		return "", ErrMissingHardwareSignature
	}
	if cmd.AxiomCode == AxiomHumanRefutation && len(cmd.DisagreementNote) < 10 {
		return "", ErrEmptyDisagreementNote
	}
	skew := now.UTC().Sub(cmd.SignedAtUTC.UTC())
	if skew < -time.Minute || skew > g.maxClockSkew {
		return "", ErrReplayOrClockSkew
	}
	digest := sha256.Sum256(cmd.CanonicalSigningBytes())
	var sig ecdsaDERSignature
	if _, err := asn1.Unmarshal(cmd.SignatureDER, &sig); err != nil {
		return "", ErrInvalidHardwareSignature
	}
	if !ecdsa.Verify(pubKey, digest[:], sig.R, sig.S) {
		return "", ErrInvalidHardwareSignature
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.seenNonces[cmd.NonceHex]; exists {
		return "", ErrReplayOrClockSkew
	}
	g.seenNonces[cmd.NonceHex] = now.UTC()
	return cmd.ComputeRecordHash(), nil
}

func SealPatientPayload(key32 []byte, plaintext []byte, associatedData []byte) ([]byte, error) {
	if len(key32) != 32 { return nil, ErrEncryptionIntegrity }
	block, err := aes.NewCipher(key32); if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block); if err != nil { return nil, err }
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return nil, err }
	return gcm.Seal(nonce, nonce, plaintext, associatedData), nil
}

type VitalSignsInput struct { RespiratoryRateBPM, SpO2Percent, SystolicBP, HeartRateBPM int; TemperatureCelsius float64 }

func EvaluateDeterministicTriage(v VitalSignsInput) error {
	if v.SpO2Percent < 92 || v.SpO2Percent > 100 || v.RespiratoryRateBPM < 8 || v.RespiratoryRateBPM > 24 || v.SystolicBP < 90 || v.SystolicBP > 219 || v.HeartRateBPM < 40 || v.HeartRateBPM > 130 || v.TemperatureCelsius < 35.0 || v.TemperatureCelsius > 39.0 {
		return ErrFailCloseTriggered
	}
	return nil
}
