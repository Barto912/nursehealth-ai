# ADR-007: FHIR R4 Extensions for Patient-Side Cognitive Gates (MVP)

## Status
Accepted

## Context
The patient-side axioms AX-004 to AX-006 (docs/philosophy/axioms.json) gate
CareEvent closure and SNOMED CT signature. The interoperability layer
(Cloud Healthcare API, FHIR R4) requires a standard wire representation of
these gates so external auditors can verify them without reading internal
Go types.

## Decision
Represent the gates with three FHIR R4 artifacts, isolated behind the
existing Anti-Corruption Layer:

1. Suffering (AX-004): Observation with category survey or mental-health.
   Code SNOMED CT 248096005 (existential suffering) versus 22253000
   (physical pain). The ActiveListeningNote travels as Observation.note;
   anti-gaming validation is enforced upstream by pkg/triage, never by
   the wire layer.
2. Formation landscape (AX-005): custom Patient extension
   http://nursehealth.ai/fhir/StructureDefinition/formation-landscape
   with valueCode GENERATIONAL, CULTURAL, EMOTIONAL or UNKNOWN.
3. Ethical confirmation (AX-006): Provenance resource attached to every
   critical CareEvent, carrying the VerificationReceipt outcome
   (Approved or Rejected) and the operator agent role.

MVP scope: schema definitions and ACL mapping tables only. No outbound
transmission of these resources until pilot go-live (December 2026).

## Reversal Criteria
- If FHIR evaluators or ANMAT require standard extensions instead of the
  custom formation-landscape extension, replace it with a standard
  Extension definition before regional scaling.
- If SNOMED CT licensing restricts code usage in exports, switch to a
  local coding system with a published mapping table.

## Consequences
- External auditors can verify gate evidence with standard FHIR queries.
- The custom extension requires public StructureDefinition publication
  (future work, pre-scaling).
- The ACL keeps internal ubiquitous language decoupled from the wire
  format; no internal type leaks into FHIR payloads.

## Sources
- docs/philosophy/axioms.json (AX-004 to AX-006, verbatim quotes and
  source_metadata, H20 verifiability standard).
- ADR-006 (Human Firewall MVP) for the operator-side gates.
