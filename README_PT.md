# NurseHealth.ai | Plataforma Pan-Americana de Enfermagem Domiciliar e Biossegurança

**Enfermagem domiciliar de precisão com IA no edge e rastreabilidade criptográfica.**

[![CI](https://github.com/Barto912/nursehealth-ai/actions/workflows/main.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Docs](https://github.com/Barto912/nursehealth-ai/actions/workflows/auditoria-docs.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)
[![Seguridad](https://github.com/Barto912/nursehealth-ai/actions/workflows/seguridad-gitleaks.yml/badge.svg)](https://github.com/Barto912/nursehealth-ai/actions)

🌐 **Leia esta documentação em outros idiomas:**
- [🇺🇸 English (Official GCP Proposal)](./README.md)
- [🇪🇸 Español (Operaciones Hispanoamérica)](./README_ES.md)

> 🚧 **Status do projeto: PESQUISA E DESIGN ARQUITETÔNICO.**
> Este repositório não é um produto clínico em produção. O código e os documentos aqui presentes são experimentais e seu uso em ambientes clínicos reais não está autorizado nem validado.

---

## 🛡️ O Problema e a Doutrina de Segurança Sanitária

Milhares de pacientes crônicos nas Américas recebem atendimento domiciliar sem rastreabilidade real. Hoje, esse trabalho é invisível: planilhas de papel, medicamentos vencidos em casa, quebra na cadeia de custódia e decisões clínicas sem evidências auditáveis.

Alinhada à **Doutrina de Segurança Sanitária das Américas**, a **NurseHealth.ai** transforma o cuidado descentralizado em uma rede de biossegurança auditável (Zero-Trust) construída nativamente no ecossistema do **Google Cloud Platform (GCP)**.

* **Piloto Inicial:** 57 condomínios fechados em Canning, Buenos Aires, Argentina.
* **Visão Continental:** Escalonamento regional para o Brasil (SUS/ANVISA), Chile, Uruguai, Colômbia, México e Estados Unidos sob um protocolo clínico e criptográfico unificado.

---

## 💡 A Solução

Plataforma extra-hospitalar de enfermagem de precisão e logística reversa farmacêutica:

- **IA Híbrida (Cloud/Edge):** Identificação de medicamentos e triagem clínica no dispositivo do enfermeiro (**Google Pixel 11 Pro**), sem expor dados sensíveis na nuvem.
- **Cadeia de Custódia Criptográfica:** Cada evento clínico é selado com **SHA-256** em um armazenamento de eventos *append-only*.
- **Logística Reversa Farmacêutica:** Coleta domiciliar de medicamentos vencidos com rastreabilidade de ponta a ponta até o descarte final bio-seguro.
- **Fail-Close by Design:** Em caso de dúvida, o motor determinístico em Go bloqueia a execução e escala para uma autoridade médica humana. Nunca improvisa.

> **Lema:** *O critério é humano, a rastreabilidade é do sistema.*

**Tese da Plataforma:** Enfermagem de precisão com suporte de IA que eleva o trabalho autônomo qualificado a uma carreira profissional verificável: educação contínua, reputação auditável e dignidade econômica para os cuidadores das Américas.

---

## 🏗️ Arquitetura Técnica e Fluxo de IA Híbrida

1. **Edge Node — Google Pixel 11 Pro:** MedGemma quantizado (2B/4B) no Tensor NPU + TFLite + Engram (armazenamento de contexto local SQLite/FTS5).
2. **Transporte Zero-Trust:** Fluxo de eventos clínicos pseudonimizados e assinados criptograficamente.
3. **Ingestão Core — Google Cloud Run:** Microsserviços 100% Go (Golang), REST API, Motor Fail-Close determinístico e Event Sourcing.
4. **Orquestração Cloud AI — Vertex AI:** MedGemma Maestro (7B/27B) + Google Harness² (guardrails dinâmicos *Propose-Probe-Compose*).
5. **Auditoria Imutável — Cloud SQL:** PostgreSQL Append-Only SHA-256 Event Store.
6. **Interoperabilidade Clínica — Cloud Healthcare API:** HL7 FHIR R4 e SNOMED CT isolados atrás de uma camada anticorrupção (ACL) Domain-Driven Design.

### Stack Técnico

| Camada | Tecnologia |
| :--- | :--- |
| **Backend Core** | 100% Go (Golang) no Google Cloud Run (Event Sourcing & Fail-Close) |
| **Frontend Mobile** | React Native com Expo |
| **Edge AI (Dispositivo)** | MedGemma Small (2B/4B) quantizado & TensorFlow Lite no Pixel 11 Pro |
| **Cloud AI (GCP)** | MedGemma Maestro (7B/27B) no Vertex AI + Google Harness² |
| **Governança IA & Memória** | Gentle-AI (ODD/RDD) + Engram (SQLite + FTS5 local-first) |
| **Banco de Dados & Ledger** | PostgreSQL no Cloud SQL (Event Store SHA-256 Append-only) |
| **Interoperabilidade** | HL7 FHIR R4 & SNOMED CT via Google Cloud Healthcare API |
| **Infraestrutura** | Google Cloud Platform (GCP) |

---

## ⚖️ Alinhamento Regulatório Multijurisdicional

Projetado desde o primeiro dia para satisfazer os padrões hemisféricos de segurança da saúde e soberania de dados:
* 🇧🇷 **Brasil:** Adesão estrita à LGPD, ANVISA (diretrizes SaMD - Software como Dispositivo Médico) e prontidão para interoperabilidade nativa com o **SUS (Rede Nacional de Dados em Saúde - RNDS)**.
* 🇺🇸 **Estados Unidos:** HIPAA, SOC 2 Type II, alinhamento de segurança CISA e conformidade NDAA.
* 🌎 **América Hispânica:** Marcos locais de privacidade de dados (Argentina, Chile, Colômbia, México) e padrões nacionais de Prontuário Eletrônico do Paciente.

---

## 📈 Status Atual e Marcos

Design arquitetônico concluído e auditado:
- **Event Storming (MVP-0):** Concluído com 14 hotspots clínicos/operacionais resolvidos.
- **Limites de Agregados v0.3:** Padrões de reparo vs. aceitação formalmente documentados.
- **Context Map (DDD):** Camada Anticorrupção (ACL) para FHIR R4 e Modelo de Referência regulatório estabelecidos.
- **Integração Contínua:** Pipelines CI/CD verdes, auditoria automatizada de documentação e verificação de segredos `gitleaks`.

**Em andamento:** Implementação do MVP e orquestração de Cloud AI para o **Google for Startups Cloud Program**.

---

## 🧭 Governança Ética

Princípios inegociáveis do manifesto do projeto:

1. **Dignidade sobre os dados.**
2. **Autonomia sobre a automação.**
3. **Transparência sobre caixas pretas.**
4. **Cuidado sobre a eficiência.**
5. **Presença sobre a conexão.**
6. **Prevenção sobre a reação.**
7. **Educação sobre a doutrinação.**

Os dados dos pacientes são processados de forma estritamente pseudonimizada. Nenhuma decisão clínica é automatizada.

---

## 🧠 Filosofia do Cuidado: O Firewall Humano

A enfermagem de precisão exige mais do que habilidade técnica: requer integridade cognitiva. O **Firewall Humano** (ADR-006) codifica três axiomas fundacionais que condicionam cada evento clínico:

- **AX-001 (Inércia e Ruído):** a clareza mental como precondição. Se `MentalClarity < 7`, o sistema bloqueia a criação do evento e registra a tentativa.
- **AX-002 (Tons / nível de trabalho):** o tom energético modula a eficácia clínica. Se `WorkLevel == LOW`, o sistema sugere realocação ou pausa.
- **AX-003 (Atenção simples):** a atenção dirigida como condição para assinar o evento SNOMED. Se `SimpleAttention == false`, o operador deve registrar um `AttentionalFinding` antes de assinar.

Estes axiomas não são aspiracionais: são **invariantes de domínio** exigidas pelo Value Object `RuidoCheck` com testes unitários. Um profissional em estado cognitivo não apto não pode criar nem assinar um evento clínico: o sistema bloqueia e registra a tentativa.

**Implementação:**

- [Axiomas fundacionais](docs/philosophy/axioms.json) (AX-001 a AX-003)
- [Value Object RuidoCheck](pkg/cognitive/ruido_check.go) (gates tipados com testes)
- [ADR-006: Firewall Humano (MVP)](docs/architecture/decisions/ADR-006-humanistic-firewall-mvp.md)
- [Context Map v0.2](docs/domain/context-map.md) (invariantes I10, I11)

---

## 👥 Equipe Core

- **Barto** — Founder & Product Owner.
- **Leonor** — Strategic Advisor.

---

## 🛡️ Arquitetura de Resiliência Multicamadas e Soberania de Dados

O NurseHealth.ai opera sob a **Doutrina de Segurança Sanitária das Américas**, garantindo continuidade clínica e conformidade regulatória em todos os momentos.

- **Camada 0 (Intra-Regional):** Nossa primeira linha de defesa. Operamos com alta disponibilidade síncrona multizona em `southamerica-east1` (São Paulo), garantindo RPO=0 e RTO<60s para a soberania de dados no Brasil, em estrita conformidade com a **LGPD (Lei Geral de Proteção de Dados)** e a **RDC 657/2022 da ANVISA** para SaMD.
- **Fail-Close Imutável (Invariantes I10/I11):** Nenhuma falha de infraestrutura permite novas tentativas se os gates cognitivos do operador (`RuidoCheck`, AX-001 a AX-003) ou do paciente (AX-004 a AX-006) tiverem bloqueado a operação. O veto clínico é determinístico e registrado de forma imutável (append-only) no ledger.
- **Motivação Operacional:** Projetado após o incidente público em `europe-west2` (outubro de 2026). Nossa topologia pan-americana opera de forma soberana, sem dependência de regiões externas.

*Para detalhes técnicos de replicação entre regiões e failover, consulte `docs/gcp-proposal/04_multiregion_dr_and_failover.md`.*

## 📄 Licença

Proprietário. Todos os direitos reservados.
