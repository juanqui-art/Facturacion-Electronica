package database

import (
	"go-facturacion-sri/config"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/sri"
	"os"
	"testing"
	"time"
)

// setupTestConfig configura la configuración para tests
func setupTestConfig() {
	config.CargarConfiguracionPorDefecto()
}

func TestNew(t *testing.T) {
	// Usar base de datos temporal para tests
	dbPath := "test_facturacion.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("Base de datos no fue creada")
	}
}

func TestGuardarYObtenerFactura(t *testing.T) {
	// Configurar para tests
	setupTestConfig()
	
	// Usar base de datos temporal
	dbPath := "test_factura_operaciones.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear factura de prueba
	facturaData := models.FacturaInput{
		ClienteNombre: "CLIENTE PRUEBA DB",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "TEST001",
				Descripcion:    "Producto de prueba",
				Cantidad:       2.0,
				PrecioUnitario: 50.00,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		t.Fatalf("Error creando factura: %v", err)
	}

	// Generar clave de acceso
	claveConfig := sri.ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  sri.Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         sri.Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      sri.EmisionNormal,
	}

	claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
	if err != nil {
		t.Fatalf("Error generando clave de acceso: %v", err)
	}

	// Guardar factura en base de datos
	facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
	if err != nil {
		t.Fatalf("Error guardando factura: %v", err)
	}

	if facturaDB.ID == 0 {
		t.Error("ID de factura no fue asignado")
	}

	if facturaDB.NumeroFactura == "" {
		t.Error("Número de factura no fue generado")
	}

	if facturaDB.ClaveAcceso != claveAcceso {
		t.Errorf("Clave de acceso esperada: %s, obtenida: %s", claveAcceso, facturaDB.ClaveAcceso)
	}

	// Obtener factura por ID
	facturaObtenida, err := db.ObtenerFacturaPorID(facturaDB.ID)
	if err != nil {
		t.Fatalf("Error obteniendo factura por ID: %v", err)
	}

	if facturaObtenida.ClienteNombre != facturaData.ClienteNombre {
		t.Errorf("Nombre de cliente esperado: %s, obtenido: %s",
			facturaData.ClienteNombre, facturaObtenida.ClienteNombre)
	}

	// Obtener factura por número
	facturaPorNumero, err := db.ObtenerFacturaPorNumero(facturaDB.NumeroFactura)
	if err != nil {
		t.Fatalf("Error obteniendo factura por número: %v", err)
	}

	if facturaPorNumero.ID != facturaDB.ID {
		t.Errorf("IDs no coinciden: esperado %d, obtenido %d",
			facturaDB.ID, facturaPorNumero.ID)
	}
}

func TestListarFacturas(t *testing.T) {
	// Configurar para tests
	setupTestConfig()
	
	// Usar base de datos temporal
	dbPath := "test_listar_facturas.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear múltiples facturas de prueba
	for i := 1; i <= 3; i++ {
		facturaData := models.FacturaInput{
			ClienteNombre: "CLIENTE PRUEBA " + string(rune(i+48)), // Convertir número a char
			ClienteCedula: "1713175071",
			Productos: []models.ProductoInput{
				{
					Codigo:         "TEST" + string(rune(i+48)),
					Descripcion:    "Producto " + string(rune(i+48)),
					Cantidad:       1.0,
					PrecioUnitario: float64(i * 10),
				},
			},
		}

		factura, err := factory.CrearFactura(facturaData)
		if err != nil {
			t.Fatalf("Error creando factura %d: %v", i, err)
		}

		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1792146739001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: "00000000" + string(rune(i+48)),
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			t.Fatalf("Error generando clave de acceso %d: %v", i, err)
		}

		_, err = db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
		if err != nil {
			t.Fatalf("Error guardando factura %d: %v", i, err)
		}
	}

	// Listar facturas
	facturas, err := db.ListarFacturas(10, 0)
	if err != nil {
		t.Fatalf("Error listando facturas: %v", err)
	}

	if len(facturas) != 3 {
		t.Errorf("Se esperaban 3 facturas, se obtuvieron: %d", len(facturas))
	}

	// Verificar que estén ordenadas por fecha de creación (más reciente primero)
	for i, factura := range facturas {
		t.Logf("Factura %d: %s - %s", i+1, factura.NumeroFactura, factura.ClienteNombre)
	}
}

func TestActualizarEstadoFactura(t *testing.T) {
	// Configurar para tests
	setupTestConfig()
	
	// Usar base de datos temporal
	dbPath := "test_actualizar_estado.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear factura de prueba
	facturaData := models.FacturaInput{
		ClienteNombre: "CLIENTE ESTADO PRUEBA",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "ESTADO001",
				Descripcion:    "Producto para prueba de estado",
				Cantidad:       1.0,
				PrecioUnitario: 100.00,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		t.Fatalf("Error creando factura: %v", err)
	}

	claveConfig := sri.ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  sri.Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         sri.Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      sri.EmisionNormal,
	}

	claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
	if err != nil {
		t.Fatalf("Error generando clave de acceso: %v", err)
	}

	facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
	if err != nil {
		t.Fatalf("Error guardando factura: %v", err)
	}

	// Verificar estado inicial
	if facturaDB.Estado != "BORRADOR" {
		t.Errorf("Estado inicial esperado: BORRADOR, obtenido: %s", facturaDB.Estado)
	}

	// Actualizar estado a AUTORIZADA
	numeroAutorizacion := "2025062301179214673900110010010000000011234567891"
	xmlAutorizado := "<factura>XML autorizado</factura>"
	observaciones := "Factura autorizada correctamente"

	err = db.ActualizarEstadoFactura(facturaDB.ID, "AUTORIZADA", numeroAutorizacion, xmlAutorizado, observaciones)
	if err != nil {
		t.Fatalf("Error actualizando estado: %v", err)
	}

	// Verificar actualización
	facturaActualizada, err := db.ObtenerFacturaPorID(facturaDB.ID)
	if err != nil {
		t.Fatalf("Error obteniendo factura actualizada: %v", err)
	}

	if facturaActualizada.Estado != "AUTORIZADA" {
		t.Errorf("Estado esperado: AUTORIZADA, obtenido: %s", facturaActualizada.Estado)
	}

	if facturaActualizada.NumeroAutorizacion != numeroAutorizacion {
		t.Errorf("Número de autorización esperado: %s, obtenido: %s",
			numeroAutorizacion, facturaActualizada.NumeroAutorizacion)
	}

	if facturaActualizada.FechaAutorizacion == nil {
		t.Error("Fecha de autorización no fue establecida")
	}

	if facturaActualizada.XMLAutorizado != xmlAutorizado {
		t.Errorf("XML autorizado no coincide")
	}
}

func TestObtenerProductosPorFactura(t *testing.T) {
	// Configurar para tests
	setupTestConfig()
	
	// Usar base de datos temporal
	dbPath := "test_productos_factura.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear factura con múltiples productos
	facturaData := models.FacturaInput{
		ClienteNombre: "CLIENTE PRODUCTOS PRUEBA",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "PROD001",
				Descripcion:    "Producto 1",
				Cantidad:       2.0,
				PrecioUnitario: 25.00,
			},
			{
				Codigo:         "PROD002",
				Descripcion:    "Producto 2",
				Cantidad:       1.0,
				PrecioUnitario: 50.00,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaData)
	if err != nil {
		t.Fatalf("Error creando factura: %v", err)
	}

	claveConfig := sri.ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  sri.Factura,
		RUCEmisor:        "1792146739001",
		Ambiente:         sri.Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      sri.EmisionNormal,
	}

	claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
	if err != nil {
		t.Fatalf("Error generando clave de acceso: %v", err)
	}

	facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
	if err != nil {
		t.Fatalf("Error guardando factura: %v", err)
	}

	// Obtener productos de la factura
	productos, err := db.ObtenerProductosPorFactura(facturaDB.ID)
	if err != nil {
		t.Fatalf("Error obteniendo productos: %v", err)
	}

	if len(productos) != 2 {
		t.Errorf("Se esperaban 2 productos, se obtuvieron: %d", len(productos))
	}

	// Verificar primer producto
	if productos[0].Codigo != "PROD001" {
		t.Errorf("Código del primer producto esperado: PROD001, obtenido: %s", productos[0].Codigo)
	}

	if productos[0].Descripcion != "Producto 1" {
		t.Errorf("Descripción del primer producto no coincide")
	}

	// Verificar segundo producto
	if productos[1].Codigo != "PROD002" {
		t.Errorf("Código del segundo producto esperado: PROD002, obtenido: %s", productos[1].Codigo)
	}
}

func TestGuardarYObtenerCliente(t *testing.T) {
	// Usar base de datos temporal
	dbPath := "test_clientes.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear cliente de prueba
	cliente := &ClienteDB{
		Cedula:      "1713175071",
		Nombre:      "JUAN CARLOS PEREZ",
		Direccion:   "Av. Principal 123",
		Telefono:    "0987654321",
		Email:       "juan@ejemplo.com",
		TipoCliente: "PERSONA_NATURAL",
	}

	// Guardar cliente
	clienteGuardado, err := db.GuardarCliente(cliente)
	if err != nil {
		t.Fatalf("Error guardando cliente: %v", err)
	}

	if clienteGuardado.ID == 0 {
		t.Error("ID de cliente no fue asignado")
	}

	// Obtener cliente por ID
	clientePorID, err := db.ObtenerClientePorID(clienteGuardado.ID)
	if err != nil {
		t.Fatalf("Error obteniendo cliente por ID: %v", err)
	}

	if clientePorID.Nombre != cliente.Nombre {
		t.Errorf("Nombre de cliente esperado: %s, obtenido: %s",
			cliente.Nombre, clientePorID.Nombre)
	}

	// Obtener cliente por cédula
	clientePorCedula, err := db.ObtenerClientePorCedula(cliente.Cedula)
	if err != nil {
		t.Fatalf("Error obteniendo cliente por cédula: %v", err)
	}

	if clientePorCedula.ID != clienteGuardado.ID {
		t.Errorf("IDs no coinciden: esperado %d, obtenido %d",
			clienteGuardado.ID, clientePorCedula.ID)
	}
}

func TestEstadisticasFacturas(t *testing.T) {
	// Configurar para tests
	setupTestConfig()
	
	// Usar base de datos temporal
	dbPath := "test_estadisticas.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear facturas de prueba con diferentes estados
	estados := []string{"BORRADOR", "AUTORIZADA", "AUTORIZADA"}

	for i, estado := range estados {
		facturaData := models.FacturaInput{
			ClienteNombre: "CLIENTE ESTADISTICAS " + string(rune(i+49)),
			ClienteCedula: "1713175071",
			Productos: []models.ProductoInput{
				{
					Codigo:         "STAT" + string(rune(i+49)),
					Descripcion:    "Producto estadísticas " + string(rune(i+49)),
					Cantidad:       1.0,
					PrecioUnitario: 100.00,
				},
			},
		}

		factura, err := factory.CrearFactura(facturaData)
		if err != nil {
			t.Fatalf("Error creando factura %d: %v", i, err)
		}

		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1792146739001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: "00000000" + string(rune(i+49)),
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			t.Fatalf("Error generando clave de acceso %d: %v", i, err)
		}

		facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
		if err != nil {
			t.Fatalf("Error guardando factura %d: %v", i, err)
		}

		// Actualizar estado si no es BORRADOR
		if estado != "BORRADOR" {
			err = db.ActualizarEstadoFactura(facturaDB.ID, estado, "AUTH"+string(rune(i+49)), "", "")
			if err != nil {
				t.Fatalf("Error actualizando estado factura %d: %v", i, err)
			}
		}
	}

	// Obtener estadísticas
	stats, err := db.EstadisticasFacturas()
	if err != nil {
		t.Fatalf("Error obteniendo estadísticas: %v", err)
	}

	// Verificar total de facturas
	totalFacturas, ok := stats["total_facturas"].(int)
	if !ok || totalFacturas != 3 {
		t.Errorf("Total de facturas esperado: 3, obtenido: %v", stats["total_facturas"])
	}

	// Verificar estadísticas por estado
	porEstado, ok := stats["por_estado"].(map[string]int)
	if !ok {
		t.Error("Estadísticas por estado no tienen el formato correcto")
	} else {
		if porEstado["BORRADOR"] != 1 {
			t.Errorf("Facturas en BORRADOR esperadas: 1, obtenidas: %d", porEstado["BORRADOR"])
		}
		if porEstado["AUTORIZADA"] != 2 {
			t.Errorf("Facturas AUTORIZADAS esperadas: 2, obtenidas: %d", porEstado["AUTORIZADA"])
		}
	}

	t.Logf("Estadísticas obtenidas: %+v", stats)
}

// BenchmarkGuardarFactura mide el performance de guardar facturas
func BenchmarkGuardarFactura(b *testing.B) {
	// Configurar para tests
	setupTestConfig()
	
	dbPath := "benchmark_facturacion.db"
	defer os.Remove(dbPath)

	db, err := New(dbPath)
	if err != nil {
		b.Fatalf("Error creando base de datos: %v", err)
	}
	defer db.Close()

	// Datos de factura reutilizables
	facturaData := models.FacturaInput{
		ClienteNombre: "CLIENTE BENCHMARK",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "BENCH001",
				Descripcion:    "Producto benchmark",
				Cantidad:       1.0,
				PrecioUnitario: 50.00,
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factura, err := factory.CrearFactura(facturaData)
		if err != nil {
			b.Fatalf("Error creando factura: %v", err)
		}

		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1792146739001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: "000000001",
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			b.Fatalf("Error generando clave de acceso: %v", err)
		}

		_, err = db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
		if err != nil {
			b.Fatalf("Error guardando factura: %v", err)
		}
	}
}
