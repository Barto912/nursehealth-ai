AVISO DE TRANSPARENCIA Y ALCANCE PROFESIONAL:
Borrador técnico preparado por el equipo de ingeniería de NurseHealth.ai con asistencia de IA (ARGON) y auditado por LEO. No constituye dictamen jurídico vinculante hasta su adopción y firma digital por un abogado matriculado.

NOTA DE AUDITORÍA TÉCNICA (SANIDAD DE CÓDIGO LEO-ARGON):
En cumplimiento estricto de los ADRs 001 a 008 de la rama main, se ha depurado el esquema base para erradicar desviaciones de arquitectura: se elimina toda referencia a lenguajes ajenos al stack oficial (100% Go en monolito modular) y se etiqueta reglamentariamente el componente critical_list.go bajo la Regla 1 de la Lista Negra.

1. ARQUITECTURA DEL MOTOR DETERMINISTA FAIL-CLOSE
El Motor determinista Fail-Close constituye la barrera primaria de seguridad clínica y jurídica. Su diseño responde a un paradigma de ejecución síncrona, determinista y de rechazo por defecto ante cualquier ambigüedad.
- Núcleo de Ejecución (pkg/triage/middleware.go): La función ValidateCareEvent actúa como interceptor central. Ningún evento alcanza la persistencia sin su validación afirmativa.
- Acoplamiento de Contexto (pkg/triage/patient_context.go, events.go, errors.go): Evalúa la estructura del evento contra el estado fisiológico y alergias del paciente. Ante violación de invariantes, instancia una excepción tipificada que aborta la ejecución.
- Flujo: Entrada (Payload firmado por Titan M3) → Validación Estructural → Evaluación de Umbral y Guardrails (pkg/cognitive/ruido_check.go) → Decisión (Commit o Veto Fail-Close).

2. UMBRALES DETERMINISTAS DE VETO Y ESCALACIÓN
- Incertidumbre Algorítmica (≥ 1%): Si la confianza es < 99%, ValidateCareEvent ejecuta un bloqueo inmediato y exige evidencia suplementaria humana.
- Intercepción de Fármacos de Alta Vigilancia: Toda operación con medicamentos de estrecho margen terapéutico requiere validación cruzada independiente, gestionada vía critical_list.go (Componente planificado para fase piloto, diseño en docs/domain/logistica-inversa.md).
- Detección de Ruido Cognitivo: Si pkg/cognitive/ruido_check.go detecta fatiga o latencia anómala en el enfermero, bloquea el cierre autónomo y fuerza la escalación jerárquica.

3. MATRIZ DE DECISIONES DEL MOTOR
| Condición Detectada | Umbral / Disparador | Acción del Motor | Destinatario | Componente Responsable |
|---|---|---|---|---|
| Confianza < 99% | Incertidumbre ≥ 1% | Bloqueo + Exigencia de foto de respaldo | Enfermero en Domicilio | pkg/triage/middleware.go, errors.go |
| Fármaco Crítico | Principio activo de alto riesgo | Bloqueo preventivo + Alerta doble chequeo | Farmacéutico Remoto | middleware.go + critical_list.go (Planificado) |
| Ruido Cognitivo / Fatiga | Anomalía conductual detectada | Bloqueo de acto + Escalación para relevo | Supervisor Médico | pkg/cognitive/ruido_check.go, ruido_gate.go |
| Negativa de Entrega | Código REFUSAL_TO_SURRENDER_WASTE | Bloqueo de cierre + Acta de exención | Auditoría Legal / Dirección Médica | pkg/triage/events.go, audit_postgres.go |

4. REGISTRO ATDP (FLIGHT RECORDER) Y PERSISTENCIA FORENSE
Cada decisión se consolida en la tupla inmutable ATDP = (o, h, a, y, r, m):
- o: Observación clínica (pkg/triage/patient_context.go).
- h: Hipótesis o propuesta del modelo en el borde.
- a: Acción determinista (Aprobación / Bloqueo) en middleware.go.
- y: Resultado clínico o resolución humana (ej. REFUSAL_TO_SURRENDER_WASTE).
- r: Señal de auditoría de riesgo (pkg/cognitive/ruido_check.go).
- m: Metadatos criptográficos (Timestamp, ID Pixel 11 Pro, firma Titan M3).
Persistencia: Almacenamiento exclusivo append-only con encadenamiento forense SHA-256 en pkg/repository/audit_postgres.go (db/migrations/001_create_cognitive_audit_log.sql).

5. INMUTABILIDAD DE GUARDRAILS Y AISLAMIENTO
- Gobernanza Estricta en 100% Go: Mecanismos de validación integrados sin depender de binarios externos no autorizados.
- Exclusión de Ciclos de Autoevolución: El código de triage y ruido_check es estático y determinista. Ningún proceso de ML tiene permisos de escritura sobre los umbrales de ValidateCareEvent.
- Guardrails Éticos Inmodificables: Las reglas del Humanistic Firewall y códigos como REFUSAL_TO_SURRENDER_WASTE están compilados de forma estática, impidiendo su desactivación por inyección de prompts.

6. CASOS DE PRUEBA Y VALIDACIÓN DETERMINISTA
- Escenario 1 (Incertidumbre ≥ 1%): Mapeo SNOMED CT con 98.4% de confianza. Resultado: ValidateCareEvent rechaza la transacción, solicita confirmación manual con foto y asienta el bloqueo en audit_postgres.go.
- Escenario 2 (Fármaco Crítico): Dosis fuera de rango en patient_context.go o en critical_list.go (Planificado). Resultado: Bloqueo inmediato en interfaz, alerta al Farmacéutico Remoto y grabación de tupla ATDP.
- Escenario 3 (Deterioro Cognitivo + Fricción): Negativa de entrega con indicadores de estrés en ruido_gate.go. Resultado: Bloqueo de cierre ordinario, guía de contención humanística, declaración REFUSAL_TO_SURRENDER_WASTE firmada por Titan M3 y escala a Supervisor Médico.
