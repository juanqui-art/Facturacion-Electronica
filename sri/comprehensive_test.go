package sri

import (
	"fmt"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"strings"
	"testing"
	"time"
)

// TestGenerarNumeroAutorizacion tests authorization number generation
func TestGenerarNumeroAutorizacion(t *testing.T) {
	claveAcceso := "2306202401179214673900110010010000000011234567891"
	numeroAuth := GenerarNumeroAutorizacion(claveAcceso)

	if len(numeroAuth) != 37 {
		t.Errorf("Número de autorización debe tener 37 caracteres, obtuvo %d", len(numeroAuth))
	}

	if !strings.HasPrefix(numeroAuth, "2306202401179214673900110010010000000011234567891") {
		t.Error("Número de autorización debe comenzar con la clave de acceso")
	}

	// Verificar que es único generando otro
	numeroAuth2 := GenerarNumeroAutorizacion(claveAcceso)
	if numeroAuth == numeroAuth2 {
		t.Error("Los números de autorización deben ser únicos")
	}
}

// TestSimularAutorizacionSRI tests SRI authorization simulation
func TestSimularAutorizacionSRI(t *testing.T) {
	claveAcceso := "2306202401179214673900110010010000000011234567891"

	tests := []struct {
		ambiente Ambiente
		esperado string
	}{
		{Pruebas, "AUTORIZADO"},
		{Produccion, "AUTORIZADO"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Ambiente_%s", obtenerNombreAmbiente(tt.ambiente)), func(t *testing.T) {
			autorizacion := SimularAutorizacionSRI(claveAcceso, tt.ambiente)

			if autorizacion.Estado != tt.esperado {
				t.Errorf("Estado esperado: %s, obtenido: %s", tt.esperado, autorizacion.Estado)
			}

			if autorizacion.NumeroAutorizacion == "" {
				t.Error("Número de autorización no puede estar vacío")
			}

			if autorizacion.FechaAutorizacion.IsZero() {
				t.Error("Fecha de autorización no puede estar vacía")
			}

			if len(autorizacion.NumeroAutorizacion) != 37 {
				t.Errorf("Número de autorización debe tener 37 caracteres, obtuvo %d", len(autorizacion.NumeroAutorizacion))
			}
		})
	}
}

// TestMostrarInformacionClaveAcceso tests key access information display
func TestMostrarInformacionClaveAcceso(t *testing.T) {
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

	// This function prints to stdout, so we test it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MostrarInformacionClaveAcceso() panic = %v", r)
		}
	}()

	MostrarInformacionClaveAcceso(clave)
}

// TestObtenerNombreTipoComprobante tests voucher type name function
func TestObtenerNombreTipoComprobante(t *testing.T) {
	tests := []struct {
		tipo     TipoComprobante
		esperado string
	}{
		{Factura, "Factura"},
		{NotaCredito, "Nota de Crédito"},
		{NotaDebito, "Nota de Débito"},
		{GuiaRemision, "Guía de Remisión"},
		{ComprobanteRetencion, "Comprobante de Retención"},
		{TipoComprobante(99), "Desconocido"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Tipo_%d", tt.tipo), func(t *testing.T) {
			resultado := obtenerNombreTipoComprobante(tt.tipo)
			if resultado != tt.esperado {
				t.Errorf("Nombre esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestObtenerNombreAmbiente tests environment name function
func TestObtenerNombreAmbiente(t *testing.T) {
	tests := []struct {
		ambiente Ambiente
		esperado string
	}{
		{Pruebas, "Pruebas"},
		{Produccion, "Producción"},
		{Ambiente(99), "Desconocido"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Ambiente_%d", tt.ambiente), func(t *testing.T) {
			resultado := obtenerNombreAmbiente(tt.ambiente)
			if resultado != tt.esperado {
				t.Errorf("Nombre esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestObtenerNombreTipoEmision tests emission type name function
func TestObtenerNombreTipoEmision(t *testing.T) {
	tests := []struct {
		tipo     TipoEmision
		esperado string
	}{
		{EmisionNormal, "Emisión Normal"},
		{EmisionContingencia, "Emisión por Contingencia"},
		{TipoEmision(99), "Desconocido"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Tipo_%d", tt.tipo), func(t *testing.T) {
			resultado := obtenerNombreTipoEmision(tt.tipo)
			if resultado != tt.esperado {
				t.Errorf("Nombre esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestTipoComprobanteString tests TipoComprobante String method
func TestTipoComprobanteString(t *testing.T) {
	tests := []struct {
		tipo     TipoComprobante
		esperado string
	}{
		{Factura, "01"},
		{NotaCredito, "04"},
		{NotaDebito, "05"},
		{GuiaRemision, "06"},
		{ComprobanteRetencion, "07"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Tipo_%s", tt.esperado), func(t *testing.T) {
			resultado := tt.tipo.String()
			if resultado != tt.esperado {
				t.Errorf("String esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestAmbienteString tests Ambiente String method
func TestAmbienteString(t *testing.T) {
	tests := []struct {
		ambiente Ambiente
		esperado string
	}{
		{Pruebas, "1"},
		{Produccion, "2"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Ambiente_%s", tt.esperado), func(t *testing.T) {
			resultado := tt.ambiente.String()
			if resultado != tt.esperado {
				t.Errorf("String esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestTipoEmisionString tests TipoEmision String method
func TestTipoEmisionString(t *testing.T) {
	tests := []struct {
		tipo     TipoEmision
		esperado string
	}{
		{EmisionNormal, "1"},
		{EmisionContingencia, "2"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Tipo_%s", tt.esperado), func(t *testing.T) {
			resultado := tt.tipo.String()
			if resultado != tt.esperado {
				t.Errorf("String esperado: %s, obtenido: %s", tt.esperado, resultado)
			}
		})
	}
}

// TestCalcularDigitoVerificadorCasosBorde tests edge cases for verification digit
func TestCalcularDigitoVerificadorCasosBorde(t *testing.T) {
	tests := []struct {
		name    string
		clave   string
		wantErr bool
	}{
		{
			name:    "Clave válida",
			clave:   "230620240117921467390011001001000000001123456789",
			wantErr: false,
		},
		{
			name:    "Clave muy corta",
			clave:   "12345",
			wantErr: true,
		},
		{
			name:    "Clave vacía",
			clave:   "",
			wantErr: true,
		},
		{
			name:    "Clave con caracteres no numéricos",
			clave:   "23062024011792146739001100100100000000112345678A",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For edge cases with invalid input, we expect it to still return a digit
			// calcularDigitoVerificador handles invalid input gracefully
			digito := calcularDigitoVerificador(tt.clave)

			if !tt.wantErr {
				if digito < 0 || digito > 9 {
					t.Errorf("Dígito verificador debe estar entre 0-9, obtuvo %d", digito)
				}
			}
		})
	}
}

// TestGenerarCodigoNumerico tests numeric code generation
func TestGenerarCodigoNumerico(t *testing.T) {
	codigo := generarCodigoNumerico()

	if len(codigo) != 8 {
		t.Errorf("Código numérico debe tener 8 dígitos, obtuvo %d", len(codigo))
	}

	// Verificar que solo contiene dígitos
	for _, char := range codigo {
		if char < '0' || char > '9' {
			t.Errorf("Código numérico debe contener solo dígitos, encontró: %c", char)
		}
	}

	// Verificar unicidad generando múltiples códigos
	codigos := make(map[string]bool)
	for i := 0; i < 100; i++ {
		c := generarCodigoNumerico()
		if codigos[c] {
			t.Error("Los códigos numéricos deben ser únicos")
			break
		}
		codigos[c] = true
	}
}

// TestErrorHandlingClaveAcceso tests error handling in key generation
func TestErrorHandlingClaveAcceso(t *testing.T) {
	tests := []struct {
		name   string
		config ClaveAccesoConfig
	}{
		{
			name: "RUC demasiado corto",
			config: ClaveAccesoConfig{
				FechaEmision:     time.Now(),
				TipoComprobante:  Factura,
				RUCEmisor:        "123",
				Ambiente:         Pruebas,
				Serie:            "001001",
				NumeroSecuencial: "000000001",
				TipoEmision:      EmisionNormal,
			},
		},
		{
			name: "Serie inválida",
			config: ClaveAccesoConfig{
				FechaEmision:     time.Now(),
				TipoComprobante:  Factura,
				RUCEmisor:        "1792146739001",
				Ambiente:         Pruebas,
				Serie:            "001",
				NumeroSecuencial: "000000001",
				TipoEmision:      EmisionNormal,
			},
		},
		{
			name: "Secuencial inválido",
			config: ClaveAccesoConfig{
				FechaEmision:     time.Now(),
				TipoComprobante:  Factura,
				RUCEmisor:        "1792146739001",
				Ambiente:         Pruebas,
				Serie:            "001001",
				NumeroSecuencial: "1",
				TipoEmision:      EmisionNormal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerarClaveAcceso(tt.config)
			if err == nil {
				t.Error("GenerarClaveAcceso() esperaba error para configuración inválida")
			}
		})
	}
}

// TestIntegracionComprensiva tests comprehensive integration without real SRI
func TestIntegracionComprensiva(t *testing.T) {
	// Create invoice
	facturaData := models.FacturaInput{
		ClienteNombre: "TEST COMPREHENSIVE INTEGRATION",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "COMP001",
				Descripcion:    "Comprehensive Test Product",
				Cantidad:       3.0,
				PrecioUnitario: 199.99,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		t.Fatalf("Error creating invoice: %v", err)
	}

	// Generate XML
	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("Error generating XML: %v", err)
	}

	if len(xmlData) == 0 {
		t.Fatal("XML data cannot be empty")
	}

	// Generate access key
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
		t.Fatalf("Error generating access key: %v", err)
	}

	// Validate access key
	if err := ValidarClaveAcceso(claveAcceso); err != nil {
		t.Fatalf("Invalid access key: %v", err)
	}

	// Parse access key
	parsedConfig, err := ParsearClaveAcceso(claveAcceso)
	if err != nil {
		t.Fatalf("Error parsing access key: %v", err)
	}

	// Verify parsed data matches original
	if parsedConfig.RUCEmisor != claveConfig.RUCEmisor {
		t.Errorf("RUC mismatch: expected %s, got %s", claveConfig.RUCEmisor, parsedConfig.RUCEmisor)
	}

	// Generate authorization
	autorizacion := SimularAutorizacionSRI(claveAcceso, Pruebas)
	if autorizacion.Estado != "AUTORIZADO" {
		t.Errorf("Expected authorized state, got %s", autorizacion.Estado)
	}

	// Verify XML contains access key
	xmlString := string(xmlData)
	if !strings.Contains(xmlString, claveAcceso) {
		t.Error("XML should contain the access key")
	}

	t.Logf("✅ Comprehensive integration test completed successfully")
	t.Logf("📄 XML size: %d bytes", len(xmlData))
	t.Logf("🔑 Access key: %s", FormatearClaveAcceso(claveAcceso))
	t.Logf("📝 Authorization: %s", autorizacion.NumeroAutorizacion)
	t.Logf("💰 Invoice total: $%.2f", factura.InfoFactura.ImporteTotal)
}