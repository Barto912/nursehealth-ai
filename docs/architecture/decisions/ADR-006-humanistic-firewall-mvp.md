# ADR-006: Human Firewall as Traceability Precondition (MVP)

## Status
Accepted

## Context
Iatrogenesis in home nursing does not arise solely from lack of technical
knowledge, but from noise interference (tensions, emotional climates,
inertia from previous events) in the operator's vigilance level (AX-001).
Technology alone cannot correct this: it must require prior human validation.

## Decision
Implement a `RuidoCheck` Value Object in the Cognitive bounded context
(`pkg/cognitive`) that acts as a domain invariant. No `CareEvent` can be
created or signed without a valid `RuidoCheck` attached as
`CognitiveSnapshot`.

MVP Gates:

- `MentalClarity >= 7` (AX-001). Conservative threshold; hypothesis to
  validate with cognitive_audit_log and 3 real nurses.
- `WorkLevel != LOW` (AX-002). In MVP, gate applies to all CareEvents.
  Reversal criterion: when procedure criticality catalog exists, gate
  will relax to critical procedures only.
- `SimpleAttention == true` or non-empty `AttentionalFinding` (AX-003,
  Autoliberación Serie 13) to sign SNOMED event.

Every blocked attempt is logged to `cognitive_audit_log` (immutable table).
Zero external hardware dependency.

## Consequences
- Clinical domain (CareEvent) depends on prior cognitive domain validation.
- Immutable traceability of operator fatigue or blocking.
- No automatic escalation in MVP: human supervisor queries
  cognitive_audit_log manually. Pattern-based alerts are future work.

## Source Correction (H20)
Initial AX-003 version cited a passage from Silo_OC2_es.pdf that was not
found in the text during verification. The axiom was re-sourced to
Autoliberacion.pdf (Serie 13, Atención simple). Domain ubiquitous language
follows the verified source: SimpleAttention and AttentionalFinding.

## Uninvoked Sources
No unverified external studies are cited in the project knowledge base.
Thresholds are MVP hypotheses, not cited truths.
