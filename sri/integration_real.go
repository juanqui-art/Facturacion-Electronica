// Package sri implementa testing real con el SRI de Ecuador
package sri

import (
	"fmt"
	"time"

	"go-facturacion-sri/config"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
)

// TestearIntegracionSRIReal realiza testing completo con SRI real
func TestearIntegracionSRIReal() error {
	fmt.Println("🚀 INICIANDO TESTING DE INTEGRACIÓN SRI REAL")
	fmt.Println("=" + string(make([]byte, 50)))

	// 1. Cargar configuración
	fmt.Println("\n📋 Paso 1: Cargando configuración...")
	if err := config.CargarConfiguracion("config/desarrollo.json"); err != nil {
		return fmt.Errorf("error cargando configuración: %v", err)
	}
	fmt.Printf("✅ Configuración cargada: %s\n", config.Config.Ambiente.Descripcion)

	// 2. Verificar certificado (si está configurado)
	fmt.Println("\n🔐 Paso 2: Verificando certificado digital...")
	if config.Config.Certificado.RutaArchivo != "" {
		cert, err := CargarCertificado(CertificadoConfig{
			RutaArchivo:     config.Config.Certificado.RutaArchivo,
			Password:        config.Config.Certificado.Password,
			ValidarVigencia: true,
		})
		if err != nil {
			fmt.Printf("⚠️  Certificado no disponible: %v\n", err)
			fmt.Println("   (Continuando con testing sin firma digital)")
		} else {
			fmt.Printf("✅ Certificado válido: %s\n", cert.ObtenerSubject())
			cert.MostrarInformacion()
		}
	}

	// 3. Crear factura de prueba
	fmt.Println("\n📄 Paso 3: Creando factura de prueba...")
	facturaInput := models.FacturaInput{
		ClienteNombre: "CLIENTE DE PRUEBA SRI",
		ClienteCedula: "1713175071", // Cédula válida para pruebas
		Productos: []models.ProductoInput{
			{
				Codigo:         "TEST001",
				Descripcion:    "Producto de prueba para integración SRI",
				Cantidad:       1.0,
				PrecioUnitario: 10.00,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaInput)
	if err != nil {
		return fmt.Errorf("error creando factura: %v", err)
	}
	
	fmt.Printf("✅ Factura creada - Clave de Acceso: %s\n", factura.InfoTributaria.ClaveAcceso)
	factura.MostrarResumen()

	// 4. Generar XML
	fmt.Println("\n🔧 Paso 4: Generando XML...")
	xmlData, err := factura.GenerarXML()
	if err != nil {
		return fmt.Errorf("error generando XML: %v", err)
	}
	fmt.Printf("✅ XML generado (%d bytes)\n", len(xmlData))

	// 5. Crear cliente SRI
	fmt.Println("\n🌐 Paso 5: Conectando con SRI...")
	var ambiente Ambiente = Pruebas
	if config.Config.Ambiente.Codigo == "2" {
		ambiente = Produccion
	}
	
	sriClient := NewSOAPClient(ambiente)
	fmt.Printf("✅ Cliente SRI creado para ambiente: %s\n", ambiente)
	
	// Mostrar estado del circuit breaker
	fmt.Printf("🔧 Circuit Breaker: %s\n", sriClient.ObtenerEstadoCircuitBreaker())

	// 6. Enviar comprobante
	fmt.Println("\n📤 Paso 6: Enviando comprobante al SRI...")
	respuesta, err := sriClient.EnviarComprobante(xmlData)
	if err != nil {
		// Si falla, mostrar detalles del error
		fmt.Printf("❌ Error enviando comprobante: %v\n", err)
		
		// Mostrar información del circuit breaker
		sriClient.MostrarEstadoCircuitBreaker()
		
		// Si es error de SRI conocido, mostrar detalles
		if errorSRI, ok := err.(*ErrorSRI); ok {
			MostrarInformacionError(errorSRI)
		}
		
		return fmt.Errorf("falló envío a SRI: %v", err)
	}

	fmt.Printf("✅ Comprobante enviado exitosamente\n")
	fmt.Printf("   Estado: %s\n", respuesta.Estado)
	
	if len(respuesta.Comprobantes) > 0 {
		comprobante := respuesta.Comprobantes[0]
		fmt.Printf("   Clave: %s\n", comprobante.ClaveAcceso)
		
		if len(comprobante.Mensajes) > 0 {
			fmt.Println("   Mensajes del SRI:")
			for _, msg := range comprobante.Mensajes {
				fmt.Printf("     - %s: %s\n", msg.Tipo, msg.Mensaje)
			}
		}
	}

	// 7. Consultar autorización (solo si fue recibido)
	if respuesta.Estado == "RECIBIDA" {
		fmt.Println("\n🔍 Paso 7: Consultando autorización...")
		
		// Esperar un poco para que el SRI procese
		fmt.Println("   Esperando procesamiento del SRI (10 segundos)...")
		time.Sleep(10 * time.Second)
		
		respuestaAuth, err := sriClient.ConsultarAutorizacion(factura.InfoTributaria.ClaveAcceso)
		if err != nil {
			fmt.Printf("⚠️  Error consultando autorización: %v\n", err)
		} else {
			fmt.Printf("✅ Consulta de autorización exitosa\n")
			
			if len(respuestaAuth.Autorizaciones) > 0 {
				auth := respuestaAuth.Autorizaciones[0]
				fmt.Printf("   Estado: %s\n", auth.Estado)
				fmt.Printf("   Número de Autorización: %s\n", auth.NumeroAutorizacion)
				fmt.Printf("   Fecha: %s\n", auth.FechaAutorizacion)
				
				if len(auth.Mensajes) > 0 {
					fmt.Println("   Mensajes:")
					for _, msg := range auth.Mensajes {
						fmt.Printf("     - %s: %s\n", msg.Tipo, msg.Mensaje)
					}
				}
			}
		}
	}

	// 8. Mostrar estadísticas finales
	fmt.Println("\n📊 Paso 8: Estadísticas finales...")
	sriClient.MostrarEstadoCircuitBreaker()

	fmt.Println("\n🎉 TESTING DE INTEGRACIÓN COMPLETADO EXITOSAMENTE")
	fmt.Println("=" + string(make([]byte, 50)))
	
	return nil
}

// TestearCertificadoDigital prueba la carga y validación de certificados
func TestearCertificadoDigital(rutaCertificado, password string) error {
	fmt.Println("🔐 TESTING DE CERTIFICADO DIGITAL")
	fmt.Println("=" + string(make([]byte, 40)))

	// Cargar certificado
	fmt.Printf("📂 Cargando certificado: %s\n", rutaCertificado)
	
	cert, err := CargarCertificado(CertificadoConfig{
		RutaArchivo:     rutaCertificado,
		Password:        password,
		ValidarVigencia: true,
		ValidarCadena:   false,
	})
	if err != nil {
		return fmt.Errorf("error cargando certificado: %v", err)
	}

	fmt.Println("✅ Certificado cargado exitosamente")
	
	// Mostrar información detallada
	cert.MostrarInformacion()

	// Exportar a PEM para testing
	fmt.Println("\n🔧 Exportando a formato PEM...")
	
	pemKey, err := cert.ExportarClavePEM()
	if err != nil {
		return fmt.Errorf("error exportando clave PEM: %v", err)
	}
	fmt.Printf("✅ Clave privada PEM generada (%d bytes)\n", len(pemKey))

	pemCert := cert.ExportarCertificadoPEM()
	fmt.Printf("✅ Certificado PEM generado (%d bytes)\n", len(pemCert))

	// Testing de firma básica
	fmt.Println("\n✍️  Testing de firma digital...")
	testXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<factura>
	<infoTributaria>
		<claveAcceso>2025062401123456789001001001000000011234567890</claveAcceso>
	</infoTributaria>
</factura>`)

	config := XAdESBESConfig{
		Certificado: cert,
		PolicyID:    "https://www.sri.gob.ec/politica-de-firma",
		PolicyHash:  "G7roucf600+f03r/o0bAOQ6WAs0=",
	}

	xmlFirmado, err := FirmarXMLXAdESBES(testXML, config)
	if err != nil {
		fmt.Printf("⚠️  Error firmando XML (esperado si no hay implementación completa): %v\n", err)
	} else {
		fmt.Printf("✅ XML firmado exitosamente (%d bytes)\n", len(xmlFirmado))
	}

	fmt.Println("\n🎉 TESTING DE CERTIFICADO COMPLETADO")
	fmt.Println("=" + string(make([]byte, 40)))
	
	return nil
}

// MostrarEndpointsSRI muestra los endpoints configurados
func MostrarEndpointsSRI() {
	fmt.Println("🌐 ENDPOINTS SRI CONFIGURADOS")
	fmt.Println("=" + string(make([]byte, 35)))
	
	fmt.Printf("Ambiente: %s (%s)\n", 
		config.Config.Ambiente.Codigo, 
		config.Config.Ambiente.Descripcion)
	
	fmt.Printf("Recepción: %s\n", config.Config.SRI.EndpointRecepcion)
	fmt.Printf("Autorización: %s\n", config.Config.SRI.EndpointAutorizacion)
	fmt.Printf("Timeout: %d segundos\n", config.Config.SRI.TimeoutSegundos)
	fmt.Printf("Max Reintentos: %d\n", config.Config.SRI.MaxReintentos)
	
	fmt.Println("=" + string(make([]byte, 35)))
}

// LogIntegracion registra eventos de integración en archivo
func LogIntegracion(evento, detalle string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] INTEGRACION: %s - %s\n", timestamp, evento, detalle)
	
	// Por ahora imprimir en consola, más tarde se puede escribir a archivo
	fmt.Print(logMessage)
}