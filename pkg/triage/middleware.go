package triage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// AuditRepository defines the contract for persisting cognitive audit events.
// Implementations must be append-only and thread-safe.
type AuditRepository interface {
	Append(ctx context.Context, event CognitiveEvent) error
}

// OperatorGate abstracts the operator-side cognitive gate (RuidoCheck,
// pkg/cognitive). It returns the typed block reason, or ReasonNone when the
// gate passes. Keeping it as an interface preserves bounded-context isolation:
// the triage package never imports the cognitive package directly.
type OperatorGate interface {
	Validate() CognitiveBlockReason
}

// BlockError wraps a typed block reason so HTTP layers can render 403 payloads.
type BlockError struct {
	Reason CognitiveBlockReason
	Cause  error
}

// Error implements the error interface.
func (e *BlockError) Error() string {
	return "triage: care event blocked: " + e.Reason.String()
}

// Unwrap exposes the underlying validation error, if any.
func (e *BlockError) Unwrap() error { return e.Cause }

// CareEventInput carries the minimal data required to validate a care event.
type CareEventInput struct {
	ID         string
	PatientID  string
	OperatorID string
	Operator   OperatorGate

	SufferingIndicator  bool
	ActiveListeningNote string
	Landscape           ContextualLandscape
	EthicalConfirmation bool
	FailedAttempts      int
	SupervisorConsulted bool
}

// ValidateCareEvent executes the cognitive firewall gates in order:
// operator-side first, then patient-side. On any block it persists a
// CareEventBlockedByCognitiveState event and returns a *BlockError.
// On success it persists a CognitiveStateValidated event and returns nil.
func ValidateCareEvent(ctx context.Context, repo AuditRepository, input CareEventInput) error {
	if input.Operator != nil {
		if reason := input.Operator.Validate(); reason != ReasonNone {
			return emitBlock(ctx, repo, input, reason, nil)
		}
	}

	patient, err := NewPatientContextCheck(
		input.SufferingIndicator,
		input.ActiveListeningNote,
		input.Landscape,
		input.EthicalConfirmation,
		input.FailedAttempts,
		input.SupervisorConsulted,
	)
	if err != nil {
		return emitBlock(ctx, repo, input, FromError(err), err)
	}
	if gateErr := patient.GateError(); gateErr != nil {
		return emitBlock(ctx, repo, input, FromError(gateErr), gateErr)
	}
	return emitSuccess(ctx, repo, input, patient)
}

func emitBlock(ctx context.Context, repo AuditRepository, input CareEventInput, reason CognitiveBlockReason, cause error) error {
	snapshot := map[string]interface{}{
		"patient": map[string]interface{}{
			"suffering": input.SufferingIndicator,
			"landscape": input.Landscape.String(),
		},
	}
	snapJSON, _ := json.Marshal(snapshot)

	event := CareEventBlockedByCognitiveState{
		ID:           newAuditID(),
		CareID:       input.ID,
		OperatorID:   input.OperatorID,
		PatientHash:  hashSHA256(input.PatientID),
		Timestamp:    time.Now().UTC(),
		Reason:       reason,
		GateSnapshot: string(snapJSON),
	}
	_ = repo.Append(ctx, event)
	return &BlockError{Reason: reason, Cause: cause}
}

func emitSuccess(ctx context.Context, repo AuditRepository, input CareEventInput, patient PatientContextCheck) error {
	snapshot := map[string]interface{}{
		"patient": map[string]interface{}{
			"suffering": input.SufferingIndicator,
			"landscape": patient.Landscape().String(),
		},
	}
	snapJSON, _ := json.Marshal(snapshot)

	event := CognitiveStateValidated{
		ID:           newAuditID(),
		CareID:       input.ID,
		OperatorID:   input.OperatorID,
		PatientHash:  hashSHA256(input.PatientID),
		Timestamp:    time.Now().UTC(),
		GateSnapshot: string(snapJSON),
	}
	_ = repo.Append(ctx, event)
	return nil
}

// newAuditID returns an opaque, collision-resistant audit event identifier.
func newAuditID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "audit-id-unavailable"
	}
	return hex.EncodeToString(b[:])
}

// hashSHA256 returns the hex-encoded SHA-256 digest of s.
func hashSHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
