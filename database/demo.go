// Demo del sistema de base de datos
package database

import (
	"fmt"
	"strings"
	"time"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/sri"
)

// DemoDatabase ejecuta una demostración completa del sistema de base de datos
func DemoDatabase() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🗄️  DEMO SISTEMA DE BASE DE DATOS")
	fmt.Println(strings.Repeat("=", 60))

	// Demo 1: Inicializar base de datos
	fmt.Println("\n1️⃣ INICIALIZACIÓN DE BASE DE DATOS")
	fmt.Println(strings.Repeat("-", 40))
	
	db, err := New("database/demo_facturacion.db")
	if err != nil {
		fmt.Printf("❌ Error inicializando base de datos: %v\n", err)
		return
	}
	defer db.Close()
	
	fmt.Println("✅ Base de datos inicializada correctamente")
	fmt.Println("📂 Archivo: database/demo_facturacion.db")
	fmt.Println("📋 Tablas creadas: facturas, productos, clientes, configuracion")

	// Demo 2: Crear y guardar facturas
	fmt.Println("\n2️⃣ CREACIÓN Y GUARDADO DE FACTURAS")
	fmt.Println(strings.Repeat("-", 40))
	
	facturas := []models.FacturaInput{
		{
			ClienteNombre: "EMPRESA TECNOLOGIA XYZ S.A.",
			ClienteCedula: "1713175071",
			Productos: []models.ProductoInput{
				{
					Codigo:         "LAPTOP001",
					Descripcion:    "Laptop Dell Inspiron 15 - 8GB RAM",
					Cantidad:       2.0,
					PrecioUnitario: 750.00,
				},
				{
					Codigo:         "MOUSE001",
					Descripcion:    "Mouse Inalámbrico Logitech MX Master",
					Cantidad:       2.0,
					PrecioUnitario: 45.00,
				},
			},
		},
		{
			ClienteNombre: "CONSULTORA BUSINESS SOLUTIONS",
			ClienteCedula: "1234567890",
			Productos: []models.ProductoInput{
				{
					Codigo:         "CONSUL001",
					Descripcion:    "Consultoría en Desarrollo de Software",
					Cantidad:       40.0,
					PrecioUnitario: 85.00,
				},
			},
		},
		{
			ClienteNombre: "DISTRIBUIDORA PRODUCTOS ECUADOR",
			ClienteCedula: "0987654321",
			Productos: []models.ProductoInput{
				{
					Codigo:         "PROD001",
					Descripcion:    "Camisetas Promocionales",
					Cantidad:       100.0,
					PrecioUnitario: 12.50,
				},
				{
					Codigo:         "PROD002",
					Descripcion:    "Gorras Bordadas",
					Cantidad:       50.0,
					PrecioUnitario: 8.75,
				},
			},
		},
	}

	var facturasGuardadas []*FacturaDB
	
	for i, facturaData := range facturas {
		fmt.Printf("\n📝 Procesando factura %d: %s\n", i+1, facturaData.ClienteNombre)
		
		// Crear factura
		factura, err := factory.CrearFactura(facturaData)
		if err != nil {
			fmt.Printf("❌ Error creando factura %d: %v\n", i+1, err)
			continue
		}

		// Generar clave de acceso
		claveConfig := sri.ClaveAccesoConfig{
			FechaEmision:     time.Now(),
			TipoComprobante:  sri.Factura,
			RUCEmisor:        "1792146739001",
			Ambiente:         sri.Pruebas,
			Serie:            "001001",
			NumeroSecuencial: fmt.Sprintf("%09d", i+1),
			TipoEmision:      sri.EmisionNormal,
		}

		claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
		if err != nil {
			fmt.Printf("❌ Error generando clave de acceso %d: %v\n", i+1, err)
			continue
		}

		// Guardar en base de datos
		facturaDB, err := db.GuardarFactura(factura, claveAcceso, facturaData.Productos)
		if err != nil {
			fmt.Printf("❌ Error guardando factura %d: %v\n", i+1, err)
			continue
		}

		facturasGuardadas = append(facturasGuardadas, facturaDB)
		
		fmt.Printf("✅ Factura guardada: %s\n", facturaDB.NumeroFactura)
		fmt.Printf("💰 Total: $%.2f\n", facturaDB.Total)
		fmt.Printf("🔑 Clave: %s\n", sri.FormatearClaveAcceso(facturaDB.ClaveAcceso))
	}

	// Demo 3: Listar facturas
	fmt.Println("\n3️⃣ LISTADO DE FACTURAS")
	fmt.Println(strings.Repeat("-", 40))
	
	todasFacturas, err := db.ListarFacturas(10, 0)
	if err != nil {
		fmt.Printf("❌ Error listando facturas: %v\n", err)
	} else {
		fmt.Printf("📊 Total de facturas en base de datos: %d\n", len(todasFacturas))
		
		for i, factura := range todasFacturas {
			fmt.Printf("\n%d. %s\n", i+1, factura.NumeroFactura)
			fmt.Printf("   👤 Cliente: %s\n", factura.ClienteNombre)
			fmt.Printf("   💰 Total: $%.2f\n", factura.Total)
			fmt.Printf("   📅 Fecha: %s\n", factura.FechaEmision.Format("02/01/2006"))
			fmt.Printf("   📊 Estado: %s\n", factura.Estado)
		}
	}

	// Demo 4: Actualizar estados de facturas (simular autorización SRI)
	fmt.Println("\n4️⃣ SIMULACIÓN DE AUTORIZACIÓN SRI")
	fmt.Println(strings.Repeat("-", 40))
	
	if len(facturasGuardadas) > 0 {
		// Autorizar las primeras dos facturas
		for i := 0; i < min(2, len(facturasGuardadas)); i++ {
			factura := facturasGuardadas[i]
			
			fmt.Printf("\n🔐 Autorizando factura: %s\n", factura.NumeroFactura)
			
			// Simular respuesta del SRI
			autorizacion := sri.SimularAutorizacionSRI(factura.ClaveAcceso, sri.Pruebas)
			xmlAutorizado := fmt.Sprintf("<facturaAutorizada>XML autorizado para %s</facturaAutorizada>", factura.NumeroFactura)
			
			err := db.ActualizarEstadoFactura(
				factura.ID,
				"AUTORIZADA",
				autorizacion.NumeroAutorizacion,
				xmlAutorizado,
				"Factura autorizada automáticamente por el SRI",
			)
			
			if err != nil {
				fmt.Printf("❌ Error actualizando estado: %v\n", err)
			} else {
				fmt.Printf("✅ Factura %s AUTORIZADA\n", factura.NumeroFactura)
				fmt.Printf("📝 N° Autorización: %s\n", autorizacion.NumeroAutorizacion)
			}
		}
	}

	// Demo 5: Consultar factura específica con productos
	fmt.Println("\n5️⃣ CONSULTA DETALLADA DE FACTURA")
	fmt.Println(strings.Repeat("-", 40))
	
	if len(facturasGuardadas) > 0 {
		facturaID := facturasGuardadas[0].ID
		
		facturaDetalle, err := db.ObtenerFacturaPorID(facturaID)
		if err != nil {
			fmt.Printf("❌ Error obteniendo factura: %v\n", err)
		} else {
			fmt.Printf("📄 Factura: %s\n", facturaDetalle.NumeroFactura)
			fmt.Printf("👤 Cliente: %s (%s)\n", facturaDetalle.ClienteNombre, facturaDetalle.ClienteCedula)
			fmt.Printf("💰 Subtotal: $%.2f\n", facturaDetalle.Subtotal)
			fmt.Printf("🧮 IVA: $%.2f\n", facturaDetalle.IVA)
			fmt.Printf("💵 Total: $%.2f\n", facturaDetalle.Total)
			fmt.Printf("📊 Estado: %s\n", facturaDetalle.Estado)
			
			if facturaDetalle.FechaAutorizacion != nil {
				fmt.Printf("🔐 Fecha Autorización: %s\n", facturaDetalle.FechaAutorizacion.Format("02/01/2006 15:04:05"))
			}
			
			// Obtener productos
			productos, err := db.ObtenerProductosPorFactura(facturaID)
			if err != nil {
				fmt.Printf("❌ Error obteniendo productos: %v\n", err)
			} else {
				fmt.Printf("\n📦 Productos (%d items):\n", len(productos))
				for j, producto := range productos {
					fmt.Printf("   %d. %s - %s\n", j+1, producto.Codigo, producto.Descripcion)
					fmt.Printf("      Cantidad: %.2f | Precio: $%.2f | Total: $%.2f\n",
						producto.Cantidad, producto.PrecioUnitario, producto.PrecioTotal)
				}
			}
		}
	}

	// Demo 6: Gestión de clientes
	fmt.Println("\n6️⃣ GESTIÓN DE CLIENTES")
	fmt.Println(strings.Repeat("-", 40))
	
	clientesDemo := []*ClienteDB{
		{
			Cedula:      "1713175071",
			Nombre:      "EMPRESA TECNOLOGIA XYZ S.A.",
			Direccion:   "Av. Amazonas N24-03 y Colón",
			Telefono:    "02-2234567",
			Email:       "facturacion@tecnologiaxyz.com",
			TipoCliente: "EMPRESA",
		},
		{
			Cedula:      "1234567890",
			Nombre:      "JUAN CARLOS PEREZ MENDOZA",
			Direccion:   "Calle García Moreno 456",
			Telefono:    "0987654321",
			Email:       "juan.perez@email.com",
			TipoCliente: "PERSONA_NATURAL",
		},
	}
	
	for i, clienteData := range clientesDemo {
		fmt.Printf("\n👤 Guardando cliente %d: %s\n", i+1, clienteData.Nombre)
		
		cliente, err := db.GuardarCliente(clienteData)
		if err != nil {
			fmt.Printf("❌ Error guardando cliente: %v\n", err)
		} else {
			fmt.Printf("✅ Cliente guardado con ID: %d\n", cliente.ID)
			fmt.Printf("📧 Email: %s\n", cliente.Email)
			fmt.Printf("📞 Teléfono: %s\n", cliente.Telefono)
		}
	}

	// Demo 7: Estadísticas del sistema
	fmt.Println("\n7️⃣ ESTADÍSTICAS DEL SISTEMA")
	fmt.Println(strings.Repeat("-", 40))
	
	estadisticas, err := db.EstadisticasFacturas()
	if err != nil {
		fmt.Printf("❌ Error obteniendo estadísticas: %v\n", err)
	} else {
		fmt.Printf("📊 RESUMEN GENERAL\n")
		fmt.Printf("   📋 Total facturas: %v\n", estadisticas["total_facturas"])
		fmt.Printf("   💰 Total facturado: $%.2f\n", estadisticas["total_facturado"])
		
		if porEstado, ok := estadisticas["por_estado"].(map[string]int); ok {
			fmt.Printf("\n📈 FACTURAS POR ESTADO:\n")
			for estado, cantidad := range porEstado {
				fmt.Printf("   %s: %d facturas\n", estado, cantidad)
			}
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Demo de base de datos completado exitosamente")
	fmt.Println("💡 Beneficios del sistema de persistencia:")
	fmt.Println("   • Almacenamiento permanente de facturas")
	fmt.Println("   • Seguimiento de estados (BORRADOR → AUTORIZADA)")
	fmt.Println("   • Gestión de clientes recurrentes")
	fmt.Println("   • Consultas rápidas con índices optimizados")
	fmt.Println("   • Estadísticas en tiempo real")
	fmt.Println("   • Integridad transaccional garantizada")
	fmt.Println(strings.Repeat("=", 60))
}

// DemoAPIDatabase demuestra el uso de la API con base de datos
func DemoAPIDatabase() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🌐 DEMO API CON BASE DE DATOS")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Println("\n📡 ENDPOINTS DISPONIBLES:")
	fmt.Println(strings.Repeat("-", 40))
	
	endpoints := []struct {
		metodo      string
		ruta        string
		descripcion string
	}{
		{"POST", "/api/facturas/db", "Crear factura en base de datos"},
		{"GET", "/api/facturas/db/list", "Listar facturas paginadas"},
		{"GET", "/api/facturas/db/{id}", "Obtener factura por ID"},
		{"PUT", "/api/facturas/db/{id}/estado", "Actualizar estado de factura"},
		{"GET", "/api/estadisticas", "Obtener estadísticas"},
		{"POST", "/api/clientes", "Crear/actualizar cliente"},
		{"GET", "/api/clientes/buscar?cedula={cedula}", "Buscar cliente por cédula"},
	}
	
	for _, endpoint := range endpoints {
		fmt.Printf("   %s %-35s - %s\n", endpoint.metodo, endpoint.ruta, endpoint.descripcion)
	}
	
	fmt.Println("\n🔧 EJEMPLOS DE USO:")
	fmt.Println(strings.Repeat("-", 40))
	
	fmt.Println("\n1. Crear factura:")
	fmt.Println(`curl -X POST http://localhost:8080/api/facturas/db \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "EMPRESA DEMO S.A.",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "DEMO001",
        "descripcion": "Producto Demo",
        "cantidad": 1,
        "precioUnitario": 100.00
      }
    ]
  }'`)
	
	fmt.Println("\n2. Listar facturas:")
	fmt.Println(`curl "http://localhost:8080/api/facturas/db/list?limit=5&offset=0"`)
	
	fmt.Println("\n3. Obtener factura específica:")
	fmt.Println(`curl "http://localhost:8080/api/facturas/db/1?includeXML=true"`)
	
	fmt.Println("\n4. Actualizar estado de factura:")
	fmt.Println(`curl -X PUT http://localhost:8080/api/facturas/db/1/estado \
  -H "Content-Type: application/json" \
  -d '{
    "estado": "AUTORIZADA",
    "numero_autorizacion": "1234567890123456789",
    "observaciones_sri": "Autorizada correctamente"
  }'`)
	
	fmt.Println("\n5. Obtener estadísticas:")
	fmt.Println(`curl http://localhost:8080/api/estadisticas`)
	
	fmt.Println("\n6. Crear cliente:")
	fmt.Println(`curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cedula": "1713175071",
    "nombre": "JUAN PEREZ",
    "direccion": "Av. Principal 123",
    "telefono": "0987654321",
    "email": "juan@ejemplo.com",
    "tipoCliente": "PERSONA_NATURAL"
  }'`)
	
	fmt.Println("\n7. Buscar cliente:")
	fmt.Println(`curl "http://localhost:8080/api/clientes/buscar?cedula=1713175071"`)
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Para probar la API con base de datos:")
	fmt.Println("   1. Ejecutar: go run main.go test_validaciones.go api")
	fmt.Println("   2. Usar los endpoints mostrados arriba")
	fmt.Println("   3. Los datos se guardan en database/facturacion.db")
	fmt.Println(strings.Repeat("=", 60))
}

// min helper function para Go < 1.21
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}