# Compliance Audit Registry (Gate 6 Preparation)

Purpose: register external expert dictámenes (legal and pharmaceutical-clinical)
without uploading PII-bearing documents to the public repository.
Zero-Trust posture per architectural mitigation Q5 (Argon-Deivis consensus).

## Method

- Each received dictamen is hashed locally with SHA-256 on the operator machine.
- This registry stores only: ID, title, issuer role, reception date, SHA-256 hash and status.
- Originals remain in offline private custody; any auditor can verify integrity
  by recomputing the hash against the original file on demand.

## Registry

| ID | Document | Issuer (role) | Reception | SHA-256 | Status |
| --- | --- | --- | --- | --- | --- |
| COMP-001 | Legal dictamen: regulatory framework and sensitive data (Argentina) | External legal counsel | Pending (deadline 2026-11-15) | To be registered upon reception | Pending |
| COMP-002 | Pharmaceutical dictamen: traceability and cold chain | Technical pharmaceutical director | Pending (deadline 2026-11-15) | To be registered upon reception | Pending |

## Local hashing procedure (never upload the source PDF)

```bash
sha256sum dictamen_legal.pdf
sha256sum dictamen_farmaceutico.pdf
```

On macOS use `shasum -a 256 <file>`. On Windows PowerShell use `Get-FileHash -Algorithm SHA256 <file>`.

## Rules

- Never commit PDFs, scans or documents containing PII to this repository
  (enforced by gitleaks plus manual review at commit time).
- Update this registry in the same commit that references a dictamen elsewhere
  in docs, preserving traceability.
- Full 64-hex hashes are recorded here upon reception; originals stay offline.