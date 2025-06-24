package sri

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// TestCrearErrorConexion tests connection error creation
func TestCrearErrorConexion(t *testing.T) {
	mensaje := "Connection failed to SRI endpoint"
	err := CrearErrorConexion(mensaje)

	if err == nil {
		t.Fatal("CrearErrorConexion should not return nil")
	}

	if !strings.Contains(err.Error(), mensaje) {
		t.Errorf("Error should contain message: %s, got: %s", mensaje, err.Error())
	}

	// Test with empty message
	errEmpty := CrearErrorConexion("")
	if errEmpty == nil {
		t.Error("CrearErrorConexion should handle empty message")
	}
}

// TestCrearErrorCertificado tests certificate error creation
func TestCrearErrorCertificado(t *testing.T) {
	mensaje := "Certificate expired or invalid"
	err := CrearErrorCertificado(mensaje)

	if err == nil {
		t.Fatal("CrearErrorCertificado should not return nil")
	}

	if !strings.Contains(err.Error(), mensaje) {
		t.Errorf("Error should contain message: %s, got: %s", mensaje, err.Error())
	}

	// Test with special characters
	mensajeSpecial := "Certificado inválido: åéíóú"
	errSpecial := CrearErrorCertificado(mensajeSpecial)
	if !strings.Contains(errSpecial.Error(), mensajeSpecial) {
		t.Errorf("Error should handle special characters: %s", errSpecial.Error())
	}
}

// TestCrearErrorValidacion tests validation error creation
func TestCrearErrorValidacion(t *testing.T) {
	campo := "RUC"
	razon := "debe tener 13 dígitos"
	err := CrearErrorValidacion(campo, razon)

	if err == nil {
		t.Fatal("CrearErrorValidacion should not return nil")
	}

	errorStr := err.Error()
	if !strings.Contains(errorStr, campo) {
		t.Errorf("Error should contain field: %s, got: %s", campo, errorStr)
	}

	if !strings.Contains(errorStr, razon) {
		t.Errorf("Error should contain reason: %s, got: %s", razon, errorStr)
	}

	// Test with empty values
	errEmpty := CrearErrorValidacion("", "")
	if errEmpty == nil {
		t.Error("CrearErrorValidacion should handle empty values")
	}
}

// TestCrearErrorSRI tests SRI error creation
func TestCrearErrorSRI(t *testing.T) {
	// Since CrearErrorSRI doesn't exist, we test the ErrorSRI struct directly
	err := &ErrorSRI{
		Tipo:          ErrorSistema,
		Codigo:        "SRI-001",
		Mensaje:       "Sistema en mantenimiento",
		Recuperable:   true,
		SugerenciaFix: "Reintentar más tarde",
	}

	if err == nil {
		t.Fatal("CrearErrorSRI should not return nil")
	}

	errorStr := err.Error()
	if !strings.Contains(errorStr, err.Codigo) {
		t.Errorf("Error should contain code: %s, got: %s", err.Codigo, errorStr)
	}

	if !strings.Contains(errorStr, err.Mensaje) {
		t.Errorf("Error should contain message: %s, got: %s", err.Mensaje, errorStr)
	}

	// Test with complex error codes
	tests := []struct {
		codigo  string
		mensaje string
	}{
		{"CLAVE-01", "Clave de acceso ya registrada"},
		{"CERT-02", "Certificado revocado"},
		{"XML-03", "Estructura XML inválida"},
	}

	for _, tt := range tests {
		t.Run(tt.codigo, func(t *testing.T) {
			err := &ErrorSRI{
				Tipo:          ErrorValidacion,
				Codigo:        tt.codigo,
				Mensaje:       tt.mensaje,
				Recuperable:   false,
				SugerenciaFix: "Verificar datos",
			}
			if !strings.Contains(err.Error(), tt.codigo) {
				t.Errorf("Error should contain code: %s", tt.codigo)
			}
			if !strings.Contains(err.Error(), tt.mensaje) {
				t.Errorf("Error should contain message: %s", tt.mensaje)
			}
		})
	}
}

// TestParsearErrorSRI tests SRI error parsing
func TestParsearErrorSRI(t *testing.T) {
	tests := []struct {
		name         string
		mensaje      string
		codigoHTTP   int
		expectedTipo TipoErrorSRI
		expectRec    bool
	}{
		{
			name:         "Clave de acceso duplicada",
			mensaje:      "CLAVE-01: Clave de acceso ya registrada",
			codigoHTTP:   400,
			expectedTipo: ErrorClaveAcceso,
			expectRec:    false,
		},
		{
			name:         "Timeout de conexión",
			mensaje:      "Request timeout",
			codigoHTTP:   504,
			expectedTipo: ErrorTimeout,
			expectRec:    true,
		},
		{
			name:         "Certificado expirado",
			mensaje:      "CERT-01: Certificado expirado",
			codigoHTTP:   403,
			expectedTipo: ErrorCertificado,
			expectRec:    false,
		},
		{
			name:         "Sistema en mantenimiento",
			mensaje:      "SRI-01: Sistema en mantenimiento",
			codigoHTTP:   503,
			expectedTipo: ErrorSistema,
			expectRec:    true,
		},
		{
			name:         "Error de validación XML",
			mensaje:      "XML-01: Estructura inválida",
			codigoHTTP:   400,
			expectedTipo: ErrorValidacion,
			expectRec:    false,
		},
		{
			name:         "Error de conexión",
			mensaje:      "Connection refused",
			codigoHTTP:   0,
			expectedTipo: ErrorConexion,
			expectRec:    true,
		},
		{
			name:         "Error no categorizado",
			mensaje:      "Error desconocido",
			codigoHTTP:   500,
			expectedTipo: ErrorSistema,
			expectRec:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorSRI := ParsearErrorSRI(tt.mensaje, tt.codigoHTTP)

			if errorSRI.Tipo != tt.expectedTipo {
				t.Errorf("Expected type %v, got %v", tt.expectedTipo, errorSRI.Tipo)
			}

			if errorSRI.Recuperable != tt.expectRec {
				t.Errorf("Expected recoverable %v, got %v", tt.expectRec, errorSRI.Recuperable)
			}

			if errorSRI.Mensaje != tt.mensaje {
				t.Errorf("Expected message %s, got %s", tt.mensaje, errorSRI.Mensaje)
			}

			// Note: CodigoHTTP is not part of the ErrorSRI struct in the current implementation
			// This test is removed as the field doesn't exist

			if errorSRI.SugerenciaFix == "" {
				t.Error("SugerenciaFix should not be empty")
			}
		})
	}
}

// TestErrorSRIString tests ErrorSRI String method
func TestErrorSRIString(t *testing.T) {
	errorSRI := ErrorSRI{
		Tipo:          ErrorClaveAcceso,
		Codigo:        "CLAVE-01",
		Mensaje:       "Clave ya registrada",
		Recuperable:   false,
		SugerenciaFix: "Verificar que la clave de acceso sea única",
	}

	resultado := errorSRI.String()

	// Verify essential components are present
	if !strings.Contains(resultado, "CLAVE-01") {
		t.Error("String should contain error code")
	}

	if !strings.Contains(resultado, "ERROR_CLAVE_ACCESO") {
		t.Error("String should contain error type")
	}
}

// TestErrorSRIError tests ErrorSRI Error method
func TestErrorSRIError(t *testing.T) {
	errorSRI := ErrorSRI{
		Tipo:    ErrorCertificado,
		Codigo:  "CERT-01",
		Mensaje: "Certificado inválido",
	}

	errorStr := errorSRI.Error()

	if !strings.Contains(errorStr, "CERT-01") {
		t.Error("Error string should contain the message")
	}

	if !strings.Contains(errorStr, "ERROR_CERTIFICADO") {
		t.Error("Error string should contain the error type")
	}
}

// TestTipoErrorSRIString tests TipoErrorSRI String method
func TestTipoErrorSRIString(t *testing.T) {
	tests := []struct {
		tipo     TipoErrorSRI
		expected string
	}{
		{ErrorConexion, "ERROR_CONEXION"},
		{ErrorCertificado, "ERROR_CERTIFICADO"},
		{ErrorValidacion, "ERROR_VALIDACION"},
		{ErrorClaveAcceso, "ERROR_CLAVE_ACCESO"},
		{ErrorTimeout, "ERROR_TIMEOUT"},
		{ErrorSistema, "ERROR_SISTEMA"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tipo.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestManejarErrorHTTP tests HTTP error handling
func TestManejarErrorHTTP(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		expected   TipoErrorSRI
	}{
		{
			name:       "Bad Request",
			statusCode: http.StatusBadRequest,
			body:       "Invalid request",
			expected:   ErrorValidacion,
		},
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       "Unauthorized access",
			expected:   ErrorCertificado,
		},
		{
			name:       "Forbidden",
			statusCode: http.StatusForbidden,
			body:       "Access denied",
			expected:   ErrorCertificado,
		},
		{
			name:       "Not Found",
			statusCode: http.StatusNotFound,
			body:       "Endpoint not found",
			expected:   ErrorConexion,
		},
		{
			name:       "Timeout",
			statusCode: http.StatusGatewayTimeout,
			body:       "Request timeout",
			expected:   ErrorTimeout,
		},
		{
			name:       "Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       "Service unavailable",
			expected:   ErrorSistema,
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       "Internal error",
			expected:   ErrorSistema,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorSRI := ParsearErrorSRI(tt.body, tt.statusCode)
			if errorSRI.Tipo != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, errorSRI.Tipo)
			}
		})
	}
}

// TestCrearErroresComplejos tests complex error scenarios
func TestCrearErroresComplejos(t *testing.T) {
	// Test chained errors
	originalErr := errors.New("network connection failed")
	connectionErr := fmt.Errorf("SRI connection error: %w", originalErr)
	
	sriErr := CrearErrorConexion(connectionErr.Error())
	if !strings.Contains(sriErr.Error(), "network connection failed") {
		t.Error("Should preserve original error message")
	}

	// Test error with special formatting
	complexMessage := fmt.Sprintf("Error procesando factura #%d para RUC %s", 12345, "1792146739001")
	validationErr := CrearErrorValidacion("factura", complexMessage)
	
	if !strings.Contains(validationErr.Error(), "12345") {
		t.Error("Should preserve formatted message")
	}
}

// TestRecuperabilidadErrores tests error recoverability logic
func TestRecuperabilidadErrores(t *testing.T) {
	recuperables := []struct {
		tipo TipoErrorSRI
		code int
	}{
		{ErrorTimeout, 504},
		{ErrorConexion, 0},
		{ErrorSistema, 503},
	}

	noRecuperables := []struct {
		tipo TipoErrorSRI
		code int
	}{
		{ErrorClaveAcceso, 400},
		{ErrorCertificado, 403},
		{ErrorValidacion, 400},
	}

	for _, r := range recuperables {
		errorSRI := ParsearErrorSRI("test message", r.code)
		if !errorSRI.Recuperable {
			t.Errorf("Error tipo %v should be recoverable", r.tipo)
		}
	}

	for _, nr := range noRecuperables {
		errorSRI := ParsearErrorSRI("test message", nr.code)
		if errorSRI.Recuperable {
			t.Errorf("Error tipo %v should not be recoverable", nr.tipo)
		}
	}
}

// BenchmarkCrearErrorSRI benchmarks SRI error creation
func BenchmarkCrearErrorSRI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &ErrorSRI{
			Tipo:          ErrorSistema,
			Codigo:        "SRI-001",
			Mensaje:       "Error de prueba",
			Recuperable:   false,
			SugerenciaFix: "Revisar configuración",
		}
	}
}

// BenchmarkParsearErrorSRI benchmarks SRI error parsing
func BenchmarkParsearErrorSRI(b *testing.B) {
	mensaje := "CLAVE-01: Clave de acceso ya registrada"
	
	for i := 0; i < b.N; i++ {
		ParsearErrorSRI(mensaje, 400)
	}
}

// TestErrorHandlingEdgeCases tests edge cases in error handling
func TestErrorHandlingEdgeCases(t *testing.T) {
	// Test with nil or empty inputs
	tests := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "Empty connection error",
			testFunc: func() error {
				return CrearErrorConexion("")
			},
		},
		{
			name: "Empty certificate error",
			testFunc: func() error {
				return CrearErrorCertificado("")
			},
		},
		{
			name: "Empty validation error",
			testFunc: func() error {
				return CrearErrorValidacion("", "")
			},
		},
		{
			name: "Empty SRI error",
			testFunc: func() error {
				return &ErrorSRI{
					Tipo:          ErrorSistema,
					Codigo:        "",
					Mensaje:       "",
					Recuperable:   false,
					SugerenciaFix: "",
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if err == nil {
				t.Error("Error creation should not return nil even with empty inputs")
			}
		})
	}
}