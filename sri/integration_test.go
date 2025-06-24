// Pruebas de integración para el sistema SRI
package sri

import (
	"fmt"
	"testing"
	"time"

	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/config"
)

// TestGenerarClaveAcceso prueba la generación de claves de acceso
func TestGenerarClaveAcceso(t *testing.T) {
	tests := []struct {
		name   string
		config ClaveAccesoConfig
		valid  bool
	}{
		{
			name: "Clave válida para factura",
			config: ClaveAccesoConfig{
				FechaEmision:     time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
				TipoComprobante:  Factura,
				RUCEmisor:        "1792146739001",
				Ambiente:         Pruebas,
				Serie:            "001001",
				NumeroSecuencial: "000000001",
				CodigoNumerico:   "12345678",
				TipoEmision:      EmisionNormal,
			},
			valid: true,
		},
		{
			name: "RUC inválido",
			config: ClaveAccesoConfig{
				FechaEmision:     time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
				TipoComprobante:  Factura,
				RUCEmisor:        "123456789", // RUC muy corto
				Ambiente:         Pruebas,
				Serie:            "001001",
				NumeroSecuencial: "000000001",
				CodigoNumerico:   "12345678",
				TipoEmision:      EmisionNormal,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clave, err := GenerarClaveAcceso(tt.config)
			
			if tt.valid {
				if err != nil {
					t.Errorf("GenerarClaveAcceso() error = %v, esperaba nil", err)
					return
				}
				
				// Validar longitud
				if len(clave) != 49 {
					t.Errorf("Clave debe tener 49 dígitos, obtuvo %d", len(clave))
				}
				
				// Validar clave
				if err := ValidarClaveAcceso(clave); err != nil {
					t.Errorf("Clave generada no es válida: %v", err)
				}
				
				t.Logf("Clave generada: %s", FormatearClaveAcceso(clave))
			} else {
				if err == nil {
					t.Errorf("GenerarClaveAcceso() esperaba error, obtuvo nil")
				}
			}
		})
	}
}

// TestParsearClaveAcceso prueba el parsing de claves de acceso
func TestParsearClaveAcceso(t *testing.T) {
	// Generar una clave válida primero
	config := ClaveAccesoConfig{
		FechaEmision:     time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		CodigoNumerico:   "12345678",
		TipoEmision:      EmisionNormal,
	}

	clave, err := GenerarClaveAcceso(config)
	if err != nil {
		t.Fatalf("Error generando clave: %v", err)
	}

	// Parsear la clave
	configParseado, err := ParsearClaveAcceso(clave)
	if err != nil {
		t.Fatalf("Error parseando clave: %v", err)
	}

	// Verificar que los datos coincidan
	if configParseado.TipoComprobante != config.TipoComprobante {
		t.Errorf("TipoComprobante esperado %v, obtuvo %v", config.TipoComprobante, configParseado.TipoComprobante)
	}

	if configParseado.RUCEmisor != config.RUCEmisor {
		t.Errorf("RUCEmisor esperado %s, obtuvo %s", config.RUCEmisor, configParseado.RUCEmisor)
	}

	if configParseado.Ambiente != config.Ambiente {
		t.Errorf("Ambiente esperado %v, obtuvo %v", config.Ambiente, configParseado.Ambiente)
	}
}

// TestIntegracionFacturaConSRI prueba la integración completa
func TestIntegracionFacturaConSRI(t *testing.T) {
	// Cargar configuración
	config.CargarConfiguracionPorDefecto()

	// Crear factura de prueba
	facturaInput := models.FacturaInput{
		ClienteNombre: "JUAN CARLOS PEREZ",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell Inspiron 15",
				Cantidad:       1.0,
				PrecioUnitario: 450.00,
			},
		},
	}

	// Crear factura usando factory
	factura, err := factory.CrearFactura(facturaInput)
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

	t.Logf("✅ Factura creada exitosamente")
	t.Logf("📋 XML generado: %d bytes", len(xmlData))
	t.Logf("🔑 Clave de acceso: %s", FormatearClaveAcceso(claveAcceso))

	// Simular autorización
	autorizacion := SimularAutorizacionSRI(claveAcceso, Pruebas)
	t.Logf("📝 Autorización simulada: %s", autorizacion.Estado)
}

// TestValidarClaveAccesoDigitoVerificador prueba el cálculo del dígito verificador
func TestValidarClaveAccesoDigitoVerificador(t *testing.T) {
	// Primero generamos una clave válida
	config := ClaveAccesoConfig{
		FechaEmision:     time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		CodigoNumerico:   "12345678",
		TipoEmision:      EmisionNormal,
	}
	
	claveValida, err := GenerarClaveAcceso(config)
	if err != nil {
		t.Fatalf("Error generando clave válida: %v", err)
	}
	
	// Crear clave inválida modificando el último dígito
	claveInvalida := claveValida[:48] + "0"

	tests := []struct {
		clave    string
		esperado bool
	}{
		{
			clave:    claveValida,
			esperado: true,
		},
		{
			clave:    claveInvalida,
			esperado: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			err := ValidarClaveAcceso(tt.clave)
			
			if tt.esperado && err != nil {
				t.Errorf("ValidarClaveAcceso() error = %v, esperaba nil", err)
			}
			
			if !tt.esperado && err == nil {
				t.Errorf("ValidarClaveAcceso() esperaba error, obtuvo nil")
			}
		})
	}
}

// BenchmarkGenerarClaveAcceso mide el rendimiento de la generación de claves
func BenchmarkGenerarClaveAcceso(b *testing.B) {
	config := ClaveAccesoConfig{
		FechaEmision:     time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC),
		TipoComprobante:  Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		CodigoNumerico:   "12345678",
		TipoEmision:      EmisionNormal,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerarClaveAcceso(config)
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// TestFormatearClaveAcceso prueba el formato de claves de acceso
func TestFormatearClaveAcceso(t *testing.T) {
	clave := "2306202401179214673900110010010000000011234567891"
	esperado := "23062024-01-1792146739001-1-001001-000000001-12345678-9-1"
	
	resultado := FormatearClaveAcceso(clave)
	
	if resultado != esperado {
		t.Errorf("FormatearClaveAcceso() = %s, esperado %s", resultado, esperado)
	}
}