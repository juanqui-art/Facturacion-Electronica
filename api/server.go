// Package api contiene la implementación del servidor HTTP REST
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server - Estructura principal del servidor HTTP
type Server struct {
	port   string
	router *http.ServeMux
}

// NewServer - Crea una nueva instancia del servidor
func NewServer(port string) *Server {
	server := &Server{
		port:   port,
		router: http.NewServeMux(),
	}
	
	// Configurar rutas
	server.setupRoutes()
	
	return server
}

// setupRoutes - Configura todas las rutas de la API
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.handleHealth)
	
	// API routes originales (en memoria)
	s.router.HandleFunc("/api/facturas", s.handleFacturas)
	s.router.HandleFunc("/api/facturas/", s.handleFacturaByID)
	
	// API routes con base de datos
	s.router.HandleFunc("/api/facturas/db", s.CrearFacturaConDB)
	s.router.HandleFunc("/api/facturas/db/list", s.ListarFacturasDB)
	s.router.HandleFunc("/api/facturas/db/", s.handleFacturaDB)
	s.router.HandleFunc("/api/estadisticas", s.EstadisticasDB)
	s.router.HandleFunc("/api/clientes", s.GuardarClienteDB)
	s.router.HandleFunc("/api/clientes/buscar", s.BuscarClienteDB)
	s.router.HandleFunc("/api/sri/estado", s.ConsultarEstadoSRI)
	s.router.HandleFunc("/api/auditoria", s.ObtenerAuditoriaDB)
	s.router.HandleFunc("/api/respaldos", s.CrearRespaldoDB)
	s.router.HandleFunc("/api/respaldos/listar", s.ListarRespaldosDB)
	
	// Documentación básica
	s.router.HandleFunc("/", s.handleRoot)
}

// handleFacturaDB maneja rutas dinámicas de facturas en DB
func (s *Server) handleFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.ObtenerFacturaDB(w, r)
	} else if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/estado") {
		s.ActualizarEstadoFacturaDB(w, r)
	} else {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// Start - Inicia el servidor HTTP
func (s *Server) Start() error {
	// Crear servidor HTTP con configuración personalizada
	httpServer := &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.middlewareChain(s.router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("🚀 Servidor iniciado en http://localhost:%s", s.port)
	log.Printf("📋 Health check: http://localhost:%s/health", s.port)
	log.Printf("📊 API docs: http://localhost:%s/", s.port)
	
	return httpServer.ListenAndServe()
}

// Router - Expone el router para testing
func (s *Server) Router() http.Handler {
	return s.middlewareChain(s.router)
}

// middlewareChain - Aplica middleware a todas las requests
func (s *Server) middlewareChain(next http.Handler) http.Handler {
	return s.corsMiddleware(s.loggingMiddleware(next))
}

// ResponseWriter helper functions
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	response := map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  status,
	}
	writeJSONResponse(w, status, response)
}