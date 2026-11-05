# NurseHealth.ai | Plataforma Pan-Americana de Enfermagem Domiciliar de Precisão e Bioseguridade Verde

**Enfermagem Domiciliar de Precisão, IA Híbrida Edge-to-Cloud, Rastreabilidade Criptográfica e Logística Reversa Farmacêutica Verde.**

🌐 **Ler esta documentação em outros idiomas:**
- [🇺🇸 English (Official Google for Startups & Pan-American Architecture)](./README.md)
- [🇪🇸 Español (Operaciones Hispanoamérica, ANMAT & Economía del Conocimiento)](./README_ES.md)

> 🚧 **Status do Projeto: PESQUISA, P&D E VALIDAÇÃO ARQUITETURAL (`v0.4.0-research-proposal`).**  
> Este repositório contém a arquitetura fundamental de P&D Experimental, o núcleo determinista em Go (`stdlib-only`) e os pipelines soberanos de destilação de IA da NurseHealth.ai. O implante clínico em ambiente real é estritamente controlado sob protocolos de validação institucional.

---

## 1. Resumo Executivo e Doutrina de Segurança Sanitária das Américas

Milhares de pacientes crônicos e pós-agudos em todo o continente recebem assistência domiciliar sem rastreabilidade verificável. Hoje, o atendimento extra-hospitalar permanece fragmentado: registros em papel, cadeias de frio comprometidas, medicamentos vencidos acumulados nas residências e decisões clínicas sem respaldo auditável.

Enquadrada na **Doutrina de Segurança Sanitária das Américas**, a **NurseHealth.ai** transforma a internação e enfermagem domiciliar descentralizada em uma rede auditável de bioseguridade Zero-Trust e logística reversa verde, construída nativamente sobre o ecossistema **Google Cloud Platform (GCP)** e **Google AI (Gemma / MedGemma)**.

- **Piloto Clínico e Logístico Inicial:** Coorte de 100 casos distribuída em 57 condomínios fechados (*gated communities*) no **corredor Canning (Esteban Echeverría / Ezeiza), Província de Buenos Aires, Argentina**, preparando a expansão para o mercado brasileiro.
- **Roteiro de Escala Continental:** Expansão faseada para o **Brasil, Colômbia e México e Estados Unidos** sob um protocolo criptográfico e clínico unificado (HL7 FHIR R4, SNOMED CT e SHA-256).
- **Tese da Plataforma:** *O julgamento é humano; a rastreabilidade pertence ao sistema.* A NurseHealth.ai eleva a enfermagem autônoma a uma carreira verificável, prevenindo a poluição farmacológica doméstica por meio de logística reversa eco-roteirizada.

---

## 2. Capacidades Nucleares e Solução de Bioseguridade Verde

1. **IA Híbrida em Três Camadas: Edge-to-Workstation-to-Cloud (Ecossistema Gemma):**
   * **Triage na Borda (Em Domicílio):** Verificação de medicamentos em tempo real, documentação de feridas e triagem clínica no dispositivo **Google Pixel 11 Pro** do profissional, sem expor dados de identificação (PII/PHI) brutos a redes externas.
   * **Orquestração em Nuvem:** Raciocínio clínico multimodal via **MedGemma Maestro (7B/27B) no GCP Vertex AI**, governado por guardrails clínicos dinâmicos.
2. **Logística Reversa Farmacêutica Verde e IoT de Cadeia de Frio:**
   * Coleta domiciliar de medicamentos vencidos, não utilizados ou termossensíveis (+2°C a +8°C) utilizando recipientes térmicos IoT (Bluetooth 6 / Thread).
   * Rotas de coleta otimizadas para menor pegada de carbono por meio do **Google Maps Platform (Routes API Eco-Friendly Routing & Fleet Engine)**, garantindo a cadeia de custódia de ponta a ponta até a destinação final biosegura certificada.
3. **Motor Determinista Fail-Close (`100% Go stdlib-only`):**
   * Zero dependências externas no núcleo de decisão clínica. Ante qualquer discrepância lógica, discordância taxonómica SNOMED CT ou alerta de fadiga do operador, o motor bloqueia mecanicamente a execução, emite um recibo de rejeição assinado e escala para uma autoridade médica humana.
4. **Livro-Razão de Auditoria Criptográfica Imutável:**
   * Cada intervenção clínica, validação de portão cognitivo e transferência de custódia farmacêutica é selada com **SHA-256** em um armazenamento de eventos *append-only* (somente anexado).

---

## 3. Arquitetura Técnica de Três Camadas e Fluxo Híbrido de IA

Para satisfazer simultaneamente a **escalabilidade em Google Cloud AI** e os **mandatos de Soberania de Dados**, a NurseHealth.ai opera em três níveis físicos:

1. **Nível 1 — Nódulo Móvel na Borda (Google Pixel 11 Pro / Android 17):**
   * Executa **MedGemma Small (2B INT8)** quantizado via LiteRT na **NPU Google Tensor G6** (`< 300ms` de latência).
   * O enclave de segurança por hardware **Titan M3** assina localmente os recibos `VerificationReceipt` usando SHA-256.
   * **Armazenamento de Contexto Local:** Base SQLite + FTS5 local que despista e pseudonimiza dados sensíveis (PHI) na borda antes de qualquer transmissão.
2. **Nível 2 — Orquestração em Nuvem e Ledger Multirregião (Google Cloud Platform):**
   * **Google Cloud Run:** Microsserviços 100% Go, API REST e ingestão por Event Sourcing.
   * **Vertex AI:** MedGemma Maestro (7B/27B) com guardrails dinâmicos.
   * **Cloud SQL Enterprise Plus:** Event Store PostgreSQL Append-Only com encadeamento SHA-256.
   * **Cloud Healthcare API:** Interoperabilidade HL7 FHIR R4 e SNOMED CT isolada atrás de uma Camada Anticorrupção (ACL) baseada em Domain-Driven Design (DDD).
   * **Google Maps Platform:** Routes API e Fleet Engine para despacho domiciliar de baixas emissões e coleta farmacêutica reversa.

| Camada Arquitetural | Tecnologia e Padrões |
| :--- | :--- |
| **Backend Core** | 100% Go (núcleo `stdlib-only`) em **Google Cloud Run** (Event Sourcing e Motor Fail-Close) |
| **Nódulo Móvel Edge** | **Google Pixel 11 Pro** (Tensor G6 NPU + Titan M3) e MedGemma Quantizado (2B INT8) |
| **Cloud AI (GCP)** | **MedGemma Maestro (7B/27B)** no **Vertex AI** com guardrails clínicos dinâmicos |
| **Logística Verde e Geo** | **Google Maps Platform** (Routes API Eco-Routing & Fleet Engine) + IoT Cadeia de Frio BLE 6/Thread |
| **Ledger de Auditoria e BD** | **Cloud SQL Enterprise Plus (PostgreSQL)** Armazenamento Append-Only de Cadeia de Custódia SHA-256 |
| **Interoperabilidade Clínica** | **HL7 FHIR R4 e SNOMED CT** via **Google Cloud Healthcare API** (Camada Anticorrupção DDD) |

---

## 4. Filosofia do Cuidado: O Firewall Humanista (ADR-006 e ADR-010)

Inspirado na ética médica humanista (*Magnifica Humanitas* — zero iatrogenia, dignidade humana acima da automação), o **Firewall Humano** exige integridade cognitiva antes de habilitar qualquer ação clínica.

Estes axiomas (clareza mental obrigatória, modulação de carga de trabalho e atenção direcionada para a assinatura SNOMED CT) são compilados como **invariantes de domínio estritas em Go** dentro do núcleo de governança (`pkg/governance/`), respaldados por testes unitários automatizados (`go test ./...`). O sistema impõe o **Desarme da IA (AX-012)**: a IA pode sugerir, mas apenas um profissional humano com registro profissional válido pode autorizar ações clínicas irreversíveis.

---

## 5. Resiliência Multirregião e Recuperação de Desastres (DR)

A NurseHealth.ai garante a continuidade clínica sob a **Doutrina de Segurança Sanitária das Américas**:
- **Camada 0 (Alta Disponibilidade Intra-Regional):** Operação síncrona multizona em `southamerica-west1` (Santiago) e `us-east1` (Virgínia do Norte), alcançando **RPO = 0** e **RTO < 60s**.
- **Fail-Close Imutável:** Nenhuma comutação por falha de infraestrutura pode omitir os portões cognitivos do operador ou as invariantes de segurança do paciente.
- **Continuidade Soberana:** Arquitetura desenhada para impedir que quedas regionais de nuvem interrompam a assistência domiciliar crítica.

---

## 6. Alinhamento Regulatório Multijurisdiccional e Marco de P&D

Desenhado desde o dia zero para satisfazer marcos hemisféricos de bioseguridade, regimes de incentivo fiscal de P&D e normativas sanitárias:

- 🇦🇷 **Argentina (Jurisdição Piloto e Hub de P&D):** Estruturado sob o regime de Economia do Conhecimiento (Lei 27.506 - Códigos ARCA CLAE 620100 e 721000) e conformidade com ANMAT (SaMD), Lei 25.326 e Lei 24.051 (Resíduos Perigosos / Logística Reversa).
- 🇺🇸 **Estados Unidos (Hub Corporativo e de Escala):** Conformidade com HIPAA, FDA SaMD (Good Machine Learning Practice), SOC 2 Type II e NDAA Sec. 889 via estrutura matriz Delaware C-Corp.
- 🇧🇷 **Brasil (Expansão e Conformidade Regulatória):**
  * **Privacidad e Saúde Digital:** Lei Geral de Proteção de Dados (LGPD - Lei nº 13.709/2018).
  * **Dispositivos Médicos e Cadeia de Frio:** Diretrizes da **ANVISA** para Software como Dispositivo Médico (SaMD conforme **RDC Nº 657/2022**) e boas práticas de distribuição de medicamentos (**RDC Nº 430/2020**).
  * **Interoperabilidade Nacional:** Prontidão nativa para integração com a Rede Nacional de Dados em Saúde (**RNDS**) do Sistema Único de Saúde (**SUS**).

---

## 7. Estrutura do Repositório, ADRs e Proposta de Pesquisa GCP

- **Architecture Decision Records (ADRs):** Consolidados do ADR-001 ao ADR-010 (incluindo Firewall Humano, Extensões FHIR R4, Ledger de Auditoria SHA-256 e Arquitetura Disarm AI).
- **Integração Contínua (CI/CD):** Pipelines 100% em verde executando testes unitários/RDD em Go, validação de links Markdown, auditoria de ADRs e escaneamento automatizado de segredos com `gitleaks`.
- 📂 **Proposta Detalhada de Pesquisa GCP (`/docs/gcp-proposal/`):**
  1. [`01_architecture_overview.md`](./docs/gcp-proposal/01_architecture_overview.md) — Topologia Híbrida de Ponta a Ponta (Edge / Workstation / Cloud).
  2. [`02_panamerican_security_and_compliance.md`](./docs/gcp-proposal/02_panamerican_security_and_compliance.md) — Governança Zero-Trust, LGPD, ANVISA e SUS/RNDS.
  3. [`03_edge_and_gentle_ai.md`](./docs/gcp-proposal/03_edge_and_gentle_ai.md) — Quantização MedGemma, NPU Pixel 11 Pro e Fluxo Gentle-AI RDD/ODD.
  4. [`04_multiregion_dr_and_failover.md`](./docs/gcp-proposal/04_multiregion_dr_and_failover.md) — Conmutação Multirregião (`southamerica-west1` / `us-east1`).

---

## 8. Equipe Fundadora e Executiva

Liderança multidisciplinar que integra engenharia de software, medicina clínica, gestão de projetos, direito da saúde e logística farmacêutica:

- **Walter Bartolomé (`@Barto912`)** — Fundador, Product Owner e Arquiteto Principal de Sistemas e Ética.
- **Leonor Figueroa** — Technical Project Manager e Assessora Estratégica (Auditoria Documental, Conformidade de Repositório e Operações do Piloto Canning).
- **Dr. Daniel Rojo** — Chief Medical Information Officer (CMIO) e Diretor Médico (Taxonomias Clínicas SNOMED CT, Extensões FHIR R4 e Protocolos de Enfermagem Domiciliar).
- **Divisão de Assuntos Jurídicos e Regulatórios** — Direito Sanitário, LGPD / HIPAA e Estruturação Societária (Conselho Assessor e Direção Fracional).
- **Divisão de Desenvolvimento de Negócios e Operações** — Unit Economics e Parcerias Público-Privadas / B2B em Saúde (Conselho Assessor e Direção Fracional).
- **Divisão de Logística Farmacêutica e Bioseguridade** — Conformidade da Cadeia de Frio (+2°C a +8°C), Logística Reversa Verde e Rastreabilidade ANVISA (Farmacêuticos Especialistas).

*Contato Institucional:* `contact@nursehealth.ai`

---

## 9. Licença e Propriedade Intelectual

**Modelo de Licenciamento Híbrido:**
- **Proprietário (Todos os Direitos Reservados):** Os motores centrais em Go, protocolos criptográficos e pipelines de destilação de IA constituem propriedade intelectual exclusiva da NurseHealth.ai.
- **Aberto / Revisão Acadêmica (Apache 2.0 / CC-BY-NC-SA 4.0):** Os adaptadores FHIR/SNOMED, a documentação ADR e os esquemas de conformidade regulatória são publicados para promover a interoperabilidade, a revisão acadêmica e a auditoria institucional.
