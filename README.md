# NurseHealth.ai — Enfermería de Precisión Cívico-Militar y Bioseguridad Soberana

[![Versión](https://img.shields.io/badge/Release-v0.6.0--defensa--civil-0A3161?style=for-the-badge)](https://github.com/Barto912/nursehealth-ai/releases/tag/v0.6.0-defensa-civil)
[![Go Stdlib Only](https://img.shields.io/badge/Go_1.24-100%25_Stdlib_Only-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Soberanía](https://img.shields.io/badge/Jurisdicción-República_Argentina-74ACDF?style=for-the-badge)](https://www.argentina.gob.ar/aaip/datospersonales)

## 1. Misión Institucional

**NurseHealth.ai** es una **startup de iniciativa privada** que desarrolla infraestructura de software de misión crítica para la **Enfermería de Precisión Cívico-Militar**. 

Actuamos como socio estratégico B2G de organismos de salud, defensa y seguridad del Estado Nacional Argentino, proveyendo trazabilidad forense inmutable y bioseguridad para enfrentar **nuevas amenazas y patologías sociales derivadas de la criminalidad transnacional** (desvío de fármacos, ruptura de cadena de frío en zonas de riesgo y gestión de residuos patogénicos).

A diferencia de soluciones comerciales basadas en IA generativa, nuestro núcleo de decisión **prescinde deliberadamente de modelos probabilísticos**. Garantizamos que cada acción clínica o logística esté respaldada por prueba matemática de autorización humana y cadena de custodia inalterable.

---

## 2. Los 5 Axiomas Blindables (*Magnifica Humanitas*)

Nuestra doctrina operativa se reduce a **5 invariantes de dominio verificables mediante criptografía**:

| Código | Axioma | Implementación Criptográfica |
| :--- | :--- | :--- |
| **AX-001** | **Dignidad del Paciente** | Cifrado autenticado **AES-256-GCM** con nonce aleatorio antes de persistir datos sensibles. |
| **AX-008** | **No Violencia Algorítmica** | Evaluador determinista **Fail-Close** basado en umbrales fisiológicos objetivos (escala NEWS2). |
| **AX-009** | **Transparencia Forense** | Ledger *append-only* en PostgreSQL con `pg_advisory_xact_lock`, `ENABLE ALWAYS TRIGGER` y encadenamiento **SHA-256**. |
| **AX-012** | **Desarme Autónomo** | `IrreversibleDecisionGuard`: exige firma digital **ECDSA P-256** válida y registro único de `Nonce`. |
| **AX-016** | **Refutación Humana** | Capacidad del profesional de sobreescribir protocolos, dejando justificación firmada criptográficamente. |

---

## 3. Arquitectura de Seguridad por Diseño (Grado Militar)

1. **Cero Dependencias Externas (`100% Go stdlib-only`):** Binario compilado con `CGO_ENABLED=0` y `GOTOOLCHAIN=local`, eliminando vectores de ataque a la cadena de suministro.
2. **Barrera HTTP/JSON Endurecida:** Uso de `StrictDecodeJSON` con límite físico de 64 KB y rechazo estricto de campos desconocidos.
3. **Ledger SQL Inmune:** Serialización de inserciones mediante `pg_advisory_xact_lock` para impedir *Ledger Forking*, y triggers `ENABLE ALWAYS` que ni un DBA comprometido puede desactivar.

---

## 4. Verificación Rápida para Auditores

```bash
# 1. Verificar suite criptográfica con detector de concurrencia
GOTOOLCHAIN=local go test -v -race ./pkg/governance ./pkg/adapter

# 2. Ejecutar demostración en vivo de los 5 axiomas blindados
GOTOOLCHAIN=local CGO_ENABLED=0 go run ./cmd/demo-governance