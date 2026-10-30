-- Migración: 002_extend_cognitive_audit_log_atdp.sql
-- Propósito: Extender cognitive_audit_log para capturar la tupla ATDP completa (ADR-009).
-- Tupla: $o observación | $h estado interno | $a acción | $y resultado | $r recompensa | $m metadatos

ALTER TABLE cognitive_audit_log
ADD COLUMN IF NOT EXISTS observation_json JSONB,
ADD COLUMN IF NOT EXISTS hidden_state_summary TEXT,
ADD COLUMN IF NOT EXISTS action_fhir JSONB,
ADD COLUMN IF NOT EXISTS outcome_human_validation VARCHAR(50),
ADD COLUMN IF NOT EXISTS outcome_clinical VARCHAR(50),
ADD COLUMN IF NOT EXISTS outcome_delay_days INTEGER,
ADD COLUMN IF NOT EXISTS reward_signal DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS model_version VARCHAR(20),
ADD COLUMN IF NOT EXISTS device_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS titan_m3_signature TEXT;

COMMENT ON COLUMN cognitive_audit_log.observation_json IS 'Datos clínicos anonimizados por Data Proxy (campo $o de ATDP)';
COMMENT ON COLUMN cognitive_audit_log.hidden_state_summary IS 'Resumen del razonamiento del modelo de borde (campo $h de ATDP)';
COMMENT ON COLUMN cognitive_audit_log.action_fhir IS 'Mapeo SNOMED CT/FHIR R4 generado (campo $a de ATDP)';
COMMENT ON COLUMN cognitive_audit_log.outcome_human_validation IS 'Validación humana: CONFIRMED, CORRECTED, ESCALATED, REFUSED (campo $y)';
COMMENT ON COLUMN cognitive_audit_log.outcome_clinical IS 'Outcome clínico a 7-30 días: IMPROVED, STABLE, DETERIORATED (campo $r)';
COMMENT ON COLUMN cognitive_audit_log.outcome_delay_days IS 'Días entre evento y captura de outcome (recompensas de lazo tardío)';
COMMENT ON COLUMN cognitive_audit_log.reward_signal IS 'Señal de recompensa normalizada (-1.0 a 1.0) para optimización Harness²';
COMMENT ON COLUMN cognitive_audit_log.model_version IS 'Versión del modelo de borde (campo $m de ATDP)';
COMMENT ON COLUMN cognitive_audit_log.device_id IS 'Identificador del Google Pixel 11 Pro (campo $m de ATDP)';
COMMENT ON COLUMN cognitive_audit_log.titan_m3_signature IS 'Firma criptográfica del enclave Titan M3 para no repudio (campo $m)';

CREATE INDEX IF NOT EXISTS idx_audit_log_model_version ON cognitive_audit_log(model_version);
CREATE INDEX IF NOT EXISTS idx_audit_log_outcome_clinical ON cognitive_audit_log(outcome_clinical) WHERE outcome_clinical IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_log_reward_signal ON cognitive_audit_log(reward_signal) WHERE reward_signal IS NOT NULL;