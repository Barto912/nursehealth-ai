# Índice de Architecture Decision Records (ADRs)

Este directorio contiene los Architecture Decision Records (ADRs) de NurseHealth-AI.
Cada ADR documenta una decisión arquitectónica significativa, su contexto, alternativas
evaluadas, consecuencias y estado.

## Convención

- Formato: MADR (Markdown Architecture Decision Records)
- Numeración: secuencial e inmutable (un ADR aceptado no se renumera)
- Estado: `Aceptado` / `Accepted` según idioma del documento
- Un ADR nuevo se agrega al final; las decisiones que reemplazan a otra se marcan como `Superseded` en el ADR original

## Tabla de ADRs

| ADR | Título | Estado | Archivo |
|-----|--------|--------|---------|
| 001 | Arquitectura General del Sistema NurseHealth-AI | Aceptado | [ADR-001-arquitectura-general.md](ADR-001-arquitectura-general.md) |
| 002 | Elección de Go como Lenguaje Backend Principal | Aceptado | [ADR-002-eleccion-go-lenguaje-backend.md](ADR-002-eleccion-go-lenguaje-backend.md) |
| 003 | Patrón de Repositorio Audit-Only Append-Only | Aceptado | [ADR-003-repositorio-audit-only-append-only.md](ADR-003-repositorio-audit-only-append-only.md) |
| 004 | Pseudonimización SHA-256 de Datos Sensibles | Aceptado | [ADR-004-pseudonimizacion-sha256-datos-sensibles.md](ADR-004-pseudonimizacion-sha256-datos-sensibles.md) |
| 005 | Encadenamiento Forense de Auditoría mediante prev_hash | Aceptado | [ADR-005-encadenamiento-forense-prev-hash.md](ADR-005-encadenamiento-forense-prev-hash.md) |
| 006 | Human Firewall as Traceability Precondition (MVP) | Accepted | [ADR-006-humanistic-firewall-mvp.md](ADR-006-humanistic-firewall-mvp.md) |
| 007 | FHIR R4 Extensions for Patient-Side Cognitive Gates (MVP) | Accepted | [ADR-007-fhir-extensions.md](ADR-007-fhir-extensions.md) |
| 008 | Marco Técnico-Legal y Compliance por Diseño | Aceptado (Ref. Interna) | [ADR-008-legal-and-pharmaceutical-compliance.md](ADR-008-legal-and-pharmaceutical-compliance.md) |

## Relaciones entre ADRs

- ADR-001 define la arquitectura general; ADR-002 justifica el lenguaje; ADR-003 el patrón de persistencia
- ADR-004 (pseudonimización) y ADR-005 (encadenamiento forense) refuerzan las garantías de integridad de ADR-003
- ADR-006 y ADR-007 especifican el Cognitive Firewall y la interoperabilidad FHIR sobre la base de ADR-001

## Cómo agregar un ADR nuevo

1. Copiar la estructura MADR de un ADR existente
2. Asignar el siguiente número secuencial
3. Completar Contexto, Decisión, Consecuencias (positivas y negativas) y Estado
4. Agregar la fila a la tabla de este índice
5. Referenciar archivos reales del repositorio (verificados, no inventados)
