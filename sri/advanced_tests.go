package sri

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
)

// TestErroresSOAPAvanzados tests avanzados de manejo de errores
func TestErroresSOAPAvanzados(t *testing.T) {
	tests := []struct {
		name           string
		mensajeError   string
		codigoHTTP     int
		esRecuperable  bool
		tipoEsperado   TipoErrorSRI
	}{
		{
			name:          "Error de clave de acceso registrada",
			mensajeError:  "CLAVE-01: Clave de acceso ya registrada",
			codigoHTTP:    400,
			esRecuperable: false,
			tipoEsperado:  ErrorClaveAcceso,
		},
		{
			name:          "Error de timeout HTTP",
			mensajeError:  "Request timeout",
			codigoHTTP:    504,
			esRecuperable: true,
			tipoEsperado:  ErrorTimeout,
		},
		{
			name:          "Error de certificado expirado",
			mensajeError:  "CERT-01: Certificado expirado",
			codigoHTTP:    403,
			esRecuperable: false,
			tipoEsperado:  ErrorCertificado,
		},
		{
			name:          "Error de sistema SRI",
			mensajeError:  "SRI-01: Sistema en mantenimiento",
			codigoHTTP:    503,
			esRecuperable: true,
			tipoEsperado:  ErrorSistema,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorSRI := ParsearErrorSRI(tt.mensajeError, tt.codigoHTTP)
			
			if errorSRI.Tipo != tt.tipoEsperado {
				t.Errorf("Tipo de error esperado: %v, obtenido: %v", tt.tipoEsperado, errorSRI.Tipo)
			}
			
			if errorSRI.Recuperable != tt.esRecuperable {
				t.Errorf("Recuperabilidad esperada: %v, obtenida: %v", tt.esRecuperable, errorSRI.Recuperable)
			}
			
			if errorSRI.SugerenciaFix == "" {
				t.Error("Sugerencia de solución no puede estar vacía")
			}
		})
	}
}

// TestReintentoLogic tests de lógica de reintentos
func TestReintentoLogic(t *testing.T) {
	// Test con función que siempre falla
	funcionSiempreFalla := func() error {
		return CrearErrorConexion("Conexión falló")
	}
	
	config := ConfigReintento{
		MaxIntentos:      3,
		TiempoBase:       10 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     100 * time.Millisecond,
		JitterMaximo:     5 * time.Millisecond,
		SoloRecuperables: true,
	}
	
	resultado := EjecutarConReintento(funcionSiempreFalla, config)
	
	if resultado.Exitoso {
		t.Error("El resultado debería ser fallido")
	}
	
	if resultado.IntentosRealizados != 3 {
		t.Errorf("Intentos esperados: 3, obtenidos: %d", resultado.IntentosRealizados)
	}
	
	if len(resultado.Errores) != 3 {
		t.Errorf("Errores esperados: 3, obtenidos: %d", len(resultado.Errores))
	}
}

// TestReintentoExitoso test de reintento que eventualmente tiene éxito
func TestReintentoExitoso(t *testing.T) {
	intentos := 0
	funcionEventualmenteExitosa := func() error {
		intentos++
		if intentos < 3 {
			return CrearErrorConexion("Fallo temporal")
		}
		return nil // Éxito en el tercer intento
	}
	
	config := ConfigReintentoDefault
	config.TiempoBase = 10 * time.Millisecond
	config.TiempoMaximo = 50 * time.Millisecond
	config.JitterMaximo = 5 * time.Millisecond
	
	resultado := EjecutarConReintento(funcionEventualmenteExitosa, config)
	
	if !resultado.Exitoso {
		t.Error("El resultado debería ser exitoso")
	}
	
	if resultado.IntentosRealizados != 3 {
		t.Errorf("Intentos esperados: 3, obtenidos: %d", resultado.IntentosRealizados)
	}
}

// TestIntegracionFacturaCompleta test de integración completa sin SRI real
func TestIntegracionFacturaCompleta(t *testing.T) {
	// Crear factura
	facturaData := models.FacturaInput{
		ClienteNombre: "EMPRESA TEST INTEGRACION",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "INTEG001",
				Descripcion:    "Producto Integración Test",
				Cantidad:       2.0,
				PrecioUnitario: 75.50,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		t.Fatalf("Error creando factura: %v", err)
	}

	// Generar XML
	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("Error generando XML: %v", err)
	}

	// Generar clave de acceso
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
		t.Fatalf("Error generando clave de acceso: %v", err)
	}

	// Validar clave de acceso
	err = ValidarClaveAcceso(claveAcceso)
	if err != nil {
		t.Fatalf("Clave de acceso inválida: %v", err)
	}

	// Crear cliente SOAP
	client := NewSOAPClient(Pruebas)
	
	// Test de estructura (sin envío real)
	if client == nil {
		t.Fatal("Cliente SOAP no fue creado")
	}

	if client.Ambiente != Pruebas {
		t.Errorf("Ambiente esperado: Pruebas, obtenido: %v", client.Ambiente)
	}

	// Verificar que el XML contiene la clave de acceso
	xmlString := string(xmlData)
	if !strings.Contains(xmlString, claveAcceso) {
		t.Error("XML no contiene la clave de acceso generada")
	}

	t.Logf("✅ Test de integración completado exitosamente")
	t.Logf("📄 XML generado: %d bytes", len(xmlData))
	t.Logf("🔑 Clave de acceso: %s", FormatearClaveAcceso(claveAcceso))
	t.Logf("💰 Total factura: $%.2f", factura.InfoFactura.ImporteTotal)
}

// TestValidacionesExtendidas tests de validaciones extendidas
func TestValidacionesExtendidas(t *testing.T) {
	tests := []struct {
		name         string
		claveAcceso  string
		debeSerValida bool
	}{
		{
			name:         "Clave válida generada",
			claveAcceso:  "",
			debeSerValida: true,
		},
		{
			name:         "Clave muy corta",
			claveAcceso:  "123456789012345678901234567890123456789012345678", // 48 dígitos
			debeSerValida: false,
		},
		{
			name:         "Clave muy larga",
			claveAcceso:  "12345678901234567890123456789012345678901234567890", // 50 dígitos
			debeSerValida: false,
		},
		{
			name:         "Clave con caracteres no numéricos",
			claveAcceso:  "1234567890123456789012345678901234567890123456789A",
			debeSerValida: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claveAcceso := tt.claveAcceso
			
			// Si no se proporciona clave, generar una válida
			if claveAcceso == "" {
				config := ClaveAccesoConfig{
					FechaEmision:     time.Now(),
					TipoComprobante:  Factura,
					RUCEmisor:        "1792146739001",
					Ambiente:         Pruebas,
					Serie:            "001001",
					NumeroSecuencial: "000000001",
					TipoEmision:      EmisionNormal,
				}
				
				var err error
				claveAcceso, err = GenerarClaveAcceso(config)
				if err != nil {
					t.Fatalf("Error generando clave de acceso: %v", err)
				}
			}
			
			err := ValidarClaveAcceso(claveAcceso)
			
			if tt.debeSerValida && err != nil {
				t.Errorf("Clave debería ser válida pero falló: %v", err)
			}
			
			if !tt.debeSerValida && err == nil {
				t.Errorf("Clave debería ser inválida pero pasó la validación")
			}
		})
	}
}

// TestConfiguracionesReintento tests de diferentes configuraciones de reintento
func TestConfiguracionesReintento(t *testing.T) {
	configs := []struct {
		nombre string
		config ConfigReintento
	}{
		{"Default", ConfigReintentoDefault},
		{"Conservador", ConfigReintentoConservador},
		{"Agresivo", ConfigReintentoAgresivo},
	}

	for _, tc := range configs {
		t.Run(tc.nombre, func(t *testing.T) {
			if tc.config.MaxIntentos <= 0 {
				t.Error("MaxIntentos debe ser positivo")
			}
			
			if tc.config.TiempoBase <= 0 {
				t.Error("TiempoBase debe ser positivo")
			}
			
			if tc.config.Multiplicador <= 1.0 {
				t.Error("Multiplicador debe ser mayor a 1.0")
			}
			
			if tc.config.TiempoMaximo <= tc.config.TiempoBase {
				t.Error("TiempoMaximo debe ser mayor a TiempoBase")
			}
		})
	}
}

// BenchmarkGeneracionClaveAcceso benchmark de generación de claves
func BenchmarkGeneracionClaveAcceso(b *testing.B) {
	config := ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      EmisionNormal,
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := GenerarClaveAcceso(config)
		if err != nil {
			b.Fatalf("Error generando clave de acceso: %v", err)
		}
	}
}

// BenchmarkValidacionClaveAcceso benchmark de validación de claves
func BenchmarkValidacionClaveAcceso(b *testing.B) {
	// Generar clave válida para el benchmark
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
		b.Fatalf("Error generando clave de acceso: %v", err)
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := ValidarClaveAcceso(claveAcceso)
		if err != nil {
			b.Fatalf("Error validando clave de acceso: %v", err)
		}
	}
}

// TestCreacionMultiplesFacturas test de creación de múltiples facturas
func TestCreacionMultiplesFacturas(t *testing.T) {
	numFacturas := 10
	
	for i := 0; i < numFacturas; i++ {
		facturaData := models.FacturaInput{
			ClienteNombre: fmt.Sprintf("CLIENTE TEST %d", i+1),
			ClienteCedula: "1713175071",
			Productos: []models.ProductoInput{
				{
					Codigo:         fmt.Sprintf("PROD%03d", i+1),
					Descripcion:    fmt.Sprintf("Producto Test %d", i+1),
					Cantidad:       float64(i + 1),
					PrecioUnitario: 100.00 + float64(i*10),
				},
			},
		}

		factura, err := factory.CrearFactura(facturaData)
		if err != nil {
			t.Fatalf("Error creando factura %d: %v", i+1, err)
		}

		// Verificar que cada factura tiene datos únicos
		expectedTotal := (100.00 + float64(i*10)) * float64(i+1) * 1.15 // Con IVA
		if abs(factura.InfoFactura.ImporteTotal - expectedTotal) > 0.01 {
			t.Errorf("Total de factura %d incorrecto: esperado %.2f, obtenido %.2f", 
				i+1, expectedTotal, factura.InfoFactura.ImporteTotal)
		}
	}
	
	t.Logf("✅ Creadas %d facturas exitosamente", numFacturas)
}

// Helper function para calcular valor absoluto
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}