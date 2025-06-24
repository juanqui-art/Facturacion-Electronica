package main

import (
	"fmt"
	"os"
	"strings"

	"go-facturacion-sri/api"
	"go-facturacion-sri/config"
	"go-facturacion-sri/database"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/sri"
)

func main() {
	// Cargar configuración al inicio
	fmt.Println("📋 Cargando configuración del sistema...")
	err := config.CargarConfiguracion("config/desarrollo.json")
	if err != nil {
		fmt.Printf("⚠️  Error cargando configuración: %v\n", err)
		fmt.Println("📦 Usando configuración por defecto...")
		config.CargarConfiguracionPorDefecto()
	} else {
		fmt.Printf("✅ Configuración cargada: %s\n", config.Config.Empresa.RazonSocial)
		fmt.Printf("🏢 Ambiente: %s\n", config.Config.Ambiente.Descripcion)
	}

	// Verificar si queremos ejecutar en modo API, SRI demo o modo demo
	if len(os.Args) > 1 && os.Args[1] == "api" {
		// Modo API: Iniciar servidor HTTP
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🌐 INICIANDO SERVIDOR HTTP API")
		fmt.Println(strings.Repeat("=", 50))

		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		server := api.NewServer(port)
		if err := server.Start(); err != nil {
			fmt.Printf("❌ Error iniciando servidor: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Modo SRI Demo: Demostrar funcionalidades de integración SRI
	if len(os.Args) > 1 && os.Args[1] == "sri" {
		sri.DemoSRI()
		return
	}

	// Modo SOAP Demo: Demostrar cliente SOAP
	if len(os.Args) > 1 && os.Args[1] == "soap" {
		sri.DemoSOAPClient()
		return
	}

	// Modo Database Demo: Demostrar sistema de base de datos
	if len(os.Args) > 1 && os.Args[1] == "database" {
		database.DemoDatabase()
		return
	}

	// Modo Test SRI: Tests de integración con SRI real
	if len(os.Args) > 1 && os.Args[1] == "test-sri" {
		sri.DemoTestIntegracion()
		return
	}

	// Modo Certificación: Guía de certificación SRI
	if len(os.Args) > 1 && os.Args[1] == "certificacion" {
		sri.MostrarGuiaCertificacion()
		return
	}

	// Modo demo: Ejecutar ejemplos y pruebas
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🧪 MODO DEMO - Ejecutando ejemplos")
	fmt.Println("💡 Para modo API: go run main.go test_validaciones.go api")
	fmt.Println("🇪🇨 Para demo SRI: go run main.go test_validaciones.go sri")
	fmt.Println("🌐 Para demo SOAP: go run main.go test_validaciones.go soap")
	fmt.Println("🗄️  Para demo DB: go run main.go test_validaciones.go database")
	fmt.Println("🧪 Para test SRI: go run main.go test_validaciones.go test-sri")
	fmt.Println("📋 Para certificación: go run main.go test_validaciones.go certificacion")
	fmt.Println(strings.Repeat("=", 50))

	// Primero, ejecutar pruebas de validación
	probarValidaciones()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🚀 GENERANDO FACTURA PRINCIPAL")
	fmt.Println(strings.Repeat("=", 50))

	// Crear datos de factura - ¡Ahora con múltiples productos!
	facturaData := models.FacturaInput{
		ClienteNombre: "JUAN CARLOS PEREZ",
		ClienteCedula: "1713175071", // Cédula válida para Ecuador
		Productos: []models.ProductoInput{
			{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell Inspiron 15",
				Cantidad:       2.0,
				PrecioUnitario: 450.00,
			},
			{
				Codigo:         "MOUSE001",
				Descripcion:    "Mouse Inalámbrico Logitech",
				Cantidad:       3.0,
				PrecioUnitario: 25.00,
			},
			{
				Codigo:         "TECLADO001",
				Descripcion:    "Teclado Mecánico RGB",
				Cantidad:       1.0,
				PrecioUnitario: 85.00,
			},
		},
	}

	// Generar factura usando nuestra función factory
	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		fmt.Printf("Error al crear la factura: %v\n", err)
		return
	}

	// Mostrar resumen usando el método de la factura
	factura.MostrarResumen()

	// Generar XML usando el método de la factura
	xmlData, err := factura.GenerarXML()
	if err != nil {
		fmt.Printf("Error generando XML: %v\n", err)
		return
	}

	fmt.Println("=== XML GENERADO ===")
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", xmlData)
}
