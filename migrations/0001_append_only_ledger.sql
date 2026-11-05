-- ============================================================================
-- NURSEHEALTH.AI v0.6.0-defensa-civil
-- MIGRATION 0001: Forensically Immutable Clinical & Logistical Ledger
-- Jurisdicción: República Argentina (Ley 25.326 / Ley 25.506 / ANMAT)
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS clinical_audit_ledger (
    sequence_id       BIGSERIAL PRIMARY KEY,
    event_id          UUID NOT NULL UNIQUE,
    axiom_code        VARCHAR(16) NOT NULL CHECK (axiom_code IN ('AX-001', 'AX-008', 'AX-009', 'AX-012', 'AX-016')),
    actor_id          VARCHAR(128) NOT NULL,
    public_key_id     VARCHAR(128) NOT NULL,
    nonce_hex         VARCHAR(64) NOT NULL UNIQUE,
    signed_at_utc     TIMESTAMPTZ NOT NULL,
    recorded_at_utc   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    payload_cipher    BYTEA NOT NULL,
    disagreement_note TEXT DEFAULT '',
    signature_der     BYTEA NOT NULL,
    prev_hash         CHAR(64) NOT NULL,
    record_hash       CHAR(64) NOT NULL UNIQUE
);

CREATE OR REPLACE FUNCTION enforce_append_only_and_chain_hash()
RETURNS TRIGGER LANGUAGE plpgsql SECURITY DEFINER SET search_path = public, pg_temp AS $$
DECLARE last_hash TEXT; payload_hex TEXT; sig_hex TEXT; canonical_string TEXT;
BEGIN
    IF TG_OP IN ('UPDATE', 'DELETE', 'TRUNCATE') THEN
        RAISE EXCEPTION 'DEFENSA_CIVIL_SECURITY_VIOLATION: Clinical ledger is strictly append-only (Operation % rejected)', TG_OP USING ERRCODE = '42501';
    END IF;
    PERFORM pg_advisory_xact_lock(9120006);
    SELECT record_hash INTO last_hash FROM clinical_audit_ledger ORDER BY sequence_id DESC LIMIT 1;
    IF last_hash IS NULL THEN last_hash := encode(digest('NURSEHEALTH_DEFENSA_CIVIL_AR_GENESIS_V0_6_0', 'sha256'), 'hex'); END IF;
    NEW.prev_hash := last_hash; NEW.recorded_at_utc := NOW();
    payload_hex := encode(NEW.payload_cipher, 'hex'); sig_hex := encode(NEW.signature_der, 'hex');
    canonical_string := format('NH_DC_V1|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s|%s:%s',
        length(NEW.event_id::text), NEW.event_id::text, length(NEW.axiom_code), NEW.axiom_code,
        length(NEW.actor_id), NEW.actor_id, length(NEW.public_key_id), NEW.public_key_id,
        length(NEW.nonce_hex), NEW.nonce_hex, length(extract(epoch from NEW.signed_at_utc)::bigint::text), extract(epoch from NEW.signed_at_utc)::bigint::text,
        length(payload_hex), payload_hex, length(COALESCE(NEW.disagreement_note, '')), COALESCE(NEW.disagreement_note, ''), length(NEW.prev_hash), NEW.prev_hash);
    NEW.record_hash := encode(digest(canonical_string || '|' || length(sig_hex)::text || ':' || sig_hex, 'sha256'), 'hex');
    RETURN NEW;
END; $$;

DROP TRIGGER IF EXISTS trg_clinical_ledger_immutable ON clinical_audit_ledger;
CREATE TRIGGER trg_clinical_ledger_immutable BEFORE INSERT OR UPDATE OR DELETE ON clinical_audit_ledger FOR EACH ROW EXECUTE FUNCTION enforce_append_only_and_chain_hash();
ALTER TABLE clinical_audit_ledger ENABLE ALWAYS TRIGGER trg_clinical_ledger_immutable;
