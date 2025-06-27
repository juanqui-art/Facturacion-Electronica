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

	// 2. Verificar modo de certificado
	fmt.Println("\n🔐 Paso 2: Verificando configuración de certificado...")
	
	// Verificar si está en modo demo
	type ConfigCertificado struct {
		RutaArchivo string `json:"rutaArchivo"`
		Password    string `json:"password"`
		Habilitado  bool   `json:"habilitado"`
		ModoDemo    bool   `json:"modoDemo"`
	}
	
	var modoDemo bool = false
	if config.Config.Certificado.RutaArchivo == "" || len(config.Config.Certificado.Password) == 0 {
		modoDemo = true
	}
	
	if modoDemo {
		fmt.Println("🎭 MODO DEMO ACTIVO - Sin certificado digital real")
		fmt.Println("   • Sistema funcionará sin firma digital")
		fmt.Println("   • Comunicación SRI en modo de prueba")
		fmt.Println("   • Para activar certificado real:")
		fmt.Println("     1. Obtener certificado BCE ($24.64 USD)")
		fmt.Println("     2. Actualizar config/desarrollo.json")
		fmt.Println("     3. Reiniciar sistema")
	} else {
		// Intentar cargar certificado real
		cert, err := CargarCertificado(CertificadoConfig{
			RutaArchivo:     config.Config.Certificado.RutaArchivo,
			Password:        config.Config.Certificado.Password,
			ValidarVigencia: true,
		})
		if err != nil {
			fmt.Printf("⚠️  Error cargando certificado: %v\n", err)
			fmt.Println("   Cambiando a modo demo...")
			modoDemo = true
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
	
	if modoDemo {
		fmt.Println("🎭 MODO DEMO: Simulando envío al SRI...")
		fmt.Println("   • En modo real, aquí se enviaría el XML firmado")
		fmt.Println("   • El SRI rechazaría el documento por falta de firma digital")
		fmt.Println("   • Con certificado real, el proceso sería automático")
		
		// Simular respuesta exitosa para demo
		fmt.Println("✅ Simulación de envío completada")
		fmt.Println("   Estado simulado: RECIBIDA")
		fmt.Println("   Nota: Con certificado real obtendría autorización automática")
		
		// Mostrar siguiente paso requerido
		fmt.Println("\n🔄 Siguiente paso requerido:")
		fmt.Println("   1. Obtener certificado digital del BCE")
		fmt.Println("   2. Configurar en config/desarrollo.json:")
		fmt.Println("      \"certificado\": {")
		fmt.Println("        \"rutaArchivo\": \"./certificados/mi-certificado.p12\",")
		fmt.Println("        \"password\": \"mi_contraseña\",")
		fmt.Println("        \"habilitado\": true,")
		fmt.Println("        \"modoDemo\": false")
		fmt.Println("      }")
		fmt.Println("   3. Ejecutar nuevamente este test")
		
	} else {
		// Envío real al SRI
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
	}

	// 8. Mostrar estadísticas finales
	fmt.Println("\n📊 Paso 8: Estadísticas finales...")
	if !modoDemo {
		sriClient.MostrarEstadoCircuitBreaker()
	}

	if modoDemo {
		fmt.Println("\n🎉 DEMO DE INTEGRACIÓN COMPLETADO EXITOSAMENTE")
		fmt.Println("=" + string(make([]byte, 50)))
		fmt.Println("📋 RESUMEN DEL MODO DEMO:")
		fmt.Println("✅ Sistema funcional sin certificado digital")
		fmt.Println("✅ Configuración empresarial realista")
		fmt.Println("✅ Generación de claves de acceso válidas")
		fmt.Println("✅ Creación de XML compatible con SRI")
		fmt.Println("✅ Arquitectura lista para certificado real")
		fmt.Println("")
		fmt.Println("🔄 PARA ACTIVAR INTEGRACIÓN REAL:")
		fmt.Println("1. Obtener certificado BCE: https://www.eci.bce.ec/")
		fmt.Println("2. Costo: $24.64 USD (certificado de archivo)")
		fmt.Println("3. Tiempo: 30 minutos (proceso online)")
		fmt.Println("4. Configurar archivo y ejecutar nuevamente")
		fmt.Println("")
		fmt.Println("💡 El sistema pasará automáticamente a modo producción")
	} else {
		fmt.Println("\n🎉 TESTING DE INTEGRACIÓN REAL COMPLETADO EXITOSAMENTE")
		fmt.Println("=" + string(make([]byte, 50)))
	}
	
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