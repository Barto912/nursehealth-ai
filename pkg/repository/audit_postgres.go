package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"nursehealth-ai/pkg/triage"
)

// DBExecutor abstracts the minimal SQL surface required by the audit
// repository. *sql.DB satisfies this interface. Keeping the dependency
// narrow preserves the stdlib-only rule and allows fake injection in tests.
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// PostgresAuditRepository implements triage.AuditRepository on top of any
// DBExecutor (production: *sql.DB with a registered PostgreSQL driver).
// Writes are append-only into cognitive_audit_log (db/migrations/001).
type PostgresAuditRepository struct {
	db DBExecutor
}

// NewPostgresAuditRepository wires the repository to an injected executor.
// Driver registration and connection lifecycle belong to the composition
// root (cmd/), not to this package.
func NewPostgresAuditRepository(db DBExecutor) *PostgresAuditRepository {
	return &PostgresAuditRepository{db: db}
}

// Append persists a cognitive event as an immutable audit row.
func (r *PostgresAuditRepository) Append(ctx context.Context, event triage.CognitiveEvent) error {
	var decision string
	var reason string
	var receiptJSON []byte

	switch e := event.(type) {
	case triage.CognitiveStateValidated:
		decision = "PASSED"
		reason = "None"
	case triage.CareEventBlockedByCognitiveState:
		decision = "BLOCKED"
		reason = e.Reason.String()
		receiptJSON = []byte(e.ReceiptJSON)
	default:
		return fmt.Errorf("repository: unknown cognitive event type %T", event)
	}

	snapshotJSON, err := json.Marshal(map[string]interface{}{
		"event_name":    event.EventName(),
		"gate_snapshot": gateSnapshotOf(event),
	})
	if err != nil {
		return fmt.Errorf("repository: marshal snapshot: %w", err)
	}

	const query = `
		INSERT INTO cognitive_audit_log
		(care_event_id, patient_hash, operator_id, block_reason, gate_snapshot,
		 payload_hash, prev_hash, decision, receipt_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Hash chaining (SHA-256 of canonical payload plus previous row hash)
	// belongs to the ledger service; placeholders keep this package free
	// of ledger concerns until that service lands.
	payloadHash := "sha256-placeholder-payload"
	prevHash := "sha256-placeholder-prev"

	if _, err := r.db.ExecContext(ctx, query,
		event.CareEventID(),
		patientHashOf(event),
		operatorIDOf(event),
		reason,
		snapshotJSON,
		payloadHash,
		prevHash,
		decision,
		receiptJSON,
	); err != nil {
		return fmt.Errorf("repository: append audit event: %w", err)
	}
	return nil
}

func patientHashOf(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.PatientHash
	case triage.CareEventBlockedByCognitiveState:
		return ev.PatientHash
	}
	return ""
}

func operatorIDOf(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.OperatorID
	case triage.CareEventBlockedByCognitiveState:
		return ev.OperatorID
	}
	return ""
}

func gateSnapshotOf(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.GateSnapshot
	case triage.CareEventBlockedByCognitiveState:
		return ev.GateSnapshot
	}
	return "{}"
}
