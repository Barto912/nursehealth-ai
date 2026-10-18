// Package cognitive implements the Cognitive Ergonomics bounded context.
// It encodes the humanistic doctrine axioms (AX-001 to AX-003) as domain
// invariants that gate clinical events.
package cognitive

// WorkLevel represents the operator's energetic tone at a given moment
// (AX-002, Autoliberacion: "TONOS").
type WorkLevel string

const (
	WorkLevelLow      WorkLevel = "LOW"
	WorkLevelModerate WorkLevel = "MODERATE"
	WorkLevelHigh     WorkLevel = "HIGH"
)

// MinMentalClarity is the minimum self-reported clarity required to create
// a CareEvent (AX-001). Values below this threshold block the event and are
// recorded in cognitive_audit_log.
const MinMentalClarity = 7

// RuidoCheck is the ephemeral Value Object attached to each CareEvent as
// its CognitiveSnapshot. It captures the operator's cognitive state at the
// exact moment of the clinical act.
type RuidoCheck struct {
	// MentalClarity is the self-reported clarity of the vigilic work
	// level, 1 to 10 (AX-001).
	MentalClarity int `json:"mental_clarity"`

	// WorkLevel is the self-reported energetic tone (AX-002).
	WorkLevel WorkLevel `json:"work_level"`

	// SimpleAttention reports whether the operator keeps simple, directed
	// attention on the current task and patient, avoiding monotonous
	// drift (AX-003, Autoliberacion Serie 13).
	SimpleAttention bool `json:"simple_attention"`

	// AttentionalFinding is the descriptive finding the operator must
	// record when SimpleAttention is false, before any SNOMED event can
	// be signed (AX-003).
	AttentionalFinding string `json:"attentional_finding,omitempty"`
}

// AllowsCareEvent reports whether the cognitive state permits creating a
// new CareEvent (gates of AX-001 and AX-002).
func (rc RuidoCheck) AllowsCareEvent() bool {
	return rc.MentalClarity >= MinMentalClarity && rc.WorkLevel != WorkLevelLow
}

// AllowsSNOMEDSignature reports whether the cognitive state permits signing
// the clinical event (gate of AX-003).
func (rc RuidoCheck) AllowsSNOMEDSignature() bool {
	return rc.SimpleAttention || rc.AttentionalFinding != ""
}