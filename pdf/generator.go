// Package pdf implementa la generación de PDFs para facturas
package pdf

import (
	"bytes"
	"fmt"
	"go-facturacion-sri/database"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// FacturaPDFGenerator genera PDFs de facturas
type FacturaPDFGenerator struct {
	db *database.Database
}

// NewFacturaPDFGenerator crea un nuevo generador de PDFs
func NewFacturaPDFGenerator(db *database.Database) *FacturaPDFGenerator {
	return &FacturaPDFGenerator{db: db}
}

// GenerarFacturaPDF genera un PDF para una factura específica
func (g *FacturaPDFGenerator) GenerarFacturaPDF(facturaID int) ([]byte, error) {
	// Obtener factura de la base de datos
	factura, err := g.db.ObtenerFacturaPorID(facturaID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo factura: %v", err)
	}

	// Obtener productos de la factura
	productos, err := g.db.ObtenerProductosPorFactura(facturaID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo productos: %v", err)
	}

	// Crear PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configurar fuentes
	pdf.SetFont("Arial", "B", 16)

	// Header - Título principal
	pdf.CellFormat(190, 10, "FACTURA ELECTRÓNICA", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Información de la empresa (hardcoded por ahora, debería venir de configuración)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(95, 8, "EMPRESA DEMO S.A.", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(95, 8, fmt.Sprintf("Factura N°: %s", factura.NumeroFactura), "1", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "RUC: 1791234567001", "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Fecha: %s", factura.FechaEmision.Format("02/01/2006")), "1", 1, "R", false, 0, "")

	pdf.CellFormat(95, 6, "Av. Principal 123, Quito", "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "Clave de Acceso:", "1", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(95, 6, "Tel: (02) 123-4567", "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, factura.ClaveAcceso, "1", 1, "R", false, 0, "")

	pdf.Ln(5)

	// Información del cliente
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(190, 8, "DATOS DEL CLIENTE", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 6, "Cliente:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, factura.ClienteNombre, "1", 1, "L", false, 0, "")

	pdf.CellFormat(60, 6, "Cédula/RUC:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(130, 6, factura.ClienteCedula, "1", 1, "L", false, 0, "")

	// Si hay información adicional del cliente
	if factura.ClienteDireccion != "" {
		pdf.CellFormat(60, 6, "Dirección:", "1", 0, "L", false, 0, "")
		pdf.CellFormat(130, 6, factura.ClienteDireccion, "1", 1, "L", false, 0, "")
	}

	if factura.ClienteTelefono != "" {
		pdf.CellFormat(60, 6, "Teléfono:", "1", 0, "L", false, 0, "")
		pdf.CellFormat(130, 6, factura.ClienteTelefono, "1", 1, "L", false, 0, "")
	}

	pdf.Ln(5)

	// Tabla de productos
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(190, 8, "DETALLE DE PRODUCTOS/SERVICIOS", "1", 1, "C", false, 0, "")

	// Headers de la tabla
	pdf.CellFormat(20, 8, "Cód.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(80, 8, "Descripción", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "Cant.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "P. Unit.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "Desc.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Total", "1", 1, "C", false, 0, "")

	// Productos
	pdf.SetFont("Arial", "", 9)
	var totalGeneral float64

	for _, producto := range productos {
		total := producto.Cantidad * producto.PrecioUnitario - producto.Descuento
		totalGeneral += total

		pdf.CellFormat(20, 6, producto.Codigo, "1", 0, "C", false, 0, "")
		
		// Descripción puede ser larga, usar MultiCell si es necesario
		if len(producto.Descripcion) > 40 {
			// Para descripciones largas, usar texto más pequeño
			pdf.SetFont("Arial", "", 8)
			pdf.CellFormat(80, 6, producto.Descripcion[:37]+"...", "1", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 9)
		} else {
			pdf.CellFormat(80, 6, producto.Descripcion, "1", 0, "L", false, 0, "")
		}
		
		pdf.CellFormat(20, 6, fmt.Sprintf("%.2f", producto.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", producto.PrecioUnitario), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("$%.2f", producto.Descuento), "1", 0, "R", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", total), "1", 1, "R", false, 0, "")
	}

	pdf.Ln(3)

	// Resumen de totales
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(140, 8, "", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 8, "RESUMEN", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(140, 6, "", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "Subtotal:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 6, fmt.Sprintf("$%.2f", factura.Subtotal), "1", 1, "R", false, 0, "")

	pdf.CellFormat(140, 6, "", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "IVA (12%):", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 6, fmt.Sprintf("$%.2f", factura.IVA), "1", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(140, 8, "", "0", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, "TOTAL:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 8, fmt.Sprintf("$%.2f", factura.Total), "1", 1, "R", false, 0, "")

	pdf.Ln(10)

	// Información adicional
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(190, 6, fmt.Sprintf("Estado: %s", factura.Estado), "0", 1, "L", false, 0, "")

	if factura.NumeroAutorizacion != "" {
		pdf.CellFormat(190, 6, fmt.Sprintf("Número de Autorización: %s", factura.NumeroAutorizacion), "0", 1, "L", false, 0, "")
	}

	if factura.FechaAutorizacion != nil {
		pdf.CellFormat(190, 6, fmt.Sprintf("Fecha de Autorización: %s", factura.FechaAutorizacion.Format("02/01/2006 15:04:05")), "0", 1, "L", false, 0, "")
	}

	pdf.CellFormat(190, 6, fmt.Sprintf("Ambiente: %s", factura.Ambiente), "0", 1, "L", false, 0, "")

	// Footer
	pdf.Ln(5)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(190, 4, "Este documento ha sido generado electrónicamente y no requiere firma autógrafa.", "0", 1, "C", false, 0, "")
	pdf.CellFormat(190, 4, fmt.Sprintf("Generado el: %s", time.Now().Format("02/01/2006 15:04:05")), "0", 1, "C", false, 0, "")

	// Convertir a bytes
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generando PDF: %v", err)
	}

	return buf.Bytes(), nil
}

// GenerarFacturaSimplePDF genera un PDF básico más simple
func (g *FacturaPDFGenerator) GenerarFacturaSimplePDF(facturaID int) ([]byte, error) {
	// Obtener factura
	factura, err := g.db.ObtenerFacturaPorID(facturaID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo factura: %v", err)
	}

	// Crear PDF simple
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(190, 15, "FACTURA", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(190, 8, fmt.Sprintf("Número: %s", factura.NumeroFactura), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Cliente: %s", factura.ClienteNombre), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Cédula/RUC: %s", factura.ClienteCedula), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Fecha: %s", factura.FechaEmision.Format("02/01/2006")), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Total: $%.2f", factura.Total), "0", 1, "L", false, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Estado: %s", factura.Estado), "0", 1, "L", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generando PDF simple: %v", err)
	}

	return buf.Bytes(), nil
}

// ValidarFacturaParaPDF valida que una factura puede generar PDF
func (g *FacturaPDFGenerator) ValidarFacturaParaPDF(facturaID int) error {
	// Verificar que la factura existe
	_, err := g.db.ObtenerFacturaPorID(facturaID)
	if err != nil {
		return fmt.Errorf("factura no encontrada: %v", err)
	}

	// Verificar que tiene productos
	productos, err := g.db.ObtenerProductosPorFactura(facturaID)
	if err != nil {
		return fmt.Errorf("error obteniendo productos: %v", err)
	}

	if len(productos) == 0 {
		return fmt.Errorf("la factura no tiene productos asociados")
	}

	return nil
}