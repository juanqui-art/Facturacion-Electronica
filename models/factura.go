// Package models contiene todas las estructuras de datos del sistema de facturación
package models

import (
	"encoding/xml"
	"fmt"
)

// ProductoInput - Datos de un producto individual
type ProductoInput struct {
	Codigo         string
	Descripcion    string
	Cantidad       float64
	PrecioUnitario float64
}

// FacturaInput - Datos simples para crear una factura
// Ahora soporta múltiples productos!
type FacturaInput struct {
	ClienteNombre string
	ClienteCedula string
	Productos     []ProductoInput // Slice de productos!
}

// InfoTributaria - Datos básicos del emisor (obligatorios SRI)
type InfoTributaria struct {
	Ambiente        string `xml:"ambiente"`
	TipoEmision     string `xml:"tipoEmision"`
	RazonSocial     string `xml:"razonSocial"`
	RUC             string `xml:"ruc"`
	ClaveAcceso     string `xml:"claveAcceso"`
	CodDoc          string `xml:"codDoc"`
	Establecimiento string `xml:"estab"`
	PuntoEmision    string `xml:"ptoEmi"`
	Secuencial      string `xml:"secuencial"`
}

// InfoFactura - Datos específicos de la factura
type InfoFactura struct {
	FechaEmision                string  `xml:"fechaEmision"`
	DirEstablecimiento          string  `xml:"dirEstablecimiento"`
	TipoIdentificacionComprador string  `xml:"tipoIdentificacionComprador"`
	IdentificacionComprador     string  `xml:"identificacionComprador"`
	RazonSocialComprador        string  `xml:"razonSocialComprador"`
	TotalSinImpuestos           float64 `xml:"totalSinImpuestos"`
	TotalDescuento              float64 `xml:"totalDescuento"`
	ImporteTotal                float64 `xml:"importeTotal"`
	Moneda                      string  `xml:"moneda"`
}

// Detalle - Item individual de la factura
type Detalle struct {
	CodigoPrincipal        string  `xml:"codigoPrincipal"`
	Descripcion            string  `xml:"descripcion"`
	Cantidad               float64 `xml:"cantidad"`
	PrecioUnitario         float64 `xml:"precioUnitario"`
	Descuento              float64 `xml:"descuento"`
	PrecioTotalSinImpuesto float64 `xml:"precioTotalSinImpuesto"`
}

// Factura - Estructura completa del documento
type Factura struct {
	XMLName        xml.Name       `xml:"factura"`
	InfoTributaria InfoTributaria `xml:"infoTributaria"`
	InfoFactura    InfoFactura    `xml:"infoFactura"`
	Detalles       []Detalle      `xml:"detalles>detalle"`
}

// GenerarXML - Método que convierte la factura a XML
// Receiver: (f Factura) significa que este método "pertenece" a cualquier Factura
func (f Factura) GenerarXML() ([]byte, error) {
	// xml.MarshalIndent formatea el XML con indentación bonita
	xmlData, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return nil, err // nil es el valor "vacío" para []byte
	}
	return xmlData, nil
}

// MostrarResumen - Método que imprime un resumen de la factura
func (f Factura) MostrarResumen() {
	fmt.Println("=== FACTURA ELECTRÓNICA ECUATORIANA ===")
	fmt.Printf("Secuencial: %s\n", f.InfoTributaria.Secuencial)
	fmt.Printf("Cliente: %s (%s)\n", 
		f.InfoFactura.RazonSocialComprador, 
		f.InfoFactura.IdentificacionComprador)
	
	// Mostrar productos
	for i, detalle := range f.Detalles {
		fmt.Printf("Producto %d: %s\n", i+1, detalle.Descripcion)
		fmt.Printf("Cantidad: %.0f x $%.2f = $%.2f\n", 
			detalle.Cantidad, 
			detalle.PrecioUnitario, 
			detalle.PrecioTotalSinImpuesto)
	}
	
	fmt.Printf("IVA 15%%: $%.2f\n", f.InfoFactura.ImporteTotal - f.InfoFactura.TotalSinImpuestos)
	fmt.Printf("TOTAL: $%.2f\n", f.InfoFactura.ImporteTotal)
	fmt.Println()
}