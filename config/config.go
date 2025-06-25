// Package config maneja toda la configuración del sistema
package config

// EmpresaConfig - Configuración de la empresa emisora
type EmpresaConfig struct {
	RazonSocial     string `json:"razonSocial"`
	RUC             string `json:"ruc"`
	Establecimiento string `json:"establecimiento"`
	PuntoEmision    string `json:"puntoEmision"`
	Direccion       string `json:"direccion"`
}

// AmbienteConfig - Configuración por ambiente (desarrollo/producción)
type AmbienteConfig struct {
	Codigo      string `json:"codigo"`      // "1" = pruebas, "2" = producción
	Descripcion string `json:"descripcion"` // "Pruebas" o "Producción"
	TipoEmision string `json:"tipoEmision"` // "1" = normal, "2" = contingencia
}

// CertificadoConfig configuración del certificado digital
type CertificadoConfig struct {
	RutaArchivo string `json:"rutaArchivo"`
	Password    string `json:"password"`
}

// SRIConfig configuración específica del SRI
type SRIConfig struct {
	TimeoutSegundos   int    `json:"timeoutSegundos"`
	MaxReintentos     int    `json:"maxReintentos"`
	PolicyID          string `json:"policyID"`
	PolicyHash        string `json:"policyHash"`
	EndpointRecepcion string `json:"endpointRecepcion"`
	EndpointAutorizacion string `json:"endpointAutorizacion"`
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	Ruta         string `json:"ruta"`
	MaxConexiones int   `json:"maxConexiones"`
}

// FacturacionConfig - Configuración completa del sistema
type FacturacionConfig struct {
	Empresa     EmpresaConfig     `json:"empresa"`
	Ambiente    AmbienteConfig    `json:"ambiente"`
	Certificado CertificadoConfig `json:"certificado"`
	SRI         SRIConfig         `json:"sri"`
	Database    DatabaseConfig    `json:"database"`
}

// Config Global configuration instance
var Config FacturacionConfig

// ContadorSecuencial contador global para secuenciales
var ContadorSecuencial int64 = 1
