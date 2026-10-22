package triage

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// AuditRepository defines the contract for persisting cognitive audit events.
// Implementations must be append-only and thread-safe.
type AuditRepository interface {
	Append(ctx context.Context, event CognitiveEvent) error
}

// CareEventInput represents the minimal data required to validate a care event.
// In a full system, this would be the domain CareEvent aggregate.
type CareEventInput struct {
	ID         string
	PatientID  string // Will be hashed before logging
	OperatorID string
	
	// Operator-side state (AX-001 to AX-003)
	MentalClarity int // 0-100
	WorkLevel     WorkLevel
	Attention     bool
	
	// Patient-side state (AX-004 to AX-006)
	SufferingIndicator  bool
	ActiveListeningNote string
	Landscape           ContextualLandscape
	EthicalConfirmation bool
	FailedAttempts      int
	SupervisorConsulted bool
}

// ValidateCareEvent executes the cognitive firewall gates.
// It returns nil if the event is allowed, or a typed error if blocked.
// Side effect: emits an audit event to the repository (success or failure).
func ValidateCareEvent(ctx context.Context, repo AuditRepository, input CareEventInput) error {
	// 1. Operator-side validation (RuidoCheck)
	ruido := NewRuidoCheck(input.MentalClarity, input.WorkLevel, input.Attention)
	if err := ruido.Validate(); err != nil {
		return emitBlock(ctx, repo, input, err)
	}
	
	// 2. Patient-side validation (PatientContextCheck)
	patient, err := NewPatientContextCheck(
		input.SufferingIndicator,
		input.ActiveListeningNote,
		input.Landscape,
		input.EthicalConfirmation,
		input.FailedAttempts,
		input.SupervisorConsulted,
	)
	if err != nil {
		return emitBlock(ctx, repo, input, err)
	}
	
	// 3. Additional gate check for patient context (escalation pending)
	if gateErr := patient.GateError(); gateErr != nil {
		return emitBlock(ctx, repo, input, gateErr)
	}
	
	// 4. Success path: emit validated event
	return emitSuccess(ctx, repo, input, ruido, patient)
}

func emitBlock(ctx context.Context, repo AuditRepository, input CareEventInput, blockErr error) error {
	reason := FromError(blockErr)
	snapshot := map[string]interface{}{
		"operator": map[string]int{"mental_clarity": input.MentalClarity},
		"patient":  map[string]bool{"suffering": input.SufferingIndicator},
	}
	snapJSON, _ := json.Marshal(snapshot)
	
	event := CareEventBlockedByCognitiveState{
		ID:          generateUUID(), // Placeholder for real UUID generation
		CareID:      input.ID,
		OperatorID:  input.OperatorID,
		PatientHash: hashSHA256(input.PatientID),
		Timestamp:   time.Now().UTC(),
		Reason:      reason,
		GateSnapshot: string(snapJSON),
	}
	
	// Fire-and-forget audit log; do not fail the request if logging fails,
	// but in production this should be robust.
	_ = repo.Append(ctx, event)
	
	return blockErr
}

func emitSuccess(ctx context.Context, repo AuditRepository, input CareEventInput, ruido RuidoCheck, patient PatientContextCheck) error {
	snapshot := map[string]interface{}{
		"ruido":   map[string]int{"mental_clarity": input.MentalClarity},
		"patient": map[string]string{"landscape": patient.Landscape().String()},
	}
	snapJSON, _ := json.Marshal(snapshot)
	
	event := CognitiveStateValidated{
		ID:          generateUUID(),
		CareID:      input.ID,
		OperatorID:  input.OperatorID,
		PatientHash: hashSHA256(input.PatientID),
		Timestamp:   time.Now().UTC(),
		GateSnapshot: string(snapJSON),
	}
	
	_ = repo.Append(ctx, event)
	return nil
}

// Helpers (placeholders for real implementation)
func generateUUID() string { return "uuid-placeholder" }
func hashSHA256(s string) string { return "hash-placeholder" }
