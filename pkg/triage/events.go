package triage

import "time"

// CognitiveEvent represents a domain event emitted by the cognitive firewall.
// These events are immutable records of decision points.
type CognitiveEvent interface {
	EventName() string
	OccurredAt() time.Time
	CareEventID() string
}

// CognitiveStateValidated is emitted when a CareEvent passes all cognitive gates.
type CognitiveStateValidated struct {
	ID          string    // UUID of this audit event
	CareID      string    // UUID of the CareEvent being validated
	OperatorID  string    // UUID of the nurse/operator
	PatientHash string    // SHA-256 of patient ID (pseudonymized)
	Timestamp   time.Time // When the validation occurred
	GateSnapshot string   // JSON representation of the gate state at validation time
}

func (e CognitiveStateValidated) EventName() string    { return "CognitiveStateValidated" }
func (e CognitiveStateValidated) OccurredAt() time.Time { return e.Timestamp }
func (e CognitiveStateValidated) CareEventID() string   { return e.CareID }

// CareEventBlockedByCognitiveState is emitted when a CareEvent is blocked by a gate.
type CareEventBlockedByCognitiveState struct {
	ID          string             // UUID of this audit event
	CareID      string             // UUID of the blocked CareEvent
	OperatorID  string             // UUID of the nurse/operator
	PatientHash string             // SHA-256 of patient ID
	Timestamp   time.Time          // When the block occurred
	Reason      CognitiveBlockReason // Typed reason for the block (from errors.go)
	GateSnapshot string            // JSON representation of the gate state at block time
	ReceiptJSON string             // JSON of VerificationReceipt if applicable
}

func (e CareEventBlockedByCognitiveState) EventName() string    { return "CareEventBlockedByCognitiveState" }
func (e CareEventBlockedByCognitiveState) OccurredAt() time.Time { return e.Timestamp }
func (e CareEventBlockedByCognitiveState) CareEventID() string   { return e.CareID }
