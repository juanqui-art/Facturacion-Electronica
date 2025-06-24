// Package database implementa persistencia de datos para el sistema de facturación
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite
	"go-facturacion-sri/models"
)

// Database estructura para manejar la base de datos
type Database struct {
	db *sql.DB
}

// FacturaDB estructura de factura para base de datos
type FacturaDB struct {
	ID                   int       `json:"id"`
	NumeroFactura        string    `json:"numeroFactura"`
	ClaveAcceso          string    `json:"claveAcceso"`
	FechaEmision         time.Time `json:"fechaEmision"`
	ClienteNombre        string    `json:"clienteNombre"`
	ClienteCedula        string    `json:"clienteCedula"`
	ClienteDireccion     string    `json:"clienteDireccion"`
	ClienteTelefono      string    `json:"clienteTelefono"`
	ClienteEmail         string    `json:"clienteEmail"`
	Subtotal             float64   `json:"subtotal"`
	IVA                  float64   `json:"iva"`
	Total                float64   `json:"total"`
	Estado               string    `json:"estado"` // BORRADOR, ENVIADA, AUTORIZADA, RECHAZADA
	NumeroAutorizacion   string    `json:"numeroAutorizacion"`
	FechaAutorizacion    *time.Time `json:"fechaAutorizacion"`
	XMLOriginal          string    `json:"xmlOriginal"`
	XMLAutorizado        string    `json:"xmlAutorizado"`
	ObservacionesSRI     string    `json:"observacionesSRI"`
	Ambiente             string    `json:"ambiente"` // PRUEBAS, PRODUCCION
	TipoEmision          string    `json:"tipoEmision"`
	FechaCreacion        time.Time `json:"fechaCreacion"`
	FechaActualizacion   time.Time `json:"fechaActualizacion"`
}

// ProductoDB estructura de producto para base de datos
type ProductoDB struct {
	ID                int     `json:"id"`
	FacturaID         int     `json:"facturaId"`
	Codigo            string  `json:"codigo"`
	CodigoPrincipal   string  `json:"codigoPrincipal"`
	CodigoAuxiliar    string  `json:"codigoAuxiliar"`
	Descripcion       string  `json:"descripcion"`
	UnidadMedida      string  `json:"unidadMedida"`
	Cantidad          float64 `json:"cantidad"`
	PrecioUnitario    float64 `json:"precioUnitario"`
	Descuento         float64 `json:"descuento"`
	PrecioTotalSinIva float64 `json:"precioTotalSinIva"`
	PrecioTotal       float64 `json:"precioTotal"`
	IVA               float64 `json:"iva"`
}

// ClienteDB estructura de cliente para base de datos
type ClienteDB struct {
	ID            int       `json:"id"`
	Cedula        string    `json:"cedula"`
	Nombre        string    `json:"nombre"`
	Direccion     string    `json:"direccion"`
	Telefono      string    `json:"telefono"`
	Email         string    `json:"email"`
	TipoCliente   string    `json:"tipoCliente"` // PERSONA_NATURAL, EMPRESA
	FechaCreacion time.Time `json:"fechaCreacion"`
	Activo        bool      `json:"activo"`
}

// ConfigDB estructura de configuración para base de datos
type ConfigDB struct {
	ID     int    `json:"id"`
	Clave  string `json:"clave"`
	Valor  string `json:"valor"`
	Tipo   string `json:"tipo"` // STRING, NUMBER, BOOLEAN, JSON
	Activo bool   `json:"activo"`
}

// New crea una nueva instancia de base de datos
func New(dbPath string) (*Database, error) {
	// Crear directorio si no existe
	if err := os.MkdirAll("database", 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio database: %v", err)
	}

	// Abrir conexión a SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo base de datos: %v", err)
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	database := &Database{db: db}

	// Crear tablas si no existen
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("error creando tablas: %v", err)
	}

	log.Printf("✅ Base de datos inicializada: %s", dbPath)
	return database, nil
}

// Close cierra la conexión a la base de datos
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// createTables crea las tablas necesarias
func (d *Database) createTables() error {
	// Tabla de facturas
	facturaSQL := `
	CREATE TABLE IF NOT EXISTS facturas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		numero_factura TEXT NOT NULL UNIQUE,
		clave_acceso TEXT NOT NULL UNIQUE,
		fecha_emision DATETIME NOT NULL,
		cliente_nombre TEXT NOT NULL,
		cliente_cedula TEXT NOT NULL,
		cliente_direccion TEXT,
		cliente_telefono TEXT,
		cliente_email TEXT,
		subtotal REAL NOT NULL,
		iva REAL NOT NULL,
		total REAL NOT NULL,
		estado TEXT NOT NULL DEFAULT 'BORRADOR',
		numero_autorizacion TEXT,
		fecha_autorizacion DATETIME,
		xml_original TEXT,
		xml_autorizado TEXT,
		observaciones_sri TEXT,
		ambiente TEXT NOT NULL DEFAULT 'PRUEBAS',
		tipo_emision TEXT NOT NULL DEFAULT 'NORMAL',
		fecha_creacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		fecha_actualizacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	// Tabla de productos
	productoSQL := `
	CREATE TABLE IF NOT EXISTS productos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		factura_id INTEGER NOT NULL,
		codigo TEXT NOT NULL,
		codigo_principal TEXT,
		codigo_auxiliar TEXT,
		descripcion TEXT NOT NULL,
		unidad_medida TEXT DEFAULT 'UNI',
		cantidad REAL NOT NULL,
		precio_unitario REAL NOT NULL,
		descuento REAL DEFAULT 0,
		precio_total_sin_iva REAL NOT NULL,
		precio_total REAL NOT NULL,
		iva REAL NOT NULL,
		FOREIGN KEY (factura_id) REFERENCES facturas (id) ON DELETE CASCADE
	);`

	// Tabla de clientes
	clienteSQL := `
	CREATE TABLE IF NOT EXISTS clientes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cedula TEXT NOT NULL UNIQUE,
		nombre TEXT NOT NULL,
		direccion TEXT,
		telefono TEXT,
		email TEXT,
		tipo_cliente TEXT NOT NULL DEFAULT 'PERSONA_NATURAL',
		fecha_creacion DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		activo BOOLEAN NOT NULL DEFAULT 1
	);`

	// Tabla de configuración
	configSQL := `
	CREATE TABLE IF NOT EXISTS configuracion (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		clave TEXT NOT NULL UNIQUE,
		valor TEXT NOT NULL,
		tipo TEXT NOT NULL DEFAULT 'STRING',
		activo BOOLEAN NOT NULL DEFAULT 1
	);`

	// Índices para mejorar performance
	indicesSQL := []string{
		"CREATE INDEX IF NOT EXISTS idx_facturas_numero ON facturas(numero_factura);",
		"CREATE INDEX IF NOT EXISTS idx_facturas_clave ON facturas(clave_acceso);",
		"CREATE INDEX IF NOT EXISTS idx_facturas_cliente ON facturas(cliente_cedula);",
		"CREATE INDEX IF NOT EXISTS idx_facturas_fecha ON facturas(fecha_emision);",
		"CREATE INDEX IF NOT EXISTS idx_facturas_estado ON facturas(estado);",
		"CREATE INDEX IF NOT EXISTS idx_productos_factura ON productos(factura_id);",
		"CREATE INDEX IF NOT EXISTS idx_clientes_cedula ON clientes(cedula);",
	}

	// Ejecutar creación de tablas
	tables := []string{facturaSQL, productoSQL, clienteSQL, configSQL}
	for _, table := range tables {
		if _, err := d.db.Exec(table); err != nil {
			return fmt.Errorf("error creando tabla: %v", err)
		}
	}

	// Ejecutar creación de índices
	for _, index := range indicesSQL {
		if _, err := d.db.Exec(index); err != nil {
			return fmt.Errorf("error creando índice: %v", err)
		}
	}

	return nil
}

// GuardarFactura guarda una factura completa en la base de datos
func (d *Database) GuardarFactura(factura models.Factura, claveAcceso string, productos []models.ProductoInput) (*FacturaDB, error) {
	// Iniciar transacción
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %v", err)
	}
	defer tx.Rollback()

	// Generar número de factura
	numeroFactura, err := d.generarNumeroFactura(tx)
	if err != nil {
		return nil, fmt.Errorf("error generando número de factura: %v", err)
	}

	// Generar XML
	xmlOriginal, err := factura.GenerarXML()
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %v", err)
	}

	// Insertar factura
	facturaSQL := `
		INSERT INTO facturas (
			numero_factura, clave_acceso, fecha_emision, cliente_nombre, cliente_cedula,
			cliente_direccion, cliente_telefono, cliente_email, subtotal, iva, total,
			estado, xml_original, ambiente, tipo_emision
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(facturaSQL,
		numeroFactura,
		claveAcceso,
		time.Now(),
		factura.InfoFactura.RazonSocialComprador,
		factura.InfoFactura.IdentificacionComprador,
		factura.InfoFactura.DirEstablecimiento, // Usamos dirección del establecimiento como placeholder
		"", // Teléfono - por implementar
		"", // Email - por implementar
		factura.InfoFactura.TotalSinImpuestos,
		factura.InfoFactura.ImporteTotal-factura.InfoFactura.TotalSinImpuestos,
		factura.InfoFactura.ImporteTotal,
		"BORRADOR",
		string(xmlOriginal),
		"PRUEBAS",
		"NORMAL",
	)
	if err != nil {
		return nil, fmt.Errorf("error insertando factura: %v", err)
	}

	facturaID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ID de factura: %v", err)
	}

	// Insertar productos
	productoSQL := `
		INSERT INTO productos (
			factura_id, codigo, descripcion, cantidad, precio_unitario,
			descuento, precio_total_sin_iva, precio_total, iva
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for i, detalle := range factura.Detalles {
		producto := productos[i] // Correspondencia con productos originales
		
		_, err := tx.Exec(productoSQL,
			facturaID,
			producto.Codigo,
			detalle.Descripcion,
			detalle.Cantidad,
			detalle.PrecioUnitario,
			detalle.Descuento,
			detalle.PrecioTotalSinImpuesto,
			detalle.PrecioTotalSinImpuesto, // Por ahora sin IVA en detalle
			0, // IVA por producto - por implementar
		)
		if err != nil {
			return nil, fmt.Errorf("error insertando producto %d: %v", i+1, err)
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %v", err)
	}

	// Retornar factura creada
	return d.ObtenerFacturaPorID(int(facturaID))
}

// generarNumeroFactura genera un número de factura secuencial
func (d *Database) generarNumeroFactura(tx *sql.Tx) (string, error) {
	var ultimoNumero int
	err := tx.QueryRow("SELECT COALESCE(MAX(CAST(SUBSTR(numero_factura, 5) AS INTEGER)), 0) FROM facturas WHERE numero_factura LIKE 'FAC-%'").Scan(&ultimoNumero)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	
	return fmt.Sprintf("FAC-%06d", ultimoNumero+1), nil
}

// ObtenerFacturaPorID obtiene una factura por su ID
func (d *Database) ObtenerFacturaPorID(id int) (*FacturaDB, error) {
	query := `
		SELECT id, numero_factura, clave_acceso, fecha_emision, cliente_nombre, cliente_cedula,
			   cliente_direccion, cliente_telefono, cliente_email, subtotal, iva, total,
			   estado, numero_autorizacion, fecha_autorizacion, xml_original, xml_autorizado,
			   observaciones_sri, ambiente, tipo_emision, fecha_creacion, fecha_actualizacion
		FROM facturas WHERE id = ?`

	row := d.db.QueryRow(query, id)
	
	factura := &FacturaDB{}
	var fechaAutorizacion sql.NullTime
	var numeroAutorizacion, xmlAutorizado, observacionesSRI sql.NullString
	var clienteDireccion, clienteTelefono, clienteEmail sql.NullString
	
	err := row.Scan(
		&factura.ID, &factura.NumeroFactura, &factura.ClaveAcceso, &factura.FechaEmision,
		&factura.ClienteNombre, &factura.ClienteCedula, &clienteDireccion,
		&clienteTelefono, &clienteEmail, &factura.Subtotal, &factura.IVA,
		&factura.Total, &factura.Estado, &numeroAutorizacion, &fechaAutorizacion,
		&factura.XMLOriginal, &xmlAutorizado, &observacionesSRI,
		&factura.Ambiente, &factura.TipoEmision, &factura.FechaCreacion, &factura.FechaActualizacion,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("factura con ID %d no encontrada", id)
		}
		return nil, fmt.Errorf("error obteniendo factura: %v", err)
	}
	
	// Asignar valores nullable
	if fechaAutorizacion.Valid {
		factura.FechaAutorizacion = &fechaAutorizacion.Time
	}
	if numeroAutorizacion.Valid {
		factura.NumeroAutorizacion = numeroAutorizacion.String
	}
	if xmlAutorizado.Valid {
		factura.XMLAutorizado = xmlAutorizado.String
	}
	if observacionesSRI.Valid {
		factura.ObservacionesSRI = observacionesSRI.String
	}
	if clienteDireccion.Valid {
		factura.ClienteDireccion = clienteDireccion.String
	}
	if clienteTelefono.Valid {
		factura.ClienteTelefono = clienteTelefono.String
	}
	if clienteEmail.Valid {
		factura.ClienteEmail = clienteEmail.String
	}
	
	return factura, nil
}

// ObtenerFacturaPorNumero obtiene una factura por su número
func (d *Database) ObtenerFacturaPorNumero(numero string) (*FacturaDB, error) {
	query := `
		SELECT id, numero_factura, clave_acceso, fecha_emision, cliente_nombre, cliente_cedula,
			   cliente_direccion, cliente_telefono, cliente_email, subtotal, iva, total,
			   estado, numero_autorizacion, fecha_autorizacion, xml_original, xml_autorizado,
			   observaciones_sri, ambiente, tipo_emision, fecha_creacion, fecha_actualizacion
		FROM facturas WHERE numero_factura = ?`

	row := d.db.QueryRow(query, numero)
	
	factura := &FacturaDB{}
	var fechaAutorizacion sql.NullTime
	var numeroAutorizacion, xmlAutorizado, observacionesSRI sql.NullString
	var clienteDireccion, clienteTelefono, clienteEmail sql.NullString
	
	err := row.Scan(
		&factura.ID, &factura.NumeroFactura, &factura.ClaveAcceso, &factura.FechaEmision,
		&factura.ClienteNombre, &factura.ClienteCedula, &clienteDireccion,
		&clienteTelefono, &clienteEmail, &factura.Subtotal, &factura.IVA,
		&factura.Total, &factura.Estado, &numeroAutorizacion, &fechaAutorizacion,
		&factura.XMLOriginal, &xmlAutorizado, &observacionesSRI,
		&factura.Ambiente, &factura.TipoEmision, &factura.FechaCreacion, &factura.FechaActualizacion,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("factura con número %s no encontrada", numero)
		}
		return nil, fmt.Errorf("error obteniendo factura: %v", err)
	}
	
	// Asignar valores nullable
	if fechaAutorizacion.Valid {
		factura.FechaAutorizacion = &fechaAutorizacion.Time
	}
	if numeroAutorizacion.Valid {
		factura.NumeroAutorizacion = numeroAutorizacion.String
	}
	if xmlAutorizado.Valid {
		factura.XMLAutorizado = xmlAutorizado.String
	}
	if observacionesSRI.Valid {
		factura.ObservacionesSRI = observacionesSRI.String
	}
	if clienteDireccion.Valid {
		factura.ClienteDireccion = clienteDireccion.String
	}
	if clienteTelefono.Valid {
		factura.ClienteTelefono = clienteTelefono.String
	}
	if clienteEmail.Valid {
		factura.ClienteEmail = clienteEmail.String
	}
	
	return factura, nil
}

// ListarFacturas obtiene una lista paginada de facturas
func (d *Database) ListarFacturas(limite, offset int) ([]*FacturaDB, error) {
	query := `
		SELECT id, numero_factura, clave_acceso, fecha_emision, cliente_nombre, cliente_cedula,
			   subtotal, iva, total, estado, numero_autorizacion, ambiente
		FROM facturas 
		ORDER BY fecha_creacion DESC 
		LIMIT ? OFFSET ?`

	rows, err := d.db.Query(query, limite, offset)
	if err != nil {
		return nil, fmt.Errorf("error listando facturas: %v", err)
	}
	defer rows.Close()

	var facturas []*FacturaDB
	
	for rows.Next() {
		factura := &FacturaDB{}
		var numeroAutorizacion sql.NullString
		
		err := rows.Scan(
			&factura.ID, &factura.NumeroFactura, &factura.ClaveAcceso, &factura.FechaEmision,
			&factura.ClienteNombre, &factura.ClienteCedula, &factura.Subtotal, &factura.IVA,
			&factura.Total, &factura.Estado, &numeroAutorizacion, &factura.Ambiente,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando factura: %v", err)
		}
		
		if numeroAutorizacion.Valid {
			factura.NumeroAutorizacion = numeroAutorizacion.String
		}
		
		facturas = append(facturas, factura)
	}
	
	return facturas, nil
}

// ActualizarEstadoFactura actualiza el estado de una factura
func (d *Database) ActualizarEstadoFactura(id int, estado string, numeroAutorizacion string, xmlAutorizado string, observaciones string) error {
	query := `
		UPDATE facturas 
		SET estado = ?, numero_autorizacion = ?, fecha_autorizacion = ?, 
		    xml_autorizado = ?, observaciones_sri = ?, fecha_actualizacion = CURRENT_TIMESTAMP
		WHERE id = ?`

	var fechaAutorizacion *time.Time
	if estado == "AUTORIZADA" {
		now := time.Now()
		fechaAutorizacion = &now
	}

	_, err := d.db.Exec(query, estado, numeroAutorizacion, fechaAutorizacion, xmlAutorizado, observaciones, id)
	if err != nil {
		return fmt.Errorf("error actualizando estado de factura: %v", err)
	}

	return nil
}

// ObtenerProductosPorFactura obtiene los productos de una factura
func (d *Database) ObtenerProductosPorFactura(facturaID int) ([]*ProductoDB, error) {
	query := `
		SELECT id, factura_id, codigo, codigo_principal, codigo_auxiliar, descripcion,
			   unidad_medida, cantidad, precio_unitario, descuento, precio_total_sin_iva,
			   precio_total, iva
		FROM productos WHERE factura_id = ? ORDER BY id`

	rows, err := d.db.Query(query, facturaID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo productos: %v", err)
	}
	defer rows.Close()

	var productos []*ProductoDB
	
	for rows.Next() {
		producto := &ProductoDB{}
		var codigoPrincipal, codigoAuxiliar sql.NullString
		
		err := rows.Scan(
			&producto.ID, &producto.FacturaID, &producto.Codigo, &codigoPrincipal,
			&codigoAuxiliar, &producto.Descripcion, &producto.UnidadMedida,
			&producto.Cantidad, &producto.PrecioUnitario, &producto.Descuento,
			&producto.PrecioTotalSinIva, &producto.PrecioTotal, &producto.IVA,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando producto: %v", err)
		}
		
		if codigoPrincipal.Valid {
			producto.CodigoPrincipal = codigoPrincipal.String
		}
		if codigoAuxiliar.Valid {
			producto.CodigoAuxiliar = codigoAuxiliar.String
		}
		
		productos = append(productos, producto)
	}
	
	return productos, nil
}

// GuardarCliente guarda un cliente en la base de datos
func (d *Database) GuardarCliente(cliente *ClienteDB) (*ClienteDB, error) {
	query := `
		INSERT OR REPLACE INTO clientes (cedula, nombre, direccion, telefono, email, tipo_cliente)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, cliente.Cedula, cliente.Nombre, cliente.Direccion, 
		cliente.Telefono, cliente.Email, cliente.TipoCliente)
	if err != nil {
		return nil, fmt.Errorf("error guardando cliente: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ID de cliente: %v", err)
	}

	return d.ObtenerClientePorID(int(id))
}

// ObtenerClientePorID obtiene un cliente por su ID
func (d *Database) ObtenerClientePorID(id int) (*ClienteDB, error) {
	query := `
		SELECT id, cedula, nombre, direccion, telefono, email, tipo_cliente, fecha_creacion, activo
		FROM clientes WHERE id = ?`

	row := d.db.QueryRow(query, id)
	
	cliente := &ClienteDB{}
	var direccion, telefono, email sql.NullString
	
	err := row.Scan(&cliente.ID, &cliente.Cedula, &cliente.Nombre, &direccion,
		&telefono, &email, &cliente.TipoCliente, &cliente.FechaCreacion, &cliente.Activo)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente con ID %d no encontrado", id)
		}
		return nil, fmt.Errorf("error obteniendo cliente: %v", err)
	}
	
	if direccion.Valid {
		cliente.Direccion = direccion.String
	}
	if telefono.Valid {
		cliente.Telefono = telefono.String
	}
	if email.Valid {
		cliente.Email = email.String
	}
	
	return cliente, nil
}

// ObtenerClientePorCedula obtiene un cliente por su cédula
func (d *Database) ObtenerClientePorCedula(cedula string) (*ClienteDB, error) {
	query := `
		SELECT id, cedula, nombre, direccion, telefono, email, tipo_cliente, fecha_creacion, activo
		FROM clientes WHERE cedula = ? AND activo = 1`

	row := d.db.QueryRow(query, cedula)
	
	cliente := &ClienteDB{}
	var direccion, telefono, email sql.NullString
	
	err := row.Scan(&cliente.ID, &cliente.Cedula, &cliente.Nombre, &direccion,
		&telefono, &email, &cliente.TipoCliente, &cliente.FechaCreacion, &cliente.Activo)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente con cédula %s no encontrado", cedula)
		}
		return nil, fmt.Errorf("error obteniendo cliente: %v", err)
	}
	
	if direccion.Valid {
		cliente.Direccion = direccion.String
	}
	if telefono.Valid {
		cliente.Telefono = telefono.String
	}
	if email.Valid {
		cliente.Email = email.String
	}
	
	return cliente, nil
}

// EstadisticasFacturas obtiene estadísticas básicas de facturas
func (d *Database) EstadisticasFacturas() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total de facturas
	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM facturas").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo total de facturas: %v", err)
	}
	stats["total_facturas"] = total
	
	// Facturas por estado
	rows, err := d.db.Query("SELECT estado, COUNT(*) FROM facturas GROUP BY estado")
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas por estado: %v", err)
	}
	defer rows.Close()
	
	estadosMap := make(map[string]int)
	for rows.Next() {
		var estado string
		var count int
		if err := rows.Scan(&estado, &count); err != nil {
			return nil, fmt.Errorf("error escaneando estadísticas: %v", err)
		}
		estadosMap[estado] = count
	}
	stats["por_estado"] = estadosMap
	
	// Total facturado
	var totalFacturado sql.NullFloat64
	err = d.db.QueryRow("SELECT SUM(total) FROM facturas WHERE estado = 'AUTORIZADA'").Scan(&totalFacturado)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo total facturado: %v", err)
	}
	
	if totalFacturado.Valid {
		stats["total_facturado"] = totalFacturado.Float64
	} else {
		stats["total_facturado"] = 0.0
	}
	
	return stats, nil
}