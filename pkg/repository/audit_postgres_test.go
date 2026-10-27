package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Barto912/nursehealth-ai/pkg/triage"
)

type fakeResult struct{}

func (f fakeResult) LastInsertId() (int64, error) { return 0, nil }
func (f fakeResult) RowsAffected() (int64, error) { return 1, nil }

type fakeDBExecutor struct {
	err   error
	calls int
	query string
	args  []interface{}
}

func (f *fakeDBExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	f.calls++
	f.query = query
	f.args = args
	if f.err != nil {
		return nil, f.err
	}
	return fakeResult{}, nil
}

func TestPostgresAuditRepository_Append_Success(t *testing.T) {
	db := &fakeDBExecutor{}
	repo := NewPostgresAuditRepository(db)

	event := triage.CognitiveStateValidated{
		ID:           "evt-001",
		CareID:       "care-101",
		OperatorID:   "nurse-42",
		PatientHash:  "sha256-patient-abc",
		Timestamp:    time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC),
		GateSnapshot: `{"state":"ok"}`,
	}

	err := repo.Append(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if db.calls != 1 {
		t.Fatalf("expected 1 ExecContext call, got %d", db.calls)
	}
}

func TestPostgresAuditRepository_Append_DBError(t *testing.T) {
	expectedErr := errors.New("db connection refused")
	db := &fakeDBExecutor{err: expectedErr}
	repo := NewPostgresAuditRepository(db)

	event := triage.CognitiveStateValidated{
		ID:           "evt-002",
		CareID:       "care-102",
		OperatorID:   "nurse-42",
		PatientHash:  "sha256-patient-xyz",
		Timestamp:    time.Date(2026, 10, 27, 12, 5, 0, 0, time.UTC),
		GateSnapshot: `{"state":"ok"}`,
	}

	err := repo.Append(context.Background(), event)
	if err == nil {
		t.Fatal("expected error from DB, got nil")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error wrapping %v, got %v", expectedErr, err)
	}
	if db.calls != 1 {
		t.Fatalf("expected 1 ExecContext call, got %d", db.calls)
	}
}
