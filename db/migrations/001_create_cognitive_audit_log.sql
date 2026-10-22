-- Migration: 001_create_cognitive_audit_log
-- Purpose: Append-only audit trail for cognitive firewall decisions (AX-001 to AX-006).
-- Constraint: No UPDATE or DELETE operations allowed on this table (enforced by app logic).

CREATE TABLE IF NOT EXISTS cognitive_audit_log (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Context identifiers (pseudonymized)
    care_event_id   UUID NOT NULL,
    patient_hash    CHAR(64) NOT NULL, -- SHA-256 of patient ID
    operator_id     UUID NOT NULL,
    
    -- Cognitive state snapshot
    block_reason    VARCHAR(50) NOT NULL, -- Maps to triage.CognitiveBlockReason
    gate_snapshot   JSONB NOT NULL,       -- Full state of RuidoCheck + PatientContextCheck
    
    -- Integrity verification
    payload_hash    CHAR(64) NOT NULL,    -- SHA-256 of the canonical JSON payload
    prev_hash       CHAR(64) NOT NULL,    -- Hash of the previous log entry (chain)
    
    -- Outcome
    decision        VARCHAR(20) NOT NULL CHECK (decision IN ('BLOCKED', 'ESCALATED', 'PASSED')),
    receipt_json    JSONB                 -- VerificationReceipt if applicable
);

-- Indexes for audit queries
CREATE INDEX idx_audit_care_event ON cognitive_audit_log(care_event_id);
CREATE INDEX idx_audit_created_at ON cognitive_audit_log(created_at DESC);
CREATE INDEX idx_audit_block_reason ON cognitive_audit_log(block_reason) WHERE block_reason != 'None';

-- Comment enforcing append-only policy
COMMENT ON TABLE cognitive_audit_log IS 'Append-only forensic log. Do not update or delete rows.';
