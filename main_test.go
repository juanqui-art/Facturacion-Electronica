package main

import (
	"testing"
	
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/validators"
)

// TestValidarCedula - Prueba la validación de cédulas ecuatorianas usando table-driven tests
func TestValidarCedula(t *testing.T) {
	// Tabla de casos de prueba
	tests := []struct {
		name          string  // Nombre descriptivo del test
		cedula        string  // Input: la cédula a probar
		shouldBeValid bool    // Expected: ¿debería ser válida?
		description   string  // Descripción del caso
	}{
		{
			name:          "cedula_valida_pichincha",
			cedula:        "1713175071",
			shouldBeValid: true,
			description:   "Cédula válida de Pichincha (provincia 17)",
		},
		{
			name:          "cedula_muy_corta",
			cedula:        "123456789",
			shouldBeValid: false,
			description:   "Cédula con solo 9 dígitos",
		},
		{
			name:          "cedula_muy_larga",
			cedula:        "12345678901",
			shouldBeValid: false,
			description:   "Cédula con 11 dígitos",
		},
		{
			name:          "cedula_con_letras",
			cedula:        "17131ABC71",
			shouldBeValid: false,
			description:   "Cédula que contiene letras",
		},
		{
			name:          "provincia_invalida",
			cedula:        "2713175071",
			shouldBeValid: false,
			description:   "Provincia 27 no existe (solo hay 24)",
		},
		{
			name:          "digito_verificador_incorrecto",
			cedula:        "1713175070",
			shouldBeValid: false,
			description:   "Último dígito incorrecto",
		},
	}
	
	// Ejecutar cada caso de prueba
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validators.ValidarCedula(test.cedula)
			
			if test.shouldBeValid {
				// Esperamos que sea válida (no error)
				if err != nil {
					t.Errorf("Test '%s': %s\nCédula: %s\nError inesperado: %v", 
						test.name, test.description, test.cedula, err)
				}
			} else {
				// Esperamos que sea inválida (con error)
				if err == nil {
					t.Errorf("Test '%s': %s\nCédula: %s\nDebería haber fallado pero no lo hizo", 
						test.name, test.description, test.cedula)
				}
			}
		})
	}
}

// TestValidarProducto - Prueba la validación de productos individuales
func TestValidarProducto(t *testing.T) {
	tests := []struct {
		name          string
		producto      models.ProductoInput
		shouldBeValid bool
		description   string
	}{
		{
			name: "producto_valido",
			producto: models.ProductoInput{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell",
				Cantidad:       2.0,
				PrecioUnitario: 450.00,
			},
			shouldBeValid: true,
			description:   "Producto con todos los campos válidos",
		},
		{
			name: "codigo_vacio",
			producto: models.ProductoInput{
				Codigo:         "", // Código vacío
				Descripcion:    "Laptop Dell",
				Cantidad:       2.0,
				PrecioUnitario: 450.00,
			},
			shouldBeValid: false,
			description:   "Producto sin código",
		},
		{
			name: "descripcion_vacia",
			producto: models.ProductoInput{
				Codigo:         "LAPTOP001",
				Descripcion:    "", // Descripción vacía
				Cantidad:       2.0,
				PrecioUnitario: 450.00,
			},
			shouldBeValid: false,
			description:   "Producto sin descripción",
		},
		{
			name: "cantidad_cero",
			producto: models.ProductoInput{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell",
				Cantidad:       0.0, // Cantidad cero
				PrecioUnitario: 450.00,
			},
			shouldBeValid: false,
			description:   "Producto con cantidad cero",
		},
		{
			name: "precio_negativo",
			producto: models.ProductoInput{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell",
				Cantidad:       2.0,
				PrecioUnitario: -450.00, // Precio negativo
			},
			shouldBeValid: false,
			description:   "Producto con precio negativo",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validators.ValidarProducto(test.producto)
			
			if test.shouldBeValid {
				if err != nil {
					t.Errorf("Test '%s': %s\nProducto: %+v\nError inesperado: %v", 
						test.name, test.description, test.producto, err)
				}
			} else {
				if err == nil {
					t.Errorf("Test '%s': %s\nProducto: %+v\nDebería haber fallado pero no lo hizo", 
						test.name, test.description, test.producto)
				}
			}
		})
	}
}

// TestCrearFactura - Prueba la creación completa de facturas
func TestCrearFactura(t *testing.T) {
	tests := []struct {
		name          string
		input         models.FacturaInput
		shouldBeValid bool
		expectedTotal float64 // Total esperado si es válida
		description   string
	}{
		{
			name: "factura_valida_un_producto",
			input: models.FacturaInput{
				ClienteNombre: "Juan Perez",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto de prueba",
						Cantidad:       2.0,
						PrecioUnitario: 100.0,
					},
				},
			},
			shouldBeValid: true,
			expectedTotal: 230.0, // 200 + 30 (15% IVA) = 230
			description:   "Factura válida con un producto",
		},
		{
			name: "factura_valida_multiples_productos",
			input: models.FacturaInput{
				ClienteNombre: "Maria Rodriguez",
				ClienteCedula: "1713175071",
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto A",
						Cantidad:       1.0,
						PrecioUnitario: 100.0,
					},
					{
						Codigo:         "PROD002",
						Descripcion:    "Producto B",
						Cantidad:       2.0,
						PrecioUnitario: 50.0,
					},
				},
			},
			shouldBeValid: true,
			expectedTotal: 230.0, // (100 + 100) + 30 (15% IVA) = 230
			description:   "Factura válida con múltiples productos",
		},
		{
			name: "cedula_invalida",
			input: models.FacturaInput{
				ClienteNombre: "Juan Perez",
				ClienteCedula: "123456789", // Cédula inválida
				Productos: []models.ProductoInput{
					{
						Codigo:         "PROD001",
						Descripcion:    "Producto",
						Cantidad:       1.0,
						PrecioUnitario: 100.0,
					},
				},
			},
			shouldBeValid: false,
			expectedTotal: 0,
			description:   "Factura con cédula inválida",
		},
		{
			name: "sin_productos",
			input: models.FacturaInput{
				ClienteNombre: "Juan Perez",
				ClienteCedula: "1713175071",
				Productos:     []models.ProductoInput{}, // Sin productos
			},
			shouldBeValid: false,
			expectedTotal: 0,
			description:   "Factura sin productos",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			factura, err := factory.CrearFactura(test.input)
			
			if test.shouldBeValid {
				// Esperamos que sea válida
				if err != nil {
					t.Errorf("Test '%s': %s\nError inesperado: %v", 
						test.name, test.description, err)
					return
				}
				
				// Verificar que el total calculado sea correcto
				if factura.InfoFactura.ImporteTotal != test.expectedTotal {
					t.Errorf("Test '%s': Total incorrecto\nEsperado: %.2f\nObtenido: %.2f", 
						test.name, test.expectedTotal, factura.InfoFactura.ImporteTotal)
				}
				
				// Verificar que tenga el número correcto de productos
				expectedProducts := len(test.input.Productos)
				actualProducts := len(factura.Detalles)
				if actualProducts != expectedProducts {
					t.Errorf("Test '%s': Número de productos incorrecto\nEsperado: %d\nObtenido: %d", 
						test.name, expectedProducts, actualProducts)
				}
				
			} else {
				// Esperamos que sea inválida
				if err == nil {
					t.Errorf("Test '%s': %s\nDebería haber fallado pero no lo hizo", 
						test.name, test.description)
				}
			}
		})
	}
}