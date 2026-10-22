package triage

// CognitiveBlockReason enumerates the typed reasons why a CareEvent was blocked
// by the cognitive firewall. Used in HTTP 403 responses and audit logs.
type CognitiveBlockReason int

const (
	// ReasonNone indicates no block (should not appear in blocked events).
	ReasonNone CognitiveBlockReason = iota

	// Operator-side gates (AX-001 to AX-003, RuidoCheck)
	ReasonMentalClarityTooLow
	ReasonWorkLevelTooLow
	ReasonAttentionNotFocused

	// Patient-side gates (AX-004 to AX-006, PatientContextCheck)
	ReasonInsufficientActiveListening
	ReasonContextNotDeterminable
	ReasonEthicalConfirmationMissing
)

// String returns the machine-readable name of the block reason.
func (r CognitiveBlockReason) String() string {
	switch r {
	case ReasonMentalClarityTooLow:
		return "MentalClarityTooLow"
	case ReasonWorkLevelTooLow:
		return "WorkLevelTooLow"
	case ReasonAttentionNotFocused:
		return "AttentionNotFocused"
	case ReasonInsufficientActiveListening:
		return "InsufficientActiveListening"
	case ReasonContextNotDeterminable:
		return "ContextNotDeterminable"
	case ReasonEthicalConfirmationMissing:
		return "EthicalConfirmationMissing"
	default:
		return "Unknown"
	}
}

// FromError maps a validation error to a typed CognitiveBlockReason.
// Returns ReasonNone if the error is not a known cognitive block.
func FromError(err error) CognitiveBlockReason {
	if err == nil {
		return ReasonNone
	}
	switch err {
	case ErrInsufficientActiveListening:
		return ReasonInsufficientActiveListening
	case ErrContextNotDeterminable, ErrContextEscalationPending:
		return ReasonContextNotDeterminable
	case ErrEthicalConfirmationMissing:
		return ReasonEthicalConfirmationMissing
	default:
		return ReasonNone
	}
}
