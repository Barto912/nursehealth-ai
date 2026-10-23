package adapter

import (
	"nursehealth-ai/pkg/cognitive"
	"nursehealth-ai/pkg/triage"
)

// RuidoGateAdapter implements triage.OperatorGate using pkg/cognitive.RuidoCheck.
// This adapter lives in a separate package to prevent circular dependencies
// between the triage domain and the cognitive domain.
type RuidoGateAdapter struct {
	check cognitive.RuidoCheck
}

// NewRuidoGateAdapter creates a new adapter instance.
func NewRuidoGateAdapter(clarity int, level cognitive.WorkLevel, attention bool) *RuidoGateAdapter {
	return &RuidoGateAdapter{
		check: cognitive.NewRuidoCheck(clarity, level, attention),
	}
}

// Validate executes the cognitive gate and returns the typed block reason.
// It satisfies the triage.OperatorGate interface.
func (a *RuidoGateAdapter) Validate() triage.CognitiveBlockReason {
	if err := a.check.Validate(); err != nil {
		return triage.FromError(err)
	}
	return triage.ReasonNone
}
