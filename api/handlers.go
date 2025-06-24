package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
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

// FacturaStorageInterface define el contrato para almacenamiento de facturas
type FacturaStorageInterface interface {
	Store(id string, factura FacturaResponse)
	Get(id string) (FacturaResponse, bool)
	GetAll() []FacturaResponse
	GetNextID() int
	Count() int
}

// FacturaStorage representa un almacenamiento thread-safe de facturas
type FacturaStorage struct {
	mu       sync.RWMutex
	facturas map[string]FacturaResponse
	nextID   int
}

// NewFacturaStorage crea una nueva instancia de almacenamiento thread-safe
func NewFacturaStorage() *FacturaStorage {
	return &FacturaStorage{
		facturas: make(map[string]FacturaResponse),
		nextID:   1,
	}
}

// Store almacena una factura de forma thread-safe
func (fs *FacturaStorage) Store(id string, factura FacturaResponse) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.facturas[id] = factura
}

// Get obtiene una factura de forma thread-safe
func (fs *FacturaStorage) Get(id string) (FacturaResponse, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	factura, exists := fs.facturas[id]
	return factura, exists
}

// GetAll obtiene todas las facturas de forma thread-safe
func (fs *FacturaStorage) GetAll() []FacturaResponse {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	facturas := make([]FacturaResponse, 0, len(fs.facturas))
	for _, factura := range fs.facturas {
		// No incluir XML en la lista para reducir payload
		factura.XML = ""
		facturas = append(facturas, factura)
	}
	return facturas
}

// GetNextID obtiene el siguiente ID de forma thread-safe
func (fs *FacturaStorage) GetNextID() int {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	id := fs.nextID
	fs.nextID++
	return id
}

// Count retorna el número total de facturas
func (fs *FacturaStorage) Count() int {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return len(fs.facturas)
}

// Almacenamiento global thread-safe
var storage FacturaStorageInterface = NewFacturaStorage()

// SetStorage permite cambiar la implementación de storage (útil para testing)
func SetStorage(s FacturaStorageInterface) {
	storage = s
}

// LoggingFacturaStorage es una implementación que añade logging al storage
type LoggingFacturaStorage struct {
	underlying FacturaStorageInterface
}

// NewLoggingFacturaStorage crea un storage con logging
func NewLoggingFacturaStorage(underlying FacturaStorageInterface) *LoggingFacturaStorage {
	return &LoggingFacturaStorage{
		underlying: underlying,
	}
}

// Store implementa FacturaStorageInterface con logging
func (lfs *LoggingFacturaStorage) Store(id string, factura FacturaResponse) {
	fmt.Printf("🗄️  Almacenando factura: %s\n", id)
	lfs.underlying.Store(id, factura)
}

// Get implementa FacturaStorageInterface con logging
func (lfs *LoggingFacturaStorage) Get(id string) (FacturaResponse, bool) {
	fmt.Printf("🔍 Buscando factura: %s\n", id)
	factura, exists := lfs.underlying.Get(id)
	if !exists {
		fmt.Printf("❌ Factura no encontrada: %s\n", id)
	}
	return factura, exists
}

// GetAll implementa FacturaStorageInterface con logging
func (lfs *LoggingFacturaStorage) GetAll() []FacturaResponse {
	fmt.Printf("📋 Listando todas las facturas\n")
	return lfs.underlying.GetAll()
}

// GetNextID implementa FacturaStorageInterface con logging
func (lfs *LoggingFacturaStorage) GetNextID() int {
	id := lfs.underlying.GetNextID()
	fmt.Printf("🔢 Generando nuevo ID: %d\n", id)
	return id
}

// Count implementa FacturaStorageInterface con logging
func (lfs *LoggingFacturaStorage) Count() int {
	count := lfs.underlying.Count()
	fmt.Printf("📊 Total de facturas: %d\n", count)
	return count
}

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
	facturas := storage.GetAll()
	total := storage.Count()

	response := map[string]interface{}{
		"facturas": facturas,
		"total":    total,
		"message":  fmt.Sprintf("Se encontraron %d facturas", total),
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
	nextID := storage.GetNextID()
	id := fmt.Sprintf("FAC-%06d", nextID)

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

	// Guardar en storage thread-safe
	storage.Store(id, response)

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
	factura, exists := storage.Get(path)
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