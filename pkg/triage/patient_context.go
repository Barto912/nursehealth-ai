// Package triage implements the clinical decision gates of NurseHealth.ai.
//
// The gates are domain invariants: pure, deterministic and fail-closed.
// No gate may be bypassed by configuration or by a generative model.
package triage

import (
	"errors"
	"strings"
	"unicode"
)

// ContextualLandscape is the formation landscape of the patient (AX-005).
type ContextualLandscape int

// ContextualLandscape values, in ubiquitous language.
const (
	LandscapeUnknown ContextualLandscape = iota
	LandscapeGenerational
	LandscapeCultural
	LandscapeEmotional
)

// String returns the ubiquitous-language name of the landscape.
func (l ContextualLandscape) String() string {
	switch l {
	case LandscapeGenerational:
		return "GENERATIONAL"
	case LandscapeCultural:
		return "CULTURAL"
	case LandscapeEmotional:
		return "EMOTIONAL"
	default:
		return "UNKNOWN"
	}
}

// EscalationStep is the next action required by the AX-005 protocol.
type EscalationStep int

// EscalationStep values, ordered by protocol progression.
const (
	StepDetermined EscalationStep = iota // landscape known: no escalation
	StepRetry                            // failed attempts < 2: operator retries
	StepSupervisor                       // 2 failed attempts: clinical supervisor consult
	StepBlock                            // 3rd failed attempt: ContextNotDeterminable
)

// Fail-closed validation errors. Each maps 1:1 to a blocked invariant.
var (
	ErrInsufficientActiveListening = errors.New("triage: ActiveListeningNote fails anti-gaming validation (AX-004)")
	ErrContextNotDeterminable      = errors.New("triage: ContextualLandscape UNKNOWN after escalation protocol (AX-005)")
	ErrContextEscalationPending    = errors.New("triage: ContextualLandscape UNKNOWN, escalation in progress (AX-005)")
	ErrEthicalConfirmationMissing  = errors.New("triage: EthicalConfirmation missing, fail-close (AX-006)")
)

// MVP hypothesis constants (see ADR-006, Uninvoked Sources).
const (
	minActiveListeningRunes = 50
	maxFailedAttempts       = 2
)

// genericTokens cannot constitute an active-listening note, alone or combined.
var genericTokens = map[string]bool{
	"ok": true, "okay": true, "listo": true, "hecho": true,
	"done": true, "si": true, "sí": true, "no": true,
	"bien": true, "nada": true,
}

// PatientContextCheck is the Value Object that gates CareEvent closure and
// SNOMED signature with the patient-side axioms AX-004, AX-005 and AX-006.
// It is immutable: invariants are validated at construction time.
type PatientContextCheck struct {
	sufferingIndicator  bool
	activeListeningNote string
	contextualLandscape ContextualLandscape
	ethicalConfirmation bool
	failedAttempts      int
	supervisorConsulted bool
}

// NewPatientContextCheck validates the invariants in deterministic order
// (AX-004, AX-005, AX-006) and returns the Value Object or the first
// violated invariant. failedAttempts counts failed landscape determinations.
func NewPatientContextCheck(
	sufferingIndicator bool,
	activeListeningNote string,
	contextualLandscape ContextualLandscape,
	ethicalConfirmation bool,
	failedAttempts int,
	supervisorConsulted bool,
) (PatientContextCheck, error) {
	c := PatientContextCheck{
		sufferingIndicator:  sufferingIndicator,
		activeListeningNote: activeListeningNote,
		contextualLandscape: contextualLandscape,
		ethicalConfirmation: ethicalConfirmation,
		failedAttempts:      failedAttempts,
		supervisorConsulted: supervisorConsulted,
	}
	if err := c.validate(); err != nil {
		return PatientContextCheck{}, err
	}
	return c, nil
}

// validate enforces the three invariants, fail-close, in axiom order.
func (c PatientContextCheck) validate() error {
	if c.sufferingIndicator && !ValidActiveListeningNote(c.activeListeningNote) {
		return ErrInsufficientActiveListening
	}
	if c.contextualLandscape == LandscapeUnknown && c.EscalationStep() == StepBlock {
		return ErrContextNotDeterminable
	}
	if !c.ethicalConfirmation {
		return ErrEthicalConfirmationMissing
	}
	return nil
}

// ValidActiveListeningNote enforces the AX-004 anti-gaming rule: at least
// minActiveListeningRunes runes, not whitespace-only, and not composed
// exclusively of generic tokens.
func ValidActiveListeningNote(note string) bool {
	trimmed := strings.TrimSpace(note)
	if len([]rune(trimmed)) < minActiveListeningRunes {
		return false
	}
	fields := strings.FieldsFunc(trimmed, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})
	if len(fields) == 0 {
		return false
	}
	for _, f := range fields {
		if !genericTokens[strings.ToLower(f)] {
			return true
		}
	}
	return false
}

// EscalationStep returns the next action required by the AX-005 protocol:
// 2 failed attempts, then clinical supervisor, then block on the 3rd failure.
func (c PatientContextCheck) EscalationStep() EscalationStep {
	if c.contextualLandscape != LandscapeUnknown {
		return StepDetermined
	}
	switch {
	case c.failedAttempts >= maxFailedAttempts+1:
		return StepBlock
	case c.failedAttempts == maxFailedAttempts && !c.supervisorConsulted:
		return StepSupervisor
	default:
		return StepRetry
	}
}

// GateError returns the blocking reason of the check, or nil when the
// patient-side invariants admit CareEvent closure and SNOMED CT signature.
// AX-005 blocks both while the landscape remains UNKNOWN.
func (c PatientContextCheck) GateError() error {
	if c.contextualLandscape != LandscapeUnknown {
		return nil
	}
	if c.EscalationStep() == StepBlock {
		return ErrContextNotDeterminable
	}
	return ErrContextEscalationPending
}

// SufferingIndicator reports whether existential suffering was detected.
func (c PatientContextCheck) SufferingIndicator() bool { return c.sufferingIndicator }

// ActiveListeningNote returns the recorded intervention note.
func (c PatientContextCheck) ActiveListeningNote() string { return c.activeListeningNote }

// Landscape returns the patient formation landscape.
func (c PatientContextCheck) Landscape() ContextualLandscape { return c.contextualLandscape }

// EthicalConfirmation reports the operator ethical consent (AX-006).
func (c PatientContextCheck) EthicalConfirmation() bool { return c.ethicalConfirmation }

// FailedAttempts returns the count of failed landscape determinations.
func (c PatientContextCheck) FailedAttempts() int { return c.failedAttempts }

// SupervisorConsulted reports whether the clinical supervisor was consulted.
func (c PatientContextCheck) SupervisorConsulted() bool { return c.supervisorConsulted }
