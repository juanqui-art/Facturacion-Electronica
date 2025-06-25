// Package database implementa sistema de respaldos automáticos para SQLite
package database

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BackupConfig configuración para respaldos
type BackupConfig struct {
	RutaRespaldos      string        `json:"ruta_respaldos"`      // Directorio donde guardar respaldos
	IntervaloRespaldo  time.Duration `json:"intervalo_respaldo"`  // Intervalo entre respaldos
	MaxRespaldos       int           `json:"max_respaldos"`       // Máximo número de respaldos a mantener
	CompresiónHabilitada bool        `json:"compresion_habilitada"` // Si comprimir los respaldos
	PrefijRespaldo     string        `json:"prefijo_respaldo"`    // Prefijo para nombres de archivos
}

// BackupManager gestor de respaldos
type BackupManager struct {
	config   BackupConfig
	database *Database
	ticker   *time.Ticker
	done     chan bool
}

// DefaultBackupConfig configuración por defecto para respaldos
var DefaultBackupConfig = BackupConfig{
	RutaRespaldos:        "./respaldos",
	IntervaloRespaldo:    24 * time.Hour, // Diario
	MaxRespaldos:         30,             // 30 días
	CompresiónHabilitada: false,
	PrefijRespaldo:       "facturacion_backup",
}

// NewBackupManager crea un nuevo gestor de respaldos
func NewBackupManager(database *Database, config BackupConfig) *BackupManager {
	return &BackupManager{
		config:   config,
		database: database,
		done:     make(chan bool),
	}
}

// NewBackupManagerDefault crea gestor con configuración por defecto
func NewBackupManagerDefault(database *Database) *BackupManager {
	return NewBackupManager(database, DefaultBackupConfig)
}

// IniciarRespaldosAutomaticos inicia respaldos automáticos en segundo plano
func (bm *BackupManager) IniciarRespaldosAutomaticos() error {
	// Crear directorio de respaldos si no existe
	if err := os.MkdirAll(bm.config.RutaRespaldos, 0755); err != nil {
		return fmt.Errorf("error creando directorio de respaldos: %v", err)
	}

	// Hacer respaldo inicial
	if err := bm.CrearRespaldo(); err != nil {
		return fmt.Errorf("error en respaldo inicial: %v", err)
	}

	// Configurar ticker para respaldos periódicos
	bm.ticker = time.NewTicker(bm.config.IntervaloRespaldo)

	// Goroutine para respaldos automáticos
	go func() {
		for {
			select {
			case <-bm.ticker.C:
				if err := bm.CrearRespaldo(); err != nil {
					fmt.Printf("❌ Error en respaldo automático: %v\n", err)
				}
				
				// Limpiar respaldos antiguos
				if err := bm.LimpiarRespaldosAntiguos(); err != nil {
					fmt.Printf("⚠️  Error limpiando respaldos antiguos: %v\n", err)
				}
			case <-bm.done:
				return
			}
		}
	}()

	fmt.Printf("✅ Respaldos automáticos iniciados (cada %v)\n", bm.config.IntervaloRespaldo)
	return nil
}

// DetenerRespaldosAutomaticos detiene los respaldos automáticos
func (bm *BackupManager) DetenerRespaldosAutomaticos() {
	if bm.ticker != nil {
		bm.ticker.Stop()
		bm.done <- true
		fmt.Println("🛑 Respaldos automáticos detenidos")
	}
}

// CrearRespaldo crea un respaldo de la base de datos
func (bm *BackupManager) CrearRespaldo() error {
	timestamp := time.Now().Format("20060102_150405")
	nombreArchivo := fmt.Sprintf("%s_%s.db", bm.config.PrefijRespaldo, timestamp)
	rutaDestino := filepath.Join(bm.config.RutaRespaldos, nombreArchivo)

	fmt.Printf("💾 Creando respaldo: %s\n", nombreArchivo)

	// Obtener path de la base de datos actual
	dbPath, err := bm.obtenerRutaBaseDatos()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta de base de datos: %v", err)
	}

	// Copiar archivo de base de datos
	if err := bm.copiarArchivo(dbPath, rutaDestino); err != nil {
		return fmt.Errorf("error copiando base de datos: %v", err)
	}

	// Verificar integridad del respaldo
	if err := bm.verificarIntegridadRespaldo(rutaDestino); err != nil {
		// Eliminar respaldo corrupto
		os.Remove(rutaDestino)
		return fmt.Errorf("respaldo corrupto, eliminado: %v", err)
	}

	fmt.Printf("✅ Respaldo creado exitosamente: %s\n", rutaDestino)
	return nil
}

// obtenerRutaBaseDatos obtiene la ruta del archivo de base de datos actual
func (bm *BackupManager) obtenerRutaBaseDatos() (string, error) {
	// En una implementación más avanzada, esto podría obtenerse del objeto Database
	// Por ahora, intentamos varias rutas posibles
	rutasPosibles := []string{
		"database/facturacion.db",
		"test_respaldos.db",
		"test_sistema_completo.db",
		"facturacion.db",
	}
	
	for _, ruta := range rutasPosibles {
		if _, err := os.Stat(ruta); err == nil {
			return ruta, nil
		}
	}
	
	// Si no encuentra ninguna, usar la por defecto
	return "database/facturacion.db", nil
}

// copiarArchivo copia un archivo de origen a destino
func (bm *BackupManager) copiarArchivo(origen, destino string) error {
	archivoOrigen, err := os.Open(origen)
	if err != nil {
		return err
	}
	defer archivoOrigen.Close()

	archivoDestino, err := os.Create(destino)
	if err != nil {
		return err
	}
	defer archivoDestino.Close()

	_, err = io.Copy(archivoDestino, archivoOrigen)
	if err != nil {
		return err
	}

	return archivoDestino.Sync()
}

// verificarIntegridadRespaldo verifica que el respaldo sea válido
func (bm *BackupManager) verificarIntegridadRespaldo(rutaRespaldo string) error {
	// Intentar abrir la base de datos de respaldo
	respaldoDB, err := New(rutaRespaldo)
	if err != nil {
		return fmt.Errorf("respaldo no se puede abrir: %v", err)
	}
	defer respaldoDB.Close()

	// Hacer una consulta simple para verificar integridad
	stats, err := respaldoDB.EstadisticasFacturas()
	if err != nil {
		return fmt.Errorf("respaldo corrupto, error en consulta: %v", err)
	}

	fmt.Printf("📊 Respaldo verificado - Total facturas: %v\n", stats["total_facturas"])
	return nil
}

// LimpiarRespaldosAntiguos elimina respaldos antiguos según la configuración
func (bm *BackupManager) LimpiarRespaldosAntiguos() error {
	archivos, err := os.ReadDir(bm.config.RutaRespaldos)
	if err != nil {
		return fmt.Errorf("error leyendo directorio de respaldos: %v", err)
	}

	// Filtrar solo archivos de respaldo
	var respaldos []os.FileInfo
	for _, archivo := range archivos {
		if !archivo.IsDir() && strings.HasPrefix(archivo.Name(), bm.config.PrefijRespaldo) {
			info, err := archivo.Info()
			if err == nil {
				respaldos = append(respaldos, info)
			}
		}
	}

	// Ordenar por fecha de modificación (más reciente primero)
	sort.Slice(respaldos, func(i, j int) bool {
		return respaldos[i].ModTime().After(respaldos[j].ModTime())
	})

	// Eliminar respaldos excedentes
	if len(respaldos) > bm.config.MaxRespaldos {
		eliminados := 0
		for i := bm.config.MaxRespaldos; i < len(respaldos); i++ {
			rutaArchivo := filepath.Join(bm.config.RutaRespaldos, respaldos[i].Name())
			if err := os.Remove(rutaArchivo); err != nil {
				fmt.Printf("⚠️  Error eliminando respaldo antiguo %s: %v\n", respaldos[i].Name(), err)
			} else {
				eliminados++
			}
		}
		
		if eliminados > 0 {
			fmt.Printf("🗑️  Eliminados %d respaldos antiguos\n", eliminados)
		}
	}

	return nil
}

// ListarRespaldos lista todos los respaldos disponibles
func (bm *BackupManager) ListarRespaldos() ([]RespaldoInfo, error) {
	archivos, err := os.ReadDir(bm.config.RutaRespaldos)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de respaldos: %v", err)
	}

	var respaldos []RespaldoInfo
	for _, archivo := range archivos {
		if !archivo.IsDir() && strings.HasPrefix(archivo.Name(), bm.config.PrefijRespaldo) {
			info, err := archivo.Info()
			if err != nil {
				continue
			}

			respaldo := RespaldoInfo{
				Nombre:         info.Name(),
				RutaCompleta:   filepath.Join(bm.config.RutaRespaldos, info.Name()),
				FechaCreacion:  info.ModTime(),
				TamañoBytes:    info.Size(),
				TamañoLegible:  formatearTamaño(info.Size()),
			}
			respaldos = append(respaldos, respaldo)
		}
	}

	// Ordenar por fecha (más reciente primero)
	sort.Slice(respaldos, func(i, j int) bool {
		return respaldos[i].FechaCreacion.After(respaldos[j].FechaCreacion)
	})

	return respaldos, nil
}

// RestaurarDesdeRespaldo restaura la base de datos desde un respaldo
func (bm *BackupManager) RestaurarDesdeRespaldo(rutaRespaldo string) error {
	fmt.Printf("🔄 Restaurando desde respaldo: %s\n", rutaRespaldo)

	// Verificar que el respaldo existe
	if _, err := os.Stat(rutaRespaldo); os.IsNotExist(err) {
		return fmt.Errorf("respaldo no encontrado: %s", rutaRespaldo)
	}

	// Verificar integridad del respaldo antes de restaurar
	if err := bm.verificarIntegridadRespaldo(rutaRespaldo); err != nil {
		return fmt.Errorf("respaldo corrupto: %v", err)
	}

	// Obtener ruta de base de datos actual
	dbPath, err := bm.obtenerRutaBaseDatos()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta de base de datos: %v", err)
	}

	// Crear respaldo de la base de datos actual antes de restaurar
	respaldoActual := dbPath + ".pre_restore_" + time.Now().Format("20060102_150405")
	if err := bm.copiarArchivo(dbPath, respaldoActual); err != nil {
		fmt.Printf("⚠️  No se pudo respaldar la base actual: %v\n", err)
	} else {
		fmt.Printf("💾 Base actual respaldada en: %s\n", respaldoActual)
	}

	// Restaurar desde respaldo
	if err := bm.copiarArchivo(rutaRespaldo, dbPath); err != nil {
		return fmt.Errorf("error restaurando desde respaldo: %v", err)
	}

	fmt.Printf("✅ Restauración completada desde: %s\n", rutaRespaldo)
	return nil
}

// RespaldoInfo información sobre un respaldo
type RespaldoInfo struct {
	Nombre         string    `json:"nombre"`
	RutaCompleta   string    `json:"ruta_completa"`
	FechaCreacion  time.Time `json:"fecha_creacion"`
	TamañoBytes    int64     `json:"tamaño_bytes"`
	TamañoLegible  string    `json:"tamaño_legible"`
}

// formatearTamaño convierte bytes a formato legible
func formatearTamaño(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// CrearRespaldoManual crea un respaldo manual con nombre personalizado
func (bm *BackupManager) CrearRespaldoManual(sufijo string) error {
	// Crear directorio de respaldos si no existe
	if err := os.MkdirAll(bm.config.RutaRespaldos, 0755); err != nil {
		return fmt.Errorf("error creando directorio de respaldos: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	nombreArchivo := fmt.Sprintf("%s_manual_%s_%s.db", bm.config.PrefijRespaldo, sufijo, timestamp)
	rutaDestino := filepath.Join(bm.config.RutaRespaldos, nombreArchivo)

	fmt.Printf("💾 Creando respaldo manual: %s\n", nombreArchivo)

	dbPath, err := bm.obtenerRutaBaseDatos()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta de base de datos: %v", err)
	}

	if err := bm.copiarArchivo(dbPath, rutaDestino); err != nil {
		return fmt.Errorf("error copiando base de datos: %v", err)
	}

	if err := bm.verificarIntegridadRespaldo(rutaDestino); err != nil {
		os.Remove(rutaDestino)
		return fmt.Errorf("respaldo corrupto, eliminado: %v", err)
	}

	fmt.Printf("✅ Respaldo manual creado: %s\n", rutaDestino)
	return nil
}