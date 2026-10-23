package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"nursehealth-ai/pkg/triage"
	"github.com/jackc/pgx/v5"
)

// PostgresAuditRepository implements triage.AuditRepository using PostgreSQL.
// It writes append-only records to the cognitive_audit_log table.
type PostgresAuditRepository struct {
	pool *pgx.Pool
}

// NewPostgresAuditRepository creates a new repository instance.
func NewPostgresAuditRepository(ctx context.Context, connString string) (*PostgresAuditRepository, error) {
	pool, err := pgx.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	return &PostgresAuditRepository{pool: pool}, nil
}

// Close releases the connection pool resources.
func (r *PostgresAuditRepository) Close() {
	r.pool.Close()
}

// Append writes a cognitive event to the audit log in an append-only manner.
func (r *PostgresAuditRepository) Append(ctx context.Context, event triage.CognitiveEvent) error {
	var decision string
	var reason string
	var receiptJSON []byte

	// Type switch to determine the event type and extract fields
	switch e := event.(type) {
	case triage.CognitiveStateValidated:
		decision = "PASSED"
		reason = "None"
		receiptJSON = nil
	case triage.CareEventBlockedByCognitiveState:
		decision = "BLOCKED"
		reason = e.Reason.String()
		receiptJSON = []byte(e.ReceiptJSON)
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}

	// Marshal the gate snapshot for JSONB storage
	snapshotJSON, err := json.Marshal(map[string]interface{}{
		"event_name":    event.EventName(),
		"gate_snapshot": getGateSnapshot(event),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Append-only insert. We do not provide prev_hash here; a DB trigger 
	// or a separate hashing service should compute it to ensure true immutability.
	query := `
		INSERT INTO cognitive_audit_log 
		(care_event_id, patient_hash, operator_id, block_reason, gate_snapshot, payload_hash, prev_hash, decision, receipt_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Placeholder hashes: in production, compute SHA-256 of canonical JSON + previous row hash
	payloadHash := "sha256-placeholder-payload"
	prevHash := "sha256-placeholder-prev"

	_, err = r.pool.Exec(ctx, query,
		event.CareEventID(),
		getPatientHash(event),
		getOperatorID(event),
		reason,
		snapshotJSON,
		payloadHash,
		prevHash,
		decision,
		receiptJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to append audit event: %w", err)
	}

	return nil
}

// Helpers to extract common fields via interface or type assertion
func getPatientHash(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.PatientHash
	case triage.CareEventBlockedByCognitiveState:
		return ev.PatientHash
	}
	return ""
}

func getOperatorID(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.OperatorID
	case triage.CareEventBlockedByCognitiveState:
		return ev.OperatorID
	}
	return ""
}

func getGateSnapshot(e triage.CognitiveEvent) string {
	switch ev := e.(type) {
	case triage.CognitiveStateValidated:
		return ev.GateSnapshot
	case triage.CareEventBlockedByCognitiveState:
		return ev.GateSnapshot
	}
	return "{}"
}
