package governance

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
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
	MaxActiveNonces           = 100000
)

var (
	ErrUnauthorizedAxiom        = errors.New("GOVERNANCE_FAIL_CLOSE: axiom not part of the 5 cryptographically verifiable invariants")
	ErrMissingHardwareSignature = errors.New("AX-012_VIOLATION: irreversible action requires hardware ECDSA P-256 signature")
	ErrInvalidHardwareSignature = errors.New("AX-012_VIOLATION: ECDSA P-256 cryptographic signature verification failed")
	ErrReplayOrClockSkew        = errors.New("AX-012_VIOLATION: signature timestamp outside allowed window or nonce replayed")
	ErrEmptyDisagreementNote    = errors.New("AX-016_VIOLATION: human refutation requires non-empty clinical justification")
	ErrFailCloseTriggered       = errors.New("AX-008_VIOLATION: deterministic clinical boundary exceeded — mechanical fail-close engaged")
	ErrEncryptionIntegrity      = errors.New("AX-001_VIOLATION: invalid AES-256 key length or ciphertext authentication failure")
	ErrMalleableOrInvalidCurve  = errors.New("AX-012_VIOLATION: non-canonical ECDSA High-S signature or invalid elliptic curve")
	ErrNonceCacheExhausted      = errors.New("AX-008_VIOLATION: anti-replay nonce cache capacity reached — fail-close engaged")
)

type ecdsaDERSignature struct{ R, S *big.Int }

type SignedClinicalCommand struct {
	EventID, AxiomCode, ClinicianID, PublicKeyID, NonceHex, PrevHashHex string
	SignedAtUTC                                                         time.Time
	PayloadCipher                                                       []byte
	DisagreementNote                                                    string
	SignatureDER                                                        []byte
}

func (c SignedClinicalCommand) CanonicalSigningBytes() []byte {
	payloadHex := hex.EncodeToString(c.PayloadCipher)
	return fmt.Appendf(nil, "NH_DC_V1|%d:%s|%d:%s|%d:%s|%d:%s|%d:%s|%d:%d|%d:%s|%d:%s|%d:%s",
		len(c.EventID), c.EventID, len(c.AxiomCode), c.AxiomCode, len(c.ClinicianID), c.ClinicianID,
		len(c.PublicKeyID), c.PublicKeyID, len(c.NonceHex), c.NonceHex,
		len(fmt.Sprintf("%d", c.SignedAtUTC.UTC().Unix())), c.SignedAtUTC.UTC().Unix(),
		len(payloadHex), payloadHex, len(c.DisagreementNote), c.DisagreementNote, len(c.PrevHashHex), c.PrevHashHex)
}

func (c SignedClinicalCommand) ComputeRecordHash() string {
	sigHex := hex.EncodeToString(c.SignatureDER)
	canonicalBytes := c.CanonicalSigningBytes()
	sumBytes := fmt.Appendf(nil, "%s|%d:%s", string(canonicalBytes), len(sigHex), sigHex)
	sum := sha256.Sum256(sumBytes)
	return hex.EncodeToString(sum[:])
}

type IrreversibleDecisionGuard struct {
	mu           sync.Mutex
	maxClockSkew time.Duration
	seenNonces   map[string]time.Time
}

func NewIrreversibleDecisionGuard(maxClockSkew time.Duration) *IrreversibleDecisionGuard {
	return &IrreversibleDecisionGuard{maxClockSkew: maxClockSkew, seenNonces: make(map[string]time.Time)}
}

// VerifyCommand exige curva P-256, firma ECDSA Low-S canónica, ventana temporal estricta y purga de nonces.
func (g *IrreversibleDecisionGuard) VerifyCommand(cmd SignedClinicalCommand, pubKey *ecdsa.PublicKey, now time.Time) (string, error) {
	switch cmd.AxiomCode {
	case AxiomPatientDignity, AxiomNoAlgorithmicHarm, AxiomForensicTransparency, AxiomDisarmAI, AxiomHumanRefutation:
	default:
		return "", ErrUnauthorizedAxiom
	}

	if pubKey == nil || len(cmd.SignatureDER) == 0 || cmd.ClinicianID == "" || cmd.NonceHex == "" {
		return "", ErrMissingHardwareSignature
	}

	// Garantizar que la clave pública esté estrictamente sobre NIST P-256
	if pubKey.Curve != elliptic.P256() {
		return "", ErrMalleableOrInvalidCurve
	}

	if cmd.AxiomCode == AxiomHumanRefutation && len(cmd.DisagreementNote) < 10 {
		return "", ErrEmptyDisagreementNote
	}

	// 1. Validación Anti-Replay temporal
	nowUTC := now.UTC()
	skew := nowUTC.Sub(cmd.SignedAtUTC.UTC())
	if skew < -time.Minute || skew > g.maxClockSkew {
		return "", ErrReplayOrClockSkew
	}

	// 2. Deserialización y bloqueo de maleabilidad ECDSA (Exigir Low-S: 0 < S <= N/2)
	var sig ecdsaDERSignature
	rest, err := asn1.Unmarshal(cmd.SignatureDER, &sig)
	if err != nil || len(rest) > 0 || sig.R == nil || sig.S == nil {
		return "", ErrInvalidHardwareSignature
	}

	curveOrder := elliptic.P256().Params().N
	halfOrder := new(big.Int).Rsh(curveOrder, 1) // N / 2
	if sig.R.Sign() <= 0 || sig.R.Cmp(curveOrder) >= 0 || sig.S.Sign() <= 0 || sig.S.Cmp(halfOrder) > 0 {
		return "", ErrMalleableOrInvalidCurve
	}

	// 3. Verificación criptográfica ECDSA P-256 sobre el digest SHA-256 canónico
	digest := sha256.Sum256(cmd.CanonicalSigningBytes())
	if !ecdsa.Verify(pubKey, digest[:], sig.R, sig.S) {
		return "", ErrInvalidHardwareSignature
	}

	// 4. Gestión de memoria acotada y purga de Nonces expirados (Anti-OOM DoS)
	g.mu.Lock()
	defer g.mu.Unlock()

	for k, ts := range g.seenNonces {
		if nowUTC.Sub(ts) > g.maxClockSkew {
			delete(g.seenNonces, k)
		}
	}

	if _, exists := g.seenNonces[cmd.NonceHex]; exists {
		return "", ErrReplayOrClockSkew
	}
	if len(g.seenNonces) >= MaxActiveNonces {
		return "", ErrNonceCacheExhausted
	}
	g.seenNonces[cmd.NonceHex] = nowUTC

	return cmd.ComputeRecordHash(), nil
}

func SealPatientPayload(key32 []byte, plaintext []byte, associatedData []byte) ([]byte, error) {
	if len(key32) != 32 {
		return nil, ErrEncryptionIntegrity
	}
	block, err := aes.NewCipher(key32)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, associatedData), nil
}

type VitalSignsInput struct {
	RespiratoryRateBPM, SpO2Percent, SystolicBP, HeartRateBPM int
	TemperatureCelsius                                        float64
}

func EvaluateDeterministicTriage(v VitalSignsInput) error {
	if v.SpO2Percent < 92 || v.SpO2Percent > 100 || v.RespiratoryRateBPM < 8 || v.RespiratoryRateBPM > 24 || v.SystolicBP < 90 || v.SystolicBP > 219 || v.HeartRateBPM < 40 || v.HeartRateBPM > 130 || v.TemperatureCelsius < 35.0 || v.TemperatureCelsius > 39.0 {
		return ErrFailCloseTriggered
	}
	return nil
}

func OpenPatientPayload(key32 []byte, ciphertext []byte, associatedData []byte) ([]byte, error) {
	if len(key32) != 32 {
		return nil, ErrEncryptionIntegrity
	}
	block, err := aes.NewCipher(key32)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrEncryptionIntegrity
	}
	nonce := ciphertext[:gcm.NonceSize()]
	enc := ciphertext[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, enc, associatedData)
	if err != nil {
		return nil, ErrEncryptionIntegrity
	}
	return plain, nil
}
