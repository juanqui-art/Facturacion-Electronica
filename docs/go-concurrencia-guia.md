# Guía de Concurrencia en Go - Aprendiendo con Ejemplos Reales

Esta guía te explica paso a paso los conceptos de concurrencia en Go usando ejemplos del proyecto de facturación SRI.

## 📚 Tabla de Contenido

1. [¿Qué es la Concurrencia?](#qué-es-la-concurrencia)
2. [Problemas sin Concurrencia](#problemas-sin-concurrencia)
3. [Goroutines Básicas](#goroutines-básicas)
4. [Race Conditions](#race-conditions)
5. [Mutex y RWMutex](#mutex-y-rwmutex)
6. [Implementación Real en el Proyecto](#implementación-real-en-el-proyecto)
7. [Ejercicios Prácticos](#ejercicios-prácticos)

---

## ¿Qué es la Concurrencia?

**Concurrencia** es la capacidad de un programa para manejar múltiples tareas al mismo tiempo. En Go, esto se logra principalmente con **goroutines**.

### Analogía Simple
Imagina una cafetería:
- **Sin concurrencia**: Un empleado atiende un cliente, prepara el café, lo entrega, y solo entonces atiende al siguiente.
- **Con concurrencia**: Un empleado toma pedidos, otro prepara cafés, otro entrega. Varios clientes son atendidos simultáneamente.

### Ejemplo Básico
```go
package main

import (
    "fmt"
    "time"
)

// Función normal (secuencial)
func hacerCafe(cliente string) {
    fmt.Printf("Preparando café para %s...\n", cliente)
    time.Sleep(2 * time.Second) // Simula tiempo de preparación
    fmt.Printf("☕ Café listo para %s\n", cliente)
}

func main() {
    fmt.Println("=== SIN CONCURRENCIA ===")
    start := time.Now()
    
    hacerCafe("Ana")
    hacerCafe("Juan")
    hacerCafe("María")
    
    fmt.Printf("Tiempo total: %v\n", time.Since(start))
    // Resultado: ~6 segundos (2+2+2)
}
```

### Con Concurrencia (Goroutines)
```go
func main() {
    fmt.Println("=== CON CONCURRENCIA ===")
    start := time.Now()
    
    // Lanzar goroutines
    go hacerCafe("Ana")
    go hacerCafe("Juan") 
    go hacerCafe("María")
    
    // Esperar a que terminen todas
    time.Sleep(3 * time.Second)
    
    fmt.Printf("Tiempo total: %v\n", time.Since(start))
    // Resultado: ~2 segundos (todas en paralelo)
}
```

---

## Problemas sin Concurrencia

En nuestro proyecto de facturación, **antes** teníamos este código problemático:

```go
// ❌ CÓDIGO PROBLEMÁTICO (ANTES)
var facturaStorage = make(map[string]FacturaResponse)
var nextID = 1

func crearFactura() {
    // Problema: ¿Qué pasa si dos usuarios llegan al mismo tiempo?
    id := fmt.Sprintf("FAC-%06d", nextID)  // Usuario A: FAC-000001
    nextID++                               // nextID = 2
    
    // Pero entre estas líneas, otro usuario puede hacer lo mismo:
    // Usuario B también obtiene nextID = 1, crea FAC-000001
    // ¡ID DUPLICADO!
    
    facturaStorage[id] = factura // Usuario B sobrescribe factura de A
}
```

### Simulación del Problema
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var contador = 0

// Función INSEGURA
func incrementarInseguro(usuario string) {
    for i := 0; i < 1000; i++ {
        temp := contador     // Lee valor actual
        time.Sleep(1 * time.Nanosecond) // Simula trabajo
        contador = temp + 1  // Escribe nuevo valor
    }
    fmt.Printf("%s terminó. Contador: %d\n", usuario, contador)
}

func main() {
    fmt.Println("=== SIMULACIÓN DE RACE CONDITION ===")
    
    var wg sync.WaitGroup
    wg.Add(3)
    
    go func() {
        defer wg.Done()
        incrementarInseguro("Usuario A")
    }()
    
    go func() {
        defer wg.Done()
        incrementarInseguro("Usuario B")
    }()
    
    go func() {
        defer wg.Done()
        incrementarInseguro("Usuario C")
    }()
    
    wg.Wait()
    fmt.Printf("Contador final: %d (debería ser 3000)\n", contador)
    // Resultado: Número menor a 3000 debido a race conditions
}
```

---

## Race Conditions

Una **race condition** ocurre cuando múltiples goroutines acceden al mismo recurso simultáneamente, y el resultado depende del orden de ejecución.

### Ejemplo Visual
```
Tiempo →  T1    T2    T3    T4    T5
Goroutine A: Lee(5) →    → Escribe(6)
Goroutine B:    → Lee(5) → Escribe(6)

Resultado: 6 (perdimos un incremento)
Esperado: 7
```

### Cómo Detectar Race Conditions
Go incluye un detector de race conditions:

```bash
# Correr con detector
go run -race main.go

# Ejemplo de salida:
# WARNING: DATA RACE
# Write at 0x... by goroutine 7:
#   main.incrementar()
# Previous read at 0x... by goroutine 6:
#   main.incrementar()
```

---

## Mutex y RWMutex

**Mutex** (Mutual Exclusion) es como un candado que solo permite a una goroutine acceder al recurso a la vez.

### Tipos de Mutex

#### 1. Mutex Básico
```go
import "sync"

var (
    contador int
    mu       sync.Mutex
)

func incrementarSeguro() {
    mu.Lock()         // 🔒 Pedir el candado
    defer mu.Unlock() // 🔓 Liberar automáticamente al final
    
    contador++        // Solo una goroutine puede estar aquí
}
```

#### 2. RWMutex (Read-Write Mutex)
Permite múltiples lectores, pero solo un escritor:

```go
var (
    datos   map[string]string
    rwMutex sync.RWMutex
)

// Múltiples goroutines pueden leer simultáneamente
func leerDatos(key string) string {
    rwMutex.RLock()         // 🔍 Candado de lectura
    defer rwMutex.RUnlock() // 🔍 Liberar lectura
    
    return datos[key]
}

// Solo una goroutine puede escribir
func escribirDatos(key, value string) {
    rwMutex.Lock()         // 🔒 Candado exclusivo
    defer rwMutex.Unlock() // 🔓 Liberar escritura
    
    datos[key] = value
}
```

### Analogía del Mutex
Imagina una biblioteca:
- **Mutex**: Solo una persona puede usar el libro a la vez
- **RWMutex**: Muchas personas pueden leer el mismo libro, pero solo una puede escribir en él

---

## Implementación Real en el Proyecto

Así es como solucionamos el problema en nuestro proyecto:

### Antes (❌ Inseguro)
```go
// Variables globales sin protección
var facturaStorage = make(map[string]FacturaResponse)
var nextID = 1

func crearFactura() {
    // Race condition aquí ↓
    id := fmt.Sprintf("FAC-%06d", nextID)
    nextID++
    facturaStorage[id] = factura
}
```

### Después (✅ Seguro)
```go
// Estructura con protección
type FacturaStorage struct {
    mu       sync.RWMutex                    // 🔒 Protección
    facturas map[string]FacturaResponse     // Datos protegidos
    nextID   int                            // ID protegido
}

func (fs *FacturaStorage) Store(id string, factura FacturaResponse) {
    fs.mu.Lock()         // 🔒 Solo yo puedo escribir
    defer fs.mu.Unlock() // 🔓 Liberar al terminar
    
    fs.facturas[id] = factura
}

func (fs *FacturaStorage) Get(id string) (FacturaResponse, bool) {
    fs.mu.RLock()         // 🔍 Muchos pueden leer
    defer fs.mu.RUnlock() // 🔍 Liberar lectura
    
    factura, exists := fs.facturas[id]
    return factura, exists
}

func (fs *FacturaStorage) GetNextID() int {
    fs.mu.Lock()         // 🔒 ID único garantizado
    defer fs.mu.Unlock()
    
    id := fs.nextID
    fs.nextID++
    return id
}
```

### Uso en la API
```go
// Storage global thread-safe
var storage = NewFacturaStorage()

func handleCreateFactura(w http.ResponseWriter, r *http.Request) {
    // ... validar datos ...
    
    // ID único garantizado, sin race conditions
    nextID := storage.GetNextID()
    id := fmt.Sprintf("FAC-%06d", nextID)
    
    response := FacturaResponse{
        ID:      id,
        Factura: factura,
        // ...
    }
    
    // Almacenar de forma thread-safe
    storage.Store(id, response)
    
    writeJSONResponse(w, http.StatusCreated, response)
}
```

---

## Ejercicios Prácticos

### Ejercicio 1: Contador Seguro
Implementa un contador thread-safe:

```go
type ContadorSeguro struct {
    // ¿Qué campos necesitas?
}

func (c *ContadorSeguro) Incrementar() {
    // Tu implementación aquí
}

func (c *ContadorSeguro) Obtener() int {
    // Tu implementación aquí
}

func (c *ContadorSeguro) Reset() {
    // Tu implementación aquí
}
```

<details>
<summary>👁️ Ver Solución</summary>

```go
type ContadorSeguro struct {
    mu    sync.RWMutex
    valor int
}

func (c *ContadorSeguro) Incrementar() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor++
}

func (c *ContadorSeguro) Obtener() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.valor
}

func (c *ContadorSeguro) Reset() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.valor = 0
}
```
</details>

### Ejercicio 2: Cache Thread-Safe
Implementa un cache simple:

```go
type Cache struct {
    // Tu implementación
}

func (c *Cache) Set(key, value string) {
    // Tu implementación
}

func (c *Cache) Get(key string) (string, bool) {
    // Tu implementación
}

func (c *Cache) Delete(key string) {
    // Tu implementación
}
```

### Ejercicio 3: Banco Thread-Safe
Simula transferencias bancarias sin race conditions:

```go
type Cuenta struct {
    // Tu implementación
}

func (c *Cuenta) Depositar(cantidad float64) {
    // Tu implementación
}

func (c *Cuenta) Retirar(cantidad float64) error {
    // Tu implementación
}

func (c *Cuenta) Saldo() float64 {
    // Tu implementación
}

func Transferir(origen, destino *Cuenta, cantidad float64) error {
    // ¿Cómo evitas deadlocks?
}
```

---

## 🎯 Puntos Clave para Recordar

1. **Goroutines** son ligeras y baratas, puedes crear miles
2. **Race conditions** son bugs difíciles de reproducir - usa `-race` para detectarlas
3. **Mutex** es para exclusión mutua (solo uno a la vez)
4. **RWMutex** permite múltiples lectores, un escritor
5. **defer unlock** siempre para evitar deadlocks
6. **Keep it simple**: No optimices prematuramente

---

## 📖 Lecturas Adicionales

- [Go Blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Race Detector](https://go.dev/doc/articles/race_detector)

---

## 🚀 Siguiente Paso

¡Ahora que entiendes concurrencia, el siguiente paso es aprender sobre **Interfaces y Desacoplamiento**! 

Ve a: [`go-interfaces-guia.md`](./go-interfaces-guia.md)