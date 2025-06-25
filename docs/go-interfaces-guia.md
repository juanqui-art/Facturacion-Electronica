# Guía de Interfaces en Go - Desacoplamiento y Flexibilidad

Esta guía te enseña interfaces en Go desde cero, usando ejemplos reales del proyecto de facturación SRI.

## 📚 Tabla de Contenido

1. [¿Qué son las Interfaces?](#qué-son-las-interfaces)
2. [Implementación Implícita](#implementación-implícita)
3. [Interfaces Pequeñas vs Grandes](#interfaces-pequeñas-vs-grandes)
4. [Desacoplamiento en Práctica](#desacoplamiento-en-práctica)
5. [Patrones de Diseño](#patrones-de-diseño)
6. [Implementación Real en el Proyecto](#implementación-real-en-el-proyecto)
7. [Testing con Interfaces](#testing-con-interfaces)
8. [Ejercicios Prácticos](#ejercicios-prácticos)

---

## ¿Qué son las Interfaces?

Una **interfaz** en Go es un contrato que define qué métodos debe tener un tipo, pero no cómo los implementa.

### Analogía Simple
Imagina que necesitas algo que "pueda volar":
- **Interfaz**: "Todo lo que vuele debe tener método `Volar()`"
- **Implementaciones**: Avión, pájaro, drone (cada uno vuela diferente)
- **Tu código**: Solo necesita saber que "puede volar", no cómo

### Ejemplo Básico
```go
// Definir la interfaz
type Volador interface {
    Volar() string
}

// Implementaciones diferentes
type Avion struct {
    Modelo string
}

func (a Avion) Volar() string {
    return fmt.Sprintf("Avión %s volando con motores", a.Modelo)
}

type Pajaro struct {
    Especie string
}

func (p Pajaro) Volar() string {
    return fmt.Sprintf("Pájaro %s volando con alas", p.Especie)
}

// Función que usa la interfaz
func hacerVolar(v Volador) {
    fmt.Println(v.Volar())
}

func main() {
    avion := Avion{Modelo: "Boeing 747"}
    pajaro := Pajaro{Especie: "Águila"}
    
    hacerVolar(avion)  // Funciona
    hacerVolar(pajaro) // También funciona
}
```

---

## Implementación Implícita

En Go, **no necesitas declarar** que implementas una interfaz. Si tienes los métodos, automáticamente la implementas.

### Comparación con Otros Lenguajes

#### Java (Explícito)
```java
interface Volador {
    String volar();
}

class Avion implements Volador {  // ← Declaración explícita
    public String volar() {
        return "Volando con motores";
    }
}
```

#### Go (Implícito)
```go
type Volador interface {
    Volar() string
}

type Avion struct {}

// Solo implementa el método, Go automáticamente
// reconoce que Avion implementa Volador
func (a Avion) Volar() string {
    return "Volando con motores"
}
```

### Ventajas de la Implementación Implícita

1. **Flexibilidad**: Puedes hacer que tipos existentes implementen nuevas interfaces
2. **Desacoplamiento**: Las interfaces se pueden definir donde se necesitan
3. **Evolución**: Agregar interfaces no requiere cambiar código existente

### Ejemplo Práctico
```go
// Interfaz estándar de Go
type Stringer interface {
    String() string
}

// Mi tipo
type Persona struct {
    Nombre string
    Edad   int
}

// Implemento String() y automáticamente implemento Stringer
func (p Persona) String() string {
    return fmt.Sprintf("%s (%d años)", p.Nombre, p.Edad)
}

func main() {
    p := Persona{Nombre: "Ana", Edad: 30}
    
    // Ahora fmt.Println usa automáticamente mi método String()
    fmt.Println(p) // Output: Ana (30 años)
}
```

---

## Interfaces Pequeñas vs Grandes

### Principio de Go: "Interfaces pequeñas"

```go
// ✅ BUENO: Interfaz pequeña, enfocada
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// ✅ BUENO: Componer interfaces pequeñas
type ReadWriter interface {
    Reader
    Writer
}

// ❌ MALO: Interfaz grande, múltiples responsabilidades
type GiganteFeasomeInterface interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    Seek(int64, int) (int64, error)
    Stat() (FileInfo, error)
    Sync() error
    Truncate(int64) error
    // ... 20 métodos más
}
```

### ¿Por qué Interfaces Pequeñas?

1. **Fáciles de implementar**: Menos métodos = menos trabajo
2. **Más reutilizables**: Más tipos pueden implementarlas
3. **Composables**: Se pueden combinar para crear interfaces más grandes
4. **Testeable**: Fácil crear mocks

### Ejemplo Comparativo
```go
// ❌ Interfaz grande - difícil de implementar
type FacturaServiceCompleto interface {
    CrearFactura(data FacturaInput) (*Factura, error)
    ActualizarFactura(id string, data FacturaInput) error
    EliminarFactura(id string) error
    ObtenerFactura(id string) (*Factura, error)
    ListarFacturas() ([]*Factura, error)
    ExportarPDF(id string) ([]byte, error)
    ExportarXML(id string) ([]byte, error)
    EnviarPorEmail(id string, email string) error
    ValidarConSRI(id string) error
    FirmarDigitalmente(id string) error
}

// ✅ Interfaces pequeñas - fáciles de implementar
type FacturaCreator interface {
    CrearFactura(data FacturaInput) (*Factura, error)
}

type FacturaReader interface {
    ObtenerFactura(id string) (*Factura, error)
    ListarFacturas() ([]*Factura, error)
}

type FacturaExporter interface {
    ExportarPDF(id string) ([]byte, error)
    ExportarXML(id string) ([]byte, error)
}

// Se pueden componer cuando necesites múltiples capacidades
type FacturaService interface {
    FacturaCreator
    FacturaReader
    FacturaExporter
}
```

---

## Desacoplamiento en Práctica

**Desacoplamiento** significa que tus componentes no dependan directamente uno del otro.

### Ejemplo: Sistema de Notificaciones

#### Sin Desacoplamiento (❌ Acoplado)
```go
type EmailService struct {
    SMTPServer string
}

func (e *EmailService) EnviarEmail(to, subject, body string) error {
    // Lógica específica de email
    return nil
}

// Esta función está ACOPLADA a EmailService
func NotificarUsuario(usuario string, mensaje string) {
    emailService := &EmailService{SMTPServer: "smtp.gmail.com"}
    emailService.EnviarEmail(usuario, "Notificación", mensaje)
    
    // ¿Qué pasa si después quiero usar SMS? ¿WhatsApp?
    // Tendría que cambiar esta función ❌
}
```

#### Con Desacoplamiento (✅ Flexible)
```go
// 1. Definir interfaz
type Notificador interface {
    Enviar(destinatario, mensaje string) error
}

// 2. Implementaciones específicas
type EmailNotificador struct {
    SMTPServer string
}

func (e *EmailNotificador) Enviar(destinatario, mensaje string) error {
    fmt.Printf("📧 Email a %s: %s\n", destinatario, mensaje)
    return nil
}

type SMSNotificador struct {
    APIKey string
}

func (s *SMSNotificador) Enviar(destinatario, mensaje string) error {
    fmt.Printf("📱 SMS a %s: %s\n", destinatario, mensaje)
    return nil
}

type WhatsAppNotificador struct {
    Token string
}

func (w *WhatsAppNotificador) Enviar(destinatario, mensaje string) error {
    fmt.Printf("💬 WhatsApp a %s: %s\n", destinatario, mensaje)
    return nil
}

// 3. Función desacoplada - acepta cualquier Notificador
func NotificarUsuario(notificador Notificador, usuario, mensaje string) {
    notificador.Enviar(usuario, mensaje)
}

// 4. Uso flexible
func main() {
    // Cambiar tipo de notificación sin modificar NotificarUsuario
    emailNotif := &EmailNotificador{SMTPServer: "smtp.gmail.com"}
    smsNotif := &SMSNotificador{APIKey: "abc123"}
    whatsappNotif := &WhatsAppNotificador{Token: "xyz789"}
    
    NotificarUsuario(emailNotif, "ana@email.com", "Bienvenida")
    NotificarUsuario(smsNotif, "+593999123456", "Código: 1234")
    NotificarUsuario(whatsappNotif, "+593999123456", "Hola!")
}
```

---

## Patrones de Diseño

### 1. Patrón Strategy (Estrategia)
Cambiar algoritmos sin modificar código:

```go
type CalculadorImpuestos interface {
    Calcular(monto float64) float64
}

type IVAEcuador struct{}
func (i IVAEcuador) Calcular(monto float64) float64 {
    return monto * 0.15 // 15% IVA Ecuador
}

type IVAColumbia struct{}
func (i IVAColumbia) Calcular(monto float64) float64 {
    return monto * 0.19 // 19% IVA Colombia
}

type FacturaProcessor struct {
    calculadora CalculadorImpuestos
}

func (fp *FacturaProcessor) ProcesarFactura(monto float64) float64 {
    impuesto := fp.calculadora.Calcular(monto)
    return monto + impuesto
}

// Uso
processor := &FacturaProcessor{calculadora: IVAEcuador{}}
total := processor.ProcesarFactura(100.0) // 115.0
```

### 2. Patrón Decorator (Decorador)
Agregar funcionalidades sin modificar código original:

```go
type Logger interface {
    Log(mensaje string)
}

// Implementación básica
type SimpleLogger struct{}
func (s SimpleLogger) Log(mensaje string) {
    fmt.Println(mensaje)
}

// Decorador que agrega timestamp
type TimestampLogger struct {
    logger Logger
}

func (t TimestampLogger) Log(mensaje string) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    t.logger.Log(fmt.Sprintf("[%s] %s", timestamp, mensaje))
}

// Decorador que agrega nivel
type LevelLogger struct {
    logger Logger
    level  string
}

func (l LevelLogger) Log(mensaje string) {
    l.logger.Log(fmt.Sprintf("[%s] %s", l.level, mensaje))
}

// Uso - se pueden componer decoradores
func main() {
    simple := SimpleLogger{}
    withTime := TimestampLogger{logger: simple}
    withTimeAndLevel := LevelLogger{logger: withTime, level: "INFO"}
    
    withTimeAndLevel.Log("Sistema iniciado")
    // Output: [INFO] [2024-06-24 15:30:45] Sistema iniciado
}
```

### 3. Patrón Dependency Injection
Inyectar dependencias a través de interfaces:

```go
type Database interface {
    GuardarFactura(factura *Factura) error
    ObtenerFactura(id string) (*Factura, error)
}

type EmailSender interface {
    EnviarFactura(email string, factura *Factura) error
}

// Servicio que depende de interfaces, no implementaciones
type FacturaService struct {
    db     Database
    mailer EmailSender
}

func NewFacturaService(db Database, mailer EmailSender) *FacturaService {
    return &FacturaService{db: db, mailer: mailer}
}

func (fs *FacturaService) CrearYEnviarFactura(data FacturaInput, email string) error {
    factura := crearFacturaDe(data)
    
    // Guardar usando la interfaz (puede ser cualquier DB)
    if err := fs.db.GuardarFactura(factura); err != nil {
        return err
    }
    
    // Enviar usando la interfaz (puede ser cualquier mailer)
    return fs.mailer.EnviarFactura(email, factura)
}

// En producción
postgresDB := &PostgresDatabase{}
sendgridMailer := &SendgridEmailSender{}
service := NewFacturaService(postgresDB, sendgridMailer)

// En testing
mockDB := &MockDatabase{}
mockMailer := &MockEmailSender{}
testService := NewFacturaService(mockDB, mockMailer)
```

---

## Implementación Real en el Proyecto

Así implementamos interfaces en nuestro proyecto de facturación SRI:

### Problema Original (❌ Acoplado)
```go
// Código directamente acoplado a implementación específica
var facturaStorage = make(map[string]FacturaResponse)
var nextID = 1

func handleCreateFactura() {
    // Dependencia directa del storage específico
    id := fmt.Sprintf("FAC-%06d", nextID)
    nextID++
    facturaStorage[id] = factura // ← Acoplado a map
}
```

### Solución con Interfaces (✅ Desacoplado)

#### 1. Definir Interfaz
```go
// Interfaz que define el contrato
type FacturaStorageInterface interface {
    Store(id string, factura FacturaResponse)
    Get(id string) (FacturaResponse, bool)
    GetAll() []FacturaResponse
    GetNextID() int
    Count() int
}
```

#### 2. Implementación en Memoria
```go
type FacturaStorage struct {
    mu       sync.RWMutex
    facturas map[string]FacturaResponse
    nextID   int
}

func NewFacturaStorage() *FacturaStorage {
    return &FacturaStorage{
        facturas: make(map[string]FacturaResponse),
        nextID:   1,
    }
}

func (fs *FacturaStorage) Store(id string, factura FacturaResponse) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    fs.facturas[id] = factura
}

func (fs *FacturaStorage) Get(id string) (FacturaResponse, bool) {
    fs.mu.RLock()
    defer fs.mu.RUnlock()
    factura, exists := fs.facturas[id]
    return factura, exists
}

// ... otros métodos
```

#### 3. Implementación con Logging (Patrón Decorator)
```go
type LoggingFacturaStorage struct {
    underlying FacturaStorageInterface
}

func NewLoggingFacturaStorage(underlying FacturaStorageInterface) *LoggingFacturaStorage {
    return &LoggingFacturaStorage{underlying: underlying}
}

func (lfs *LoggingFacturaStorage) Store(id string, factura FacturaResponse) {
    fmt.Printf("🗄️ Almacenando factura: %s\n", id)
    lfs.underlying.Store(id, factura)
}

func (lfs *LoggingFacturaStorage) Get(id string) (FacturaResponse, bool) {
    fmt.Printf("🔍 Buscando factura: %s\n", id)
    factura, exists := lfs.underlying.Get(id)
    if !exists {
        fmt.Printf("❌ Factura no encontrada: %s\n", id)
    }
    return factura, exists
}

// ... otros métodos con logging
```

#### 4. Uso Flexible
```go
// Variable global usa la interfaz
var storage FacturaStorageInterface = NewFacturaStorage()

// Función para cambiar implementación
func SetStorage(s FacturaStorageInterface) {
    storage = s
}

// En código de producción
func handleCreateFactura() {
    // El código no sabe qué implementación se usa
    nextID := storage.GetNextID()
    id := fmt.Sprintf("FAC-%06d", nextID)
    
    response := FacturaResponse{ID: id, /* ... */}
    storage.Store(id, response) // ← Funciona con cualquier implementación
}

// Configuración flexible
func main() {
    if os.Getenv("ENABLE_LOGGING") == "true" {
        basicStorage := NewFacturaStorage()
        storage = NewLoggingFacturaStorage(basicStorage)
    } else {
        storage = NewFacturaStorage()
    }
}
```

### Beneficios Obtenidos

1. **Flexibilidad**: Cambiar de memoria a base de datos sin tocar handlers
2. **Testing**: Usar mock storage para tests
3. **Logging**: Agregar logging sin modificar código existente
4. **Mantenimiento**: Cambios en storage no afectan resto del código

---

## Testing con Interfaces

Las interfaces hacen el testing mucho más fácil:

### Sin Interfaces (❌ Difícil de Testear)
```go
func ProcessarPago(monto float64) error {
    // Dependencia directa - difícil de testear
    paypalClient := &PayPalClient{
        APIKey: "real_api_key_12345",
        URL:    "https://api.paypal.com",
    }
    
    return paypalClient.CargarTarjeta(monto)
}

// Test problemático
func TestProcessarPago(t *testing.T) {
    // ¿Cómo testeo esto sin hacer pagos reales? 😱
    err := ProcessarPago(100.0)
    // Este test haría cargo real a PayPal
}
```

### Con Interfaces (✅ Fácil de Testear)
```go
// 1. Definir interfaz
type PagosProcessor interface {
    CargarTarjeta(monto float64) error
}

// 2. Implementación real
type PayPalClient struct {
    APIKey string
    URL    string
}

func (p *PayPalClient) CargarTarjeta(monto float64) error {
    // Lógica real de PayPal
    return nil
}

// 3. Función que usa interfaz
func ProcessarPago(processor PagosProcessor, monto float64) error {
    return processor.CargarTarjeta(monto)
}

// 4. Mock para testing
type MockPagosProcessor struct {
    ShouldFail bool
    AmountPaid float64
}

func (m *MockPagosProcessor) CargarTarjeta(monto float64) error {
    m.AmountPaid = monto
    if m.ShouldFail {
        return errors.New("pago falló")
    }
    return nil
}

// 5. Test fácil y controlado
func TestProcessarPago(t *testing.T) {
    tests := []struct {
        name       string
        monto      float64
        shouldFail bool
        expectErr  bool
    }{
        {"pago exitoso", 100.0, false, false},
        {"pago fallido", 50.0, true, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &MockPagosProcessor{ShouldFail: tt.shouldFail}
            
            err := ProcessarPago(mock, tt.monto)
            
            if tt.expectErr && err == nil {
                t.Error("esperaba error")
            }
            if !tt.expectErr && err != nil {
                t.Errorf("no esperaba error: %v", err)
            }
            if mock.AmountPaid != tt.monto {
                t.Errorf("monto esperado: %v, obtuvo: %v", tt.monto, mock.AmountPaid)
            }
        })
    }
}
```

### Testing en Nuestro Proyecto
```go
func TestFacturaStorage(t *testing.T) {
    // Usar implementación real para test
    storage := NewFacturaStorage()
    
    factura := FacturaResponse{ID: "TEST-001"}
    
    // Test Store
    storage.Store("TEST-001", factura)
    
    // Test Get
    retrieved, exists := storage.Get("TEST-001")
    if !exists {
        t.Error("factura debería existir")
    }
    if retrieved.ID != "TEST-001" {
        t.Error("ID no coincide")
    }
    
    // Test Count
    if storage.Count() != 1 {
        t.Error("count debería ser 1")
    }
}

func TestLoggingStorage(t *testing.T) {
    // Test del decorador
    base := NewFacturaStorage()
    logging := NewLoggingFacturaStorage(base)
    
    factura := FacturaResponse{ID: "LOG-001"}
    
    // Capturar output de logging
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    logging.Store("LOG-001", factura)
    
    w.Close()
    os.Stdout = oldStdout
    
    output, _ := io.ReadAll(r)
    if !strings.Contains(string(output), "Almacenando factura: LOG-001") {
        t.Error("debería logear el almacenamiento")
    }
}
```

---

## Ejercicios Prácticos

### Ejercicio 1: Sistema de Archivos
Implementa un sistema de archivos con diferentes tipos de storage:

```go
type FileStorage interface {
    Save(filename string, content []byte) error
    Load(filename string) ([]byte, error)
    Delete(filename string) error
    List() ([]string, error)
}

// Implementa:
// 1. LocalFileStorage (guarda en disco)
// 2. MemoryFileStorage (guarda en memoria)
// 3. LoggingFileStorage (decorador que logea operaciones)

type LocalFileStorage struct {
    basePath string
}

type MemoryFileStorage struct {
    files map[string][]byte
}

type LoggingFileStorage struct {
    underlying FileStorage
}
```

<details>
<summary>👁️ Ver Solución</summary>

```go
import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

type LocalFileStorage struct {
    basePath string
}

func NewLocalFileStorage(basePath string) *LocalFileStorage {
    os.MkdirAll(basePath, 0755)
    return &LocalFileStorage{basePath: basePath}
}

func (lfs *LocalFileStorage) Save(filename string, content []byte) error {
    path := filepath.Join(lfs.basePath, filename)
    return ioutil.WriteFile(path, content, 0644)
}

func (lfs *LocalFileStorage) Load(filename string) ([]byte, error) {
    path := filepath.Join(lfs.basePath, filename)
    return ioutil.ReadFile(path)
}

func (lfs *LocalFileStorage) Delete(filename string) error {
    path := filepath.Join(lfs.basePath, filename)
    return os.Remove(path)
}

func (lfs *LocalFileStorage) List() ([]string, error) {
    files, err := ioutil.ReadDir(lfs.basePath)
    if err != nil {
        return nil, err
    }
    
    var names []string
    for _, file := range files {
        if !file.IsDir() {
            names = append(names, file.Name())
        }
    }
    return names, nil
}

type MemoryFileStorage struct {
    files map[string][]byte
}

func NewMemoryFileStorage() *MemoryFileStorage {
    return &MemoryFileStorage{
        files: make(map[string][]byte),
    }
}

func (mfs *MemoryFileStorage) Save(filename string, content []byte) error {
    mfs.files[filename] = make([]byte, len(content))
    copy(mfs.files[filename], content)
    return nil
}

func (mfs *MemoryFileStorage) Load(filename string) ([]byte, error) {
    content, exists := mfs.files[filename]
    if !exists {
        return nil, errors.New("archivo no encontrado")
    }
    result := make([]byte, len(content))
    copy(result, content)
    return result, nil
}

func (mfs *MemoryFileStorage) Delete(filename string) error {
    if _, exists := mfs.files[filename]; !exists {
        return errors.New("archivo no encontrado")
    }
    delete(mfs.files, filename)
    return nil
}

func (mfs *MemoryFileStorage) List() ([]string, error) {
    var names []string
    for filename := range mfs.files {
        names = append(names, filename)
    }
    return names, nil
}

type LoggingFileStorage struct {
    underlying FileStorage
}

func NewLoggingFileStorage(underlying FileStorage) *LoggingFileStorage {
    return &LoggingFileStorage{underlying: underlying}
}

func (lfs *LoggingFileStorage) Save(filename string, content []byte) error {
    fmt.Printf("💾 Guardando archivo: %s (%d bytes)\n", filename, len(content))
    return lfs.underlying.Save(filename, content)
}

func (lfs *LoggingFileStorage) Load(filename string) ([]byte, error) {
    fmt.Printf("📖 Cargando archivo: %s\n", filename)
    return lfs.underlying.Load(filename)
}

func (lfs *LoggingFileStorage) Delete(filename string) error {
    fmt.Printf("🗑️ Eliminando archivo: %s\n", filename)
    return lfs.underlying.Delete(filename)
}

func (lfs *LoggingFileStorage) List() ([]string, error) {
    fmt.Printf("📋 Listando archivos\n")
    return lfs.underlying.List()
}
```
</details>

### Ejercicio 2: Sistema de Autenticación
Crea un sistema flexible de autenticación:

```go
type Authenticator interface {
    Authenticate(username, password string) (User, error)
    IsValid(token string) bool
}

type User struct {
    ID       string
    Username string
    Email    string
    Roles    []string
}

// Implementa:
// 1. DatabaseAuthenticator (verifica contra BD)
// 2. LDAPAuthenticator (verifica contra LDAP)
// 3. MockAuthenticator (para testing)
```

### Ejercicio 3: Sistema de Métricas
Implementa un sistema de métricas desacoplado:

```go
type MetricsCollector interface {
    Counter(name string, value int)
    Gauge(name string, value float64)
    Histogram(name string, value float64)
}

// Implementa:
// 1. PrometheusCollector
// 2. CloudWatchCollector  
// 3. ConsoleCollector (para desarrollo)
// 4. MultiCollector (envía a múltiples destinos)
```

---

## 🎯 Puntos Clave para Recordar

1. **Interfaces definen comportamiento**, no datos
2. **Implementación implícita** - no declaras que implementas
3. **Interfaces pequeñas** son mejores que grandes
4. **Define interfaces donde las uses**, no donde las implementes
5. **Testing se vuelve trivial** con interfaces
6. **Desacoplamiento** permite flexibilidad y mantenimiento
7. **Patrones como Decorator** son naturales con interfaces

---

## 📖 Lecturas Adicionales

- [Go Blog: Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [Effective Go: Interfaces](https://go.dev/doc/effective_go#interfaces)
- [Interface Segregation Principle](https://en.wikipedia.org/wiki/Interface_segregation_principle)

---

## 🚀 Siguiente Paso

¡Excelente! Ahora que dominas interfaces, aprende sobre **Mejores Prácticas Generales** en Go.

Ve a: [`go-mejores-practicas.md`](./go-mejores-practicas.md)