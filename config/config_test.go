package config

import (
	"encoding/json"
	"testing"
)

// TestEmpresaConfig_Structure verifica la estructura de configuración de empresa
func TestEmpresaConfig_Structure(t *testing.T) {
	empresa := EmpresaConfig{
		RazonSocial:     "TEST EMPRESA S.A.",
		RUC:             "1234567890001",
		Establecimiento: "001",
		PuntoEmision:    "001", 
		Direccion:       "Av. Test 123",
	}

	// Verificar que todos los campos están presentes
	if empresa.RazonSocial != "TEST EMPRESA S.A." {
		t.Errorf("RazonSocial = %v, quería 'TEST EMPRESA S.A.'", empresa.RazonSocial)
	}
	if empresa.RUC != "1234567890001" {
		t.Errorf("RUC = %v, quería '1234567890001'", empresa.RUC)
	}
	if empresa.Establecimiento != "001" {
		t.Errorf("Establecimiento = %v, quería '001'", empresa.Establecimiento)
	}
	if empresa.PuntoEmision != "001" {
		t.Errorf("PuntoEmision = %v, quería '001'", empresa.PuntoEmision)
	}
	if empresa.Direccion != "Av. Test 123" {
		t.Errorf("Direccion = %v, quería 'Av. Test 123'", empresa.Direccion)
	}
}

// TestAmbienteConfig_Structure verifica la estructura de configuración de ambiente
func TestAmbienteConfig_Structure(t *testing.T) {
	ambiente := AmbienteConfig{
		Codigo:      "1",
		Descripcion: "Ambiente de Pruebas",
		TipoEmision: "1",
	}

	if ambiente.Codigo != "1" {
		t.Errorf("Codigo = %v, quería '1'", ambiente.Codigo)
	}
	if ambiente.Descripcion != "Ambiente de Pruebas" {
		t.Errorf("Descripcion = %v, quería 'Ambiente de Pruebas'", ambiente.Descripcion)
	}
	if ambiente.TipoEmision != "1" {
		t.Errorf("TipoEmision = %v, quería '1'", ambiente.TipoEmision)
	}
}

// TestFacturacionConfig_Structure verifica la estructura completa de configuración
func TestFacturacionConfig_Structure(t *testing.T) {
	config := FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "TEST EMPRESA",
			RUC:         "1234567890001",
		},
		Ambiente: AmbienteConfig{
			Codigo:      "1",
			Descripcion: "Pruebas",
		},
	}

	if config.Empresa.RazonSocial != "TEST EMPRESA" {
		t.Error("Empresa.RazonSocial no se asignó correctamente")
	}
	if config.Ambiente.Codigo != "1" {
		t.Error("Ambiente.Codigo no se asignó correctamente")
	}
}

// TestEmpresaConfig_JSONSerialization verifica serialización JSON
func TestEmpresaConfig_JSONSerialization(t *testing.T) {
	empresa := EmpresaConfig{
		RazonSocial:     "TEST JSON S.A.",
		RUC:             "1234567890001",
		Establecimiento: "002",
		PuntoEmision:    "003",
		Direccion:       "Calle JSON 456",
	}

	// Serializar a JSON
	data, err := json.Marshal(empresa)
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Deserializar de JSON
	var empresaDeserialized EmpresaConfig
	err = json.Unmarshal(data, &empresaDeserialized)
	if err != nil {
		t.Fatalf("Error unmarshaling from JSON: %v", err)
	}

	// Verificar que los datos se mantuvieron
	if empresaDeserialized.RazonSocial != empresa.RazonSocial {
		t.Errorf("RazonSocial después de JSON = %v, quería %v", empresaDeserialized.RazonSocial, empresa.RazonSocial)
	}
	if empresaDeserialized.RUC != empresa.RUC {
		t.Errorf("RUC después de JSON = %v, quería %v", empresaDeserialized.RUC, empresa.RUC)
	}
}

// TestAmbienteConfig_JSONSerialization verifica serialización JSON del ambiente
func TestAmbienteConfig_JSONSerialization(t *testing.T) {
	ambiente := AmbienteConfig{
		Codigo:      "2",
		Descripcion: "Ambiente de Producción",
		TipoEmision: "1",
	}

	// Serializar y deserializar
	data, err := json.Marshal(ambiente)
	if err != nil {
		t.Fatalf("Error marshaling ambiente: %v", err)
	}

	var ambienteDeserialized AmbienteConfig
	err = json.Unmarshal(data, &ambienteDeserialized)
	if err != nil {
		t.Fatalf("Error unmarshaling ambiente: %v", err)
	}

	// Verificar datos
	if ambienteDeserialized.Codigo != ambiente.Codigo {
		t.Errorf("Codigo después de JSON = %v, quería %v", ambienteDeserialized.Codigo, ambiente.Codigo)
	}
	if ambienteDeserialized.Descripcion != ambiente.Descripcion {
		t.Errorf("Descripcion después de JSON = %v, quería %v", ambienteDeserialized.Descripcion, ambiente.Descripcion)
	}
}

// TestFacturacionConfig_CompleteJSONSerialization verifica serialización completa
func TestFacturacionConfig_CompleteJSONSerialization(t *testing.T) {
	config := FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial:     "EMPRESA COMPLETA S.A.",
			RUC:             "9876543210001",
			Establecimiento: "002",
			PuntoEmision:    "003",
			Direccion:       "Av. Completa 789",
		},
		Ambiente: AmbienteConfig{
			Codigo:      "2",
			Descripcion: "Producción",
			TipoEmision: "2",
		},
	}

	// Serializar configuración completa
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Error marshaling config completa: %v", err)
	}

	// Deserializar
	var configDeserialized FacturacionConfig
	err = json.Unmarshal(data, &configDeserialized)
	if err != nil {
		t.Fatalf("Error unmarshaling config completa: %v", err)
	}

	// Verificar empresa
	if configDeserialized.Empresa.RazonSocial != config.Empresa.RazonSocial {
		t.Errorf("Empresa.RazonSocial = %v, quería %v", configDeserialized.Empresa.RazonSocial, config.Empresa.RazonSocial)
	}
	if configDeserialized.Empresa.RUC != config.Empresa.RUC {
		t.Errorf("Empresa.RUC = %v, quería %v", configDeserialized.Empresa.RUC, config.Empresa.RUC)
	}

	// Verificar ambiente
	if configDeserialized.Ambiente.Codigo != config.Ambiente.Codigo {
		t.Errorf("Ambiente.Codigo = %v, quería %v", configDeserialized.Ambiente.Codigo, config.Ambiente.Codigo)
	}
	if configDeserialized.Ambiente.Descripcion != config.Ambiente.Descripcion {
		t.Errorf("Ambiente.Descripcion = %v, quería %v", configDeserialized.Ambiente.Descripcion, config.Ambiente.Descripcion)
	}
}

// TestConfig_GlobalVariable verifica el comportamiento de la variable global
func TestConfig_GlobalVariable(t *testing.T) {
	// Guardar estado original
	originalConfig := Config

	// Modificar configuración global
	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "GLOBAL TEST",
			RUC:         "1111111111111",
		},
		Ambiente: AmbienteConfig{
			Codigo: "1",
		},
	}

	// Verificar que la modificación tomó efecto
	if Config.Empresa.RazonSocial != "GLOBAL TEST" {
		t.Errorf("Config global no se modificó correctamente")
	}

	// Restaurar estado original
	Config = originalConfig
}

// TestEmpresaConfig_EmptyValues verifica comportamiento con valores vacíos
func TestEmpresaConfig_EmptyValues(t *testing.T) {
	empresa := EmpresaConfig{}

	if empresa.RazonSocial != "" {
		t.Errorf("RazonSocial vacía = %v, quería cadena vacía", empresa.RazonSocial)
	}
	if empresa.RUC != "" {
		t.Errorf("RUC vacío = %v, quería cadena vacía", empresa.RUC)
	}
}

// TestAmbienteConfig_ValidCodes verifica códigos válidos de ambiente
func TestAmbienteConfig_ValidCodes(t *testing.T) {
	testCases := []struct {
		name        string
		codigo      string
		descripcion string
	}{
		{
			name:        "ambiente pruebas",
			codigo:      "1",
			descripcion: "Pruebas",
		},
		{
			name:        "ambiente producción",
			codigo:      "2", 
			descripcion: "Producción",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ambiente := AmbienteConfig{
				Codigo:      tc.codigo,
				Descripcion: tc.descripcion,
				TipoEmision: "1",
			}

			if ambiente.Codigo != tc.codigo {
				t.Errorf("Codigo = %v, quería %v", ambiente.Codigo, tc.codigo)
			}
			if ambiente.Descripcion != tc.descripcion {
				t.Errorf("Descripcion = %v, quería %v", ambiente.Descripcion, tc.descripcion)
			}
		})
	}
}

// TestFacturacionConfig_JSONTagsValidation verifica que los tags JSON funcionen
func TestFacturacionConfig_JSONTagsValidation(t *testing.T) {
	config := FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial: "Test Tags",
			RUC:         "1234567890001",
		},
		Ambiente: AmbienteConfig{
			Codigo: "1",
		},
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Error marshaling: %v", err)
	}

	jsonString := string(jsonData)

	// Verificar que los nombres JSON están presentes
	expectedJSONKeys := []string{
		"\"empresa\":",
		"\"ambiente\":",
		"\"razonSocial\":",
		"\"ruc\":",
		"\"codigo\":",
	}

	for _, key := range expectedJSONKeys {
		if !containsSubstring(jsonString, key) {
			t.Errorf("JSON no contiene clave esperada: %s", key)
		}
	}
}

// containsSubstring es una función helper para verificar si una cadena contiene otra
func containsSubstring(str, substr string) bool {
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

// TestConfig_StructCopy verifica que las estructuras se puedan copiar correctamente
func TestConfig_StructCopy(t *testing.T) {
	original := EmpresaConfig{
		RazonSocial:     "ORIGINAL",
		RUC:             "1234567890001",
		Establecimiento: "001",
		PuntoEmision:    "001",
		Direccion:       "Dirección Original",
	}

	// Copiar estructura
	copia := original

	// Modificar copia
	copia.RazonSocial = "COPIA"

	// Verificar que el original no cambió
	if original.RazonSocial != "ORIGINAL" {
		t.Errorf("Original se modificó cuando no debería: %v", original.RazonSocial)
	}
	if copia.RazonSocial != "COPIA" {
		t.Errorf("Copia no se modificó correctamente: %v", copia.RazonSocial)
	}
}

// Benchmark para serialización JSON de EmpresaConfig
func BenchmarkEmpresaConfig_JSONMarshal(b *testing.B) {
	empresa := EmpresaConfig{
		RazonSocial:     "BENCHMARK EMPRESA S.A.",
		RUC:             "1234567890001",
		Establecimiento: "001",
		PuntoEmision:    "001",
		Direccion:       "Av. Benchmark 123",
	}

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(empresa)
		if err != nil {
			b.Fatalf("Error marshaling: %v", err)
		}
	}
}

// Benchmark para deserialización JSON de FacturacionConfig
func BenchmarkFacturacionConfig_JSONUnmarshal(b *testing.B) {
	jsonData := []byte(`{
		"empresa": {
			"razonSocial": "BENCHMARK EMPRESA",
			"ruc": "1234567890001",
			"establecimiento": "001",
			"puntoEmision": "001",
			"direccion": "Dirección benchmark"
		},
		"ambiente": {
			"codigo": "1",
			"descripcion": "Benchmark",
			"tipoEmision": "1"
		}
	}`)

	for i := 0; i < b.N; i++ {
		var config FacturacionConfig
		err := json.Unmarshal(jsonData, &config)
		if err != nil {
			b.Fatalf("Error unmarshaling: %v", err)
		}
	}
}

// Benchmark para asignación de estructuras
func BenchmarkConfig_StructAssignment(b *testing.B) {
	empresa := EmpresaConfig{
		RazonSocial:     "BENCHMARK",
		RUC:             "1234567890001",
		Establecimiento: "001",
		PuntoEmision:    "001",
		Direccion:       "Benchmark Dir",
	}

	for i := 0; i < b.N; i++ {
		copia := empresa
		_ = copia // Evitar optimización del compilador
	}
}