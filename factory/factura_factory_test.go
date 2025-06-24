package factory

import (
	"go-facturacion-sri/config"
	"go-facturacion-sri/models"
	"math"
	"testing"
	"time"
)

// almostEqual verifica si dos números float64 son casi iguales (considera precisión)
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

// setUp inicializa la configuración para los tests
func setUp() {
	config.CargarConfiguracionPorDefecto()
}

// TestCrearFactura prueba la creación de facturas
func TestCrearFactura(t *testing.T) {
	setUp()

	tests := []struct {
		name      string
		input     models.FacturaInput
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, models.Factura)
	}{
		// Casos válidos
		{
			name: "factura válida con un producto",
			input: models.FacturaInput{
				ClienteNombre: "Juan Carlos Pérez",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "LAPTOP001",
						Descripcion:    "Laptop Dell Inspiron 15",
						Cantidad:       1.0,
						PrecioUnitario: 450.00,
					},
				},
			},
			wantErr: false,
			checkFunc: func(t *testing.T, factura models.Factura) {
				// Verificar cálculos
				expectedSubtotal := 450.00
				expectedIVA := expectedSubtotal * 0.15          // 67.50
				expectedTotal := expectedSubtotal + expectedIVA // 517.50

				if factura.InfoFactura.TotalSinImpuestos != expectedSubtotal {
					t.Errorf("TotalSinImpuestos = %v, quería %v", factura.InfoFactura.TotalSinImpuestos, expectedSubtotal)
				}
				if factura.InfoFactura.ImporteTotal != expectedTotal {
					t.Errorf("ImporteTotal = %v, quería %v", factura.InfoFactura.ImporteTotal, expectedTotal)
				}

				// Verificar estructura
				if len(factura.Detalles) != 1 {
					t.Errorf("Número de detalles = %v, quería 1", len(factura.Detalles))
				}

				// Verificar datos del cliente
				if factura.InfoFactura.RazonSocialComprador != "Juan Carlos Pérez" {
					t.Errorf("RazonSocialComprador = %v, quería 'Juan Carlos Pérez'", factura.InfoFactura.RazonSocialComprador)
				}
				if factura.InfoFactura.IdentificacionComprador != "1713175071" {
					t.Errorf("IdentificacionComprador = %v, quería '1713175071'", factura.InfoFactura.IdentificacionComprador)
				}

				// Verificar fecha (debe ser hoy)
				expectedDate := time.Now().Format("02/01/2006")
				if factura.InfoFactura.FechaEmision != expectedDate {
					t.Errorf("FechaEmision = %v, quería %v", factura.InfoFactura.FechaEmision, expectedDate)
				}

				// Verificar información tributaria
				if factura.InfoTributaria.CodDoc != "01" {
					t.Errorf("CodDoc = %v, quería '01'", factura.InfoTributaria.CodDoc)
				}
				if factura.InfoTributaria.TipoEmision != "1" {
					t.Errorf("TipoEmision = %v, quería '1'", factura.InfoTributaria.TipoEmision)
				}
			},
		},
		{
			name: "factura válida con múltiples productos",
			input: models.FacturaInput{
				ClienteNombre: "María González",
				ClienteCedula: "0926687856",
				Productos: []models.ProductoInput{
					{
						Codigo:         "LAPTOP001",
						Descripcion:    "Laptop Dell Inspiron 15",
						Cantidad:       2.0,
						PrecioUnitario: 450.00,
					},
					{
						Codigo:         "MOUSE001",
						Descripcion:    "Mouse Inalámbrico",
						Cantidad:       3.0,
						PrecioUnitario: 25.00,
					},
					{
						Codigo:         "TECLADO001",
						Descripcion:    "Teclado Mecánico",
						Cantidad:       1.0,
						PrecioUnitario: 85.00,
					},
				},
			},
			wantErr: false,
			checkFunc: func(t *testing.T, factura models.Factura) {
				// Verificar cálculos
				// Laptop: 2 * 450 = 900
				// Mouse: 3 * 25 = 75
				// Teclado: 1 * 85 = 85
				// Subtotal: 900 + 75 + 85 = 1060
				expectedSubtotal := 1060.00
				expectedIVA := expectedSubtotal * 0.15          // 159.00
				expectedTotal := expectedSubtotal + expectedIVA // 1219.00

				if factura.InfoFactura.TotalSinImpuestos != expectedSubtotal {
					t.Errorf("TotalSinImpuestos = %v, quería %v", factura.InfoFactura.TotalSinImpuestos, expectedSubtotal)
				}
				if factura.InfoFactura.ImporteTotal != expectedTotal {
					t.Errorf("ImporteTotal = %v, quería %v", factura.InfoFactura.ImporteTotal, expectedTotal)
				}

				// Verificar número de detalles
				if len(factura.Detalles) != 3 {
					t.Errorf("Número de detalles = %v, quería 3", len(factura.Detalles))
				}

				// Verificar detalles individuales
				detalles := factura.Detalles
				if detalles[0].PrecioTotalSinImpuesto != 900.00 {
					t.Errorf("Precio total primer producto = %v, quería 900.00", detalles[0].PrecioTotalSinImpuesto)
				}
				if detalles[1].PrecioTotalSinImpuesto != 75.00 {
					t.Errorf("Precio total segundo producto = %v, quería 75.00", detalles[1].PrecioTotalSinImpuesto)
				}
				if detalles[2].PrecioTotalSinImpuesto != 85.00 {
					t.Errorf("Precio total tercer producto = %v, quería 85.00", detalles[2].PrecioTotalSinImpuesto)
				}
			},
		},
		{
			name: "factura con cantidades decimales",
			input: models.FacturaInput{
				ClienteNombre: "Empresa ABC",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "CABLE001",
						Descripcion:    "Cable de red por metro",
						Cantidad:       2.5,
						PrecioUnitario: 12.00,
					},
				},
			},
			wantErr: false,
			checkFunc: func(t *testing.T, factura models.Factura) {
				// Verificar cálculos con decimales
				// 2.5 * 12.00 = 30.00
				expectedSubtotal := 30.00
				expectedIVA := expectedSubtotal * 0.15          // 4.50
				expectedTotal := expectedSubtotal + expectedIVA // 34.50

				if factura.InfoFactura.TotalSinImpuestos != expectedSubtotal {
					t.Errorf("TotalSinImpuestos = %v, quería %v", factura.InfoFactura.TotalSinImpuestos, expectedSubtotal)
				}
				if factura.InfoFactura.ImporteTotal != expectedTotal {
					t.Errorf("ImporteTotal = %v, quería %v", factura.InfoFactura.ImporteTotal, expectedTotal)
				}
			},
		},
		// Casos inválidos (que deberían retornar error)
		{
			name: "cliente sin nombre",
			input: models.FacturaInput{
				ClienteNombre: "",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto test",
						Cantidad:       1.0,
						PrecioUnitario: 100.00,
					},
				},
			},
			wantErr: true,
			errMsg:  "el nombre del cliente no puede estar vacío",
		},
		{
			name: "cédula inválida",
			input: models.FacturaInput{
				ClienteNombre: "Cliente Test",
				ClienteCedula: "123", // cédula muy corta
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto test",
						Cantidad:       1.0,
						PrecioUnitario: 100.00,
					},
				},
			},
			wantErr: true,
			errMsg:  "cédula inválida: la cédula debe tener exactamente 10 dígitos",
		},
		{
			name: "sin productos",
			input: models.FacturaInput{
				ClienteNombre: "Cliente Test",
				ClienteCedula: "1713175071",
				Productos:     []models.ProductoInput{}, // slice vacío
			},
			wantErr: true,
			errMsg:  "debe incluir al menos un producto",
		},
		{
			name: "producto con precio cero",
			input: models.FacturaInput{
				ClienteNombre: "Cliente Test",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto con precio cero",
						Cantidad:       1.0,
						PrecioUnitario: 0.0, // precio inválido
					},
				},
			},
			wantErr: true,
			errMsg:  "producto 1 inválido: el precio unitario debe ser mayor a cero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factura, err := CrearFactura(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CrearFactura() error = nil, quería error")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("CrearFactura() error = %v, quería %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("CrearFactura() error = %v, no quería error", err)
					return
				}
				// Ejecutar verificaciones adicionales si no hay error
				if tt.checkFunc != nil {
					tt.checkFunc(t, factura)
				}
			}
		})
	}
}

// TestCrearFactura_ConfiguracionTributaria verifica que se use la configuración correcta
func TestCrearFactura_ConfiguracionTributaria(t *testing.T) {
	setUp()

	input := models.FacturaInput{
		ClienteNombre: "Test Cliente",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "TEST001",
				Descripcion:    "Producto test",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
		},
	}

	factura, err := CrearFactura(input)
	if err != nil {
		t.Fatalf("CrearFactura() error = %v, no quería error", err)
	}

	// Verificar que se usen los valores de configuración
	if factura.InfoTributaria.RazonSocial == "" {
		t.Error("RazonSocial no debe estar vacía")
	}
	if factura.InfoTributaria.RUC == "" {
		t.Error("RUC no debe estar vacío")
	}
	if factura.InfoTributaria.Ambiente == "" {
		t.Error("Ambiente no debe estar vacío")
	}
	if factura.InfoTributaria.ClaveAcceso == "" {
		t.Error("ClaveAcceso no debe estar vacía")
	}
	if factura.InfoTributaria.Secuencial == "" {
		t.Error("Secuencial no debe estar vacío")
	}

	// Verificar valores específicos de la configuración por defecto
	if factura.InfoFactura.TipoIdentificacionComprador != "05" {
		t.Errorf("TipoIdentificacionComprador = %v, quería '05'", factura.InfoFactura.TipoIdentificacionComprador)
	}
	if factura.InfoFactura.Moneda != "DOLAR" {
		t.Errorf("Moneda = %v, quería 'DOLAR'", factura.InfoFactura.Moneda)
	}
	if factura.InfoFactura.TotalDescuento != 0.00 {
		t.Errorf("TotalDescuento = %v, quería 0.00", factura.InfoFactura.TotalDescuento)
	}
}

// TestCrearFactura_CalculosIVA verifica específicamente los cálculos de IVA
func TestCrearFactura_CalculosIVA(t *testing.T) {
	setUp()

	testCases := []struct {
		name             string
		cantidad         float64
		precio           float64
		expectedSubtotal float64
		expectedTotal    float64
	}{
		{"precio entero", 1.0, 100.00, 100.00, 115.00},
		{"precio con decimales", 1.0, 99.99, 99.99, 114.9885}, // 99.99 * 1.15
		{"cantidad múltiple", 3.0, 50.00, 150.00, 172.50},
		{"cantidad decimal", 2.5, 40.00, 100.00, 115.00},
		{"precio alto", 1.0, 1000.00, 1000.00, 1150.00},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := models.FacturaInput{
				ClienteNombre: "Test Cliente",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "TEST001",
						Descripcion:    "Producto test",
						Cantidad:       tc.cantidad,
						PrecioUnitario: tc.precio,
					},
				},
			}

			factura, err := CrearFactura(input)
			if err != nil {
				t.Fatalf("CrearFactura() error = %v, no quería error", err)
			}

			if !almostEqual(factura.InfoFactura.TotalSinImpuestos, tc.expectedSubtotal) {
				t.Errorf("TotalSinImpuestos = %v, quería %v", factura.InfoFactura.TotalSinImpuestos, tc.expectedSubtotal)
			}
			if !almostEqual(factura.InfoFactura.ImporteTotal, tc.expectedTotal) {
				t.Errorf("ImporteTotal = %v, quería %v", factura.InfoFactura.ImporteTotal, tc.expectedTotal)
			}
		})
	}
}

// Benchmark para CrearFactura con un producto
func BenchmarkCrearFactura_UnProducto(b *testing.B) {
	setUp()

	input := models.FacturaInput{
		ClienteNombre: "Cliente Benchmark",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "BENCH001",
				Descripcion:    "Producto benchmark",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := CrearFactura(input)
		if err != nil {
			b.Fatalf("CrearFactura() error = %v", err)
		}
	}
}

// Benchmark para CrearFactura con múltiples productos
func BenchmarkCrearFactura_MultipleProductos(b *testing.B) {
	setUp()

	productos := make([]models.ProductoInput, 10)
	for i := 0; i < 10; i++ {
		productos[i] = models.ProductoInput{
			Codigo:         "PROD" + string(rune(i+1)),
			Descripcion:    "Producto " + string(rune(i+1)),
			Cantidad:       1.0,
			PrecioUnitario: 100.00,
		}
	}

	input := models.FacturaInput{
		ClienteNombre: "Cliente Benchmark",
		ClienteCedula: "1713175071",
		Productos:     productos,
	}

	for i := 0; i < b.N; i++ {
		_, err := CrearFactura(input)
		if err != nil {
			b.Fatalf("CrearFactura() error = %v", err)
		}
	}
}
