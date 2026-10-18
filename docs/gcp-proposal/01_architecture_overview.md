# Architecture Overview: Master/Student AI Topology & Zero-Trust Governance

## 1. System Paradigm
NurseHealth.ai operates on a decentralized, hybrid Edge-to-Cloud architecture natively built on Google Cloud Platform (GCP). The system pairs a lightweight Edge AI ("Student") with a robust Cloud AI ("Master") to ensure continuous clinical triage without compromising patient data sovereignty, even in offline environments.

## 2. The Master/Student MedGemma Topology
* **The Student (Edge AI):** Google Pixel 11 Pro running quantized MedGemma Small (2B/4B) via TensorFlow Lite directly on the Tensor NPU. It handles immediate, offline-first tasks such as wound image classification, medication verification, and initial vital sign parsing.
* **The Master (Cloud AI):** GCP Vertex AI running MedGemma Maestro (7B/27B). It provides complex asynchronous clinical analysis, longitudinal patient context, and acts as the ultimate validator for the Student's preliminary findings.
* **Google Harness² Orchestration:** Connects Master and Student using a *Propose-Probe-Compose* paradigm, enforcing deterministic guardrails that physically prevent clinical hallucinations.

## 3. Zero-Trust Contribution & Execution Engine
To maintain enterprise-grade security and allow scalable collaboration without exposing the platform to unknown or untrusted third-party code, NurseHealth.ai enforces a strict Zero-Trust governance model:
* **100% Go (Golang) Microservices:** Ultra-low latency, memory-safe execution ensuring isolated process boundaries.
* **Fail-Close Determinism:** The engine is programmed to block execution and escalate to human medical authorities upon any logical ambiguity.
* **Immutable Event Sourcing:** All state changes are appended to a PostgreSQL Cloud SQL database with SHA-256 cryptographic hashing.
* **Strict PR Governance:** The repository architecture isolates all CI/CD secrets. External contributions cannot alter core clinical logic or access production pipelines without multi-signature cryptographic approval from the core team.

## 4. Clinical Interoperability (DDD ACL)
All communications with external Electronic Health Records (EHR) and national health directories occur via the **Google Cloud Healthcare API** (native HL7 FHIR R4 & SNOMED CT semantic alignment). To protect the core logic from external data corruption, this integration is strictly shielded by a Domain-Driven Design (DDD) Anti-Corruption Layer (ACL).
