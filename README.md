# NurseHealth.ai | Pan-American Precision Home Healthcare & Green Biosecurity Platform

**Precision Home Nursing, Hybrid Edge-to-Cloud AI, Cryptographic Traceability & Green Pharmaceutical Reverse Logistics.**

🌐 **Read this documentation in other languages:**
- [🇪🇸 Español (Operaciones Hispanoamérica, ANMAT & Economía del Conocimiento)](./README_ES.md)
- [🇧🇷 Português (Operações Brasil / ANVISA & SUS)](./README_PT.md)

> 🚧 **Project Status: RESEARCH, R&D & ARCHITECTURAL VALIDATION (`v0.4.0-research-proposal`).**  
> This repository contains the foundational R&D architecture, deterministic Go core (`stdlib-only`), and sovereign AI distillation pipelines of NurseHealth.ai. Live clinical deployment is strictly gated behind institutional validation protocols.

---

## 1. Executive Summary & Pan-American Health Security Doctrine

Thousands of chronic and post-acute patients across the Americas receive home healthcare without verifiable traceability. Today, extra-hospital care remains fragmented: paper records, compromised cold chains, expired medications accumulating inside homes, and clinical decisions lacking cryptographic auditability.

Framed within the **Pan-American Health Security Doctrine**, **NurseHealth.ai** transforms decentralized home care into an auditable, zero-trust biosecurity and green reverse-logistics grid built natively on the **Google Cloud Platform (GCP)** and **Google AI (Gemma / MedGemma)** ecosystem.

- **Initial Clinical & Logistics Pilot:** 100-case cohort across 57 gated communities in the **Canning corridor (Esteban Echeverría / Ezeiza), Buenos Aires Province, Argentina**.
- **Continental Scaling Roadmap:** Phased expansion across **Uruguay, Chile, Colombia, Mexico, Brazil (SUS/RNDS & ANVISA), and the United States** under a unified HL7 FHIR R4, SNOMED CT, and SHA-256 cryptographic protocol.
- **Platform Thesis:** *Judgment is human; traceability belongs to the system.* NurseHealth.ai elevates autonomous nursing into a verifiable profession while preventing domestic pharmaceutical pollution through eco-routed reverse logistics.

---

## 2. Core Capabilities & Green Biosecurity Solution

1. **Hybrid Edge-to-Workstation-to-Cloud AI (Gemma Ecosystem):**
   * **In-Home Edge Triage:** Real-time medication verification, wound documentation, and clinical triage on the nurse's **Google Pixel 11 Pro** without exposing raw patient PII/PHI to external networks.
   * **Sovereign Local R&D Distillation Node:** On-premise **MIL-STD-810H Certified Workstation** dedicated to local dataset sanitization, QLoRA fine-tuning, and INT8 quantization of **Google Gemma / MedGemma** models under strict national data sovereignty laws.
   * **Cloud Orchestration:** Multimodal clinical reasoning via **MedGemma Maestro (7B/27B) on GCP Vertex AI** governed by dynamic clinical guardrails.
2. **Green Pharmaceutical Reverse Logistics & Cold-Chain IoT:**
   * Home collection of expired, unused, or thermo-sensitive medications (+2°C to +8°C) using IoT thermal containers (Bluetooth 6 / Thread).
   * Carbon-optimized collection routes powered by **Google Maps Platform (Routes API Eco-Friendly Routing & Fleet Engine)**, ensuring full chain-of-custody up to certified bio-safe destruction.
3. **Deterministic Fail-Close Engine (`100% Go stdlib-only`):**
   * Zero external dependencies in the clinical decision core. Upon any logical discrepancy, SNOMED CT taxonomy mismatch, or operator fatigue alert, the engine mechanically blocks execution, emits a signed rejection receipt, and escalates to a human medical authority.
4. **Immutable Cryptographic Audit Ledger:**
   * Every clinical intervention, cognitive gate check, and pharmaceutical custody transfer is sealed with **SHA-256** in an append-only event store.

---

## 3. Three-Tier Technical Architecture & Hybrid AI Flow

To satisfy both **Google Cloud AI scalability** and **Argentine / Pan-American Data Sovereignty mandates**, NurseHealth.ai operates across three physical tiers:

1. **Tier 1 — Mobile Edge Node (Google Pixel 11 Pro / Android 17):**
   * Runs quantized **MedGemma Small (2B/4B INT8)** via LiteRT on the **Google Tensor G6 NPU** (`< 300ms` latency).
   * Hardware-backed **Titan M3** security enclave signs local `VerificationReceipt` payloads using SHA-256.
   * **Local Context Store:** Edge-first SQLite + FTS5 context store that strips and pseudonymizes PHI at the edge prior to transmission.
2. **Tier 2 — Sovereign R&D & Air-Gapped Distillation Node (MIL-STD-810H Certified Workstation):**
   * Dedicated on-premise workstation (**Intel Xeon W ECC / NVIDIA RTX PRO, MIL-STD-810G/H certified**) located in Buenos Aires Province, Argentina.
   * Executes local model distillation (QLoRA / weight quantization), synthetic clinical data generation, and deterministic Go (`stdlib-only`) binary compilation and verification before edge deployment.
3. **Tier 3 — Cloud Orchestration & Multi-Region Ledger (Google Cloud Platform):**
   * **Google Cloud Run:** 100% Go microservices, REST API, and Event Sourcing ingestion.
   * **Vertex AI:** MedGemma Maestro (7B/27B) with dynamic guardrails.
   * **Cloud SQL Enterprise Plus:** PostgreSQL Append-Only SHA-256 Event Store.
   * **Cloud Healthcare API:** HL7 FHIR R4 & SNOMED CT isolated behind a Domain-Driven Design Anti-Corruption Layer (ACL).
   * **Google Maps Platform:** Routes API & Fleet Engine for low-emission home nursing dispatch and reverse pharmaceutical collection.

| Architectural Layer | Technology & Standards |
| :--- | :--- |
| **Backend Core** | 100% Go (`stdlib-only` core) on **Google Cloud Run** (Event Sourcing & Fail-Close Engine) |
| **On-Premise R&D Node** | **MIL-STD-810H Certified Workstation** (Intel Xeon W ECC + NVIDIA RTX PRO on Ubuntu Linux) |
| **Mobile Edge Node** | **Google Pixel 11 Pro** (Tensor G6 NPU + Titan M3) & Quantized MedGemma (2B/4B INT8) |
| **Cloud AI (GCP)** | **MedGemma Maestro (7B/27B)** on **Vertex AI** with dynamic clinical guardrails |
| **Green Logistics & Geo** | **Google Maps Platform** (Routes API Eco-Routing & Fleet Engine) + BLE 6/Thread Cold-Chain IoT |
| **Audit Ledger & DB** | **Cloud SQL Enterprise Plus (PostgreSQL)** Append-Only SHA-256 Chain-of-Custody Store |
| **Clinical Interoperability** | **HL7 FHIR R4 & SNOMED CT** via **Google Cloud Healthcare API** (DDD Anti-Corruption Layer) |

---

## 4. Philosophy of Care: The Humanistic Firewall (ADR-006 & ADR-010)

Inspired by humanist medical ethics (*Magnifica Humanitas* — zero iatrogenesis, human dignity over automation), the **Human Firewall** enforces cognitive integrity before any clinical action is permitted.

These axioms (mandatory mental clarity, workload modulation, and directed attention for SNOMED CT signing) are compiled as strict **Go domain invariants** inside the governance core (`pkg/governance/`), backed by automated unit tests (`go test ./...`). The system enforces **AI Disarmament (AX-012)**: the AI can suggest, but only a credentialed human can authorize irreversible clinical actions.

---

## 5. Multi-Region Resilience & Disaster Recovery

NurseHealth.ai guarantees clinical continuity under the **Doctrine of Health Security of the Americas**:
- **Layer 0 (Intra-Regional HA):** Synchronous multi-zone high availability across `southamerica-west1` (Santiago) and `us-east1` (Northern Virginia), achieving **RPO = 0** and **RTO < 60s**.
- **Immutable Fail-Close:** Infrastructure failovers never bypass cognitive gates or patient safety invariants.
- **Sovereign Continuity:** Architected to prevent single-region cloud outages from ever disrupting active home care.

---

## 6. Multi-Jurisdictional Regulatory & R&D Alignment

Engineered from day one to satisfy hemispheric biosecurity, tax-incentive R&D frameworks, and health regulations:

- 🇦🇷 **Argentina (Pilot Jurisdiction & R&D Hub):**
  * **Data Privacy & Health Law:** Ley 25.326 (Protección de Datos Personales), Ley 26.529 (Derechos del Paciente), and Ley 27.553 (Receta Electrónica y Teleasistencia).
  * **Clinical & Pharmaceutical Biosecurity:** Aligned with **ANMAT** Software as a Medical Device (SaMD) guidelines, **Ley Nacional 24.051** (Residuos Peligrosos / Farmacéuticos), and **Ley Provincial PBA 11.347** (Residuos Patogénicos) for reverse logistics traceability.
  * **Knowledge Economy R&D Framework (Ley 27.506 - LEC):** Structured under ARCA CLAE Codes **620100** (Software Architecture) and **721000** (Experimental R&D in Engineering and Health Sciences) via our Argentine SAS operating entity.
- 🇺🇸 **United States (Corporate & Scaling Hub):** HIPAA, FDA SaMD (Good Machine Learning Practice), SOC 2 Type II, and NDAA Sec. 889 compliance via a Delaware C-Corp parent structure.
- 🇧🇷 **Brazil:** LGPD, ANVISA (RDC Nº 657/2022 for SaMD & RDC Nº 430/2020 for pharmaceutical cold chain), and native interoperability readiness for **SUS (RNDS)**.

---

## 7. Repository Structure, ADRs & GCP Research Proposal

- **Architecture Decision Records (ADRs):** Consolidated ADR-001 through ADR-010 (including Human Firewall, FHIR R4 Extensions, SHA-256 Audit Ledger, and Disarm AI Architecture).
- **Continuous Integration:** 100% green CI/CD pipelines across Go unit/RDD test suites, Markdown link validation, ADR linting, and automated `gitleaks` secret scanning.
- 📂 **Detailed GCP Research Proposal (`/docs/gcp-proposal/`):**
  1. [`01_architecture_overview.md`](./docs/gcp-proposal/01_architecture_overview.md) — End-to-End Hybrid Edge/Workstation/Cloud Topology.
  2. [`02_panamerican_security_and_compliance.md`](./docs/gcp-proposal/02_panamerican_security_and_compliance.md) — HIPAA, ANMAT, ANVISA & Zero-Trust Governance.
  3. [`03_edge_and_gentle_ai.md`](./docs/gcp-proposal/03_edge_and_gentle_ai.md) — MedGemma Quantization, Pixel 11 Pro NPU & Gentle-AI RDD/ODD.
  4. [`04_multiregion_dr_and_failover.md`](./docs/gcp-proposal/04_multiregion_dr_and_failover.md) — Multi-Region Failover (`southamerica-west1` / `us-east1`).

---

## 8. Founding & Executive Team

Multidisciplinary leadership combining software engineering, clinical medicine, project governance, healthcare law, and pharmaceutical logistics:

- **Walter Bartolomé (`@Barto912`)** — Founder, Product Owner & Lead Ethical/Systems Architect.
- **Leonor Figueroa** — Technical Project Manager & Strategic Advisor (Documentation Audit, Repository Compliance & Canning Pilot Operations).
- **Dr. Daniel Rojo** — Chief Medical Information Officer (CMIO) & Medical Director (SNOMED CT Clinical Taxonomies, FHIR R4 Extensions & Home Nursing Protocols).
- **Legal & Regulatory Affairs Division** — Healthcare Law, Data Privacy (Ley 25.326 / HIPAA) & Corporate Structuring (Advisory & Fractional Counsel).
- **Business Development & Operations Division** — Unit Economics, B2B/Public-Private Health Partnerships (Advisory & Fractional Lead).
- **Pharmaceutical Logistics & Biosecurity Division** — Cold-Chain Compliance (+2°C to +8°C), Green Reverse Logistics & ANMAT Traceability (Specialist Pharmacists).

*Institutional Contact:* `contact@nursehealth.ai`

---

## 9. License & Intellectual Property

**Hybrid Licensing Model:**
- **Proprietary (All Rights Reserved):** Core Go engines, cryptographic protocols, and AI distillation pipelines remain the exclusive intellectual property of NurseHealth.ai.
- **Open / Academic Review (Apache 2.0 / CC-BY-NC-SA 4.0):** FHIR/SNOMED adapters, ADR documentation, and regulatory compliance schemas are published to promote interoperability, academic review, and institutional audit.
