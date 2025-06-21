// Package factory contiene las funciones para crear facturas
package factory

import (
	"time"
	
	"go-facturacion-sri/models"
	"go-facturacion-sri/validators"
)

// CrearFactura - Función factory que crea una factura completa
// Recibe datos simples y devuelve una estructura completa lista para XML
// Ahora devuelve (Factura, error) - dos valores!
func CrearFactura(input models.FacturaInput) (models.Factura, error) {
	// Primero validamos los datos de entrada
	if err := validators.ValidarFacturaInput(input); err != nil {
		return models.Factura{}, err // Devolvemos factura vacía y el error
	}
	
	// Calcular totales de TODOS los productos
	var subtotal float64 = 0
	var detalles []models.Detalle // Slice vacío para ir agregando productos
	
	// Procesar cada producto
	for _, producto := range input.Productos {
		// Calcular subtotal de este producto
		subtotalProducto := producto.Cantidad * producto.PrecioUnitario
		subtotal += subtotalProducto // Sumar al total general
		
		// Crear detalle para este producto
		detalle := models.Detalle{
			CodigoPrincipal:        producto.Codigo,
			Descripcion:            producto.Descripcion,
			Cantidad:               producto.Cantidad,
			PrecioUnitario:         producto.PrecioUnitario,
			Descuento:              0.00,
			PrecioTotalSinImpuesto: subtotalProducto,
		}
		
		// Agregar al slice de detalles
		detalles = append(detalles, detalle)
	}
	
	// Calcular IVA sobre el subtotal total
	iva := subtotal * 0.15  // 15% IVA Ecuador
	total := subtotal + iva
	
	// Crear la factura completa con valores por defecto
	factura := models.Factura{
		InfoTributaria: models.InfoTributaria{
			Ambiente:        "1", // 1=pruebas, 2=producción
			TipoEmision:     "1", // 1=normal
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			ClaveAcceso:     generarClaveAcceso(),
			CodDoc:          "01", // 01=factura
			Establecimiento: "001",
			PuntoEmision:    "001",
			Secuencial:      "000000001",
		},
		InfoFactura: models.InfoFactura{
			FechaEmision:                time.Now().Format("02/01/2006"), // DD/MM/YYYY
			DirEstablecimiento:          "Av. Amazonas y Naciones Unidas",
			TipoIdentificacionComprador: "05", // 05=cédula
			IdentificacionComprador:     input.ClienteCedula,
			RazonSocialComprador:        input.ClienteNombre,
			TotalSinImpuestos:           subtotal,
			TotalDescuento:              0.00,
			ImporteTotal:                total,
			Moneda:                      "DOLAR",
		},
		Detalles: detalles, // Usar el slice que construimos en el loop
	}
	
	return factura, nil // nil significa "no hay error"
}

// generarClaveAcceso - función privada (minúscula)
// Por ahora generamos una clave fake - en semana 4 implementaremos el algoritmo real del SRI
func generarClaveAcceso() string {
	return "2025062001123456789000110010010000000011234567890"
}