package main

import "fmt"

// Función para probar diferentes casos de validación
func probarValidaciones() {
	fmt.Println("🧪 PROBANDO VALIDACIONES DE ERRORES")
	fmt.Println("=====================================")
	
	// Test 1: Cédula muy corta
	fmt.Println("\n1. Cédula muy corta:")
	datosError1 := FacturaInput{
		ClienteNombre:       "Juan Perez",
		ClienteCedula:       "123456789", // Solo 9 dígitos
		ProductoCodigo:      "PROD001",
		ProductoDescripcion: "Producto de prueba",
		Cantidad:            1.0,
		PrecioUnitario:      100.0,
	}
	_, err := CrearFactura(datosError1)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 2: Cédula con letras
	fmt.Println("\n2. Cédula con letras:")
	datosError2 := FacturaInput{
		ClienteNombre:       "Juan Perez",
		ClienteCedula:       "17131ABC71", // Contiene letras
		ProductoCodigo:      "PROD001",
		ProductoDescripcion: "Producto de prueba",
		Cantidad:            1.0,
		PrecioUnitario:      100.0,
	}
	_, err = CrearFactura(datosError2)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 3: Cantidad cero
	fmt.Println("\n3. Cantidad inválida:")
	datosError3 := FacturaInput{
		ClienteNombre:       "Juan Perez",
		ClienteCedula:       "1713175071", // Cédula válida
		ProductoCodigo:      "PROD001",
		ProductoDescripcion: "Producto de prueba",
		Cantidad:            0.0, // Cantidad inválida
		PrecioUnitario:      100.0,
	}
	_, err = CrearFactura(datosError3)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 4: Nombre vacío
	fmt.Println("\n4. Nombre vacío:")
	datosError4 := FacturaInput{
		ClienteNombre:       "", // Nombre vacío
		ClienteCedula:       "1713175071",
		ProductoCodigo:      "PROD001",
		ProductoDescripcion: "Producto de prueba",
		Cantidad:            1.0,
		PrecioUnitario:      100.0,
	}
	_, err = CrearFactura(datosError4)
	if err != nil {
		fmt.Printf("   ❌ Error esperado: %v\n", err)
	}
	
	// Test 5: Datos válidos
	fmt.Println("\n5. Datos válidos:")
	datosValidos := FacturaInput{
		ClienteNombre:       "Maria Rodriguez",
		ClienteCedula:       "1713175071", // Cédula válida
		ProductoCodigo:      "PROD001",
		ProductoDescripcion: "Producto de prueba",
		Cantidad:            3.0,
		PrecioUnitario:      50.0,
	}
	factura, err := CrearFactura(datosValidos)
	if err != nil {
		fmt.Printf("   ❌ Error inesperado: %v\n", err)
	} else {
		fmt.Printf("   ✅ Factura creada exitosamente para %s\n", factura.InfoFactura.RazonSocialComprador)
		fmt.Printf("   💰 Total: $%.2f\n", factura.InfoFactura.ImporteTotal)
	}
}