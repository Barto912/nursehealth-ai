# NurseHealth.ai — Deterministic Clinical-Logistical Operating System for Pan-American Home Healthcare & Green Biosecurity

[![Release](https://img.shields.io/badge/release-v0.4.0--research--proposal-blue.svg)](https://github.com/Barto912/nursehealth-ai/releases)
[![Go Backend](https://img.shields.io/badge/Go-100%25%20stdlib--only-00ADD8.svg)](https://golang.org/)
[![Edge AI](https://img.shields.io/badge/Edge%20AI-MedGemma%202B%20INT8%20%7C%20LiteRT-4285F4.svg)](https://ai.google.dev/gemma)
[![Cloud AI](https://img.shields.io/badge/Cloud%20AI-GCP%20Vertex%20AI%20%7C%20MedGemma%2027B-34A853.svg)](https://cloud.google.com/vertex-ai)
[![Interoperability](https://img.shields.io/badge/HL7-FHIR%20R4%20%7C%20SNOMED%20CT-orange.svg)](https://www.hl7.org/fhir/)
[![Security](https://img.shields.io/badge/Zero--Trust-SHA--256%20Ledger%20%7C%20Gitleaks-critical.svg)](https://github.com/gitleaks/gitleaks)

> **🌍 Multilingual Documentation (Pan-American Health Security Doctrine):**
> * 🇺🇸 **[English (Official Main Interface — US/Global/GCP Review)](./README.md)**
> * 🇦🇷 **[Español (Hispanoamérica & Operación Piloto Canning)](./README_ES.md)**
> * 🇧🇷 **[Português (Expansão Brasil, ANVISA RDC Nº 657/2022 & SUS/RNDS)](./README_PT.md)**

---

## 1. Executive Summary & Foundational Postulate

> **"Judgment is human; traceability belongs to the system."**

**NurseHealth.ai** is NOT a consumer medical chatbot. It is a deterministic, Fail-Close clinical-logistical Operating System designed for precision home healthcare and green pharmaceutical inverse biosecurity across the Americas, structured as a **Delaware C-Corporation** with hybrid clinical validation operations across the United States and Argentina.

In the Americas, extra-hospital home care is severely fragmented, lacking cryptographic chain-of-custody verification, leaving biological cold chains (`+2°C to +8°C`) unverified at the last mile, and generating hazardous domestic pharmaceutical waste. NurseHealth.ai resolves this systemic failure by deploying a **Zero-Trust Hybrid Edge-to-Cloud Infrastructure** built natively on **Google Cloud Platform (GCP)** and **Google Health AI Developer Foundations (MedGemma)**.

The platform enforces deterministic **Fail-Close** decision gates over probabilistic AI outputs, guarantees patient privacy via on-device PII/PHI redaction prior to cloud ingestion, seals every clinical and logistical transition inside an append-only **SHA-256 cryptographic ledger**, and integrates eco-routed inverse logistics via **Google Maps Platform** to safely reclaim and dispose of expired or temperature-sensitive biologicals.

---

## 2. Three-Tier Hybrid Edge-to-Cloud Architecture (GCP Native)

```text
+-----------------------------------------------------------------------------------+
| TIER 1: MOBILE EDGE NODE (Google Pixel 11 Pro / Android 17 / Offline-First)       |
|                                                                                   |
|  [Field RN Input / Wound Telemetry / Rechargeable IoT Thermal Box (+2C to +8C)]   |
|         |                                                                         |
|         v                                                                         |
|  [MedGemma Small (2B INT8) via LiteRT on Google Tensor G6 NPU (<300ms TTFT)]      |
|         |                                                                         |
|         +--> [Local SQLite + FTS5 Engram Store (100% Decentralized)]              |
|         +--> [Titan M3 Security Enclave (Post-Quantum VerificationReceipt Sign)]  |
|         +--> [On-Device Data Proxy (HIPAA Safe Harbor PII/PHI Censoring)]         |
+-----------------------------------------+-----------------------------------------+
                                          |
                     (Pseudonymized & Signed Event Envelope)
                                          |
                                          v
+-----------------------------------------------------------------------------------+
| TIER 2: DETERMINISTIC INGESTION & FAIL-CLOSE CORE (Google Cloud Run)              |
|                                                                                   |
|  [Model Context Protocol (MCP) Deterministic Barrier (Go stdlib-only)]            |
|         |                                                                         |
|         +--> IF Invariant Violated OR High Entropy --> [FAIL-CLOSE VETO]          |
|         |                                              (Block DB & Escalate RN)   |
|         |                                                                         |
|         +--> IF Validated --> [100% Go stdlib-only Event-Sourcing Core Engine]    |
+-----------------------------------------+-----------------------------------------+
                                          |
                                          v
+-----------------------------------------------------------------------------------+
| TIER 3: CLOUD ORCHESTRATION, MULTIREGION LEDGER & GREEN BIOSECURITY (GCP)         |
|  Regions: us-east1 (Primary US Control) <---> southamerica-west1 (LatAm Operations)|
|                                                                                   |
|  +-----------------------------+  +--------------------------------------------+  |
|  | GCP Vertex AI               |  | DDD Anti-Corruption Layer (ACL)            |  |
|  | - MedGemma Maestro (7B/27B) |  | - Google Cloud Healthcare API              |  |
|  | - Google Harness² Policies  |  | - Canonical HL7 FHIR R4 & SNOMED CT Mapper |  |
|  | - QLoRA Distillation to 2B  |  +--------------------------------------------+  |
|  +-----------------------------+                                                  |
|  +-----------------------------+  +--------------------------------------------+  |
|  | Cloud SQL Enterprise Plus   |  | Google Maps Platform (Blind Tokenization)  |  |
|  | - Append-Only PostgreSQL    |  | - Routes API Eco-Friendly Routing          |  |
|  | - SHA-256 Hash-Chain Ledger |  | - Fleet Engine (+2C..+8C Inverse Logistics)|  |
|  +-----------------------------+  +--------------------------------------------+  |
+-----------------------------------------------------------------------------------+
```

### Tier 1: Mobile Edge Component (Local-First Clinical Execution)

- **Hardware Target & Silicon:** Engineered for Google Pixel 11 Pro (Android 17) leveraging the Google Tensor G6 NPU and the Titan M3 hardware security enclave for cryptographic attestation (`VerificationReceipt`), maintaining backward compatibility with current Android StrongBox/Titan M2 Pixel deployments.
- **Bounded Edge Model:** Executes a quantized MedGemma Small (2B INT8) runtime optimized via LiteRT (formerly TensorFlow Lite) strictly constrained to achieve <300ms Time-To-First-Token (TTFT) on invariant classification tasks while eliminating thermal throttling and battery drain during extended 12-hour nursing shifts.
- **Decentralized Persistent Memory (Engram):** Structured atop a local SQLite + FTS5 engine. Preserves longitudinal clinical context and procedural history 100% offline in low-connectivity rural or suburban zones.
- **Zero-Trust Edge Privacy Proxy:** An on-device Data Proxy automatically redacts and censors all 18 HIPAA Safe Harbor identifiers (PII/PHI) locally on the mobile device before any pseudonymized telemetry crosses the network boundary.

### Tier 2: Central Ingestion & Deterministic Fail-Close Engine

- **100% Go stdlib-only Backend Core:** Deployed on Google Cloud Run using an Event-Sourcing architecture written strictly with the Go standard library (`CGO_ENABLED=0`, zero external third-party dependencies in the clinical domain core). This architectural constraint mathematically eliminates third-party software supply-chain attack vectors.
- **Deterministic MCP Hooks (Fail-Close Barrier):** Probabilistic AI outputs are intercepted by a non-negotiable Model Context Protocol (MCP) gatekeeper barrier implemented in pure Go. If the model output violates hardcoded clinical invariants (e.g., dosage limits, missing SNOMED CT taxonomy, or cold-chain excursions outside +2°C to +8°C) or exhibits uncertainty above the entropy threshold, the engine mechanically executes a **Fail-Close Veto**: it blocks state persistence, emits a signed rejection receipt, and escalates execution to a licensed human professional.
- **Hermetic Development vs. Cloud Compute:** Local engineering workstations operate strictly as hermetic compilation and data-hygiene validation nodes. **100% of heavy model distillation, QLoRA fine-tuning, and synthetic clinical data generation execute exclusively on GCP Vertex AI.**

### Tier 3: Cloud Orchestration, Multiregion Ledger & Green Biosecurity

- **Cloud AI & Google Harness² Scaffolding:** Powered by MedGemma Maestro (7B/27B) on GCP Vertex AI. Governed by dynamic clinical guardrails using the recursive Google Harness² (Propose-Probe-Compose) pattern, adapting safety policies to long-tail clinical scenarios via lightweight declarative configurations without mutating underlying neural weights.
- **DDD Anti-Corruption Layer (ACL) & Interoperability:** Google Cloud Healthcare API is isolated behind a Domain-Driven Design (DDD) Anti-Corruption Layer that deterministically maps normalized nursing narratives into canonical HL7 FHIR R4 resources (Observation, MedicationAdministration, SupplyDelivery) and SNOMED CT clinical ontologies.
- **Cryptographic Append-Only Audit Ledger:** Hosted on multiregion Cloud SQL Enterprise Plus (PostgreSQL) across `southamerica-west1` and `us-east1`. Implements an immutable, append-only Event Store enforced by relational database triggers and cryptographically linked via SHA-256 hash chains, eliminating healthcare billing fraud and providing forensic-grade auditability for insurers and health ministries.
- **Green Biosecurity & Inverse Logistics:** Integrates Google Maps Platform (Routes API Eco-Friendly Routing & Fleet Engine) with rechargeable IoT thermal transport containers (+2°C to +8°C). Operating under **Blind Logistics Tokenization** (zero PHI shared with routing APIs), the system optimizes low-emission home delivery of biologicals and orchestrates the safe reverse-logistics collection of expired or hazardous domestic medications.

---

## 3. Intellectual Property, Licensing & DeepMind Data Flywheel

Built upon the **Health AI Developer Foundations** license, NurseHealth.ai implements a **Hybrid Co-Development and Data Sovereignty Model** that protects Delaware C-Corp intellectual property while generating high-value validation assets for the Google Cloud and Google DeepMind ecosystem:

```text
 ┌─────────────────────────────────────────────────────────────────────────┐
 │        PROPRIETARY INTELLECTUAL PROPERTY (Delaware C-Corp)              │
 │ ──► Deterministic Go Core (stdlib-only) + Fail-Close Triage Engine      │
 │ ──► Decentralized Persistent Memory Architecture (Engram / SQLite FTS5) │
 │ ──► Domain-Driven Design (DDD) Anti-Corruption Layer (ACL) & MCP Hooks  │
 └────────────────────────────────────┬────────────────────────────────────┘
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────┐
 │        GOOGLE CLOUD & AI ECOSYSTEM INTEGRATION LAYER                    │
 │ ──► Open Foundational Weights: MedGemma (Health AI Dev Foundations)     │
 │ ──► QLoRA Distillation & INT8 Quantization for Google Tensor G6 NPU     │
 │ ──► Edge De-Identification Pipeline via Local Data Proxy                │
 └────────────────────────────────────┬────────────────────────────────────┘
                                      ▼
 ┌─────────────────────────────────────────────────────────────────────────┐
 │               ECOSYSTEM RETURN VALUE FOR GOOGLE DEEPMIND                │
 │ ──► Multilingual, De-Identified Real-World Datasets (RWD) from Field    │
 │ ──► Empirical Validation of Gemma Architectures in Offline Edge Care    │
 └─────────────────────────────────────────────────────────────────────────┘
```

### Vertex AI Compute & Scale AI Justification

NurseHealth.ai protects Pan-American healthcare data sovereignty through its on-device Data Proxy while utilizing GCP Vertex AI compute infrastructure for continuous teacher-student model distillation:

- **QLoRA Alignment & Distillation:** Distilling complex clinical reasoning and Harness² safety guardrails from MedGemma Maestro (27B) on Vertex AI down to our edge-optimized MedGemma Small (2B INT8) runtime on LiteRT.
- **Real-World Evidence (RWE) Delta Capture:** Capturing de-identified pairs of `(AI_Proposal, Human_RN_FailClose_Correction)` and IoT cold-chain telemetry from our initial 100-case cohort in Canning, Argentina, scaling seamlessly to our Pan-American multiregion cloud topology (`southamerica-west1` and `us-east1`).

---

## 4. Regulatory Compliance, Humanist Ethics & Security Doctrine

- **Non-Device CDS Regulatory Posture (US FDA):** Engineered in alignment with Section 520(o)(1)(E) of the FD&C Act (21st Century Cures Act). The platform operates exclusively as a deterministic support tool for licensed healthcare professionals who independently review the basis of every structured recommendation prior to clinical action.
- **HIPAA, SOC 2 Type II & Pan-American Data Sovereignty:** Enforces HIPAA Safe Harbor (45 CFR § 164.514(b)(2)) de-identification at the edge, Zero-Trust BeyondCorp access controls, NDAA Section 889 hardware supply-chain compliance, ANMAT (Argentina) clinical traceability, and ANVISA RDC Nº 657/2022 (Brazil) software governance.
- **Humanistic Firewall (Magnifica Humanitas):** Formalized in ADR-001 through ADR-010, prioritizing human dignity, zero algorithmic violence, and the strict prevention of physical or psychological iatrogenesis through deterministic human-in-the-loop escalation. The system enforces **AI Disarmament (AX-012)**: the AI can suggest, but only a credentialed human can authorize irreversible clinical actions.

### 🇦 Argentina (Pilot Jurisdiction & R&D Hub)

- **Data Privacy & Health Law:** Ley 25.326 (Protección de Datos Personales), Ley 26.529 (Derechos del Paciente), and Ley 27.553 (Receta Electrónica y Teleasistencia).
- **Clinical & Pharmaceutical Biosecurity:** Aligned with **ANMAT** Software as a Medical Device (SaMD) guidelines, **Ley Nacional 24.051** (Residuos Peligrosos / Farmacéuticos), and **Ley Provincial PBA 11.347** (Residuos Patogénicos) for reverse logistics traceability.
- **Knowledge Economy R&D Framework (Ley 27.506 - LEC):** Structured under ARCA CLAE Codes **620100** (Software Architecture) and **721000** (Experimental R&D in Engineering and Health Sciences) via our Argentine SAS operating entity.

### 🇺 United States (Corporate & Scaling Hub)

HIPAA, FDA SaMD (Good Machine Learning Practice), SOC 2 Type II, and NDAA Sec. 889 compliance via a Delaware C-Corp parent structure.

### 🇧🇷 Brazil

LGPD, ANVISA (RDC Nº 657/2022 for SaMD & RDC Nº 430/2020 for pharmaceutical cold chain), and native interoperability readiness for **SUS (RNDS)**.

---

## 5. Multi-Region Resilience & Disaster Recovery

NurseHealth.ai guarantees clinical continuity under the **Doctrine of Health Security of the Americas**:

- **Layer 0 (Intra-Regional HA):** Synchronous multi-zone high availability across `southamerica-west1` (Santiago) and `us-east1` (Northern Virginia), achieving **RPO = 0** and **RTO < 60s**.
- **Immutable Fail-Close:** Infrastructure failovers never bypass cognitive gates or patient safety invariants. The `IrreversibleDecisionGuard` validates before any multi-region router; no cloud switch skips the human firewall.
- **Sovereign Continuity:** Architected to prevent single-region cloud outages from ever disrupting active home care.

---

## 6. Repository Structure & Verification Suite

```text
barto912/nursehealth-ai/
├── README.md                                  # Primary English Documentation (US/GCP Review)
├── README_ES.md                               # Spanish Documentation (LatAm / Canning Pilot)
├── README_PT.md                               # Portuguese Documentation (Brazil / ANVISA / SUS)
├── go.mod                                     # Go Module Definition (100% stdlib-only core)
├── cmd/
│   ├── demo-governance/
│   │   └── main.go                            # End-to-End Governance CLI Demo (AX-012 & AX-016)
│   └── validate-axioms/
│       └── main.go                            # Normative Axiom Validator (Single Source of Truth)
├── pkg/
│   ├── governance/
│   │   ├── irreversible_decision_guard.go     # Deterministic Fail-Close Engine (AX-008/009/011/012)
│   │   ├── irreversible_decision_guard_test.go
│   │   ├── disagreement_logger.go             # Human Refutation Logger (AX-016)
│   │   └── disagreement_logger_test.go
│   └── adapter/
│       ├── data_proxy.go                      # PII/PHI Redaction Proxy (HIPAA Safe Harbor)
│       ├── data_proxy_test.go
│       ├── phi_sanitizer_adapter.go           # Adapter Pattern for Governance Integration
│       └── phi_sanitizer_adapter_test.go
├── docs/
│   ├── architecture/decisions/                # ADR-001 to ADR-010 (Architectural Decision Records)
│   ├── ethics/
│   │   └── governance-axioms.json             # Normative Axiom Set (AX-001 to AX-016)
│   └── gcp-proposal/                          # Google for Startups Cloud Program Technical Dossier
│       ├── 01_architecture_overview.md
│       ├── 02_panamerican_security_and_compliance.md
│       ├── 03_edge_and_gentle_ai.md
│       └── 04_multiregion_dr_and_failover.md
```

### Quick Verification (Hermetic Local Check)

Verify the zero-dependency deterministic engine and run the automated Fail-Close unit test suite:

```bash
# 1. Verify zero external dependencies in go.mod
cat go.mod

# 2. Execute deterministic Fail-Close unit and RDD test suites
CGO_ENABLED=0 go test -v ./...

# 3. Run the end-to-end governance demo (AI Disarmament + Human Refutation)
go run cmd/demo-governance/main.go

# 4. Run Zero-Trust secret scanning check
gitleaks detect --source . -v
```

---

## 7. Founding & Executive Leadership Team

Multidisciplinary leadership combining software engineering, clinical medicine, project governance, healthcare law, and pharmaceutical logistics:

- **Walter Ruben Bartolomé (`@Barto912`)** — Founder, Product Owner & Lead Ethical/Systems Architect.
- **Leonor Figueroa** — Co-Founder, Technical Project Manager & Strategic Compliance Advisor (Documentation Audit, Repository Compliance & Canning Pilot Operations).
- **Dr. Daniel Rojo, MD** — Co-Founder, Chief Medical Information Officer (CMIO) & Medical Director (SNOMED CT Clinical Taxonomies, FHIR R4 Extensions & Home Nursing Protocols).
- **Legal & Regulatory Affairs Division** — Healthcare Law, Data Privacy (Ley 25.326 / HIPAA) & Corporate Structuring (Advisory & Fractional Counsel).
- **Business Development & Operations Division** — Unit Economics, B2B/Public-Private Health Partnerships (Advisory & Fractional Lead).
- **Pharmaceutical Logistics & Biosecurity Division** — Cold-Chain Compliance (+2°C to +8°C), Green Reverse Logistics & ANMAT Traceability (Specialist Pharmacists).

*Institutional Contact:* `contact@nursehealth.ai`

---

## 8. License & Intellectual Property

**Hybrid Licensing Model:**

- **Proprietary (All Rights Reserved):** Core Go engines, cryptographic protocols, and AI distillation pipelines remain the exclusive intellectual property of NurseHealth.ai.
- **Open / Academic Review (Apache 2.0 / CC-BY-NC-SA 4.0):** FHIR/SNOMED adapters, ADR documentation, and regulatory compliance schemas are published to promote interoperability, academic review, and institutional audit.

Foundational medical language models operate under the Google Health AI Developer Foundations terms of use.
