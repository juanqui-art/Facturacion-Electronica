package sri

import (
	"errors"
	"testing"
	"time"
)

// TestConfigReintentoDefault tests default retry configuration
func TestConfigReintentoDefault(t *testing.T) {
	config := ConfigReintentoDefault

	if config.MaxIntentos <= 0 {
		t.Error("MaxIntentos should be positive")
	}

	if config.TiempoBase <= 0 {
		t.Error("TiempoBase should be positive")
	}

	if config.Multiplicador <= 1.0 {
		t.Error("Multiplicador should be greater than 1.0")
	}

	if config.TiempoMaximo <= config.TiempoBase {
		t.Error("TiempoMaximo should be greater than TiempoBase")
	}

	if config.JitterMaximo < 0 {
		t.Error("JitterMaximo should not be negative")
	}
}

// TestConfigReintentoConservador tests conservative retry configuration
func TestConfigReintentoConservador(t *testing.T) {
	config := ConfigReintentoConservador

	if config.MaxIntentos <= 0 {
		t.Error("MaxIntentos should be positive")
	}

	if config.TiempoBase <= 0 {
		t.Error("TiempoBase should be positive")
	}

	// Conservative should be more cautious than default
	if config.MaxIntentos > ConfigReintentoDefault.MaxIntentos {
		t.Error("Conservative config should have fewer or equal max attempts")
	}
}

// TestConfigReintentoAgresivo tests aggressive retry configuration
func TestConfigReintentoAgresivo(t *testing.T) {
	config := ConfigReintentoAgresivo

	if config.MaxIntentos <= 0 {
		t.Error("MaxIntentos should be positive")
	}

	if config.TiempoBase <= 0 {
		t.Error("TiempoBase should be positive")
	}

	// Aggressive should try more than default
	if config.MaxIntentos < ConfigReintentoDefault.MaxIntentos {
		t.Error("Aggressive config should have more or equal max attempts")
	}
}

// TestEjecutarConReintentoExitoso tests successful retry execution
func TestEjecutarConReintentoExitoso(t *testing.T) {
	attempts := 0
	
	// Function that succeeds on third attempt
	testFunc := func() error {
		attempts++
		if attempts < 3 {
			return CrearErrorConexion("temporary failure")
		}
		return nil
	}

	config := ConfigReintento{
		MaxIntentos:      5,
		TiempoBase:       1 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     10 * time.Millisecond,
		JitterMaximo:     1 * time.Millisecond,
		SoloRecuperables: true,
	}

	resultado := EjecutarConReintento(testFunc, config)

	if !resultado.Exitoso {
		t.Error("Retry should succeed eventually")
	}

	if resultado.IntentosRealizados != 3 {
		t.Errorf("Expected 3 attempts, got %d", resultado.IntentosRealizados)
	}

	if len(resultado.Errores) != 2 {
		t.Errorf("Expected 2 errors (from first 2 attempts), got %d", len(resultado.Errores))
	}

	if resultado.TiempoTotal <= 0 {
		t.Error("TiempoTotal should be positive")
	}
}

// TestEjecutarConReintentoFallido tests failed retry execution
func TestEjecutarConReintentoFallido(t *testing.T) {
	// Function that always fails
	testFunc := func() error {
		return CrearErrorConexion("persistent failure")
	}

	config := ConfigReintento{
		MaxIntentos:      3,
		TiempoBase:       1 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     5 * time.Millisecond,
		JitterMaximo:     1 * time.Millisecond,
		SoloRecuperables: true,
	}

	resultado := EjecutarConReintento(testFunc, config)

	if resultado.Exitoso {
		t.Error("Retry should fail when function always fails")
	}

	if resultado.IntentosRealizados != 3 {
		t.Errorf("Expected 3 attempts, got %d", resultado.IntentosRealizados)
	}

	if len(resultado.Errores) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(resultado.Errores))
	}
}

// TestEjecutarConReintentoErrorNoRecuperable tests non-recoverable error handling
func TestEjecutarConReintentoErrorNoRecuperable(t *testing.T) {
	// Function that returns non-recoverable error
	testFunc := func() error {
		return CrearErrorValidacion("campo", "error no recuperable")
	}

	config := ConfigReintento{
		MaxIntentos:      5,
		TiempoBase:       1 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     10 * time.Millisecond,
		JitterMaximo:     1 * time.Millisecond,
		SoloRecuperables: true,
	}

	resultado := EjecutarConReintento(testFunc, config)

	if resultado.Exitoso {
		t.Error("Retry should fail immediately for non-recoverable error")
	}

	if resultado.IntentosRealizados != 1 {
		t.Errorf("Expected 1 attempt for non-recoverable error, got %d", resultado.IntentosRealizados)
	}

	if len(resultado.Errores) != 1 {
		t.Errorf("Expected 1 error, got %d", len(resultado.Errores))
	}
}

// TestEjecutarConReintentoSinFiltroRecuperables tests retry without recoverable filter
func TestEjecutarConReintentoSinFiltroRecuperables(t *testing.T) {
	attempts := 0
	
	// Function that returns non-recoverable error but should retry anyway
	testFunc := func() error {
		attempts++
		if attempts < 3 {
			return CrearErrorValidacion("campo", "error no recuperable")
		}
		return nil
	}

	config := ConfigReintento{
		MaxIntentos:      5,
		TiempoBase:       1 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     10 * time.Millisecond,
		JitterMaximo:     1 * time.Millisecond,
		SoloRecuperables: false, // Allow retry of non-recoverable errors
	}

	resultado := EjecutarConReintento(testFunc, config)

	if !resultado.Exitoso {
		t.Error("Retry should succeed when not filtering by recoverability")
	}

	if resultado.IntentosRealizados != 3 {
		t.Errorf("Expected 3 attempts, got %d", resultado.IntentosRealizados)
	}
}

// TestCalcularTiempoEspera tests wait time calculation
func TestCalcularTiempoEspera(t *testing.T) {
	config := ConfigReintento{
		TiempoBase:    100 * time.Millisecond,
		Multiplicador: 2.0,
		TiempoMaximo:  1 * time.Second,
		JitterMaximo:  50 * time.Millisecond,
	}

	tests := []struct {
		intento  int
		minTime  time.Duration
		maxTime  time.Duration
	}{
		{1, 50 * time.Millisecond, 150 * time.Millisecond},   // 100ms base ± 50ms jitter
		{2, 150 * time.Millisecond, 250 * time.Millisecond}, // 200ms ± 50ms jitter
		{3, 350 * time.Millisecond, 450 * time.Millisecond}, // 400ms ± 50ms jitter
	}

	for _, tt := range tests {
		t.Run("Intento_"+string(rune(tt.intento+48)), func(t *testing.T) {
			// Calculate wait time multiple times to account for jitter randomness
			for i := 0; i < 10; i++ {
				waitTime := calcularTiempoEspera(tt.intento, config)
				
				if waitTime < tt.minTime || waitTime > tt.maxTime {
					t.Errorf("Wait time %v not in expected range [%v, %v] for attempt %d",
						waitTime, tt.minTime, tt.maxTime, tt.intento)
				}
			}
		})
	}
}

// TestCalcularTiempoEsperaTiempoMaximo tests maximum wait time enforcement
func TestCalcularTiempoEsperaTiempoMaximo(t *testing.T) {
	config := ConfigReintento{
		TiempoBase:    1 * time.Second,
		Multiplicador: 10.0, // Very high multiplier
		TiempoMaximo:  2 * time.Second,
		JitterMaximo:  0, // No jitter for predictable testing
	}

	// High attempt number should be capped by TiempoMaximo
	waitTime := calcularTiempoEspera(10, config)
	
	if waitTime > config.TiempoMaximo {
		t.Errorf("Wait time %v should not exceed maximum %v", waitTime, config.TiempoMaximo)
	}
}

// TestEsErrorRecuperable tests error recoverability detection
func TestEsErrorRecuperable(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		recuperable bool
	}{
		{
			name:        "Connection error (recoverable)",
			err:         CrearErrorConexion("network failure"),
			recuperable: true,
		},
		{
			name:        "Timeout error (recoverable)",
			err:         ParsearErrorSRI("timeout", 504),
			recuperable: true,
		},
		{
			name:        "System error (recoverable)",
			err:         ParsearErrorSRI("maintenance", 503),
			recuperable: true,
		},
		{
			name:        "Validation error (not recoverable)",
			err:         CrearErrorValidacion("campo", "invalid"),
			recuperable: false,
		},
		{
			name:        "Certificate error (not recoverable)",
			err:         CrearErrorCertificado("expired"),
			recuperable: false,
		},
		{
			name:        "Access key error (not recoverable)",
			err:         ParsearErrorSRI("CLAVE-01: duplicate", 400),
			recuperable: false,
		},
		{
			name:        "Generic error (not recoverable by default)",
			err:         errors.New("generic error"),
			recuperable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EsErrorRecuperable(tt.err)
			if result != tt.recuperable {
				t.Errorf("Expected recoverability %v for error %v, got %v",
					tt.recuperable, tt.err, result)
			}
		})
	}
}

// TestResultadoReintentoString tests ResultadoReintento String method
func TestResultadoReintentoString(t *testing.T) {
	resultado := ResultadoReintento{
		Exitoso:            true,
		IntentosRealizados: 3,
		TiempoTotal:        150 * time.Millisecond,
		Errores: []error{
			CrearErrorConexion("error 1"),
			CrearErrorConexion("error 2"),
		},
	}

	str := resultado.String()

	if str == "" {
		t.Error("String representation should not be empty")
	}

	// Should contain key information
	if !containsSubstr(str, "Exitoso: true") {
		t.Error("Should contain success status")
	}

	if !containsSubstr(str, "Intentos: 3") {
		t.Error("Should contain attempt count")
	}
}

// Helper function to check if string contains substring
func containsSubstr(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr) >= 0)
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestReintentoConTimeout tests retry with timeout scenarios
func TestReintentoConTimeout(t *testing.T) {
	// Function that takes time but eventually succeeds
	testFunc := func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	config := ConfigReintento{
		MaxIntentos:      3,
		TiempoBase:       10 * time.Millisecond,
		Multiplicador:    2.0,
		TiempoMaximo:     100 * time.Millisecond,
		JitterMaximo:     5 * time.Millisecond,
		SoloRecuperables: true,
	}

	start := time.Now()
	resultado := EjecutarConReintento(testFunc, config)
	elapsed := time.Since(start)

	if !resultado.Exitoso {
		t.Error("Function should succeed")
	}

	if elapsed < 50*time.Millisecond {
		t.Error("Should take at least as long as the function execution time")
	}

	if resultado.TiempoTotal <= 0 {
		t.Error("TiempoTotal should be positive")
	}
}

// BenchmarkEjecutarConReintento benchmarks retry execution
func BenchmarkEjecutarConReintento(b *testing.B) {
	config := ConfigReintento{
		MaxIntentos:      3,
		TiempoBase:       1 * time.Microsecond, // Very short for benchmarking
		Multiplicador:    2.0,
		TiempoMaximo:     10 * time.Microsecond,
		JitterMaximo:     0, // No jitter for consistent benchmarking
		SoloRecuperables: true,
	}

	// Function that succeeds immediately
	successFunc := func() error {
		return nil
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		EjecutarConReintento(successFunc, config)
	}
}

// BenchmarkCalcularTiempoEspera benchmarks wait time calculation
func BenchmarkCalcularTiempoEspera(b *testing.B) {
	config := ConfigReintentoDefault
	
	for i := 0; i < b.N; i++ {
		calcularTiempoEspera(3, config)
	}
}

// TestEdgeCasesReintento tests edge cases for retry logic
func TestEdgeCasesReintento(t *testing.T) {
	tests := []struct {
		name   string
		config ConfigReintento
		valid  bool
	}{
		{
			name: "Zero max attempts",
			config: ConfigReintento{
				MaxIntentos: 0,
				TiempoBase:  100 * time.Millisecond,
			},
			valid: false,
		},
		{
			name: "Negative time base",
			config: ConfigReintento{
				MaxIntentos: 3,
				TiempoBase:  -100 * time.Millisecond,
			},
			valid: false,
		},
		{
			name: "Zero multiplier",
			config: ConfigReintento{
				MaxIntentos:   3,
				TiempoBase:    100 * time.Millisecond,
				Multiplicador: 0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFunc := func() error {
				return CrearErrorConexion("test")
			}

			resultado := EjecutarConReintento(testFunc, tt.config)

			if tt.valid && resultado.IntentosRealizados == 0 {
				t.Error("Valid config should result in at least one attempt")
			}

			if !tt.valid && resultado.Exitoso {
				t.Error("Invalid config should not succeed")
			}
		})
	}
}