# ADR-001: Arquitectura General del Sistema NurseHealth-AI

## Estado

Aceptado

## Contexto

NurseHealth-AI requiere una arquitectura que garantice trazabilidad completa de decisiones clínicas asistidas por IA, manteniendo cumplimiento normativo sanitario (HIPAA/GDPR) y preservando la autonomía del profesional de enfermería.

El sistema debe operar como herramienta asistencial de apoyo (filosofía "Gentle Care"), nunca como sistema autónomo de toma de decisiones clínicas.

Se necesita un backend que procese eventos cognitivos (triage, alertas, recomendaciones) con latencia predecible y auditabilidad total.

La persistencia debe garantizar inmutabilidad de registros: ningún evento clínico asistido por IA puede ser modificado o eliminado post-escritura.

El PatientID debe estar pseudonimizado en toda la capa de persistencia y transmisión para minimizar superficie de exposición de datos sensibles.

## Decisión

Se adopta una arquitectura monolítica modular en Go, con separación clara de capas: HTTP handlers, servicios de dominio, repositorios y motor cognitivo.

La persistencia utiliza PostgreSQL en modo audit-only append-only: todas las tablas son inmutables tras inserción, sin UPDATE ni DELETE a nivel de aplicación.

El motor cognitivo (Cognitive Firewall) opera como middleware de triage que intercepta requests clínicos, aplica reglas de seguridad y genera eventos inmutables. Su implementación vive en `pkg/triage/middleware.go` mediante la función `ValidateCareEvent()`, que ejecuta los gates en orden: operator-side primero, luego patient-side.

Los identificadores de paciente se pseudonimizan mediante SHA-256 antes de cualquier persistencia, manteniendo un mapping service separado y restringido.

La comunicación entre componentes sigue un patrón de eventos: cada interacción clínica genera un CognitiveEvent con timestamp, hash de integridad y metadata de contexto.

No se implementa cola de mensajes externa (Kafka/RabbitMQ) en MVP; la serialización de eventos se realiza directamente en PostgreSQL con garantías ACID.

### Hash de integridad y timestamp

- La capa de persistencia calcula `payload_hash` (SHA-256 del JSON canónico del evento) antes de la inserción — ver columna `payload_hash` en `db/migrations/001_create_cognitive_audit_log.sql`.
- El timestamp de auditoría (`created_at`) es generado por el servidor PostgreSQL mediante `DEFAULT NOW()`, no por el cliente.

### Encadenamiento forense (prev_hash)

Cada fila de `cognitive_audit_log` almacena `prev_hash` (hash de la entrada anterior), formando una cadena tamper-evident: alterar cualquier fila histórica rompe la cadena desde ese punto. El procedimiento completo de verificación se documenta en ADR-005.

### Triage clínico complementario

El sistema incluye además un módulo de triage clínico tradicional en `pkg/triage/triage.go` con la función `Classify(systolicBP, heartRate, spo2)` que asigna códigos Celeste/Amarillo/Rojo según umbrales de signos vitales. Este componente es determinista, puro y está separado del Cognitive Firewall cognitivo.

## Consecuencias

### Positivas

- Trazabilidad completa: cada decisión asistida por IA queda registrada como evento inmutable con contexto completo
- Simplicidad operacional en MVP: un solo binario Go + PostgreSQL, sin dependencias de infraestructura de mensajería
- Cumplimiento normativo facilitado: la inmutabilidad append-only elimina riesgos de manipulación de registros clínicos
- Minimización de superficie de ataque: pseudonimización en capa de persistencia reduce impacto de posibles brechas de BD
- Latencia predecible: sin saltos de red a brokers externos, el procesamiento cognitivo es síncrono y medible
- Escalabilidad vertical clara: Go + PostgreSQL permite crecimiento hasta límites conocidos sin re-arquitectura
- Verificación forense independiente: el encadenamiento prev_hash permite auditorías externas sin acceso al código fuente

### Negativas

- Monolito modular limita despliegue independiente de componentes; un fallo en el motor cognitivo afecta todo el sistema
- Ausencia de cola de mensajes externa implica que picos de carga pueden saturar el procesamiento síncrono de eventos
- PostgreSQL como único store de eventos limita capacidades de replay temporal y análisis complejo sobre streams
- La pseudonimización SHA-256 sin salt por-paciente puede ser vulnerable a ataques de diccionario sobre IDs conocidos
- El acoplamiento temporal entre HTTP request y generación de evento cognitivo añade latencia al response time del usuario
- Migración futura a arquitectura event-driven distribuida requerirá refactor significativo de la capa de repositorios
- El encadenamiento forense introduce dependencia serial en las escrituras (limita paralelismo en alta concurrencia)

## Referencias

- `pkg/triage/middleware.go` (Cognitive Firewall: función ValidateCareEvent)
- `pkg/triage/triage.go` (triage clínico Celeste/Amarillo/Rojo)
- `pkg/triage/events.go` (definición de CognitiveStateValidated y CareEventBlockedByCognitiveState)
- `pkg/triage/errors.go` (tipos CognitiveBlockReason y mapeo desde errores)
- `pkg/repository/audit_postgres.go` (repositorio append-only)
- `db/migrations/001_create_cognitive_audit_log.sql` (schema de tabla de auditoría)
- `docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md`
- `docs/architecture/decisions/ADR-007-fhir-extensions.md`
