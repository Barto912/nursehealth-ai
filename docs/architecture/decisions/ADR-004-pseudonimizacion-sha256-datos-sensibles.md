# ADR-004: Pseudonimización SHA-256 de Datos Sensibles

## Estado

Aceptado

## Contexto

NurseHealth-AI persiste eventos cognitivos (CognitiveStateValidated) que documentan cada interacción entre el profesional de enfermería y el sistema de triage asistido por IA. Estos eventos se almacenan en la tabla de auditoría definida en `db/migrations/001_create_cognitive_audit_log.sql` y gestionada por `pkg/repository/audit_postgres.go`.

Cada evento cognitivo necesita referenciar al paciente sobre el cual se ejecutó el triage para permitir trazabilidad clínica. Sin embargo, persistir el identificador real del paciente (PatientID) en la capa de auditoría introduce riesgos significativos:

- Violación del principio de minimización de datos exigido por GDPR (Art. 5.1.c): la tabla de auditoría no requiere el identificador real del paciente para cumplir su función de trazabilidad del proceso cognitivo
- Superficie de exposición ampliada: en caso de compromiso de la base de datos de auditoría, un atacante obtendría acceso directo a identificadores de pacientes reales
- Separación de responsabilidades insuficiente: la capa de auditoría no debe tener capacidad de resolver identidades clínicas reales, solo debe registrar que un proceso cognitivo se ejecutó sobre un sujeto determinado
- Cumplimiento HIPAA: la regla de seguridad requiere salvaguardas técnicas que limiten el acceso a PHI (Protected Health Information) únicamente a sistemas que lo necesiten funcionalmente

El struct `CognitiveStateValidated` definido en `pkg/triage/events.go` contiene los campos verificados: ID, CareID, OperatorID, PatientHash, Timestamp y GateSnapshot. El campo PatientHash es el mecanismo designado para referenciar al paciente sin exponer su identificador real.

## Decisión

Se adopta la pseudonimización del PatientID mediante SHA-256 como mecanismo de protección de datos sensibles en la capa de persistencia de auditoría.

### Mecanismo de pseudonimización

El PatientID original se transforma mediante SHA-256 antes de cualquier persistencia en la tabla de auditoría. El resultado se almacena en el campo PatientHash del struct CognitiveStateValidated. Esta transformación es:

- Unidireccional: no existe operación criptográfica que recupere el PatientID original a partir del PatientHash
- Determinista: el mismo PatientID produce siempre el mismo PatientHash, permitiendo correlación de eventos sobre un mismo paciente sin exponer su identidad
- De longitud fija: 256 bits (64 caracteres hexadecimales) independientemente de la longitud del PatientID original

### Separación del mapping service

La capacidad de resolver PatientHash a PatientID real se mantiene en un servicio separado y restringido (mapping service), fuera del alcance de la capa de auditoría. Este servicio:

- No es accesible desde `pkg/repository/audit_postgres.go` ni desde `pkg/triage/events.go`
- Requiere autenticación y autorización específica para operaciones de de-pseudonimización
- Mantiene su propio log de acceso para auditoría de quién y cuándo resolvió identidades
- Opera bajo principio de least-privilege: solo componentes con necesidad funcional explícita pueden invocarlo

### Alcance de la pseudonimización

La pseudonimización aplica exclusivamente al PatientID. Otros campos del struct CognitiveStateValidated no se pseudonimizan:

- ID: identificador del evento cognitivo (no sensible)
- CareID (columna `care_event_id`): identificador del episodio de cuidado; clasificado como no-PHI directo bajo criterio de minimización (confirmar con responsable de cumplimiento)
- OperatorID: identificador del profesional de enfermería (no se pseudonimiza; el operador tiene legitimidad para ser identificado en auditoría) (clasificación asumida; confirmar con responsable de cumplimiento)
- Timestamp: momento de la validación (no sensible)
- GateSnapshot: estado del Cognitive Firewall en el momento de la validación (no contiene PHI)

### Correspondencia struct→tabla

Correspondencia struct→tabla (migration 001): `PatientHash`→`patient_hash`, `CareID`→`care_event_id`, `Timestamp`→`created_at`.

El contexto del paciente gestionado en `pkg/triage/patient_context.go` contiene el PatientID original durante el procesamiento en memoria, pero este identificador nunca se propaga a la capa de persistencia de auditoría.

### Ausencia de salt por paciente

Se decide deliberadamente no aplicar salt criptográfico por paciente al cálculo SHA-256. Esta decisión se documenta explícitamente por sus implicaciones (ver Consecuencias Negativas).

## Consecuencias

### Positivas

- Minimización de datos by-design: la tabla de auditoría nunca contiene PHI directo, cumpliendo GDPR Art. 5.1.c y reduciendo scope de cumplimiento HIPAA para este componente
- Separación de responsabilidades verificable: la capa de persistencia (`pkg/repository/audit_postgres.go`) no tiene dependencia alguna del mapping service, verificable por análisis estático de imports
- Superficie de ataque reducida: compromiso de la base de datos de auditoría no expone identidades de pacientes reales sin acceso adicional al mapping service
- Correlación de eventos preservada: la naturaleza determinista de SHA-256 permite agrupar eventos cognitivos por paciente para análisis clínico sin de-pseudonimizar
- Compatibilidad con derecho al olvido (GDPR Art. 17): eliminar la entrada correspondiente en el mapping service efectivamente desvincula los eventos de auditoría del paciente real
- Integridad del modelo de dominio: el campo PatientHash en CognitiveStateValidated refleja explícitamente que la capa de triage opera sobre identidades pseudonimizadas, no sobre PHI
- Auditoría de accesos simplificada: cualquier operación de de-pseudonimización pasa por el mapping service, generando un único punto de control y logging

### Negativas

- Vulnerabilidad a ataques de diccionario: al no utilizar salt por paciente, un atacante con acceso al PatientHash y conocimiento del espacio de posibles PatientIDs (ej: formato secuencial, número de seguridad social) puede computar SHA-256 sobre candidatos y comparar. Mitigación requerida: el mapping service debe aplicar rate-limiting y los PatientIDs originales deben tener entropía suficiente (no secuenciales predecibles)
- Imposibilidad de de-pseudonimización criptográfica: si se pierde el mapping service o su base de datos de correspondencias, los PatientHash en la tabla de auditoría se vuelven permanentemente irresolubles. Mitigación: el mapping service requiere backup y plan de recuperación con RPO definido
- Queries de auditoría limitadas: no es posible ejecutar queries SQL directas sobre PatientID real en la tabla de auditoría. Cualquier búsqueda por paciente requiere resolución previa vía mapping service o búsqueda por PatientHash conocido. Esto incrementa latencia en flujos de investigación clínica
- Overhead computacional por evento: cada evento cognitivo requiere un cálculo SHA-256 adicional sobre PatientID antes de persistencia. Aunque SHA-256 es eficiente, en escenarios de alto throughput (miles de eventos/segundo) este overhead debe monitorizarse
- Irreversibilidad ante error de pseudonimización: si se calcula PatientHash sobre un PatientID incorrecto (error en capa de entrada), el evento queda permanentemente asociado a un hash que no corresponde a ningún paciente real. No existe mecanismo de corrección criptográfica; solo puede insertarse evento de invalidación (patrón append-only de ADR-003)
- Complejidad operativa añadida: el mapping service es un componente adicional que requiere despliegue, monitorización, backup y gestión de acceso independiente, incrementando la superficie operacional del sistema
- Testing más complejo: los tests deben verificar que PatientID nunca llega a la capa de persistencia en su forma original, requiriendo fixtures con mapping service mockeado y aserciones explícitas sobre el contenido de PatientHash

## Referencias

- `pkg/triage/events.go` (definición de CognitiveStateValidated con campo PatientHash)
- `pkg/triage/patient_context.go` (contexto del paciente con PatientID original en memoria)
- `pkg/repository/audit_postgres.go` (persistencia append-only que recibe PatientHash, nunca PatientID)
- `db/migrations/001_create_cognitive_audit_log.sql` (schema de tabla audit con columna patient_hash)
- `docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md` (Cognitive Firewall como origen de eventos)
- `docs/architecture/decisions/ADR-007-fhir-extensions.md` (extensiones FHIR que interactúan con PatientHash)
