# Mejores Prácticas en Go - Código Profesional y Mantenible

Esta guía te enseña las mejores prácticas aplicadas en el proyecto de facturación SRI, con ejemplos antes/después.

## 📚 Tabla de Contenido

1. [Principios Fundamentales](#principios-fundamentales)
2. [Estructura de Proyecto](#estructura-de-proyecto)
3. [Manejo de Errores](#manejo-de-errores)
4. [Testing Estratégico](#testing-estratégico)
5. [Concurrencia Segura](#concurrencia-segura)
6. [Interfaces Efectivas](#interfaces-efectivas)
7. [Rendimiento y Memoria](#rendimiento-y-memoria)
8. [Seguridad](#seguridad)
9. [Documentación](#documentación)
10. [Herramientas de Desarrollo](#herramientas-de-desarrollo)

---

## Principios Fundamentales

### 1. Simplicidad sobre Cleverness

```go
// ❌ MALO: Código "inteligente" pero difícil de leer
func calcTotal(items []Item) float64 {
    return funk.Reduce(items, func(acc, item interface{}) interface{} {
        return acc.(float64) + item.(Item).Price*item.(Item).Qty
    }, 0.0).(float64)
}

// ✅ BUENO: Código simple y claro
func calcularTotal(items []Item) float64 {
    var total float64
    for _, item := range items {
        total += item.Price * item.Quantity
    }
    return total
}
```

### 2. Nombres Descriptivos

```go
// ❌ MALO: Nombres confusos
func proc(d []byte) ([]byte, error) {
    var r []byte
    for _, b := range d {
        if b > 32 {
            r = append(r, b)
        }
    }
    return r, nil
}

// ✅ BUENO: Nombres que explican la intención
func eliminarCaracteresDeControl(data []byte) ([]byte, error) {
    var resultado []byte
    for _, byte := range data {
        if byte > 32 { // 32 = espacio ASCII
            resultado = append(resultado, byte)
        }
    }
    return resultado, nil
}
```

### 3. Funciones Pequeñas con Una Responsabilidad

```go
// ❌ MALO: Función que hace demasiado
func procesarFactura(data FacturaInput) error {
    // Validar datos
    if data.ClienteNombre == "" {
        return errors.New("nombre requerido")
    }
    if !validarCedula(data.ClienteCedula) {
        return errors.New("cédula inválida")
    }
    
    // Calcular totales
    var subtotal float64
    for _, producto := range data.Productos {
        subtotal += producto.Cantidad * producto.PrecioUnitario
    }
    iva := subtotal * 0.15
    total := subtotal + iva
    
    // Crear XML
    xml := fmt.Sprintf("<factura><total>%.2f</total></factura>", total)
    
    // Guardar en base de datos
    db.Save(xml)
    
    // Enviar por email
    sendEmail(data.ClienteEmail, xml)
    
    return nil
}

// ✅ BUENO: Funciones pequeñas con responsabilidades específicas
func procesarFactura(data FacturaInput) error {
    if err := validarDatosFactura(data); err != nil {
        return fmt.Errorf("validación falló: %w", err)
    }
    
    factura := crearFactura(data)
    
    if err := guardarFactura(factura); err != nil {
        return fmt.Errorf("error guardando: %w", err)
    }
    
    if err := enviarFacturaPorEmail(factura, data.ClienteEmail); err != nil {
        // Log error pero no fallar todo el proceso
        log.Printf("Error enviando email: %v", err)
    }
    
    return nil
}

func validarDatosFactura(data FacturaInput) error {
    if data.ClienteNombre == "" {
        return errors.New("nombre del cliente es requerido")
    }
    
    if !validarCedula(data.ClienteCedula) {
        return errors.New("cédula inválida")
    }
    
    if len(data.Productos) == 0 {
        return errors.New("debe incluir al menos un producto")
    }
    
    return nil
}

func crearFactura(data FacturaInput) *Factura {
    subtotal := calcularSubtotal(data.Productos)
    iva := calcularIVA(subtotal)
    
    return &Factura{
        Cliente:     data.ClienteNombre,
        Subtotal:    subtotal,
        IVA:         iva,
        Total:       subtotal + iva,
        FechaEmision: time.Now(),
    }
}
```

---

## Estructura de Proyecto

### Organización por Dominio (DDD)

Nuestro proyecto sigue Domain-Driven Design:

```
go-facturacion-sri/
├── api/              # HTTP handlers y middleware (interfaz)
├── config/           # Configuración externa
├── database/         # Persistencia de datos
├── factory/          # Creación de objetos de negocio
├── models/           # Entidades de dominio
├── sri/              # Lógica específica del SRI (dominio)
├── validators/       # Validaciones de negocio
└── docs/             # Documentación
```

### Principios de Organización

1. **Paquetes por funcionalidad**, no por tipo de archivo
2. **Dependencias hacia adentro** (Clean Architecture)
3. **Interfaces definidas donde se usan**
4. **Configuración externa al código**

### Ejemplo: Evolución de la Estructura

#### Antes (❌ Organizado por tipo)
```
proyecto/
├── handlers/
│   ├── factura_handler.go
│   ├── cliente_handler.go
│   └── producto_handler.go
├── models/
│   ├── factura.go
│   ├── cliente.go
│   └── producto.go
├── services/
│   ├── factura_service.go
│   ├── cliente_service.go
│   └── producto_service.go
└── repositories/
    ├── factura_repo.go
    ├── cliente_repo.go
    └── producto_repo.go
```

#### Después (✅ Organizado por dominio)
```
proyecto/
├── api/              # Capa de presentación
├── facturas/         # Dominio de facturas
│   ├── factura.go
│   ├── service.go
│   ├── repository.go
│   └── handler.go
├── clientes/         # Dominio de clientes
└── productos/        # Dominio de productos
```

---

## Manejo de Errores

### 1. Errores Descriptivos con Contexto

```go
// ❌ MALO: Error genérico sin contexto
func cargarCertificado(rutaArchivo string) (*Certificado, error) {
    data, err := os.ReadFile(rutaArchivo)
    if err != nil {
        return nil, err // ¿Qué archivo? ¿Por qué falló?
    }
    
    cert, err := parsearCertificado(data)
    if err != nil {
        return nil, err // ¿Qué tipo de error de parsing?
    }
    
    return cert, nil
}

// ✅ BUENO: Errores con contexto específico
func cargarCertificado(rutaArchivo string) (*Certificado, error) {
    data, err := os.ReadFile(rutaArchivo)
    if err != nil {
        return nil, fmt.Errorf("error leyendo certificado desde %s: %w", rutaArchivo, err)
    }
    
    cert, err := parsearCertificado(data)
    if err != nil {
        return nil, fmt.Errorf("error parseando certificado PKCS#12 desde %s: %w", rutaArchivo, err)
    }
    
    return cert, nil
}
```

### 2. Tipos de Error Específicos

```go
// ✅ BUENO: Errores tipados para diferentes escenarios
type ErrorValidacion struct {
    Campo   string
    Valor   interface{}
    Mensaje string
}

func (e ErrorValidacion) Error() string {
    return fmt.Sprintf("validación falló en campo '%s' con valor '%v': %s", 
        e.Campo, e.Valor, e.Mensaje)
}

type ErrorNegocio struct {
    Codigo  string
    Mensaje string
    Detalle string
}

func (e ErrorNegocio) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.Codigo, e.Mensaje, e.Detalle)
}

// Uso específico
func validarCedula(cedula string) error {
    if len(cedula) != 10 {
        return ErrorValidacion{
            Campo:   "cedula",
            Valor:   cedula,
            Mensaje: "debe tener exactamente 10 dígitos",
        }
    }
    
    if !algoritmoValidacionCedula(cedula) {
        return ErrorNegocio{
            Codigo:  "CEDULA_INVALIDA",
            Mensaje: "Cédula ecuatoriana inválida",
            Detalle: "El dígito verificador no coincide",
        }
    }
    
    return nil
}
```

### 3. Manejo de Errores en Capas

```go
// Capa de dominio: errores específicos del negocio
func (s *FacturaService) CrearFactura(input FacturaInput) (*Factura, error) {
    if err := s.validator.ValidarInput(input); err != nil {
        return nil, fmt.Errorf("datos inválidos: %w", err)
    }
    
    factura, err := s.factory.CrearFactura(input)
    if err != nil {
        return nil, fmt.Errorf("error creando factura: %w", err)
    }
    
    return factura, nil
}

// Capa de API: errores HTTP apropiados
func (h *FacturaHandler) CrearFactura(w http.ResponseWriter, r *http.Request) {
    var input FacturaInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    
    factura, err := h.service.CrearFactura(input)
    if err != nil {
        // Convertir errores de dominio a códigos HTTP
        var validationErr ErrorValidacion
        if errors.As(err, &validationErr) {
            writeErrorResponse(w, http.StatusBadRequest, validationErr.Error())
            return
        }
        
        var businessErr ErrorNegocio
        if errors.As(err, &businessErr) {
            writeErrorResponse(w, http.StatusUnprocessableEntity, businessErr.Error())
            return
        }
        
        // Error interno no esperado
        log.Printf("Error inesperado creando factura: %v", err)
        writeErrorResponse(w, http.StatusInternalServerError, "Error interno del servidor")
        return
    }
    
    writeJSONResponse(w, http.StatusCreated, factura)
}
```

---

## Testing Estratégico

### 1. Pirámide de Testing

```
       /\     E2E Tests
      /  \    (Pocos, lentos, frágiles)
     /____\   
    /      \  Integration Tests  
   /        \ (Algunos, medios)
  /__________\
 /            \ Unit Tests
/              \ (Muchos, rápidos, estables)
```

### 2. Unit Tests Efectivos

```go
// ✅ BUENO: Test con tabla de casos
func TestValidarCedula(t *testing.T) {
    tests := []struct {
        name     string
        cedula   string
        wantErr  bool
        errType  string
    }{
        {
            name:    "cédula válida",
            cedula:  "1713175071",
            wantErr: false,
        },
        {
            name:    "cédula muy corta",
            cedula:  "123456789",
            wantErr: true,
            errType: "validacion",
        },
        {
            name:    "cédula con letras",
            cedula:  "171317507A",
            wantErr: true,
            errType: "validacion",
        },
        {
            name:    "dígito verificador incorrecto",
            cedula:  "1713175070",
            wantErr: true,
            errType: "negocio",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidarCedula(tt.cedula)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("ValidarCedula() error = nil, esperaba error")
                    return
                }
                
                // Verificar tipo de error específico
                switch tt.errType {
                case "validacion":
                    var validationErr ErrorValidacion
                    if !errors.As(err, &validationErr) {
                        t.Errorf("esperaba ErrorValidacion, obtuvo %T", err)
                    }
                case "negocio":
                    var businessErr ErrorNegocio
                    if !errors.As(err, &businessErr) {
                        t.Errorf("esperaba ErrorNegocio, obtuvo %T", err)
                    }
                }
            } else {
                if err != nil {
                    t.Errorf("ValidarCedula() error = %v, esperaba nil", err)
                }
            }
        })
    }
}
```

### 3. Integration Tests con Dependencias Reales

```go
func TestFacturaAPI_Integration(t *testing.T) {
    // Setup: crear storage temporal para test
    testStorage := NewFacturaStorage()
    originalStorage := storage
    SetStorage(testStorage)
    defer SetStorage(originalStorage)
    
    // Setup: servidor de test
    server := NewServer("0") // puerto 0 = asignación automática
    go server.Start()
    defer server.Stop()
    
    // Test completo de flujo
    t.Run("crear y obtener factura", func(t *testing.T) {
        // 1. Crear factura
        facturaData := models.FacturaInput{
            ClienteNombre: "Test Cliente",
            ClienteCedula: "1713175071",
            Productos: []models.ProductoInput{
                {
                    Codigo:         "TEST001",
                    Descripcion:    "Producto Test",
                    Cantidad:       1,
                    PrecioUnitario: 100.0,
                },
            },
        }
        
        reqBody, _ := json.Marshal(CreateFacturaRequest{
            FacturaInput: facturaData,
            IncludeXML:   true,
        })
        
        resp, err := http.Post(
            fmt.Sprintf("http://localhost:%s/api/facturas", server.Port),
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        
        require.NoError(t, err)
        require.Equal(t, http.StatusCreated, resp.StatusCode)
        
        var facturaResp FacturaResponse
        err = json.NewDecoder(resp.Body).Decode(&facturaResp)
        require.NoError(t, err)
        
        // Verificaciones
        assert.NotEmpty(t, facturaResp.ID)
        assert.Equal(t, "Test Cliente", facturaResp.Factura.InfoFactura.RazonSocialComprador)
        assert.NotEmpty(t, facturaResp.XML)
        
        // 2. Obtener factura creada
        getResp, err := http.Get(
            fmt.Sprintf("http://localhost:%s/api/facturas/%s", server.Port, facturaResp.ID),
        )
        
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, getResp.StatusCode)
        
        var facturaObtenida FacturaResponse
        err = json.NewDecoder(getResp.Body).Decode(&facturaObtenida)
        require.NoError(t, err)
        
        assert.Equal(t, facturaResp.ID, facturaObtenida.ID)
    })
}
```

### 4. Benchmarks para Rendimiento

```go
func BenchmarkCrearFactura(b *testing.B) {
    facturaData := models.FacturaInput{
        ClienteNombre: "Cliente Benchmark",
        ClienteCedula: "1713175071",
        Productos: []models.ProductoInput{
            {
                Codigo:         "BENCH001",
                Descripcion:    "Producto Benchmark",
                Cantidad:       1,
                PrecioUnitario: 100.0,
            },
        },
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := factory.CrearFactura(facturaData)
        if err != nil {
            b.Fatalf("Error creando factura: %v", err)
        }
    }
}

func BenchmarkStorageOperations(b *testing.B) {
    storage := NewFacturaStorage()
    factura := FacturaResponse{
        ID: "BENCH-001",
        Factura: models.Factura{
            // ... datos de factura
        },
    }
    
    b.Run("Store", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            id := fmt.Sprintf("BENCH-%d", i)
            storage.Store(id, factura)
        }
    })
    
    b.Run("Get", func(b *testing.B) {
        // Pre-poblar storage
        for i := 0; i < 1000; i++ {
            id := fmt.Sprintf("BENCH-%d", i)
            storage.Store(id, factura)
        }
        
        b.ResetTimer()
        
        for i := 0; i < b.N; i++ {
            id := fmt.Sprintf("BENCH-%d", i%1000)
            _, _ = storage.Get(id)
        }
    })
}
```

---

## Concurrencia Segura

### 1. Principios de Diseño Concurrente

```go
// ✅ BUENO: Diseño que evita compartir estado
type FacturaProcessor struct {
    // Dependencias inmutables
    validator validators.FacturaValidator
    sriClient sri.Client
    storage   FacturaStorageInterface
}

func (fp *FacturaProcessor) ProcesarFactura(ctx context.Context, input FacturaInput) (*Factura, error) {
    // Cada procesamiento es independiente
    // No compartimos estado mutable entre goroutines
    
    if err := fp.validator.Validar(input); err != nil {
        return nil, err
    }
    
    factura := crear FacturaFromInput(input)
    
    // Usar contexto para timeouts y cancelación
    xmlFirmado, err := fp.sriClient.FirmarXML(ctx, factura.XML)
    if err != nil {
        return nil, err
    }
    
    factura.XMLFirmado = xmlFirmado
    
    // Storage thread-safe
    if err := fp.storage.Store(factura.ID, factura); err != nil {
        return nil, err
    }
    
    return factura, nil
}
```

### 2. Worker Pools para Procesamiento Masivo

```go
type BatchProcessor struct {
    numWorkers int
    storage    FacturaStorageInterface
}

func (bp *BatchProcessor) ProcesarLote(facturas []FacturaInput) error {
    // Canal para trabajo
    jobs := make(chan FacturaInput, len(facturas))
    
    // Canal para resultados
    results := make(chan error, len(facturas))
    
    // Iniciar workers
    for i := 0; i < bp.numWorkers; i++ {
        go bp.worker(jobs, results)
    }
    
    // Enviar trabajos
    for _, factura := range facturas {
        jobs <- factura
    }
    close(jobs)
    
    // Recoger resultados
    var errores []error
    for i := 0; i < len(facturas); i++ {
        if err := <-results; err != nil {
            errores = append(errores, err)
        }
    }
    
    if len(errores) > 0 {
        return fmt.Errorf("errores procesando lote: %v", errores)
    }
    
    return nil
}

func (bp *BatchProcessor) worker(jobs <-chan FacturaInput, results chan<- error) {
    for facturaInput := range jobs {
        _, err := factory.CrearFactura(facturaInput)
        results <- err
    }
}
```

### 3. Rate Limiting para APIs Externas

```go
type RateLimitedSRIClient struct {
    client sri.Client
    limiter *rate.Limiter
}

func NewRateLimitedSRIClient(client sri.Client, requestsPerSecond int) *RateLimitedSRIClient {
    return &RateLimitedSRIClient{
        client:  client,
        limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), 1),
    }
}

func (rl *RateLimitedSRIClient) AutorizarComprobante(ctx context.Context, claveAcceso string) error {
    // Esperar por permiso del rate limiter
    if err := rl.limiter.Wait(ctx); err != nil {
        return fmt.Errorf("rate limiting: %w", err)
    }
    
    return rl.client.AutorizarComprobante(ctx, claveAcceso)
}
```

---

## Interfaces Efectivas

### 1. Interfaces en el Lugar de Uso

```go
// ❌ MALO: Interfaz definida junto a implementación
// archivo: database/repository.go
type FacturaRepository interface {
    Save(*Factura) error
    FindByID(string) (*Factura, error)
}

type PostgresFacturaRepository struct {
    db *sql.DB
}

func (r *PostgresFacturaRepository) Save(f *Factura) error { ... }
func (r *PostgresFacturaRepository) FindByID(id string) (*Factura, error) { ... }

// ✅ BUENO: Interfaz definida donde se usa
// archivo: services/factura_service.go
type FacturaRepository interface {
    Save(*Factura) error
    FindByID(string) (*Factura, error)
}

type FacturaService struct {
    repo FacturaRepository // ← Usa la interfaz
}

// archivo: database/postgres_repository.go
// Implementa la interfaz implícitamente
type PostgresFacturaRepository struct {
    db *sql.DB
}

func (r *PostgresFacturaRepository) Save(f *Factura) error { ... }
func (r *PostgresFacturaRepository) FindByID(id string) (*Factura, error) { ... }
```

### 2. Interfaces Composables

```go
// Interfaces pequeñas y enfocadas
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// Componer según necesidades
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// Aplicado a nuestro dominio
type FacturaReader interface {
    GetFactura(id string) (*Factura, error)
    ListFacturas() ([]*Factura, error)
}

type FacturaWriter interface {
    SaveFactura(*Factura) error
    UpdateFactura(*Factura) error
}

type FacturaDeleter interface {
    DeleteFactura(id string) error
}

// Servicios componen según necesidad
type ReadOnlyFacturaService struct {
    reader FacturaReader // ← Solo puede leer
}

type FullFacturaService struct {
    FacturaReader
    FacturaWriter
    FacturaDeleter
}
```

---

## Rendimiento y Memoria

### 1. Evitar Allocaciones Innecesarias

```go
// ❌ MALO: Muchas allocaciones
func formatearFacturas(facturas []*Factura) string {
    var resultado string
    for _, factura := range facturas {
        resultado += fmt.Sprintf("ID: %s, Total: %.2f\n", factura.ID, factura.Total)
    }
    return resultado
}

// ✅ BUENO: Usar strings.Builder
func formatearFacturas(facturas []*Factura) string {
    var builder strings.Builder
    builder.Grow(len(facturas) * 50) // Pre-allocar espacio estimado
    
    for _, factura := range facturas {
        fmt.Fprintf(&builder, "ID: %s, Total: %.2f\n", factura.ID, factura.Total)
    }
    
    return builder.String()
}
```

### 2. Reutilizar Buffers

```go
// Pool de buffers para reutilización
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024) // 1KB inicial
    },
}

func procesarXML(xmlData []byte) ([]byte, error) {
    // Obtener buffer del pool
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer[:0]) // Devolver al pool (con longitud 0)
    
    // Usar buffer para procesamiento
    buffer = append(buffer, xmlData...)
    
    // ... procesar ...
    
    // Crear copia para retornar (el buffer vuelve al pool)
    result := make([]byte, len(buffer))
    copy(result, buffer)
    
    return result, nil
}
```

### 3. Lazy Loading y Cacheing

```go
type FacturaService struct {
    repo  FacturaRepository
    cache map[string]*Factura
    mu    sync.RWMutex
}

func (fs *FacturaService) GetFactura(id string) (*Factura, error) {
    // Verificar cache primero
    fs.mu.RLock()
    if factura, exists := fs.cache[id]; exists {
        fs.mu.RUnlock()
        return factura, nil
    }
    fs.mu.RUnlock()
    
    // No está en cache, cargar de repositorio
    factura, err := fs.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Agregar a cache
    fs.mu.Lock()
    fs.cache[id] = factura
    fs.mu.Unlock()
    
    return factura, nil
}
```

---

## Seguridad

### 1. Validación de Input

```go
type FacturaInput struct {
    ClienteNombre string    `json:"clienteNombre" validate:"required,max=100"`
    ClienteCedula string    `json:"clienteCedula" validate:"required,cedula"`
    ClienteEmail  string    `json:"clienteEmail" validate:"omitempty,email"`
    Productos     []Producto `json:"productos" validate:"required,min=1,dive"`
}

type Producto struct {
    Codigo         string  `json:"codigo" validate:"required,alphanum,max=20"`
    Descripcion    string  `json:"descripcion" validate:"required,max=200"`
    Cantidad       float64 `json:"cantidad" validate:"required,gt=0"`
    PrecioUnitario float64 `json:"precioUnitario" validate:"required,gt=0"`
}

func validarInput(input FacturaInput) error {
    validate := validator.New()
    
    // Registrar validador personalizado para cédula
    validate.RegisterValidation("cedula", validarCedulaEcuatoriana)
    
    return validate.Struct(input)
}

func validarCedulaEcuatoriana(fl validator.FieldLevel) bool {
    cedula := fl.Field().String()
    return validators.ValidarCedula(cedula) == nil
}
```

### 2. Sanitización y Escape

```go
import "html"

func sanitizarInput(input string) string {
    // Escape HTML
    sanitized := html.EscapeString(input)
    
    // Remover caracteres de control
    sanitized = strings.Map(func(r rune) rune {
        if r < 32 && r != '\t' && r != '\n' && r != '\r' {
            return -1 // Remover carácter
        }
        return r
    }, sanitized)
    
    return strings.TrimSpace(sanitized)
}

func (h *FacturaHandler) CrearFactura(w http.ResponseWriter, r *http.Request) {
    var input FacturaInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        writeErrorResponse(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    
    // Sanitizar inputs
    input.ClienteNombre = sanitizarInput(input.ClienteNombre)
    for i := range input.Productos {
        input.Productos[i].Descripcion = sanitizarInput(input.Productos[i].Descripcion)
    }
    
    // ... continuar procesamiento
}
```

### 3. Manejo Seguro de Certificados

```go
type CertificadoManager struct {
    certificados map[string]*CertificadoDigital
    mu          sync.RWMutex
}

func (cm *CertificadoManager) CargarCertificado(rutaArchivo, password string) error {
    // Validar ruta de archivo (prevenir path traversal)
    cleanPath := filepath.Clean(rutaArchivo)
    if !strings.HasPrefix(cleanPath, "/ruta/segura/certificados/") {
        return errors.New("ruta de certificado no autorizada")
    }
    
    // Leer archivo de forma segura
    data, err := os.ReadFile(cleanPath)
    if err != nil {
        return fmt.Errorf("error leyendo certificado: %w", err)
    }
    
    // Limpiar password de memoria después de uso
    defer func() {
        // Sobreescribir password en memoria
        for i := range password {
            password = password[:i] + "x" + password[i+1:]
        }
    }()
    
    cert, err := pkcs12.DecodeChain(data, password)
    if err != nil {
        return fmt.Errorf("error decodificando certificado: %w", err)
    }
    
    // Verificar validez del certificado
    if time.Now().After(cert.NotAfter) {
        return errors.New("certificado expirado")
    }
    
    cm.mu.Lock()
    cm.certificados[cleanPath] = &CertificadoDigital{
        Archivo: cleanPath,
        Cert:    cert,
    }
    cm.mu.Unlock()
    
    return nil
}
```

---

## Documentación

### 1. Comentarios de Código Efectivos

```go
// ✅ BUENO: Explica el "por qué", no el "qué"

// GenerarClaveAcceso crea una clave de acceso de 49 dígitos según especificaciones SRI.
// La clave incluye fecha, RUC, tipo de documento, serie, secuencial y dígito verificador
// calculado con algoritmo módulo 11 como requiere el SRI de Ecuador.
func GenerarClaveAcceso(fecha time.Time, ruc, tipoDoc, establecimiento, ptoEmision, secuencial string) string {
    // Formato: ddMMyyyyTTrrrrrrrrrrrraeeeeeeNNNNNNNNNdv
    // TT = tipo documento, rrr = RUC, a = ambiente, etc.
    
    fechaStr := fecha.Format("02012006")
    
    // El ambiente se obtiene de configuración porque cambia entre desarrollo y producción
    ambiente := config.Config.Ambiente.Codigo
    
    clave := fechaStr + tipoDoc + ruc + ambiente + establecimiento + ptoEmision + secuencial
    
    // Calcular dígito verificador usando módulo 11 (algoritmo oficial SRI)
    digitoVerificador := calcularModulo11(clave)
    
    return clave + strconv.Itoa(digitoVerificador)
}

// calcularModulo11 implementa el algoritmo de módulo 11 específico del SRI.
// Este algoritmo es diferente al módulo 11 estándar porque usa una secuencia
// específica de multiplicadores: 2,3,4,5,6,7,2,3,4...
func calcularModulo11(clave string) int {
    multiplicadores := []int{2, 3, 4, 5, 6, 7}
    suma := 0
    
    // Recorrer desde la derecha
    for i := len(clave) - 1; i >= 0; i-- {
        digito, _ := strconv.Atoi(string(clave[i]))
        multiplicador := multiplicadores[(len(clave)-1-i)%len(multiplicadores)]
        suma += digito * multiplicador
    }
    
    residuo := suma % 11
    
    // Casos especiales del algoritmo SRI
    switch residuo {
    case 0:
        return 0
    case 1:
        return 1
    default:
        return 11 - residuo
    }
}
```

### 2. Documentación de API

```go
// Package api proporciona handlers HTTP para el sistema de facturación SRI.
//
// La API sigue principios REST y maneja:
//   - Creación de facturas electrónicas
//   - Consulta de facturas existentes
//   - Generación de XML firmado para SRI
//
// Todos los endpoints requieren content-type application/json y retornan
// respuestas en formato JSON con códigos HTTP estándar.
//
// Ejemplo de uso:
//
//	POST /api/facturas
//	{
//	  "clienteNombre": "Juan Pérez",
//	  "clienteCedula": "1713175071",
//	  "productos": [
//	    {
//	      "codigo": "PROD001",
//	      "descripcion": "Producto ejemplo",
//	      "cantidad": 1,
//	      "precioUnitario": 100.00
//	    }
//	  ]
//	}
package api

// CreateFacturaRequest define la estructura para crear una nueva factura.
//
// Campos requeridos:
//   - ClienteNombre: Nombre completo del cliente
//   - ClienteCedula: Cédula ecuatoriana válida (10 dígitos)
//   - Productos: Al menos un producto con código, descripción, cantidad y precio
//
// Campos opcionales:
//   - IncludeXML: Si true, incluye el XML generado en la respuesta
type CreateFacturaRequest struct {
    models.FacturaInput
    IncludeXML bool `json:"includeXML,omitempty"`
}
```

### 3. README Comprensivo

```markdown
# Go Facturación SRI

Sistema de facturación electrónica para el SRI de Ecuador desarrollado en Go.

## Características

- ✅ Generación de facturas según especificaciones SRI
- ✅ Firmas digitales XAdES-BES
- ✅ API REST completa
- ✅ Validación de cédulas ecuatorianas
- ✅ Manejo de certificados PKCS#12
- ✅ Thread-safe storage
- ✅ Testing extensivo (45.5% cobertura)

## Inicio Rápido

```bash
# Clonar repositorio
git clone https://github.com/usuario/go-facturacion-sri
cd go-facturacion-sri

# Instalar dependencias
go mod download

# Ejecutar tests
go test ./...

# Iniciar API
go run main.go test_validaciones.go api
```

## Arquitectura

El proyecto sigue principios de Clean Architecture con separación clara de responsabilidades:

- `api/`: Handlers HTTP y middleware
- `models/`: Entidades de dominio
- `factory/`: Creación de objetos de negocio
- `validators/`: Validaciones de negocio
- `sri/`: Integración específica con SRI
- `config/`: Configuración externa

## Configuración

Crear `config/desarrollo.json`:

```json
{
  "empresa": {
    "razonSocial": "Mi Empresa SAS",
    "ruc": "1792146739001",
    "establecimiento": "001",
    "puntoEmision": "001",
    "direccion": "Quito, Ecuador"
  },
  "ambiente": {
    "codigo": "1",
    "descripcion": "Pruebas",
    "tipoEmision": "1"
  }
}
```

## Uso de la API

### Crear Factura

```bash
curl -X POST http://localhost:8080/api/facturas \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "Juan Pérez",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "PROD001",
        "descripcion": "Laptop Dell",
        "cantidad": 1,
        "precioUnitario": 800.00
      }
    ],
    "includeXML": true
  }'
```

### Obtener Factura

```bash
curl http://localhost:8080/api/facturas/FAC-000001
```

## Contribuir

1. Fork el proyecto
2. Crear rama de feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.
```

---

## Herramientas de Desarrollo

### 1. Makefile para Automatización

```makefile
# Makefile para proyecto Go

.PHONY: build test lint fmt vet coverage clean

# Configuración
BINARY_NAME=facturacion-sri
BUILD_DIR=./bin
COVERAGE_FILE=coverage.out

# Build
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go test_validaciones.go

# Testing
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE)

test-race:
	go test -race ./...

benchmark:
	go test -bench=. -benchmem ./...

# Code Quality
lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

# Development
dev-api:
	go run main.go test_validaciones.go api

dev-sri:
	go run main.go test_validaciones.go sri

# Cleanup
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE)

# Dependencies
deps:
	go mod download
	go mod tidy

# Install tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Help
help:
	@echo "Comandos disponibles:"
	@echo "  build         - Compilar binario"
	@echo "  test          - Ejecutar tests"
	@echo "  test-coverage - Tests con reporte de cobertura"
	@echo "  test-race     - Tests con detector de race conditions"
	@echo "  benchmark     - Ejecutar benchmarks"
	@echo "  lint          - Ejecutar linter"
	@echo "  fmt           - Formatear código"
	@echo "  vet           - Ejecutar go vet"
	@echo "  dev-api       - Ejecutar en modo API"
	@echo "  dev-sri       - Ejecutar demo SRI"
	@echo "  clean         - Limpiar archivos generados"
	@echo "  deps          - Instalar dependencias"
	@echo "  install-tools - Instalar herramientas de desarrollo"
```

### 2. Configuración de CI/CD (.github/workflows/ci.yml)

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.24'
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
    
    - name: Run go vet
      run: go vet ./...
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
    
    - name: Build
      run: go build -v ./...

  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: './...'
```

### 3. Configuración de Linter (.golangci.yml)

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - varcheck
    - deadcode
    - structcheck
    - misspell
    - unconvert
    - gofmt
    - goimports
    - golint
    - gocritic
    - goprintffuncname
    - gosec

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  govet:
    check-shadowing: true
  
  golint:
    min-confidence: 0.8
  
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance
    
    disabled-checks:
      - unnamedResult
      - whyNoLint

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - errcheck
```

---

## 🎯 Resumen de Mejores Prácticas

### Código
- ✅ Simplicidad sobre cleverness
- ✅ Nombres descriptivos
- ✅ Funciones pequeñas con una responsabilidad
- ✅ Manejo de errores con contexto
- ✅ Evitar allocaciones innecesarias

### Arquitectura
- ✅ Organización por dominio
- ✅ Interfaces pequeñas y composables
- ✅ Dependencias hacia adentro
- ✅ Configuración externa

### Concurrencia
- ✅ Evitar compartir estado mutable
- ✅ Usar channels para comunicación
- ✅ Proteger estado compartido con mutex
- ✅ Contextos para timeouts y cancelación

### Testing
- ✅ Pirámide de testing (muchos unit, pocos E2E)
- ✅ Tests con tabla de casos
- ✅ Mocks e interfaces para dependencias
- ✅ Benchmarks para código crítico

### Herramientas
- ✅ Makefile para automatización
- ✅ CI/CD pipeline
- ✅ Linters y formatters
- ✅ Detector de race conditions

**¡Estas prácticas convierten código funcional en código profesional y mantenible!**