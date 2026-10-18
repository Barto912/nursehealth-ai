# ADR-006: Firewall Humano como precondición de trazabilidad (MVP)

## Estado
Aceptado

## Contexto
La iatrogenia en enfermería domiciliaria no surge solo de la falta de
conocimiento técnico, sino de la interferencia de ruido (tensiones, climas
emotivos, inercia de eventos anteriores) en el nivel de trabajo vigílico del
profesional (AX-001). La tecnología no puede corregir esto por sí sola: debe
exigir una validación humana previa.

## Decisión
Implementar un Value Object `RuidoCheck` en el bounded context Cognitivo
(`pkg/cognitive`) que actúa como invariante de dominio. Ningún `CareEvent`
puede crearse ni firmarse sin un `RuidoCheck` válido adjunto como
`CognitiveSnapshot`.

Gates del MVP:

- `MentalClarity >= 7` (AX-001). Umbral conservador; hipótesis a validar
  con cognitive_audit_log y 3 enfermeros reales.
- `WorkLevel != LOW` (AX-002). En MVP el gate aplica a todo CareEvent.
  Criterio de reversión: cuando exista catálogo de criticidad de
  procedimientos, el gate se relajará a solo procedimientos críticos.
- `SimpleAttention == true` o `AttentionalFinding` no vacío (AX-003,
  Autoliberación Serie 13) para firmar el evento SNOMED.

Todo intento bloqueado se registra en `cognitive_audit_log` (tabla
inmutable). Cero dependencia de hardware externo.

## Consecuencias
- El dominio clínico (CareEvent) depende de una validación previa del
  dominio cognitivo.
- Trazabilidad inmutable de fatiga o bloqueo del operador.
- Sin escalación automática en MVP: el supervisor humano consulta
  cognitive_audit_log manualmente. Las alertas por patrones son future work.

## Corrección de fuente (H20)
La versión inicial de AX-003 citaba un pasaje de Silo_OC2_es.pdf que no fue
hallado en el texto durante la verificación. El axioma fue re-sourced a
Autoliberacion.pdf (Serie 13, Atención simple). El lenguaje ubicuo del
dominio sigue a la fuente verificada: SimpleAttention y AttentionalFinding.

## Fuentes no invocadas
No se citan estudios externos no verificados en la base de conocimientos
del proyecto. Los umbrales son hipótesis MVP, no verdades citadas.