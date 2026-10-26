package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Barto912/nursehealth-ai/pkg/triage"
)

// fakeDBExecutor implementa DBExecutor con stdlib-only (sql.Result fake)
type fakeDBExecutor struct {
	executedQueries []string
	executedArgs    [][]interface{}
	execError       error
}

// fakeResult implementa sql.Result para tests sin driver real
type fakeResult struct{}

func (f fakeResult) LastInsertId() (int64, error) { return 0, nil }
func (f fakeResult) RowsAffected() (int64, error) { return 1, nil }

func (f *fakeDBExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	f.executedQueries = append(f.executedQueries, query)
	f.executedArgs = append(f.executedArgs, args)
	if f.execError != nil {
		return nil, f.execError
	}
	return fakeResult{}, nil
}

func TestPostgresAuditRepository_Append(t *testing.T) {
	tests := []struct {
		name          string
		event         triage.CognitiveEvent
		execError     error
		wantErr       bool
		wantQueryPart string
	}{
		{
			name: "successful append-only insert",
			event: triage.CognitiveEvent{
				ID:        "evt-123",
				PatientID: "patient-456",
				Reason:    "AX-001",
				Timestamp: 1698777600,
			},
			execError:     nil,
			wantErr:       false,
			wantQueryPart: "INSERT INTO audit_log",
		},
		{
			name: "database error propagates",
			event: triage.CognitiveEvent{
				ID:        "evt-789",
				PatientID: "patient-012",
				Reason:    "AX-002",
				Timestamp: 1698777700,
			},
			execError:     errors.New("connection refused"),
			wantErr:       true,
			wantQueryPart: "INSERT INTO audit_log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExec := &fakeDBExecutor{execError: tt.execError}
			repo := NewPostgresAuditRepository(fakeExec)

			err := repo.Append(context.Background(), tt.event)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if err.Error() != tt.execError.Error() {
					t.Errorf("expected error %q, got %q", tt.execError.Error(), err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}

			if len(fakeExec.executedQueries) != 1 {
				t.Fatalf("expected 1 query, got %d", len(fakeExec.executedQueries))
			}

			query := fakeExec.executedQueries[0]
			if !contains(query, tt.wantQueryPart) {
				t.Errorf("expected query containing %q, got %q", tt.wantQueryPart, query)
			}

			if len(fakeExec.executedArgs) != 1 {
				t.Fatalf("expected 1 args slice, got %d", len(fakeExec.executedArgs))
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}