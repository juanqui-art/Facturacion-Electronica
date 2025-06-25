// Package sri implementa sistema de logging básico para debugging
package sri

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel representa el nivel de logging
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelCritical
)

// String implementa Stringer para LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Logger estructura para logging básico
type Logger struct {
	nivel       LogLevel
	archivo     *os.File
	logger      *log.Logger
	habilitado  bool
}

// logger global
var logger *Logger

// init inicializa el logger por defecto
func init() {
	logger = &Logger{
		nivel:      LogLevelInfo,
		logger:     log.New(os.Stdout, "", log.LstdFlags),
		habilitado: true,
	}
}

// ConfigurarLogger configura el logger global
func ConfigurarLogger(nivel LogLevel, archivoLog string) error {
	logger.nivel = nivel
	
	if archivoLog != "" {
		// Crear directorio si no existe
		dir := filepath.Dir(archivoLog)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creando directorio de logs: %v", err)
		}
		
		// Abrir archivo de log
		file, err := os.OpenFile(archivoLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("error abriendo archivo de log: %v", err)
		}
		
		// Cerrar archivo anterior si existe
		if logger.archivo != nil {
			logger.archivo.Close()
		}
		
		logger.archivo = file
		logger.logger = log.New(file, "", log.LstdFlags)
	}
	
	return nil
}

// HabilitarLogger habilita o deshabilita el logging
func HabilitarLogger(habilitado bool) {
	logger.habilitado = habilitado
}

// escribir escribe un mensaje de log si el nivel es apropiado
func escribir(nivel LogLevel, formato string, args ...interface{}) {
	if !logger.habilitado || nivel < logger.nivel {
		return
	}
	
	mensaje := fmt.Sprintf(formato, args...)
	prefijo := fmt.Sprintf("[%s] ", nivel.String())
	logger.logger.Printf("%s%s", prefijo, mensaje)
}

// Debug escribe un mensaje de debug
func Debug(formato string, args ...interface{}) {
	escribir(LogLevelDebug, formato, args...)
}

// Info escribe un mensaje informativo
func Info(formato string, args ...interface{}) {
	escribir(LogLevelInfo, formato, args...)
}

// Warning escribe un mensaje de advertencia
func Warning(formato string, args ...interface{}) {
	escribir(LogLevelWarning, formato, args...)
}

// Error escribe un mensaje de error
func Error(formato string, args ...interface{}) {
	escribir(LogLevelError, formato, args...)
}

// Critical escribe un mensaje crítico
func Critical(formato string, args ...interface{}) {
	escribir(LogLevelCritical, formato, args...)
}

// LogValidacion registra eventos de validación
func LogValidacion(operacion string, exito bool, detalle string) {
	if exito {
		Info("Validación exitosa - %s: %s", operacion, detalle)
	} else {
		Warning("Validación fallida - %s: %s", operacion, detalle)
	}
}

// LogSRI registra eventos de comunicación con SRI
func LogSRI(operacion string, exito bool, tiempoMs int64, detalle string) {
	if exito {
		Info("SRI OK - %s (%dms): %s", operacion, tiempoMs, detalle)
	} else {
		Error("SRI ERROR - %s (%dms): %s", operacion, tiempoMs, detalle)
	}
}

// LogCircuitBreaker registra eventos del circuit breaker
func LogCircuitBreaker(evento string, estado string, detalle string) {
	Info("Circuit Breaker - %s [%s]: %s", evento, estado, detalle)
}

// LogReintento registra eventos de reintentos
func LogReintento(operacion string, intento int, maxIntentos int, exito bool, detalle string) {
	if exito {
		Info("Reintento exitoso - %s (%d/%d): %s", operacion, intento, maxIntentos, detalle)
	} else {
		Warning("Reintento fallido - %s (%d/%d): %s", operacion, intento, maxIntentos, detalle)
	}
}

// LogFactura registra eventos de procesamiento de facturas
func LogFactura(claveAcceso string, operacion string, exito bool, detalle string) {
	if exito {
		Info("Factura %s - %s: %s", claveAcceso, operacion, detalle)
	} else {
		Error("Factura %s ERROR - %s: %s", claveAcceso, operacion, detalle)
	}
}

// LogSeguridad registra eventos de seguridad (inputs maliciosos, etc.)
func LogSeguridad(evento string, detalle string, origen string) {
	Critical("SEGURIDAD - %s desde %s: %s", evento, origen, detalle)
}

// LogPerformance registra métricas de performance
func LogPerformance(operacion string, tiempoMs int64, memoria string) {
	Debug("Performance - %s: %dms, memoria: %s", operacion, tiempoMs, memoria)
}

// CerrarLogger cierra el logger y libera recursos
func CerrarLogger() {
	if logger.archivo != nil {
		logger.archivo.Close()
	}
}

// MostrarEstadisticasLogging muestra estadísticas básicas de logging
func MostrarEstadisticasLogging() {
	fmt.Printf("\n📊 ESTADO DEL LOGGING\n")
	fmt.Printf("====================\n")
	fmt.Printf("🎚️  Nivel: %s\n", logger.nivel)
	fmt.Printf("✅ Habilitado: %t\n", logger.habilitado)
	
	if logger.archivo != nil {
		if stat, err := logger.archivo.Stat(); err == nil {
			fmt.Printf("📁 Archivo: %s\n", logger.archivo.Name())
			fmt.Printf("📏 Tamaño: %d bytes\n", stat.Size())
			fmt.Printf("📅 Modificado: %s\n", stat.ModTime().Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("📁 Salida: Consola\n")
	}
	
	fmt.Printf("====================\n")
}

// InicializarLoggingProduccion configura logging para producción
func InicializarLoggingProduccion() error {
	timestamp := time.Now().Format("2006-01-02")
	archivoLog := fmt.Sprintf("logs/sri-facturacion-%s.log", timestamp)
	
	return ConfigurarLogger(LogLevelInfo, archivoLog)
}

// InicializarLoggingDesarrollo configura logging para desarrollo
func InicializarLoggingDesarrollo() error {
	return ConfigurarLogger(LogLevelDebug, "")
}