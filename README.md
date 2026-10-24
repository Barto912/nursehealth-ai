# NurseHealth.ai | Pan-American Precision Home Healthcare & Biosecurity Platform

**Precision Home Nursing with Edge AI & Cryptographic Traceability.**

[![CI](https://github.com/Barto912/nursehealth-ai/actions/workflows/main.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Docs](https://github.com/Barto912/nursehealth-ai/actions/workflows/auditoria-docs.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Seguridad](https://github.com/Barto912/nursehealth-ai/actions/workflows/seguridad-gitleaks.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)

🌐 **Read this documentation in other languages:**
- [🇪🇸 Español (Operaciones Hispanoamérica)](./README_ES.md)
- [🇧🇷 Português (Operações Brasil / ANVISA & SUS)](./README_PT.md)

> 🚧 **Project Status: RESEARCH & ARCHITECTURAL DESIGN.**
> This repository is not a production clinical product yet. The code and documentation presented here are experimental and their use in live clinical environments is neither authorized nor validated at this stage.

---

## 🛡️ The Problem & Pan-American Health Security Doctrine

Thousands of chronic patients across the Americas receive home healthcare without real traceability. Today, this critical work remains invisible: paper spreadsheets, expired medications inside homes, broken chains of custody, and clinical decisions lacking auditable evidence.

Framed within the **Pan-American Health Security Doctrine**, **NurseHealth.ai** transforms decentralized home care into an auditable, zero-trust biosecurity grid built natively on the **Google Cloud Platform (GCP)** ecosystem.

* **Initial Pilot:** 57 gated communities (*countries y barrios cerrados*) in Canning, Buenos Aires, Argentina.
* **Continental Vision:** Regional scaling across Chile, Uruguay, Colombia, Mexico, Brazil (SUS/ANVISA), and the United States under a unified cryptographic and clinical protocol.

---

## 💡 The Solution

An extra-hospital precision nursing and pharmaceutical reverse logistics platform:

- **Hybrid Edge-to-Cloud AI:** Medication identification and clinical triage directly on the nurse's device (**Google Pixel 11 Pro**), without exposing raw patient PII to the cloud.
- **Cryptographic Chain of Custody:** Every clinical event and sample movement is sealed with **SHA-256** in an append-only Event Store.
- **Pharmaceutical Reverse Logistics:** Home collection of expired or unused medications with end-to-end traceability up to final bio-safe disposal.
- **Fail-Close by Design:** When in doubt, the deterministic Go engine blocks execution and escalates to a human medical authority. It never improvises.

> **Motto:** *Judgment is human; traceability belongs to the system.*

**Platform Thesis:** AI-supported precision nursing that elevates skilled autonomous work into a verifiable professional career—enabling continuous education, auditable reputation, and economic dignity for caregivers across the Americas.

---

## 🏗️ Technical Architecture & Hybrid AI Flow

1. **Edge Node — Google Pixel 11 Pro:** Quantized MedGemma (2B/4B) on Tensor NPU + TFLite + Engram (SQLite/FTS5 local context store).
2. **Zero-Trust Transport:** Pseudonymized, cryptographically signed clinical event stream.
3. **Core Ingestion — Google Cloud Run:** 100% Go (Golang) Microservices, REST API, Deterministic Fail-Close Engine, and Event Sourcing.
4. **Cloud AI Orchestration — Vertex AI:** MedGemma Maestro (7B/27B) + Google Harness² (*Propose-Probe-Compose* dynamic guardrails).
5. **Immutable Audit Ledger — Cloud SQL:** PostgreSQL Append-Only SHA-256 Event Store.
6. **Clinical Interoperability — Cloud Healthcare API:** HL7 FHIR R4 & SNOMED CT isolated behind a Domain-Driven Design Anti-Corruption Layer (ACL).

### Technical Stack

| Layer | Technology |
| :--- | :--- |
| **Backend Core** | 100% Go (Golang) on Google Cloud Run (Event Sourcing & Fail-Close Engine) |
| **Frontend Mobile** | React Native with Expo |
| **Edge AI (Device)** | Quantized MedGemma Small (2B/4B) & TensorFlow Lite on Google Pixel 11 Pro (Tensor NPU) |
| **Cloud AI (GCP)** | MedGemma Maestro (7B/27B) on Vertex AI + Google Harness² guardrails |
| **AI Governance & Memory** | Gentle-AI (ODD/RDD) + Engram (Local-first SQLite + FTS5 context store) |
| **Database & Ledger** | PostgreSQL on Cloud SQL (Append-only SHA-256 Event Store) |
| **Interoperability** | HL7 FHIR R4 & SNOMED CT via Google Cloud Healthcare API (behind DDD ACL) |
| **Infrastructure** | Google Cloud Platform (GCP) |

---

## ⚖️ Multi-Jurisdictional Regulatory Alignment

Engineered from day one to satisfy hemispheric health security and data sovereignty standards:
* 🇺🇸 **United States:** HIPAA, SOC 2 Type II, CISA security alignment, and NDAA compliance.
* 🇧🇷 **Brazil:** LGPD, ANVISA (SaMD guidelines), and native interoperability readiness for **SUS (RNDS)**.
* 🌎 **Hispanic America:** Local data privacy frameworks (Argentina Ley 25.326, Chile, Uruguay, Colombia, Mexico) and national Electronic Health Record standards.

---

## 📈 Current Status & Milestones

Architectural design completed and audited:
- **Event Storming (MVP-0):** Completed with 14 critical clinical/operational hotspots resolved.
- **Aggregate Boundaries v0.3:** Repair vs. acceptance patterns formally documented.
- **DDD Context Map:** Anti-Corruption Layer (ACL) for FHIR R4 and regulatory Reference Model established.
- **Continuous Integration:** Green CI/CD pipelines, automated documentation auditing, and `gitleaks` secret scanning on every commit.

**In Progress:** MVP implementation and cloud AI orchestration applying for the **Google for Startups Cloud Program**.

📂 **Detailed GCP Research Proposal:** Technical deep-dives are located in [`/docs/gcp-proposal/`](./docs/gcp-proposal/):
1. [`01_architecture_overview.md`](./docs/gcp-proposal/01_architecture_overview.md)
2. [`02_panamerican_security_and_compliance.md`](./docs/gcp-proposal/02_panamerican_security_and_compliance.md)
3. [`03_edge_and_gentle_ai.md`](./docs/gcp-proposal/03_edge_and_gentle_ai.md)

---

## 🧭 Ethical Governance

Non-negotiable principles from the project manifesto:

1. **Dignity over data.**
2. **Autonomy over automation.**
3. **Transparency over black boxes.**
4. **Care over efficiency.**
5. **Presence over connection.**
6. **Prevention over reaction.**
7. **Education over indoctrination.**

Patient data is processed strictly pseudonymized. No clinical decision is ever automated.

---

## 🧠 Philosophy of Care: The Human Firewall

Precision nursing requires more than technical skill—it demands cognitive integrity. The **Human Firewall** (ADR-006) codifies three foundational axioms that gate every clinical event:

- **AX-001 (Inertia and Noise):** Mental clarity as precondition. If `MentalClarity < 7`, the system blocks event creation and logs the attempt.
- **AX-002 (Tones / Work Level):** Energetic tone modulates clinical efficacy. If `WorkLevel == LOW`, the system suggests reassignment or pause.
- **AX-003 (Simple Attention):** Directed attention as condition for SNOMED signature. If `SimpleAttention == false`, the operator must record an `AttentionalFinding` before signing.

These axioms are not aspirational—they are **domain invariants** enforced by the `RuidoCheck` Value Object with unit tests. A professional in unfit cognitive state cannot create or sign a clinical event: the system blocks and records.

**Implementation:**

- [Foundational Axioms](docs/philosophy/axioms.json) (AX-001 to AX-003)
- [RuidoCheck Value Object](pkg/cognitive/ruido_check.go) (typed gates with tests)
- [ADR-006: Human Firewall (MVP)](docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md)
- [Context Map v0.2](docs/domain/context-map.md) (invariants I10, I11)

---
## 👥 Core Team

- **Barto** — Founder & Product Owner.
- **Leonor** — Strategic Advisor.

---

## 📄 License

Proprietary. All rights reserved.
## 🛡️ Multi-Layer Resilience and Data Sovereignty Architecture

NurseHealth.ai operates under the **Doctrine of Health Security of the Americas**, ensuring clinical continuity and regulatory compliance at all times.

- **Layer 0 (Intra-Regional):** Our first line of defense. We operate with synchronous multi-zone high availability in `southamerica-west1` and `us-east1`, guaranteeing RPO=0 and RTO<60s. This design ensures strict adherence to **HIPAA**, **FDA SaMD** guidelines, and **SOC 2 Type II** standards.
- **Immutable Fail-Close (Invariants I10/I11):** No infrastructure failure enables automatic retries if the operator's cognitive gates (`RuidoCheck`, AX-001 to AX-003) or patient gates (AX-004 to AX-006) have blocked the operation. Clinical vetoes are deterministic and recorded immutably (append-only) in the audit ledger.
- **Operational Motivation:** Architected in response to the public `europe-west2` incident (October 2026). Our Pan-American topology operates sovereignly, with zero dependency on external regions, ensuring precision nursing tools never fail during critical care.

*For technical details on cross-region replication and failover mechanisms, refer to `docs/gcp-proposal/04_multiregion_dr_and_failover.md`.*
