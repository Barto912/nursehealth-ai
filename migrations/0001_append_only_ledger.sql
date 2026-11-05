-- =============================================================================
-- NurseHealth.ai: Append-Only Cryptographic Audit Ledger
-- Purpose: Ensure forensic-grade immutability for clinical and logistical events.
-- Compliance: HIPAA, SOC 2 Type II, ANMAT Traceability, Ley 24.051 (AR).
-- =============================================================================

-- 1. Enable cryptographic functions in PostgreSQL
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2. Create the immutable ledger table
CREATE TABLE clinical_event_ledger (
    sequence_id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    patient_pseudonym_id VARCHAR(64) NOT NULL, -- Never stores raw PHI
    edge_device_attestation VARCHAR(128) NOT NULL, -- e.g., Titan M3 signature
    snomed_ct_code VARCHAR(32) NOT NULL,
    fhir_resource_type VARCHAR(64) NOT NULL,
    gate_status VARCHAR(32) NOT NULL CHECK (gate_status IN ('COMMIT_ALLOWED', 'FAIL_CLOSE_VETO', 'HUMAN_REFUTATION')),
    payload_sanitized JSONB NOT NULL, -- Sanitized via pkg/adapter DataProxy
    iot_thermal_telemetry JSONB,      -- e.g., {"temp_c": 4.5, "sscc": "XYZ"}
    prev_hash CHAR(64) NOT NULL,      -- SHA-256 of the previous record
    curr_hash CHAR(64) NOT NULL UNIQUE, -- SHA-256 of the current record
    created_at TIMESTAMPTZ NOT NULL DEFAULT CLOCK_TIMESTAMP()
);

-- 3. Create Indexes for efficient audit queries (without compromising immutability)
CREATE INDEX idx_ledger_patient_pseudonym ON clinical_event_ledger(patient_pseudonym_id);
CREATE INDEX idx_ledger_created_at ON clinical_event_ledger(created_at);
CREATE INDEX idx_ledger_gate_status ON clinical_event_ledger(gate_status);

-- 4. Trigger Function: Enforce SHA-256 Chain Integrity on INSERT
CREATE OR REPLACE FUNCTION enforce_sha256_chain()
RETURNS TRIGGER AS $$
DECLARE
    last_recorded_hash CHAR(64);
    expected_hash CHAR(64);
BEGIN
    -- Get the hash of the most recent record (the "tip" of the chain)
    SELECT curr_hash INTO last_recorded_hash
    FROM clinical_event_ledger
    ORDER BY sequence_id DESC
    LIMIT 1;

    -- If it's the first record (Genesis Block), use a zero-hash
    IF last_recorded_hash IS NULL THEN
        last_recorded_hash := REPEAT('0', 64);
    END IF;

    -- Invariant 1: The provided prev_hash MUST match the actual tip of the chain
    IF NEW.prev_hash != last_recorded_hash THEN
        RAISE EXCEPTION 'CRYPTOGRAPHIC_CHAIN_BROKEN: prev_hash (%) does not match chain tip (%)', NEW.prev_hash, last_recorded_hash;
    END IF;

    -- Invariant 2: Compute the expected curr_hash based on the new record's immutable fields
    expected_hash := ENCODE(
        DIGEST(
            NEW.prev_hash || 
            NEW.event_id::TEXT || 
            NEW.snomed_ct_code || 
            NEW.gate_status || 
            NEW.payload_sanitized::TEXT,
            'sha256'
        ),
        'hex'
    );

    -- Invariant 3: The provided curr_hash MUST match the computed hash
    IF NEW.curr_hash != expected_hash THEN
        RAISE EXCEPTION 'HASH_INTEGRITY_MISMATCH: Computed SHA-256 (%) differs from provided curr_hash (%)', expected_hash, NEW.curr_hash;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 5. Attach the INSERT trigger
CREATE TRIGGER trg_verify_sha256_chain
BEFORE INSERT ON clinical_event_ledger
FOR EACH ROW EXECUTE FUNCTION enforce_sha256_chain();

-- 6. Trigger Function: Absolute Prohibition of UPDATE or DELETE (Append-Only)
CREATE OR REPLACE FUNCTION forbid_ledger_mutation()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'IMMUTABLE_LEDGER_VIOLATION: UPDATE and DELETE operations are strictly prohibited on clinical_event_ledger to preserve forensic auditability.';
END;
$$ LANGUAGE plpgsql;

-- 7. Attach the UPDATE/DELETE trigger
CREATE TRIGGER trg_append_only_guard
BEFORE UPDATE OR DELETE ON clinical_event_ledger
FOR EACH ROW EXECUTE FUNCTION forbid_ledger_mutation();

-- =============================================================================
-- Reasoning Steps for Auditors:
-- 1. pgcrypto provides native, highly optimized SHA-256 hashing within the DB engine.
-- 2. The BEFORE INSERT trigger ensures that no application-layer bug can insert 
--    a broken chain. The database itself enforces the cryptographic invariant.
-- 3. The BEFORE UPDATE/DELETE trigger guarantees append-only semantics at the 
--    storage layer, satisfying SOC 2 Type II and ANMAT non-repudiation requirements.
-- =============================================================================
