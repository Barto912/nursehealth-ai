package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Barto912/nursehealth-ai/pkg/triage"
)

// fakeDBExecutor implementa DBExecutor con stdlib-only
type fakeDBExecutor struct {
	executedQueries []string
	executedArgs    [][]interface{}
	execError       error
}

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

// fakeCognitiveEvent implementa triage.CognitiveEvent para tests
type fakeCognitiveEvent struct {
	eventName string
	occurred  time.Time
	careID    string
}

func (f fakeCognitiveEvent) EventName() string     { return f.eventName }
func (f fakeCognitiveEvent) OccurredAt() time.Time { return f.occurred }
func (f fakeCognitiveEvent) CareEventID() string   { return f.careID }

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
			event: fakeCognitiveEvent{
				eventName: "RuidoCheckPassed",
				occurred:  time.Unix(1698777600, 0),
				careID:    "evt-123",
			},
			execError:     nil,
			wantErr:       false,
			wantQueryPart: "INSERT INTO cognitive_audit_log",
		},
		{
			name: "database error propagates",
			event: fakeCognitiveEvent{
				eventName: "RuidoCheckFailed",
				occurred:  time.Unix(1698777700, 0),
				careID:    "evt-789",
			},
			execError:     errors.New("connection refused"),
			wantErr:       true,
			wantQueryPart: "INSERT INTO cognitive_audit_log",
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
