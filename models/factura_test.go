package models

import (
	"encoding/xml"
	"strings"
	"testing"
)

// TestProductoInput_CreacionBasica verifica la creación de productos
func TestProductoInput_CreacionBasica(t *testing.T) {
	producto := ProductoInput{
		Codigo:         "LAPTOP001",
		Descripcion:    "Laptop Dell Inspiron 15",
		Cantidad:       2.0,
		PrecioUnitario: 450.00,
	}

	if producto.Codigo != "LAPTOP001" {
		t.Errorf("Codigo = %v, quería 'LAPTOP001'", producto.Codigo)
	}
	if producto.Descripcion != "Laptop Dell Inspiron 15" {
		t.Errorf("Descripcion = %v, quería 'Laptop Dell Inspiron 15'", producto.Descripcion)
	}
	if producto.Cantidad != 2.0 {
		t.Errorf("Cantidad = %v, quería 2.0", producto.Cantidad)
	}
	if producto.PrecioUnitario != 450.00 {
		t.Errorf("PrecioUnitario = %v, quería 450.00", producto.PrecioUnitario)
	}
}

// TestFacturaInput_CreacionBasica verifica la creación de entrada de facturas
func TestFacturaInput_CreacionBasica(t *testing.T) {
	input := FacturaInput{
		ClienteNombre: "Juan Carlos Pérez",
		ClienteCedula: "1713175071",
		Productos: []ProductoInput{
			{
				Codigo:         "LAPTOP001",
				Descripcion:    "Laptop Dell Inspiron 15",
				Cantidad:       1.0,
				PrecioUnitario: 450.00,
			},
			{
				Codigo:         "MOUSE001",
				Descripcion:    "Mouse Inalámbrico",
				Cantidad:       2.0,
				PrecioUnitario: 25.00,
			},
		},
	}

	if input.ClienteNombre != "Juan Carlos Pérez" {
		t.Errorf("ClienteNombre = %v, quería 'Juan Carlos Pérez'", input.ClienteNombre)
	}
	if input.ClienteCedula != "1713175071" {
		t.Errorf("ClienteCedula = %v, quería '1713175071'", input.ClienteCedula)
	}
	if len(input.Productos) != 2 {
		t.Errorf("Número de productos = %v, quería 2", len(input.Productos))
	}
}

// TestInfoTributaria_EstructuraCompleta verifica la estructura de información tributaria
func TestInfoTributaria_EstructuraCompleta(t *testing.T) {
	info := InfoTributaria{
		Ambiente:        "1",
		TipoEmision:     "1",
		RazonSocial:     "EMPRESA DEMO S.A.",
		RUC:             "1234567890001",
		ClaveAcceso:     "2306202501179214673900110010010000000019152728411",
		CodDoc:          "01",
		Establecimiento: "001",
		PuntoEmision:    "001",
		Secuencial:      "000000001",
	}

	// Verificar todos los campos
	if info.Ambiente != "1" {
		t.Errorf("Ambiente = %v, quería '1'", info.Ambiente)
	}
	if info.TipoEmision != "1" {
		t.Errorf("TipoEmision = %v, quería '1'", info.TipoEmision)
	}
	if info.RazonSocial != "EMPRESA DEMO S.A." {
		t.Errorf("RazonSocial = %v, quería 'EMPRESA DEMO S.A.'", info.RazonSocial)
	}
	if info.RUC != "1234567890001" {
		t.Errorf("RUC = %v, quería '1234567890001'", info.RUC)
	}
	if len(info.ClaveAcceso) != 49 {
		t.Errorf("ClaveAcceso longitud = %v, quería 49", len(info.ClaveAcceso))
	}
	if info.CodDoc != "01" {
		t.Errorf("CodDoc = %v, quería '01'", info.CodDoc)
	}
}

// TestInfoFactura_EstructuraCompleta verifica la estructura de información de factura
func TestInfoFactura_EstructuraCompleta(t *testing.T) {
	info := InfoFactura{
		FechaEmision:                "23/06/2025",
		DirEstablecimiento:          "Av. Amazonas y Naciones Unidas, Quito, Ecuador",
		TipoIdentificacionComprador: "05",
		IdentificacionComprador:     "1713175071",
		RazonSocialComprador:        "Juan Carlos Pérez",
		TotalSinImpuestos:           450.00,
		TotalDescuento:              0.00,
		ImporteTotal:                517.50,
		Moneda:                      "DOLAR",
	}

	// Verificar campos críticos
	if info.TipoIdentificacionComprador != "05" {
		t.Errorf("TipoIdentificacionComprador = %v, quería '05'", info.TipoIdentificacionComprador)
	}
	if info.Moneda != "DOLAR" {
		t.Errorf("Moneda = %v, quería 'DOLAR'", info.Moneda)
	}
	if info.TotalSinImpuestos != 450.00 {
		t.Errorf("TotalSinImpuestos = %v, quería 450.00", info.TotalSinImpuestos)
	}
	if info.ImporteTotal != 517.50 {
		t.Errorf("ImporteTotal = %v, quería 517.50", info.ImporteTotal)
	}
}

// TestDetalle_EstructuraCompleta verifica la estructura de detalle
func TestDetalle_EstructuraCompleta(t *testing.T) {
	detalle := Detalle{
		CodigoPrincipal:        "LAPTOP001",
		Descripcion:            "Laptop Dell Inspiron 15",
		Cantidad:               2.0,
		PrecioUnitario:         450.00,
		Descuento:              0.00,
		PrecioTotalSinImpuesto: 900.00,
	}

	if detalle.CodigoPrincipal != "LAPTOP001" {
		t.Errorf("CodigoPrincipal = %v, quería 'LAPTOP001'", detalle.CodigoPrincipal)
	}
	if detalle.Cantidad != 2.0 {
		t.Errorf("Cantidad = %v, quería 2.0", detalle.Cantidad)
	}
	if detalle.PrecioTotalSinImpuesto != 900.00 {
		t.Errorf("PrecioTotalSinImpuesto = %v, quería 900.00", detalle.PrecioTotalSinImpuesto)
	}
}

// TestFactura_EstructuraCompleta verifica la estructura completa de factura
func TestFactura_EstructuraCompleta(t *testing.T) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:        "1",
			TipoEmision:     "1",
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			ClaveAcceso:     "2306202501179214673900110010010000000019152728411",
			CodDoc:          "01",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Secuencial:      "000000001",
		},
		InfoFactura: InfoFactura{
			FechaEmision:                "23/06/2025",
			DirEstablecimiento:          "Av. Amazonas y Naciones Unidas, Quito, Ecuador",
			TipoIdentificacionComprador: "05",
			IdentificacionComprador:     "1713175071",
			RazonSocialComprador:        "Juan Carlos Pérez",
			TotalSinImpuestos:           450.00,
			TotalDescuento:              0.00,
			ImporteTotal:                517.50,
			Moneda:                      "DOLAR",
		},
		Detalles: []Detalle{
			{
				CodigoPrincipal:        "LAPTOP001",
				Descripcion:            "Laptop Dell Inspiron 15",
				Cantidad:               1.0,
				PrecioUnitario:         450.00,
				Descuento:              0.00,
				PrecioTotalSinImpuesto: 450.00,
			},
		},
	}

	// Verificar que la estructura está completa
	if len(factura.Detalles) != 1 {
		t.Errorf("Número de detalles = %v, quería 1", len(factura.Detalles))
	}
	if factura.InfoTributaria.RUC == "" {
		t.Error("RUC no debe estar vacío")
	}
	if factura.InfoFactura.IdentificacionComprador == "" {
		t.Error("IdentificacionComprador no debe estar vacío")
	}
}

// TestFactura_GenerarXML verifica la generación de XML
func TestFactura_GenerarXML(t *testing.T) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:        "1",
			TipoEmision:     "1",
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			ClaveAcceso:     "2306202501179214673900110010010000000019152728411",
			CodDoc:          "01",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Secuencial:      "000000001",
		},
		InfoFactura: InfoFactura{
			FechaEmision:                "23/06/2025",
			DirEstablecimiento:          "Av. Amazonas y Naciones Unidas, Quito, Ecuador",
			TipoIdentificacionComprador: "05",
			IdentificacionComprador:     "1713175071",
			RazonSocialComprador:        "Juan Carlos Pérez",
			TotalSinImpuestos:           450.00,
			TotalDescuento:              0.00,
			ImporteTotal:                517.50,
			Moneda:                      "DOLAR",
		},
		Detalles: []Detalle{
			{
				CodigoPrincipal:        "LAPTOP001",
				Descripcion:            "Laptop Dell Inspiron 15",
				Cantidad:               1.0,
				PrecioUnitario:         450.00,
				Descuento:              0.00,
				PrecioTotalSinImpuesto: 450.00,
			},
		},
	}

	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("GenerarXML() error = %v, no quería error", err)
	}

	xmlString := string(xmlData)

	// Verificar que el XML contiene elementos esperados
	expectedElements := []string{
		"<factura>",
		"<infoTributaria>",
		"<infoFactura>",
		"<detalles>",
		"<detalle>",
		"<ambiente>1</ambiente>",
		"<ruc>1234567890001</ruc>",
		"<razonSocial>EMPRESA DEMO S.A.</razonSocial>",
		"<identificacionComprador>1713175071</identificacionComprador>",
		"<codigoPrincipal>LAPTOP001</codigoPrincipal>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(xmlString, element) {
			t.Errorf("XML no contiene elemento esperado: %s", element)
		}
	}

	// Verificar que es XML válido intentando unmarshalling
	var facturaUnmarshaled Factura
	err = xml.Unmarshal(xmlData, &facturaUnmarshaled)
	if err != nil {
		t.Errorf("XML generado no es válido: %v", err)
	}
}

// TestFactura_GenerarXML_MultipleProductos verifica XML con múltiples productos
func TestFactura_GenerarXML_MultipleProductos(t *testing.T) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:        "1",
			TipoEmision:     "1",
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			ClaveAcceso:     "2306202501179214673900110010010000000019152728411",
			CodDoc:          "01",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Secuencial:      "000000002",
		},
		InfoFactura: InfoFactura{
			FechaEmision:                "23/06/2025",
			DirEstablecimiento:          "Av. Amazonas y Naciones Unidas, Quito, Ecuador",
			TipoIdentificacionComprador: "05",
			IdentificacionComprador:     "0926687856",
			RazonSocialComprador:        "María González",
			TotalSinImpuestos:           975.00,
			TotalDescuento:              0.00,
			ImporteTotal:                1121.25,
			Moneda:                      "DOLAR",
		},
		Detalles: []Detalle{
			{
				CodigoPrincipal:        "LAPTOP001",
				Descripcion:            "Laptop Dell Inspiron 15",
				Cantidad:               2.0,
				PrecioUnitario:         450.00,
				Descuento:              0.00,
				PrecioTotalSinImpuesto: 900.00,
			},
			{
				CodigoPrincipal:        "MOUSE001",
				Descripcion:            "Mouse Inalámbrico",
				Cantidad:               3.0,
				PrecioUnitario:         25.00,
				Descuento:              0.00,
				PrecioTotalSinImpuesto: 75.00,
			},
		},
	}

	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("GenerarXML() error = %v, no quería error", err)
	}

	xmlString := string(xmlData)

	// Verificar que contiene múltiples productos
	laptopCount := strings.Count(xmlString, "LAPTOP001")
	mouseCount := strings.Count(xmlString, "MOUSE001")
	detalleCount := strings.Count(xmlString, "<detalle>")

	if laptopCount != 1 {
		t.Errorf("LAPTOP001 aparece %d veces, quería 1", laptopCount)
	}
	if mouseCount != 1 {
		t.Errorf("MOUSE001 aparece %d veces, quería 1", mouseCount)
	}
	if detalleCount != 2 {
		t.Errorf("<detalle> aparece %d veces, quería 2", detalleCount)
	}
}

// TestFactura_GenerarXML_XMLStructureTags verifica que los tags XML funcionen correctamente
func TestFactura_GenerarXML_XMLStructureTags(t *testing.T) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:    "1",
			RazonSocial: "Test Company",
			RUC:         "1234567890001",
		},
		InfoFactura: InfoFactura{
			FechaEmision:         "23/06/2025",
			TotalSinImpuestos:    100.00,
			ImporteTotal:         115.00,
		},
		Detalles: []Detalle{
			{
				CodigoPrincipal: "TEST001",
				Descripcion:     "Test Product",
				Cantidad:        1.0,
			},
		},
	}

	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("GenerarXML() error = %v, no quería error", err)
	}

	xmlString := string(xmlData)

	// Verificar que los tags XML específicos estén presentes
	expectedTags := []string{
		"<ambiente>1</ambiente>",
		"<razonSocial>Test Company</razonSocial>",
		"<ruc>1234567890001</ruc>",
		"<fechaEmision>23/06/2025</fechaEmision>",
		"<totalSinImpuestos>100</totalSinImpuestos>",
		"<importeTotal>115</importeTotal>",
		"<codigoPrincipal>TEST001</codigoPrincipal>",
		"<descripcion>Test Product</descripcion>",
		"<cantidad>1</cantidad>",
	}

	for _, tag := range expectedTags {
		if !strings.Contains(xmlString, tag) {
			t.Errorf("XML no contiene tag esperado: %s", tag)
		}
	}
}

// TestFactura_MostrarResumen verifica que el método MostrarResumen no falle
func TestFactura_MostrarResumen(t *testing.T) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Secuencial: "000000001",
		},
		InfoFactura: InfoFactura{
			RazonSocialComprador:    "Juan Pérez",
			IdentificacionComprador: "1713175071",
			TotalSinImpuestos:       450.00,
			ImporteTotal:            517.50,
		},
		Detalles: []Detalle{
			{
				Descripcion:            "Laptop Dell",
				Cantidad:               1.0,
				PrecioUnitario:         450.00,
				PrecioTotalSinImpuesto: 450.00,
			},
		},
	}

	// MostrarResumen imprime a stdout, no retorna valores
	// Solo verificamos que no lance panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MostrarResumen() causó panic: %v", r)
		}
	}()

	factura.MostrarResumen()
}

// TestFactura_GenerarXML_FacturaVacia verifica manejo de factura vacía
func TestFactura_GenerarXML_FacturaVacia(t *testing.T) {
	factura := Factura{}

	xmlData, err := factura.GenerarXML()
	if err != nil {
		t.Fatalf("GenerarXML() con factura vacía error = %v, no quería error", err)
	}

	xmlString := string(xmlData)

	// Debe generar XML válido aunque esté vacío
	if !strings.Contains(xmlString, "<factura>") {
		t.Error("XML no contiene elemento raíz <factura>")
	}
	if !strings.Contains(xmlString, "</factura>") {
		t.Error("XML no contiene cierre de elemento </factura>")
	}
}

// Benchmark para GenerarXML con una factura simple
func BenchmarkFactura_GenerarXML_Simple(b *testing.B) {
	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:        "1",
			RazonSocial:     "Test Company",
			RUC:             "1234567890001",
			ClaveAcceso:     "2306202501179214673900110010010000000019152728411",
		},
		InfoFactura: InfoFactura{
			FechaEmision:      "23/06/2025",
			TotalSinImpuestos: 100.00,
			ImporteTotal:      115.00,
		},
		Detalles: []Detalle{
			{
				CodigoPrincipal:        "BENCH001",
				Descripcion:            "Benchmark Product",
				Cantidad:               1.0,
				PrecioUnitario:         100.00,
				PrecioTotalSinImpuesto: 100.00,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := factura.GenerarXML()
		if err != nil {
			b.Fatalf("GenerarXML() error = %v", err)
		}
	}
}

// Benchmark para GenerarXML con múltiples productos
func BenchmarkFactura_GenerarXML_MultipleProductos(b *testing.B) {
	detalles := make([]Detalle, 10)
	for i := 0; i < 10; i++ {
		detalles[i] = Detalle{
			CodigoPrincipal:        "PROD" + string(rune(i+48)), // ASCII 48 = '0'
			Descripcion:            "Producto " + string(rune(i+48)),
			Cantidad:               1.0,
			PrecioUnitario:         100.00,
			PrecioTotalSinImpuesto: 100.00,
		}
	}

	factura := Factura{
		InfoTributaria: InfoTributaria{
			Ambiente:    "1",
			RazonSocial: "Test Company",
			RUC:         "1234567890001",
		},
		InfoFactura: InfoFactura{
			FechaEmision:      "23/06/2025",
			TotalSinImpuestos: 1000.00,
			ImporteTotal:      1150.00,
		},
		Detalles: detalles,
	}

	for i := 0; i < b.N; i++ {
		_, err := factura.GenerarXML()
		if err != nil {
			b.Fatalf("GenerarXML() error = %v", err)
		}
	}
}