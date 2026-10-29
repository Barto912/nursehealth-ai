# ADR-008: Marco Técnico-Legal y Compliance por Diseño (Sanitario, Ambiental, Datos y Cloud)

## Estado

Aceptado como referencia técnica interna. Validez legal sujeta a revisión, corrección y firma de profesional matriculado (dictamen externo). Este documento no constituye asesoramiento legal.

## Contexto

NurseHealth-AI opera en la intersección de cuatro dominios regulatorios: sanitario (ANMAT, Ley 26.529), ambiental (Ley 24.051, Ley PBA 11.347, Ley CABA 154), protección de datos (Ley 25.326, Convenio 108+) y firma digital/cloud (Ley 25.506, Ley 27.706). El piloto Canning (diciembre 2026) y la auditoría del Google for Startups Cloud Program (AI-First Tier) exigen un marco consolidado de compliance por diseño.

### Origen del documento y resolución de defectos

El borrador original fue generado con asistencia de IA y revisado internamente. Se detectaron y resolvieron cuatro defectos antes de consolidar esta versión:

- Duplicidad de versiones contradictorias en el borrador: se consolida una única versión (este documento)
- Evento de negativa del paciente con dos nombres distintos: se unifica como `REFUSAL_TO_SURRENDER_WASTE`
- Terminología "Event Sourcing SHA-256" contradictoria con ADR-003: se reemplaza por registro append-only con encadenamiento forense (ADR-003 y ADR-005)
- Referencias a código inexistente en el repositorio: se mapean a estado real versus planificado (tabla siguiente)

## Decisión

Se adopta el marco técnico-legal de siete bloques (A a G) como referencia de compliance por diseño del sistema, vinculando cada mitigación con componentes reales o planificados del repositorio.

### Mapeo de componentes: estado real versus planificado

| Componente del marco legal | Estado en repositorio | Referencia |
|---|---|---|
| Motor Fail-Close y gates deterministas | Existente | `pkg/triage/middleware.go` (ValidateCareEvent), `pkg/triage/patient_context.go`, `pkg/triage/errors.go` |
| Gate de ruido (Capa 2, gobernanza determinista) | Existente | `pkg/cognitive/ruido_check.go`, `pkg/adapter/ruido_gate.go` |
| Registro append-only con hashes forenses | Existente | `pkg/repository/audit_postgres.go`, `db/migrations/001_create_cognitive_audit_log.sql` |
| Eventos de auditoría inmutables | Existente | `pkg/triage/events.go` |
| Clasificadores de residuos y fallback óptico (`classifier.go`, `critical_list.go`, `biological_classifier.go`, `sharps_handler.go`, `barcode_fallback.go`) | Planificado (fase piloto) | Diseño en `docs/domain/logistica-inversa.md` |
| Matrices CSV de clasificación (`medicamentos_criticos`, `clasificacion_biologicos`, `clasificacion_farmaceuticos`, `especificaciones_cortopunzantes`) | Planificado (artefactos de configuración) | Diseño en `docs/domain/logistica-inversa.md` |
| Data Proxy y memoria local de borde (Engram) | Planificado (edge) | Propuesta en `docs/gcp-proposal/` |

### Matriz de riesgos consolidada

| Bloque | Riesgo jurídico principal | Marco normativo | Probabilidad × Impacto | Mitigación por diseño |
|---|---|---|---|---|
| A: Responsabilidad civil y consentimiento | Contagio a terceros, repudio del retiro, mala praxis algorítmica | CCyCN Arts. 59, 1743, 1757, 1758, 1768; Ley 26.529 | Media × Crítico | Juicio clínico humano; Fail-Close determinista ante incertidumbre; doble firma digital |
| B: Residuos peligrosos y patogénicos | Secuestro vehicular por falta de habilitación de carga peligrosa | Ley 24.051; Ley PBA 11.347; Ley CABA 154 | Alta × Alto | Encuadre in itinere como acto asistencial domiciliario; doble contención; manifiesto digital |
| C: Residuos farmacéuticos (ANMAT) | Recepción ilegal de patogénicos en farmacias; desvío de psicotrópicos | Ley 16.463; Disp. ANMAT 6693/2010 y 4886/2003 | Media × Alto | Taxonomía de tres niveles; farmacias solo Categoría 1; acta co-firmada |
| D: Datos sensibles versus inmutabilidad | Conflicto append-only versus derecho de supresión (ARCO) | Ley 25.326 (Arts. 2, 7, 8, 9, 14, 16); Convenio 108+ (Ley 27.699); Res. AAIP 47/2018; Ley 26.529 Art. 18 | Media × Crítico | Sanitización en borde; disociación definitiva tras plazo legal de 10 años |
| E: Firma criptográfica dual | Impugnación de autoría o falta de fecha cierta | Ley 25.506; Dec. 182/2019 y 743/2024; CCyCN Art. 288 | Baja × Alto | Firma digital del profesional más firma electrónica robusta del paciente; sellado RFC 3161 |
| F: Armonización interjurisdiccional | Fricción PBA versus CABA; incompatibilidad regional | Ley PBA 11.347 versus Ley CABA 154 | Media × Medio | Regla de la cota más estricta (24 hs); mapeo FHIR R4 y SNOMED CT |
| G: Cloud, SaMD y datasets | Exigencia de registro SaMD; transferencia internacional ilegítima | Disp. ANMAT 64/2025 (Res. GMC 25/21); Disp. DNPDP 60-E/2016; Ley 11.723 | Baja × Crítico | Exención SaMD por finalidad logística; DPA con cláusulas tipo AAIP; solo eventos seudonimizados en nube |

### BLOQUE A — Consentimiento informado y responsabilidad civil

Principio rector: el juicio clínico y la manipulación física de residuos son actos exclusivamente humanos de profesional matriculado; la plataforma provee trazabilidad criptográfica y bloqueo preventivo, sin ejercer guarda material ni dirigir praxis médica.

Fraccionamiento de guarda material: en Etapas 1 a 6 el Enfermero Autónomo es guardián material del residuo (Art. 1758 CCyCN) con responsabilidad subjetiva (Art. 1768 CCyCN); en Etapa 7 la Segunda Firma Digital Dual transfiere la guarda al nodo receptor, liberando a enfermero y plataforma de responsabilidad ulterior por manipulación en acopio o disposición final.

Consentimiento de doble capa (Dual Opt-In): (a) consentimiento informado clínico para el procedimiento domiciliario (Ley 26.529 Art. 5, CCyCN Art. 59) y (b) autorización expresa de retiro y cesión de custodia de residuos domiciliarios y fluidos biológicos, en flujo electrónico inalterable.

Protocolo ante negativa de entrega (`REFUSAL_TO_SURRENDER_WASTE`): el enfermero no puede ejercer coacción en domicilio; activa el protocolo de excepción, el motor determinista exige confinamiento en descartador sellado, toma fotografía probatoria y emite Acta Digital de Deslinde sellada con hash SHA-256, geolocalización y marca temporal de hardware (ATDP). Si el residuo pertenece a la matriz de medicamentos críticos (artefacto planificado), se alerta al Director Médico y Farmacéutico Remoto para auditoría de desvío ante ANMAT.

Eximente por bloqueo determinista: los disclaimers genéricos son nulos si encubren dolo o culpa grave (Arts. 1743 y 988 CCyCN); el blindaje real es la barrera técnica Fail-Close que interrumpe el nexo causal: modelo congelado, bloqueo mecánico ante incertidumbre o fármaco crítico, y escalada obligatoria a humano. El registro ATDP prueba en sede judicial quién tomó la decisión final.

### BLOQUE B — Residuos peligrosos y patogénicos

Encuadre in itinere: las Leyes 24.051, 11.347 (PBA) y 154 (CABA) regulan gravosamente a generadores institucionales y transportistas comerciales a granel. El domicilio particular no es establecimiento sanitario habilitado; el traslado de micro-volúmenes en el vehículo del enfermero encuadra como tramo final accesorio del acto de asistencia sanitaria domiciliaria (prevención ambiental y responsabilidad extendida), no como transporte comercial de cargas peligrosas.

Condiciones de blindaje en tránsito: (1) volumen limitado a la micro-generación de la hoja de ruta diaria del propio profesional; (2) contenedores primarios (Verde, Rojo, Amarillo Rígido, Azul) en baúl, dentro de contenedor secundario rígido, estanco, lavable y con anclaje antiderrame; (3) manifiesto digital en tránsito con QR verificable, hash SHA-256, categorías y destino, exhibible junto a credencial matriculada.

Contención y retención: etiquetado con pictograma normativo y QR sin PII ni PHI en el exterior del envase (Ley 25.326 Arts. 2, 9 y 10); descartadores amarillos bloqueados operativamente al alcanzar dos tercios o tres cuartos de capacidad; retención vehicular máxima de 24 horas con bloqueo de nuevas prestaciones sin descarga acreditada; doble contención sellada en origen para VIH, Hepatitis B/C o quimioterapia domiciliaria.

### BLOQUE C — Residuos farmacéuticos y ANMAT

Taxonomía de tres niveles: Críticos o Alto Riesgo (con bloqueo Fail-Close), Riesgo Medio o zonas grises (AINEs, corticoides tópicos), y Asimilables. Supera la clasificación binaria y se alinea con Ley 16.463 y Disposiciones ANMAT 4886/2003 y 6693/2010, impidiendo que principios activos con impacto ecotóxico o riesgo de reingreso al mercado ilegal terminen en basura común. El protocolo de fallback óptico (ingreso manual de GTIN, serie, lote y vencimiento con umbral de confianza alto) asegura trazabilidad unitaria.

Farmacias Partners: solo reciben Categoría 1 (Contenedor Verde cerrado). Requieren habilitación sanitaria vigente, Dirección Técnica farmacéutica matriculada, y adhesión a programa registrado de recolección o contrato con operador de residuos especiales (Categorías Y2/Y3). Prohibido por diseño en el motor de ruteo derivar a farmacias las Categorías 2 (biológicos) y 3 (cortopunzantes), rutadas exclusivamente a plantas de tratamiento patogénico.

Cierre de auditoría: en Etapa 7 se suscribe Manifiesto Electrónico de Recepción con Segunda Firma Digital Dual; el convenio con farmacias incluye Cláusula de Derecho de Auditoría e Indemnidad que obliga a cargar el Certificado de Tratamiento y Disposición Final de la planta habilitada, cerrando la Etapa 8 ante inspecciones de ANMAT u OPDS.

### BLOQUE D — Datos personales y trazabilidad inmutable

Conflicto: el Art. 16 de la Ley 25.326 otorga rectificación y supresión (5 días hábiles; acceso en 10 días corridos, Art. 14), pero el derecho de supresión no es absoluto en sanidad: el Art. 18 de la Ley 26.529 impone conservación de documentación clínica por mínimo 10 años, plazo coincidente con trazabilidad ambiental (OPDS) y vigilancia sanitaria (ANMAT). Los datos de salud son sensibles (Arts. 2 y 7, Decreto 1558/2001, Convenio 108+ por Ley 27.699).

Reconciliación arquitectónica en tres tiempos: (1) en tiempo real, el Data Proxy y la memoria local de borde censuran PII y PHI antes de transmitir; el registro append-only solo almacena hashes y metadatos seudonimizados, coherente con ADR-004; (2) ante rectificación, no se sobrescribe historia: se añade evento compensatorio `RECTIFICATION_EVENT` firmado, alineado al mecanismo de corrección de ADR-003; (3) ante supresión o vencimiento del plazo legal, se ejecuta disociación definitiva (Art. 2 y Art. 4 inc. 4 Ley 25.326) mediante borrado de la tabla de vinculación de identidad o crypto-shredding de claves: los hashes remanentes dejan de ser datos personales y conservan integridad matemática y valor probatorio para ANMAT u OPDS.
