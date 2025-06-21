package main

import (
	"fmt"
	
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
)

// Función para probar diferentes casos de validación
func probarValidaciones() {
	fmt.Println("🧪 PROBANDO VALIDACIONES DE ERRORES")
	fmt.Println("=====================================")
	
	// Test 1: Cédula muy corta
	fmt.Println("\n1. Cédula muy corta:")
	datosError1 := models.FacturaInput{
		ClienteNombre: "Juan Perez",
		ClienteCedula: "123456789", // Solo 9 dígitos
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto de prueba",
				Cantidad:       1.0,
				PrecioUnitario: 100.0,
			},
		},
	}
	_, err := factory.CrearFactura(datosError1)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 2: Cédula con letras
	fmt.Println("\n2. Cédula con letras:")
	datosError2 := models.FacturaInput{
		ClienteNombre: "Juan Perez",
		ClienteCedula: "17131ABC71", // Contiene letras
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto de prueba",
				Cantidad:       1.0,
				PrecioUnitario: 100.0,
			},
		},
	}
	_, err = factory.CrearFactura(datosError2)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 3: Cantidad cero
	fmt.Println("\n3. Cantidad inválida:")
	datosError3 := models.FacturaInput{
		ClienteNombre: "Juan Perez",
		ClienteCedula: "1713175071", // Cédula válida
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto de prueba",
				Cantidad:       0.0, // Cantidad inválida
				PrecioUnitario: 100.0,
			},
		},
	}
	_, err = factory.CrearFactura(datosError3)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 4: Nombre vacío
	fmt.Println("\n4. Nombre vacío:")
	datosError4 := models.FacturaInput{
		ClienteNombre: "", // Nombre vacío
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto de prueba",
				Cantidad:       1.0,
				PrecioUnitario: 100.0,
			},
		},
	}
	_, err = factory.CrearFactura(datosError4)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 5: Sin productos
	fmt.Println("\n5. Sin productos:")
	datosError5 := models.FacturaInput{
		ClienteNombre: "Juan Perez",
		ClienteCedula: "1713175071",
		Productos:     []models.ProductoInput{}, // Lista vacía
	}
	_, err = factory.CrearFactura(datosError5)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 6: Datos válidos con múltiples productos
	fmt.Println("\n6. Datos válidos con múltiples productos:")
	datosValidos := models.FacturaInput{
		ClienteNombre: "Maria Rodriguez",
		ClienteCedula: "1713175071", // Cédula válida
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto A",
				Cantidad:       2.0,
				PrecioUnitario: 30.0,
			},
			{
				Codigo:         "PROD002", 
				Descripcion:    "Producto B",
				Cantidad:       1.0,
				PrecioUnitario: 15.0,
			},
		},
	}
	factura, err := factory.CrearFactura(datosValidos)
	if err != nil {
		fmt.Printf("   ❌ Error inesperado: %v\n", err)
	} else {
		fmt.Printf("   ✅ Factura creada exitosamente para %s\n", factura.InfoFactura.RazonSocialComprador)
		fmt.Printf("   💰 Total: $%.2f\n", factura.InfoFactura.ImporteTotal)
	}
}