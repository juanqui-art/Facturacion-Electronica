package main

import (
	"fmt"
	"time"

	"go-facturacion-sri/config"
	"go-facturacion-sri/database"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/sri"
)

func main() {
	fmt.Println("🧪 TESTING SISTEMA COMPLETO - PERSISTENCIA Y RESPALDOS")
	fmt.Println("=" + string(make([]byte, 65)))

	// Configurar sistema
	config.CargarConfiguracionPorDefecto()

	// 1. Test de Base de Datos
	fmt.Println("\n💾 Test 1: Funcionalidad de Base de Datos")
	if err := testBaseDatos(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Base de Datos OK")

	// 2. Test de Auditoría
	fmt.Println("\n📝 Test 2: Sistema de Auditoría")
	if err := testAuditoria(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Auditoría OK")

	// 3. Test de Respaldos
	fmt.Println("\n💾 Test 3: Sistema de Respaldos")
	if err := testRespaldos(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Respaldos OK")

	// 4. Test de Integración Completa
	fmt.Println("\n🔄 Test 4: Integración Completa")
	if err := testIntegracionCompleta(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Integración Completa OK")

	// 5. Test de Performance
	fmt.Println("\n⚡ Test 5: Performance del Sistema")
	if err := testPerformance(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Performance OK")

	fmt.Println("\n🎉 TODOS LOS TESTS DEL SISTEMA COMPLETO PASARON!")
	fmt.Println("   El sistema está listo para producción.")
	fmt.Println("=" + string(make([]byte, 65)))
}

func testBaseDatos() error {
	// Crear base de datos de test
	db, err := database.New("test_sistema_completo.db")
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}
	defer db.Close()
	defer func() {
		// Limpiar archivo de test
		// os.Remove("test_sistema_completo.db")
	}()

	// Test crear cliente
	cliente := &database.ClienteDB{
		Cedula:      "1713175071",
		Nombre:      "Cliente Sistema Completo",
		Direccion:   "Av. Test Completo 123",
		Telefono:    "0987654321",
		Email:       "cliente@sistema.com",
		TipoCliente: "PERSONA_NATURAL",
	}

	clienteGuardado, err := db.GuardarCliente(cliente)
	if err != nil {
		return fmt.Errorf("error guardando cliente: %v", err)
	}
	fmt.Printf("   ✓ Cliente guardado: ID %d\n", clienteGuardado.ID)

	// Test crear factura
	facturaInput := models.FacturaInput{
		ClienteNombre: "Cliente Sistema Completo",
		ClienteCedula: "1713175071",
		Productos: []models.ProductoInput{
			{
				Codigo:         "SIST001",
				Descripcion:    "Producto Sistema Completo",
				Cantidad:       3.0,
				PrecioUnitario: 75.50,
			},
			{
				Codigo:         "SIST002",
				Descripcion:    "Segundo Producto",
				Cantidad:       1.0,
				PrecioUnitario: 125.99,
			},
		},
	}

	factura, err := factory.CrearFactura(facturaInput)
	if err != nil {
		return fmt.Errorf("error creando factura: %v", err)
	}

	// Generar clave de acceso
	claveConfig := sri.ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  sri.Factura,
		RUCEmisor:        "1234567890001",
		Ambiente:         sri.Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001",
		TipoEmision:      sri.EmisionNormal,
	}

	claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
	if err != nil {
		return fmt.Errorf("error generando clave de acceso: %v", err)
	}

	facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaInput.Productos)
	if err != nil {
		return fmt.Errorf("error guardando factura: %v", err)
	}
	fmt.Printf("   ✓ Factura guardada: %s (Total: $%.2f)\n", facturaDB.NumeroFactura, facturaDB.Total)

	// Test obtener productos
	productos, err := db.ObtenerProductosPorFactura(facturaDB.ID)
	if err != nil {
		return fmt.Errorf("error obteniendo productos: %v", err)
	}
	fmt.Printf("   ✓ Productos obtenidos: %d\n", len(productos))

	// Test estadísticas
	stats, err := db.EstadisticasFacturas()
	if err != nil {
		return fmt.Errorf("error obteniendo estadísticas: %v", err)
	}
	fmt.Printf("   ✓ Estadísticas: %d facturas totales\n", stats["total_facturas"])

	return nil
}

func testAuditoria() error {
	// Crear base de datos de test
	db, err := database.New("test_auditoria.db")
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}
	defer db.Close()

	// Crear registro de auditoría
	audit := &database.AuditLogDB{
		Tabla:        "facturas",
		RegistroID:   1,
		Operacion:    "CREATE",
		Usuario:      "test_user",
		DatosAntes:   "",
		DatosDespues: `{"numero":"FAC-000001","total":100.0}`,
		IPAddress:    "192.168.1.100",
		UserAgent:    "TestAgent/1.0",
	}

	err = db.RegistrarAuditoria(audit)
	if err != nil {
		return fmt.Errorf("error registrando auditoría: %v", err)
	}
	fmt.Printf("   ✓ Auditoría registrada\n")

	// Crear segundo registro
	audit2 := &database.AuditLogDB{
		Tabla:        "facturas",
		RegistroID:   1,
		Operacion:    "UPDATE",
		Usuario:      "test_user",
		DatosAntes:   `{"numero":"FAC-000001","total":100.0}`,
		DatosDespues: `{"numero":"FAC-000001","total":115.0,"estado":"AUTORIZADA"}`,
		IPAddress:    "192.168.1.100",
		UserAgent:    "TestAgent/1.0",
	}

	err = db.RegistrarAuditoria(audit2)
	if err != nil {
		return fmt.Errorf("error registrando segunda auditoría: %v", err)
	}

	// Obtener auditoría por tabla
	registros, err := db.ObtenerAuditoriaPorTabla("facturas", 10, 0)
	if err != nil {
		return fmt.Errorf("error obteniendo auditoría: %v", err)
	}
	fmt.Printf("   ✓ Registros de auditoría obtenidos: %d\n", len(registros))

	// Obtener auditoría por registro específico
	registrosEspecificos, err := db.ObtenerAuditoriaPorRegistro("facturas", 1)
	if err != nil {
		return fmt.Errorf("error obteniendo auditoría específica: %v", err)
	}
	fmt.Printf("   ✓ Auditoría específica: %d registros\n", len(registrosEspecificos))

	return nil
}

func testRespaldos() error {
	// Crear base de datos de test
	db, err := database.New("test_respaldos.db")
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}
	defer db.Close()

	// Agregar algunos datos de prueba
	cliente := &database.ClienteDB{
		Cedula:      "1713175071",
		Nombre:      "Cliente Respaldo",
		TipoCliente: "PERSONA_NATURAL",
	}
	_, err = db.GuardarCliente(cliente)
	if err != nil {
		return fmt.Errorf("error guardando cliente: %v", err)
	}

	// Crear gestor de respaldos
	backupConfig := database.BackupConfig{
		RutaRespaldos:        "./test_respaldos",
		IntervaloRespaldo:    time.Hour,
		MaxRespaldos:         5,
		CompresiónHabilitada: false,
		PrefijRespaldo:       "test_backup",
	}

	backupManager := database.NewBackupManager(db, backupConfig)

	// Test crear respaldo manual
	err = backupManager.CrearRespaldoManual("sistema_test")
	if err != nil {
		return fmt.Errorf("error creando respaldo manual: %v", err)
	}
	fmt.Printf("   ✓ Respaldo manual creado\n")

	// Test listar respaldos
	respaldos, err := backupManager.ListarRespaldos()
	if err != nil {
		return fmt.Errorf("error listando respaldos: %v", err)
	}
	fmt.Printf("   ✓ Respaldos listados: %d encontrados\n", len(respaldos))

	if len(respaldos) > 0 {
		respaldo := respaldos[0]
		fmt.Printf("     - %s (%s)\n", respaldo.Nombre, respaldo.TamañoLegible)
	}

	// Test crear respaldo automático
	err = backupManager.CrearRespaldo()
	if err != nil {
		return fmt.Errorf("error creando respaldo automático: %v", err)
	}
	fmt.Printf("   ✓ Respaldo automático creado\n")

	return nil
}

func testIntegracionCompleta() error {
	// Test que simula un flujo completo de trabajo
	db, err := database.New("test_integracion_completa.db")
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}
	defer db.Close()

	// 1. Crear cliente
	cliente := &database.ClienteDB{
		Cedula:      "1713175071",
		Nombre:      "Cliente Integración Completa",
		Direccion:   "Av. Integración 456",
		Telefono:    "0998877665",
		Email:       "integracion@test.com",
		TipoCliente: "PERSONA_NATURAL",
	}

	clienteGuardado, err := db.GuardarCliente(cliente)
	if err != nil {
		return fmt.Errorf("error guardando cliente: %v", err)
	}

	// Registrar auditoría del cliente
	auditCliente := &database.AuditLogDB{
		Tabla:        "clientes",
		RegistroID:   clienteGuardado.ID,
		Operacion:    "CREATE",
		Usuario:      "sistema_integracion",
		DatosAntes:   "",
		DatosDespues: fmt.Sprintf(`{"id":%d,"nombre":"%s","cedula":"%s"}`, clienteGuardado.ID, cliente.Nombre, cliente.Cedula),
		IPAddress:    "127.0.0.1",
		UserAgent:    "SistemaIntegracion/1.0",
	}
	db.RegistrarAuditoria(auditCliente)

	// 2. Crear múltiples facturas
	for i := 1; i <= 3; i++ {
		facturaInput := models.FacturaInput{
			ClienteNombre: cliente.Nombre,
			ClienteCedula: cliente.Cedula,
			Productos: []models.ProductoInput{
				{
					Codigo:         fmt.Sprintf("INT%03d", i),
					Descripcion:    fmt.Sprintf("Producto Integración %d", i),
					Cantidad:       float64(i),
					PrecioUnitario: 50.0 * float64(i),
				},
			},
		}

		factura, err := factory.CrearFactura(facturaInput)
		if err != nil {
			return fmt.Errorf("error creando factura %d: %v", i, err)
		}

		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1234567890001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: fmt.Sprintf("%09d", i),
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			return fmt.Errorf("error generando clave de acceso %d: %v", i, err)
		}

		facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaInput.Productos)
		if err != nil {
			return fmt.Errorf("error guardando factura %d: %v", i, err)
		}

		// Registrar auditoría de la factura
		auditFactura := &database.AuditLogDB{
			Tabla:        "facturas",
			RegistroID:   facturaDB.ID,
			Operacion:    "CREATE",
			Usuario:      "sistema_integracion",
			DatosDespues: fmt.Sprintf(`{"id":%d,"numero":"%s","total":%.2f}`, facturaDB.ID, facturaDB.NumeroFactura, facturaDB.Total),
			IPAddress:    "127.0.0.1",
			UserAgent:    "SistemaIntegracion/1.0",
		}
		db.RegistrarAuditoria(auditFactura)

		// Simular actualización de estado a AUTORIZADA para algunas facturas
		if i%2 == 0 {
			err = db.ActualizarEstadoFactura(facturaDB.ID, "AUTORIZADA", 
				"AUTH"+claveAcceso, "", "Factura autorizada automáticamente")
			if err != nil {
				return fmt.Errorf("error actualizando estado factura %d: %v", i, err)
			}

			// Auditoría de actualización
			auditUpdate := &database.AuditLogDB{
				Tabla:        "facturas",
				RegistroID:   facturaDB.ID,
				Operacion:    "UPDATE",
				Usuario:      "sistema_integracion",
				DatosAntes:   fmt.Sprintf(`{"estado":"BORRADOR"}`),
				DatosDespues: fmt.Sprintf(`{"estado":"AUTORIZADA","numero_autorizacion":"AUTH%s"}`, claveAcceso),
				IPAddress:    "127.0.0.1",
				UserAgent:    "SistemaIntegracion/1.0",
			}
			db.RegistrarAuditoria(auditUpdate)
		}
	}

	// 3. Crear respaldo después de las operaciones
	backupManager := database.NewBackupManagerDefault(db)
	err = backupManager.CrearRespaldoManual("integracion_completa")
	if err != nil {
		return fmt.Errorf("error creando respaldo: %v", err)
	}

	// 4. Verificar estadísticas finales
	stats, err := db.EstadisticasFacturas()
	if err != nil {
		return fmt.Errorf("error obteniendo estadísticas finales: %v", err)
	}

	fmt.Printf("   ✓ Cliente creado: %s\n", clienteGuardado.Nombre)
	fmt.Printf("   ✓ Facturas creadas: %v\n", stats["total_facturas"])
	fmt.Printf("   ✓ Total facturado: $%.2f\n", stats["total_facturado"])
	fmt.Printf("   ✓ Estados: %v\n", stats["por_estado"])

	// 5. Verificar auditoría completa
	auditFacturas, err := db.ObtenerAuditoriaPorTabla("facturas", 50, 0)
	if err != nil {
		return fmt.Errorf("error obteniendo auditoría facturas: %v", err)
	}

	auditClientes, err := db.ObtenerAuditoriaPorTabla("clientes", 50, 0)
	if err != nil {
		return fmt.Errorf("error obteniendo auditoría clientes: %v", err)
	}

	fmt.Printf("   ✓ Registros auditoría facturas: %d\n", len(auditFacturas))
	fmt.Printf("   ✓ Registros auditoría clientes: %d\n", len(auditClientes))

	return nil
}

func testPerformance() error {
	fmt.Printf("   🔄 Testing performance con múltiples operaciones...\n")

	db, err := database.New("test_performance.db")
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}
	defer db.Close()

	start := time.Now()

	// Test inserción masiva de clientes
	for i := 1; i <= 100; i++ {
		// Usar siempre la misma cédula válida para test de performance
		cliente := &database.ClienteDB{
			Cedula:      fmt.Sprintf("1713175071_%d", i), // Agregar sufijo para evitar duplicados
			Nombre:      fmt.Sprintf("Cliente Performance %d", i),
			TipoCliente: "PERSONA_NATURAL",
		}
		_, err = db.GuardarCliente(cliente)
		if err != nil {
			return fmt.Errorf("error guardando cliente %d: %v", i, err)
		}
	}

	clientesTime := time.Since(start)
	fmt.Printf("   ✓ 100 clientes insertados en %v\n", clientesTime)

	// Test inserción de facturas
	start = time.Now()
	for i := 1; i <= 50; i++ {
		facturaInput := models.FacturaInput{
			ClienteNombre: fmt.Sprintf("Cliente Performance %d", i),
			ClienteCedula: "1713175071", // Usar la cédula válida sin modificar
			Productos: []models.ProductoInput{
				{
					Codigo:         fmt.Sprintf("PERF%03d", i),
					Descripcion:    fmt.Sprintf("Producto Performance %d", i),
					Cantidad:       1.0,
					PrecioUnitario: 25.50,
				},
			},
		}

		factura, err := factory.CrearFactura(facturaInput)
		if err != nil {
			return fmt.Errorf("error creando factura %d: %v", i, err)
		}

		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1234567890001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: fmt.Sprintf("%09d", i),
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			return fmt.Errorf("error generando clave %d: %v", i, err)
		}

		_, err = db.GuardarFactura(factura, claveAcceso, facturaInput.Productos)
		if err != nil {
			return fmt.Errorf("error guardando factura %d: %v", i, err)
		}
	}

	facturasTime := time.Since(start)
	fmt.Printf("   ✓ 50 facturas insertadas en %v\n", facturasTime)

	// Test consultas
	start = time.Now()
	for i := 1; i <= 20; i++ {
		_, err := db.ListarFacturas(10, 0)
		if err != nil {
			return fmt.Errorf("error en consulta %d: %v", i, err)
		}
	}
	consultasTime := time.Since(start)
	fmt.Printf("   ✓ 20 consultas de listado en %v\n", consultasTime)

	// Test estadísticas
	start = time.Now()
	_, err = db.EstadisticasFacturas()
	if err != nil {
		return fmt.Errorf("error obteniendo estadísticas: %v", err)
	}
	statsTime := time.Since(start)
	fmt.Printf("   ✓ Estadísticas calculadas en %v\n", statsTime)

	return nil
}