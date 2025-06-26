// Package database - Funciones adicionales para CRUD completo
package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// ActualizarCliente actualiza un cliente existente
func (d *Database) ActualizarCliente(cliente *ClienteDB) (*ClienteDB, error) {
	// Registrar auditoría - obtener datos antes
	clienteAntes, err := d.ObtenerClientePorID(cliente.ID)
	if err != nil {
		return nil, fmt.Errorf("cliente no encontrado: %v", err)
	}

	query := `
		UPDATE clientes 
		SET cedula = ?, nombre = ?, direccion = ?, telefono = ?, email = ?, tipo_cliente = ?
		WHERE id = ? AND activo = 1`

	_, err = d.db.Exec(query, cliente.Cedula, cliente.Nombre, cliente.Direccion, 
		cliente.Telefono, cliente.Email, cliente.TipoCliente, cliente.ID)
	if err != nil {
		return nil, fmt.Errorf("error actualizando cliente: %v", err)
	}

	// Obtener cliente actualizado
	clienteActualizado, err := d.ObtenerClientePorID(cliente.ID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo cliente actualizado: %v", err)
	}

	// Registrar auditoría
	datosAntes, _ := json.Marshal(clienteAntes)
	datosDespues, _ := json.Marshal(clienteActualizado)
	audit := &AuditLogDB{
		Tabla:        "clientes",
		RegistroID:   cliente.ID,
		Operacion:    "UPDATE",
		Usuario:      "system", // TODO: obtener usuario real
		DatosAntes:   string(datosAntes),
		DatosDespues: string(datosDespues),
	}
	d.RegistrarAuditoria(audit)

	return clienteActualizado, nil
}

// DesactivarCliente marca un cliente como inactivo (soft delete)
func (d *Database) DesactivarCliente(id int) error {
	// Verificar que el cliente existe
	cliente, err := d.ObtenerClientePorID(id)
	if err != nil {
		return fmt.Errorf("cliente no encontrado: %v", err)
	}

	query := `UPDATE clientes SET activo = 0 WHERE id = ?`
	_, err = d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error desactivando cliente: %v", err)
	}

	// Registrar auditoría
	datosAntes, _ := json.Marshal(cliente)
	clienteDesactivado := *cliente
	clienteDesactivado.Activo = false
	datosDespues, _ := json.Marshal(clienteDesactivado)
	
	audit := &AuditLogDB{
		Tabla:        "clientes",
		RegistroID:   id,
		Operacion:    "DEACTIVATE",
		Usuario:      "system", // TODO: obtener usuario real
		DatosAntes:   string(datosAntes),
		DatosDespues: string(datosDespues),
	}
	d.RegistrarAuditoria(audit)

	return nil
}

// EliminarCliente elimina completamente un cliente (hard delete)
func (d *Database) EliminarCliente(id int) error {
	// Obtener cliente para auditoría
	cliente, err := d.ObtenerClientePorID(id)
	if err != nil {
		return fmt.Errorf("cliente no encontrado: %v", err)
	}

	query := `DELETE FROM clientes WHERE id = ?`
	_, err = d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando cliente: %v", err)
	}

	// Registrar auditoría
	datosAntes, _ := json.Marshal(cliente)
	audit := &AuditLogDB{
		Tabla:      "clientes",
		RegistroID: id,
		Operacion:  "DELETE",
		Usuario:    "system", // TODO: obtener usuario real
		DatosAntes: string(datosAntes),
	}
	d.RegistrarAuditoria(audit)

	return nil
}

// ListarFacturasPorCliente obtiene facturas de un cliente específico
func (d *Database) ListarFacturasPorCliente(cedula string, limite, offset int) ([]*FacturaDB, error) {
	query := `
		SELECT id, numero_factura, clave_acceso, cliente_cedula, cliente_nombre,
			   subtotal, iva, total, estado, fecha_emision, fecha_creacion,
			   numero_autorizacion, fecha_autorizacion, observaciones_sri, xml_original, xml_autorizado
		FROM facturas 
		WHERE cliente_cedula = ?
		ORDER BY fecha_creacion DESC 
		LIMIT ? OFFSET ?`

	rows, err := d.db.Query(query, cedula, limite, offset)
	if err != nil {
		return nil, fmt.Errorf("error listando facturas por cliente: %v", err)
	}
	defer rows.Close()

	var facturas []*FacturaDB

	for rows.Next() {
		factura := &FacturaDB{}
		var numeroAutorizacion, fechaAutorizacion, observacionesSRI, xmlOriginal, xmlAutorizado sql.NullString

		err := rows.Scan(
			&factura.ID, &factura.NumeroFactura, &factura.ClaveAcceso,
			&factura.ClienteCedula, &factura.ClienteNombre,
			&factura.Subtotal, &factura.IVA, &factura.Total,
			&factura.Estado, &factura.FechaEmision, &factura.FechaCreacion,
			&numeroAutorizacion, &fechaAutorizacion, &observacionesSRI,
			&xmlOriginal, &xmlAutorizado,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando factura: %v", err)
		}

		// Manejar campos nullables
		if numeroAutorizacion.Valid {
			factura.NumeroAutorizacion = numeroAutorizacion.String
		}
		if fechaAutorizacion.Valid {
			// Parse fecha autorización
			if fechaTime, err := time.Parse(time.RFC3339, fechaAutorizacion.String); err == nil {
				factura.FechaAutorizacion = &fechaTime
			}
		}
		if observacionesSRI.Valid {
			factura.ObservacionesSRI = observacionesSRI.String
		}
		if xmlOriginal.Valid {
			factura.XMLOriginal = xmlOriginal.String
		}
		if xmlAutorizado.Valid {
			factura.XMLAutorizado = xmlAutorizado.String
		}

		facturas = append(facturas, factura)
	}

	return facturas, nil
}

// ActualizarFactura actualiza una factura completa (solo en estado BORRADOR)
func (d *Database) ActualizarFactura(id int, clienteCedula, clienteNombre string, productos []ProductoDB, observaciones string) (*FacturaDB, error) {
	// Verificar que la factura existe y está en estado BORRADOR
	facturaAntes, err := d.ObtenerFacturaPorID(id)
	if err != nil {
		return nil, fmt.Errorf("factura no encontrada: %v", err)
	}

	if facturaAntes.Estado != "BORRADOR" {
		return nil, fmt.Errorf("solo se pueden actualizar facturas en estado BORRADOR")
	}

	// Calcular totales
	var subtotal, total float64
	for _, producto := range productos {
		subtotalProducto := producto.Cantidad * producto.PrecioUnitario
		subtotal += subtotalProducto
	}
	total = subtotal // Para simplificar, sin IVA por ahora

	// Iniciar transacción
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error iniciando transacción: %v", err)
	}
	defer tx.Rollback()

	// Actualizar factura
	queryFactura := `
		UPDATE facturas 
		SET cliente_cedula = ?, cliente_nombre = ?, subtotal = ?, total = ?
		WHERE id = ?`

	_, err = tx.Exec(queryFactura, clienteCedula, clienteNombre, subtotal, total, id)
	if err != nil {
		return nil, fmt.Errorf("error actualizando factura: %v", err)
	}

	// Eliminar productos existentes
	_, err = tx.Exec("DELETE FROM productos WHERE factura_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("error eliminando productos existentes: %v", err)
	}

	// Insertar nuevos productos
	queryProducto := `
		INSERT INTO productos (factura_id, codigo, descripcion, cantidad, precio_unitario, descuento, subtotal)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, producto := range productos {
		subtotalProducto := producto.Cantidad * producto.PrecioUnitario
		_, err = tx.Exec(queryProducto, id, producto.Codigo, producto.Descripcion,
			producto.Cantidad, producto.PrecioUnitario, producto.Descuento, subtotalProducto)
		if err != nil {
			return nil, fmt.Errorf("error insertando producto: %v", err)
		}
	}

	// Confirmar transacción
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error confirmando transacción: %v", err)
	}

	// Obtener factura actualizada
	facturaActualizada, err := d.ObtenerFacturaPorID(id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo factura actualizada: %v", err)
	}

	// Registrar auditoría
	datosAntes, _ := json.Marshal(facturaAntes)
	datosDespues, _ := json.Marshal(facturaActualizada)
	audit := &AuditLogDB{
		Tabla:        "facturas",
		RegistroID:   id,
		Operacion:    "UPDATE",
		Usuario:      "system", // TODO: obtener usuario real
		DatosAntes:   string(datosAntes),
		DatosDespues: string(datosDespues),
	}
	d.RegistrarAuditoria(audit)

	return facturaActualizada, nil
}

// EliminarFactura elimina completamente una factura y sus productos
func (d *Database) EliminarFactura(id int) error {
	// Obtener factura para auditoría
	factura, err := d.ObtenerFacturaPorID(id)
	if err != nil {
		return fmt.Errorf("factura no encontrada: %v", err)
	}

	// Iniciar transacción
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %v", err)
	}
	defer tx.Rollback()

	// Eliminar productos primero (FK constraint)
	_, err = tx.Exec("DELETE FROM productos WHERE factura_id = ?", id)
	if err != nil {
		return fmt.Errorf("error eliminando productos: %v", err)
	}

	// Eliminar factura
	_, err = tx.Exec("DELETE FROM facturas WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error eliminando factura: %v", err)
	}

	// Confirmar transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error confirmando transacción: %v", err)
	}

	// Registrar auditoría
	datosAntes, _ := json.Marshal(factura)
	audit := &AuditLogDB{
		Tabla:      "facturas",
		RegistroID: id,
		Operacion:  "DELETE",
		Usuario:    "system", // TODO: obtener usuario real
		DatosAntes: string(datosAntes),
	}
	d.RegistrarAuditoria(audit)

	return nil
}