# 04. Multi-Layer Resilience, Intra-Regional HA (Layer 0) and Sovereign Fail-Close Architecture

- **Project:** [Barto912/nursehealth-ai](https://github.com/Barto912/nursehealth-ai)
- **Verified baseline:** `module github.com/Barto912/nursehealth-ai` (`go 1.23`, stdlib-only, zero external dependencies in `go.mod`) | 81 commits at PR #20 merge | tag `v0.4.0-research-proposal`
- **Doctrine framework:** Doctrina de Seguridad Sanitaria de las Américas (EN, ES, PT), application to the Google for Startups Cloud Program
- **Ecosystem entities:** `MissionHealth.ai` (Delaware C-Corp, global matrix) and `NurseCanning.health` (ACE Argentina, Joint Venture Clínica Canning Health)
- **Confirmed database (Phase 1):** [Cloud SQL Enterprise Plus for PostgreSQL with Cross-Region DR Replica](https://docs.cloud.google.com/sql/docs/postgres/replication/cross-region-replicas) between `southamerica-west1` and `southamerica-east1`, to optimize Google for Startups Cloud Program credits
- **Operational motivation:** public incident in `europe-west2` (London, October 2026) that affected payment systems in Europe; NurseHealth.ai had zero exposure in that region. The event motivated the formalization of Layer 0 and the containment cascade (section 7)
- **Consensus record:** Argon (architectural gemelo) verdicts Q1-Q6 approved with corrections; Deivis (technical gemelo) evidence audit against `go.mod`, merged sources and CI (rules P-12, P-14, P-15)
- **Resilience targets:** RPO = 0 seconds at the sovereign edge (Pixel 11 Pro with Titan M3 signed Event Sourcing), RTO < 15 seconds for cross-region failover (Layer 0 intra-regional + Layer 1 inter-regional cascade with deterministic Fail-Close).

## 1. Verified repository baseline and zero-friction guarantee

The design integrates with the real state of the repository (no assumed files):

| Real component in the repository | Risk before single-region SPOF | Zero-friction integration strategy |
| --- | --- | --- |
| `pkg/cognitive/ruido_check.go` and ADR-006 (Human Firewall, AX-001 to AX-003) | A regional failover that retries a clinical transaction without operator state would violate invariants I10/I11 | `RuidoCheck` is evaluated before any multi-region router; no cloud switch skips the human firewall (invariant I10) |
| `pkg/triage` real contracts: `ValidateCareEvent`, `OperatorGate`, `AuditRepository`, `CognitiveBlockReason`, `PatientContextCheck` (AX-004 to AX-006), `middleware_test.go` | Wrapping non-existent files would break the build and committee trust | `pkg/triage/resilient_evaluator.go` will be born in PASO 3 as a self-contained deep module (Argon Q4 decision), with own types, wrapping real contracts; today `evaluator.go` does not exist and is not referenced |
| `db/migrations/001_create_cognitive_audit_log.sql` (append-only, SHA-256, decision BLOCKED/ESCALATED/PASSED) | Audit lost during regional promotion would break clinical non-repudiation | Every failover event and clinical veto is recorded append-only with region and authority; no UPDATE/DELETE in schema |
| CI pipeline (Go tests + RDD, markdown/ADR, broken links, gitleaks) | Markdown format errors or heavy dependencies would turn checks red in front of reviewers | This document uses tables compatible with `.markdownlint-cli2.jsonc`, public links only, and reaffirms stdlib-only in `go.mod` |

## 2. Layer 0: intra-regional resilience (zone to region)

The first defense line never moves data across jurisdictions and is the cheapest and fastest:

| Level | Scope | Mechanism | RPO | RTO | Escalation rule |
| --- | --- | --- | --- | --- | --- |
| L0 zonal | Zone within `southamerica-west1` | Cloud SQL Enterprise Plus with synchronous multi-zone regional HA; Cloud Run with multi-zone redundancy and `min-instances = 2`; Global ALB with Outlier Detection ejecting sick instances | 0 (synchronous replica) | 60-120 s (SQL failover); compute in seconds | Escalate to L1 only if the whole region fails |
| L1 regional | Full region `southamerica-west1` | Promotion of Cross-Region DR Replica in `southamerica-east1` plus ALB switch | Under 15-60 s WAL (effectively 0 via Engram Outbox at the edge) | 2-3 min | Edge covers the gap with RPO 0 |
| L2 sovereign edge | Patient home (Canning) | Local-first mode on Pixel 11 Pro: SQLite/FTS5 plus Outbox CMEK plus Merkle root; events signed with Titan M3 queued until reconnection | 0 (on device) | Under 300 ms | Escalate to L3 if edge also fails |
| L3 veto | Everything down | Rule `FAIL_CLOSE_L0_EXHAUSTED`: never improvise; refuse clinical transactions without traceability | Not applicable | Not applicable | End of cascade |

Cascade invariants (Context Map v0.2, reinforced here per Argon Q3):

- I10 (pre-routing gate): `cognitive.RuidoCheck` (AX-001 to AX-003) and `PatientContextCheck` (AX-004 to AX-006) run before any routing or regional failover attempt.
- I11 (no retry on clinical veto): no infrastructure failure enables retries once `OperatorGate` or `PatientContextCheck` returned a clinical block; a clinical `BlockError` (`ReasonMentalClarityTooLow`, `ReasonWorkLevelTooLow`, `ReasonAttentionNotFocused`, patient-side reasons) is deterministic, is recorded append-only in `cognitive_audit_log`, and closes the gate (immutable Fail-Close).

## 3. Federated and sovereign multi-region topology (Shield of the Americas)

Two geopolitical rings on [official Google Cloud regions](https://docs.cloud.google.com/compute/docs/regions-zones), articulated with the sovereign edge node, comply with data sovereignty laws per [DR architecture for locality-restricted workloads](https://docs.cloud.google.com/architecture/architecting-disaster-recovery-for-locality-restricted-workloads):

```text
===============================================================================
               SOVEREIGN EDGE LAYER (LOCAL-FIRST / ZERO-TRUST EDGE)
 [Pixel 11 Pro: Tensor G6 + Titan M3 PQC + NurseGemma INT8 (Stanford SwiLA)]
 [Human Firewall (RuidoCheck AX-001..003) + Gentle Engram (SQLite/FTS5 + PHI
                 filter + Outbox CMEK + Merkle root)]
===============================================================================
                                     |
                  Zero-Trust transport: BeyondCorp Enterprise + mTLS
                                     v
===============================================================================
                 GLOBAL ANYCAST ENTRY PLANE (NO REGIONAL SPOF)
   [Global External Application Load Balancer (single Anycast IP) + WAF Armor]
   [Cloud Run Service Health + Outlier Detection (switch under 15 s)]
===============================================================================
                     /                             \
                    v                               v
+-------------------------------------+   +-------------------------------------+
| SOUTH RING (Southern Cone & Brazil) |   | NORTH RING (US, Mexico & Colombia)  |
| southamerica-west1 (Chile) primary  |   | us-east1 (South Carolina) primary   |
| - Cloud Run Go Core + ACL + MCP     |   | - Cloud Run Go Core + ACL + MCP     |
| - Cloud SQL Enterprise+ (ledger)    |   | - Vertex AI MedGemma-27B + Harness2 |
| - FHIR R4 Store (primary AR)        |   | - FHIR R4 Store (HIPAA US)          |
|               <==>                  |   |               <==>                  |
| southamerica-east1 (Sao Paulo) peer |   | us-central1 (Iowa) peer             |
| - Cloud Run Go Core (active)        |   | - Cloud Run Go Core (active)        |
| - Cloud SQL Cross-Region DR Replica |   | - Vertex AI (failover)              |
| - FHIR R4 Store (sovereign BR)      |   | - BigQuery Managed DR (RWD)         |
+-------------------------------------+   +-------------------------------------+
```

### 3.1 Trilingual regional assignment and legal sovereignty

| Ring / documentary facade | GCP region | Physical location | Regulatory framework covered | Operational and failover role |
| --- | --- | --- | --- | --- |
| Southern Cone primary (`README_ES.md`) | `southamerica-west1` | Santiago, Chile | Argentina (Ley 25.326 art. 2, Ley 26.529, Ley 12.245 PBA, ANMAT), Chile, Uruguay | Primary ingestion of the Canning Health pilot; Cloud SQL Enterprise Plus primary instance |
| South America peer and Brazil sovereignty (`README_PT.md`) | `southamerica-east1` | Osasco, São Paulo, Brazil | Brazil (ANVISA RDC 657/2022 SaMD, LGPD/ANPD, RNDS/SUS) | Active-active Go compute; mandatory residence for Brazilian clinical data; Cross-Region DR Replica for Southern Cone |
| North primary and Andean region (`README.md`) | `us-east1` | Moncks Corner, South Carolina, US | US (HIPAA, FDA SaMD, SOC 2 Type II), Colombia (INVIMA, Ley 1581) | Primary Vertex AI inference node (MedGemma-27B + Harness2); transcontinental backup |
| North peer and Mexico (`README.md`) | `us-central1` | Council Bluffs, Iowa, US | US (HIPAA), Mexico (COFEPRIS, LFPDPPP) | Vertex AI failover; RWD analytics with differential privacy (epsilon = 0.5) |

## 4. Confirmed data replication strategy (credit optimization)

| Data layer | Service and configuration | Inter-region replication | RPO | RTO | Justification for the committee |
| --- | --- | --- | --- | --- | --- |
| 1. Edge clinical memory (`pkg/engram`, pilot) | SQLite + FTS5 with CMEK on Pixel 11 Pro | Local Write-Ahead Outbox plus Merkle root until cloud ACK | 0 s (on device) | Under 300 ms local | Continuous clinical operation without internet; effective end-to-end RPO 0 on resend |
| 2. Event bus | [Cloud Pub/Sub](https://docs.cloud.google.com/pubsub/docs/reliability-intro) with [message storage policy](https://docs.cloud.google.com/pubsub/docs/resource-location-restriction) | Synchronous replication across zones in `southamerica-west1` and `southamerica-east1` with CMEK | 0 s (acknowledged messages) | Under 5 s | Decouples Cloud Run Go services from the database during promotion |
| 3. SHA-256 ledger and ATDP | [Cloud SQL Enterprise Plus](https://docs.cloud.google.com/sql/docs/postgres/intro-to-cloud-sql-disaster-recovery) primary `southamerica-west1` with [DR Replica](https://docs.cloud.google.com/sql/docs/postgres/replication/cross-region-replicas) in `southamerica-east1` | Regional synchronous HA plus continuous asynchronous WAL streaming | Effectively 0 s via Engram Outbox rehydration (15-60 s pure WAL) | 2-3 min (automated promotion) | Reduces monthly credit consumption by about 65-75% versus Spanner in MVP (100 to 10,000 pickups per month), 100% PostgreSQL compatible |
| 4. HL7 FHIR R4 and SNOMED CT interoperability | [Cloud Healthcare API](https://docs.cloud.google.com/healthcare-api/docs/regions) dual store plus [Cloud Storage Turbo Replication](https://docs.cloud.google.com/storage/docs/managing-turbo-replication) | Dual-Write from Go ACL via Pub/Sub to both FHIR stores plus dual-region bucket with [import/export](https://docs.cloud.google.com/healthcare-api/docs/how-tos/fhir-import-export) and [21-day PITR](https://docs.cloud.google.com/healthcare-api/docs/fhir-pitr) | Under 5 s (Pub/Sub) | Under 30 s (ACL switch) | Overcomes the regional limitation of Healthcare API without duplicating licensing |
| 5. IoT GS1 telemetry and RWD | [BigQuery Managed Disaster Recovery](https://docs.cloud.google.com/bigquery/docs/managed-disaster-recovery) | Inter-region Turbo Replication of pharmacovigilance datasets (epsilon = 0.5) | Under 15 min | Under 1 min | Protects the commercial asset for laboratories against regional contingencies |

Phase 2 scalability note: above 50,000 monthly multi-country events (Series A/B), the Go repository interfaces (`AuditRepository`, `DBExecutor`) allow migrating the PostgreSQL adapter to [Cloud Spanner multi-region](https://docs.cloud.google.com/spanner/docs/create-manage-instances) with [99.999% SLA](https://cloud.google.com/spanner/sla) without touching domain logic.

## 5. Automated failover and deep modules in Go

### 5.1 Network and compute layer: Global ALB plus Cloud Run Service Health (RTO under 15 s)

- Single Anycast entry: [Global External Application Load Balancer](https://docs.cloud.google.com/load-balancing/docs/https/setup-global-ext-https-serverless) with [Google Cloud Armor](https://docs.cloud.google.com/armor/docs/security-policy-overview) and BeyondCorp Enterprise.
- Automated evacuation: [Cloud Run Service Health and Outlier Detection](https://docs.cloud.google.com/run/docs/configuring/configure-service-health) (`consecutiveErrors: 3`, `interval: 5s`, `maxEjectionPercent: 100`) redirects 100% of traffic to serverless NEGs in the peer region in under 15 s, without waiting for DNS propagation, per [Global External ALB failover](https://docs.cloud.google.com/load-balancing/docs/https/global-applb-failover-overview).
- Code-Red pre-warm (`PROTO-003`): `min-instances = 2` in both southern regions so the under-5-minute emergency sequence (UTIM/SAME 107 dispatch and Shockroom alert at Clínica Canning Health) never suffers cold-start latency.
- Multi-region compute reference: [serve traffic from multiple regions in Cloud Run](https://docs.cloud.google.com/run/docs/multiple-regions) and [deployment archetypes](https://docs.cloud.google.com/architecture/deployment-archetypes/multiregional), complemented with [DR for cloud infrastructure outages](https://docs.cloud.google.com/architecture/disaster-recovery).

### 5.2 Application layer: pkg/triage/resilient_evaluator.go (self-contained deep module)

PASO 3 decision (Argon Q4, double-gate P-15): the module is born inside `pkg/triage` with self-contained types, wrapping real contracts (`ValidateCareEvent`, `OperatorGate`, `AuditRepository`, `CognitiveBlockReason`) and `cognitive.RuidoCheck`. Design sketch (not compiled code; implementation will re-verify every symbol against `go.mod` and merged sources):

```go
package triage

// Design sketch for PASO 3. Self-contained types only (no external deps).

type ExecutionRegion string

const (
    RegionSAWest1   ExecutionRegion = "southamerica-west1"
    RegionSAEast1   ExecutionRegion = "southamerica-east1"
    RegionUSEast1   ExecutionRegion = "us-east1"
    RegionEdgePixel ExecutionRegion = "edge-pixel11-npu"
)

type RegionalAuditOutcome struct {
    IsSafe      bool            `json:"is_safe"`
    RiskReason  string          `json:"risk_reason"`
    JudgeID     string          `json:"judge_id"`
    RegionUsed  ExecutionRegion `json:"region_used"`
    AuthorityID string          `json:"authority_id"` // Titan M3 / CMEK URI
    DegradedDR  bool            `json:"degraded_dr"`
    LatencyMs   int64           `json:"latency_ms"`
}

type RegionalSafetyJudge interface {
    Region() ExecutionRegion
    Audit(ctx context.Context, draft CanonicalDraft) (RegionalAuditOutcome, error)
}

// Cascade rule (I10/I11):
//   infrastructure error    -> next region in prelation order;
//   clinical veto (!IsSafe) -> immediate Fail-Close, no cross-region retry;
//   all regions exhausted   -> FAIL_CLOSE_L0_EXHAUSTED veto.
```

Regional cascade behavior:

| Phase | Behavior | Guarantee |
| --- | --- | --- |
| 1. ODD structuring | Orchestrator generates canonical clinical draft; nil orchestrator is Fail-Close | No drafting without authority |
| 2. Blind regional audit | Judges in order west1, east1, us-east1, edge-pixel11-npu with bounded budget per hop (800 ms cloud, 300 ms edge) | Infrastructure timeout switches without false Fail-Close |
| 3. Certification | Canonical SHA-256 of payload plus region plus authority plus degraded mode; append-only record in `cognitive_audit_log` | Clinical and operational non-repudiation with executing region |
| 4. ATDP telemetry | Tuples (o, h, a, y, r, m) to Flight Recorder for Harness2 loop | Evolvable routing without touching frozen clinical weights |

### 5.3 Self-learning improvements (Harness2 plus ATDP plus SEW plus SwiLA)

| Improvement | Foundation | Function in multi-region architecture | Non-negotiable safety barrier |
| --- | --- | --- | --- |
| Predictive routing via ATDP and Harness2 (C4 tool guidance, C8 middleware plugin) | [Harness2: Recursive Agent Harnessing for an Open World (arXiv:2610.14224v1)](https://arxiv.org/html/2610.14224v1) | ATDP captures real latencies and micro-outages per time band and country; Propose-Probe-Compose adjusts `PerRegionTimeout` and regional prelation as signed light config | MedGemma weights and Fail-Close rules frozen; only routing and latency budget evolve |
| Counterfactual replay of regional outages for NurseGemma distillation | SEW trajectory distillation protocol (internal research notebook) | After failover, offline counterfactual query: did local NurseGemma INT8 emit the same SNOMED CT verdict as cloud MedGemma-27B; discrepancies feed the distillation set | Raises offline concordance (over 98-99.5%) so edge equals cloud quality during total outage |
| Constant memory O(J d^2) with Stanford SwiLA during Outbox reconciliation | SwiLA technical spec (internal document): Forget Gate and block-Newtonian iteration | During multi-hour outages, edge processes the full shift (IoT thermal telemetry every 60 s plus clinical events) without saturating 16 GB RAM or thermal throttling | Forget Gate discards irrelevant thermal fluctuations but deterministically retains excursions over 180 s outside +2 to +8 C |
| 7-day late-loop reward tied to region and SSCC container | ATDP late reward signals; PROTO-006-LIF scope | Ledger ties original receipt (receipt_id plus RegionUsed plus SSCC) to 7-day clinical outcome and pharmacy reception record | Deterministic hooks (RuidoCheck AX-001/002/003 and PROTO-001 triggers) out of evolutionary reach |

## 6. Master plan of five atomic steps

| Step | Technical goal | State |
| --- | --- | --- |
| PASO 1 | Branch from current main (Git Safety First) | Completed: `feature/gcp-proposal-04` |
| PASO 2 | This document in `docs/gcp-proposal/` with verified baseline and Layer 0 | This file |
| PASO 3 | `pkg/triage/resilient_evaluator.go` plus tests (regional failover, immediate clinical veto, `FAIL_CLOSE_L0_EXHAUSTED`), under double-gate P-15 | Pending, after this merge |
| PASO 4 | Trilingual sync in `README.md`, `README_ES.md`, `README_PT.md` (Gate 5) | Pending |
| PASO 5 | Final PR plus immutable tag `v0.5.0-gfs-ready` (Gate 8b) | Pending, after Gates 5 and 6 |

## 7. Triggering incident and lessons learned (October 2026)

- Event: public outage in `europe-west2` (London) in October 2026, with press impact on payment and authentication systems in Europe.
- Real exposure of NurseHealth.ai: none. No compute, data, or managed dependency in that region; the Pan-American topology operates in `southamerica-west1`, `southamerica-east1`, `us-east1`, `us-central1` plus the sovereign edge.
- Lessons adopted: (1) formalize Layer 0 as explicit first defense line (section 2); (2) document the zone-region-multiregion-edge-L0-veto cascade with RPO/RTO per level; (3) register the incident as architectural motivation so DR design responds to observed events, not fashion; (4) schedule quarterly regional failover drills after pilot go-live.

## 8. Verified technical references (public)

Google Cloud Platform:

- [Multi-regional deployment archetype](https://docs.cloud.google.com/architecture/deployment-archetypes/multiregional)
- [Architecting disaster recovery for cloud infrastructure outages](https://docs.cloud.google.com/architecture/disaster-recovery)
- [DR for locality-restricted workloads](https://docs.cloud.google.com/architecture/architecting-disaster-recovery-for-locality-restricted-workloads)
- [Regions and zones](https://docs.cloud.google.com/compute/docs/regions-zones) and [Cloud Run locations](https://docs.cloud.google.com/run/docs/locations)
- [Serve traffic from multiple regions in Cloud Run](https://docs.cloud.google.com/run/docs/multiple-regions) and [Cloud Run service health](https://docs.cloud.google.com/run/docs/configuring/configure-service-health)
- [Global external Application Load Balancer with Cloud Run](https://docs.cloud.google.com/load-balancing/docs/https/setup-global-ext-https-serverless) and [failover overview](https://docs.cloud.google.com/load-balancing/docs/https/global-applb-failover-overview)
- [Cloud Armor security policy overview](https://docs.cloud.google.com/armor/docs/security-policy-overview)
- [Cloud SQL cross-region replicas](https://docs.cloud.google.com/sql/docs/postgres/replication/cross-region-replicas) and [Cloud SQL DR overview](https://docs.cloud.google.com/sql/docs/postgres/intro-to-cloud-sql-disaster-recovery)
- [Cloud Storage Turbo Replication](https://docs.cloud.google.com/storage/docs/managing-turbo-replication)
- [Cloud Healthcare API regions](https://docs.cloud.google.com/healthcare-api/docs/regions), [FHIR import/export](https://docs.cloud.google.com/healthcare-api/docs/how-tos/fhir-import-export) and [FHIR PITR](https://docs.cloud.google.com/healthcare-api/docs/fhir-pitr)
- [Pub/Sub reliability](https://docs.cloud.google.com/pubsub/docs/reliability-intro) and [message storage policies](https://docs.cloud.google.com/pubsub/docs/resource-location-restriction)
- [BigQuery managed disaster recovery](https://docs.cloud.google.com/bigquery/docs/managed-disaster-recovery)
- [Cloud Spanner multi-region](https://docs.cloud.google.com/spanner/docs/create-manage-instances) and [Spanner SLA](https://cloud.google.com/spanner/sla)

Research and repository:

- [Harness2: Recursive Agent Harnessing for an Open World (arXiv:2610.14224v1)](https://arxiv.org/html/2610.14224v1)
- [Barto912/nursehealth-ai](https://github.com/Barto912/nursehealth-ai)

Internal artifacts referenced by title only (no public URL, by design, to keep the link-checker job green):

- Documento 1, Consulta Legal de NurseHealth.ai (Google Docs, internal access)
- Documento 2, Consulta Farmacéutica (Google Docs, internal access)
- Nota de Ingeniería RDD / PROTO-006-LIF (Google Docs, internal access)
- Manifiesto de Enfermería de Precisión, 00_Master_Architecture_GentleAI_Pixel11_Compliance.md (internal notebook)
- Visión Integral de la Arquitectura NurseHealth-AI, Harness2 y Gentle-AI (internal PDF)
- Especificación Técnica de los 6 Protocolos Operativos e Integración Stanford SwiLA (Google Docs, internal access)
- Cuaderno de investigación en agentes autoevolutivos (ATDP, SEW, CoRE, Harness2; internal notebook)
