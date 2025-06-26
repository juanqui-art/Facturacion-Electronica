// Package api Handlers adicionales para base de datos en la API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-facturacion-sri/database"
	"go-facturacion-sri/factory"
	"go-facturacion-sri/models"
	"go-facturacion-sri/pdf"
	"go-facturacion-sri/sri"
)

// CrearFacturaConDB crea una factura y la guarda en base de datos
func (s *Server) CrearFacturaConDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear input JSON
	var input models.FacturaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error parseando JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Crear factura
	factura, err := factory.CrearFactura(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creando factura: %v", err), http.StatusBadRequest)
		return
	}

	// Generar clave de acceso
	claveConfig := sri.ClaveAccesoConfig{
		FechaEmision:     time.Now(),
		TipoComprobante:  sri.Factura,
		RUCEmisor:        "1792146739001", // TODO: Obtener de configuración
		Ambiente:         sri.Pruebas,
		Serie:            "001001",
		NumeroSecuencial: "000000001", // TODO: Generar secuencial automático
		TipoEmision:      sri.EmisionNormal,
	}

	claveAcceso, err := sri.GenerarClaveAcceso(claveConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generando clave de acceso: %v", err), http.StatusInternalServerError)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Guardar en base de datos
	facturaDB, err := db.GuardarFactura(factura, claveAcceso, input.Productos)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error guardando factura: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Factura creada y guardada exitosamente",
		"data": map[string]interface{}{
			"id":             facturaDB.ID,
			"numero_factura": facturaDB.NumeroFactura,
			"clave_acceso":   facturaDB.ClaveAcceso,
			"cliente_nombre": facturaDB.ClienteNombre,
			"total":          facturaDB.Total,
			"estado":         facturaDB.Estado,
			"fecha_creacion": facturaDB.FechaCreacion.Format(time.RFC3339),
		},
	}

	// Incluir XML si se solicita
	includeXML := r.URL.Query().Get("includeXML") == "true"
	if includeXML {
		response["data"].(map[string]interface{})["xml"] = facturaDB.XMLOriginal
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarFacturasDB lista facturas desde la base de datos
func (s *Server) ListarFacturasDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parámetros de paginación
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Por defecto
	offset := 0 // Por defecto

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Obtener facturas
	facturas, err := db.ListarFacturas(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listando facturas: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"facturas": facturas,
			"count":    len(facturas),
			"limit":    limit,
			"offset":   offset,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ObtenerFacturaDB obtiene una factura específica por ID
func (s *Server) ObtenerFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/facturas/db/"):]
	if idStr == "" {
		http.Error(w, "ID de factura requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de factura inválido", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Obtener factura
	factura, err := db.ObtenerFacturaPorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo factura: %v", err), http.StatusNotFound)
		return
	}

	// Obtener productos asociados
	productos, err := db.ObtenerProductosPorFactura(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo productos: %v", err), http.StatusInternalServerError)
		return
	}

	// Incluir XML si se solicita
	includeXML := r.URL.Query().Get("includeXML") == "true"

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"factura":   factura,
			"productos": productos,
		},
	}

	if includeXML {
		response["data"].(map[string]interface{})["xml_original"] = factura.XMLOriginal
		if factura.XMLAutorizado != "" {
			response["data"].(map[string]interface{})["xml_autorizado"] = factura.XMLAutorizado
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ActualizarEstadoFacturaDB actualiza el estado de una factura
func (s *Server) ActualizarEstadoFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/facturas/db/"):]
	idStr = idStr[:len(idStr)-len("/estado")] // Remover "/estado" del final

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de factura inválido", http.StatusBadRequest)
		return
	}

	// Parsear input JSON
	var input struct {
		Estado             string `json:"estado"`
		NumeroAutorizacion string `json:"numero_autorizacion"`
		XMLAutorizado      string `json:"xml_autorizado"`
		ObservacionesSRI   string `json:"observaciones_sri"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error parseando JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Actualizar estado
	err = db.ActualizarEstadoFactura(id, input.Estado, input.NumeroAutorizacion, input.XMLAutorizado, input.ObservacionesSRI)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error actualizando estado: %v", err), http.StatusInternalServerError)
		return
	}

	// Obtener factura actualizada
	factura, err := db.ObtenerFacturaPorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo factura actualizada: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Estado de factura actualizado exitosamente",
		"data":    factura,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EstadisticasDB obtiene estadísticas de facturas
func (s *Server) EstadisticasDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Obtener estadísticas
	estadisticas, err := db.EstadisticasFacturas()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo estadísticas: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data":    estadisticas,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GuardarClienteDB guarda un cliente en la base de datos
func (s *Server) GuardarClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear input JSON
	var input database.ClienteDB
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error parseando JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Guardar cliente
	cliente, err := db.GuardarCliente(&input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error guardando cliente: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Cliente guardado exitosamente",
		"data":    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BuscarClienteDB busca un cliente por cédula
func (s *Server) BuscarClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener cédula de query parameters
	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		http.Error(w, "Cédula requerida", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Buscar cliente
	cliente, err := db.ObtenerClientePorCedula(cedula)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cliente no encontrado: %v", err), http.StatusNotFound)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data":    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarClientesDB lista todos los clientes con filtros opcionales
func (s *Server) ListarClientesDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener parámetros de query
	nombre := r.URL.Query().Get("nombre")
	tipoCliente := r.URL.Query().Get("tipo")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Parsear limit y offset
	limit := 50 // Default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0 // Default
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Listar clientes
	clientes, err := db.ListarClientes(nombre, tipoCliente, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listando clientes: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"clientes": clientes,
			"count":    len(clientes),
			"limit":    limit,
			"offset":   offset,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConsultarEstadoSRI consulta el estado de una factura en el SRI
func (s *Server) ConsultarEstadoSRI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener clave de acceso de query parameters
	claveAcceso := r.URL.Query().Get("clave")
	if claveAcceso == "" {
		http.Error(w, "Clave de acceso requerida", http.StatusBadRequest)
		return
	}

	// Crear cliente SRI
	sriClient := sri.NewSOAPClient(sri.Pruebas)

	// Consultar autorización
	respuesta, err := sriClient.ConsultarAutorizacion(claveAcceso)
	if err != nil {
		// Respuesta con error, pero no falla el endpoint
		response := map[string]interface{}{
			"success":      false,
			"clave_acceso": claveAcceso,
			"error":        err.Error(),
			"estado":       "ERROR",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Respuesta exitosa
	response := map[string]interface{}{
		"success":      true,
		"clave_acceso": claveAcceso,
		"data":         respuesta,
	}

	if len(respuesta.Autorizaciones) > 0 {
		auth := respuesta.Autorizaciones[0]
		response["estado"] = auth.Estado
		response["numero_autorizacion"] = auth.NumeroAutorizacion
		response["fecha_autorizacion"] = auth.FechaAutorizacion
		
		if len(auth.Mensajes) > 0 {
			response["mensajes"] = auth.Mensajes
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EstadoGeneralSRI verifica el estado general de los servicios del SRI
func (s *Server) EstadoGeneralSRI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Crear cliente SRI
	sriClient := sri.NewSOAPClient(sri.Pruebas)

	// Intentar una consulta simple para verificar conectividad
	// Usamos una clave de acceso de prueba conocida
	claveTestPruebas := "2501202401179214673900110010010000000011234567893"
	
	response := map[string]interface{}{
		"disponible": false,
		"mensaje":    "Verificando estado...",
		"ambiente":   "PRUEBAS",
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	// Intentar consulta con timeout
	done := make(chan bool, 1)
	var sriError error

	go func() {
		_, err := sriClient.ConsultarAutorizacion(claveTestPruebas)
		sriError = err
		done <- true
	}()

	// Timeout de 5 segundos
	select {
	case <-done:
		if sriError != nil {
			response["disponible"] = false
			response["mensaje"] = fmt.Sprintf("Servicio no disponible: %v", sriError)
		} else {
			response["disponible"] = true
			response["mensaje"] = "Servicio SRI operativo"
		}
	case <-time.After(5 * time.Second):
		response["disponible"] = false
		response["mensaje"] = "Timeout en la consulta al SRI"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ObtenerAuditoriaDB obtiene registros de auditoría
func (s *Server) ObtenerAuditoriaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parámetros de consulta
	tabla := r.URL.Query().Get("tabla")
	registroIDStr := r.URL.Query().Get("registro_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // Por defecto
	offset := 0 // Por defecto

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var registros []*database.AuditLogDB

	// Determinar tipo de consulta
	if registroIDStr != "" && tabla != "" {
		// Consulta específica por tabla y registro
		registroID, err := strconv.Atoi(registroIDStr)
		if err != nil {
			http.Error(w, "ID de registro inválido", http.StatusBadRequest)
			return
		}
		
		registros, err = db.ObtenerAuditoriaPorRegistro(tabla, registroID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error obteniendo auditoría: %v", err), http.StatusInternalServerError)
			return
		}
	} else if tabla != "" {
		// Consulta por tabla
		registros, err = db.ObtenerAuditoriaPorTabla(tabla, limit, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error obteniendo auditoría: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Parámetro 'tabla' requerido", http.StatusBadRequest)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"registros": registros,
			"count":     len(registros),
			"limit":     limit,
			"offset":    offset,
			"tabla":     tabla,
		},
	}

	if registroIDStr != "" {
		response["data"].(map[string]interface{})["registro_id"] = registroIDStr
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CrearRespaldoDB crea un respaldo manual de la base de datos
func (s *Server) CrearRespaldoDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear input JSON (opcional)
	var input struct {
		Sufijo string `json:"sufijo"`
	}
	
	// Es opcional, por defecto usará timestamp
	json.NewDecoder(r.Body).Decode(&input)
	
	if input.Sufijo == "" {
		input.Sufijo = "api_request"
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Crear gestor de respaldos
	backupManager := database.NewBackupManagerDefault(db)

	// Crear respaldo manual
	err = backupManager.CrearRespaldoManual(input.Sufijo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creando respaldo: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Respaldo creado exitosamente",
		"sufijo":  input.Sufijo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarRespaldosDB lista todos los respaldos disponibles
func (s *Server) ListarRespaldosDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Crear gestor de respaldos
	backupManager := database.NewBackupManagerDefault(db)

	// Listar respaldos
	respaldos, err := backupManager.ListarRespaldos()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listando respaldos: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"respaldos": respaldos,
			"count":     len(respaldos),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ObtenerClienteDB obtiene un cliente específico por ID
func (s *Server) ObtenerClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/clientes/"):]
	if idStr == "" {
		http.Error(w, "ID de cliente requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de cliente inválido", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Obtener cliente
	cliente, err := db.ObtenerClientePorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cliente no encontrado: %v", err), http.StatusNotFound)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"data":    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ActualizarClienteDB actualiza un cliente existente
func (s *Server) ActualizarClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/clientes/"):]
	if idStr == "" {
		http.Error(w, "ID de cliente requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de cliente inválido", http.StatusBadRequest)
		return
	}

	// Parsear input JSON
	var input database.ClienteDB
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error parseando JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Verificar que el cliente existe
	_, err = db.ObtenerClientePorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cliente no encontrado: %v", err), http.StatusNotFound)
		return
	}

	// Actualizar cliente (establecer ID para update)
	input.ID = id
	cliente, err := db.ActualizarCliente(&input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error actualizando cliente: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Cliente actualizado exitosamente",
		"data":    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EliminarClienteDB elimina un cliente (soft delete)
func (s *Server) EliminarClienteDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/clientes/"):]
	if idStr == "" {
		http.Error(w, "ID de cliente requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de cliente inválido", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Verificar que el cliente existe
	cliente, err := db.ObtenerClientePorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cliente no encontrado: %v", err), http.StatusNotFound)
		return
	}

	// Verificar si el cliente tiene facturas asociadas
	facturas, err := db.ListarFacturasPorCliente(cliente.Cedula, 1, 0)
	if err == nil && len(facturas) > 0 {
		// Cliente tiene facturas, no se puede eliminar completamente
		// En su lugar, marcamos como inactivo
		err = db.DesactivarCliente(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error desactivando cliente: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"message": "Cliente desactivado (tiene facturas asociadas)",
			"action":  "desactivado",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// El cliente no tiene facturas, se puede eliminar
	err = db.EliminarCliente(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error eliminando cliente: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Cliente eliminado exitosamente",
		"action":  "eliminado",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ActualizarFacturaDB actualiza una factura completa
func (s *Server) ActualizarFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/facturas/db/"):]
	if idStr == "" {
		http.Error(w, "ID de factura requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de factura inválido", http.StatusBadRequest)
		return
	}

	// Parsear input JSON
	var input struct {
		ClienteCedula string                   `json:"clienteCedula"`
		ClienteNombre string                   `json:"clienteNombre"`
		Productos     []database.ProductoDB    `json:"productos"`
		Observaciones string                   `json:"observaciones"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error parseando JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Verificar que la factura existe y está en estado BORRADOR
	factura, err := db.ObtenerFacturaPorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Factura no encontrada: %v", err), http.StatusNotFound)
		return
	}

	if factura.Estado != "BORRADOR" {
		http.Error(w, "Solo se pueden actualizar facturas en estado BORRADOR", http.StatusBadRequest)
		return
	}

	// Actualizar factura
	facturaActualizada, err := db.ActualizarFactura(id, input.ClienteCedula, input.ClienteNombre, input.Productos, input.Observaciones)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error actualizando factura: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": "Factura actualizada exitosamente",
		"data":    facturaActualizada,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EliminarFacturaDB elimina/anula una factura
func (s *Server) EliminarFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/facturas/db/"):]
	if idStr == "" {
		http.Error(w, "ID de factura requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de factura inválido", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Obtener factura
	factura, err := db.ObtenerFacturaPorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Factura no encontrada: %v", err), http.StatusNotFound)
		return
	}

	var mensaje string
	var action string

	// Determinar acción según el estado
	switch factura.Estado {
	case "BORRADOR":
		// Se puede eliminar completamente
		err = db.EliminarFactura(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error eliminando factura: %v", err), http.StatusInternalServerError)
			return
		}
		mensaje = "Factura eliminada exitosamente"
		action = "eliminada"

	case "ENVIADA", "AUTORIZADA":
		// No se puede eliminar, solo anular
		err = db.ActualizarEstadoFactura(id, "ANULADA", "", "", "Anulada por solicitud del usuario")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error anulando factura: %v", err), http.StatusInternalServerError)
			return
		}
		mensaje = "Factura anulada exitosamente"
		action = "anulada"

	case "RECHAZADA":
		// Se puede eliminar
		err = db.EliminarFactura(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error eliminando factura: %v", err), http.StatusInternalServerError)
			return
		}
		mensaje = "Factura eliminada exitosamente"
		action = "eliminada"

	default:
		http.Error(w, "Estado de factura no válido para eliminación", http.StatusBadRequest)
		return
	}

	// Respuesta
	response := map[string]interface{}{
		"success": true,
		"message": mensaje,
		"action":  action,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerarPDFFacturaDB genera un PDF para una factura específica
func (s *Server) GenerarPDFFacturaDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Path[len("/api/facturas/db/"):]
	idStr = idStr[:len(idStr)-len("/pdf")] // Remover "/pdf" del final

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de factura inválido", http.StatusBadRequest)
		return
	}

	// Conectar a base de datos
	db, err := database.New("database/facturacion.db")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando a base de datos: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Crear generador de PDF
	pdfGenerator := pdf.NewFacturaPDFGenerator(db)

	// Validar que la factura puede generar PDF
	err = pdfGenerator.ValidarFacturaParaPDF(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validando factura para PDF: %v", err), http.StatusBadRequest)
		return
	}

	// Verificar si se quiere PDF simple
	simple := r.URL.Query().Get("simple") == "true"

	var pdfBytes []byte
	if simple {
		pdfBytes, err = pdfGenerator.GenerarFacturaSimplePDF(id)
	} else {
		pdfBytes, err = pdfGenerator.GenerarFacturaPDF(id)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Error generando PDF: %v", err), http.StatusInternalServerError)
		return
	}

	// Obtener información de la factura para el nombre del archivo
	factura, err := db.ObtenerFacturaPorID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo factura: %v", err), http.StatusInternalServerError)
		return
	}

	// Configurar headers para descarga de PDF
	filename := fmt.Sprintf("factura_%s.pdf", factura.NumeroFactura)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	// Escribir PDF
	w.Write(pdfBytes)
}
