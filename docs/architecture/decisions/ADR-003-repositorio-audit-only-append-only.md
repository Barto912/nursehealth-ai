# ADR-003: Patrón de Repositorio Audit-Only Append-Only

## Estado

Aceptado

## Contexto

NurseHealth-AI opera en dominio sanitario donde cada decisión asistida por IA debe ser trazable, auditable e inmutable. Los requisitos normativos (HIPAA, GDPR) exigen que ningún registro clínico generado por el sistema pueda ser alterado o eliminado tras su creación.

El sistema genera eventos cognitivos (CognitiveEvent) que capturan triage, alertas, recomendaciones y validaciones del Cognitive Firewall. Estos eventos constituyen evidencia forense de la interacción entre profesional de enfermería e IA asistencial.

Un repositorio CRUD tradicional (Create-Read-Update-Delete) introduce riesgos inaceptables:

- Posibilidad técnica de modificar eventos cognitivos post-escritura, comprometiendo integridad forense
- DELETE físico elimina evidencia de decisiones clínicas previas
- Auditorías externas requieren garantía de que registros no han sido alterados
- Filosofía "Gentle Care" exige que el profesional pueda confiar en que el historial de asistencias IA es fiel a lo ocurrido

El backend en Go debe persistir estos eventos en PostgreSQL con garantías de inmutabilidad a nivel de aplicación, complementando (no reemplazando) las constraints de integridad de la base de datos.

## Decisión

Se adopta el patrón de repositorio audit-only append-only para toda persistencia de eventos cognitivos y registros clínicos asistidos.

### Definición del patrón

El repositorio expone exclusivamente operaciones de inserción (INSERT). No se implementan métodos de actualización (UPDATE) ni eliminación (DELETE) en la capa de aplicación. Cada evento cognitivo se escribe una sola vez y permanece inmutable durante todo su ciclo de vida.

La interfaz `DBExecutor` expone el método `ExecContext` como mecanismo único de escritura. Esta interfaz abstrae la conexión a PostgreSQL y garantiza que todas las operaciones de persistencia pasen por un punto de control centralizado donde se valida la naturaleza append-only de cada operación.

### Implementación en pkg/repository/audit_postgres.go

El repositorio de auditoría implementa las siguientes garantías:

- Cada inserción de CognitiveEvent incluye timestamp generado por servidor (`DEFAULT NOW()` en la tabla `cognitive_audit_log`)
- La capa de persistencia calcula `payload_hash` (SHA-256 del JSON canónico del evento) antes de la inserción — ver columna `payload_hash` en `db/migrations/001_create_cognitive_audit_log.sql`
- El repositorio no expone métodos Update ni Delete en su interfaz pública
- Las queries de lectura (SELECT) son las únicas operaciones de consulta permitidas
- La conexión a PostgreSQL se gestiona mediante `DBExecutor.ExecContext` para todas las escrituras

### Diferenciación con event sourcing clásico

Este patrón no es event sourcing completo. No se implementa:

- Proyecciones materializadas actualizables
- Snapshotting de estado
- Replay de eventos para reconstrucción de estado
- Separación entre write model y read model

El patrón es estrictamente audit-only: los eventos se escriben para trazabilidad y cumplimiento normativo, no como mecanismo de reconstrucción de estado de aplicación.

### Mecanismo de corrección de errores

Cuando un evento cognitivo contiene información errónea (ej: triage incorrecto por bug), no se modifica el evento original. Se inserta un nuevo evento de tipo "corrección" que referencia el evento original y contiene la información correcta. Esto mantiene la cadena de auditoría intacta y documenta el proceso de corrección.

### Encadenamiento forense (prev_hash)

Cada fila de `cognitive_audit_log` almacena `prev_hash` (hash de la entrada anterior), formando una cadena tamper-evident: alterar cualquier fila histórica rompe la cadena desde ese punto. El procedimiento completo de verificación se documenta en ADR-005.

## Consecuencias

### Positivas

- Integridad forense garantizada: ningún evento clínico puede ser alterado post-escritura, cumpliendo requisitos HIPAA/GDPR de inmutabilidad de registros sanitarios
- Auditoría simplificada: auditores externos pueden verificar que la base de datos no contiene operaciones UPDATE/DELETE sobre eventos cognitivos mediante revisión de queries y logs de PostgreSQL
- Trazabilidad completa de correcciones: los errores se documentan como nuevos eventos que referencian el original, preservando el contexto de qué se corrigió, cuándo y por qué
- Simplificación del modelo de datos: al no existir estado mutable, no se requieren mecanismos de control de concurrencia optimista (versiones, locks) para escrituras
- Recuperación ante desastres facilitada: el log append-only es inherentemente compatible con backup incremental y point-in-time recovery de PostgreSQL
- Debugging de producción mejorado: cada evento cognitivo contiene contexto completo (timestamp, hash, metadata) sin riesgo de que haya sido modificado
- Cumplimiento normativo verificable: la arquitectura append-only es evidence-by-design para reguladores sanitarios
- Verificación forense independiente: el encadenamiento prev_hash permite auditorías externas sin acceso al código fuente

### Negativas

- Crecimiento indefinido del almacenamiento: al no eliminarse registros, el volumen de datos crece monótonamente, requiriendo estrategia de archivado o particionado por tiempo
- Queries de estado actual más complejas: obtener el "estado vigente" de un paciente requiere filtrar eventos por tipo y ordenar por timestamp, no es un simple SELECT directo
- Imposibilidad de corregir datos sensibles erróneos: si se persiste accidentalmente información incorrecta de PatientID pseudonimizado, no puede eliminarse sin violar el patrón (solo puede insertarse evento de invalidación)
- Overhead de almacenamiento por duplicación: eventos de corrección contienen información redundante con los eventos originales que referencian
- Migración de esquema compleja: cambiar la estructura de CognitiveEvent requiere estrategia de versionado de eventos (event versioning) ya que eventos antiguos no pueden actualizarse
- Testing de scenarios edge más verboso: tests deben verificar que operaciones UPDATE/DELETE son rechazadas, no solo que INSERT funciona correctamente
- Performance degradation a largo plazo: tablas con millones de registros append-only requieren indexación cuidadosa y mantenimiento de VACUUM en PostgreSQL

## Referencias

- `pkg/repository/audit_postgres.go` (implementación del repositorio append-only)
- `pkg/repository/audit_postgres_test.go` (tests del repositorio con fake DBExecutor)
- `pkg/triage/middleware.go` (Cognitive Firewall que genera eventos a persistir)
- `pkg/triage/events.go` (definición de CognitiveEvent inmutable)
- `db/migrations/001_create_cognitive_audit_log.sql` (schema de tabla de auditoría con columnas payload_hash y prev_hash)
- `docs/architecture/decisions/ADR-001-arquitectura-general.md` (arquitectura general del sistema)
- `docs/architecture/decisions/ADR-005-encadenamiento-forense-prev-hash.md` (mecanismo de encadenamiento)
- `docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md` (Cognitive Firewall genera estos eventos)
