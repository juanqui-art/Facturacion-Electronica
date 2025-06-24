# Arquitectura Multi-Tenant: Sistema SaaS de Facturación

## 🏗️ Visión General

### ¿Qué es Multi-Tenant?
Sistema donde **múltiples empresas (tenants)** comparten la **misma infraestructura** pero con **datos completamente aislados**.

### Analogía: Hotel vs Apartamentos
```
🏨 HOTEL (Multi-Tenant SaaS):
├── 🏢 Un edificio compartido (tu servidor)
├── 🗝️ Cada huésped tiene su habitación privada (empresa)
├── 🚿 Servicios compartidos (base de datos, APIs)
├── 🛎️ Recepción común (tu sistema de autenticación)
├── 🧹 Mantenimiento centralizado (tú actualizas todo)
└── 💰 Cada uno paga por su habitación (suscripción)

vs

🏠 APARTAMENTOS (On-Premise):
├── 🏘️ Muchos edificios separados (servidor por empresa)
├── 🔑 Cada uno es dueño de su casa (instalación propia)
├── 🔧 Cada uno mantiene su casa (ellos se encargan)
└── 💸 Cada uno paga todo completo (licencia + infraestructura)
```

---

## 🏢 Arquitectura de Datos

### Modelo de Aislamiento

#### **Estrategia: Shared Database, Isolated Schemas**

```sql
-- Base de datos PostgreSQL compartida
DATABASE: facturacion_saas

-- Schema por empresa
SCHEMAS:
├── empresa_1792146739001/  -- RUC como identificador
│   ├── facturas
│   ├── clientes  
│   ├── productos
│   └── configuracion
│
├── empresa_0992345678001/
│   ├── facturas
│   ├── clientes
│   ├── productos  
│   └── configuracion
│
└── shared/  -- Datos compartidos del sistema
    ├── empresas
    ├── usuarios_sistema
    ├── planes_suscripcion
    └── audit_logs
```

#### **Estructura Completa de Base de Datos**

```sql
-- ===============================
-- SCHEMA COMPARTIDO (shared)
-- ===============================

CREATE SCHEMA shared;

-- Tabla principal de empresas registradas
CREATE TABLE shared.empresas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ruc VARCHAR(13) UNIQUE NOT NULL,
    razon_social VARCHAR(255) NOT NULL,
    nombre_comercial VARCHAR(255),
    email_principal VARCHAR(100) NOT NULL,
    telefono VARCHAR(20),
    direccion TEXT,
    
    -- Información de suscripción
    plan_activo VARCHAR(20) NOT NULL DEFAULT 'BASICO',
    estado VARCHAR(20) NOT NULL DEFAULT 'ACTIVA', -- ACTIVA, SUSPENDIDA, CANCELADA
    fecha_registro TIMESTAMP NOT NULL DEFAULT NOW(),
    fecha_ultimo_pago TIMESTAMP,
    fecha_vencimiento TIMESTAMP,
    
    -- Configuración de certificado
    certificado_ruta VARCHAR(500),
    certificado_hash VARCHAR(64),
    certificado_vencimiento DATE,
    certificado_valido BOOLEAN DEFAULT false,
    
    -- Límites por plan
    limite_facturas_mes INTEGER NOT NULL DEFAULT 100,
    limite_establecimientos INTEGER NOT NULL DEFAULT 1,
    limite_usuarios INTEGER NOT NULL DEFAULT 2,
    
    -- Configuración general
    configuracion JSONB DEFAULT '{}',
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by UUID,
    updated_by UUID
);

-- Usuarios del sistema (pueden gestionar múltiples empresas)
CREATE TABLE shared.usuarios_sistema (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    telefono VARCHAR(20),
    
    -- Estado del usuario
    activo BOOLEAN DEFAULT true,
    email_verificado BOOLEAN DEFAULT false,
    ultimo_acceso TIMESTAMP,
    intentos_login_fallidos INTEGER DEFAULT 0,
    bloqueado_hasta TIMESTAMP,
    
    -- Configuración
    timezone VARCHAR(50) DEFAULT 'America/Guayaquil',
    idioma VARCHAR(10) DEFAULT 'es',
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Relación usuario-empresa (un usuario puede manejar varias empresas)
CREATE TABLE shared.usuarios_empresas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id UUID NOT NULL REFERENCES shared.usuarios_sistema(id),
    empresa_id UUID NOT NULL REFERENCES shared.empresas(id),
    
    -- Rol específico en esta empresa
    rol VARCHAR(20) NOT NULL, -- OWNER, ADMIN, USER, READONLY
    activo BOOLEAN DEFAULT true,
    
    -- Permisos específicos
    permisos JSONB DEFAULT '{}',
    
    -- Auditoría  
    created_at TIMESTAMP DEFAULT NOW(),
    created_by UUID,
    
    UNIQUE(usuario_id, empresa_id)
);

-- Planes de suscripción disponibles
CREATE TABLE shared.planes_suscripcion (
    id VARCHAR(20) PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    descripcion TEXT,
    precio_mensual DECIMAL(10,2) NOT NULL,
    
    -- Límites del plan
    limite_facturas_mes INTEGER NOT NULL,
    limite_establecimientos INTEGER NOT NULL,
    limite_usuarios INTEGER NOT NULL,
    api_incluida BOOLEAN DEFAULT false,
    soporte_nivel VARCHAR(20) DEFAULT 'EMAIL', -- EMAIL, TELEFONO, PRIORITARIO
    
    -- Características
    caracteristicas JSONB DEFAULT '[]',
    
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Auditoría de acciones del sistema
CREATE TABLE shared.audit_logs (
    id BIGSERIAL PRIMARY KEY,
    empresa_id UUID REFERENCES shared.empresas(id),
    usuario_id UUID REFERENCES shared.usuarios_sistema(id),
    
    accion VARCHAR(100) NOT NULL,
    entidad VARCHAR(50), -- facturas, clientes, etc.
    entidad_id VARCHAR(100),
    
    datos_anteriores JSONB,
    datos_nuevos JSONB,
    
    ip_address INET,
    user_agent TEXT,
    
    timestamp TIMESTAMP DEFAULT NOW()
);

-- ===============================
-- FUNCIÓN PARA CREAR SCHEMA EMPRESA
-- ===============================

CREATE OR REPLACE FUNCTION crear_schema_empresa(p_ruc VARCHAR(13))
RETURNS VOID AS $$
DECLARE
    schema_name VARCHAR(50);
BEGIN
    schema_name := 'empresa_' || p_ruc;
    
    -- Crear schema
    EXECUTE format('CREATE SCHEMA %I', schema_name);
    
    -- Crear tablas específicas de la empresa
    EXECUTE format('
        CREATE TABLE %I.facturas (
            id BIGSERIAL PRIMARY KEY,
            numero_factura VARCHAR(50) NOT NULL UNIQUE,
            clave_acceso VARCHAR(49) NOT NULL UNIQUE,
            secuencial_interno INTEGER NOT NULL,
            
            -- Información temporal
            fecha_emision TIMESTAMP NOT NULL,
            fecha_vencimiento DATE,
            
            -- Cliente
            cliente_tipo_identificacion VARCHAR(10) NOT NULL,
            cliente_identificacion VARCHAR(20) NOT NULL,
            cliente_razon_social VARCHAR(255) NOT NULL,
            cliente_direccion TEXT,
            cliente_telefono VARCHAR(20),
            cliente_email VARCHAR(100),
            
            -- Establecimiento
            establecimiento_codigo VARCHAR(3) NOT NULL DEFAULT ''001'',
            punto_emision_codigo VARCHAR(3) NOT NULL DEFAULT ''001'',
            
            -- Totales
            subtotal_sin_impuestos DECIMAL(12,2) NOT NULL,
            subtotal_0 DECIMAL(12,2) DEFAULT 0,
            subtotal_12 DECIMAL(12,2) DEFAULT 0,
            subtotal_15 DECIMAL(12,2) DEFAULT 0,
            total_descuento DECIMAL(12,2) DEFAULT 0,
            ice DECIMAL(12,2) DEFAULT 0,
            iva_12 DECIMAL(12,2) DEFAULT 0,
            iva_15 DECIMAL(12,2) DEFAULT 0,
            propina DECIMAL(12,2) DEFAULT 0,
            importe_total DECIMAL(12,2) NOT NULL,
            
            -- Detalles como JSON
            productos JSONB NOT NULL,
            informacion_adicional JSONB DEFAULT ''[]'',
            
            -- Estado SRI
            estado_sri VARCHAR(20) NOT NULL DEFAULT ''BORRADOR'',
            numero_autorizacion VARCHAR(50),
            fecha_autorizacion TIMESTAMP,
            observaciones_sri TEXT,
            
            -- XMLs
            xml_generado TEXT,
            xml_autorizado TEXT,
            xml_firmado TEXT,
            
            -- Configuración
            ambiente VARCHAR(20) NOT NULL DEFAULT ''PRUEBAS'',
            tipo_emision VARCHAR(10) NOT NULL DEFAULT ''NORMAL'',
            
            -- Auditoría
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            created_by UUID,
            updated_by UUID,
            
            -- Índices
            CONSTRAINT fk_facturas_created_by FOREIGN KEY (created_by) REFERENCES shared.usuarios_sistema(id),
            CONSTRAINT fk_facturas_updated_by FOREIGN KEY (updated_by) REFERENCES shared.usuarios_sistema(id)
        )', schema_name);
    
    -- Tabla de clientes
    EXECUTE format('
        CREATE TABLE %I.clientes (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            
            -- Identificación
            tipo_identificacion VARCHAR(10) NOT NULL,
            identificacion VARCHAR(20) NOT NULL UNIQUE,
            razon_social VARCHAR(255) NOT NULL,
            nombre_comercial VARCHAR(255),
            
            -- Contacto
            direccion TEXT,
            telefono VARCHAR(20),
            email VARCHAR(100),
            
            -- Clasificación
            tipo_cliente VARCHAR(20) DEFAULT ''PERSONA_NATURAL'',
            categoria VARCHAR(50),
            
            -- Estado
            activo BOOLEAN DEFAULT true,
            
            -- Información adicional
            observaciones TEXT,
            datos_adicionales JSONB DEFAULT ''{}}'',
            
            -- Auditoría
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            created_by UUID,
            updated_by UUID,
            
            CONSTRAINT fk_clientes_created_by FOREIGN KEY (created_by) REFERENCES shared.usuarios_sistema(id),
            CONSTRAINT fk_clientes_updated_by FOREIGN KEY (updated_by) REFERENCES shared.usuarios_sistema(id)
        )', schema_name);
    
    -- Tabla de productos/servicios
    EXECUTE format('
        CREATE TABLE %I.productos (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            codigo VARCHAR(50) NOT NULL UNIQUE,
            descripcion VARCHAR(500) NOT NULL,
            
            -- Clasificación
            tipo VARCHAR(20) DEFAULT ''PRODUCTO'', -- PRODUCTO, SERVICIO
            categoria VARCHAR(100),
            unidad_medida VARCHAR(20) DEFAULT ''UNIDAD'',
            
            -- Precios
            precio_unitario DECIMAL(12,4) NOT NULL,
            precio_minimo DECIMAL(12,4),
            precio_maximo DECIMAL(12,4),
            
            -- Impuestos
            gravado_iva BOOLEAN DEFAULT true,
            porcentaje_iva INTEGER DEFAULT 15,
            codigo_ice VARCHAR(10),
            porcentaje_ice DECIMAL(5,2) DEFAULT 0,
            
            -- Stock (opcional)
            maneja_stock BOOLEAN DEFAULT false,
            stock_actual DECIMAL(12,4) DEFAULT 0,
            stock_minimo DECIMAL(12,4) DEFAULT 0,
            
            -- Estado
            activo BOOLEAN DEFAULT true,
            
            -- Información adicional
            observaciones TEXT,
            datos_adicionales JSONB DEFAULT ''{}}'',
            
            -- Auditoría
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            created_by UUID,
            updated_by UUID,
            
            CONSTRAINT fk_productos_created_by FOREIGN KEY (created_by) REFERENCES shared.usuarios_sistema(id),
            CONSTRAINT fk_productos_updated_by FOREIGN KEY (updated_by) REFERENCES shared.usuarios_sistema(id)
        )', schema_name);
    
    -- Tabla de establecimientos
    EXECUTE format('
        CREATE TABLE %I.establecimientos (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            codigo VARCHAR(3) NOT NULL UNIQUE,
            nombre VARCHAR(100) NOT NULL,
            
            -- Ubicación
            direccion TEXT NOT NULL,
            telefono VARCHAR(20),
            email VARCHAR(100),
            
            -- Configuración
            por_defecto BOOLEAN DEFAULT false,
            punto_emision_defecto VARCHAR(3) DEFAULT ''001'',
            
            -- Estado
            activo BOOLEAN DEFAULT true,
            
            -- Auditoría
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            created_by UUID,
            
            CONSTRAINT fk_establecimientos_created_by FOREIGN KEY (created_by) REFERENCES shared.usuarios_sistema(id)
        )', schema_name);
    
    -- Tabla de configuración específica de la empresa
    EXECUTE format('
        CREATE TABLE %I.configuracion (
            clave VARCHAR(100) PRIMARY KEY,
            valor JSONB NOT NULL,
            descripcion TEXT,
            
            updated_at TIMESTAMP DEFAULT NOW(),
            updated_by UUID,
            
            CONSTRAINT fk_configuracion_updated_by FOREIGN KEY (updated_by) REFERENCES shared.usuarios_sistema(id)
        )', schema_name);
    
    -- Crear índices importantes
    EXECUTE format('CREATE INDEX idx_%I_facturas_fecha_emision ON %I.facturas(fecha_emision DESC)', 
                   replace(schema_name, 'empresa_', ''), schema_name);
    EXECUTE format('CREATE INDEX idx_%I_facturas_estado_sri ON %I.facturas(estado_sri)', 
                   replace(schema_name, 'empresa_', ''), schema_name);
    EXECUTE format('CREATE INDEX idx_%I_facturas_cliente ON %I.facturas(cliente_identificacion)', 
                   replace(schema_name, 'empresa_', ''), schema_name);
    EXECUTE format('CREATE INDEX idx_%I_clientes_identificacion ON %I.clientes(identificacion)', 
                   replace(schema_name, 'empresa_', ''), schema_name);
    EXECUTE format('CREATE INDEX idx_%I_productos_codigo ON %I.productos(codigo)', 
                   replace(schema_name, 'empresa_', ''), schema_name);

    -- Insertar configuración por defecto
    EXECUTE format('
        INSERT INTO %I.configuracion (clave, valor, descripcion) VALUES
        (''secuencial_actual'', ''{"factura": 1}'', ''Números secuenciales por tipo de documento''),
        (''ambientes'', ''{"actual": "PRUEBAS"}'', ''Ambiente actual de trabajo''),
        (''numeracion'', ''{"establecimiento": "001", "punto_emision": "001"}'', ''Configuración de numeración''),
        (''sri_endpoints'', ''{"recepcion": "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline", "autorizacion": "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline"}'', ''Endpoints del SRI''),
        (''empresa_datos'', ''{"logo_url": "", "eslogan": "", "website": ""}'', ''Datos adicionales de la empresa'')
    ', schema_name);
    
    RAISE NOTICE 'Schema creado exitosamente: %', schema_name;
END;
$$ LANGUAGE plpgsql;
```

---

## 🔐 Sistema de Autenticación y Autorización

### Flujo de Autenticación

```go
// JWT Token con información de empresa
type JWTClaims struct {
    UserID      string   `json:"user_id"`
    Email       string   `json:"email"`
    EmpresaID   string   `json:"empresa_id"`
    EmpresaRUC  string   `json:"empresa_ruc"`
    Rol         string   `json:"rol"`
    Permisos    []string `json:"permisos"`
    Plan        string   `json:"plan"`
    jwt.StandardClaims
}

// Middleware de autenticación multi-tenant
func MiddlewareAutenticacion(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extraer token JWT
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Token de autorización requerido", 401)
            return
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        // 2. Validar y parsear token
        claims := &JWTClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        
        if err != nil || !token.Valid {
            http.Error(w, "Token inválido", 401)
            return
        }
        
        // 3. Verificar que la empresa está activa
        empresa, err := db.ObtenerEmpresa(claims.EmpresaID)
        if err != nil || empresa.Estado != "ACTIVA" {
            http.Error(w, "Empresa no autorizada", 403)
            return
        }
        
        // 4. Verificar límites del plan
        if !verificarLimitesPlan(empresa, r.URL.Path) {
            http.Error(w, "Límite del plan excedido", 429)
            return
        }
        
        // 5. Agregar información al contexto
        ctx := context.WithValue(r.Context(), "usuario", claims)
        ctx = context.WithValue(ctx, "empresa", empresa)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware específico para aislamiento de datos
func MiddlewareAislamiento(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value("usuario").(*JWTClaims)
        
        // Configurar conexión de base de datos para el schema específico
        schemaName := fmt.Sprintf("empresa_%s", claims.EmpresaRUC)
        
        // Crear conexión con search_path específico
        db := database.GetConnection()
        _, err := db.Exec(fmt.Sprintf("SET search_path TO %s, shared", schemaName))
        if err != nil {
            http.Error(w, "Error configurando acceso a datos", 500)
            return
        }
        
        // Agregar conexión configurada al contexto
        ctx := context.WithValue(r.Context(), "db_schema", schemaName)
        ctx = context.WithValue(ctx, "db_connection", db)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Verificación de Límites por Plan

```go
func verificarLimitesPlan(empresa *Empresa, endpoint string) bool {
    plan := obtenerPlan(empresa.PlanActivo)
    fechaInicio := time.Now().AddDate(0, 0, -30) // Último mes
    
    switch {
    case strings.Contains(endpoint, "/facturas"):
        // Verificar límite de facturas
        count, err := contarFacturasMes(empresa.ID, fechaInicio)
        if err != nil {
            return false
        }
        return count < plan.LimiteFacturasMes
        
    case strings.Contains(endpoint, "/establecimientos"):
        // Verificar límite de establecimientos
        count, err := contarEstablecimientos(empresa.ID)
        if err != nil {
            return false
        }
        return count < plan.LimiteEstablecimientos
        
    case strings.Contains(endpoint, "/api/") && plan.APIIncluida == false:
        // Verificar si el plan incluye API
        return false
        
    default:
        return true
    }
}
```

---

## 📊 APIs Multi-Tenant

### Controladores con Aislamiento

```go
// Controlador de facturas con aislamiento automático
type FacturaController struct {
    db *database.Connection
}

func (fc *FacturaController) CrearFactura(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener información del contexto
    claims := r.Context().Value("usuario").(*JWTClaims)
    empresa := r.Context().Value("empresa").(*Empresa)
    dbConn := r.Context().Value("db_connection").(*sql.DB)
    
    // 2. Parsear datos de entrada
    var facturaData models.FacturaInput
    if err := json.NewDecoder(r.Body).Decode(&facturaData); err != nil {
        http.Error(w, "Datos inválidos", 400)
        return
    }
    
    // 3. Validar permisos específicos
    if !tienePermiso(claims.Permisos, "crear_facturas") {
        http.Error(w, "Sin permisos para crear facturas", 403)
        return
    }
    
    // 4. Crear factura usando configuración de la empresa
    factura, err := fc.crearFacturaEmpresa(facturaData, empresa, claims.UserID, dbConn)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creando factura: %v", err), 500)
        return
    }
    
    // 5. Registrar en auditoría
    fc.registrarAuditoria("CREAR_FACTURA", claims.UserID, empresa.ID, factura.ID)
    
    // 6. Responder
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "factura": factura,
        "empresa": empresa.RazonSocial,
    })
}

func (fc *FacturaController) crearFacturaEmpresa(data models.FacturaInput, empresa *Empresa, userID string, db *sql.DB) (*models.Factura, error) {
    // 1. Obtener secuencial único para esta empresa
    secuencial, err := fc.obtenerProximoSecuencial(empresa.RUC, db)
    if err != nil {
        return nil, err
    }
    
    // 2. Generar número de factura
    numeroFactura := fmt.Sprintf("%s-%s-%09d", 
        empresa.EstablecimientoPorDefecto, 
        empresa.PuntoEmisionPorDefecto, 
        secuencial)
    
    // 3. Crear factura usando factory con configuración de empresa
    config := factory.ConfigEmpresa{
        RUC:              empresa.RUC,
        RazonSocial:      empresa.RazonSocial,
        Establecimiento:  empresa.EstablecimientoPorDefecto,
        PuntoEmision:     empresa.PuntoEmisionPorDefecto,
        Ambiente:         empresa.AmbienteSRI,
        CertificadoRuta:  empresa.CertificadoRuta,
    }
    
    factura, err := factory.CrearFacturaConConfig(data, config)
    if err != nil {
        return nil, err
    }
    
    // 4. Asignar identificadores únicos de la empresa
    factura.NumeroFactura = numeroFactura
    factura.SecuencialInterno = secuencial
    
    // 5. Guardar en schema específico de la empresa
    err = fc.guardarFacturaEnSchema(factura, empresa.RUC, userID, db)
    if err != nil {
        return nil, err
    }
    
    return factura, nil
}

func (fc *FacturaController) guardarFacturaEnSchema(factura *models.Factura, ruc string, userID string, db *sql.DB) error {
    // La conexión ya tiene el search_path configurado al schema correcto
    query := `
        INSERT INTO facturas (
            numero_factura, clave_acceso, secuencial_interno,
            fecha_emision, cliente_identificacion, cliente_razon_social,
            subtotal_sin_impuestos, importe_total, productos,
            estado_sri, ambiente, created_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `
    
    productosJSON, _ := json.Marshal(factura.Detalles)
    
    err := db.QueryRow(query,
        factura.NumeroFactura,
        factura.InfoTributaria.ClaveAcceso,
        factura.SecuencialInterno,
        factura.InfoFactura.FechaEmision,
        factura.InfoFactura.IdentificacionComprador,
        factura.InfoFactura.RazonSocialComprador,
        factura.InfoFactura.TotalSinImpuestos,
        factura.InfoFactura.ImporteTotal,
        productosJSON,
        "BORRADOR",
        factura.InfoTributaria.Ambiente,
        userID,
    ).Scan(&factura.ID)
    
    return err
}

// Listar facturas con filtros automáticos por empresa
func (fc *FacturaController) ListarFacturas(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("usuario").(*JWTClaims)
    dbConn := r.Context().Value("db_connection").(*sql.DB)
    
    // Parsear parámetros de consulta
    filtros := ParsearFiltros(r.URL.Query())
    
    // La consulta automáticamente se ejecuta en el schema correcto
    query := `
        SELECT id, numero_factura, fecha_emision, cliente_razon_social, 
               importe_total, estado_sri, created_at
        FROM facturas 
        WHERE fecha_emision BETWEEN $1 AND $2
        ORDER BY fecha_emision DESC
        LIMIT $3 OFFSET $4
    `
    
    rows, err := dbConn.Query(query, 
        filtros.FechaInicio, 
        filtros.FechaFin, 
        filtros.Limit, 
        filtros.Offset)
    if err != nil {
        http.Error(w, "Error consultando facturas", 500)
        return
    }
    defer rows.Close()
    
    var facturas []models.FacturaResumen
    for rows.Next() {
        var f models.FacturaResumen
        err := rows.Scan(&f.ID, &f.NumeroFactura, &f.FechaEmision, 
                        &f.ClienteNombre, &f.Total, &f.EstadoSRI, &f.FechaCreacion)
        if err != nil {
            continue
        }
        facturas = append(facturas, f)
    }
    
    // Contar total para paginación
    var total int
    countQuery := `SELECT COUNT(*) FROM facturas WHERE fecha_emision BETWEEN $1 AND $2`
    dbConn.QueryRow(countQuery, filtros.FechaInicio, filtros.FechaFin).Scan(&total)
    
    response := map[string]interface{}{
        "facturas": facturas,
        "total":    total,
        "empresa":  claims.EmpresaRUC,
        "filtros":  filtros,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

## 🔒 Gestión de Certificados Multi-Tenant

### Almacenamiento Seguro por Empresa

```go
type GestorCertificadosMultiTenant struct {
    basePath   string
    masterKey  []byte
    db         *database.Connection
}

func (g *GestorCertificadosMultiTenant) GuardarCertificado(empresaID, ruc string, certificadoP12 []byte, password string) error {
    // 1. Validar que el certificado corresponde a la empresa
    cert, err := validarCertificadoP12(certificadoP12, password)
    if err != nil {
        return fmt.Errorf("certificado inválido: %v", err)
    }
    
    // Verificar que el RUC del certificado coincide con la empresa
    if cert.RUCEmisor != ruc {
        return fmt.Errorf("el certificado no corresponde al RUC de la empresa")
    }
    
    // 2. Crear directorio específico para la empresa
    dirEmpresa := filepath.Join(g.basePath, "certificados", ruc)
    err = os.MkdirAll(dirEmpresa, 0700)
    if err != nil {
        return fmt.Errorf("error creando directorio: %v", err)
    }
    
    // 3. Encriptar certificado con clave específica de empresa
    claveEmpresa := g.derivarClaveEmpresa(empresaID, ruc)
    certificadoEncriptado, err := encriptarAES256(certificadoP12, claveEmpresa)
    if err != nil {
        return fmt.Errorf("error encriptando certificado: %v", err)
    }
    
    // 4. Guardar archivo encriptado
    nombreArchivo := fmt.Sprintf("cert_%s_%d.p12.enc", ruc, time.Now().Unix())
    rutaCompleta := filepath.Join(dirEmpresa, nombreArchivo)
    
    err = ioutil.WriteFile(rutaCompleta, certificadoEncriptado, 0600)
    if err != nil {
        return fmt.Errorf("error guardando archivo: %v", err)
    }
    
    // 5. Calcular hash para verificación de integridad
    hash := sha256.Sum256(certificadoP12)
    hashHex := hex.EncodeToString(hash[:])
    
    // 6. Actualizar base de datos
    query := `
        UPDATE shared.empresas 
        SET certificado_ruta = $1,
            certificado_hash = $2,
            certificado_vencimiento = $3,
            certificado_valido = true,
            updated_at = NOW()
        WHERE id = $4
    `
    
    _, err = g.db.Exec(query, rutaCompleta, hashHex, cert.FechaVencimiento, empresaID)
    if err != nil {
        // Si falla la DB, eliminar archivo
        os.Remove(rutaCompleta)
        return fmt.Errorf("error actualizando base de datos: %v", err)
    }
    
    // 7. Eliminar certificado anterior si existe
    g.limpiarCertificadosAnteriores(dirEmpresa, nombreArchivo)
    
    return nil
}

func (g *GestorCertificadosMultiTenant) CargarCertificado(empresaID string) (*CertificadoDigital, error) {
    // 1. Obtener información del certificado de la DB
    var rutaCert, hashEsperado, ruc string
    var fechaVencimiento time.Time
    
    query := `
        SELECT certificado_ruta, certificado_hash, certificado_vencimiento, ruc
        FROM shared.empresas 
        WHERE id = $1 AND certificado_valido = true
    `
    
    err := g.db.QueryRow(query, empresaID).Scan(&rutaCert, &hashEsperado, &fechaVencimiento, &ruc)
    if err != nil {
        return nil, fmt.Errorf("certificado no encontrado para empresa: %v", err)
    }
    
    // 2. Verificar que el certificado no ha expirado
    if time.Now().After(fechaVencimiento) {
        return nil, fmt.Errorf("certificado ha expirado el %v", fechaVencimiento)
    }
    
    // 3. Leer archivo encriptado
    certificadoEncriptado, err := ioutil.ReadFile(rutaCert)
    if err != nil {
        return nil, fmt.Errorf("error leyendo certificado: %v", err)
    }
    
    // 4. Desencriptar con clave específica de empresa
    claveEmpresa := g.derivarClaveEmpresa(empresaID, ruc)
    certificadoP12, err := desencriptarAES256(certificadoEncriptado, claveEmpresa)
    if err != nil {
        return nil, fmt.Errorf("error desencriptando certificado: %v", err)
    }
    
    // 5. Verificar integridad
    hash := sha256.Sum256(certificadoP12)
    hashActual := hex.EncodeToString(hash[:])
    
    if hashActual != hashEsperado {
        return nil, fmt.Errorf("certificado corrupto - hash no coincide")
    }
    
    // 6. Cargar y validar certificado
    cert, err := cargarCertificadoP12(certificadoP12, g.obtenerPasswordCertificado(empresaID))
    if err != nil {
        return nil, fmt.Errorf("error cargando certificado: %v", err)
    }
    
    return cert, nil
}

func (g *GestorCertificadosMultiTenant) derivarClaveEmpresa(empresaID, ruc string) []byte {
    // Derivar clave única por empresa usando PBKDF2
    salt := []byte(fmt.Sprintf("facturacion_sri_%s_%s", empresaID, ruc))
    return pbkdf2.Key(g.masterKey, salt, 10000, 32, sha256.New)
}

// Validación periódica de certificados
func (g *GestorCertificadosMultiTenant) ValidarCertificadosPeriodicamente() {
    ticker := time.NewTicker(24 * time.Hour) // Revisar diariamente
    
    for range ticker.C {
        query := `
            SELECT id, ruc, razon_social, certificado_vencimiento
            FROM shared.empresas 
            WHERE certificado_valido = true 
            AND estado = 'ACTIVA'
        `
        
        rows, err := g.db.Query(query)
        if err != nil {
            continue
        }
        
        for rows.Next() {
            var empresaID, ruc, razonSocial string
            var fechaVencimiento time.Time
            
            rows.Scan(&empresaID, &ruc, &razonSocial, &fechaVencimiento)
            
            // Notificar si el certificado vence en menos de 30 días
            diasRestantes := int(time.Until(fechaVencimiento).Hours() / 24)
            if diasRestantes <= 30 && diasRestantes > 0 {
                g.enviarNotificacionVencimiento(empresaID, razonSocial, diasRestantes)
            } else if diasRestantes <= 0 {
                // Marcar certificado como expirado
                g.marcarCertificadoExpirado(empresaID)
            }
        }
        rows.Close()
    }
}
```

---

## 📈 Monitoreo y Métricas por Tenant

### Dashboard de Administración

```go
type MetricasMultiTenant struct {
    db *database.Connection
}

func (m *MetricasMultiTenant) ObtenerResumenGeneral() (*ResumenGeneral, error) {
    resumen := &ResumenGeneral{}
    
    // Estadísticas generales
    query := `
        SELECT 
            COUNT(*) as total_empresas,
            COUNT(CASE WHEN estado = 'ACTIVA' THEN 1 END) as empresas_activas,
            COUNT(CASE WHEN certificado_valido = true THEN 1 END) as empresas_con_certificado,
            SUM(CASE WHEN plan_activo = 'BASICO' THEN 1 ELSE 0 END) as plan_basico,
            SUM(CASE WHEN plan_activo = 'PROFESIONAL' THEN 1 ELSE 0 END) as plan_profesional,
            SUM(CASE WHEN plan_activo = 'EMPRESARIAL' THEN 1 ELSE 0 END) as plan_empresarial
        FROM shared.empresas
    `
    
    err := m.db.QueryRow(query).Scan(
        &resumen.TotalEmpresas,
        &resumen.EmpresasActivas,
        &resumen.EmpresasConCertificado,
        &resumen.PlanBasico,
        &resumen.PlanProfesional,
        &resumen.PlanEmpresarial,
    )
    
    if err != nil {
        return nil, err
    }
    
    // Calcular ingresos estimados
    ingresos := (resumen.PlanBasico * 29) + 
                (resumen.PlanProfesional * 59) + 
                (resumen.PlanEmpresarial * 119)
    resumen.IngresosEstimados = ingresos
    
    // Obtener estadísticas de uso
    resumen.FacturasUltimoMes = m.contarFacturasUltimoMes()
    resumen.EmpresasMasActivas = m.obtenerEmpresasMasActivas(10)
    
    return resumen, nil
}

func (m *MetricasMultiTenant) contarFacturasUltimoMes() int {
    var total int
    fechaInicio := time.Now().AddDate(0, -1, 0)
    
    // Consultar en todos los schemas de empresas
    query := `
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name LIKE 'empresa_%'
    `
    
    rows, err := m.db.Query(query)
    if err != nil {
        return 0
    }
    defer rows.Close()
    
    for rows.Next() {
        var schemaName string
        rows.Scan(&schemaName)
        
        // Contar facturas en este schema
        countQuery := fmt.Sprintf(`
            SELECT COUNT(*) 
            FROM %s.facturas 
            WHERE fecha_emision >= $1
        `, schemaName)
        
        var count int
        m.db.QueryRow(countQuery, fechaInicio).Scan(&count)
        total += count
    }
    
    return total
}

func (m *MetricasMultiTenant) ObtenerMetricasEmpresa(empresaID string) (*MetricasEmpresa, error) {
    // Obtener información básica de la empresa
    var ruc, razonSocial, plan string
    query := `
        SELECT ruc, razon_social, plan_activo
        FROM shared.empresas 
        WHERE id = $1
    `
    
    err := m.db.QueryRow(query, empresaID).Scan(&ruc, &razonSocial, &plan)
    if err != nil {
        return nil, err
    }
    
    metricas := &MetricasEmpresa{
        EmpresaID:    empresaID,
        RUC:          ruc,
        RazonSocial:  razonSocial,
        Plan:         plan,
    }
    
    // Configurar schema específico
    schemaName := fmt.Sprintf("empresa_%s", ruc)
    
    // Métricas del último mes
    fechaInicio := time.Now().AddDate(0, -1, 0)
    
    // Facturas emitidas
    countQuery := fmt.Sprintf(`
        SELECT 
            COUNT(*) as total,
            COUNT(CASE WHEN estado_sri = 'AUTORIZADA' THEN 1 END) as autorizadas,
            COALESCE(SUM(importe_total), 0) as total_facturado
        FROM %s.facturas 
        WHERE fecha_emision >= $1
    `, schemaName)
    
    err = m.db.QueryRow(countQuery, fechaInicio).Scan(
        &metricas.FacturasEmitidas,
        &metricas.FacturasAutorizadas,
        &metricas.TotalFacturado,
    )
    
    if err != nil {
        return nil, err
    }
    
    // Clientes únicos
    clientesQuery := fmt.Sprintf(`
        SELECT COUNT(DISTINCT cliente_identificacion)
        FROM %s.facturas 
        WHERE fecha_emision >= $1
    `, schemaName)
    
    m.db.QueryRow(clientesQuery, fechaInicio).Scan(&metricas.ClientesUnicos)
    
    // Productos más vendidos
    metricas.ProductosMasVendidos = m.obtenerProductosMasVendidos(schemaName, fechaInicio)
    
    return metricas, nil
}

// Alertas automáticas del sistema
func (m *MetricasMultiTenant) VerificarAlertas() {
    // 1. Empresas cerca del límite de su plan
    query := `
        SELECT id, ruc, razon_social, plan_activo, limite_facturas_mes
        FROM shared.empresas 
        WHERE estado = 'ACTIVA'
    `
    
    rows, err := m.db.Query(query)
    if err != nil {
        return
    }
    defer rows.Close()
    
    fechaInicio := time.Now().AddDate(0, 0, -30)
    
    for rows.Next() {
        var empresaID, ruc, razonSocial, plan string
        var limite int
        
        rows.Scan(&empresaID, &ruc, &razonSocial, &plan, &limite)
        
        // Contar facturas del último mes
        countQuery := fmt.Sprintf(`
            SELECT COUNT(*) 
            FROM empresa_%s.facturas 
            WHERE fecha_emision >= $1
        `, ruc)
        
        var count int
        m.db.QueryRow(countQuery, fechaInicio).Scan(&count)
        
        // Alertar si está al 80% del límite
        if float64(count) >= float64(limite)*0.8 {
            m.enviarAlertaLimite(empresaID, razonSocial, count, limite)
        }
    }
    
    // 2. Certificados por vencer
    certQuery := `
        SELECT id, ruc, razon_social, certificado_vencimiento
        FROM shared.empresas 
        WHERE certificado_valido = true 
        AND certificado_vencimiento <= $1
    `
    
    fechaAlerta := time.Now().AddDate(0, 0, 30) // 30 días
    
    certRows, err := m.db.Query(certQuery, fechaAlerta)
    if err != nil {
        return
    }
    defer certRows.Close()
    
    for certRows.Next() {
        var empresaID, ruc, razonSocial string
        var fechaVencimiento time.Time
        
        certRows.Scan(&empresaID, &ruc, &razonSocial, &fechaVencimiento)
        
        diasRestantes := int(time.Until(fechaVencimiento).Hours() / 24)
        m.enviarAlertaCertificado(empresaID, razonSocial, diasRestantes)
    }
}
```

---

## 🚀 Despliegue y Escalabilidad

### Configuración Docker Multi-Tenant

```dockerfile
# Dockerfile para aplicación multi-tenant
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o facturacion-saas main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Crear estructura de directorios
RUN mkdir -p /data/certificados
RUN mkdir -p /data/backups
RUN mkdir -p /logs

COPY --from=builder /app/facturacion-saas .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./facturacion-saas"]
```

```yaml
# docker-compose.yml para ambiente completo
version: '3.8'

services:
  # Base de datos principal
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: facturacion_saas
      POSTGRES_USER: facturacion_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U facturacion_user"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Aplicación principal
  facturacion-app:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: facturacion_saas
      DB_USER: facturacion_user
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      ENCRYPTION_KEY: ${ENCRYPTION_KEY}
      SRI_TIMEOUT: 30
      LOG_LEVEL: info
    volumes:
      - cert_storage:/data/certificados
      - backup_storage:/data/backups
      - app_logs:/logs
    ports:
      - "8080:8080"
    restart: unless-stopped

  # Redis para caché y sesiones
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  # Nginx como proxy reverso
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - facturacion-app
    restart: unless-stopped

  # Backup automático
  backup:
    image: postgres:15-alpine
    depends_on:
      - postgres
    environment:
      PGPASSWORD: ${DB_PASSWORD}
    volumes:
      - backup_storage:/backups
      - ./backup-script.sh:/backup-script.sh
    command: |
      sh -c '
        while true; do
          echo "Iniciando backup..."
          pg_dump -h postgres -U facturacion_user -d facturacion_saas > /backups/backup_$$(date +%Y%m%d_%H%M%S).sql
          echo "Backup completado"
          # Limpiar backups antiguos (mantener últimos 7 días)
          find /backups -name "backup_*.sql" -mtime +7 -delete
          sleep 86400  # 24 horas
        done
      '

volumes:
  postgres_data:
  redis_data:
  cert_storage:
  backup_storage:
  app_logs:
```

### Configuración Nginx

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app_servers {
        server facturacion-app:8080;
        # Para escalamiento horizontal:
        # server facturacion-app-2:8080;
        # server facturacion-app-3:8080;
    }

    # Rate limiting por empresa
    map $http_x_empresa_id $rate_limit_key {
        default $http_x_empresa_id;
        "" $binary_remote_addr;
    }

    limit_req_zone $rate_limit_key zone=api_limit:10m rate=10r/s;
    limit_req_zone $rate_limit_key zone=auth_limit:10m rate=1r/s;

    server {
        listen 80;
        server_name your-domain.com;
        
        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        ssl_certificate /etc/ssl/certs/your-domain.crt;
        ssl_certificate_key /etc/ssl/certs/your-domain.key;

        # Security headers
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

        # API endpoints
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;
            
            proxy_pass http://app_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # Timeouts para SRI
            proxy_connect_timeout 30s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Authentication endpoints
        location /auth/ {
            limit_req zone=auth_limit burst=5 nodelay;
            
            proxy_pass http://app_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Static files
        location /static/ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }

        # Health check
        location /health {
            proxy_pass http://app_servers;
            access_log off;
        }
    }
}
```

---

## 📊 Conclusión

Esta arquitectura multi-tenant proporciona:

### ✅ **Beneficios Técnicos**
- **Aislamiento completo** de datos por empresa
- **Escalabilidad horizontal** mediante schemas separados
- **Seguridad robusta** con encriptación por empresa
- **Monitoreo detallado** por tenant y global
- **Backup y recuperación** granular

### ✅ **Beneficios de Negocio**
- **Costos optimizados** - una infraestructura para todos
- **Mantenimiento centralizado** - actualizaciones simultáneas
- **Onboarding rápido** - nuevas empresas en minutos
- **Escalamiento automático** - crece con la demanda
- **Métricas de negocio** - visibilidad completa del SaaS

### 🎯 **Preparado para Producción**
- Configuración Docker completa
- Proxy reverso con SSL
- Backup automático
- Rate limiting por empresa
- Logs centralizados
- Monitoreo de salud

**¡Tu sistema está listo para servir a cientos de empresas ecuatorianas! 🇪🇨**