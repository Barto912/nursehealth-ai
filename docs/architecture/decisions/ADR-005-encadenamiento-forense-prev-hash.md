# ADR-005: Encadenamiento Forense de Auditoría mediante prev_hash

## Estado

Aceptado

## Contexto

NurseHealth-AI persiste eventos cognitivos en la tabla `cognitive_audit_log` definida en `db/migrations/001_create_cognitive_audit_log.sql`. El patrón append-only documentado en ADR-003 garantiza inmutabilidad a nivel de aplicación, pero no proporciona evidencia criptográfica de que las filas no han sido modificadas tras su inserción.

Un atacante con acceso privilegiado a la base de datos podría alterar el contenido de una fila (ej: modificar `gate_snapshot`, `block_reason` o `decision`) sin que el sistema detecte la manipulación, ya que PostgreSQL no ofrece garantías nativas de integridad criptográfica sobre el contenido de las celdas.

Los axiomas del sistema documentados en `docs/philosophy/axioms.json` (AX-001 a AX-006) establecen que la trazabilidad de decisiones clínicas asistidas por IA debe ser forensemente verificable por terceros independientes. Esto requiere un mecanismo que permita a un auditor externo detectar cualquier alteración posterior al momento de la inserción.

La tabla `cognitive_audit_log` incluye dos columnas diseñadas para este propósito:

- `payload_hash` (CHAR(64)): hash SHA-256 del contenido canónico del evento
- `prev_hash` (CHAR(64)): hash de la fila anterior en la cadena

Estas columnas permiten construir una cadena tamper-evident análoga a la estructura de bloques de un ledger, donde cada entrada referencia criptográficamente a su predecesora.

## Decisión

Se adopta un mecanismo de encadenamiento forense basado en hashes criptográficos que vincula cada fila de `cognitive_audit_log` con su predecesora, permitiendo verificación de integridad por auditores externos.

### Mecanismo de encadenamiento

Cada inserción en `cognitive_audit_log` calcula dos hashes:

**1. `payload_hash`:** SHA-256 sobre la representación canónica del evento. El JSON canónico incluye los campos del evento en orden alfabético de claves, sin espacios ni saltos de línea innecesarios, codificado en UTF-8. Los campos incluidos en el cálculo son: `care_event_id`, `patient_hash`, `operator_id`, `block_reason`, `gate_snapshot`, `decision` y `receipt_json`. 

Los campos `created_at` e `id` se excluyen del `payload_hash` porque son generados por PostgreSQL (`DEFAULT NOW()` y `BIGSERIAL` respectivamente) y no son deterministas desde la aplicación.

**2. `prev_hash`:** SHA-256 sobre la concatenación del `payload_hash` y `prev_hash` de la fila inmediatamente anterior en la cadena, ordenada por `id` (BIGSERIAL). 

Para la primera fila de la tabla (genesis), `prev_hash` se establece a un valor constante documentado: **64 caracteres `'0'`** (`0000000000000000000000000000000000000000000000000000000000000000`).

El orden de cálculo es estricto: primero se computa `payload_hash` del evento actual, luego se recupera el par (`payload_hash`, `prev_hash`) de la última fila insertada, y finalmente se calcula `prev_hash` del evento actual como `SHA-256(payload_hash_anterior || prev_hash_anterior)`.

### Procedimiento de verificación de integridad

Un auditor externo puede verificar la integridad de la cadena siguiendo este procedimiento:

1. Obtener todas las filas de `cognitive_audit_log` ordenadas por `id` ASC
2. Para la primera fila (genesis): verificar que `prev_hash` coincide con el valor constante `0000...0000` (64 ceros); recalcular `payload_hash` a partir de los campos del evento y comparar con el valor almacenado
3. Para cada fila subsiguiente (i): verificar que `prev_hash[i]` coincide con `SHA-256(payload_hash[i-1] || prev_hash[i-1])`; recalcular `payload_hash[i]` a partir de los campos y comparar con el valor almacenado
4. Si todas las verificaciones pasan, la cadena es íntegra desde el punto de vista criptográfico

Este procedimiento puede ejecutarse de forma independiente sin acceso al código fuente de la aplicación, únicamente requiriendo el schema de la tabla y el valor constante de genesis.

### Comportamiento ante detección de ruptura de cadena

Si durante la verificación se detecta una discrepancia (ya sea en `payload_hash` o en `prev_hash`), se considera incidente de seguridad de nivel crítico:

- Se genera alerta inmediata al equipo de seguridad y cumplimiento normativo
- Se inicia investigación forense para determinar el alcance de la manipulación
- Se documenta el incidente en el log de seguridad del sistema
- Se notifica a las autoridades competentes si la manipulación afecta registros clínicos de pacientes (requisito HIPAA breach notification)
- Las filas posteriores a la ruptura se consideran comprometidas hasta que se demuestre lo contrario mediante análisis independiente

### Límites del mecanismo

El encadenamiento forense proporciona garantías específicas y tiene límites explícitos:

**Protege contra:**
- Modificación del contenido de campos (`gate_snapshot`, `block_reason`, `decision`, `receipt_json`, `patient_hash`, `operator_id`, `care_event_id`) en filas existentes
- Inserción de filas falsas en medio de la cadena (rompería `prev_hash` de filas subsiguientes)
- Reordenamiento de filas (rompería la secuencia de `prev_hash`)

**NO protege contra:**
- Eliminación de filas completas: un atacante con privilegios DELETE puede eliminar filas y recalcular `prev_hash` de la fila siguiente para mantener la cadena consistente. Esta protección requiere refuerzo a nivel de base de datos (triggers anti-DELETE, replicación inmutable, o WORM storage)
- Colusión con acceso total a la base de datos: si un atacante tiene privilegios para modificar cualquier fila y puede recalcular hashes, puede construir una cadena alternativa consistente que pase la verificación. La detección requiere comparación contra réplicas externas o logs de escritura inmutables fuera de PostgreSQL
- Manipulación del valor de genesis: si el atacante puede modificar la primera fila y su `prev_hash` constante, toda la cadena subsiguiente puede reconstruirse. Mitigación: el valor de genesis debe publicarse en un medio externo inmutable (ej: firma en blockchain público o notario digital)
- Ataques al JSON canónico: si la implementación de serialización canónica tiene ambigüedades (orden de claves en objetos anidados de `gate_snapshot`, representación de números, manejo de nulls), dos implementaciones pueden producir hashes diferentes para el mismo contenido lógico. Mitigación: especificación estricta del formato canónico documentada y testeada
- Manipulación del campo `created_at`: al no incluirse en `payload_hash`, un atacante con privilegios UPDATE puede modificar timestamps sin romper la cadena. Esto es aceptable porque `created_at` es metadata de auditoría secundaria, no contenido clínico

## Consecuencias

### Positivas

- Verificación forense independiente: auditores externos pueden verificar integridad de la cadena sin acceso al código fuente de la aplicación, únicamente con el schema y el valor de genesis
- Detección de manipulación de contenido: cualquier alteración de `gate_snapshot`, `block_reason`, `decision` u otros campos clínicos rompe la cadena y es detectable
- Evidencia criptográfica de secuencia: el encadenamiento prueba que los eventos ocurrieron en el orden registrado por `id`, previniendo reordenamiento temporal
- Cumplimiento de axiomas AX-001 a AX-006: el mecanismo proporciona la verificabilidad forense exigida por la filosofía del sistema documentada en `docs/philosophy/axioms.json`
- Compatibilidad con ADR-003: el encadenamiento complementa el patrón append-only añadiendo garantías criptográficas sobre la inmutabilidad a nivel de aplicación
- Auditoría de cumplimiento normativo: reguladores HIPAA/GDPR pueden solicitar verificación de integridad y el sistema puede demostrarla mediante el procedimiento documentado
- Detección de bugs de aplicación: si un error en el código produce hashes inconsistentes, se detecta inmediatamente en la siguiente verificación, facilitando debugging

### Negativas

- No protege contra eliminación de filas: el mecanismo asume que todas las filas permanecen en la tabla; un atacante con DELETE puede eliminar filas y reconstruir la cadena. Requiere refuerzo adicional (triggers anti-DELETE, replicación inmutable)
- Vulnerable a colusión con acceso total: si un atacante controla completamente PostgreSQL, puede modificar filas y recalcular hashes, construyendo una cadena falsa consistente. Mitigación requiere réplicas externas o logs inmutables fuera de la BD
- Overhead computacional por inserción: cada evento requiere dos cálculos SHA-256 adicionales (payload_hash y prev_hash) más una lectura de la última fila para obtener sus hashes. En throughput alto (>1000 eventos/segundo), esto puede convertirse en bottleneck
- Complejidad de verificación: el procedimiento de auditoría requiere recorrer toda la tabla secuencialmente; para tablas con millones de filas, la verificación completa puede tomar minutos u horas
- Dependencia de serialización canónica estricta: cualquier ambigüedad en la representación JSON de `gate_snapshot` (orden de claves, manejo de nulls, precisión de floats) produce hashes inconsistentes entre inserción y verificación. Requiere especificación rigurosa y testeada
- Primer evento (genesis) es punto único de fallo: si se compromete la fila de genesis, toda la cadena puede reconstruirse falsamente. Mitigación requiere publicación externa del valor de genesis
- Migración de esquema compleja: si se añaden columnas a `cognitive_audit_log` en el futuro, el cálculo de `payload_hash` debe actualizarse, pero las filas antiguas mantendrán el hash calculado con el schema anterior. Requiere estrategia de versionado de hashes o recálculo masivo
- Latencia añadida en escritura: la necesidad de leer la última fila para calcular `prev_hash` introduce una dependencia serial en las inserciones, limitando paralelismo. En escenarios de alta concurrencia, esto puede requerir locking optimista o serialización explícita

## Referencias

- `db/migrations/001_create_cognitive_audit_log.sql` (schema de la tabla con columnas payload_hash y prev_hash)
- `pkg/repository/audit_postgres.go` (implementación de persistencia append-only)
- `pkg/triage/events.go` (definición de CognitiveStateValidated y estructura del evento)
- `docs/philosophy/axioms.json` (axiomas AX-001 a AX-006 que exigen verificabilidad forense)
- `docs/architecture/decisions/ADR-003-repositorio-audit-only-append-only.md` (patrón append-only complementado por este encadenamiento)
- `docs/architecture/decisions/ADR-004-pseudonimizacion-sha256-datos-sensibles.md` (PatientHash como campo incluido en payload_hash)
