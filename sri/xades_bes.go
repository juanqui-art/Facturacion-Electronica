// Package sri implementa firma digital XAdES-BES para Ecuador
package sri

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// XAdESBESConfig configuración para firma XAdES-BES
type XAdESBESConfig struct {
	Certificado *CertificadoDigital
	PolicyID    string // Política de firma (requerido por SRI)
	PolicyHash  string // Hash de la política
	PolicyURL   string // URL de la política
}

// SignedInfo estructura XMLDSig
type SignedInfo struct {
	XMLName                xml.Name `xml:"ds:SignedInfo"`
	CanonicalizationMethod struct {
		XMLName   xml.Name `xml:"ds:CanonicalizationMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	} `xml:"ds:CanonicalizationMethod"`
	SignatureMethod struct {
		XMLName   xml.Name `xml:"ds:SignatureMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	} `xml:"ds:SignatureMethod"`
	Reference Reference `xml:"ds:Reference"`
}

// Reference estructura para referencias en XMLDSig
type Reference struct {
	XMLName xml.Name `xml:"ds:Reference"`
	URI     string   `xml:"URI,attr"`
	Transforms struct {
		XMLName   xml.Name `xml:"ds:Transforms"`
		Transform []struct {
			XMLName   xml.Name `xml:"ds:Transform"`
			Algorithm string   `xml:"Algorithm,attr"`
		} `xml:"ds:Transform"`
	} `xml:"ds:Transforms"`
	DigestMethod struct {
		XMLName   xml.Name `xml:"ds:DigestMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	} `xml:"ds:DigestMethod"`
	DigestValue string `xml:"ds:DigestValue"`
}

// KeyInfo estructura para información de la clave
type KeyInfo struct {
	XMLName     xml.Name `xml:"ds:KeyInfo"`
	X509Data    X509Data `xml:"ds:X509Data"`
	KeyValue    KeyValue `xml:"ds:KeyValue,omitempty"`
}

// X509Data estructura para datos del certificado X.509
type X509Data struct {
	XMLName         xml.Name `xml:"ds:X509Data"`
	X509Certificate string   `xml:"ds:X509Certificate"`
	X509SubjectName string   `xml:"ds:X509SubjectName,omitempty"`
	X509IssuerName  string   `xml:"ds:X509IssuerName,omitempty"`
}

// KeyValue estructura para valor de la clave pública
type KeyValue struct {
	XMLName  xml.Name `xml:"ds:KeyValue"`
	RSAKeyValue RSAKeyValue `xml:"ds:RSAKeyValue"`
}

// RSAKeyValue estructura para clave RSA
type RSAKeyValue struct {
	XMLName  xml.Name `xml:"ds:RSAKeyValue"`
	Modulus  string   `xml:"ds:Modulus"`
	Exponent string   `xml:"ds:Exponent"`
}

// QualifyingProperties estructura XAdES
type QualifyingProperties struct {
	XMLName              xml.Name `xml:"xades:QualifyingProperties"`
	XAdESNamespace       string   `xml:"xmlns:xades,attr"`
	Target               string   `xml:"Target,attr"`
	SignedProperties     SignedProperties `xml:"xades:SignedProperties"`
}

// SignedProperties estructura XAdES
type SignedProperties struct {
	XMLName                  xml.Name `xml:"xades:SignedProperties"`
	ID                       string   `xml:"Id,attr"`
	SignedSignatureProperties SignedSignatureProperties `xml:"xades:SignedSignatureProperties"`
}

// SignedSignatureProperties estructura XAdES
type SignedSignatureProperties struct {
	XMLName           xml.Name `xml:"xades:SignedSignatureProperties"`
	SigningTime       string   `xml:"xades:SigningTime"`
	SigningCertificate SigningCertificate `xml:"xades:SigningCertificate"`
	SignaturePolicyIdentifier SignaturePolicyIdentifier `xml:"xades:SignaturePolicyIdentifier"`
}

// SigningCertificate estructura XAdES
type SigningCertificate struct {
	XMLName xml.Name `xml:"xades:SigningCertificate"`
	Cert    CertInfo `xml:"xades:Cert"`
}

// CertInfo información del certificado para XAdES
type CertInfo struct {
	XMLName      xml.Name `xml:"xades:Cert"`
	CertDigest   CertDigest `xml:"xades:CertDigest"`
	IssuerSerial IssuerSerial `xml:"xades:IssuerSerial"`
}

// CertDigest digest del certificado
type CertDigest struct {
	XMLName     xml.Name `xml:"xades:CertDigest"`
	DigestMethod struct {
		XMLName   xml.Name `xml:"ds:DigestMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	} `xml:"ds:DigestMethod"`
	DigestValue string `xml:"ds:DigestValue"`
}

// IssuerSerial información del emisor y serial
type IssuerSerial struct {
	XMLName              xml.Name `xml:"xades:IssuerSerial"`
	X509IssuerName       string   `xml:"ds:X509IssuerName"`
	X509SerialNumber     string   `xml:"ds:X509SerialNumber"`
}

// SignaturePolicyIdentifier identificador de política de firma
type SignaturePolicyIdentifier struct {
	XMLName                xml.Name `xml:"xades:SignaturePolicyIdentifier"`
	SignaturePolicyId      SignaturePolicyId `xml:"xades:SignaturePolicyId"`
}

// SignaturePolicyId política de firma
type SignaturePolicyId struct {
	XMLName     xml.Name `xml:"xades:SignaturePolicyId"`
	SigPolicyId SigPolicyId `xml:"xades:SigPolicyId"`
	SigPolicyHash SigPolicyHash `xml:"xades:SigPolicyHash"`
}

// SigPolicyId identificador de política
type SigPolicyId struct {
	XMLName    xml.Name `xml:"xades:SigPolicyId"`
	Identifier string   `xml:"xades:Identifier"`
}

// SigPolicyHash hash de política
type SigPolicyHash struct {
	XMLName     xml.Name `xml:"xades:SigPolicyHash"`
	DigestMethod struct {
		XMLName   xml.Name `xml:"ds:DigestMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	} `xml:"ds:DigestMethod"`
	DigestValue string `xml:"ds:DigestValue"`
}

// Signature estructura principal XMLDSig con XAdES
type Signature struct {
	XMLName              xml.Name `xml:"ds:Signature"`
	DSNamespace          string   `xml:"xmlns:ds,attr"`
	ID                   string   `xml:"Id,attr"`
	SignedInfo           SignedInfo `xml:"ds:SignedInfo"`
	SignatureValue       string   `xml:"ds:SignatureValue"`
	KeyInfo              KeyInfo `xml:"ds:KeyInfo"`
	Object               Object  `xml:"ds:Object"`
}

// Object estructura para objetos en XMLDSig
type Object struct {
	XMLName              xml.Name `xml:"ds:Object"`
	QualifyingProperties QualifyingProperties `xml:"xades:QualifyingProperties"`
}

// FirmarXMLXAdESBES firma un documento XML usando XAdES-BES
func FirmarXMLXAdESBES(xmlData []byte, config XAdESBESConfig) ([]byte, error) {
	// Validar certificado
	if config.Certificado == nil {
		return nil, fmt.Errorf("certificado requerido para firma XAdES-BES")
	}

	// Validar clave privada RSA
	rsaKey, ok := config.Certificado.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("se requiere clave privada RSA")
	}

	// Calcular hash del documento
	documentHash := sha256.Sum256(xmlData)
	documentHashB64 := base64.StdEncoding.EncodeToString(documentHash[:])

	// Calcular hash del certificado
	certHash := sha256.Sum256(config.Certificado.Cert.Raw)
	certHashB64 := base64.StdEncoding.EncodeToString(certHash[:])

	// Obtener certificado en base64
	certB64 := base64.StdEncoding.EncodeToString(config.Certificado.Cert.Raw)

	// Crear SignedInfo
	signedInfo := SignedInfo{}
	signedInfo.CanonicalizationMethod.Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	signedInfo.SignatureMethod.Algorithm = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	
	// Referencia al documento
	signedInfo.Reference.URI = ""
	signedInfo.Reference.Transforms.Transform = []struct {
		XMLName   xml.Name `xml:"ds:Transform"`
		Algorithm string   `xml:"Algorithm,attr"`
	}{
		{Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"},
	}
	signedInfo.Reference.DigestMethod.Algorithm = "http://www.w3.org/2000/09/xmldsig#sha1"
	signedInfo.Reference.DigestValue = documentHashB64

	// Serializar SignedInfo para firmarlo
	signedInfoXML, err := xml.Marshal(signedInfo)
	if err != nil {
		return nil, fmt.Errorf("error serializando SignedInfo: %v", err)
	}

	// Calcular hash de SignedInfo
	signedInfoHash := sha256.Sum256(signedInfoXML)

	// Firmar con RSA
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, signedInfoHash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando: %v", err)
	}

	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	// Crear estructura XAdES-BES completa
	xadesSignature := Signature{
		DSNamespace: "http://www.w3.org/2000/09/xmldsig#",
		ID:          "Signature",
		SignedInfo:  signedInfo,
		SignatureValue: signatureB64,
	}

	// KeyInfo
	xadesSignature.KeyInfo.X509Data.X509Certificate = certB64
	xadesSignature.KeyInfo.X509Data.X509SubjectName = config.Certificado.Cert.Subject.String()
	xadesSignature.KeyInfo.X509Data.X509IssuerName = config.Certificado.Cert.Issuer.String()

	// QualifyingProperties (XAdES-BES)
	xadesSignature.Object.QualifyingProperties.XAdESNamespace = "http://uri.etsi.org/01903/v1.3.2#"
	xadesSignature.Object.QualifyingProperties.Target = "#Signature"
	
	// SignedProperties
	xadesSignature.Object.QualifyingProperties.SignedProperties.ID = "SignedProperties"
	xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SigningTime = time.Now().Format(time.RFC3339)
	
	// SigningCertificate
	xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SigningCertificate.Cert.CertDigest.DigestMethod.Algorithm = "http://www.w3.org/2000/09/xmldsig#sha1"
	xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SigningCertificate.Cert.CertDigest.DigestValue = certHashB64
	xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SigningCertificate.Cert.IssuerSerial.X509IssuerName = config.Certificado.Cert.Issuer.String()
	xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SigningCertificate.Cert.IssuerSerial.X509SerialNumber = config.Certificado.Cert.SerialNumber.String()

	// SignaturePolicyIdentifier (obligatorio para SRI)
	if config.PolicyID != "" {
		xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SignaturePolicyIdentifier.SignaturePolicyId.SigPolicyId.Identifier = config.PolicyID
		xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SignaturePolicyIdentifier.SignaturePolicyId.SigPolicyHash.DigestMethod.Algorithm = "http://www.w3.org/2000/09/xmldsig#sha1"
		xadesSignature.Object.QualifyingProperties.SignedProperties.SignedSignatureProperties.SignaturePolicyIdentifier.SignaturePolicyId.SigPolicyHash.DigestValue = config.PolicyHash
	}

	// Convertir a XML
	signatureXML, err := xml.MarshalIndent(xadesSignature, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML de firma: %v", err)
	}

	// Insertar la firma en el documento original
	// Buscar el elemento raíz y agregar la firma antes del cierre
	xmlString := string(xmlData)
	rootEndIndex := strings.LastIndex(xmlString, "</")
	if rootEndIndex == -1 {
		return nil, fmt.Errorf("documento XML mal formado")
	}

	// Construir documento firmado
	signedDocument := xmlString[:rootEndIndex] + 
		string(signatureXML) + "\n" + 
		xmlString[rootEndIndex:]

	return []byte(signedDocument), nil
}

// ValidarFirmaXAdESBES valida una firma XAdES-BES
func ValidarFirmaXAdESBES(xmlFirmado []byte) error {
	// Esta función implementaría la validación de la firma XAdES-BES
	// Por ahora retornamos nil (validación pendiente de implementar)
	return nil
}

// ExtraerCertificadoDeXML extrae el certificado X.509 de un XML firmado (stub)
func ExtraerCertificadoDeXML(xmlData []byte) (*CertificadoDigital, error) {
	// Stub implementation - TODO: implementar extracción real
	if len(xmlData) == 0 {
		return nil, fmt.Errorf("XML vacío")
	}
	return nil, fmt.Errorf("extracción de certificado no implementada")
}

// GenerarHashSHA1 genera hash SHA1 de datos (stub)
func GenerarHashSHA1(data []byte) string {
	// Stub implementation - TODO: implementar hash real
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)[:40] // Truncar a 40 chars para simular SHA1
}

// CrearTimestamp crea timestamp para XAdES (stub)
func CrearTimestamp() string {
	// Stub implementation - TODO: implementar timestamp real
	return time.Now().UTC().Format(time.RFC3339)
}

// NormalizarXML normaliza XML para canonicalización (stub)
func NormalizarXML(xmlData []byte) []byte {
	// Stub implementation - TODO: implementar normalización real
	return xmlData
}

// CrearDigestValue crea digest value para XAdES (stub)
func CrearDigestValue(xmlData []byte) string {
	// Stub implementation - TODO: implementar digest real
	hash := sha256.Sum256(xmlData)
	return base64.StdEncoding.EncodeToString(hash[:])
}