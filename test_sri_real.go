package main

import (
	"fmt"
	"go-facturacion-sri/sri"
)

func main() {
	fmt.Println("=======================================================")
	fmt.Println("🚀 TESTING DE INTEGRACIÓN SRI REAL")
	fmt.Println("=======================================================")
	
	err := sri.TestearIntegracionSRIReal()
	if err != nil {
		fmt.Printf("❌ Error en integración SRI: %v\n", err)
		return
	}
	
	fmt.Println("\n✅ INTEGRACIÓN SRI COMPLETADA EXITOSAMENTE")
}