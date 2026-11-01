package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
)

type AxiomsDocument struct {
	Meta      map[string]interface{}   `json:"meta"`
	Axioms    []map[string]interface{} `json:"axioms"`
	Integrity map[string]interface{}   `json:"integrity"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: validate-axioms <path-to-json>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("❌ Error leyendo archivo: %v\n", err)
		os.Exit(1)
	}

	var doc AxiomsDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		fmt.Printf("❌ JSON inválido: %v\n", err)
		os.Exit(1)
	}

	// Validar cantidad de axiomas
	if len(doc.Axioms) != 7 {
		fmt.Printf("❌ Se esperaban 7 axiomas, encontrados %d\n", len(doc.Axioms))
		os.Exit(1)
	}

	// Validar IDs esperados
	expectedIDs := []string{"AX-007", "AX-008", "AX-009", "AX-010", "AX-011", "AX-012", "AX-016"}
	var actualIDs []string
	for _, ax := range doc.Axioms {
		if id, ok := ax["id"].(string); ok {
			actualIDs = append(actualIDs, id)
		}
	}
	sort.Strings(actualIDs)
	sort.Strings(expectedIDs)
	if !reflect.DeepEqual(actualIDs, expectedIDs) {
		fmt.Printf("❌ IDs no coinciden. Esperados: %v, Encontrados: %v\n", expectedIDs, actualIDs)
		os.Exit(1)
	}

	// Validar que NO existan IDs excluidos (fuentes no ratificadas por el Founder)
	excludedIDs := []string{"AX-013", "AX-014", "AX-015"}
	for _, ax := range doc.Axioms {
		if id, ok := ax["id"].(string); ok {
			for _, ex := range excludedIDs {
				if id == ex {
					fmt.Printf("❌ ID excluido encontrado: %s (fuente no ratificada)\n", id)
					os.Exit(1)
				}
			}
		}
	}

	// Validar presencia de firmas humanas
	meta := doc.Meta
	if sigs, ok := meta["human_signatures"].(map[string]interface{}); ok {
		requiredRoles := []string{"founder_and_ethical_architect", "cmio", "technical_project_manager"}
		for _, role := range requiredRoles {
			if _, exists := sigs[role]; !exists {
				fmt.Printf("❌ Falta firma humana requerida: %s\n", role)
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("❌ Falta bloque human_signatures en meta")
		os.Exit(1)
	}

	// Validar homologación CMIO
	if cmio, ok := meta["clinical_homologation"].(map[string]interface{}); ok {
		if status, ok := cmio["status"].(string); !ok || status != "HOMOLOGATED" {
			fmt.Println("❌ Homologación CMIO no está en estado HOMOLOGATED")
			os.Exit(1)
		}
	} else {
		fmt.Println("❌ Falta bloque clinical_homologation en meta")
		os.Exit(1)
	}

	// Validar registro ARGON (única AI del ecosistema permitida)
	if ai, ok := meta["ai_ecosystem_registry"].(map[string]interface{}); ok {
		if model, ok := ai["current_model"].(string); !ok || model != "ARGON" {
			fmt.Println("❌ Modelo AI debe ser ARGON (no otros modelos)")
			os.Exit(1)
		}
	} else {
		fmt.Println("❌ Falta bloque ai_ecosystem_registry en meta")
		os.Exit(1)
	}

	fmt.Println("✅ Validación completa: 7 axiomas válidos, firmas humanas presentes, homologación CMIO OK, registro ARGON OK")
	fmt.Println("   IDs válidos: AX-007, AX-008, AX-009, AX-010, AX-011, AX-012, AX-016")
	fmt.Println("   IDs excluidos correctamente: AX-013, AX-014, AX-015 (Evangelium Vitae)")
	os.Exit(0)
}
