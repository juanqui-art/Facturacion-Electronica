// Package factory contiene las funciones para crear facturas
package factory

import (
	"fmt"
	"log"
	"time"

	"go-facturacion-sri/config"
	"go-facturacion-sri/models"
	"go-facturacion-sri/validators"
)

// CrearFactura - Función factory que crea una factura completa con protección contra panics
// Recibe datos simples y devuelve una estructura completa lista para XML
// Ahora devuelve (Factura, error) - dos valores!
func CrearFactura(input models.FacturaInput) (factura models.Factura, err error) {
	// Protección contra panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CRITICAL] Panic recovered in CrearFactura: %v", r)
			factura = models.Factura{}
			err = fmt.Errorf("error crítico creando factura: %v", r)
		}
	}()
	// Primero validamos los datos de entrada
	if err := validators.ValidarFacturaInput(input); err != nil {
		return models.Factura{}, err // Devolvemos factura vacía y el error
	}

	// Validar que tengamos productos antes de procesar
	if len(input.Productos) == 0 {
		return models.Factura{}, fmt.Errorf("no se pueden procesar facturas sin productos")
	}

	// Calcular totales de TODOS los productos
	var subtotal float64 = 0
	var detalles []models.Detalle // Slice vacío para ir agregando productos

	// Procesar cada producto
	for i, producto := range input.Productos {
		// Validar producto antes de calcular
		if producto.Cantidad <= 0 || producto.PrecioUnitario <= 0 {
			return models.Factura{}, fmt.Errorf("producto %d tiene valores inválidos (cantidad: %.2f, precio: %.2f)", i+1, producto.Cantidad, producto.PrecioUnitario)
		}

		// Calcular subtotal de este producto con protección overflow
		subtotalProducto := producto.Cantidad * producto.PrecioUnitario
		
		// Verificar overflow
		if subtotalProducto > 99999999.99 {
			return models.Factura{}, fmt.Errorf("subtotal del producto %d excede límite máximo", i+1)
		}
		
		subtotal += subtotalProducto // Sumar al total general
		
		// Verificar overflow del subtotal total
		if subtotal > 99999999.99 {
			return models.Factura{}, fmt.Errorf("subtotal total excede límite máximo permitido")
		}

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

	// Validar subtotal antes de calcular IVA
	if subtotal <= 0 {
		return models.Factura{}, fmt.Errorf("subtotal inválido: %.2f", subtotal)
	}

	// Calcular IVA sobre el subtotal total
	iva := subtotal * 0.15 // 15% IVA Ecuador
	total := subtotal + iva
	
	// Validar que el total no exceda límites
	if total > 99999999.99 {
		return models.Factura{}, fmt.Errorf("total de factura excede límite máximo permitido: %.2f", total)
	}

	// Validar configuración antes de crear factura
	if config.Config.Empresa.RUC == "" {
		return models.Factura{}, fmt.Errorf("configuración incompleta: RUC de empresa no configurado")
	}
	if config.Config.Empresa.RazonSocial == "" {
		return models.Factura{}, fmt.Errorf("configuración incompleta: razón social no configurada")
	}

	// Crear la factura completa usando configuración externa
	facturaResult := models.Factura{
		InfoTributaria: models.InfoTributaria{
			Ambiente:        config.Config.Ambiente.Codigo,        // Desde configuración
			TipoEmision:     config.Config.Ambiente.TipoEmision,   // Desde configuración
			RazonSocial:     config.Config.Empresa.RazonSocial,    // Desde configuración
			RUC:             config.Config.Empresa.RUC,            // Desde configuración
			ClaveAcceso:     config.GenerarClaveAcceso(),          // Función del config
			CodDoc:          "01",                                 // 01=factura
			Establecimiento: config.Config.Empresa.Establecimiento, // Desde configuración
			PuntoEmision:    config.Config.Empresa.PuntoEmision,   // Desde configuración
			Secuencial:      config.ObtenerSecuencialSiguiente(),  // Función del config
		},
		InfoFactura: models.InfoFactura{
			FechaEmision:                time.Now().Format("02/01/2006"), // DD/MM/YYYY
			DirEstablecimiento:          config.Config.Empresa.Direccion,  // Desde configuración
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

	return facturaResult, nil // nil significa "no hay error"
}

// La función generarClaveAcceso() se movió al package config
