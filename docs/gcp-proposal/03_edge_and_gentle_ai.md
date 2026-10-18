# Edge AI, Gentle-AI Governance & Deterministic Fail-Close Engine

## 1. Edge Execution on Google Pixel 11 Pro
NurseHealth.ai leverages the advanced neural processing capabilities of the Google Pixel 11 Pro. By running quantized MedGemma (2B/4B) models directly on the Tensor NPU via TensorFlow Lite, the platform ensures that critical triage—such as wound assessment and vital sign validation—can occur in real-time, even in disconnected or low-bandwidth environments.

## 2. Gentle-AI & Engram Local Store
To guarantee patient data sovereignty and localized context retention:
* **Gentle-AI Methodology:** Employs RDD (Readiness-Driven Development) and ODD (Observation-Driven Development) to ensure the AI behaves transparently and predictably.
* **Engram Context Store:** An offline-first, local-only database utilizing SQLite and FTS5. It stores pseudonymized patient context directly on the edge device, preventing unnecessary data exposure to the cloud and ensuring the Edge AI has immediate access to historical case data.

## 3. Go Fail-Close Deterministic Engine
At the core of the NurseHealth.ai architecture is the Go (Golang) Fail-Close Engine. Designed for high-reliability medical environments:
* **Deterministic Execution:** The engine evaluates every AI proposition against hardcoded, deterministic clinical guardrails.
* **Fail-Close Paradigm:** If the AI output exhibits low confidence, logical contradictions, or violates any safety boundary, the system automatically defaults to a "closed" or "safe" state. It halts automated processing and instantly escalates the case to the human Chief Medical Officer (CMO) or designated clinical authority. No clinical action is ever taken autonomously under uncertainty.
