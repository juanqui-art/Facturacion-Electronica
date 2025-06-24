package sri

import (
	"strings"
	"testing"
	"time"
)

// TestCrearFirmaXAdESBES tests XAdES-BES signature creation
func TestCrearFirmaXAdESBES(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <infoTributaria>
        <ambiente>1</ambiente>
        <tipoEmision>1</tipoEmision>
        <razonSocial>EMPRESA TEST</razonSocial>
        <ruc>1792146739001</ruc>
        <claveAcceso>2306202401179214673900110010010000000011234567891</claveAcceso>
    </infoTributaria>
</factura>`

	// Crear un certificado mock para testing
	certificado := &CertificadoDigital{
		Archivo:    "test_cert.p12",
		Password:   "test_password",
		PrivateKey: nil, // En test real sería la clave privada
		Cert:       nil, // En test real sería el certificado
		CACerts:    nil, // En test real serían los certificados de CA
	}

	config := XAdESBESConfig{
		Certificado: certificado,
		PolicyID:    "https://www.sri.gob.ec/politica-firma",
		PolicyHash:  "test_hash",
		PolicyURL:   "https://www.sri.gob.ec/politica-firma",
	}

	xmlFirmado, err := FirmarXMLXAdESBES([]byte(xmlData), config)

	if err != nil {
		t.Fatalf("FirmarXMLXAdESBES() error = %v", err)
	}

	if len(xmlFirmado) == 0 {
		t.Fatal("XML firmado no puede estar vacío")
	}

	xmlFirmadoStr := string(xmlFirmado)

	// Verificar que contiene elementos de firma XAdES
	expectedElements := []string{
		"<ds:Signature",
		"<ds:SignedInfo",
		"<ds:CanonicalizationMethod",
		"<ds:SignatureMethod",
		"<ds:Reference",
		"<ds:DigestMethod",
		"<ds:DigestValue",
		"<ds:SignatureValue",
		"<ds:KeyInfo",
		"<xades:QualifyingProperties",
		"<xades:SignedProperties",
		"<xades:SignedSignatureProperties",
		"<xades:SigningTime",
		"<xades:SigningCertificate",
	}

	for _, element := range expectedElements {
		if !strings.Contains(xmlFirmadoStr, element) {
			t.Errorf("XML firmado debería contener elemento: %s", element)
		}
	}

	// Verificar que el XML original está preservado
	if !strings.Contains(xmlFirmadoStr, "EMPRESA TEST") {
		t.Error("XML firmado debería preservar el contenido original")
	}

	if !strings.Contains(xmlFirmadoStr, "1792146739001") {
		t.Error("XML firmado debería preservar el RUC")
	}
}

// TestCrearFirmaXAdESBESErrores tests error handling in XAdES creation
func TestCrearFirmaXAdESBESErrores(t *testing.T) {
	validXML := `<?xml version="1.0" encoding="UTF-8"?><factura><test>content</test></factura>`
	validCert := &CertificadoDigital{
		Archivo:    "test_cert.p12",
		Password:   "test_password",
		PrivateKey: nil,
		Cert:       nil,
		CACerts:    nil,
	}

	tests := []struct {
		name        string
		xmlData     []byte
		certificado *CertificadoDigital
		expectError bool
	}{
		{
			name:        "XML vacío",
			xmlData:     []byte(""),
			certificado: validCert,
			expectError: true,
		},
		{
			name:        "XML nil",
			xmlData:     nil,
			certificado: validCert,
			expectError: true,
		},
		{
			name:        "XML malformado",
			xmlData:     []byte("<invalid><xml>"),
			certificado: validCert,
			expectError: true,
		},
		{
			name:        "Certificado nil",
			xmlData:     []byte(validXML),
			certificado: nil,
			expectError: true,
		},
		{
			name:    "Certificado inválido",
			xmlData: []byte(validXML),
			certificado: &CertificadoDigital{
				Archivo:  "invalid_cert.p12",
				Password: "wrong_password",
			},
			expectError: true,
		},
		{
			name:    "Certificado expirado",
			xmlData: []byte(validXML),
			certificado: &CertificadoDigital{
				Archivo:    "expired_cert.p12",
				Password:   "test_password",
				PrivateKey: nil,
				Cert:       nil, // Certificado expirado sería nil o inválido
				CACerts:    nil,
			},
			expectError: true,
		},
		{
			name:    "Certificado sin datos",
			xmlData: []byte(validXML),
			certificado: &CertificadoDigital{
				Archivo:    "",
				Password:   "",
				PrivateKey: nil,
				Cert:       nil,
				CACerts:    nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := XAdESBESConfig{
				Certificado: tt.certificado,
				PolicyID:    "https://www.sri.gob.ec/politica-firma",
				PolicyHash:  "test_hash",
				PolicyURL:   "https://www.sri.gob.ec/politica-firma",
			}
			xmlFirmado, err := FirmarXMLXAdESBES(tt.xmlData, config)

			if tt.expectError {
				if err == nil {
					t.Error("FirmarXMLXAdESBES() esperaba error, obtuvo nil")
				}
				if len(xmlFirmado) > 0 {
					t.Error("XML firmado debería estar vacío en caso de error")
				}
			} else {
				if err != nil {
					t.Errorf("FirmarXMLXAdESBES() error inesperado = %v", err)
				}
				if len(xmlFirmado) == 0 {
					t.Error("XML firmado no debería estar vacío en caso de éxito")
				}
			}
		})
	}
}

// TestValidarFirmaXAdES tests XAdES signature validation
func TestValidarFirmaXAdES(t *testing.T) {
	// XML con firma simulada
	xmlConFirma := `<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <infoTributaria>
        <ambiente>1</ambiente>
        <ruc>1792146739001</ruc>
    </infoTributaria>
    <ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <ds:SignedInfo>
            <ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/>
            <ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
        </ds:SignedInfo>
        <ds:SignatureValue>mock signature value</ds:SignatureValue>
        <xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#">
            <xades:SignedProperties>
                <xades:SignedSignatureProperties>
                    <xades:SigningTime>2024-06-24T10:30:00Z</xades:SigningTime>
                </xades:SignedSignatureProperties>
            </xades:SignedProperties>
        </xades:QualifyingProperties>
    </ds:Signature>
</factura>`

	xmlSinFirma := `<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <infoTributaria>
        <ambiente>1</ambiente>
        <ruc>1792146739001</ruc>
    </infoTributaria>
</factura>`

	tests := []struct {
		name        string
		xmlData     []byte
		expectValid bool
	}{
		{
			name:        "XML con firma XAdES",
			xmlData:     []byte(xmlConFirma),
			expectValid: true,
		},
		{
			name:        "XML sin firma",
			xmlData:     []byte(xmlSinFirma),
			expectValid: false,
		},
		{
			name:        "XML vacío",
			xmlData:     []byte(""),
			expectValid: false,
		},
		{
			name:        "XML malformado",
			xmlData:     []byte("<invalid><xml>"),
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarFirmaXAdESBES(tt.xmlData)

			if tt.expectValid {
				if err != nil {
					t.Errorf("ValidarFirmaXAdESBES() error inesperado = %v", err)
				}
			} else {
				// Para casos inválidos, esperamos un error
				if err == nil {
					t.Error("ValidarFirmaXAdESBES() debería retornar error para XML inválido")
				}
			}
		})
	}
}

// TestExtraerCertificadoDeXML tests certificate extraction from XML
func TestExtraerCertificadoDeXML(t *testing.T) {
	xmlConCertificado := `<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <ds:KeyInfo>
            <ds:X509Data>
                <ds:X509Certificate>TW9jayBjZXJ0aWZpY2F0ZSBkYXRh</ds:X509Certificate>
            </ds:X509Data>
        </ds:KeyInfo>
    </ds:Signature>
</factura>`

	xmlSinCertificado := `<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <infoTributaria>
        <ruc>1792146739001</ruc>
    </infoTributaria>
</factura>`

	tests := []struct {
		name          string
		xmlData       []byte
		expectCert    bool
		expectError   bool
	}{
		{
			name:        "XML con certificado",
			xmlData:     []byte(xmlConCertificado),
			expectCert:  true,
			expectError: false,
		},
		{
			name:        "XML sin certificado",
			xmlData:     []byte(xmlSinCertificado),
			expectCert:  false,
			expectError: true,
		},
		{
			name:        "XML vacío",
			xmlData:     []byte(""),
			expectCert:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certificado, err := ExtraerCertificadoDeXML(tt.xmlData)

			if tt.expectError {
				if err == nil {
					t.Error("ExtraerCertificadoDeXML() esperaba error, obtuvo nil")
				}
			}

			if tt.expectCert {
				if certificado == nil {
					t.Error("ExtraerCertificadoDeXML() debería retornar certificado")
				}
				if err != nil {
					t.Errorf("ExtraerCertificadoDeXML() error inesperado = %v", err)
				}
			} else {
				if certificado != nil {
					t.Error("ExtraerCertificadoDeXML() no debería retornar certificado")
				}
			}
		})
	}
}

// TestGenerarHashSHA1 tests SHA1 hash generation
func TestGenerarHashSHA1(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "Datos básicos",
			data:     []byte("test data"),
			expected: "916f0027a575074ce72a331777c3478d6513f786", // SHA1 de "test data"
		},
		{
			name:     "Datos vacíos",
			data:     []byte(""),
			expected: "da39a3ee5e6b4b0d3255bfef95601890afd80709", // SHA1 de string vacío
		},
		{
			name:     "XML simple",
			data:     []byte("<test>content</test>"),
			expected: "", // No verificamos hash específico, solo que no sea vacío
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := GenerarHashSHA1(tt.data)

			if len(hash) != 40 { // SHA1 hash is always 40 hex characters
				t.Errorf("Hash SHA1 debería tener 40 caracteres, obtuvo %d", len(hash))
			}

			if tt.expected != "" && hash != tt.expected {
				t.Errorf("Hash SHA1 esperado: %s, obtuvo: %s", tt.expected, hash)
			}

			// Verificar que es hexadecimal válido
			for _, char := range hash {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
					t.Errorf("Hash contiene carácter no hexadecimal: %c", char)
				}
			}
		})
	}
}

// TestCrearTimestamp tests timestamp creation for XAdES
func TestCrearTimestamp(t *testing.T) {
	timestamp := CrearTimestamp()

	if timestamp == "" {
		t.Error("Timestamp no puede estar vacío")
	}

	// Verificar formato ISO 8601
	if !strings.Contains(timestamp, "T") {
		t.Error("Timestamp debería contener 'T' (formato ISO 8601)")
	}

	if !strings.HasSuffix(timestamp, "Z") {
		t.Error("Timestamp debería terminar con 'Z' (UTC)")
	}

	// Verificar que es una fecha válida reciente
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Errorf("Timestamp no es una fecha RFC3339 válida: %v", err)
	}
}

// TestNormalizarXML tests XML normalization
func TestNormalizarXML(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "XML con espacios",
			input:    []byte("  <test>  content  </test>  "),
			expected: []byte("<test>content</test>"),
		},
		{
			name:     "XML con saltos de línea",
			input:    []byte("<test>\n  <child>value</child>\n</test>"),
			expected: []byte("<test><child>value</child></test>"),
		},
		{
			name:     "XML ya normalizado",
			input:    []byte("<test><child>value</child></test>"),
			expected: []byte("<test><child>value</child></test>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := NormalizarXML(tt.input)
			
			if string(resultado) != string(tt.expected) {
				t.Errorf("NormalizarXML() = %s, esperado %s", string(resultado), string(tt.expected))
			}
		})
	}
}

// TestCrearDigestValue tests digest value creation
func TestCrearDigestValue(t *testing.T) {
	xmlData := []byte("<test>content</test>")
	digest := CrearDigestValue(xmlData)

	if digest == "" {
		t.Error("Digest value no puede estar vacío")
	}

	// Verificar que es Base64 válido
	if len(digest)%4 != 0 {
		t.Error("Digest value debería ser Base64 válido (longitud múltiplo de 4)")
	}

	// El mismo input debería producir el mismo digest
	digest2 := CrearDigestValue(xmlData)
	if digest != digest2 {
		t.Error("El mismo XML debería producir el mismo digest")
	}

	// Diferente input debería producir diferente digest
	digest3 := CrearDigestValue([]byte("<test>different</test>"))
	if digest == digest3 {
		t.Error("Diferente XML debería producir diferente digest")
	}
}

// BenchmarkCrearFirmaXAdESBES benchmarks XAdES signature creation
func BenchmarkCrearFirmaXAdESBES(b *testing.B) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?><factura><test>content</test></factura>`)
	certificado := &CertificadoDigital{
		Archivo:    "test_cert.p12",
		Password:   "test_password",
		PrivateKey: nil,
		Cert:       nil,
		CACerts:    nil,
	}

	config := XAdESBESConfig{
		Certificado: certificado,
		PolicyID:    "https://www.sri.gob.ec/politica-firma",
		PolicyHash:  "test_hash",
		PolicyURL:   "https://www.sri.gob.ec/politica-firma",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = FirmarXMLXAdESBES(xmlData, config)
	}
}

// BenchmarkGenerarHashSHA1 benchmarks SHA1 hash generation
func BenchmarkGenerarHashSHA1(b *testing.B) {
	data := []byte("test data for hashing performance")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GenerarHashSHA1(data)
	}
}

// TestIntegracionXAdESCompleta tests complete XAdES integration
func TestIntegracionXAdESCompleta(t *testing.T) {
	// Crear XML de factura
	xmlOriginal := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<factura>
    <infoTributaria>
        <ambiente>1</ambiente>
        <tipoEmision>1</tipoEmision>
        <razonSocial>EMPRESA TEST XADES</razonSocial>
        <ruc>1792146739001</ruc>
        <claveAcceso>2306202401179214673900110010010000000011234567891</claveAcceso>
    </infoTributaria>
    <infoFactura>
        <fechaEmision>23/06/2024</fechaEmision>
        <totalSinImpuestos>100.00</totalSinImpuestos>
        <totalConImpuestos>115.00</totalConImpuestos>
        <importeTotal>115.00</importeTotal>
    </infoFactura>
</factura>`)

	// Crear certificado mock
	certificado := &CertificadoDigital{
		Archivo:    "integration_test_cert.p12",
		Password:   "test_password",
		PrivateKey: nil,
		Cert:       nil,
		CACerts:    nil,
	}

	config := XAdESBESConfig{
		Certificado: certificado,
		PolicyID:    "https://www.sri.gob.ec/politica-firma",
		PolicyHash:  "test_hash",
		PolicyURL:   "https://www.sri.gob.ec/politica-firma",
	}

	// Paso 1: Crear firma XAdES
	xmlFirmado, err := FirmarXMLXAdESBES(xmlOriginal, config)
	if err != nil {
		t.Fatalf("Error creando firma XAdES: %v", err)
	}

	// Paso 2: Validar que el XML firmado contiene los elementos necesarios
	xmlFirmadoStr := string(xmlFirmado)
	
	requiredElements := []string{
		"<ds:Signature",
		"<xades:QualifyingProperties",
		"<xades:SigningTime",
		"EMPRESA TEST XADES", // Contenido original preservado
		"1792146739001",      // RUC preservado
	}

	for _, element := range requiredElements {
		if !strings.Contains(xmlFirmadoStr, element) {
			t.Errorf("XML firmado no contiene elemento requerido: %s", element)
		}
	}

	// Paso 3: Validar firma
	err = ValidarFirmaXAdESBES(xmlFirmado)
	if err != nil {
		t.Errorf("Error validando firma XAdES: %v", err)
	}

	// Paso 4: Extraer certificado del XML firmado
	certExtraido, err := ExtraerCertificadoDeXML(xmlFirmado)
	if err != nil {
		t.Errorf("Error extrayendo certificado: %v", err)
	}

	if certExtraido == nil {
		t.Error("Debería poder extraer certificado del XML firmado")
	}

	t.Logf("✅ Integración XAdES completada exitosamente")
	t.Logf("📄 XML original: %d bytes", len(xmlOriginal))
	t.Logf("🔒 XML firmado: %d bytes", len(xmlFirmado))
	t.Logf("📋 Certificado extraído: %v", certExtraido != nil)
}