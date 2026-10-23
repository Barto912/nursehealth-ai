package adapter

import (
	"github.com/Barto912/nursehealth-ai/pkg/cognitive"
	"github.com/Barto912/nursehealth-ai/pkg/triage"
)

// RuidoGateAdapter implements triage.OperatorGate on top of the cognitive
// RuidoCheck Value Object (AX-001 to AX-003). It lives in pkg/adapter to
// keep pkg/triage and pkg/cognitive mutually independent (bounded contexts).
type RuidoGateAdapter struct {
	check cognitive.RuidoCheck
}

// NewRuidoGateAdapter wraps an already built RuidoCheck snapshot.
// The caller (wiring layer) builds the Value Object by composite literal.
func NewRuidoGateAdapter(check cognitive.RuidoCheck) *RuidoGateAdapter {
	return &RuidoGateAdapter{check: check}
}

// Validate translates the cognitive gates into a typed block reason.
// Gate order follows the axioms: AX-001 (clarity), AX-002 (work level),
// AX-003 (simple attention or recorded attentional finding).
// It satisfies triage.OperatorGate.
func (a RuidoGateAdapter) Validate() triage.CognitiveBlockReason {
	rc := a.check
	if !rc.AllowsCareEvent() {
		if rc.MentalClarity < cognitive.MinMentalClarity {
			return triage.ReasonMentalClarityTooLow
		}
		return triage.ReasonWorkLevelTooLow
	}
	if !rc.AllowsSNOMEDSignature() {
		return triage.ReasonAttentionNotFocused
	}
	return triage.ReasonNone
}