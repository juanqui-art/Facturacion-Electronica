package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-facturacion-sri/api"
	"go-facturacion-sri/config"
)

// setupAPITest - Configura el servidor para tests
func setupAPITest() *api.Server {
	config.CargarConfiguracionPorDefecto()
	return api.NewServer("8080")
}

// TestHealthEndpoint - Prueba el endpoint de health check
func TestHealthEndpoint(t *testing.T) {
	server := setupAPITest()
	
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := server.Router() // Necesitaremos exponer el router
	
	handler.ServeHTTP(rr, req)
	
	// Verificar status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler devolvió status incorrecto: obtuvo %v esperaba %v",
			status, http.StatusOK)
	}
	
	// Verificar que la respuesta sea JSON válido
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Respuesta no es JSON válido: %v", err)
	}
	
	// Verificar campos requeridos
	if response["status"] != "healthy" {
		t.Errorf("Status incorrecto: obtuvo %v esperaba %v", response["status"], "healthy")
	}
}

// TestCreateFacturaEndpoint - Prueba crear factura via API
func TestCreateFacturaEndpoint(t *testing.T) {
	server := setupAPITest()
	
	// Crear request body
	requestBody := map[string]interface{}{
		"clienteNombre": "Juan Perez Test",
		"clienteCedula": "1713175071",
		"productos": []map[string]interface{}{
			{
				"codigo":         "TEST001",
				"descripcion":    "Producto de prueba",
				"cantidad":       2.0,
				"precioUnitario": 50.0,
			},
		},
		"includeXML": true,
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	
	req, err := http.NewRequest("POST", "/api/facturas", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := server.Router()
	
	handler.ServeHTTP(rr, req)
	
	// Verificar status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler devolvió status incorrecto: obtuvo %v esperaba %v",
			status, http.StatusCreated)
	}
	
	// Verificar respuesta
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Respuesta no es JSON válido: %v", err)
	}
	
	// Verificar que se creó la factura
	if response["id"] == nil {
		t.Error("Respuesta no contiene ID de factura")
	}
	
	if response["status"] != "created" {
		t.Errorf("Status incorrecto: obtuvo %v esperaba %v", response["status"], "created")
	}
	
	// Verificar que incluye XML
	if response["xml"] == nil {
		t.Error("Respuesta no contiene XML solicitado")
	}
}

// TestCreateFacturaWithInvalidData - Prueba validación de datos inválidos
func TestCreateFacturaWithInvalidData(t *testing.T) {
	server := setupAPITest()
	
	// Request con cédula inválida
	requestBody := map[string]interface{}{
		"clienteNombre": "Juan Perez Test",
		"clienteCedula": "123", // Cédula inválida
		"productos": []map[string]interface{}{
			{
				"codigo":         "TEST001",
				"descripcion":    "Producto de prueba", 
				"cantidad":       2.0,
				"precioUnitario": 50.0,
			},
		},
	}
	
	jsonBody, _ := json.Marshal(requestBody)
	
	req, err := http.NewRequest("POST", "/api/facturas", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := server.Router()
	
	handler.ServeHTTP(rr, req)
	
	// Debe devolver Bad Request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler devolvió status incorrecto: obtuvo %v esperaba %v",
			status, http.StatusBadRequest)
	}
	
	// Verificar que contiene mensaje de error
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Respuesta no es JSON válido: %v", err)
	}
	
	if response["error"] != true {
		t.Error("Respuesta debería indicar error")
	}
}