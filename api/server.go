// Package api contiene la implementación del servidor HTTP REST
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	s.router.HandleFunc("/api/clientes/list", s.ListarClientesDB)
	s.router.HandleFunc("/api/clientes/", s.handleClienteDB)
	s.router.HandleFunc("/api/sri/estado", s.ConsultarEstadoSRI)
	s.router.HandleFunc("/api/sri/status", s.EstadoGeneralSRI)
	s.router.HandleFunc("/api/auditoria", s.ObtenerAuditoriaDB)
	s.router.HandleFunc("/api/respaldos", s.CrearRespaldoDB)
	s.router.HandleFunc("/api/respaldos/listar", s.ListarRespaldosDB)
	
	// Servir archivos estáticos del frontend (Astro build)
	s.setupStaticFiles()
	
	// Documentación básica - Solo para desarrollo cuando no hay frontend build
	s.router.HandleFunc("/api", s.handleRoot)
}

// handleFacturaDB maneja rutas dinámicas de facturas en DB
func (s *Server) handleFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/pdf") {
		s.GenerarPDFFacturaDB(w, r)
	} else if r.Method == http.MethodGet {
		s.ObtenerFacturaDB(w, r)
	} else if r.Method == http.MethodPut && strings.HasSuffix(r.URL.Path, "/estado") {
		s.ActualizarEstadoFacturaDB(w, r)
	} else if r.Method == http.MethodPut {
		s.ActualizarFacturaDB(w, r)
	} else if r.Method == http.MethodDelete {
		s.EliminarFacturaDB(w, r)
	} else {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Método no permitido")
	}
}

// handleClienteDB maneja rutas dinámicas de clientes en DB
func (s *Server) handleClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.ObtenerClienteDB(w, r)
	} else if r.Method == http.MethodPut {
		s.ActualizarClienteDB(w, r)
	} else if r.Method == http.MethodDelete {
		s.EliminarClienteDB(w, r)
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
	log.Printf("🌐 Frontend: http://localhost:%s/ (requiere build)", s.port)
	log.Printf("📊 API docs: http://localhost:%s/api", s.port)
	log.Printf("🛠️  Desarrollo frontend: cd web && pnpm dev (puerto 4321)")
	
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

// setupStaticFiles configura el servidor para servir archivos estáticos del frontend
func (s *Server) setupStaticFiles() {
	// Directorio donde Astro genera los archivos build
	staticDir := "./web/dist"
	
	// Verificar si el directorio existe
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("⚠️  Frontend no encontrado en %s. Ejecuta 'cd web && pnpm build' para generar archivos estáticos", staticDir)
		return
	}
	
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir(staticDir))
	
	// Manejar rutas del frontend (SPA routing)
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Si es una request de API, no procesar aquí
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		
		// Construir ruta del archivo
		filePath := filepath.Join(staticDir, r.URL.Path)
		
		// Si el archivo existe, servirlo directamente
		if _, err := os.Stat(filePath); err == nil {
			fs.ServeHTTP(w, r)
			return
		}
		
		// Si no existe, servir index.html para SPA routing
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
		} else {
			// Fallback: mostrar mensaje de desarrollo
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Facturación SRI - Desarrollo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; }
        .api-list { background: #ecf0f1; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .endpoint { margin: 5px 0; font-family: monospace; }
        .frontend-info { background: #e8f5e8; padding: 15px; border-radius: 5px; border-left: 4px solid #27ae60; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="header">🧾 Sistema de Facturación Electrónica SRI</h1>
        <h2>🚀 Servidor API Activo</h2>
        
        <div class="frontend-info">
            <h3>📱 Frontend en Desarrollo</h3>
            <p>Para ver la aplicación web completa:</p>
            <ol>
                <li>Abre otra terminal</li>
                <li>Ejecuta: <code>cd web && pnpm dev</code></li>
                <li>Visita: <a href="http://localhost:4321">http://localhost:4321</a></li>
            </ol>
        </div>
        
        <div class="api-list">
            <h3>🔗 Endpoints API Disponibles</h3>
            <div class="endpoint">GET /health - Health check</div>
            <div class="endpoint">GET /api/facturas/db/list - Listar facturas</div>
            <div class="endpoint">POST /api/facturas/db - Crear factura</div>
            <div class="endpoint">GET /api/estadisticas - Estadísticas</div>
            <div class="endpoint">POST /api/clientes - Guardar cliente</div>
            <div class="endpoint">GET /api/clientes/buscar - Buscar cliente</div>
            <div class="endpoint">GET /api/sri/estado - Estado SRI</div>
        </div>
        
        <p><strong>Modo:</strong> Desarrollo | <strong>Puerto:</strong> ` + s.port + `</p>
    </div>
</body>
</html>
			`))
		}
	})
	
	log.Printf("📁 Archivos estáticos configurados desde: %s", staticDir)
}