// Package sri implementa lógica de reintentos avanzada para comunicación con SRI
package sri

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ConfigReintento configuración para lógica de reintentos
type ConfigReintento struct {
	MaxIntentos     int           `json:"max_intentos"`
	TiempoBase      time.Duration `json:"tiempo_base"`
	Multiplicador   float64       `json:"multiplicador"`
	TiempoMaximo    time.Duration `json:"tiempo_maximo"`
	JitterMaximo    time.Duration `json:"jitter_maximo"`
	SoloRecuperables bool         `json:"solo_recuperables"`
}

// ConfigReintentoDefault configuración por defecto para reintentos
var ConfigReintentoDefault = ConfigReintento{
	MaxIntentos:     5,
	TiempoBase:      2 * time.Second,
	Multiplicador:   2.0,
	TiempoMaximo:    30 * time.Second,
	JitterMaximo:    1 * time.Second,
	SoloRecuperables: true,
}

// ConfigReintentoConservador configuración conservadora para producción
var ConfigReintentoConservador = ConfigReintento{
	MaxIntentos:     3,
	TiempoBase:      5 * time.Second,
	Multiplicador:   2.0,
	TiempoMaximo:    60 * time.Second,
	JitterMaximo:    2 * time.Second,
	SoloRecuperables: true,
}

// ConfigReintentoAgresivo configuración agresiva para desarrollo
var ConfigReintentoAgresivo = ConfigReintento{
	MaxIntentos:     7,
	TiempoBase:      1 * time.Second,
	Multiplicador:   1.5,
	TiempoMaximo:    20 * time.Second,
	JitterMaximo:    500 * time.Millisecond,
	SoloRecuperables: false,
}

// ResultadoReintento resultado de la ejecución con reintentos
type ResultadoReintento struct {
	Exitoso         bool          `json:"exitoso"`
	IntentosRealizados int        `json:"intentos_realizados"`
	TiempoTotal     time.Duration `json:"tiempo_total"`
	UltimoError     error         `json:"ultimo_error"`
	Errores         []error       `json:"errores"`
}

// FuncionReintentable función que puede ser reintentada
type FuncionReintentable func() error

// EjecutarConReintento ejecuta una función con lógica de reintentos
func EjecutarConReintento(fn FuncionReintentable, config ConfigReintento) *ResultadoReintento {
	inicio := time.Now()
	resultado := &ResultadoReintento{
		Exitoso:            false,
		IntentosRealizados: 0,
		Errores:            make([]error, 0),
	}

	for intento := 1; intento <= config.MaxIntentos; intento++ {
		resultado.IntentosRealizados = intento
		
		fmt.Printf("🔄 Intento %d/%d...\n", intento, config.MaxIntentos)
		
		err := fn()
		if err == nil {
			// Éxito!
			resultado.Exitoso = true
			resultado.TiempoTotal = time.Since(inicio)
			fmt.Printf("✅ Operación exitosa en intento %d\n", intento)
			return resultado
		}

		// Registrar error
		resultado.UltimoError = err
		resultado.Errores = append(resultado.Errores, err)
		
		// Mostrar error
		fmt.Printf("❌ Error en intento %d: %v\n", intento, err)
		
		// Verificar si el error es recuperable
		if config.SoloRecuperables && !EsErrorRecuperable(err) {
			fmt.Printf("🚫 Error no recuperable, abortando reintentos\n")
			break
		}

		// Si no es el último intento, esperar antes del siguiente
		if intento < config.MaxIntentos {
			tiempoEspera := calcularTiempoEspera(intento, config)
			fmt.Printf("⏳ Esperando %v antes del siguiente intento...\n", tiempoEspera)
			time.Sleep(tiempoEspera)
		}
	}

	resultado.TiempoTotal = time.Since(inicio)
	fmt.Printf("💥 Operación falló después de %d intentos\n", resultado.IntentosRealizados)
	
	return resultado
}

// calcularTiempoEspera calcula el tiempo de espera con backoff exponencial
func calcularTiempoEspera(intento int, config ConfigReintento) time.Duration {
	// Backoff exponencial: tiempo_base * multiplicador^(intento-1)
	tiempoBase := float64(config.TiempoBase)
	factor := math.Pow(config.Multiplicador, float64(intento-1))
	tiempoCalculado := time.Duration(tiempoBase * factor)

	// Aplicar límite máximo
	if tiempoCalculado > config.TiempoMaximo {
		tiempoCalculado = config.TiempoMaximo
	}

	// Agregar jitter aleatorio para evitar thundering herd
	jitter := time.Duration(float64(config.JitterMaximo) * (0.5 + 0.5*float64(time.Now().UnixNano()%2)))
	
	return tiempoCalculado + jitter
}

// ReintentarEnvioSRI envía comprobante al SRI con reintentos
func (c *SOAPClient) ReintentarEnvioSRI(xmlComprobante []byte, config ConfigReintento) (*RespuestaSolicitud, *ResultadoReintento) {
	var respuesta *RespuestaSolicitud
	var resultadoReintento *ResultadoReintento

	fn := func() error {
		resp, err := c.EnviarComprobante(xmlComprobante)
		if err != nil {
			return err
		}
		respuesta = resp
		return nil
	}

	resultadoReintento = EjecutarConReintento(fn, config)
	
	return respuesta, resultadoReintento
}

// ReintentarConsultaAutorizacion consulta autorización con reintentos
func (c *SOAPClient) ReintentarConsultaAutorizacion(claveAcceso string, config ConfigReintento) (*RespuestaComprobante, *ResultadoReintento) {
	var respuesta *RespuestaComprobante
	var resultadoReintento *ResultadoReintento

	fn := func() error {
		resp, err := c.ConsultarAutorizacion(claveAcceso)
		if err != nil {
			return err
		}
		
		// Verificar si el comprobante ya está autorizado
		if len(resp.Autorizaciones) > 0 {
			autorizacion := resp.Autorizaciones[0]
			if autorizacion.Estado == "AUTORIZADO" {
				respuesta = resp
				return nil
			} else if autorizacion.Estado == "NO_AUTORIZADO" {
				return CrearErrorValidacion("ESTADO", "Comprobante no autorizado por SRI")
			}
		}
		
		// Si aún está en proceso, considerarlo como error temporal
		return CrearErrorConexion("Comprobante aún en procesamiento")
	}

	resultadoReintento = EjecutarConReintento(fn, config)
	
	return respuesta, resultadoReintento
}

// ProcesarComprobanteCompletoConReintento procesa comprobante completo con reintentos avanzados
func (c *SOAPClient) ProcesarComprobanteCompletoConReintento(xmlComprobante []byte, claveAcceso string) (*AutorizacionSRI, error) {
	fmt.Println("🚀 Iniciando procesamiento completo con reintentos avanzados...")
	
	// Paso 1: Enviar comprobante con reintentos
	fmt.Println("\n📤 PASO 1: Enviando comprobante al SRI")
	fmt.Println(strings.Repeat("-", 50))
	
	respRecepcion, resultadoEnvio := c.ReintentarEnvioSRI(xmlComprobante, ConfigReintentoDefault)
	if !resultadoEnvio.Exitoso {
		return nil, fmt.Errorf("error enviando comprobante después de %d intentos: %v", 
			resultadoEnvio.IntentosRealizados, resultadoEnvio.UltimoError)
	}
	
	fmt.Printf("✅ Comprobante enviado exitosamente en %d intentos (tiempo: %v)\n", 
		resultadoEnvio.IntentosRealizados, resultadoEnvio.TiempoTotal)
	
	// Verificar estado de recepción
	if respRecepcion.Estado != "RECIBIDA" {
		return nil, fmt.Errorf("comprobante no fue recibido por SRI. Estado: %s", respRecepcion.Estado)
	}

	// Paso 2: Esperar procesamiento inicial
	fmt.Println("\n⏳ PASO 2: Esperando procesamiento inicial del SRI")
	fmt.Println(strings.Repeat("-", 50))
	tiempoEsperaInicial := 5 * time.Second
	fmt.Printf("Esperando %v para permitir procesamiento...\n", tiempoEsperaInicial)
	time.Sleep(tiempoEsperaInicial)

	// Paso 3: Consultar autorización con reintentos
	fmt.Println("\n🔍 PASO 3: Consultando autorización")
	fmt.Println(strings.Repeat("-", 50))
	
	// Configuración especial para consulta de autorización
	configConsulta := ConfigReintento{
		MaxIntentos:      8,  // Más intentos para autorización
		TiempoBase:       3 * time.Second,
		Multiplicador:    1.5,
		TiempoMaximo:     45 * time.Second,
		JitterMaximo:     1 * time.Second,
		SoloRecuperables: true,
	}
	
	respAutorizacion, resultadoConsulta := c.ReintentarConsultaAutorizacion(claveAcceso, configConsulta)
	if !resultadoConsulta.Exitoso {
		return nil, fmt.Errorf("error consultando autorización después de %d intentos: %v", 
			resultadoConsulta.IntentosRealizados, resultadoConsulta.UltimoError)
	}
	
	fmt.Printf("✅ Autorización obtenida exitosamente en %d intentos (tiempo: %v)\n", 
		resultadoConsulta.IntentosRealizados, resultadoConsulta.TiempoTotal)

	// Paso 4: Verificar resultado final
	if len(respAutorizacion.Autorizaciones) == 0 {
		return nil, fmt.Errorf("no se encontraron autorizaciones para la clave de acceso")
	}

	autorizacion := respAutorizacion.Autorizaciones[0]
	
	fmt.Println("\n🎉 RESULTADO FINAL")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("📊 Estado: %s\n", autorizacion.Estado)
	fmt.Printf("📝 Número de Autorización: %s\n", autorizacion.NumeroAutorizacion)
	fmt.Printf("📅 Fecha de Autorización: %s\n", autorizacion.FechaAutorizacion)
	fmt.Printf("⏱️  Tiempo total de procesamiento: %v\n", 
		resultadoEnvio.TiempoTotal + resultadoConsulta.TiempoTotal + tiempoEsperaInicial)

	return &autorizacion, nil
}

// MostrarEstadisticasReintento muestra estadísticas de reintentos
func MostrarEstadisticasReintento(resultado *ResultadoReintento) {
	fmt.Println("\n📊 ESTADÍSTICAS DE REINTENTOS")
	fmt.Println(strings.Repeat("-", 40))
	
	if resultado.Exitoso {
		fmt.Printf("✅ Resultado: EXITOSO\n")
	} else {
		fmt.Printf("❌ Resultado: FALLIDO\n")
	}
	
	fmt.Printf("🔢 Intentos realizados: %d\n", resultado.IntentosRealizados)
	fmt.Printf("⏱️  Tiempo total: %v\n", resultado.TiempoTotal)
	
	if resultado.UltimoError != nil {
		fmt.Printf("💥 Último error: %v\n", resultado.UltimoError)
		
		// Mostrar información detallada si es un ErrorSRI
		MostrarInformacionError(resultado.UltimoError)
	}
	
	if len(resultado.Errores) > 1 {
		fmt.Printf("\n📋 Historial de errores:\n")
		for i, err := range resultado.Errores {
			fmt.Printf("   %d. %v\n", i+1, err)
		}
	}
}