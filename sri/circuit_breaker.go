// Package sri implementa circuit breaker para comunicación robusta con SRI
package sri

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// EstadoCircuitBreaker representa el estado del circuit breaker
type EstadoCircuitBreaker int

const (
	EstadoCerrado EstadoCircuitBreaker = iota // Funcionando normalmente
	EstadoAbierto                             // Bloqueando peticiones por fallos
	EstadoMedioCerrado                        // Permitiendo peticiones de prueba
)

// String implementa Stringer para EstadoCircuitBreaker
func (e EstadoCircuitBreaker) String() string {
	switch e {
	case EstadoCerrado:
		return "CERRADO"
	case EstadoAbierto:
		return "ABIERTO"
	case EstadoMedioCerrado:
		return "MEDIO_CERRADO"
	default:
		return "DESCONOCIDO"
	}
}

// ConfigCircuitBreaker configuración del circuit breaker
type ConfigCircuitBreaker struct {
	MaxErrores        int           `json:"max_errores"`         // Errores antes de abrir el circuito
	TiempoAbierto     time.Duration `json:"tiempo_abierto"`      // Tiempo que permanece abierto
	TiempoEvaluacion  time.Duration `json:"tiempo_evaluacion"`   // Ventana de tiempo para evaluar errores
	MaxPeticionesTest int           `json:"max_peticiones_test"` // Peticiones de prueba en estado medio cerrado
}

// ConfigCircuitBreakerDefault configuración por defecto
var ConfigCircuitBreakerDefault = ConfigCircuitBreaker{
	MaxErrores:        5,
	TiempoAbierto:     30 * time.Second,
	TiempoEvaluacion:  60 * time.Second,
	MaxPeticionesTest: 3,
}

// ConfigCircuitBreakerConservador configuración conservadora para producción
var ConfigCircuitBreakerConservador = ConfigCircuitBreaker{
	MaxErrores:        3,
	TiempoAbierto:     60 * time.Second,
	TiempoEvaluacion:  120 * time.Second,
	MaxPeticionesTest: 2,
}

// CircuitBreaker implementa el patrón circuit breaker para SRI
type CircuitBreaker struct {
	config              ConfigCircuitBreaker
	estado              EstadoCircuitBreaker
	errores             int
	ultimoError         time.Time
	ultimoCambioEstado  time.Time
	peticionesTest      int
	mutex               sync.RWMutex
	estadisticas        EstadisticasCircuitBreaker
}

// EstadisticasCircuitBreaker estadísticas del circuit breaker
type EstadisticasCircuitBreaker struct {
	TotalPeticiones      int64     `json:"total_peticiones"`
	PeticionesExitosas   int64     `json:"peticiones_exitosas"`
	PeticionesFallidas   int64     `json:"peticiones_fallidas"`
	PeticionesBloqueadas int64     `json:"peticiones_bloqueadas"`
	VecesAbierto         int64     `json:"veces_abierto"`
	UltimaApertura       time.Time `json:"ultima_apertura"`
	UltimoCierre         time.Time `json:"ultimo_cierre"`
}

// NuevoCircuitBreaker crea una nueva instancia de circuit breaker
func NuevoCircuitBreaker(config ConfigCircuitBreaker) *CircuitBreaker {
	return &CircuitBreaker{
		config:             config,
		estado:             EstadoCerrado,
		ultimoCambioEstado: time.Now(),
		estadisticas:       EstadisticasCircuitBreaker{},
	}
}

// NuevoCircuitBreakerDefault crea circuit breaker con configuración por defecto
func NuevoCircuitBreakerDefault() *CircuitBreaker {
	return NuevoCircuitBreaker(ConfigCircuitBreakerDefault)
}

// Ejecutar ejecuta una función con protección de circuit breaker
func (cb *CircuitBreaker) Ejecutar(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.estadisticas.TotalPeticiones++

	// Verificar si podemos ejecutar la petición
	if !cb.puedeEjecutar() {
		cb.estadisticas.PeticionesBloqueadas++
		return fmt.Errorf("circuit breaker ABIERTO: SRI no disponible, reintente en %v", 
			cb.tiempoRestanteAbierto())
	}

	// Ejecutar la función
	err := fn()

	// Actualizar estado según resultado
	if err != nil {
		cb.registrarError(err)
		cb.estadisticas.PeticionesFallidas++
		return err
	}

	cb.registrarExito()
	cb.estadisticas.PeticionesExitosas++
	return nil
}

// puedeEjecutar determina si se puede ejecutar una petición
func (cb *CircuitBreaker) puedeEjecutar() bool {
	switch cb.estado {
	case EstadoCerrado:
		return true
	case EstadoAbierto:
		// Verificar si es tiempo de pasar a medio cerrado
		if time.Since(cb.ultimoCambioEstado) >= cb.config.TiempoAbierto {
			cb.cambiarEstado(EstadoMedioCerrado)
			cb.peticionesTest = 0
			log.Printf("[CIRCUIT_BREAKER] Cambiando a estado MEDIO_CERRADO para pruebas")
			return true
		}
		return false
	case EstadoMedioCerrado:
		return cb.peticionesTest < cb.config.MaxPeticionesTest
	default:
		return false
	}
}

// registrarError registra un error y actualiza el estado
func (cb *CircuitBreaker) registrarError(err error) {
	cb.ultimoError = time.Now()
	
	switch cb.estado {
	case EstadoCerrado:
		cb.errores++
		// Verificar si debemos abrir el circuito
		if cb.errores >= cb.config.MaxErrores {
			cb.cambiarEstado(EstadoAbierto)
			cb.estadisticas.VecesAbierto++
			cb.estadisticas.UltimaApertura = time.Now()
			log.Printf("[CIRCUIT_BREAKER] ABRIENDO circuito después de %d errores. Último error: %v", 
				cb.errores, err)
		}
	case EstadoMedioCerrado:
		// En estado medio cerrado, cualquier error nos regresa a abierto
		cb.cambiarEstado(EstadoAbierto)
		cb.estadisticas.VecesAbierto++
		cb.estadisticas.UltimaApertura = time.Now()
		log.Printf("[CIRCUIT_BREAKER] Regresando a estado ABIERTO desde medio cerrado. Error: %v", err)
	}
}

// registrarExito registra un éxito y actualiza el estado
func (cb *CircuitBreaker) registrarExito() {
	switch cb.estado {
	case EstadoCerrado:
		// Limpiar errores en ventana de tiempo
		if time.Since(cb.ultimoError) >= cb.config.TiempoEvaluacion {
			cb.errores = 0
		}
	case EstadoMedioCerrado:
		cb.peticionesTest++
		// Si completamos todas las peticiones de prueba exitosamente, cerrar circuito
		if cb.peticionesTest >= cb.config.MaxPeticionesTest {
			cb.cambiarEstado(EstadoCerrado)
			cb.errores = 0
			cb.estadisticas.UltimoCierre = time.Now()
			log.Printf("[CIRCUIT_BREAKER] CERRANDO circuito después de %d peticiones exitosas", 
				cb.config.MaxPeticionesTest)
		}
	}
}

// cambiarEstado cambia el estado del circuit breaker
func (cb *CircuitBreaker) cambiarEstado(nuevoEstado EstadoCircuitBreaker) {
	estadoAnterior := cb.estado
	cb.estado = nuevoEstado
	cb.ultimoCambioEstado = time.Now()
	
	log.Printf("[CIRCUIT_BREAKER] Cambio de estado: %s -> %s", 
		estadoAnterior, nuevoEstado)
}

// tiempoRestanteAbierto calcula el tiempo restante en estado abierto
func (cb *CircuitBreaker) tiempoRestanteAbierto() time.Duration {
	if cb.estado != EstadoAbierto {
		return 0
	}
	
	transcurrido := time.Since(cb.ultimoCambioEstado)
	restante := cb.config.TiempoAbierto - transcurrido
	
	if restante < 0 {
		return 0
	}
	
	return restante
}

// ObtenerEstado obtiene el estado actual del circuit breaker (thread-safe)
func (cb *CircuitBreaker) ObtenerEstado() EstadoCircuitBreaker {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.estado
}

// ObtenerEstadisticas obtiene las estadísticas del circuit breaker (thread-safe)
func (cb *CircuitBreaker) ObtenerEstadisticas() EstadisticasCircuitBreaker {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.estadisticas
}

// Reiniciar reinicia el circuit breaker al estado inicial
func (cb *CircuitBreaker) Reiniciar() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.estado = EstadoCerrado
	cb.errores = 0
	cb.peticionesTest = 0
	cb.ultimoCambioEstado = time.Now()
	
	log.Printf("[CIRCUIT_BREAKER] Circuit breaker reiniciado")
}

// MostrarEstado muestra información detallada del estado actual
func (cb *CircuitBreaker) MostrarEstado() {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	fmt.Printf("\n🔧 ESTADO CIRCUIT BREAKER\n")
	fmt.Printf("========================\n")
	fmt.Printf("📊 Estado Actual: %s\n", cb.estado)
	fmt.Printf("❌ Errores Actuales: %d/%d\n", cb.errores, cb.config.MaxErrores)
	fmt.Printf("🧪 Peticiones Test: %d/%d\n", cb.peticionesTest, cb.config.MaxPeticionesTest)
	fmt.Printf("⏱️  Último Cambio: %v\n", cb.ultimoCambioEstado.Format("15:04:05"))
	
	if cb.estado == EstadoAbierto {
		fmt.Printf("⏳ Tiempo Restante Abierto: %v\n", cb.tiempoRestanteAbierto())
	}
	
	fmt.Printf("\n📈 ESTADÍSTICAS GENERALES\n")
	fmt.Printf("========================\n")
	fmt.Printf("📊 Total Peticiones: %d\n", cb.estadisticas.TotalPeticiones)
	fmt.Printf("✅ Exitosas: %d\n", cb.estadisticas.PeticionesExitosas)
	fmt.Printf("❌ Fallidas: %d\n", cb.estadisticas.PeticionesFallidas)
	fmt.Printf("🚫 Bloqueadas: %d\n", cb.estadisticas.PeticionesBloqueadas)
	fmt.Printf("🔓 Veces Abierto: %d\n", cb.estadisticas.VecesAbierto)
	
	if !cb.estadisticas.UltimaApertura.IsZero() {
		fmt.Printf("🕐 Última Apertura: %v\n", cb.estadisticas.UltimaApertura.Format("15:04:05"))
	}
	if !cb.estadisticas.UltimoCierre.IsZero() {
		fmt.Printf("🕐 Último Cierre: %v\n", cb.estadisticas.UltimoCierre.Format("15:04:05"))
	}
	
	fmt.Printf("========================\n")
}

// EsOperacional indica si el circuit breaker permite operaciones
func (cb *CircuitBreaker) EsOperacional() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.puedeEjecutar()
}