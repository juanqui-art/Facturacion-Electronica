package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
)

// FacturaResponse - Estructura para respuestas de factura
type FacturaResponse struct {
	ID          string                 `json:"id"`
	Factura     models.Factura         `json:"factura"`
	XML         string                 `json:"xml,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	Status      string                 `json:"status"`
}

// CreateFacturaRequest - Estructura para crear facturas via API
type CreateFacturaRequest struct {
	models.FacturaInput
	IncludeXML bool `json:"includeXML,omitempty"`
}

// HealthResponse - Respuesta del health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Service   string    `json:"service"`
}

// Almacenamiento temporal en memoria (en el futuro será base de datos)
var facturaStorage = make(map[string]FacturaResponse)
var nextID = 1

// handleHealth - Health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Solo método GET permitido")
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Service:   "SRI Facturación Electrónica API",
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// handleRoot - Documentación básica de la API
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Solo método GET permitido")
		return
	}

	docs := map[string]interface{}{
		"service": "SRI Facturación Electrónica API",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"GET /health": "Health check del servicio",
			"GET /api/facturas": "Listar todas las facturas",
			"POST /api/facturas": "Crear nueva factura",
			"GET /api/facturas/{id}": "Obtener factura por ID",
		},
		"example_request": map[string]interface{}{
			"url": "/api/facturas",
			"method": "POST",
			"body": CreateFacturaRequest{
				FacturaInput: models.FacturaInput{
					ClienteNombre: "Juan Perez",
					ClienteCedula: "1713175071",
					Productos: []models.ProductoInput{
						{
							Codigo:         "PROD001",
							Descripcion:    "Producto de ejemplo",
							Cantidad:       1.0,
							PrecioUnitario: 100.0,
						},
					},
				},
				IncludeXML: true,
			},
		},
	}

	writeJSONResponse(w, http.StatusOK, docs)
}

// handleFacturas - Maneja /api/facturas (GET para listar, POST para crear)
func (s *Server) handleFacturas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListFacturas(w, r)
	case http.MethodPost:
		s.handleCreateFactura(w, r)
	default:
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Métodos permitidos: GET, POST")
	}
}

// handleListFacturas - Lista todas las facturas
func (s *Server) handleListFacturas(w http.ResponseWriter, r *http.Request) {
	facturas := make([]FacturaResponse, 0, len(facturaStorage))
	
	for _, factura := range facturaStorage {
		// No incluir XML en la lista para reducir payload
		factura.XML = ""
		facturas = append(facturas, factura)
	}

	response := map[string]interface{}{
		"facturas": facturas,
		"total":    len(facturas),
		"message":  fmt.Sprintf("Se encontraron %d facturas", len(facturas)),
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// handleCreateFactura - Crea una nueva factura
func (s *Server) handleCreateFactura(w http.ResponseWriter, r *http.Request) {
	var request CreateFacturaRequest

	// Decodificar JSON del request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	// Crear la factura usando nuestro factory
	factura, err := factory.CrearFactura(request.FacturaInput)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Error creando factura: "+err.Error())
		return
	}

	// Generar ID único
	id := fmt.Sprintf("FAC-%06d", nextID)
	nextID++

	// Crear respuesta
	response := FacturaResponse{
		ID:        id,
		Factura:   factura,
		CreatedAt: time.Now(),
		Status:    "created",
	}

	// Incluir XML si se solicita
	if request.IncludeXML {
		xmlData, err := factura.GenerarXML()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Error generando XML: "+err.Error())
			return
		}
		response.XML = string(xmlData)
	}

	// Guardar en storage temporal
	facturaStorage[id] = response

	writeJSONResponse(w, http.StatusCreated, response)
}

// handleFacturaByID - Maneja /api/facturas/{id}
func (s *Server) handleFacturaByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Solo método GET permitido")
		return
	}

	// Extraer ID de la URL
	path := strings.TrimPrefix(r.URL.Path, "/api/facturas/")
	if path == "" {
		writeErrorResponse(w, http.StatusBadRequest, "ID de factura requerido")
		return
	}

	// Buscar factura
	factura, exists := facturaStorage[path]
	if !exists {
		writeErrorResponse(w, http.StatusNotFound, "Factura no encontrada")
		return
	}

	// Incluir XML si se solicita via query parameter
	includeXML := r.URL.Query().Get("includeXML") == "true"
	if includeXML && factura.XML == "" {
		xmlData, err := factura.Factura.GenerarXML()
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Error generando XML: "+err.Error())
			return
		}
		factura.XML = string(xmlData)
	}

	writeJSONResponse(w, http.StatusOK, factura)
}