// Demo del sistema de integración SRI
package sri

import (
	"fmt"
	"strings"
	"time"
)

// DemoSRI ejecuta una demostración completa del sistema SRI
func DemoSRI() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🇪🇨 DEMO SISTEMA INTEGRACIÓN SRI ECUADOR")
	fmt.Println(strings.Repeat("=", 60))

	// Demo 1: Generar clave de acceso
	fmt.Println("\n1️⃣ GENERACIÓN DE CLAVE DE ACCESO")
	fmt.Println(strings.Repeat("-", 40))
	
	config := ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      EmisionNormal,
	}

	claveAcceso, err := GenerarClaveAcceso(config)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Clave de acceso generada: %s\n", FormatearClaveAcceso(claveAcceso))
	
	// Mostrar información detallada
	MostrarInformacionClaveAcceso(claveAcceso)

	// Demo 2: Simular autorización SRI
	fmt.Println("\n2️⃣ SIMULACIÓN AUTORIZACIÓN SRI")
	fmt.Println(strings.Repeat("-", 40))
	
	autorizacion := SimularAutorizacionSRI(claveAcceso, Pruebas)
	fmt.Printf("📝 Número de Autorización: %s\n", autorizacion.NumeroAutorizacion)
	fmt.Printf("📅 Fecha de Autorización: %s\n", autorizacion.FechaAutorizacion.Format("02/01/2006 15:04:05"))
	fmt.Printf("✅ Estado: %s\n", autorizacion.Estado)
	fmt.Printf("🌍 Ambiente: %s\n", autorizacion.Ambiente)

	// Demo 3: Información sobre certificados
	fmt.Println("\n3️⃣ SISTEMA DE CERTIFICADOS DIGITALES")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📋 Para firma electrónica se requiere:")
	fmt.Println("   • Certificado PKCS#12 (.p12) del Banco Central del Ecuador")
	fmt.Println("   • Contraseña del certificado")
	fmt.Println("   • Validación de vigencia")
	fmt.Println("   • Implementación XAdES-BES")
	
	// Ejemplo de configuración de certificado
	certConfig := CertificadoConfig{
		RutaArchivo:     "/ruta/al/certificado.p12",
		Password:        "password_certificado",
		ValidarVigencia: true,
		ValidarCadena:   true,
	}
	
	fmt.Printf("\n📂 Configuración de certificado ejemplo:\n")
	fmt.Printf("   Archivo: %s\n", certConfig.RutaArchivo)
	fmt.Printf("   Validar vigencia: %v\n", certConfig.ValidarVigencia)
	fmt.Printf("   Validar cadena: %v\n", certConfig.ValidarCadena)

	// Demo 4: Tipos de comprobantes soportados
	fmt.Println("\n4️⃣ TIPOS DE COMPROBANTES SOPORTADOS")
	fmt.Println(strings.Repeat("-", 40))
	
	tiposComprobantes := []struct {
		tipo   TipoComprobante
		nombre string
	}{
		{Factura, "Factura"},
		{NotaCredito, "Nota de Crédito"},
		{NotaDebito, "Nota de Débito"},
		{GuiaRemision, "Guía de Remisión"},
		{ComprobanteRetencion, "Comprobante de Retención"},
		{LiquidacionCompra, "Liquidación de Compra"},
	}

	for _, tc := range tiposComprobantes {
		fmt.Printf("   %s: %s\n", tc.tipo.String(), tc.nombre)
	}

	// Demo 5: Ambientes disponibles
	fmt.Println("\n5️⃣ AMBIENTES DISPONIBLES")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   %s: Ambiente de Pruebas\n", Pruebas.String())
	fmt.Printf("   %s: Ambiente de Producción\n", Produccion.String())

	// Demo 6: Cliente SOAP SRI
	fmt.Println("\n6️⃣ CLIENTE SOAP PARA COMUNICACIÓN CON SRI")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("🌐 Endpoints disponibles:")
	fmt.Printf("   Certificación: %s\n", EndpointRecepcionCertificacion)
	fmt.Printf("   Producción: %s\n", EndpointRecepcionProduccion)
	
	fmt.Println("\n📡 Servicios SOAP implementados:")
	fmt.Println("   • Recepción de comprobantes")
	fmt.Println("   • Autorización de comprobantes")
	fmt.Println("   • Consulta de estado")
	fmt.Println("   • Procesamiento completo automático")

	// Demo 7: Proceso completo de facturación electrónica
	fmt.Println("\n7️⃣ PROCESO COMPLETO FACTURACIÓN ELECTRÓNICA")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("📋 Pasos del proceso:")
	fmt.Println("   1. Crear factura con datos del cliente y productos")
	fmt.Println("   2. Generar XML según especificaciones SRI")
	fmt.Println("   3. Generar clave de acceso única")
	fmt.Println("   4. Firmar electrónicamente con XAdES-BES")
	fmt.Println("   5. Enviar al SRI para autorización (SOAP)")
	fmt.Println("   6. Recibir autorización y número de autorización")
	fmt.Println("   7. Generar RIDE (Representación Impresa)")
	fmt.Println("   8. Enviar RIDE al cliente")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Demo completado exitosamente")
	fmt.Println("💡 Para implementación completa se requiere:")
	fmt.Println("   • Certificado digital válido del BCE")
	fmt.Println("   • Acceso a servicios web del SRI")
	fmt.Println("   • Configuración de ambiente de producción")
	fmt.Println(strings.Repeat("=", 60))
}

// DemoClaveAccesoPersonalizada permite generar claves con parámetros personalizados
func DemoClaveAccesoPersonalizada(ruc string, serie string, secuencial string) {
	fmt.Println("\n🎯 DEMO CLAVE DE ACCESO PERSONALIZADA")
	fmt.Println(strings.Repeat("-", 45))

	config := ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  Factura,
		RUCEmisor:        ruc,
		Ambiente:         Pruebas,
		Serie:            serie,
		NumeroSecuencial: secuencial,
		TipoEmision:      EmisionNormal,
	}

	clave, err := GenerarClaveAcceso(config)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Clave generada: %s\n", FormatearClaveAcceso(clave))
	MostrarInformacionClaveAcceso(clave)
}

// DemoValidacionClaves demuestra la validación de claves de acceso
func DemoValidacionClaves() {
	fmt.Println("\n🔍 DEMO VALIDACIÓN DE CLAVES DE ACCESO")
	fmt.Println(strings.Repeat("-", 45))

	// Ejemplos de claves para validar
	claves := []struct {
		clave       string
		descripcion string
	}{
		{
			clave:       "2306202401179214673900110010010000000011234567891",
			descripcion: "Clave válida de ejemplo",
		},
		{
			clave:       "230620240117921467390011001001000000001123456789",
			descripcion: "Clave muy corta (48 dígitos)",
		},
		{
			clave:       "2306202401179214673900110010010000000011234567890",
			descripcion: "Clave con dígito verificador incorrecto",
		},
		{
			clave:       "23/06/2024-01-1792146739001-1-001001-000000001-12345678-9-1",
			descripcion: "Clave con formato incorrecto (contiene caracteres)",
		},
	}

	for i, ejemplo := range claves {
		fmt.Printf("\n%d. %s\n", i+1, ejemplo.descripcion)
		fmt.Printf("   Clave: %s\n", ejemplo.clave)
		
		err := ValidarClaveAcceso(ejemplo.clave)
		if err != nil {
			fmt.Printf("   ❌ Resultado: %v\n", err)
		} else {
			fmt.Printf("   ✅ Resultado: Clave válida\n")
			fmt.Printf("   📋 Formato: %s\n", FormatearClaveAcceso(ejemplo.clave))
		}
	}
}