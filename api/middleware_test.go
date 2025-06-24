package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCorsMiddleware verifica el middleware CORS
func TestCorsMiddleware(t *testing.T) {
	server := NewServer("8080")

	// Handler simple para testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   func(*testing.T, http.Header)
	}{
		{
			name:           "GET request con headers CORS",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, headers http.Header) {
				expectedHeaders := map[string]string{
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type, Authorization",
				}

				for header, expectedValue := range expectedHeaders {
					actualValue := headers.Get(header)
					if actualValue != expectedValue {
						t.Errorf("Header %s = %v, quería %v", header, actualValue, expectedValue)
					}
				}
			},
		},
		{
			name:           "POST request con headers CORS",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, headers http.Header) {
				if headers.Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Access-Control-Allow-Origin debe ser '*'")
				}
			},
		},
		{
			name:           "OPTIONS preflight request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, headers http.Header) {
				if headers.Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Access-Control-Allow-Origin debe ser '*' para preflight")
				}
				if headers.Get("Access-Control-Allow-Methods") == "" {
					t.Error("Access-Control-Allow-Methods no debe estar vacío")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Aplicar middleware CORS al handler de test
			handlerWithCORS := server.corsMiddleware(testHandler)

			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			handlerWithCORS.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			if tt.checkHeaders != nil {
				tt.checkHeaders(t, w.Header())
			}

			// Para OPTIONS, el body debe estar vacío (solo headers)
			if tt.method == http.MethodOptions {
				body := w.Body.String()
				if body != "" {
					t.Errorf("OPTIONS response body debe estar vacío, obtuvo: %s", body)
				}
			}
		})
	}
}

// TestLoggingMiddleware verifica el middleware de logging
func TestLoggingMiddleware(t *testing.T) {
	server := NewServer("8080")

	// Handler simple para testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET request logging",
			method:         http.MethodGet,
			path:           "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request logging",
			method:         http.MethodPost,
			path:           "/api/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "request con path largo",
			method:         http.MethodGet,
			path:           "/api/facturas/FAC-000001?includeXML=true",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Aplicar middleware de logging al handler de test
			handlerWithLogging := server.loggingMiddleware(testHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			// Simular remote address
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()

			handlerWithLogging.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			// Verificar que el response es correcto
			if w.Body.String() != "test response" {
				t.Errorf("response body = %v, quería 'test response'", w.Body.String())
			}
		})
	}
}

// TestResponseWriter verifica el wrapper de ResponseWriter
func TestResponseWriter(t *testing.T) {

	tests := []struct {
		name               string
		writeHeader        int
		expectedStatusCode int
	}{
		{
			name:               "status code 200",
			writeHeader:        http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "status code 404",
			writeHeader:        http.StatusNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "status code 500",
			writeHeader:        http.StatusInternalServerError,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear nuevo recorder y wrapper para cada test
			recorder := httptest.NewRecorder()
			wrapper := &responseWriter{
				ResponseWriter: recorder,
				statusCode:     http.StatusOK,
			}

			wrapper.WriteHeader(tt.writeHeader)

			if wrapper.statusCode != tt.expectedStatusCode {
				t.Errorf("statusCode = %v, quería %v", wrapper.statusCode, tt.expectedStatusCode)
			}

			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("ResponseWriter.Code = %v, quería %v", recorder.Code, tt.expectedStatusCode)
			}
		})
	}
}

// TestMiddlewareChain verifica que los middlewares se apliquen en orden correcto
func TestMiddlewareChain(t *testing.T) {
	server := NewServer("8080")

	// Handler que retorna diferentes status codes para testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error response"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok response"))
		}
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkHeaders   bool
		checkLogging   bool
	}{
		{
			name:           "request exitoso con todos los middlewares",
			method:         http.MethodGet,
			path:           "/test",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
			checkLogging:   true,
		},
		{
			name:           "request de error con middlewares",
			method:         http.MethodPost,
			path:           "/error",
			expectedStatus: http.StatusInternalServerError,
			checkHeaders:   true,
			checkLogging:   true,
		},
		{
			name:           "OPTIONS request con middlewares",
			method:         http.MethodOptions,
			path:           "/test",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
			checkLogging:   false, // OPTIONS es interceptado por CORS middleware
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Aplicar toda la cadena de middlewares
			fullHandler := server.middlewareChain(testHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()

			fullHandler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			// Verificar headers CORS si se solicita
			if tt.checkHeaders {
				corsHeader := w.Header().Get("Access-Control-Allow-Origin")
				if corsHeader != "*" {
					t.Errorf("CORS header = %v, quería '*'", corsHeader)
				}
			}

			// Para OPTIONS, verificar que se intercepta correctamente
			if tt.method == http.MethodOptions {
				body := w.Body.String()
				if body != "" {
					t.Error("OPTIONS response debe tener body vacío")
				}
			}
		})
	}
}

// TestWriteJSONResponse verifica la función helper writeJSONResponse
func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		data           interface{}
		expectedStatus int
		checkBody      func(*testing.T, string)
	}{
		{
			name:           "response exitoso con data",
			status:         http.StatusOK,
			data:           map[string]string{"message": "success"},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				if !strings.Contains(body, "success") {
					t.Error("Body debe contener 'success'")
				}
				if !strings.Contains(body, "message") {
					t.Error("Body debe contener 'message'")
				}
			},
		},
		{
			name:           "response con data nil",
			status:         http.StatusNoContent,
			data:           nil,
			expectedStatus: http.StatusNoContent,
			checkBody: func(t *testing.T, body string) {
				if strings.TrimSpace(body) != "" {
					t.Errorf("Body debe estar vacío con data nil, obtuvo: %s", body)
				}
			},
		},
		{
			name:   "response con estructura compleja",
			status: http.StatusCreated,
			data: map[string]interface{}{
				"id":      "FAC-000001",
				"status":  "created",
				"details": map[string]interface{}{"total": 115.0},
			},
			expectedStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				expectedStrings := []string{"FAC-000001", "created", "total", "115"}
				for _, expected := range expectedStrings {
					if !strings.Contains(body, expected) {
						t.Errorf("Body debe contener '%s', body: %s", expected, body)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			writeJSONResponse(w, tt.status, tt.data)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			// Verificar Content-Type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %v, quería 'application/json'", contentType)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.String())
			}
		})
	}
}

// TestWriteErrorResponse verifica la función helper writeErrorResponse
func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		message        string
		expectedStatus int
	}{
		{
			name:           "error 400",
			status:         http.StatusBadRequest,
			message:        "Datos inválidos",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error 404",
			status:         http.StatusNotFound,
			message:        "Recurso no encontrado",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error 500",
			status:         http.StatusInternalServerError,
			message:        "Error interno del servidor",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			writeErrorResponse(w, tt.status, tt.message)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			// Verificar Content-Type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %v, quería 'application/json'", contentType)
			}

			// Verificar estructura del error
			body := w.Body.String()
			
			// Verificar que contiene los elementos básicos del error
			if !strings.Contains(body, "\"error\":true") {
				t.Errorf("Body debe contener 'error':true, body: %s", body)
			}
			if !strings.Contains(body, "\"message\":\""+tt.message+"\"") {
				t.Errorf("Body debe contener el mensaje de error, body: %s", body)
			}
			if !strings.Contains(body, "\"status\":") {
				t.Errorf("Body debe contener status, body: %s", body)
			}
		})
	}
}

// Benchmark para middleware chain completo
func BenchmarkMiddlewareChain(b *testing.B) {
	server := NewServer("8080")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("benchmark response"))
	})

	fullHandler := server.middlewareChain(testHandler)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/benchmark", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		fullHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", w.Code)
		}
	}
}

// Benchmark para CORS middleware solo
func BenchmarkCorsMiddleware(b *testing.B) {
	server := NewServer("8080")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := server.corsMiddleware(testHandler)

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/benchmark", nil)
		w := httptest.NewRecorder()

		corsHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("expected status 200, got %d", w.Code)
		}
	}
}