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

// FacturacionConfig - Configuración completa del sistema
type FacturacionConfig struct {
	Empresa  EmpresaConfig  `json:"empresa"`
	Ambiente AmbienteConfig `json:"ambiente"`
}

// Config Global configuration instance
var Config FacturacionConfig
