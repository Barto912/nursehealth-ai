# ADR-002: Elección de Go como Lenguaje Backend Principal

## Estado

Aceptado

## Contexto

NurseHealth-AI requiere un lenguaje backend que cumpla con los siguientes requisitos técnicos y operacionales:

- Latencia predecible para procesamiento síncrono de eventos cognitivos (Cognitive Firewall, Triage System)
- Concurrencia eficiente para manejar múltiples requests clínicos simultáneos sin degradación de rendimiento
- Binarios estáticos para despliegue simple en entornos sanitarios con restricciones de infraestructura
- Integración robusta con PostgreSQL (driver maduro, connection pooling, transacciones ACID)
- Equipo pequeño en fase MVP que requiere productividad sin sacrificar rendimiento
- Filosofía "Gentle Care" que prioriza simplicidad operacional y mantenibilidad a largo plazo

El sistema procesa eventos cognitivos de forma síncrona (sin colas de mensajes externas en MVP), por lo que el tiempo de respuesta del backend impacta directamente en la experiencia del profesional de enfermería.

## Decisión

Se selecciona Go como lenguaje backend principal sobre las alternativas evaluadas.

### Alternativas evaluadas

**Java (Spring Boot)**

Descartado por:
- Overhead de memoria JVM significativo (típicamente 256MB-512MB baseline) incompatible con despliegues en entornos sanitarios con recursos limitados
- Tiempo de arranque lento (5-15 segundos) problemático para escalado horizontal y deployments frecuentes
- Complejidad de configuración XML/annotations que añade fricción para equipo pequeño
- Ecosistema de dependencias pesado (Maven/Gradle) que complica gestión de vulnerabilidades en entornos regulados

**Python (FastAPI/Django)**

Descartado por:
- Rendimiento insuficiente para procesamiento síncrono de eventos cognitivos bajo carga concurrente
- Global Interpreter Lock (GIL) limita paralelismo real en operaciones CPU-bound (hashing SHA-256, validaciones)
- Tipado dinámico incrementa riesgo de errores en runtime en sistema crítico sanitario
- Ausencia de compilación a binario estático complica despliegue en entornos con restricciones de instalación

**Node.js (Express/Fastify)**

Considerado pero descartado por:
- Modelo single-threaded con event loop puede saturarse con operaciones CPU-bound (pseudonimización SHA-256 masiva)
- Ecosistema JavaScript/TypeScript con tipado menos estricto que Go para sistemas sanitarios
- Menor madurez en drivers PostgreSQL con garantías ACID comparado con ecosistema Go
- Curva de aprendizaje similar a Go, pero sin ventajas de rendimiento en este caso de uso

**Rust**

Considerado pero descartado por:
- Curva de aprendizaje pronunciada (ownership, borrowing, lifetimes) incompatible con timeline de MVP y equipo pequeño
- Tiempo de compilación excesivo (minutos vs segundos en Go) reduce productividad en desarrollo iterativo
- Ecosistema de librerías menos maduro para integración PostgreSQL comparado con Go
- Overkill para requisitos de rendimiento del sistema (Go ofrece latencia sub-milisegundo suficiente)

### Razones de selección de Go

- Compilación a binario estático único sin dependencias runtime, facilitando despliegue en entornos sanitarios restringidos
- Goroutines y channels proporcionan concurrencia eficiente con overhead mínimo (2KB por goroutine vs 1MB en threads Java)
- Tiempo de arranque sub-segundo permite escalado horizontal rápido y deployments zero-downtime
- Tipado estático con inferencia balancea seguridad de tipos y productividad de desarrollo
- Integración con PostgreSQL vía librería estándar `database/sql` mediante la interfaz `DBExecutor` (ver `pkg/repository/audit_postgres.go`), manteniendo dependencias externas al mínimo en `go.mod`
- Curva de aprendizaje moderada (sintaxis minimalista, ~50 keywords) permite onboarding rápido de equipo pequeño
- Herramientas integradas (go fmt, go vet, go test) estandarizan calidad de código sin configuración adicional
- Gestión de dependencias nativa (go modules) simplifica auditoría de seguridad en entornos regulados

## Consecuencias

### Positivas

- Despliegue simplificado: binario único de ~10-20MB sin necesidad de JVM, runtime Python o Node.js en servidores de producción
- Latencia predecible: garbage collector concurrente de Go minimiza pausas (típicamente <1ms) vs stop-the-world de Java
- Concurrencia nativa eficiente: procesamiento paralelo de eventos cognitivos sin overhead de threads del SO
- Productividad de equipo pequeño: código verboso pero explícito reduce bugs, herramientas integradas estandarizan calidad
- Seguridad de tipos: compilador detecta errores en build time, crítico para sistema sanitario sin margen de error
- Integración PostgreSQL verificada: interfaz `DBExecutor` con método `ExecContext` abstrae el driver subyacente, permitiendo swaps sin afectar dominio
- Observabilidad nativa: package net/http/pprof y tracing distribuido facilitan debugging en producción
- Cross-compilation trivial: generación de binarios para Linux/Windows/macOS desde single codebase sin toolchains adicionales
- Superficie de dependencias mínima: `go.mod` reducido facilita auditoría de seguridad y reduce riesgo de vulnerabilidades transitivas

### Negativas

- Verbosidad en manejo de errores (if err != nil) incrementa boilerplate vs excepciones de Java/Python
- Ausencia de generics avanzados (hasta Go 1.18 limitados) puede requerir código duplicado para operaciones genéricas
- Ecosistema de librerías menos extenso que Java/Python para casos de uso edge (aunque suficiente para este proyecto)
- Curva de aprendizaje de concurrencia (goroutines, channels, context) requiere disciplina para evitar leaks y deadlocks
- Garbage collector introduce latencia no determinista en microsegundos, problemático solo si se requieren garantías hard real-time (no aplica aquí)
- Menor productividad inicial vs Python para prototipado rápido (trade-off aceptado por rendimiento y seguridad en producción)
- Debugging de race conditions requiere herramientas específicas (go race detector) y disciplina de testing

## Referencias

- `pkg/repository/audit_postgres.go` (integración PostgreSQL con interfaz `DBExecutor`)
- `internal/cognitive/firewall.go` (procesamiento concurrente de eventos) — nota: implementación actual vive en `pkg/triage/middleware.go`
- `cmd/server/main.go` (binario estático compilado)
- `go.mod` (gestión de dependencias mínima)
- `docs/architecture/decisions/ADR-001-arquitectura-general.md` (arquitectura general)
