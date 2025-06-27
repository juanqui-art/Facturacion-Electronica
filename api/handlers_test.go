package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-facturacion-sri/config"
	"go-facturacion-sri/models"
)

// setUp inicializa configuración para tests
func setUp() {
	config.CargarConfiguracionPorDefecto()
	// Crear nuevo storage para cada test
	storage = NewFacturaStorage()
}

// TestNewServer verifica la creación del servidor
func TestNewServer(t *testing.T) {
	server := NewServer("8080")
	
	if server == nil {
		t.Fatal("NewServer() retornó nil")
	}
	if server.port != "8080" {
		t.Errorf("port = %v, quería '8080'", server.port)
	}
	if server.router == nil {
		t.Error("router no debe ser nil")
	}
}

// TestHandleHealth verifica el endpoint de health check
func TestHandleHealth(t *testing.T) {
	setUp()
	server := NewServer("8080")

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "GET health check válido",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["status"] != "healthy" {
					t.Errorf("status = %v, quería 'healthy'", response["status"])
				}
				if response["service"] != "SRI Facturación Electrónica API" {
					t.Errorf("service = %v, quería 'SRI Facturación Electrónica API'", response["service"])
				}
				if response["version"] != "1.0.0" {
					t.Errorf("version = %v, quería '1.0.0'", response["version"])
				}
			},
		},
		{
			name:           "POST no permitido",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			server.handleHealth(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("error unmarshaling response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestHandleRoot verifica el endpoint de documentación
func TestHandleRoot(t *testing.T) {
	setUp()
	server := NewServer("8080")

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET documentación válida",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST no permitido",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			server.handleRoot(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("error unmarshaling response: %v", err)
				}

				if response["service"] != "SRI Facturación Electrónica API" {
					t.Error("Respuesta debe contener información del servicio")
				}
				if response["endpoints"] == nil {
					t.Error("Respuesta debe contener información de endpoints")
				}
			}
		})
	}
}

// TestHandleCreateFactura verifica la creación de facturas
func TestHandleCreateFactura(t *testing.T) {
	setUp()
	server := NewServer("8080")

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "crear factura válida sin XML",
			requestBody: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Juan Pérez",
					ClienteCedula: "1713175071",
					Productos: []models.ProductoInput{
						{
							Codigo:         "LAPTOP001",
							Descripcion:    "Laptop Dell",
							Cantidad:       1.0,
							PrecioUnitario: 450.00,
						},
					},
				},
				IncludeXML: false,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["id"] == nil {
					t.Error("ID no debe ser nil")
				}
				if response["status"] != "created" {
					t.Errorf("status = %v, quería 'created'", response["status"])
				}
				if response["xml"] != nil && response["xml"] != "" {
					t.Error("XML no debe estar presente cuando IncludeXML es false")
				}
			},
		},
		{
			name: "crear factura válida con XML",
			requestBody: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "María González",
					ClienteCedula: "0926687856",
					Productos: []models.ProductoInput{
						{
							Codigo:         "MOUSE001",
							Descripcion:    "Mouse Inalámbrico",
							Cantidad:       2.0,
							PrecioUnitario: 25.00,
						},
					},
				},
				IncludeXML: true,
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["xml"] == nil || response["xml"] == "" {
					t.Error("XML debe estar presente cuando IncludeXML es true")
				}
				xmlStr, ok := response["xml"].(string)
				if !ok {
					t.Error("XML debe ser string")
				} else if !strings.Contains(xmlStr, "<factura>") {
					t.Error("XML debe contener elemento factura")
				}
			},
		},
		{
			name: "crear factura con cédula inválida",
			requestBody: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Cliente Test",
					ClienteCedula: "123", // cédula inválida
					Productos: []models.ProductoInput{
						{
							Codigo:         "TEST001",
							Descripcion:    "Producto test",
							Cantidad:       1.0,
							PrecioUnitario: 100.00,
						},
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
				message, ok := response["message"].(string)
				if !ok || !strings.Contains(message, "cédula") {
					t.Error("mensaje debe mencionar problema con cédula")
				}
			},
		},
		{
			name: "crear factura sin productos",
			requestBody: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Cliente Test",
					ClienteCedula: "1713175071",
					Productos:     []models.ProductoInput{}, // sin productos
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
		{
			name:           "JSON inválido",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar request body
			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/facturas", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleCreateFactura(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("error unmarshaling response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestHandleListFacturas verifica el listado de facturas
func TestHandleListFacturas(t *testing.T) {
	setUp()
	server := NewServer("8080")

	// Agregar algunas facturas al storage
	storage.Store("FAC-000001", FacturaResponse{
		ID:        "FAC-000001",
		Status:    "created",
		CreatedAt: time.Now(),
	})
	storage.Store("FAC-000002", FacturaResponse{
		ID:        "FAC-000002",
		Status:    "created",
		CreatedAt: time.Now(),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/facturas", nil)
	w := httptest.NewRecorder()

	server.handleListFacturas(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, quería %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response: %v", err)
	}

	if response["total"] != float64(2) { // JSON numbers are float64
		t.Errorf("total = %v, quería 2", response["total"])
	}

	facturas, ok := response["facturas"].([]interface{})
	if !ok {
		t.Fatal("facturas debe ser array")
	}

	if len(facturas) != 2 {
		t.Errorf("número de facturas = %v, quería 2", len(facturas))
	}
}

// TestHandleFacturaByID verifica obtener factura por ID
func TestHandleFacturaByID(t *testing.T) {
	setUp()
	server := NewServer("8080")

	// Agregar factura al storage
	testFactura := FacturaResponse{
		ID:        "FAC-000001",
		Status:    "created",
		CreatedAt: time.Now(),
		Factura: models.Factura{
			InfoTributaria: models.InfoTributaria{
				Ambiente:        "1",
				TipoEmision:     "1",
				RazonSocial:     "INNOVATECH SOLUTIONS CIA. LTDA.",
				RUC:             "1791000005001",
				ClaveAcceso:     "27062025011791000005001100100100000000175853228818",
				CodDoc:          "01",
				Establecimiento: "001",
				PuntoEmision:    "001",
				Secuencial:      "000000001",
			},
			InfoFactura: models.InfoFactura{
				FechaEmision:                "27/06/2025",
				DirEstablecimiento:          "Av. República del Salvador N36-84",
				TipoIdentificacionComprador: "05",
				IdentificacionComprador:     "1713175071",
				RazonSocialComprador:        "Test Cliente",
				TotalSinImpuestos:           100.0,
				TotalDescuento:              0.0,
				ImporteTotal:                115.0,
				Moneda:                      "DOLAR",
			},
			Detalles: []models.Detalle{
				{
					CodigoPrincipal:        "PROD001",
					Descripcion:            "Producto de prueba",
					Cantidad:               1.0,
					PrecioUnitario:         100.0,
					Descuento:              0.0,
					PrecioTotalSinImpuesto: 100.0,
				},
			},
		},
	}
	storage.Store("FAC-000001", testFactura)

	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "obtener factura existente",
			path:           "/api/facturas/FAC-000001",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["id"] != "FAC-000001" {
					t.Errorf("id = %v, quería 'FAC-000001'", response["id"])
				}
			},
		},
		{
			name:           "obtener factura con XML",
			path:           "/api/facturas/FAC-000001?includeXML=true",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["xml"] == nil || response["xml"] == "" {
					t.Error("XML debe estar presente cuando includeXML=true")
				}
			},
		},
		{
			name:           "factura no encontrada",
			path:           "/api/facturas/FAC-999999",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
		{
			name:           "ID vacío",
			path:           "/api/facturas/",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
		{
			name:           "método no permitido",
			path:           "/api/facturas/FAC-000001",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse: func(t *testing.T, response map[string]interface{}) {
				if response["error"] != true {
					t.Error("error debe ser true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			server.handleFacturaByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("error unmarshaling response: %v", err)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, response)
			}
		})
	}
}

// TestHandleFacturas verifica el router principal de facturas
func TestHandleFacturas(t *testing.T) {
	setUp()
	server := NewServer("8080")

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET facturas",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "método no soportado",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			if tt.method == http.MethodPost {
				validRequest := CreateFacturaRequest{
					FacturaInput: models.FacturaInput{
						ClienteNombre: "Test",
						ClienteCedula: "1713175071",
						Productos: []models.ProductoInput{
							{
								Codigo:         "TEST001",
								Descripcion:    "Test",
								Cantidad:       1.0,
								PrecioUnitario: 100.0,
							},
						},
					},
				}
				json.NewEncoder(&body).Encode(validRequest)
			}

			req := httptest.NewRequest(tt.method, "/api/facturas", &body)
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			server.handleFacturas(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("status code = %v, quería %v", w.Code, tt.expectedStatus)
			}
		})
	}
}

// TestCreateFacturaRequest_StructureValidation verifica la estructura de request
func TestCreateFacturaRequest_StructureValidation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateFacturaRequest
	}{
		{
			name: "request completo",
			request: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Test Cliente",
					ClienteCedula: "1713175071",
					Productos: []models.ProductoInput{
						{
							Codigo:         "TEST001",
							Descripcion:    "Test Product",
							Cantidad:       1.0,
							PrecioUnitario: 100.0,
						},
					},
				},
				IncludeXML: true,
			},
		},
		{
			name: "request sin IncludeXML",
			request: CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Test Cliente",
					ClienteCedula: "1713175071",
					Productos: []models.ProductoInput{
						{
							Codigo:         "TEST001",
							Descripcion:    "Test Product",
							Cantidad:       1.0,
							PrecioUnitario: 100.0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verificar que la estructura se puede serializar/deserializar
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("error marshaling request: %v", err)
			}

			var unmarshaled CreateFacturaRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("error unmarshaling request: %v", err)
			}

			if unmarshaled.ClienteNombre != tt.request.ClienteNombre {
				t.Errorf("ClienteNombre = %v, quería %v", unmarshaled.ClienteNombre, tt.request.ClienteNombre)
			}
		})
	}
}

// TestFacturaResponse_StructureValidation verifica la estructura de response
func TestFacturaResponse_StructureValidation(t *testing.T) {
	response := FacturaResponse{
		ID:        "FAC-000001",
		Status:    "created",
		CreatedAt: time.Now(),
		Factura: models.Factura{
			InfoTributaria: models.InfoTributaria{
				Secuencial: "000000001",
			},
		},
		XML: "<factura></factura>",
	}

	// Verificar serialización JSON
	data, err := json.Marshal(response)
	if err != nil {
		t.Errorf("error marshaling response: %v", err)
	}

	var unmarshaled FacturaResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("error unmarshaling response: %v", err)
	}

	if unmarshaled.ID != response.ID {
		t.Errorf("ID = %v, quería %v", unmarshaled.ID, response.ID)
	}
	if unmarshaled.Status != response.Status {
		t.Errorf("Status = %v, quería %v", unmarshaled.Status, response.Status)
	}
}

// Benchmark para crear facturas
func BenchmarkHandleCreateFactura(b *testing.B) {
	setUp()
	server := NewServer("8080")

	request := CreateFacturaRequest{
		FacturaInput: models.FacturaInput{
			ClienteNombre: "Benchmark Cliente",
			ClienteCedula: "1713175071",
			Productos: []models.ProductoInput{
				{
					Codigo:         "BENCH001",
					Descripcion:    "Benchmark Product",
					Cantidad:       1.0,
					PrecioUnitario: 100.0,
				},
			},
		},
		IncludeXML: false,
	}

	for i := 0; i < b.N; i++ {
		var body bytes.Buffer
		json.NewEncoder(&body).Encode(request)

		req := httptest.NewRequest(http.MethodPost, "/api/facturas", &body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleCreateFactura(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	}
}