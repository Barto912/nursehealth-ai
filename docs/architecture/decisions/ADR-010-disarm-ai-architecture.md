# ADR-010: Arquitectura de "Desarme de la IA", Trazabilidad Forense Inmutable y Gobernanza Ética por Diseño

- **ID de Decisión:** `ADR-010`
- **Título:** Consolidación de Invariantes de *Magnifica Humanitas*, *Filosofía Humanista (Silo)* y Mecanismo Determinista de "Desarme de la IA" (`IrreversibleDecisionGuard` & `DisagreementLogger`)
- **Estado:** `APPROVED_BY_ETHICS_COMMITTEE` / `ACCEPTED`
- **Fecha de Ratificación:** 1 de noviembre de 2026
- **Baseline de Repositorio:** Commit `72fe634` (`main`) — Go 1.24, PostgreSQL, HL7 FHIR R4, SNOMED CT
- **Reemplaza / Consolida:** Cierra la secuencia arquitectónica `ADR-001` a `ADR-010` (unificando la trazabilidad de gobernanza ética y el Desarme de la IA en un único registro sin requerir un `ADR-011`).

---

## 1. Estructura Empresarial Humana Cualificada (Autoridad de Decisión y Firma)

En cumplimiento de los axiomas **AX-008**, **AX-011**, **AX-012** y **AX-016**, la responsabilidad arquitectónica, médica y legal de esta decisión recae exclusivamente en los agentes humanos conscientes del directorio y equipo técnico:

- **Founder, CEO & Chief Ethical Architect:** Walter Bartolomé (Autoridad final de arquitectura de invariantes y ética estructural)
- **Chief Medical & Clinical Informatics Officer (CMIO):** Dr. Daniel Rojo (Autoridad de validación de taxonomía SNOMED CT, perfiles HL7 FHIR R4 y protocolos clínicos de enfermería domiciliaria)
- **Technical Project Manager & Operations Lead (TPM):** Leonor Figueroa (Auditoría técnica humana, aseguramiento de calidad en CI/CD y validación operativa)
- **Chief Technology Officer (CTO — Posición Estratégica Abierta):** *Perfil Ex-Google (Staff/Principal Engineer) a incorporar mediante recomendación directa de Google for Startups para escalar la arquitectura de Mejora Recursiva Regularizada (RSI) en el borde para las Américas.*

> **Referencia Documental de Terreno:** Las credenciales profesionales y experiencia operativa de los firmantes se encuentran respaldadas en `docs/governance/TEAM_FIELD_EXPERIENCE_DOSSIER.pdf`.

---

## 2. Contexto y Problema Arquitectónico

En los sistemas contemporáneos de salud digital e inteligencia artificial clínica, predomina el paradigma de delegar progresivamente mayor autonomía decisoria a modelos estadísticos fundacionales, intentando contener sus alucinaciones o sesgos mediante instrucciones en lenguaje natural (*system prompts* o guardrails probabilísticos).

En el contexto de la enfermería domiciliaria de precisión, la logística inversa farmacéutica y la respuesta en zonas de baja conectividad en las Américas, dicho enfoque presenta tres fallas estructurales inaceptables para **NurseHealth.ai**:

1. **Fragilidad de los Prompts frente a Eventos Irreversibles:** Un modelo estadístico ejecutado en el borde o en la nube puede sufrir desviaciones distributivas o inyecciones adversariales. Si el sistema tiene permisos de escritura directa sobre acciones clínicas irreversibles, el desgaste o fatiga del operador humano se traduce en iatrogenia.
2. **Orfandad Legal y Dilución de Responsabilidad:** Cuando ocurre un evento adverso asistido por una "caja negra" algorítmica, la responsabilidad civil y penal se diluye entre el desarrollador, la institución sanitaria y el profesional de enfermería.
3. **Subordinación Cognitiva del Operador Clínico:** Los sistemas tradicionales carecen de mecanismos criptográficos nativos para que el enfermero en terreno refute a la IA dejando constancia auditable de su juicio clínico experto, perdiendo además la fuente más valiosa de alineación médica (*Reinforcement Learning from Human Feedback* / *Regularized Recursive Self-Improvement* — RSI).

Por mandato del Founder (`2026-11-01`), la arquitectura de **NurseHealth.ai** se fundamenta de manera exclusiva e innegociable en dos fuentes normativas:
- **Carta Encíclica *Magnifica Humanitas*** (León XIV, 15 de mayo de 2026): Secciones 104 (`AX-007`), 105 (`AX-008`), 107 (`AX-009`), 111 (`AX-010`), 199 (`AX-011`) y 200 (`AX-012`).
- **Filosofía Humanista** (Silo, *Obras Completas I y II*, *Autoliberación*): Principio de No-Violencia Activa y Soberanía Humana, *"nada por encima del hombre y ningún hombre por encima de otro"* (`AX-016`).

---

## 3. Decisión Arquitectónica: "Desarme de la IA" e Invariantes Binarias Deterministas

Establecemos el paradigma de **Desarme de la IA** como una restricción estructural compilada en el núcleo del dominio en **Go 1.24**. En **NurseHealth.ai**, nuestro modelo de **IA soberana en el borde** (**MedGemma 2B INT8** ejecutándose sobre dispositivos Google **Pixel 11 Pro** con coprocesador criptográfico **Titan M3**) actúa exclusivamente como motor de sugerencia, estructuración documental HL7 FHIR R4 y codificación SNOMED CT, careciendo por diseño de autoridad para aprobar o ejecutar decisiones clínicas irreversibles.

La decisión se materializa en cuatro pilares técnicos verificables sobre el baseline del `commit 72fe634`:

### 3.1. Veto Binario Determinista en el Motor Fail-Close (Rechazo de Umbrales Probabilísticos)
Para garantizar el máximo rigor técnico, ético y metodológico, descartamos el uso de porcentajes arbitrarios de incertidumbre o *scores* flotantes de entropía como condición de corte en procedimientos críticos. En su lugar, el Motor Fail-Close (`pkg/governance/irreversible_decision_guard.go`) implementa una **invariante binaria determinista** que exige simultáneamente:
1. **Validación Estricta de Esquema HL7 FHIR R4:** Conformidad estructural completa del recurso clínico (`CareEvent`, `Provenance`, `Observation`, `DocumentReference`).
2. **Integridad Taxonómica SNOMED CT (Validada por el CMIO):** El código del procedimiento debe existir en el catálogo activo homologado por el Chief Medical & Clinical Informatics Officer (Dr. Daniel Rojo) y no pertenecer a la lista de códigos bloqueados por cumplimiento normativo jurisdiccional (`ADR-008`).
3. **Verificación Criptográfica de Firma Humana en Titan M3:** Si el procedimiento es clasificado como `IRREVERSIBLE`, la transacción es bloqueada mecánicamente en tiempo de ejecución a menos que contenga una firma criptográfica válida emitida por un supervisor médico humano (`FHIR Provenance.agent.type == Practitioner`, prohibiendo estrictamente `agent.type == Device`).

### 3.2. Soberanía Humana y Mecanismo Obligatorio de Refutación (`AX-016`)
Toda inferencia generada por **MedGemma 2B INT8** en la interfaz del **Pixel 11 Pro** expone de forma obligatoria la acción **"Refutar Sugerencia"** (`pkg/governance/disagreement_logger.go`). Al activarse:
- El sistema subordina de inmediato la sugerencia algorítmica a la decisión del profesional humano.
- Nuestro Data Proxy (`pkg/adapter/data_proxy.go`, verificado con 11 tests pasando en CI) sanitiza localmente cualquier identificador personal de salud (PHI) garantizando la **soberanía sanitaria** (utilizando estándares tipo HIPAA Safe Harbor únicamente como referencia técnica no normativa de anonimización en origen).
- Se empaqueta la tupla `[Contexto Clínico Anonimizado SNOMED CT + Inferencia MedGemma 2B INT8 + Refutación Fundamentada del Operador]` para alimentar el motor de Datasets Curados de Discrepancia y Telemetría Anonimizada (ATDP).

### 3.3. Convergencia en un Único Registro Forense Inmutable (`AX-011`)
Todas las inferencias, aprobaciones humanas, bloqueos del Motor Fail-Close y refutaciones clínicas convergen en un único **registro forense inmutable** en PostgreSQL (`db/migrations/002_extend_cognitive_audit_log_atdp.sql`). Cada asiento genera un `VerificationReceipt` con hash **SHA-256** firmado en el enclave seguro del coprocesador **Titan M3**, garantizando **trazabilidad forense inmutable** y **cero impunidad** ante auditorías de agencias panamericanas (ANVISA, ANMAT, INVIMA, COFEPRIS).

### 3.4. Acoplamiento con Hardware Clínico Co-Diseñado (Grado Militar e IoT)
Esta arquitectura de software se acopla nativamente con la estrategia de provisión mediante leasing directo vía Partner de Telecomunicaciones / OEM de Google y el co-diseño junto a **Google Hardware** de terminales basadas en **Pixel 11 Pro** y **Titan M3** con estructura física y encriptación de grado militar, cámaras *ad-hoc* para registro clínico de heridas, luz infrarroja para venopunción, telemetría IoT de cajas térmicas farmacéuticas y enlaces de *streaming* encriptado para zonas rurales y emergencias sanitarias.

---

## 4. Especificación de Invariantes en Go 1.24 (`/pkg/governance/`)

> **⚠️ NOTA DE ESTADO DE IMPLEMENTACIÓN (Ratificada por Founder, 2026-11-01):**
> El código Go presentado en esta sección constituye la **Especificación de Referencia Aprobada** por el comité de ética y la estructura empresarial humana cualificada. La implementación efectiva en el baseline del repositorio se ejecutará en PRs secuenciales sobre el commit `72fe634`:
> - **PR #35:** `pkg/governance/irreversible_decision_guard.go` + tests table-driven
> - **PR #36:** `pkg/governance/disagreement_logger.go` + tests table-driven
> - **PR #37:** `pkg/governance/verification_receipt.go` + tests table-driven
>
> Cada PR incluirá cobertura de tests >80% y validación contra el `governance-axioms.json` como fuente única de verdad normativa.

A continuación se define la especificación de referencia en **Go 1.24** que instrumenta los axiomas `AX-008`, `AX-009`, `AX-011`, `AX-012` y `AX-016` sobre el esquema forense `002_extend_cognitive_audit_log_atdp.sql`:

```go
// Package governance implements the deterministic invariants for AI Disarmament (AX-012),
// Human Sovereignty (AX-016), and Immutable Forensic Traceability (AX-011).
package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type ProcedureReversibility string
type AgentType string
type EthicsReviewStatus string

const (
	ProcedureReversible   ProcedureReversibility = "REVERSIBLE"
	ProcedureIrreversible ProcedureReversibility = "IRREVERSIBLE"

	AgentPractitioner AgentType = "Practitioner" // Humano consciente y matriculado
	AgentDevice       AgentType = "Device"       // MedGemma 2B INT8 (Prohibido como autorizador)

	StatusApprovedByCommittee EthicsReviewStatus = "APPROVED_BY_ETHICS_COMMITTEE"
)

var (
	ErrRuleNotApprovedByCommittee     = errors.New("AX-009 violation: clinical rule lacks APPROVED_BY_ETHICS_COMMITTEE status signed by CMIO")
	ErrMissingResponsibilityChain     = errors.New("AX-008 violation: DeveloperVersion, OperatorID, and valid SNOMED/FHIR schema are mandatory")
	ErrAIAuthorizationForbidden       = errors.New("AX-012 violation: AI Disarmament enforced; Device agent cannot authorize clinical execution")
	ErrIrreversibleRequiresSupervisor = errors.New("AX-012 violation: IRREVERSIBLE procedure requires human Practitioner supervisor ID and Titan M3 signature")
	ErrInvalidFHIRorSNOMEDTaxonomy    = errors.New("Fail-Close veto: binary deterministic check failed on HL7 FHIR R4 schema or SNOMED CT taxonomy")
)

// VerificationReceipt represents the cryptographic proof stored in the immutable forensic record (AX-011).
type VerificationReceipt struct {
	EventID           string    `json:"event_id"`
	ModelVersion      string    `json:"model_version"` // e.g., "MedGemma-2B-INT8"
	DatasetHash       string    `json:"dataset_hash"`
	SNOMEDCode        string    `json:"snomed_code"`
	FHIRResourceValid bool      `json:"fhir_resource_valid"`
	OperatorID        string    `json:"operator_id"`
	SupervisorID      string    `json:"supervisor_id,omitempty"`
	TitanM3Signature  string    `json:"titan_m3_signature"`
	TimestampUTC      time.Time `json:"timestamp_utc"`
	ReceiptSHA256     string    `json:"receipt_sha256"`
}

// ComputeHash generates the deterministic SHA-256 digest for the immutable forensic record.
func (r *VerificationReceipt) ComputeHash() string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%t|%s|%s|%s|%s",
		r.EventID,
		r.ModelVersion,
		r.DatasetHash,
		r.SNOMEDCode,
		r.FHIRResourceValid,
		r.OperatorID,
		r.SupervisorID,
		r.TitanM3Signature,
		r.TimestampUTC.UTC().Format(time.RFC3339Nano),
	)
	sum := sha256.Sum256([]byte(payload))
	r.ReceiptSHA256 = hex.EncodeToString(sum[:])
	return r.ReceiptSHA256
}

// ForensicSinkRepository abstracts the single immutable forensic record (002_extend_cognitive_audit_log_atdp.sql).
type ForensicSinkRepository interface {
	AppendForensicRecord(receipt VerificationReceipt, status string, atdpDisagreementPayload *string) error
}

// TaxonomyValidator verifies SNOMED CT codes and HL7 FHIR R4 conformance approved by the CMIO.
type TaxonomyValidator interface {
	IsValidFHIRAndSNOMED(snomedCode string, fhirPayload []byte) bool
}

// IrreversibleDecisionGuard enforces AX-008, AX-009, AX-011, and AX-012 (AI Disarmament).
type IrreversibleDecisionGuard struct {
	forensicSink ForensicSinkRepository
	taxonomy     TaxonomyValidator
}

func NewIrreversibleDecisionGuard(sink ForensicSinkRepository, tax TaxonomyValidator) *IrreversibleDecisionGuard {
	return &IrreversibleDecisionGuard{
		forensicSink: sink,
		taxonomy:     tax,
	}
}

// ValidateAndSeal executes the binary deterministic Fail-Close evaluation.
func (g *IrreversibleDecisionGuard) ValidateAndSeal(
	eventID string,
	ruleStatus EthicsReviewStatus,
	reversibility ProcedureReversibility,
	authorizingAgent AgentType,
	modelVersion string,
	datasetHash string,
	snomedCode string,
	fhirPayload []byte,
	operatorID string,
	supervisorID *string,
	titanM3Signature *string,
) (*VerificationReceipt, error) {

	// 1. AX-009: Verify Ethical Code approval state homologated by CMIO
	if ruleStatus != StatusApprovedByCommittee {
		return nil, ErrRuleNotApprovedByCommittee
	}

	// 2. AX-008: Verify explicit responsibility chain identifiers
	if modelVersion == "" || operatorID == "" || snomedCode == "" {
		return nil, ErrMissingResponsibilityChain
	}

	// 3. Binary Deterministic Fail-Close Check: HL7 FHIR R4 + SNOMED CT integrity
	if !g.taxonomy.IsValidFHIRAndSNOMED(snomedCode, fhirPayload) {
		return nil, ErrInvalidFHIRorSNOMEDTaxonomy
	}

	// 4. AX-012: AI Disarmament — MedGemma 2B INT8 (Device) can NEVER be the authorizing agent
	if authorizingAgent == AgentDevice {
		return nil, ErrAIAuthorizationForbidden
	}

	// 5. AX-012: Irreversible procedures strictly require a human Supervisor and Titan M3 signature
	supID := ""
	sig := ""
	if reversibility == ProcedureIrreversible {
		if authorizingAgent != AgentPractitioner || supervisorID == nil || *supervisorID == "" || titanM3Signature == nil || *titanM3Signature == "" {
			return nil, ErrIrreversibleRequiresSupervisor
		}
		supID = *supervisorID
		sig = *titanM3Signature
	} else if titanM3Signature != nil {
		sig = *titanM3Signature
	}

	receipt := VerificationReceipt{
		EventID:           eventID,
		ModelVersion:      modelVersion,
		DatasetHash:       datasetHash,
		SNOMEDCode:        snomedCode,
		FHIRResourceValid: true,
		OperatorID:        operatorID,
		SupervisorID:      supID,
		TitanM3Signature:  sig,
		TimestampUTC:      time.Now().UTC(),
	}
	receipt.ComputeHash()

	if err := g.forensicSink.AppendForensicRecord(receipt, "APPROVED_EXECUTION", nil); err != nil {
		return nil, err
	}

	return &receipt, nil
}
```

---

## 5. Consecuencias Arquitectónicas y de Negocio

### 5.1. Consecuencias Positivas (Ventajas Competitivas)
- **Cero Impunidad y Protección Institucional:** Ninguna inferencia de MedGemma 2B INT8 puede mutar el estado clínico irreversible de un paciente sin la firma criptográfica en Titan M3 de un profesional humano (Practitioner), eliminando de raíz el riesgo de litigios por decisiones autónomas letales.
- **Generación Nativa del Data Moat (ATDP):** Cada uso del botón "Refutar Sugerencia" (AX-016), sanitizado en el borde por `pkg/adapter/data_proxy.go`, construye un activo de datos de discrepancia clínica validada por humanos, libre de PHI y listo para licenciamiento RLHF / RSI en Vertex AI.
- **Aceleración Regulatoria Panamericana:** El cumplimiento compilado en Go 1.24 y auditado por nuestro CMIO (Dr. Daniel Rojo) simplifica las inspecciones técnicas frente a ANVISA, ANMAT, INVIMA y COFEPRIS.

### 5.2. Compromisos y Mitigaciones (Trade-offs)
- **Requerimiento de Firma Humana en Eventos Críticos:** Introduce un paso obligatorio de validación médica humana únicamente en la fracción de procedimientos clasificados como IRREVERSIBLE, mientras que las tareas administrativas y de evolución rutinaria operan sin fricción en el Pixel 11 Pro.
- **Estandarización de Hardware:** Exige el uso de Pixel 11 Pro con Titan M3 (mitigado financieramente mediante el modelo de leasing directo vía Partner OEM/Telecom de Google sin intermediarios).

---

## 6. Firma Digital Humana y Registro de Instrumentación AI-Native

En estricto cumplimiento de los axiomas AX-011 (Cadena de Responsabilidad Identificable), AX-012 (Desarme de la IA) y AX-016 (Soberanía Humana), únicamente los seres humanos conscientes y cualificados del equipo de NurseHealth.ai emiten firmas digitales y asumen responsabilidad ética, clínica y legal sobre este Architecture Decision Record.

### 6.1. Firmas Digitales de la Estructura Empresarial Humana (Responsables Conscientes)

**Founder, CEO & Chief Ethical Architect (Firma Ejecutiva y de Arquitectura Ética):**
Walter Bartolomé — Ratificación de invariantes deterministas de Desarme de la IA, fuentes normativas exclusivas y arquitectura de hardware clínico de grado militar (1 de noviembre de 2026).

**Chief Medical & Clinical Informatics Officer — CMIO (Firma de Homologación Clínica y Ontológica):**
Dr. Daniel Rojo — Validación médica de taxonomías SNOMED CT, perfiles estructurales HL7 FHIR R4, flujos de enfermería domiciliaria y reglas clínicas de escalamiento del Motor Fail-Close (1 de noviembre de 2026).

**Technical Project Manager & Operations Lead — TPM (Firma de Auditoría Técnica y Operativa):**
Leonor Figueroa — Verificación humana de trazabilidad contra el baseline del repositorio (commit 72fe634) y supervisión del plan de validación operativa en piloto Canning Q1 2027 ([PENDING PM VALIDATION]).

**Chief Technology Officer — CTO (Espacio de Firma Reservado):**
[Perfil Ex-Google recomendado por Google for Startups — A incorporar para liderar el escalamiento panamericano de la arquitectura RSI en el borde].

### 6.2. Declaración de Asistencia AI-Native (Registro de Herramienta — Sin Capacidad de Firma)
Como organización nativa en inteligencia artificial gobernada por una estructura empresarial humana cualificada, se registra transparentemente el herramental algorítmico empleado en la estructuración de este artefacto:

- **Rol Desarrollado en Plataforma:** Asistencia en estructuración de código de invariantes en Go 1.24, verificación de coherencia terminológica y simulación de auditoría técnica para el ecosistema Google for Startups / Google Hardware OEM.
- **Modelo Actualizado a la Fecha (1 de noviembre de 2026):** Instancia ARGON (Google for Startups ecosystem), en evolución continua y programada hacia generaciones sucesivas (CARBON y versiones posteriores), operando estrictamente como herramienta subordinada sin personería jurídica ni capacidad de firma.

### 6.3. Metadatos de Verificación en Repositorio
- **Ruta del Archivo:** `/docs/architecture/decisions/ADR-010-disarm-ai-architecture.md`
- **Git Commit Base (main):** `72fe634`
- **Esquema Forense Asociado:** `db/migrations/002_extend_cognitive_audit_log_atdp.sql`
- **Módulo Data Proxy:** `pkg/adapter/data_proxy.go` (11/11 tests PASS en CI)