package main

import (
	"fmt"
	"os"

	"github.com/Barto912/nursehealth-ai/pkg/adapter"
	"github.com/Barto912/nursehealth-ai/pkg/governance"
)

type mockSink struct{}

func (m *mockSink) AppendForensicRecord(receipt governance.VerificationReceipt, status string, payload *string) error {
	fmt.Printf("   [Sink Forense] Asiento registrado: Status=%s | Hash=%s...\n", status, receipt.ReceiptSHA256[:16])
	return nil
}

type mockTaxonomy struct{}

func (m *mockTaxonomy) IsValidFHIRAndSNOMED(snomedCode string, fhirPayload []byte) bool {
	return true
}

func main() {
	fmt.Println("🛡️ NurseHealth.ai - Demo de Gobernanza Ética (Desarme de la IA)")
	fmt.Println("================================================================")

	sink := &mockSink{}
	tax := &mockTaxonomy{}
	
	dataProxy := adapter.NewDataProxy()
	sanitizer := adapter.NewPHISanitizerAdapter(dataProxy)

	guard := governance.NewIrreversibleDecisionGuard(sink, tax)
	logger := governance.NewDisagreementLogger(sink, sanitizer)

	fmt.Println("\n[Escenario 1] MedGemma 2B INT8 intenta autorizar un procedimiento IRREVERSIBLE...")
	
	supID := "DR-ROJO-001"
	sig := "TITAN-M3-SIG-999"
	
	_, err := guard.ValidateAndSeal(
		"DEMO-EVT-001",
		governance.StatusApprovedByCommittee,
		governance.ProcedureIrreversible,
		governance.AgentDevice,
		"MedGemma-2B-INT8",
		"HASH-DATASET-DEMO",
		"123456",
		[]byte(`{"resourceType":"CareEvent"}`),
		"NURSE-001",
		&supID,
		&sig,
	)

	if err != nil {
		fmt.Printf("   ✅ BLOQUEADO por Motor Fail-Close: %v\n", err)
	} else {
		fmt.Println("   ❌ ERROR CRÍTICO: El sistema permitió a la IA autorizar.")
		os.Exit(1)
	}

	fmt.Println("\n[Escenario 2] Enfermero en terreno refuta sugerencia de MedGemma...")

	tuple, err := logger.LogRefutation(
		"DEMO-EVT-002",
		"123456",
		"Sugerencia: Administrar antibiótico. Paciente: Juan Pérez, DNI 12345678.",
		"Refutación: Paciente sin fiebre, descartar sepsis. Mantener observación.",
		"NURSE-001",
	)

	if err != nil {
		fmt.Printf("   ❌ ERROR al registrar refutación: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("   ✅ Refutación registrada. Tupla ATDP generada.\n")
	fmt.Printf("   ✅ Hash SHA-256 (Tupla): %s\n", tuple.TupleSHA256)
	fmt.Println("\n🛡️ Demo completada. Soberanía Humana y Desarme de la IA verificados.")
}
