package triage

import (
	"errors"
	"strings"
	"testing"
)

// validNote is a clinically meaningful note above the 50-rune threshold.
const validNote = "Se escucha al paciente describir su angustia por el diagnóstico y se acuerda seguimiento conjunto con el equipo de salud mental."

// genericNote exceeds the length threshold but uses only generic tokens.
const genericNote = "ok listo hecho ok listo hecho ok listo hecho ok listo hecho sí no bien"

func TestValidActiveListeningNote(t *testing.T) {
	tests := []struct {
		name string
		note string
		want bool
	}{
		{"49 runes rejected", strings.Repeat("a", 49), false},
		{"50 runes accepted", strings.Repeat("a", 50), true},
		{"whitespace only rejected", "   \t\n  ", false},
		{"punctuation only rejected", "....,;: ... ,,, ;;;", false},
		{"generic tokens only rejected", genericNote, false},
		{"valid clinical note accepted", validNote, true},
		{"empty rejected", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidActiveListeningNote(tt.note); got != tt.want {
				t.Fatalf("ValidActiveListeningNote(%q) = %v, want %v", tt.note, got, tt.want)
			}
		})
	}
}

func TestNewPatientContextCheck_AX004(t *testing.T) {
	if _, err := NewPatientContextCheck(true, strings.Repeat("a", 49), LandscapeCultural, true, 0, false); !errors.Is(err, ErrInsufficientActiveListening) {
		t.Fatalf("suffering with short note: err = %v, want ErrInsufficientActiveListening", err)
	}
	if _, err := NewPatientContextCheck(true, genericNote, LandscapeCultural, true, 0, false); !errors.Is(err, ErrInsufficientActiveListening) {
		t.Fatalf("suffering with generic note: err = %v, want ErrInsufficientActiveListening", err)
	}
	if _, err := NewPatientContextCheck(true, validNote, LandscapeCultural, true, 0, false); err != nil {
		t.Fatalf("suffering with valid note: unexpected err = %v", err)
	}
	if _, err := NewPatientContextCheck(false, "", LandscapeCultural, true, 0, false); err != nil {
		t.Fatalf("no suffering without note: unexpected err = %v", err)
	}
}

func TestPatientContextCheck_AX005_EscalationLadder(t *testing.T) {
	tests := []struct {
		name       string
		attempts   int
		supervisor bool
		wantStep   EscalationStep
		wantGate   error
	}{
		{"first failure retries", 1, false, StepRetry, ErrContextEscalationPending},
		{"second failure requires supervisor", 2, false, StepSupervisor, ErrContextEscalationPending},
		{"third attempt in progress after supervisor", 2, true, StepRetry, ErrContextEscalationPending},
		{"determined landscape needs no escalation", 0, false, StepDetermined, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			landscape := LandscapeUnknown
			if tt.wantStep == StepDetermined {
				landscape = LandscapeCultural
			}
			c, err := NewPatientContextCheck(false, "", landscape, true, tt.attempts, tt.supervisor)
			if err != nil {
				t.Fatalf("construction: unexpected err = %v", err)
			}
			if got := c.EscalationStep(); got != tt.wantStep {
				t.Fatalf("EscalationStep() = %v, want %v", got, tt.wantStep)
			}
			if got := c.GateError(); !errors.Is(got, tt.wantGate) {
				t.Fatalf("GateError() = %v, want %v", got, tt.wantGate)
			}
		})
	}
}

func TestNewPatientContextCheck_AX005_ThirdFailureBlocks(t *testing.T) {
	if _, err := NewPatientContextCheck(false, "", LandscapeUnknown, true, 3, true); !errors.Is(err, ErrContextNotDeterminable) {
		t.Fatalf("third failed attempt: err = %v, want ErrContextNotDeterminable", err)
	}
}

func TestGateError_BlockState(t *testing.T) {
	c := PatientContextCheck{contextualLandscape: LandscapeUnknown, failedAttempts: 3, supervisorConsulted: true}
	if got := c.GateError(); !errors.Is(got, ErrContextNotDeterminable) {
		t.Fatalf("GateError() = %v, want ErrContextNotDeterminable", got)
	}
}

func TestNewPatientContextCheck_AX006_FailClose(t *testing.T) {
	if _, err := NewPatientContextCheck(false, "", LandscapeCultural, false, 0, false); !errors.Is(err, ErrEthicalConfirmationMissing) {
		t.Fatalf("missing ethical confirmation: err = %v, want ErrEthicalConfirmationMissing", err)
	}
}

func TestContextualLandscape_String(t *testing.T) {
	if LandscapeUnknown.String() != "UNKNOWN" || LandscapeGenerational.String() != "GENERATIONAL" ||
		LandscapeCultural.String() != "CULTURAL" || LandscapeEmotional.String() != "EMOTIONAL" {
		t.Fatal("landscape String() names do not match ubiquitous language")
	}
}
