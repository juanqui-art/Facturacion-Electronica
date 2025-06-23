package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// CargarConfiguracion - Carga la configuración desde archivos JSON
func CargarConfiguracion(archivoConfig string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(archivoConfig)
	if err != nil {
		return fmt.Errorf("error leyendo archivo de configuración %s: %v", archivoConfig, err)
	}
	
	// Parsear JSON a nuestra estructura
	err = json.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("error parseando JSON de configuración: %v", err)
	}
	
	// Validar que la configuración esté completa
	if err := validarConfiguracion(); err != nil {
		return fmt.Errorf("configuración inválida: %v", err)
	}
	
	return nil
}

// validarConfiguracion - Valida que todos los campos requeridos estén presentes
func validarConfiguracion() error {
	if Config.Empresa.RazonSocial == "" {
		return fmt.Errorf("razón social de la empresa es requerida")
	}
	
	if Config.Empresa.RUC == "" {
		return fmt.Errorf("RUC de la empresa es requerido")
	}
	
	if len(Config.Empresa.RUC) != 13 {
		return fmt.Errorf("RUC debe tener exactamente 13 dígitos, tiene %d", len(Config.Empresa.RUC))
	}
	
	if Config.Ambiente.Codigo == "" {
		return fmt.Errorf("código de ambiente es requerido")
	}
	
	if Config.Ambiente.Codigo != "1" && Config.Ambiente.Codigo != "2" {
		return fmt.Errorf("código de ambiente debe ser '1' (pruebas) o '2' (producción)")
	}
	
	return nil
}

// CargarConfiguracionPorDefecto - Carga configuración de desarrollo si no existe archivo
func CargarConfiguracionPorDefecto() {
	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Direccion:       "Av. Amazonas y Naciones Unidas",
		},
		Ambiente: AmbienteConfig{
			Codigo:      "1", // Pruebas
			Descripcion: "Ambiente de Pruebas",
			TipoEmision: "1", // Normal
		},
	}
}

// ObtenerSecuencialSiguiente - Genera el próximo secuencial (por ahora simple)
func ObtenerSecuencialSiguiente() string {
	// En el futuro esto vendría de una base de datos
	// Por ahora devolvemos un secuencial fijo
	return "000000001"
}

// GenerarClaveAcceso - Genera clave de acceso (por ahora placeholder)
func GenerarClaveAcceso() string {
	// En el futuro implementaremos el algoritmo real del SRI
	return "2025062201123456789000110010010000000011234567890"
}