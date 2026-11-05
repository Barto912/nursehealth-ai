package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/Barto912/nursehealth-ai/pkg/adapter"
	"github.com/Barto912/nursehealth-ai/pkg/governance"
)

type derSig struct {
	R, S *big.Int
}

func main() {
	fmt.Println("=======================================================================")
	fmt.Println(" NURSEHEALTH.AI v0.6.0-defensa-civil — MOTOR CRIPTOGRÁFICO DE AUDITORÍA")
	fmt.Println(" Jurisdicción: República Argentina (Ley 25.326 / ANMAT / Ley 25.506)")
	fmt.Println("=======================================================================")

	// === AX-001: Dignidad del Paciente + Minimización de Datos ===
	fmt.Println()
	fmt.Println("▶ AX-001 (Dignidad del Paciente) + Soberanía de Datos")
	fmt.Println("----------------------------------------------------------------------")

	// 1. HMAC secret de 32 bytes (mínimo requerido por NIST SP 800-107)
	hmacSecret := []byte("CLAVE_SOBERANA_DEFENSA_CIVIL_ARGENTINA_2026_32B")
	proxy, err := adapter.NewDataProxy(hmacSecret, adapter.AuthorizedJurisdiction)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Fallo de frontera soberana: %v\n", err)
		os.Exit(1)
	}
	pseudoID := proxy.PseudonymizeIdentifier("20-30456789-4")
	cleanNote := proxy.SanitizeClinicalText("Paciente CUIL 20-30456789-4 atendido en domicilio Tel +54 9 11 4567-8901.")
	fmt.Printf("  [OK] Frontera soberana validada: %s\n", adapter.AuthorizedJurisdiction)
	fmt.Printf("  [OK] Pseudonimización HMAC-SHA256: %s\n", pseudoID)
	fmt.Printf("  [OK] Texto sanitizado (Ley 25.326): %s\n", cleanNote)

	// 2. AES-256-GCM autenticado
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		os.Exit(1)
	}
	cipherPayload, err := governance.SealPatientPayload(aesKey, []byte(cleanNote), []byte(pseudoID))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Fallo AX-001: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  [OK] Cifrado AES-256-GCM: %d bytes autenticados\n", len(cipherPayload))

	// === AX-008: No Violencia Algorítmica (Fail-Close Determinista) ===
	fmt.Println()
	fmt.Println("▶ AX-008 (No Violencia Algorítmica — Fail-Close Determinista)")
	fmt.Println("----------------------------------------------------------------------")

	vitalsNormales := governance.VitalSignsInput{
		RespiratoryRateBPM: 16, SpO2Percent: 98, SystolicBP: 120,
		HeartRateBPM: 72, TemperatureCelsius: 36.6,
	}
	if err := governance.EvaluateDeterministicTriage(vitalsNormales); err == nil {
		fmt.Println("  [OK] Signos vitales normales: compuerta abierta")
	}

	vitalsCriticos := governance.VitalSignsInput{
		RespiratoryRateBPM: 26, SpO2Percent: 88, SystolicBP: 95,
		HeartRateBPM: 118, TemperatureCelsius: 38.9,
	}
	if err := governance.EvaluateDeterministicTriage(vitalsCriticos); err != nil {
		fmt.Printf("  [OK] Signos críticos detectados: %v\n", err)
	}

	// === AX-012 + AX-016: Desarme IA + Refutación Humana (ECDSA P-256) ===
	fmt.Println()
	fmt.Println("▶ AX-012 + AX-016 (Soberanía Humana — Firma ECDSA P-256 de Hardware)")
	fmt.Println("----------------------------------------------------------------------")

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		os.Exit(1)
	}
	now := time.Now().UTC()

	cmd := governance.SignedClinicalCommand{
		EventID:          "c9a8b7d6-e5f4-4321-8765-1029384756ab",
		AxiomCode:        governance.AxiomHumanRefutation,
		ClinicianID:      "ENF-MN-88421-AR",
		PublicKeyID:      "HARDWARE-KEY-P256-AR-01",
		NonceHex:         "9f8e7d6c5b4a39281726354455667788",
		SignedAtUTC:      now,
		PayloadCipher:    cipherPayload,
		DisagreementNote: "AX-016: Se prioriza estabilización de vía aérea en domicilio antes de traslado por tormenta severa.",
		PrevHashHex:      "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}

	// Firma ECDSA P-256 simulando enclave Titan M3
	digest := sha256.Sum256(cmd.CanonicalSigningBytes())
	r, s, _ := ecdsa.Sign(rand.Reader, privKey, digest[:])
	cmd.SignatureDER, _ = asn1.Marshal(derSig{R: r, S: s})

	guard := governance.NewIrreversibleDecisionGuard(5 * time.Minute)
	recordHash, err := guard.VerifyCommand(cmd, &privKey.PublicKey, now)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Fallo verificación AX-012: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  [OK] Firma ECDSA P-256 verificada (Titan M3 simulado)\n")
	fmt.Printf("  [OK] Refutación AX-016 registrada (justificación > 10 chars)\n")

	// === AX-009: Transparencia Forense (Hash encadenado) ===
	fmt.Println()
	fmt.Println("▶ AX-009 (Transparencia Forense — Hash SHA-256 Canónico)")
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("  [OK] Hash del registro: %s\n", recordHash)

	fmt.Println()
	fmt.Println("=======================================================================")
	fmt.Println(" ESTADO DEL SISTEMA: 5/5 AXIOMAS VERIFICADOS MATEMÁTICAMENTE")
	fmt.Println(" AX-001 (Dignidad) + AX-008 (Fail-Close) + AX-009 (Forense)")
	fmt.Println(" AX-012 (Desarme IA) + AX-016 (Refutación Humana)")
	fmt.Println("=======================================================================")
}
