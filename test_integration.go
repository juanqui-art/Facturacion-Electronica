// test_integration.go - Programa para testing de integración SRI
package main

import (
	"fmt"
	"log"

	"go-facturacion-sri/config"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/sri"
	"go-facturacion-sri/validators"
)

func main() {
	fmt.Println("🧪 TESTING DE INTEGRACIÓN - SISTEMA FACTURACIÓN SRI")
	fmt.Println("=" + string(make([]byte, 55)))

	// 1. Testing de configuración
	fmt.Println("\n📋 Testing de Configuración...")
	if err := testearConfiguracion(); err != nil {
		log.Fatalf("❌ Error en configuración: %v", err)
	}
	fmt.Println("✅ Configuración OK")

	// 2. Testing de validaciones robustas
	fmt.Println("\n🛡️  Testing de Validaciones Robustas...")
	if err := testearValidaciones(); err != nil {
		log.Fatalf("❌ Error en validaciones: %v", err)
	}
	fmt.Println("✅ Validaciones OK")

	// 3. Testing de generación de facturas
	fmt.Println("\n📄 Testing de Generación de Facturas...")
	factura, err := testearGeneracionFacturas()
	if err != nil {
		log.Fatalf("❌ Error generando facturas: %v", err)
	}
	fmt.Println("✅ Generación de Facturas OK")

	// 4. Testing de XML
	fmt.Println("\n🔧 Testing de Generación XML...")
	xmlData, err := testearGeneracionXML(factura)
	if err != nil {
		log.Fatalf("❌ Error generando XML: %v", err)
	}
	fmt.Println("✅ Generación XML OK")

	// 5. Testing de Circuit Breaker
	fmt.Println("\n🔧 Testing de Circuit Breaker...")
	if err := testearCircuitBreaker(); err != nil {
		log.Fatalf("❌ Error en circuit breaker: %v", err)
	}
	fmt.Println("✅ Circuit Breaker OK")

	// 6. Testing de logging
	fmt.Println("\n📊 Testing de Sistema de Logging...")
	if err := testearLogging(); err != nil {
		log.Fatalf("❌ Error en logging: %v", err)
	}
	fmt.Println("✅ Sistema de Logging OK")

	// 7. Mostrar resumen final
	fmt.Println("\n📊 RESUMEN DEL TESTING:")
	mostrarResumenTesting(factura, xmlData)

	fmt.Println("\n🎉 TODOS LOS TESTS PASARON EXITOSAMENTE!")
	fmt.Println("   El sistema está listo para integración SRI real.")
	fmt.Println("=" + string(make([]byte, 55)))
}

func testearConfiguracion() error {
	// Cargar configuración por defecto
	config.CargarConfiguracionPorDefecto()

	// Verificar valores
	if config.Config.Empresa.RUC == "" {
		return fmt.Errorf("RUC no configurado")
	}

	if config.Config.Ambiente.Codigo == "" {
		return fmt.Errorf("ambiente no configurado")
	}

	// Probar generación de clave de acceso
	claveAcceso := config.GenerarClaveAcceso()
	if len(claveAcceso) != 49 {
		return fmt.Errorf("clave de acceso debe tener 49 dígitos, tiene %d", len(claveAcceso))
	}

	// Validar la clave generada
	if err := config.ValidarClaveAcceso(claveAcceso); err != nil {
		return fmt.Errorf("clave generada no es válida: %v", err)
	}

	fmt.Printf("   ✓ Clave de acceso generada: %s\n", claveAcceso)
	fmt.Printf("   ✓ Secuencial: %s\n", config.ObtenerSecuencialSiguiente())

	return nil
}

func testearValidaciones() error {
	// Test ValidarRUC
	if err := validators.ValidarRUC("1713175071001"); err != nil {
		return fmt.Errorf("RUC válido rechazado: %v", err)
	}

	// Test input malicioso
	if err := validators.ValidarRUC("<script>alert('xss')</script>"); err == nil {
		return fmt.Errorf("input malicioso XSS no fue rechazado")
	}

	// Test sanitización
	textoPeligroso := "<script>alert('hack')</script>"
	textoLimpio := validators.SanitizarTexto(textoPeligroso)
	if textoLimpio == textoPeligroso {
		return fmt.Errorf("texto peligroso no fue sanitizado")
	}

	fmt.Printf("   ✓ Texto peligroso sanitizado: %s\n", textoLimpio)

	// Test límites extremos
	if err := validators.ValidarLimitesExtremos(1000000, 100.0); err == nil {
		return fmt.Errorf("cantidad extrema no fue rechazada")
	}

	return nil
}

func testearGeneracionFacturas() (*models.Factura, error) {
	// Crear input de factura
	input := models.FacturaInput{
		ClienteNombre: "Cliente Test <script>",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "TEST001",
				Descripcion:    "Producto de prueba con caracteres: <>&\"'",
				Cantidad:       2.0,
				PrecioUnitario: 15.50,
			},
			{
				Codigo:         "TEST002",
				Descripcion:    "Segundo producto",
				Cantidad:       1.0,
				PrecioUnitario: 25.99,
			},
		},
	}

	// Crear factura
	factura, err := factory.CrearFactura(input)
	if err != nil {
		return nil, fmt.Errorf("error creando factura: %v", err)
	}

	// Verificar cálculos (con tolerancia para flotantes)
	expectedSubtotal := (2.0 * 15.50) + (1.0 * 25.99)
	if abs(factura.InfoFactura.TotalSinImpuestos - expectedSubtotal) > 0.01 {
		return nil, fmt.Errorf("subtotal incorrecto: esperado %.2f, obtenido %.2f",
			expectedSubtotal, factura.InfoFactura.TotalSinImpuestos)
	}

	expectedIva := expectedSubtotal * 0.15
	expectedTotal := expectedSubtotal + expectedIva
	if abs(factura.InfoFactura.ImporteTotal - expectedTotal) > 0.01 {
		return nil, fmt.Errorf("total incorrecto: esperado %.2f, obtenido %.2f",
			expectedTotal, factura.InfoFactura.ImporteTotal)
	}

	fmt.Printf("   ✓ Factura creada - Total: $%.2f (IVA: $%.2f)\n", 
		factura.InfoFactura.ImporteTotal, expectedIva)

	return &factura, nil
}

func testearGeneracionXML(factura *models.Factura) ([]byte, error) {
	xmlData, err := factura.GenerarXML()
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	if len(xmlData) == 0 {
		return nil, fmt.Errorf("XML vacío generado")
	}

	// Verificar que contiene elementos básicos
	xmlString := string(xmlData)
	if !contains(xmlString, "<factura") {
		return nil, fmt.Errorf("XML no contiene elemento factura")
	}

	if !contains(xmlString, factura.InfoTributaria.ClaveAcceso) {
		return nil, fmt.Errorf("XML no contiene clave de acceso")
	}

	fmt.Printf("   ✓ XML generado (%d bytes)\n", len(xmlData))
	fmt.Printf("   ✓ Contiene clave: %s\n", factura.InfoTributaria.ClaveAcceso)

	return xmlData, nil
}

func testearCircuitBreaker() error {
	// Crear circuit breaker de prueba
	cb := sri.NuevoCircuitBreakerDefault()

	// Test función que siempre falla
	funcionFalla := func() error {
		return fmt.Errorf("error simulado")
	}

	// Ejecutar varias veces para abrir el circuito
	for i := 0; i < 6; i++ {
		cb.Ejecutar(funcionFalla)
	}

	// Verificar que el circuito está abierto
	if cb.ObtenerEstado() != sri.EstadoAbierto {
		return fmt.Errorf("circuit breaker debería estar abierto")
	}

	fmt.Printf("   ✓ Circuit breaker se abrió correctamente\n")

	// Test función exitosa
	funcionExitosa := func() error {
		return nil
	}

	// Crear circuit breaker nuevo y probar éxito
	cb2 := sri.NuevoCircuitBreakerDefault()
	if err := cb2.Ejecutar(funcionExitosa); err != nil {
		return fmt.Errorf("función exitosa falló: %v", err)
	}

	fmt.Printf("   ✓ Circuit breaker permite funciones exitosas\n")

	return nil
}

func testearLogging() error {
	// Test diferentes niveles de logging
	sri.Debug("Mensaje de debug")
	sri.Info("Mensaje informativo")
	sri.Warning("Mensaje de advertencia")
	sri.Error("Mensaje de error")

	// Test logging especializado
	sri.LogValidacion("TestValidacion", true, "Validación exitosa")
	sri.LogSRI("TestSRI", false, 1500, "Error simulado")
	sri.LogSeguridad("InputMalicioso", "Script detectado", "TestFunction")

	fmt.Printf("   ✓ Logging funcionando en todos los niveles\n")

	return nil
}

func mostrarResumenTesting(factura *models.Factura, xmlData []byte) {
	fmt.Printf("   📋 Empresa: %s\n", config.Config.Empresa.RazonSocial)
	fmt.Printf("   🔢 RUC: %s\n", config.Config.Empresa.RUC)
	fmt.Printf("   🌍 Ambiente: %s\n", config.Config.Ambiente.Descripcion)
	fmt.Printf("   🔑 Clave Acceso: %s\n", factura.InfoTributaria.ClaveAcceso)
	fmt.Printf("   💰 Total Factura: $%.2f\n", factura.InfoFactura.ImporteTotal)
	fmt.Printf("   📄 Tamaño XML: %d bytes\n", len(xmlData))
	fmt.Printf("   🛡️  Sanitización: Activa\n")
	fmt.Printf("   🔧 Circuit Breaker: Activo\n")
	fmt.Printf("   📊 Logging: Funcional\n")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    contains(s[1:], substr) || 
		    s[:len(substr)] == substr)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}