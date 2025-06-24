// Package sri contiene las funcionalidades específicas para integración con el SRI de Ecuador
package sri

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"software.sslmate.com/src/go-pkcs12"
	"time"
)

// CertificadoDigital representa un certificado digital PKCS#12 para firma electrónica
type CertificadoDigital struct {
	Archivo    string              // Ruta al archivo .p12
	Password   string              // Contraseña del certificado
	PrivateKey interface{}         // Clave privada extraída
	Cert       *x509.Certificate   // Certificado X.509
	CACerts    []*x509.Certificate // Certificados de la CA
}

// CertificadoConfig configuración para certificados digitales
type CertificadoConfig struct {
	RutaArchivo     string `json:"rutaArchivo"`
	Password        string `json:"password"`
	ValidarVigencia bool   `json:"validarVigencia"`
	ValidarCadena   bool   `json:"validarCadena"`
}

// CargarCertificado carga un certificado PKCS#12 desde archivo
func CargarCertificado(config CertificadoConfig) (*CertificadoDigital, error) {
	// Leer el archivo .p12
	data, err := ioutil.ReadFile(config.RutaArchivo)
	if err != nil {
		return nil, fmt.Errorf("error leyendo certificado: %v", err)
	}

	// Decodificar PKCS#12
	// IMPORTANTE: En Ecuador, los certificados del Banco Central tienen 2 claves privadas
	// Necesitamos asegurarnos de tomar la correcta (generalmente la segunda)
	privateKey, cert, caCerts, err := pkcs12.DecodeChain(data, config.Password)
	if err != nil {
		return nil, fmt.Errorf("error decodificando PKCS#12: %v", err)
	}

	certificado := &CertificadoDigital{
		Archivo:    config.RutaArchivo,
		Password:   config.Password,
		PrivateKey: privateKey,
		Cert:       cert,
		CACerts:    caCerts,
	}

	// Validaciones opcionales
	if config.ValidarVigencia {
		if err := certificado.ValidarVigencia(); err != nil {
			return nil, fmt.Errorf("certificado no válido: %v", err)
		}
	}

	return certificado, nil
}

// ValidarVigencia verifica que el certificado esté vigente
func (cd *CertificadoDigital) ValidarVigencia() error {
	// Verificar que el certificado no haya expirado
	if cd.Cert.NotAfter.Before(time.Now()) {
		return fmt.Errorf("certificado expirado el %v", cd.Cert.NotAfter)
	}

	// Verificar que el certificado ya esté vigente
	if cd.Cert.NotBefore.After(time.Now()) {
		return fmt.Errorf("certificado no vigente hasta %v", cd.Cert.NotBefore)
	}

	return nil
}

// ObtenerSubject obtiene el subject del certificado
func (cd *CertificadoDigital) ObtenerSubject() string {
	return cd.Cert.Subject.CommonName
}

// ObtenerIssuer obtiene el emisor del certificado
func (cd *CertificadoDigital) ObtenerIssuer() string {
	return cd.Cert.Issuer.CommonName
}

// ObtenerSerialNumber obtiene el número de serie del certificado
func (cd *CertificadoDigital) ObtenerSerialNumber() string {
	return cd.Cert.SerialNumber.String()
}

// ExportarClavePEM exporta la clave privada en formato PEM
// Útil para trabajar con librerías que requieren PEM en lugar de PKCS#12
func (cd *CertificadoDigital) ExportarClavePEM() ([]byte, error) {
	// Convertir clave privada a bytes DER
	derBytes, err := x509.MarshalPKCS8PrivateKey(cd.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error marshalling clave privada: %v", err)
	}

	// Crear bloque PEM
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// ExportarCertificadoPEM exporta el certificado en formato PEM
func (cd *CertificadoDigital) ExportarCertificadoPEM() []byte {
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cd.Cert.Raw,
	}

	return pem.EncodeToMemory(pemBlock)
}

// MostrarInformacion muestra información detallada del certificado
func (cd *CertificadoDigital) MostrarInformacion() {
	fmt.Println("\n🔐 INFORMACIÓN DEL CERTIFICADO DIGITAL")
	fmt.Println("=====================================")
	fmt.Printf("📋 Subject: %s\n", cd.ObtenerSubject())
	fmt.Printf("🏢 Emisor: %s\n", cd.ObtenerIssuer())
	fmt.Printf("🔢 Número de Serie: %s\n", cd.ObtenerSerialNumber())
	fmt.Printf("📅 Válido desde: %v\n", cd.Cert.NotBefore)
	fmt.Printf("📅 Válido hasta: %v\n", cd.Cert.NotAfter)
	fmt.Printf("🔧 Algoritmo: %v\n", cd.Cert.SignatureAlgorithm)
	fmt.Printf("📁 Archivo: %s\n", cd.Archivo)

	if len(cd.CACerts) > 0 {
		fmt.Printf("\n🔗 Certificados de CA en cadena: %d\n", len(cd.CACerts))
		for i, caCert := range cd.CACerts {
			fmt.Printf("   %d. %s\n", i+1, caCert.Subject.CommonName)
		}
	}
}
