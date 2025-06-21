package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FacturaInput - Datos simples para crear una factura
// Esto es lo que el usuario final proporcionará
type FacturaInput struct {
	ClienteNombre       string
	ClienteCedula       string
	ProductoCodigo      string
	ProductoDescripcion string
	Cantidad            float64
	PrecioUnitario      float64
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

// validarCedula - Valida que una cédula ecuatoriana sea correcta
// Devuelve un error si la cédula no es válida
func validarCedula(cedula string) error {
	// Verificar longitud
	if len(cedula) != 10 {
		return errors.New("la cédula debe tener exactamente 10 dígitos")
	}
	
	// Verificar que todos sean números
	for _, char := range cedula {
		if char < '0' || char > '9' {
			return errors.New("la cédula solo puede contener números")
		}
	}
	
	// Verificar que los dos primeros dígitos sean válidos (01-24)
	provincia, err := strconv.Atoi(cedula[:2])
	if err != nil {
		return errors.New("error al procesar los primeros dos dígitos de la cédula")
	}
	
	if provincia < 1 || provincia > 24 {
		return errors.New("los dos primeros dígitos de la cédula deben estar entre 01 y 24")
	}
	
	// Algoritmo de validación del dígito verificador
	coeficientes := []int{2, 1, 2, 1, 2, 1, 2, 1, 2}
	suma := 0
	
	for i := 0; i < 9; i++ {
		digito, _ := strconv.Atoi(string(cedula[i]))
		resultado := digito * coeficientes[i]
		
		if resultado >= 10 {
			resultado = resultado - 9
		}
		
		suma += resultado
	}
	
	digitoVerificador := suma % 10
	if digitoVerificador != 0 {
		digitoVerificador = 10 - digitoVerificador
	}
	
	ultimoDigito, _ := strconv.Atoi(string(cedula[9]))
	
	if digitoVerificador != ultimoDigito {
		return errors.New("el dígito verificador de la cédula no es válido")
	}
	
	return nil // nil significa "no hay error"
}

// validarFacturaInput - Valida todos los datos de entrada
func validarFacturaInput(input FacturaInput) error {
	// Validar nombre del cliente
	if input.ClienteNombre == "" {
		return errors.New("el nombre del cliente no puede estar vacío")
	}
	
	// Validar cédula
	if err := validarCedula(input.ClienteCedula); err != nil {
		return fmt.Errorf("cédula inválida: %v", err)
	}
	
	// Validar código de producto
	if input.ProductoCodigo == "" {
		return errors.New("el código del producto no puede estar vacío")
	}
	
	// Validar descripción
	if input.ProductoDescripcion == "" {
		return errors.New("la descripción del producto no puede estar vacía")
	}
	
	// Validar cantidad
	if input.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	
	// Validar precio
	if input.PrecioUnitario <= 0 {
		return errors.New("el precio unitario debe ser mayor a cero")
	}
	
	return nil
}

// CrearFactura - Función factory que crea una factura completa
// Recibe datos simples y devuelve una estructura completa lista para XML
// Ahora devuelve (Factura, error) - dos valores!
func CrearFactura(input FacturaInput) (Factura, error) {
	// Primero validamos los datos de entrada
	if err := validarFacturaInput(input); err != nil {
		return Factura{}, err // Devolvemos factura vacía y el error
	}
	// Calcular totales automáticamente
	subtotal := input.Cantidad * input.PrecioUnitario
	iva := subtotal * 0.15  // 15% IVA Ecuador
	total := subtotal + iva
	
	// Crear la factura completa con valores por defecto
	factura := Factura{
		InfoTributaria: InfoTributaria{
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
		InfoFactura: InfoFactura{
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
		Detalles: []Detalle{
			{
				CodigoPrincipal:        input.ProductoCodigo,
				Descripcion:            input.ProductoDescripcion,
				Cantidad:               input.Cantidad,
				PrecioUnitario:         input.PrecioUnitario,
				Descuento:              0.00,
				PrecioTotalSinImpuesto: subtotal,
			},
		},
	}
	
	return factura, nil // nil significa "no hay error"
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

func main() {
	// Primero, ejecutar pruebas de validación
	probarValidaciones()
	
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🚀 GENERANDO FACTURA PRINCIPAL")
	fmt.Println(strings.Repeat("=", 50))
	
	// Crear datos de factura - ¡Mucho más simple!
	facturaData := FacturaInput{
		ClienteNombre:       "JUAN CARLOS PEREZ",
		ClienteCedula:       "1713175071", // Cédula válida para Ecuador
		ProductoCodigo:      "LAPTOP001",
		ProductoDescripcion: "Laptop Dell Inspiron 15",
		Cantidad:            2.0,
		PrecioUnitario:      450.00,
	}
	
	// Generar factura usando nuestra función factory
	factura, err := CrearFactura(facturaData)
	if err != nil {
		fmt.Printf("Error al crear la factura: %v\n", err)
		return
	}
	
	// Mostrar resumen usando el método de la factura
	factura.MostrarResumen()
	
	// Generar XML usando el método de la factura
	xmlData, err := factura.GenerarXML()
	if err != nil {
		fmt.Printf("Error generando XML: %v\n", err)
		return
	}

	fmt.Println("=== XML GENERADO ===")
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", xmlData)
}

// Por ahora generamos una clave fake - en semana 4 implementaremos el algoritmo real del SRI
func generarClaveAcceso() string {
	return "2025062001123456789000110010010000000011234567890"
}
