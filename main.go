package main

import (
	"fmt"
	"strings"
	
	"go-facturacion-sri/config"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
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