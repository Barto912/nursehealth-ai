package triage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Barto912/nursehealth-ai/pkg/triage"
)

// fakeAuditRepo is an in-memory, append-only audit repository for tests.
type fakeAuditRepo struct {
	events []triage.CognitiveEvent
}

func (f *fakeAuditRepo) Append(_ context.Context, e triage.CognitiveEvent) error {
	f.events = append(f.events, e)
	return nil
}

// fakeOperatorGate is a deterministic OperatorGate that injects a fixed
// block reason (or ReasonNone) into the middleware validation path.
type fakeOperatorGate struct {
	reason triage.CognitiveBlockReason
}

func (g fakeOperatorGate) Validate() triage.CognitiveBlockReason { return g.reason }

// validInput returns a CareEventInput that satisfies all patient-side
// invariants (AX-004, AX-005, AX-006): the active-listening note is 89
// runes (above the 50-rune minimum) with non-generic tokens, the landscape
// is determined, and ethical confirmation is present.
func validInput() triage.CareEventInput {
	return triage.CareEventInput{
		ID:                  "care-evt-001",
		PatientID:           "patient-abc",
		OperatorID:          "operator-xyz",
		Operator:            fakeOperatorGate{reason: triage.ReasonNone},
		SufferingIndicator:  true,
		ActiveListeningNote: "La paciente expresa angustia existencial recurrente por aislamiento familiar prolongado.",
		Landscape:           triage.LandscapeEmotional,
		EthicalConfirmation: true,
		FailedAttempts:      0,
		SupervisorConsulted: false,
	}
}

func TestValidateCareEvent_AllGatesPass(t *testing.T) {
	repo := &fakeAuditRepo{}
	if err := triage.ValidateCareEvent(context.Background(), repo, validInput()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(repo.events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(repo.events))
	}
	if got := repo.events[0].EventName(); got != "CognitiveStateValidated" {
		t.Fatalf("expected CognitiveStateValidated, got %q", got)
	}
}

func TestValidateCareEvent_OperatorBlock_MentalClarity(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.Operator = fakeOperatorGate{reason: triage.ReasonMentalClarityTooLow}

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	var be *triage.BlockError
	if !errors.As(err, &be) || be.Reason != triage.ReasonMentalClarityTooLow {
		t.Fatalf("expected BlockError with ReasonMentalClarityTooLow, got %v", err)
	}
	if len(repo.events) != 1 || repo.events[0].EventName() != "CareEventBlockedByCognitiveState" {
		t.Fatalf("expected blocked audit event, got %v", repo.events)
	}
}

func TestValidateCareEvent_OperatorBlock_WorkLevel(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.Operator = fakeOperatorGate{reason: triage.ReasonWorkLevelTooLow}

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	var be *triage.BlockError
	if !errors.As(err, &be) || be.Reason != triage.ReasonWorkLevelTooLow {
		t.Fatalf("expected BlockError with ReasonWorkLevelTooLow, got %v", err)
	}
}

func TestValidateCareEvent_OperatorBlock_Attention(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.Operator = fakeOperatorGate{reason: triage.ReasonAttentionNotFocused}

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	var be *triage.BlockError
	if !errors.As(err, &be) || be.Reason != triage.ReasonAttentionNotFocused {
		t.Fatalf("expected BlockError with ReasonAttentionNotFocused, got %v", err)
	}
}

func TestValidateCareEvent_PatientBlock_ActiveListening(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.ActiveListeningNote = "ok" // AX-004 anti-gaming: under 50 runes and generic token

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	if !errors.Is(err, triage.ErrInsufficientActiveListening) {
		t.Fatalf("expected ErrInsufficientActiveListening, got %v", err)
	}
	var be *triage.BlockError
	if !errors.As(err, &be) {
		t.Fatalf("expected *BlockError wrapper, got %T", err)
	}
	if len(repo.events) != 1 || repo.events[0].EventName() != "CareEventBlockedByCognitiveState" {
		t.Fatalf("expected blocked audit event, got %v", repo.events)
	}
}

func TestValidateCareEvent_PatientBlock_EthicalMissing(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.EthicalConfirmation = false // AX-006 fail-close

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	if !errors.Is(err, triage.ErrEthicalConfirmationMissing) {
		t.Fatalf("expected ErrEthicalConfirmationMissing, got %v", err)
	}
}

func TestValidateCareEvent_PatientBlock_ContextNotDeterminable(t *testing.T) {
	repo := &fakeAuditRepo{}
	input := validInput()
	input.Landscape = triage.LandscapeUnknown
	input.FailedAttempts = 3 // maxFailedAttempts+1 -> StepBlock (AX-005)

	err := triage.ValidateCareEvent(context.Background(), repo, input)
	if !errors.Is(err, triage.ErrContextNotDeterminable) {
		t.Fatalf("expected ErrContextNotDeterminable, got %v", err)
	}
}

func TestBlockError_Unwrap(t *testing.T) {
	cause := errors.New("underlying validation cause")
	be := &triage.BlockError{Reason: triage.ReasonAttentionNotFocused, Cause: cause}
	if !errors.Is(be, cause) {
		t.Fatal("BlockError.Unwrap must expose the underlying cause to errors.Is")
	}
}

func TestCognitiveBlockReason_String_Stable(t *testing.T) {
	reasons := []triage.CognitiveBlockReason{
		triage.ReasonNone,
		triage.ReasonMentalClarityTooLow,
		triage.ReasonWorkLevelTooLow,
		triage.ReasonAttentionNotFocused,
	}
	seen := map[string]bool{}
	for _, r := range reasons {
		s := r.String()
		if s == "" {
			t.Fatalf("reason %d returned empty string", r)
		}
		if seen[s] {
			t.Fatalf("duplicate String() value %q", s)
		}
		seen[s] = true
	}
}