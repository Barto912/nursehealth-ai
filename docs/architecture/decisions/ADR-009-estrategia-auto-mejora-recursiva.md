# ADR-009: Estrategia de Auto-Mejora Recursiva para Modelos de Borde (Harness² Simplificado)

## Estado
Propuesto

## Contexto
NurseHealth-AI operará con cientos de enfermeros generando datos clínicos de alto valor durante el piloto en Canning y fases posteriores. Existe la necesidad de mejorar continuamente el modelo de borde (MedGemma Small INT8) sin reentrenar pesos (técnicamente inviable en borde y regulatormente complejo) y cumpliendo estrictamente con HIPAA, Ley 25.326 y ANMAT.

Las investigaciones sobre agentes autoevolutivos (Google Harness², protocolo ATDP) proporcionan conceptos valiosos, pero requieren adaptación realista a las limitaciones de una startup en fase piloto con recursos limitados y requisitos de soberanía de datos.

### Elementos Implementables vs. No Implementables a Corto Plazo

**✅ IMPLEMENTABLE (Fase Piloto - Meses 1-6):**
1. **EXPERIENCE.md**: Biblioteca de lecciones persistentes en markdown.
2. **Protocolo ATDP Simplificado**: Extensión de `cognitive_audit_log` con tupla de 6 componentes ($o, h, a, y, r, m$).
3. **Data Proxy**: Middleware de redactado PII/PHI antes de transmisión a nube (Planificado).
4. **Engram Local**: SQLite/FTS5 en el borde para contexto del paciente (Planificado).
5. **Esquema CoRE**: Estructuración de prompts de mapeo SNOMED CT para mejorar la Tasa de Éxito de Generación (GSR).

**❌ NO IMPLEMENTABLE A CORTO PLAZO:**
1. **Reentrenamiento de Pesos en Borde**: Requiere clusters de GPUs masivos; inviable en Google Pixel 11 Pro.
2. **Autoevolución Darwiniana Completa**: Violaría los principios de Harness² (pesos congelados, solo evoluciona el arnés) y los requisitos de congelamiento de modelos de ANMAT.
3. **Replay Contrafáctico Masivo**: Requiere infraestructura de simulación que excede el alcance del piloto inicial.

## Decisión
Adoptar un enfoque de **Harness² Simplificado** con cinco componentes principales que permiten auto-mejora recursiva sin reentrenar pesos del modelo, manteniendo coherencia con ADR-003 (append-only), ADR-004 (pseudonimización), ADR-005 (encadenamiento forense) y ADR-008 (Compliance Legal).

### Componente 1: Biblioteca de Experiencias (EXPERIENCE.md)
- **Ubicación:** `docs/experience/EXPERIENCE.md`
- **Propósito:** Archivo markdown que acumula lecciones de ejecuciones exitosas y fallidas del mapeo SNOMED CT/FHIR, permitiendo que nuevos casos clínicos se beneficien de experiencias previas sin reentrenar el modelo.
- **Actualización:** Semi-automática vía job cron semanal que extrae casos validados de `cognitive_audit_log` y genera entradas en EXPERIENCE.md.

### Componente 2: Protocolo ATDP Extendido
- **Ubicación:** `pkg/triage/events.go` y `db/migrations/002_extend_cognitive_audit_log_atdp.sql` (Planificado)
- **Propósito:** Extender la tupla de auditoría para capturar no solo el evento, sino el ciclo completo de aprendizaje.
- **Campos a agregar en `cognitive_audit_log`:**
  - `observation_json` (JSONB): Datos clínicos de entrada (redactados por Data Proxy).
  - `hidden_state_summary` (TEXT): Resumen del razonamiento del modelo (truncado para privacidad).
  - `action_fhir` (JSONB): Mapeo SNOMED CT/FHIR generado.
  - `outcome_clinical` (VARCHAR): Outcome clínico (mejoría/empeoramiento) capturado con delay de 7-30 días.
- **Tupla ATDP de 6 Componentes:** ($o$) Observación, ($h$) Estado Interno, ($a$) Acción, ($y$) Resultado/Validación humana, ($r$) Recompensa/Outcome, ($m$) Metadatos (versión del modelo, timestamp, firma Titan M3).

### Componente 3: Ciclo Propose-Probe-Compose Simplificado
- **Ubicación:** `pkg/cognitive/harness_optimizer.go` (Planificado)
- **Propósito:** Optimizar prompts y guardrails mediante exploración controlada, seleccionando variantes con mayor GSR (Tasa de Éxito de Generación).
- **Flujo Semanal:** 
  1. **Propose:** Generar 2-3 variantes de prompts para casos ambiguos.
  2. **Probe:** Ejecutar variantes en paralelo sobre casos históricos anonimizados de EXPERIENCE.md.
  3. **Compose:** Seleccionar la variante con mayor GSR (>80% según estándar CoRE) y guardar el arnés ganador.

### Componente 4: Data Proxy para Privacidad
- **Ubicación:** `pkg/adapter/data_proxy.go` (Planificado - Fase Piloto)
- **Propósito:** Middleware que censura automáticamente PII/PHI antes de transmitir datos a Vertex AI o cualquier servicio en la nube, cumpliendo con Ley 25.326 y HIPAA.
- **Funcionamiento:** Intercepta observaciones clínicas, aplica reglas de redactado (regex/NER ligero) y preserva únicamente códigos SNOMED CT, timestamps relativos y relaciones temporales.

### Componente 5: Engram Local (SQLite/FTS5)
- **Ubicación:** `pkg/engram/engram.go` (Planificado - Fase Piloto)
- **Propósito:** Base de datos local en el Google Pixel 11 Pro que almacena contexto del paciente para operación offline-first, permitiendo que la enfermera opere sin conectividad y sincronice cuando haya red.
- **Seguridad:** Cifrado AES-256 en reposo, aprovechando el enclave Titan M3 del dispositivo.

## Consecuencias

### Positivas
1. **Mejora Continua sin Reentrenamiento:** Permite optimizar el sistema mediante evolución de prompts/guardrails sin tocar pesos del modelo (técnicamente viable en borde).
2. **Privacidad Garantizada:** El Data Proxy asegura cumplimiento de HIPAA/Ley 25.326 antes de cualquier transmisión a nube.
3. **Operación Offline:** Engram permite trabajo continuo sin conectividad, crítico para enfermería domiciliaria en zonas con cobertura intermitente.
4. **Auditoría Completa:** El protocolo ATDP extendido proporciona trazabilidad end-to-end desde observación clínica hasta outcome, esencial para validación regulatoria (ANMAT/OPDS).

### Negativas
1. **Complejidad Operativa:** Implementar 5 componentes nuevos añade carga de desarrollo estimada en 6-8 semanas.
2. **Latencia de Recompensas Tardías:** El outcome clínico ($r$) puede tardar 7-30 días en capturarse, retrasando el ciclo de aprendizaje.
3. **Dependencia de Validación Humana:** EXPERIENCE.md solo puede actualizarse con casos validados por enfermeros, lo que introduce fricción operativa.

## Referencias
- `docs/experience/EXPERIENCE.md` (Biblioteca de lecciones persistentes)
- `pkg/adapter/data_proxy.go` (Middleware de redactado PII/PHI - Planificado)
- `pkg/cognitive/harness_optimizer.go` (Ciclo Propose-Probe-Compose - Planificado)
- `pkg/engram/engram.go` (Base de datos local SQLite/FTS5 - Planificado)
- `pkg/triage/events.go` (Extensión de tupla ATDP)
- `db/migrations/002_extend_cognitive_audit_log_atdp.sql` (Schema ATDP extendido - Planificado)
- `docs/domain/logistica-inversa.md` (Diseño de componentes planificados, ej. `critical_list.go`)
- `docs/architecture/decisions/ADR-003-repositorio-audit-only-append-only.md`
- `docs/architecture/decisions/ADR-004-pseudonimizacion-sha256-datos-sensibles.md`
- `docs/architecture/decisions/ADR-005-encadenamiento-forense-prev-hash.md`
- `docs/architecture/decisions/ADR-008-legal-and-pharmaceutical-compliance.md`
