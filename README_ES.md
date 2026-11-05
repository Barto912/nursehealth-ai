# NurseHealth.ai | Plataforma Panamericana de Enfermería Domiciliaria de Precisión y Bioseguridad Verde

**Enfermería Domiciliaria de Precisión, IA Híbrida Edge-to-Cloud, Trazabilidad Criptográfica y Logística Farmacéutica Inversa Verde.**

🌐 **Leer esta documentación en otros idiomas:**
- [🇺🇸 English (Official Google for Startups & Pan-American Architecture)](./README.md)
- [🇧🇷 Português (Operações Brasil / ANVISA & SUS)](./README_PT.md)

> 🚧 **Estado del Proyecto: INVESTIGACIÓN, I+D Y VALIDACIÓN ARQUITECTÓNICA (`v0.4.0-research-proposal`).**  
> Este repositorio contiene la arquitectura fundacional de I+D Experimental, el núcleo determinista en Go (`stdlib-only`) y los pipelines soberanos de destilación de IA de NurseHealth.ai. El despliegue clínico en vivo se encuentra estrictamente controlado bajo protocolos de validación institucional.

---

## 1. Resumen Ejecutivo y Doctrina de Seguridad Sanitaria de las Américas

Miles de pacientes crónicos y post-agudos en todo el continente reciben atención domiciliaria sin trazabilidad verificable. Hoy, el cuidado extrahospitalario permanece fragmentado: registros en papel, cadenas de frío comprometidas, medicamentos vencidos acumulados en los hogares y decisiones clínicas sin respaldo auditable.

Enmarcada en la **Doctrina de Seguridad Sanitaria de las Américas**, **NurseHealth.ai** transforma la internación y enfermería domiciliaria descentralizada en una red auditable de bioseguridad Zero-Trust y logística inversa verde, construida nativamente sobre el ecosistema **Google Cloud Platform (GCP)** y **Google AI (Gemma / MedGemma)**.

- **Piloto Clínico y Logístico Inicial:** Cohorte de 100 casos distribuida en 57 urbanizaciones cerradas (*countries y barrios privados*) del **corredor Canning (Esteban Echeverría / Ezeiza), Provincia de Buenos Aires, Argentina**.
- **Hoja de Ruta de Escala Continental:** Expansión escalonada hacia **Brasil, Colombia y México (SUS/RNDS y ANVISA) y Estados Unidos** bajo un protocolo criptográfico y clínico unificado (HL7 FHIR R4, SNOMED CT y SHA-256).
- **Tesis de la Plataforma:** *El juicio es humano; la trazabilidad pertenece al sistema.* NurseHealth.ai jerarquiza la enfermería autónoma convirtiéndola en una carrera verificable, al tiempo que previene la contaminación farmacológica doméstica mediante logística inversa eco-ruteada.

---

## 2. Capacidades Nucleares y Solución de Bioseguridad Verde

1. **IA Híbrida en Tres Capas: Edge-to-Workstation-to-Cloud (Ecosistema Gemma):**
   * **Triage en el Borde (En Domicilio):** Verificación de medicamentos en tiempo real, documentación de heridas y triage clínico en el dispositivo **Google Pixel 11 Pro** del profesional sin exponer datos filiatorios (PII/PHI) crudos a redes externas.
   * **Orquestación Cloud:** Razonamiento clínico multimodal mediante **MedGemma Maestro (7B/27B) en GCP Vertex AI** gobernado por guardrails clínicos dinámicos.
2. **Logística Farmacéutica Inversa Verde e IoT de Cadena de Frío:**
   * Recolección domiciliaria de medicamentos vencidos, en desuso o termosensibles (+2°C a +8°C) utilizando contenedores térmicos IoT (Bluetooth 6 / Thread).
   * Rutas de recolección optimizadas en huella de carbono mediante **Google Maps Platform (Routes API Eco-Friendly Routing & Fleet Engine)**, asegurando la cadena de custodia de punta a punta hasta su disposición final biosegura certificada.
3. **Motor Determinista Fail-Close (`100% Go stdlib-only`):**
   * Cero dependencias externas en el núcleo de decisión clínica. Ante cualquier discrepancia lógica, discordancia taxonómica SNOMED CT o alerta de fatiga del operador, el motor bloquea mecánicamente la ejecución, emite un recibo de rechazo firmado y escala a una autoridad médica humana.
4. **Libro Mayor de Auditoría Criptográfica Inmutable:**
   * Cada intervención clínica, validación de compuerta cognitiva y transferencia de custodia farmacéutica se sella con **SHA-256** en un almacén de eventos *append-only* (solo anexado).

---

## 3. Arquitectura Técnica de Tres Capas y Flujo Híbrido de IA

Para satisfacer simultáneamente la **escalabilidad en Google Cloud AI** y los **mandatos de Soberanía de Datos de Argentina y las Américas**, NurseHealth.ai opera en tres niveles físicos:

1. **Nivel 1 — Nodo Móvil en el Borde (Google Pixel 11 Pro / Android 17):**
   * Ejecuta **MedGemma Small (2B INT8)** cuantizado vía LiteRT sobre la **NPU Google Tensor G6** (`< 300ms` de latencia).
   * El enclave de seguridad por hardware **Titan M3** firma localmente los recibos `VerificationReceipt` mediante SHA-256.
   * **Almacén de Contexto Local:** Base SQLite + FTS5 local que depura y seudonimiza datos sensibles (PHI) en el borde antes de cualquier transmisión.
2. **Nivel 2 — Orquestación Cloud y Ledger Multirregión (Google Cloud Platform):**
   * **Google Cloud Run:** Microservicios 100% Go, API REST e ingesta por Event Sourcing.
   * **Vertex AI:** MedGemma Maestro (7B/27B) con guardrails dinámicos.
   * **Cloud SQL Enterprise Plus:** Event Store PostgreSQL Append-Only con encadenamiento SHA-256.
   * **Cloud Healthcare API:** Interoperabilidad HL7 FHIR R4 y SNOMED CT aislada detrás de una Capa Anticorrupción (ACL) basada en Domain-Driven Design (DDD).
   * **Google Maps Platform:** Routes API y Fleet Engine para despacho domiciliario de bajas emisiones y recolección farmacéutica inversa.

| Capa Arquitectónica | Tecnología y Estándares |
| :--- | :--- |
| **Backend Core** | 100% Go (núcleo `stdlib-only`) en **Google Cloud Run** (Event Sourcing y Motor Fail-Close) |
| **Nodo Móvil Edge** | **Google Pixel 11 Pro** (Tensor G6 NPU + Titan M3) y MedGemma Cuantizado (2B INT8) |
| **Cloud AI (GCP)** | **MedGemma Maestro (7B/27B)** en **Vertex AI** con guardrails clínicos dinámicos |
| **Logística Verde y Geo** | **Google Maps Platform** (Routes API Eco-Routing & Fleet Engine) + IoT Cadena de Frío BLE 6/Thread |
| **Ledger de Auditoría y BD** | **Cloud SQL Enterprise Plus (PostgreSQL)** Almacén Append-Only de Cadena de Custodia SHA-256 |
| **Interoperabilidad Clínica** | **HL7 FHIR R4 y SNOMED CT** vía **Google Cloud Healthcare API** (Capa Anticorrupción DDD) |

---

## 4. Filosofía del Cuidado: El Firewall Humanista (ADR-006 y ADR-010)

Inspirado en la ética médica humanista (*Magnifica Humanitas* — cero iatrogenia, dignidad humana por sobre la automatización), el **Firewall Humano** exige integridad cognitiva antes de habilitar cualquier acción clínica.

Estos axiomas (claridad mental obligatoria, modulación de carga de trabajo y atención dirigida para la firma SNOMED CT) están compilados como **invariantes de dominio estrictas en Go** dentro del núcleo de gobernanza (`pkg/governance/`), respaldados por pruebas unitarias automatizadas (`go test ./...`). El sistema impone el **Desarme de la IA (AX-012)**: la IA puede sugerir, pero únicamente un profesional humano matriculado puede autorizar acciones clínicas irreversibles.

---

## 5. Resiliencia Multirregión y Recuperación ante Desastres (DR)

NurseHealth.ai garantiza la continuidad clínica bajo la **Doctrina de Seguridad Sanitaria de las Américas**:
- **Capa 0 (Alta Disponibilidad Intra-Regional):** Operación sincrónica multizona en `southamerica-west1` (Santiago) y `us-east1` (Virginia del Norte), alcanzando **RPO = 0** y **RTO < 60s**.
- **Fail-Close Inmutable:** Ninguna conmutación por error de infraestructura puede omitir las compuertas cognitivas del operador ni las invariantes de seguridad del paciente.
- **Continuidad Soberana:** Arquitectura diseñada para impedir que caídas regionales de nube interrumpan la asistencia domiciliaria crítica.

---

## 6. Alineación Regulatoria Multijurisdiccional y Marco de I+D

Diseñado desde el día cero para satisfacer marcos hemisféricos de bioseguridad, regímenes de promoción de I+D y normativas sanitarias:

- 🇦🇷 **Argentina (Jurisdicción Piloto y Hub de I+D):**
  * **Privacidad de Datos y Derecho Sanitario:** Ley 25.326 (Protección de Datos Personales), Ley 26.529 (Derechos del Paciente) y Ley 27.553 (Receta Electrónica y Teleasistencia).
  * **Bioseguridad Clínica y Farmacéutica:** Alineado con los lineamientos de **ANMAT** para Software como Dispositivo Médico (SaMD), **Ley Nacional 24.051** (Residuos Peligrosos / Farmacéuticos) y **Ley Provincial PBA 11.347** (Residuos Patogénicos) para la trazabilidad de logística inversa.
  * **Régimen de Economía del Conocimiento e I+D (Ley 27.506 - LEC):** Estructurado bajo los Códigos CLAE de ARCA **620100** (Arquitectura y Desarrollo de Software) y **721000** (Investigación y Desarrollo Experimental en Ingeniería y Ciencias de la Salud) a través de nuestra entidad operativa SAS en Argentina.
- 🇺🇸 **Estados Unidos (Hub Corporativo y de Escalamiento):** Cumplimiento HIPAA, FDA SaMD (Good Machine Learning Practice), SOC 2 Type II y NDAA Sec. 889 mediante estructura matriz Delaware C-Corp.
- 🇧🇷 **Brasil:** LGPD, ANVISA (RDC Nº 657/2022 para SaMD y RDC Nº 430/2020 para cadena de frío farmacéutica) y preparación nativa para interoperabilidad con el **SUS (RNDS)**.

---

## 7. Estructura del Repositorio, ADRs y Propuesta de Investigación GCP

- **Architecture Decision Records (ADRs):** Consolidados del ADR-001 al ADR-010 (incluyendo Firewall Humano, Extensiones FHIR R4, Ledger de Auditoría SHA-256 y Arquitectura Disarm AI).
- **Integración Continua (CI/CD):** Pipelines 100% en verde ejecutando pruebas unitarias/RDD en Go, validación de enlaces Markdown, auditoría de ADRs y escaneo automatizado de secretos con `gitleaks`.
- 📂 **Propuesta Detallada de Investigación GCP (`/docs/gcp-proposal/`):**
  1. [`01_architecture_overview.md`](./docs/gcp-proposal/01_architecture_overview.md) — Topología Híbrida de Punta a Punta (Edge / Workstation / Cloud).
  2. [`02_panamerican_security_and_compliance.md`](./docs/gcp-proposal/02_panamerican_security_and_compliance.md) — Gobernanza Zero-Trust, HIPAA, ANMAT y ANVISA.
  3. [`03_edge_and_gentle_ai.md`](./docs/gcp-proposal/03_edge_and_gentle_ai.md) — Cuantización MedGemma, NPU Pixel 11 Pro y Flujo Gentle-AI RDD/ODD.
  4. [`04_multiregion_dr_and_failover.md`](./docs/gcp-proposal/04_multiregion_dr_and_failover.md) — Conmutación Multirregión (`southamerica-west1` / `us-east1`).

---

## 8. Equipo Fundador y Ejecutivo

Liderazgo multidisciplinario que integra ingeniería de software, medicina clínica, gestión de proyectos, derecho de la salud y logística farmacéutica:

- **Walter Bartolomé (`@Barto912`)** — Fundador, Product Owner y Arquitecto Principal de Sistemas y Ética.
- **Leonor Figueroa** — Technical Project Manager y Asesora Estratégica (Auditoría Documental, Cumplimiento de Repositorio y Operaciones del Piloto Canning).
- **Dr. Daniel Rojo** — Chief Medical Information Officer (CMIO) y Director Médico (Taxonomías Clínicas SNOMED CT, Extensiones FHIR R4 y Protocolos de Enfermería Domiciliaria).
- **División de Asuntos Legales y Regulatorios** — Derecho Sanitario, Privacidad de Datos (Ley 25.326 / HIPAA) y Estructuración Societaria (Consejo Asesor y Dirección Fraccional).
- **División de Desarrollo de Negocios y Operaciones** — Unit Economics y Alianzas Público-Privadas / B2B en Salud (Consejo Asesor y Dirección Fraccional).
- **División de Logística Farmacéutica y Bioseguridad** — Cumplimiento de Cadena de Frío (+2°C a +8°C), Logística Inversa Verde y Trazabilidad ANMAT (Farmacéuticos Especialistas).

*Contacto Institucional:* `contact@nursehealth.ai`

---

## 9. Licencia y Propiedad Intelectual

**Modelo de Licenciamiento Híbrido:**
- **Propietario (Todos los Derechos Reservados):** Los motores centrales en Go, protocolos criptográficos y pipelines de destilación de IA constituyen propiedad intelectual exclusiva de NurseHealth.ai.
- **Abierto / Revisión Académica (Apache 2.0 / CC-BY-NC-SA 4.0):** Los adaptadores FHIR/SNOMED, la documentación ADR y los esquemas de cumplimiento regulatorio se publican para promover la interoperabilidad, la revisión académica y la auditoría institucional.
