package validators

import (
	"go-facturacion-sri/models"
	"strings"
	"testing"
)

// TestValidarRUC prueba la validación de RUCs ecuatorianos
func TestValidarRUC(t *testing.T) {
	tests := []struct {
		name    string
		ruc     string
		wantErr bool
		errMsg  string
	}{
		// RUCs válidos - Persona Natural
		{
			name:    "RUC persona natural válido - Pichincha",
			ruc:     "1713175071001",
			wantErr: false,
		},
		{
			name:    "RUC persona natural válido - Guayas",
			ruc:     "0926687856001",
			wantErr: false,
		},
		// Note: Using basic validation for now - algorithms need specific valid RUCs for Ecuador
		// Casos inválidos - longitud
		{
			name:    "RUC muy corto",
			ruc:     "123456789012",
			wantErr: true,
			errMsg:  "el RUC debe tener exactamente 13 dígitos",
		},
		{
			name:    "RUC muy largo",
			ruc:     "12345678901234",
			wantErr: true,
			errMsg:  "el RUC debe tener exactamente 13 dígitos",
		},
		{
			name:    "RUC vacío",
			ruc:     "",
			wantErr: true,
			errMsg:  "el RUC debe tener exactamente 13 dígitos",
		},
		// Casos inválidos - caracteres
		{
			name:    "RUC con letras",
			ruc:     "171317507100A",
			wantErr: true,
			errMsg:  "el RUC solo puede contener números",
		},
		// Casos inválidos - provincia
		{
			name:    "provincia inválida - 00",
			ruc:     "0013175071001",
			wantErr: true,
			errMsg:  "los dos primeros dígitos del RUC deben estar entre 01 y 24",
		},
		{
			name:    "provincia inválida - 25",
			ruc:     "2513175071001",
			wantErr: true,
			errMsg:  "los dos primeros dígitos del RUC deben estar entre 01 y 24",
		},
		// Casos inválidos - tercer dígito
		{
			name:    "tercer dígito inválido - 7",
			ruc:     "1773175071001",
			wantErr: true,
			errMsg:  "el tercer dígito del RUC debe ser menor a 6, o igual a 6 (sector público) o 9 (empresa privada)",
		},
		// Casos inválidos - terminación incorrecta
		{
			name:    "RUC persona natural con terminación incorrecta",
			ruc:     "1713175071002",
			wantErr: true,
			errMsg:  "RUC persona natural debe terminar en 001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarRUC(tt.ruc)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidarRUC() error = nil, quería error")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidarRUC() error = %v, quería %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidarRUC() error = %v, no quería error", err)
				}
			}
		})
	}
}

// TestValidarCedula prueba la validación de cédulas ecuatorianas
func TestValidarCedula(t *testing.T) {
	tests := []struct {
		name    string
		cedula  string
		wantErr bool
		errMsg  string
	}{
		// Casos válidos (usando cédulas que ya sabemos que funcionan en el código)
		{
			name:    "cédula válida - Pichincha",
			cedula:  "1713175071",
			wantErr: false,
		},
		{
			name:    "cédula válida - Guayas",
			cedula:  "0926687856",
			wantErr: false,
		},
		// Casos inválidos - longitud
		{
			name:    "cédula muy corta",
			cedula:  "123456789",
			wantErr: true,
			errMsg:  "la cédula debe tener exactamente 10 dígitos",
		},
		{
			name:    "cédula muy larga",
			cedula:  "12345678901",
			wantErr: true,
			errMsg:  "la cédula debe tener exactamente 10 dígitos",
		},
		{
			name:    "cédula vacía",
			cedula:  "",
			wantErr: true,
			errMsg:  "la cédula debe tener exactamente 10 dígitos",
		},
		// Casos inválidos - caracteres
		{
			name:    "cédula con letras",
			cedula:  "171317507A",
			wantErr: true,
			errMsg:  "la cédula solo puede contener números",
		},
		{
			name:    "cédula con espacios",
			cedula:  "1713175 71",
			wantErr: true,
			errMsg:  "la cédula solo puede contener números",
		},
		{
			name:    "cédula con guiones",
			cedula:  "1713-17507",
			wantErr: true,
			errMsg:  "la cédula solo puede contener números",
		},
		// Casos inválidos - provincia
		{
			name:    "provincia inválida - 00",
			cedula:  "0013175071",
			wantErr: true,
			errMsg:  "los dos primeros dígitos de la cédula deben estar entre 01 y 24",
		},
		{
			name:    "provincia inválida - 25",
			cedula:  "2513175071",
			wantErr: true,
			errMsg:  "los dos primeros dígitos de la cédula deben estar entre 01 y 24",
		},
		{
			name:    "provincia inválida - 99",
			cedula:  "9913175071",
			wantErr: true,
			errMsg:  "los dos primeros dígitos de la cédula deben estar entre 01 y 24",
		},
		// Casos inválidos - dígito verificador
		{
			name:    "dígito verificador incorrecto",
			cedula:  "1713175072", // último dígito cambiado
			wantErr: true,
			errMsg:  "el dígito verificador de la cédula no es válido",
		},
		{
			name:    "otra cédula con dígito verificador incorrecto",
			cedula:  "0926687855", // último dígito cambiado
			wantErr: true,
			errMsg:  "el dígito verificador de la cédula no es válido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarCedula(tt.cedula)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidarCedula() error = nil, quería error")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidarCedula() error = %v, quería %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidarCedula() error = %v, no quería error", err)
				}
			}
		})
	}
}

// TestValidarProducto prueba la validación de productos individuales
func TestValidarProducto(t *testing.T) {
	tests := []struct {
		name     string
		producto models.ProductoInput
		wantErr  bool
		errMsg   string
	}{
		// Casos válidos
		{
			name: "producto válido básico",
			producto: models.ProductoInput{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell Inspiron 15",
				Cantidad:       1.0,
				PrecioUnitario: 450.00,
			},
			wantErr: false,
		},
		{
			name: "producto válido con cantidad decimal",
			producto: models.ProductoInput{
				Codigo:         "MOUSE001",
				Descripcion:    "Mouse Inalámbrico",
				Cantidad:       2.5,
				PrecioUnitario: 25.99,
			},
			wantErr: false,
		},
		{
			name: "producto válido con precio alto",
			producto: models.ProductoInput{
				Codigo:         "SERVER001",
				Descripcion:    "Servidor Dell PowerEdge",
				Cantidad:       1.0,
				PrecioUnitario: 5000.00,
			},
			wantErr: false,
		},
		// Casos inválidos - código
		{
			name: "código vacío",
			producto: models.ProductoInput{
				Codigo:         "",
				Descripcion:    "Producto sin código",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
			wantErr: true,
			errMsg:  "el código del producto no puede estar vacío",
		},
		// Casos inválidos - descripción
		{
			name: "descripción vacía",
			producto: models.ProductoInput{
				Codigo:         "PROD001",
				Descripcion:    "",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
			wantErr: true,
			errMsg:  "la descripción del producto no puede estar vacía",
		},
		// Casos inválidos - cantidad
		{
			name: "cantidad cero",
			producto: models.ProductoInput{
				Codigo:         "PROD001",
				Descripcion:    "Producto con cantidad cero",
				Cantidad:       0.0,
				PrecioUnitario: 100.00,
			},
			wantErr: true,
			errMsg:  "la cantidad debe ser mayor a cero",
		},
		{
			name: "cantidad negativa",
			producto: models.ProductoInput{
				Codigo:         "PROD001",
				Descripcion:    "Producto con cantidad negativa",
				Cantidad:       -1.0,
				PrecioUnitario: 100.00,
			},
			wantErr: true,
			errMsg:  "la cantidad debe ser mayor a cero",
		},
		// Casos inválidos - precio
		{
			name: "precio cero",
			producto: models.ProductoInput{
				Codigo:         "PROD001",
				Descripcion:    "Producto con precio cero",
				Cantidad:       1.0,
				PrecioUnitario: 0.0,
			},
			wantErr: true,
			errMsg:  "el precio unitario debe ser mayor a cero",
		},
		{
			name: "precio negativo",
			producto: models.ProductoInput{
				Codigo:         "PROD001",
				Descripcion:    "Producto con precio negativo",
				Cantidad:       1.0,
				PrecioUnitario: -100.00,
			},
			wantErr: true,
			errMsg:  "el precio unitario debe ser mayor a cero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarProducto(tt.producto)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidarProducto() error = nil, quería error")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidarProducto() error = %v, quería %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidarProducto() error = %v, no quería error", err)
				}
			}
		})
	}
}

// TestValidarFacturaInput prueba la validación completa de entrada de facturas
func TestValidarFacturaInput(t *testing.T) {
	tests := []struct {
		name    string
		input   models.FacturaInput
		wantErr bool
		errMsg  string
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
				},
			},
			wantErr: false,
		},
		// Casos inválidos - cliente
		{
			name: "nombre de cliente vacío",
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
			name: "cédula de cliente inválida",
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
		// Casos inválidos - productos
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
			name: "producto inválido en la lista",
			input: models.FacturaInput{
				ClienteNombre: "Cliente Test",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto válido",
						Cantidad:       1.0,
						PrecioUnitario: 100.00,
					},
					{
						Codigo:         "", // código vacío - inválido
						Descripcion:    "Producto inválido",
						Cantidad:       1.0,
						PrecioUnitario: 100.00,
					},
				},
			},
			wantErr: true,
			errMsg:  "producto 2 inválido: el código del producto no puede estar vacío",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarFacturaInput(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidarFacturaInput() error = nil, quería error")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidarFacturaInput() error = %v, quería %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidarFacturaInput() error = %v, no quería error", err)
				}
			}
		})
	}
}

// Benchmark para ValidarCedula
func BenchmarkValidarCedula(b *testing.B) {
	cedula := "1713175071"
	for i := 0; i < b.N; i++ {
		ValidarCedula(cedula)
	}
}

// Benchmark para ValidarProducto
func BenchmarkValidarProducto(b *testing.B) {
	producto := models.ProductoInput{
		Codigo:         "LAPTOP001",
		Descripcion:    "Laptop Dell Inspiron 15",
		Cantidad:       1.0,
		PrecioUnitario: 450.00,
	}
	for i := 0; i < b.N; i++ {
		ValidarProducto(producto)
	}
}

// Benchmark para ValidarFacturaInput
func BenchmarkValidarFacturaInput(b *testing.B) {
	input := models.FacturaInput{
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
	}
	for i := 0; i < b.N; i++ {
		ValidarFacturaInput(input)
	}
}

// TestValidarRUC_MaliciousInputs tests for malicious RUC inputs
func TestValidarRUC_MaliciousInputs(t *testing.T) {
	tests := []struct {
		name    string
		ruc     string
		wantErr bool
	}{
		{"XSS attempt", "<script>alert('xss')</script>", true},
		{"SQL injection", "'; DROP TABLE usuarios; --", true},
		{"Extremely long RUC", "1234567890123456789012345678901234567890", true},
		{"Null byte injection", "1713175071001\x00", true},
		{"Empty string", "", true},
		{"Only whitespaces", "             ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarRUC(tt.ruc)
			if tt.wantErr && err == nil {
				t.Errorf("ValidarRUC() expected error for malicious input %q", tt.ruc)
			}
		})
	}
}

// TestSanitizarTexto_Security tests text sanitization
func TestSanitizarTexto_Security(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"XSS script tags", "<script>alert('xss')</script>"},
		{"HTML img with onerror", "<img src=x onerror=alert(1)>"},
		{"Control characters", "Normal\x00\x01\x02Text"},
		{"Very long string", string(make([]byte, 2000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizarTexto(tt.input)
			
			// Verify no script tags remain
			if strings.Contains(result, "<script>") {
				t.Errorf("SanitizarTexto() failed to sanitize script tags")
			}
			
			// Verify length limit
			if len(result) > 1000 {
				t.Errorf("SanitizarTexto() returned string longer than 1000 chars: %d", len(result))
			}
		})
	}
}
