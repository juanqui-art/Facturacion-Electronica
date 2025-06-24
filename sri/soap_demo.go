// Package sri Demo específico del cliente SOAP para SRI Ecuador
package sri

import (
	"fmt"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"strings"
	"time"
)

// DemoSOAPClient demuestra el uso del cliente SOAP
func DemoSOAPClient() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🌐 DEMO CLIENTE SOAP SRI ECUADOR")
	fmt.Println(strings.Repeat("=", 60))

	// Demo 1: Crear cliente SOAP
	fmt.Println("\n1️⃣ CREACIÓN DE CLIENTE SOAP")
	fmt.Println(strings.Repeat("-", 40))

	client := NewSOAPClient(Pruebas)
	fmt.Printf("✅ Cliente SOAP creado para ambiente: %s\n", obtenerNombreAmbiente(client.Ambiente))
	fmt.Printf("⏱️  Timeout configurado: %d segundos\n", client.TimeoutSegundos)

	// Demo 2: Mostrar endpoints
	fmt.Println("\n2️⃣ ENDPOINTS DISPONIBLES")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("🔗 Recepción (Certificación): %s\n", EndpointRecepcionCertificacion)
	fmt.Printf("🔗 Autorización (Certificación): %s\n", EndpointAutorizacionCertificacion)
	fmt.Printf("🔗 Recepción (Producción): %s\n", EndpointRecepcionProduccion)
	fmt.Printf("🔗 Autorización (Producción): %s\n", EndpointAutorizacionProduccion)

	// Demo 3: Crear factura de ejemplo para envío
	fmt.Println("\n3️⃣ GENERACIÓN DE FACTURA PARA ENVÍO")
	fmt.Println(strings.Repeat("-", 40))

	facturaData := models.FacturaInput{
		ClienteNombre: "EMPRESA DEMO SRI",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "DEMO001",
				Descripcion:    "Producto Demo SOAP",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		fmt.Printf("❌ Error creando factura: %v\n", err)
		return
	}

	xmlFactura, err := factura.GenerarXML()
	if err != nil {
		fmt.Printf("❌ Error generando XML: %v\n", err)
		return
	}

	fmt.Printf("✅ Factura generada para cliente: %s\n", facturaData.ClienteNombre)
	fmt.Printf("💰 Total factura: $%.2f\n", factura.InfoFactura.ImporteTotal)

	// Demo 4: Generar clave de acceso
	fmt.Println("\n4️⃣ GENERACIÓN DE CLAVE DE ACCESO")
	fmt.Println(strings.Repeat("-", 40))

	claveConfig := ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      EmisionNormal,
	}

	claveAcceso, err := GenerarClaveAcceso(claveConfig)
	if err != nil {
		fmt.Printf("❌ Error generando clave: %v\n", err)
		return
	}

	fmt.Printf("🔑 Clave de acceso: %s\n", FormatearClaveAcceso(claveAcceso))

	// Demo 5: Simulación de envío (sin envío real)
	fmt.Println("\n5️⃣ SIMULACIÓN DE COMUNICACIÓN SOAP")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📤 Proceso que se ejecutaría con SRI real:")
	fmt.Println("   1. Codificar XML en Base64")
	fmt.Println("   2. Crear sobre SOAP")
	fmt.Println("   3. Enviar a endpoint de recepción")
	fmt.Println("   4. Procesar respuesta de recepción")
	fmt.Println("   5. Consultar autorización")
	fmt.Println("   6. Obtener comprobante autorizado")

	// Demo 6: Mostrar estructura XML que se enviaría
	fmt.Println("\n6️⃣ ESTRUCTURA XML QUE SE ENVIARÍA")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("📄 Tamaño del XML: %d bytes\n", len(xmlFactura))
	fmt.Println("📋 Contenido XML (primeros 200 caracteres):")
	xmlPreview := string(xmlFactura)
	if len(xmlPreview) > 200 {
		xmlPreview = xmlPreview[:200] + "..."
	}
	fmt.Println(xmlPreview)

	// Demo 7: Flujo completo teórico
	fmt.Println("\n7️⃣ FLUJO COMPLETO TEÓRICO")
	fmt.Println(strings.Repeat("-", 40))

	fmt.Println("🔄 Simulando flujo completo...")

	// Simular envío
	fmt.Println("📤 1. Enviando comprobante al SRI... ⏳")
	time.Sleep(1 * time.Second)
	fmt.Println("✅    Comprobante RECIBIDO por SRI")

	// Simular procesamiento
	fmt.Println("⚙️  2. SRI procesando comprobante... ⏳")
	time.Sleep(2 * time.Second)
	fmt.Println("✅    Comprobante PROCESADO")

	// Simular autorización
	fmt.Println("🔐 3. Consultando autorización... ⏳")
	time.Sleep(1 * time.Second)
	fmt.Println("✅    Comprobante AUTORIZADO")

	// Resultado final
	fmt.Println("\n🎉 RESULTADO FINAL (SIMULADO)")
	fmt.Println(strings.Repeat("-", 40))
	autorizacion := SimularAutorizacionSRI(claveAcceso, Pruebas)
	fmt.Printf("📝 Número de Autorización: %s\n", autorizacion.NumeroAutorizacion)
	fmt.Printf("📅 Fecha de Autorización: %s\n", autorizacion.FechaAutorizacion.Format("02/01/2006 15:04:05"))
	fmt.Printf("✅ Estado: %s\n", autorizacion.Estado)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Demo SOAP Cliente completado")
	fmt.Println("💡 Para usar con SRI real:")
	fmt.Println("   • Certificado digital válido (.p12)")
	fmt.Println("   • Conexión a internet estable")
	fmt.Println("   • RUC registrado en SRI")
	fmt.Println("   • Cambiar ambiente a Producción para validez fiscal")
	fmt.Println(strings.Repeat("=", 60))
}

// DemoSOAPOperaciones demuestra cada operación SOAP individualmente
func DemoSOAPOperaciones() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("⚙️  DEMO OPERACIONES SOAP INDIVIDUALES")
	fmt.Println(strings.Repeat("=", 50))

	client := NewSOAPClient(Pruebas)

	// Demo operación de recepción (mock)
	fmt.Println("\n📥 OPERACIÓN: RECEPCIÓN DE COMPROBANTES")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📋 Método: EnviarComprobante()")
	fmt.Println("📡 Endpoint: RecepcionComprobantesOffline")
	fmt.Println("📄 Input: XML del comprobante en Base64")
	fmt.Println("📄 Output: RespuestaSolicitud con estado")

	fmt.Println("\n🔧 Estructura de la petición SOAP:")
	fmt.Println(`
	<soap:Envelope xmlns:soap="...">
	  <soap:Body>
	    <sri:validarComprobante>
	      <xml>[XML_EN_BASE64]</xml>
	    </sri:validarComprobante>
	  </soap:Body>
	</soap:Envelope>`)

	// Demo operación de autorización (mock)
	fmt.Println("\n🔐 OPERACIÓN: AUTORIZACIÓN DE COMPROBANTES")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📋 Método: ConsultarAutorizacion()")
	fmt.Println("📡 Endpoint: AutorizacionComprobantesOffline")
	fmt.Println("📄 Input: Clave de acceso (49 dígitos)")
	fmt.Println("📄 Output: RespuestaComprobante con autorización")

	fmt.Println("\n🔧 Estructura de la petición SOAP:")
	fmt.Println(`
	<soap:Envelope xmlns:soap="...">
	  <soap:Body>
	    <sri:autorizacionComprobante>
	      <claveAccesoComprobante>[CLAVE_49_DIGITOS]</claveAccesoComprobante>
	    </sri:autorizacionComprobante>
	  </soap:Body>
	</soap:Envelope>`)

	// Demo operación completa (mock)
	fmt.Println("\n🔄 OPERACIÓN: PROCESAMIENTO COMPLETO")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📋 Método: ProcesarComprobanteCompleto()")
	fmt.Println("📡 Combina: Recepción + Autorización + Reintentos")
	fmt.Println("📄 Input: XML del comprobante + Clave de acceso")
	fmt.Println("📄 Output: AutorizacionSRI con comprobante firmado")
	fmt.Println("🔄 Incluye: Reintentos automáticos y manejo de errores")

	fmt.Printf("\n✅ Cliente configurado para ambiente: %s\n", obtenerNombreAmbiente(client.Ambiente))
	fmt.Printf("⏱️  Timeout configurado: %d segundos\n", client.TimeoutSegundos)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ Demo operaciones SOAP completado")
	fmt.Println(strings.Repeat("=", 50))
}

// DemoSOAPTesting demuestra las capacidades de testing del cliente SOAP
func DemoSOAPTesting() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🧪 DEMO TESTING CLIENTE SOAP")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("\n📊 Tests implementados:")
	fmt.Println("   ✅ TestNewSOAPClient - Creación de cliente")
	fmt.Println("   ✅ TestSOAPClientEndpoints - Validación de endpoints")
	fmt.Println("   ✅ TestParsearRespuestaRecepcionMock - Parsing recepción")
	fmt.Println("   ✅ TestParsearRespuestaAutorizacionMock - Parsing autorización")
	fmt.Println("   ✅ TestSOAPClientTimeout - Configuración timeout")
	fmt.Println("   ✅ TestConstantesEndpoints - Validación constantes")
	fmt.Println("   ✅ BenchmarkNewSOAPClient - Performance")

	fmt.Println("\n🔧 Mock responses para testing:")
	fmt.Println("   • Respuesta SOAP de recepción simulada")
	fmt.Println("   • Respuesta SOAP de autorización simulada")
	fmt.Println("   • Validación de parsing XML")
	fmt.Println("   • Verificación de estructura de datos")

	fmt.Println("\n🚀 Para ejecutar tests:")
	fmt.Println("   go test ./sri -v")
	fmt.Println("   go test ./sri -bench=.")
	fmt.Println("   go test ./sri -cover")

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ Demo testing completado")
	fmt.Println(strings.Repeat("=", 50))
}
