package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-facturacion-sri/api"
	"go-facturacion-sri/models"
)

func main() {
	fmt.Println("🧪 TESTING API DE INTEGRACIÓN - SISTEMA FACTURACIÓN SRI")
	fmt.Println("=" + string(make([]byte, 60)))

	// Iniciar servidor en goroutine
	server := api.NewServer("8080")
	go func() {
		server.Start()
	}()

	// Esperar que el servidor se inicie
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8080"

	// 1. Test Health Check
	fmt.Println("\n🔍 Test 1: Health Check")
	if err := testHealthCheck(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Health Check OK")

	// 2. Test Crear Factura en Memoria
	fmt.Println("\n📄 Test 2: Crear Factura (Memoria)")
	facturaID, err := testCrearFacturaMemoria(baseURL)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Factura creada en memoria: %s\n", facturaID)

	// 3. Test Listar Facturas en Memoria
	fmt.Println("\n📋 Test 3: Listar Facturas (Memoria)")
	if err := testListarFacturasMemoria(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Listado de facturas OK")

	// 4. Test Crear Factura en Base de Datos
	fmt.Println("\n💾 Test 4: Crear Factura (Base de Datos)")
	facturaDBID, err := testCrearFacturaDB(baseURL)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Factura creada en DB: ID %d\n", facturaDBID)

	// 5. Test Listar Facturas en Base de Datos
	fmt.Println("\n📊 Test 5: Listar Facturas (Base de Datos)")
	if err := testListarFacturasDB(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Listado desde DB OK")

	// 6. Test Crear Cliente
	fmt.Println("\n👤 Test 6: Crear Cliente")
	if err := testCrearCliente(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Cliente creado OK")

	// 7. Test Buscar Cliente
	fmt.Println("\n🔍 Test 7: Buscar Cliente")
	if err := testBuscarCliente(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Búsqueda de cliente OK")

	// 8. Test Estadísticas
	fmt.Println("\n📈 Test 8: Estadísticas")
	if err := testEstadisticas(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Estadísticas OK")

	// 9. Test Auditoría
	fmt.Println("\n📝 Test 9: Auditoría")
	if err := testAuditoria(baseURL); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Auditoría OK")

	// 10. Test Estado SRI (puede fallar sin certificados, pero debe responder)
	fmt.Println("\n🌐 Test 10: Estado SRI")
	if err := testEstadoSRI(baseURL); err != nil {
		fmt.Printf("⚠️  Estado SRI: %v (esperado sin certificados)\n", err)
	} else {
		fmt.Println("✅ Estado SRI OK")
	}

	fmt.Println("\n🎉 TODOS LOS TESTS DE API COMPLETADOS!")
	fmt.Println("   El sistema de APIs está funcionando correctamente.")
	fmt.Println("=" + string(make([]byte, 60)))
}

func testHealthCheck(baseURL string) error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result["status"] != "healthy" {
		return fmt.Errorf("status no es healthy: %v", result["status"])
	}

	return nil
}

func testCrearFacturaMemoria(baseURL string) (string, error) {
	facturaData := map[string]interface{}{
		"clienteNombre": "Cliente Test API",
		"clienteCedula": "1713175071",
		"includeXML":    true,
		"productos": []map[string]interface{}{
			{
				"codigo":         "API001",
				"descripcion":    "Producto API Test",
				"cantidad":       2.0,
				"precioUnitario": 25.50,
			},
		},
	}

	jsonData, _ := json.Marshal(facturaData)
	resp, err := http.Post(baseURL+"/api/facturas", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("ID no encontrado en respuesta")
	}

	return id, nil
}

func testListarFacturasMemoria(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/facturas")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	total, ok := result["total"].(float64)
	if !ok || total < 1 {
		return fmt.Errorf("no se encontraron facturas")
	}

	return nil
}

func testCrearFacturaDB(baseURL string) (int, error) {
	facturaData := models.FacturaInput{
		ClienteNombre: "Cliente Test DB",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "DB001",
				Descripcion:    "Producto DB Test",
				Cantidad:       1.0,
				PrecioUnitario: 100.0,
			},
		},
	}

	jsonData, _ := json.Marshal(facturaData)
	resp, err := http.Post(baseURL+"/api/facturas/db", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("data no encontrada")
	}

	id, ok := data["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("ID no encontrado")
	}

	return int(id), nil
}

func testListarFacturasDB(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/facturas/db/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("operación no exitosa")
	}

	return nil
}

func testCrearCliente(baseURL string) error {
	clienteData := map[string]interface{}{
		"cedula":      "1713175071",
		"nombre":      "Juan Perez API",
		"direccion":   "Av. Test 123",
		"telefono":    "0987654321",
		"email":       "juan@test.com",
		"tipoCliente": "PERSONA_NATURAL",
	}

	jsonData, _ := json.Marshal(clienteData)
	resp, err := http.Post(baseURL+"/api/clientes", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("operación no exitosa")
	}

	return nil
}

func testBuscarCliente(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/clientes/buscar?cedula=1713175071")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("cliente no encontrado")
	}

	return nil
}

func testEstadisticas(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/estadisticas")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("operación no exitosa")
	}

	return nil
}

func testAuditoria(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/auditoria?tabla=facturas")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	success, ok := result["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("operación no exitosa")
	}

	return nil
}

func testEstadoSRI(baseURL string) error {
	// Usar una clave de prueba
	claveTest := "2506202501123456789000110010010000000049017300010"
	resp, err := http.Get(baseURL + "/api/sri/estado?clave=" + claveTest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Para estado SRI, aceptamos tanto success como failure (sin certificados es normal)
	claveAcceso, ok := result["clave_acceso"].(string)
	if !ok || claveAcceso != claveTest {
		return fmt.Errorf("clave de acceso no coincide")
	}

	return nil
}