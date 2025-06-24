package api

import (
	"log"
	"net/http"
	"time"
)

// loggingMiddleware - Middleware para logging de requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Crear wrapper para capturar status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Llamar al siguiente handler
		next.ServeHTTP(wrapped, r)
		
		// Log de la request
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// corsMiddleware - Middleware para CORS (Cross-Origin Resource Sharing)
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Manejar preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Llamar al siguiente handler
		next.ServeHTTP(w, r)
	})
}

// responseWriter - Wrapper para capturar el status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader - Sobrescribe WriteHeader para capturar status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}