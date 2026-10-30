# WS5 — Filosofía del Cuidado: Paradigma Humanístico de la Plataforma Determinista

## Estado
Propuesto. Documento de entrenamiento interno obligatorio y referencia arquitectónica. No constituye guía clínica ni reemplaza el juicio del profesional matriculado.

## Fuente Normativa
`docs/philosophy/axioms.json` (Axiomas AX-001 a AX-006). Este README es la capa interpretativa y pedagógica de esos axiomas; en caso de discrepancia, prevalece el archivo JSON como fuente de la verdad.

## 1. Propósito
Traducir los axiomas ejecutables de NurseHealth-AI a un paradigma académico del cuidado, sirviendo a tres audiencias:
1. **Enfermeros/as:** Currículum obligatorio de humanización previo a la operación de la plataforma.
2. **Ingeniería:** Requisitos éticos expresados como contratos ejecutables (gates, vetos, registros).
3. **Auditoría:** Trazabilidad directa: Principio → Axioma → Componente Técnico.

## 2. Los Cuatro Pilares
1. **Prevención absoluta de la iatrogenia:** El daño inducido por el sistema (incluida la iatrogenia algorítmica o por omisión emocional) se trata como defecto de diseño, no como fatalidad.
2. **Soberanía del dato y custodia de la persona:** El dato clínico es extensión de la dignidad del paciente; nunca mercancía.
3. **Humanización previa del profesional:** Ningún enfermero opera la plataforma sin completar el currículum de axiomas. La tecnología amplifica la presencia humana; no la sustituye.
4. **Seguridad Sanitaria y Custodia de la Vida en las Américas:** Estándar regional de protección epidemiológica, bioseguridad y prevención del sufrimiento evitable.

## 3. Mapeo Determinista de Axiomas a Código
La ética no es un prompt; es una restricción estructural inmutable codificada en el dominio.

| Axioma (Fuente) | Principio Filosófico / Clínico | Implementación Técnica (Invariantes DDD) | Estado |
|---|---|---|---|
| **AX-001**: Inercia y Ruido *(Obras Completas II)* | Prevención de la iatrogenia por fatiga o transferencia emocional no gestionada. | `pkg/cognitive/ruido_check.go`: Value Object `RuidoCheck` exige `MentalClarity` (1-10). Si `< 7`, bloquea `CareEvent` y registra en `cognitive_audit_log`. | ✅ Existente |
| **AX-002**: Tonos *(Autoliberación)* | La eficacia clínica es modulada por el estado cognitivo/emotivo actual, no es una habilidad estática. | `pkg/cognitive/ruido_check.go`: Exige `WorkLevel` (LOW, MODERATE, HIGH). Si `LOW`, bloquea el evento y sugiere pausa/reasignación. | ✅ Existente |
| **AX-003**: Atención Simple *(Autoliberación)* | La enfermería de precisión exige atención dirigida, evitando la pasividad perceptiva por monotonía. | `pkg/cognitive/ruido_check.go`: Exige `SimpleAttention` (bool). Si `false`, bloquea firma SNOMED hasta completar `AttentionalFinding`. | ✅ Existente |
| **AX-004**: Dolor vs. Sufrimiento *(Obras Completas I)* | El sufrimiento existencial exige respuesta clínica distinta al dolor físico; omitir escucha activa es iatrogenia. | `pkg/triage/patient_context.go`: Value Object `PatientContextCheck` exige `SufferingIndicator`. Si `true`, bloquea cierre hasta registrar `ActiveListeningNote` validado. | ✅ Existente |
| **AX-005**: Paisaje de Formación *(Autoliberación)* | El contexto generacional/cultural condiciona la percepción de la enfermedad y el triage. | `pkg/triage/patient_context.go`: Exige `ContextualLandscape`. Si permanece `UNKNOWN` tras escalación, bloquea firma con razón `ContextNotDeterminable`. | ✅ Existente |
| **AX-006**: Acción Válida (Regla de Oro) *(Obras Completas I)* | Toda acción debe ser éticamente confirmada por el operador; de lo contrario, se escala. | `pkg/triage/patient_context.go`: Exige `EthicalConfirmation`. Si `false`, emite `VerificationReceipt` (Fail-Close) y escala a supervisión. | ✅ Existente |

## 4. Fundamentos Filosóficos de Respaldo
- **El sentido en el sufrimiento (Viktor Frankl):** El domicilio del paciente es un escenario de sufrimiento con sentido potencial. La plataforma existe para liberar tiempo humano de tareas administrativas y devolverlo al encuentro (Frankl, 2006).
- **La ciencia del cuidar (Jean Watson):** El cuidado es un ideal moral. Los procesos Caritas (práctica amorosa, presencia auténtica) son el estándar contra el cual se mide cualquier automatización (Watson, 2008).
- **El error como problema de sistema (Institute of Medicine):** La iatrogenia es predominantemente un fallo de diseño de sistemas. Por eso la respuesta de NurseHealth-AI es vetar determinísticamente el acto inseguro (Institute of Medicine, 1999).

## 5. Currículum de Humanización (Onboarding Obligatorio)
- **Módulo 1 — Gestión del Ruido y Tonos (AX-001, AX-002):** 2 h. Reconocimiento de la inercia emocional y el nivel de trabajo propio antes de cada visita.
- **Módulo 2 — Atención Simple y Paisaje de Formación (AX-003, AX-005):** 2 h. Técnicas de presencia dirigida y evaluación del contexto cultural/generacional del paciente.
- **Módulo 3 — Dolor, Sufrimiento y Acción Válida (AX-004, AX-006):** 2 h. Diferenciación clínica entre dolor y sufrimiento; protocolo de escucha activa y confirmación ética (Regla de Oro).
- **Módulo 4 — Lectura del Veto Determinista:** 3 h prácticas. Cómo interpretar un Fail-Close, un gate de `RuidoCheck` y un recibo ATDP sin frustración ni automatización del criterio.
- **Evaluación:** Caso clínico simulado con veto determinista + firma de acta de compromiso ético.

## 6. Glosario Técnico-Filosófico
- **Veto determinista:** Bloqueo ejecutado por código (no por probabilidad) ante incertidumbre, fatiga o falta de confirmación ética.
- **RuidoCheck:** Telemetría cognitiva del profesional (fatiga, sesgo de automatización, estrés comunicacional) que actúa como gate de seguridad.
- **Recibo ATDP:** Registro inmutable ($o, h, a, y, r, m$) que prueba quién decidió qué, cuándo y por qué, incluyendo las señales de ruido y confirmación ética.
- **Custodia:** Responsabilidad ética y legal por la persona y sus datos, desde el domicilio hasta la disposición final.

## 7. Referencias (APA 7.ª ed.)
Frankl, V. E. (2006). *Man's search for meaning* (Rev. ed.). Beacon Press. (Obra original publicada en 1946).

Institute of Medicine. (1999). *To err is human: Building a safer health system*. National Academy Press. https://doi.org/10.17226/9728

Silo. (2006). *Obras completas* (Vols. 1-2). Editorial Planeta. (Incluye *Humanizar la Tierra* y *Psicología de la Imagen*).

Silo. (1990). *Autoliberación*. Editorial Planeta.

Watson, J. (2008). *Nursing: The philosophy and science of caring* (Rev. ed.). University Press of Colorado. (Obra original publicada en 1979).

## 8. Integridad Académica
Las citas canónicas (Frankl, Watson, IOM, Silo) son verificables en sus ediciones citadas. Este documento no sustituye protocolos clínicos vigentes ni el juicio profesional. La implementación técnica de los axiomas está sujeta a la revisión de código y pruebas unitarias en el repositorio `github.com/Barto912/nursehealth-ai`.