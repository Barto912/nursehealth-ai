# NurseHealth.ai | Plataforma Panamericana de Enfermería Domiciliaria y Bioseguridad

**Enfermería domiciliaria de precisión con IA en el borde y trazabilidad criptográfica.**

[![CI](https://github.com/Barto912/nursehealth-ai/actions/workflows/main.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Docs](https://github.com/Barto912/nursehealth-ai/actions/workflows/auditoria-docs.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Seguridad](https://github.com/Barto912/nursehealth-ai/actions/workflows/seguridad-gitleaks.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)

🌐 **Leer esta documentación en otros idiomas:**
- [🇺🇸 English (Official GCP Proposal)](./README.md)
- [🇧🇷 Português (Operações Brasil / ANVISA & SUS)](./README_PT.md)

> 🚧 **Estado del proyecto: INVESTIGACIÓN Y DISEÑO.**
> Este repositorio no es un producto clínico en producción. El código y los documentos aquí presentes son experimentales y su uso en entornos clínicos reales no está autorizado ni validado.

---

## 🛡️ El Problema y la Doctrina de Seguridad Sanitaria

Miles de pacientes crónicos en las Américas reciben atención domiciliaria sin trazabilidad real. Ese trabajo hoy es invisible: planillas de papel, medicamentos vencidos en el hogar, cadenas de custodia rotas y decisiones clínicas sin evidencia auditable.

Enmarcado en la **Doctrina de Seguridad Sanitaria de las Américas**, **NurseHealth.ai** transforma el cuidado descentralizado en una red de bioseguridad auditable (Zero-Trust) construida nativamente sobre el ecosistema de **Google Cloud Platform (GCP)**.

* **Piloto Inicial:** 57 countries y barrios cerrados en Canning, Buenos Aires, Argentina.
* **Visión Continental:** Escalamiento regional a Chile, Uruguay, Colombia, México, Brasil (SUS/ANVISA) y Estados Unidos bajo un protocolo clínico y criptográfico unificado.

---

## 💡 La Solución

Plataforma extra-hospitalaria de enfermería de precisión y logística inversa farmacéutica:

- **IA Híbrida (Cloud/Edge):** Identificación de medicamentos y triaje clínico en el dispositivo del enfermero (**Google Pixel 11 Pro**), sin exponer datos sensibles a la nube.
- **Cadena de Custodia Criptográfica:** Cada evento clínico se sella con **SHA-256** en un almacén de eventos *append-only*.
- **Logística Inversa Farmacéutica:** Retiro domiciliario de medicamentos vencidos con trazabilidad hasta su disposición final.
- **Fail-Close por Diseño:** Ante la duda, el motor determinista en Go bloquea la ejecución y escala a una autoridad médica humana. Nunca improvisa.

> **Lema:** *El criterio es humano, la trazabilidad es del sistema.*

**Tesis de plataforma:** Enfermería de precisión con soporte de IA que convierte el trabajo autónomo en una carrera profesional verificable: formación continua, reputación auditable y dignidad económica para los cuidadores de las Américas.

---

## 🏗️ Arquitectura Técnica y Flujo de IA Híbrida

1. **Nodo Edge — Google Pixel 11 Pro:** MedGemma cuantizado (2B/4B) en Tensor NPU + TFLite + Engram (SQLite/FTS5 almacenamiento local).
2. **Transporte Zero-Trust:** Flujo de eventos clínicos seudonimizados y firmados criptográficamente.
3. **Ingesta Core — Google Cloud Run:** Microservicios 100% Go (Golang), API REST, Motor Fail-Close determinista y Event Sourcing.
4. **Orquestación IA Cloud — Vertex AI:** MedGemma Maestro (7B/27B) + Google Harness² (guardrails dinámicos *Propose-Probe-Compose*).
5. **Auditoría Inmutable — Cloud SQL:** PostgreSQL, almacén de eventos SHA-256 (Append-Only).
6. **Interoperabilidad Clínica — Cloud Healthcare API:** HL7 FHIR R4 y SNOMED CT aislados tras una capa anticorrupción (ACL) bajo Domain-Driven Design.

### Stack Técnico

| Capa | Tecnología |
| :--- | :--- |
| **Backend Core** | 100% Go (Golang) en Google Cloud Run (Event Sourcing & Fail-Close) |
| **Frontend Mobile** | React Native con Expo |
| **Edge AI (Dispositivo)** | MedGemma Small (2B/4B) cuantizado & TensorFlow Lite en Pixel 11 Pro |
| **Cloud AI (GCP)** | MedGemma Maestro (7B/27B) en Vertex AI + Google Harness² |
| **Gobernanza IA & Memoria** | Gentle-AI (ODD/RDD) + Engram (SQLite + FTS5 local-first) |
| **Base de Datos & Ledger** | PostgreSQL en Cloud SQL (Event Store SHA-256 Append-only) |
| **Interoperabilidad** | HL7 FHIR R4 & SNOMED CT vía Google Cloud Healthcare API |
| **Infraestructura** | Google Cloud Platform (GCP) |

---

## ⚖️ Alineación Regulatoria Multijurisdiccional

Diseñado desde el primer día para cumplir con los estándares hemisféricos de seguridad de la salud:
* 🇺🇸 **Estados Unidos:** HIPAA, SOC 2 Type II, alineación CISA y cumplimiento NDAA.
* 🇧🇷 **Brasil:** LGPD, ANVISA (SaMD) y preparación para interoperabilidad con **SUS (RNDS)**.
* 🌎 **Hispanoamérica:** Marcos de privacidad locales (Argentina Ley 25.326, Chile, Colombia, México) y normativas de Historias Clínicas Electrónicas.

---

## 📈 Estado Actual e Hitos

Diseño de arquitectura completado y auditado:
- **Event Storming (MVP-0):** Completado con 14 hotspots clínicos/operativos resueltos.
- **Límites de Agregados v0.3:** Patrones de reparación vs aceptación documentados.
- **Context Map (DDD):** Capa Anticorrupción (ACL) para FHIR R4 y Modelo de Referencia Regulatorio.
- **Integración Continua:** CI/CD en verde, auditoría automática de docs y escaneo `gitleaks`.

**En curso:** Implementación del MVP y orquestación de Cloud AI para postular al **Google for Startups Cloud Program**.

---

## 🧭 Gobernanza Ética

Principios no negociables del manifiesto del proyecto:

1. **Dignidad sobre datos.**
2. **Autonomía sobre automatización.**
3. **Transparencia sobre caja negra.**
4. **Cuidado sobre eficiencia.**
5. **Presencia sobre conexión.**
6. **Prevención sobre reacción.**
7. **Formación sobre adoctrinamiento.**

Los datos de los pacientes se procesan de forma estrictamente seudonimizada. Ninguna decisión clínica se automatiza.

---

## 🧠 Filosofía del Cuidado: El Firewall Humano

La enfermería de precisión requiere más que habilidad técnica: exige integridad cognitiva. El **Firewall Humano** (ADR-006) codifica tres axiomas fundacionales que condicionan cada evento clínico:

- **AX-001 (Inercia y Ruido):** la claridad mental como precondición. Si `MentalClarity < 7`, el sistema bloquea la creación del evento y registra el intento.
- **AX-002 (Tonos / nivel de trabajo):** el tono energético modula la eficacia clínica. Si `WorkLevel == LOW`, el sistema sugiere reasignación o pausa.
- **AX-003 (Atención simple):** la atención dirigida como condición para firmar el evento SNOMED. Si `SimpleAttention == false`, el operador debe registrar un `AttentionalFinding` antes de firmar.

Estos axiomas no son aspiracionales: son **invariantes de dominio** exigidas por el Value Object `RuidoCheck` con tests unitarios. Un profesional en estado cognitivo no apto no puede crear ni firmar un evento clínico: el sistema bloquea y registra el intento.

**Implementación:**

- [Axiomas fundacionales](docs/philosophy/axioms.json) (AX-001 a AX-003)
- [Value Object RuidoCheck](pkg/cognitive/ruido_check.go) (gates tipados con tests)
- [ADR-006: Firewall Humano (MVP)](docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md)
- [Context Map v0.2](docs/domain/context-map.md) (invariantes I10, I11)

---

## 👥 Equipo Core

- **Barto** — Founder & Product Owner.
- **Leonor** — Strategic Advisor.

---

## 📄 Licencia

Propietario. Todos los derechos reservados.
## 🛡️ Arquitectura de Resiliencia Multi-Capa y Soberanía de Datos

NurseHealth.ai opera bajo la **Doctrina de Seguridad Sanitaria de las Américas**, garantizando la continuidad clínica y el cumplimiento normativo en todo momento.

- **Capa 0 (Intra-Regional):** Nuestra primera línea de defensa. Operamos con alta disponibilidad síncrona multi-zona en `southamerica-west1` (Chile), garantizando RPO=0 y RTO<60s para el piloto de Clínica Canning Health, cumpliendo con la **Ley 25.326 (Protección de Datos Personales)** y normativas de **ANMAT**.
- **Fail-Close Inmutable (Invariantes I10/I11):** Ningún fallo de infraestructura habilita reintentos automáticos si los gates cognitivos del operador (`RuidoCheck`, AX-001 a AX-003) o del paciente (AX-004 a AX-006) han bloqueado la operación. El veto clínico es determinista y se registra de forma inmutable (append-only) en el ledger.
- **Motivación Operativa:** Diseñado tras el incidente público de `europe-west2` (octubre 2026). Nuestra topología panamericana opera de forma soberana, sin dependencia de regiones externas, asegurando que la atención de enfermería de precisión nunca se detenga.

*Para detalles técnicos de replicación cruzada y conmutación por error, consulte `docs/gcp-proposal/04_multiregion_dr_and_failover.md`.*
