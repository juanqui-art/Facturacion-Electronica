package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCargarConfiguracionPorDefecto verifica la configuración por defecto
func TestCargarConfiguracionPorDefecto(t *testing.T) {
	// Guardar configuración actual
	originalConfig := Config

	// Cargar configuración por defecto
	CargarConfiguracionPorDefecto()

	// Verificar empresa por defecto
	if Config.Empresa.RazonSocial != "EMPRESA DEMO S.A." {
		t.Errorf("RazonSocial por defecto = %v, quería 'EMPRESA DEMO S.A.'", Config.Empresa.RazonSocial)
	}
	if Config.Empresa.RUC != "1234567890001" {
		t.Errorf("RUC por defecto = %v, quería '1234567890001'", Config.Empresa.RUC)
	}
	if Config.Empresa.Establecimiento != "001" {
		t.Errorf("Establecimiento por defecto = %v, quería '001'", Config.Empresa.Establecimiento)
	}
	if Config.Empresa.PuntoEmision != "001" {
		t.Errorf("PuntoEmision por defecto = %v, quería '001'", Config.Empresa.PuntoEmision)
	}

	// Verificar ambiente por defecto
	if Config.Ambiente.Codigo != "1" {
		t.Errorf("Ambiente.Codigo por defecto = %v, quería '1'", Config.Ambiente.Codigo)
	}
	if Config.Ambiente.Descripcion != "Ambiente de Pruebas" {
		t.Errorf("Ambiente.Descripcion por defecto = %v, quería 'Ambiente de Pruebas'", Config.Ambiente.Descripcion)
	}
	if Config.Ambiente.TipoEmision != "1" {
		t.Errorf("Ambiente.TipoEmision por defecto = %v, quería '1'", Config.Ambiente.TipoEmision)
	}

	// Restaurar configuración original
	Config = originalConfig
}

// TestCargarConfiguracion_ArchivoValido verifica carga desde archivo válido
func TestCargarConfiguracion_ArchivoValido(t *testing.T) {
	// Crear archivo temporal de configuración
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.json")

	configJSON := `{
		"empresa": {
			"razonSocial": "EMPRESA TEST S.A.",
			"ruc": "9876543210001",
			"establecimiento": "002",
			"puntoEmision": "003",
			"direccion": "Av. Test 123"
		},
		"ambiente": {
			"codigo": "2",
			"descripcion": "Ambiente de Producción",
			"tipoEmision": "1"
		}
	}`

	err := os.WriteFile(configFile, []byte(configJSON), 0644)
	if err != nil {
		t.Fatalf("Error creando archivo de configuración de prueba: %v", err)
	}

	// Guardar configuración actual
	originalConfig := Config

	// Cargar configuración desde archivo
	err = CargarConfiguracion(configFile)
	if err != nil {
		t.Fatalf("Error cargando configuración: %v", err)
	}

	// Verificar que la configuración se cargó correctamente
	if Config.Empresa.RazonSocial != "EMPRESA TEST S.A." {
		t.Errorf("RazonSocial = %v, quería 'EMPRESA TEST S.A.'", Config.Empresa.RazonSocial)
	}
	if Config.Empresa.RUC != "9876543210001" {
		t.Errorf("RUC = %v, quería '9876543210001'", Config.Empresa.RUC)
	}
	if Config.Ambiente.Codigo != "2" {
		t.Errorf("Ambiente.Codigo = %v, quería '2'", Config.Ambiente.Codigo)
	}

	// Restaurar configuración original
	Config = originalConfig
}

// TestCargarConfiguracion_ArchivoInexistente verifica error con archivo inexistente
func TestCargarConfiguracion_ArchivoInexistente(t *testing.T) {
	err := CargarConfiguracion("/archivo/que/no/existe.json")
	if err == nil {
		t.Error("CargarConfiguracion() debería retornar error con archivo inexistente")
	}

	// Verificar que el mensaje de error es apropiado
	expectedSubstring := "error leyendo archivo de configuración"
	if !contains(err.Error(), expectedSubstring) {
		t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
	}
}

// TestCargarConfiguracion_JSONInvalido verifica error con JSON inválido
func TestCargarConfiguracion_JSONInvalido(t *testing.T) {
	// Crear archivo temporal con JSON inválido
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid_config.json")

	invalidJSON := `{
		"empresa": {
			"razonSocial": "TEST",
			"ruc": "1234567890001"
		// JSON inválido - falta cerrar llaves
	`

	err := os.WriteFile(configFile, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Error creando archivo con JSON inválido: %v", err)
	}

	// Intentar cargar configuración
	err = CargarConfiguracion(configFile)
	if err == nil {
		t.Error("CargarConfiguracion() debería retornar error con JSON inválido")
	}

	// Verificar tipo de error
	expectedSubstring := "error parseando JSON"
	if !contains(err.Error(), expectedSubstring) {
		t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
	}
}

// TestValidarConfiguracion_ConfiguracionCompleta verifica validación exitosa
func TestValidarConfiguracion_ConfiguracionCompleta(t *testing.T) {
	// Guardar configuración actual
	originalConfig := Config

	// Configurar datos válidos
	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial:     "EMPRESA VALIDA S.A.",
			RUC:             "1234567890001",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Direccion:       "Dirección válida",
		},
		Ambiente: AmbienteConfig{
			Codigo:      "1",
			Descripcion: "Pruebas",
			TipoEmision: "1",
		},
	}

	// Validar configuración
	err := validarConfiguracion()
	if err != nil {
		t.Errorf("validarConfiguracion() retornó error con configuración válida: %v", err)
	}

	// Restaurar configuración original
	Config = originalConfig
}

// TestValidarConfiguracion_RazonSocialVacia verifica error con razón social vacía
func TestValidarConfiguracion_RazonSocialVacia(t *testing.T) {
	originalConfig := Config

	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "", // Vacía
			RUC:         "1234567890001",
		},
		Ambiente: AmbienteConfig{
			Codigo: "1",
		},
	}

	err := validarConfiguracion()
	if err == nil {
		t.Error("validarConfiguracion() debería retornar error con razón social vacía")
	}

	expectedSubstring := "razón social de la empresa es requerida"
	if !contains(err.Error(), expectedSubstring) {
		t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
	}

	Config = originalConfig
}

// TestValidarConfiguracion_RUCVacio verifica error con RUC vacío
func TestValidarConfiguracion_RUCVacio(t *testing.T) {
	originalConfig := Config

	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "EMPRESA TEST",
			RUC:         "", // Vacío
		},
		Ambiente: AmbienteConfig{
			Codigo: "1",
		},
	}

	err := validarConfiguracion()
	if err == nil {
		t.Error("validarConfiguracion() debería retornar error con RUC vacío")
	}

	expectedSubstring := "RUC de la empresa es requerido"
	if !contains(err.Error(), expectedSubstring) {
		t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
	}

	Config = originalConfig
}

// TestValidarConfiguracion_RUCLongitudIncorrecta verifica error con RUC de longitud incorrecta
func TestValidarConfiguracion_RUCLongitudIncorrecta(t *testing.T) {
	originalConfig := Config

	testCases := []struct {
		name string
		ruc  string
	}{
		{"RUC muy corto", "123456789"},
		{"RUC muy largo", "12345678901234"},
		{"RUC de 12 dígitos", "123456789012"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Config = FacturacionConfig{
				Empresa: EmpresaConfig{
					RazonSocial: "EMPRESA TEST",
					RUC:         tc.ruc,
				},
				Ambiente: AmbienteConfig{
					Codigo: "1",
				},
			}

			err := validarConfiguracion()
			if err == nil {
				t.Errorf("validarConfiguracion() debería retornar error con RUC '%s'", tc.ruc)
			}

			expectedSubstring := "RUC debe tener exactamente 13 dígitos"
			if !contains(err.Error(), expectedSubstring) {
				t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
			}
		})
	}

	Config = originalConfig
}

// TestValidarConfiguracion_CodigoAmbienteVacio verifica error con código de ambiente vacío
func TestValidarConfiguracion_CodigoAmbienteVacio(t *testing.T) {
	originalConfig := Config

	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "EMPRESA TEST",
			RUC:         "1234567890001",
		},
		Ambiente: AmbienteConfig{
			Codigo: "", // Vacío
		},
	}

	err := validarConfiguracion()
	if err == nil {
		t.Error("validarConfiguracion() debería retornar error con código de ambiente vacío")
	}

	expectedSubstring := "código de ambiente es requerido"
	if !contains(err.Error(), expectedSubstring) {
		t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
	}

	Config = originalConfig
}

// TestValidarConfiguracion_CodigoAmbienteInvalido verifica error con código de ambiente inválido
func TestValidarConfiguracion_CodigoAmbienteInvalido(t *testing.T) {
	originalConfig := Config

	codigosInvalidos := []string{"0", "3", "9", "A", "pruebas", "produccion"}

	for _, codigo := range codigosInvalidos {
		t.Run("codigo_"+codigo, func(t *testing.T) {
			Config = FacturacionConfig{
				Empresa: EmpresaConfig{
					RazonSocial: "EMPRESA TEST",
					RUC:         "1234567890001",
				},
				Ambiente: AmbienteConfig{
					Codigo: codigo,
				},
			}

			err := validarConfiguracion()
			if err == nil {
				t.Errorf("validarConfiguracion() debería retornar error con código '%s'", codigo)
			}

			expectedSubstring := "código de ambiente debe ser '1' (pruebas) o '2' (producción)"
			if !contains(err.Error(), expectedSubstring) {
				t.Errorf("Error message = %v, debería contener '%v'", err.Error(), expectedSubstring)
			}
		})
	}

	Config = originalConfig
}

// TestObtenerSecuencialSiguiente verifica la función de secuencial
func TestObtenerSecuencialSiguiente(t *testing.T) {
	secuencial := ObtenerSecuencialSiguiente()

	if secuencial != "000000001" {
		t.Errorf("ObtenerSecuencialSiguiente() = %v, quería '000000001'", secuencial)
	}

	// Verificar que es consistente
	secuencial2 := ObtenerSecuencialSiguiente()
	if secuencial2 != secuencial {
		t.Errorf("ObtenerSecuencialSiguiente() no es consistente: %v != %v", secuencial2, secuencial)
	}
}

// TestGenerarClaveAcceso verifica la función de clave de acceso
func TestGenerarClaveAcceso(t *testing.T) {
	claveAcceso := GenerarClaveAcceso()

	// Verificar que no está vacía
	if claveAcceso == "" {
		t.Error("GenerarClaveAcceso() no debería retornar cadena vacía")
	}

	// Verificar longitud (debería ser 49 dígitos)
	if len(claveAcceso) != 49 {
		t.Errorf("GenerarClaveAcceso() longitud = %d, quería 49", len(claveAcceso))
	}

	// Verificar que es consistente
	claveAcceso2 := GenerarClaveAcceso()
	if claveAcceso2 != claveAcceso {
		t.Errorf("GenerarClaveAcceso() no es consistente: %v != %v", claveAcceso2, claveAcceso)
	}
}

// TestCargarConfiguracion_IntegracionCompleta prueba el flujo completo
func TestCargarConfiguracion_IntegracionCompleta(t *testing.T) {
	// Crear archivo de configuración completo
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "integracion_config.json")

	configJSON := `{
		"empresa": {
			"razonSocial": "EMPRESA INTEGRACION S.A.",
			"ruc": "1111222233001",
			"establecimiento": "003",
			"puntoEmision": "004",
			"direccion": "Av. Integración 999"
		},
		"ambiente": {
			"codigo": "2",
			"descripcion": "Ambiente de Producción Integración",
			"tipoEmision": "2"
		}
	}`

	err := os.WriteFile(configFile, []byte(configJSON), 0644)
	if err != nil {
		t.Fatalf("Error creando archivo de integración: %v", err)
	}

	// Guardar estado original
	originalConfig := Config

	// Cargar y verificar
	err = CargarConfiguracion(configFile)
	if err != nil {
		t.Fatalf("Error en integración completa: %v", err)
	}

	// Verificar todos los campos
	if Config.Empresa.RazonSocial != "EMPRESA INTEGRACION S.A." {
		t.Error("RazonSocial no se cargó correctamente en integración")
	}
	if Config.Empresa.RUC != "1111222233001" {
		t.Error("RUC no se cargó correctamente en integración")
	}
	if Config.Empresa.Establecimiento != "003" {
		t.Error("Establecimiento no se cargó correctamente en integración")
	}
	if Config.Ambiente.Codigo != "2" {
		t.Error("Ambiente.Codigo no se cargó correctamente en integración")
	}

	// Restaurar
	Config = originalConfig
}

// contains función helper mejorada
func contains(str, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(str) < len(substr) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark para CargarConfiguracionPorDefecto
func BenchmarkCargarConfiguracionPorDefecto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CargarConfiguracionPorDefecto()
	}
}

// Benchmark para validarConfiguracion
func BenchmarkValidarConfiguracion(b *testing.B) {
	// Configurar datos válidos
	originalConfig := Config
	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "BENCHMARK EMPRESA",
			RUC:         "1234567890001",
		},
		Ambiente: AmbienteConfig{
			Codigo: "1",
		},
	}

	for i := 0; i < b.N; i++ {
		err := validarConfiguracion()
		if err != nil {
			b.Fatalf("Error en validación durante benchmark: %v", err)
		}
	}

	Config = originalConfig
}

// Benchmark para ObtenerSecuencialSiguiente
func BenchmarkObtenerSecuencialSiguiente(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ObtenerSecuencialSiguiente()
	}
}

// Benchmark para GenerarClaveAcceso
func BenchmarkGenerarClaveAcceso(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerarClaveAcceso()
	}
}