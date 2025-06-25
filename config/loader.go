package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// CargarConfiguracion - Carga la configuración desde archivos JSON
func CargarConfiguracion(archivoConfig string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(archivoConfig)
	if err != nil {
		return fmt.Errorf("error leyendo archivo de configuración %s: %v", archivoConfig, err)
	}
	
	// Parsear JSON a nuestra estructura
	err = json.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("error parseando JSON de configuración: %v", err)
	}
	
	// Validar que la configuración esté completa
	if err := validarConfiguracion(); err != nil {
		return fmt.Errorf("configuración inválida: %v", err)
	}
	
	return nil
}

// validarConfiguracion - Valida que todos los campos requeridos estén presentes
func validarConfiguracion() error {
	// Validar empresa
	if Config.Empresa.RazonSocial == "" {
		return fmt.Errorf("razón social de la empresa es requerida")
	}
	
	if Config.Empresa.RUC == "" {
		return fmt.Errorf("RUC de la empresa es requerido")
	}
	
	if len(Config.Empresa.RUC) != 13 {
		return fmt.Errorf("RUC debe tener exactamente 13 dígitos, tiene %d", len(Config.Empresa.RUC))
	}
	
	if Config.Empresa.Establecimiento == "" {
		return fmt.Errorf("establecimiento es requerido")
	}
	
	if Config.Empresa.PuntoEmision == "" {
		return fmt.Errorf("punto de emisión es requerido")
	}
	
	// Validar ambiente
	if Config.Ambiente.Codigo == "" {
		return fmt.Errorf("código de ambiente es requerido")
	}
	
	if Config.Ambiente.Codigo != "1" && Config.Ambiente.Codigo != "2" {
		return fmt.Errorf("código de ambiente debe ser '1' (pruebas) o '2' (producción)")
	}
	
	// Aplicar valores por defecto para campos opcionales
	aplicarValoresPorDefecto()
	
	return nil
}

// aplicarValoresPorDefecto aplica valores por defecto para configuraciones opcionales
func aplicarValoresPorDefecto() {
	// SRI defaults
	if Config.SRI.TimeoutSegundos == 0 {
		Config.SRI.TimeoutSegundos = 30
	}
	if Config.SRI.MaxReintentos == 0 {
		Config.SRI.MaxReintentos = 3
	}
	if Config.SRI.PolicyID == "" {
		Config.SRI.PolicyID = "https://www.sri.gob.ec/politica-de-firma"
	}
	
	// Database defaults
	if Config.Database.Ruta == "" {
		Config.Database.Ruta = "./facturacion.db"
	}
	if Config.Database.MaxConexiones == 0 {
		Config.Database.MaxConexiones = 10
	}
	
	// Endpoints según ambiente
	if Config.Ambiente.Codigo == "1" {
		// Ambiente de pruebas
		if Config.SRI.EndpointRecepcion == "" {
			Config.SRI.EndpointRecepcion = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline"
		}
		if Config.SRI.EndpointAutorizacion == "" {
			Config.SRI.EndpointAutorizacion = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline"
		}
	} else {
		// Ambiente de producción
		if Config.SRI.EndpointRecepcion == "" {
			Config.SRI.EndpointRecepcion = "https://cel.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline"
		}
		if Config.SRI.EndpointAutorizacion == "" {
			Config.SRI.EndpointAutorizacion = "https://cel.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline"
		}
	}
}

// CargarConfiguracionPorDefecto - Carga configuración de desarrollo si no existe archivo
func CargarConfiguracionPorDefecto() {
	Config = FacturacionConfig{
		Empresa: EmpresaConfig{
			RazonSocial:     "EMPRESA DEMO S.A.",
			RUC:             "1234567890001",
			Establecimiento: "001",
			PuntoEmision:    "001",
			Direccion:       "Av. Amazonas y Naciones Unidas",
		},
		Ambiente: AmbienteConfig{
			Codigo:      "1", // Pruebas
			Descripcion: "Ambiente de Pruebas",
			TipoEmision: "1", // Normal
		},
		SRI: SRIConfig{
			TimeoutSegundos: 30,
			MaxReintentos:   3,
			PolicyID:        "https://www.sri.gob.ec/politica-de-firma",
		},
		Database: DatabaseConfig{
			Ruta:          "./facturacion.db",
			MaxConexiones: 10,
		},
	}
	
	// Aplicar valores por defecto
	aplicarValoresPorDefecto()
}

// ObtenerSecuencialSiguiente - Genera el próximo secuencial de forma thread-safe
func ObtenerSecuencialSiguiente() string {
	// Incrementar contador de forma atómica
	nuevoSecuencial := atomic.AddInt64(&ContadorSecuencial, 1)
	
	// Formatear a 9 dígitos con ceros a la izquierda
	return fmt.Sprintf("%09d", nuevoSecuencial)
}

// GenerarClaveAcceso - Genera clave de acceso según algoritmo SRI
func GenerarClaveAcceso() string {
	// Obtener fecha actual en formato ddmmaaaa
	fecha := time.Now().Format("02012006")
	
	// Tipo de comprobante (01 = factura)
	tipoComprobante := "01"
	
	// RUC de la empresa
	ruc := Config.Empresa.RUC
	
	// Ambiente
	ambiente := Config.Ambiente.Codigo
	
	// Serie (establecimiento + punto emisión)
	serie := Config.Empresa.Establecimiento + Config.Empresa.PuntoEmision
	
	// Secuencial
	secuencial := ObtenerSecuencialSiguiente()
	
	// Código numérico (8 dígitos aleatorios)
	codigoNumerico := fmt.Sprintf("%08d", time.Now().UnixNano()%100000000)
	
	// Tipo de emisión
	tipoEmision := Config.Ambiente.TipoEmision
	
	// Construir clave sin dígito verificador
	claveSinDV := fecha + tipoComprobante + ruc + ambiente + serie + secuencial + codigoNumerico + tipoEmision
	
	// Calcular dígito verificador
	digitoVerificador := calcularDigitoVerificador(claveSinDV)
	
	// Clave completa
	return claveSinDV + strconv.Itoa(digitoVerificador)
}

// calcularDigitoVerificador calcula el dígito verificador de la clave de acceso
func calcularDigitoVerificador(clave string) int {
	// Algoritmo módulo 11 usado por SRI
	factor := 7
	suma := 0
	
	for i := 0; i < len(clave); i++ {
		digito, _ := strconv.Atoi(string(clave[i]))
		suma += digito * factor
		
		factor--
		if factor < 2 {
			factor = 7
		}
	}
	
	resto := suma % 11
	digitoVerificador := 11 - resto
	
	if digitoVerificador == 10 {
		digitoVerificador = 1
	} else if digitoVerificador == 11 {
		digitoVerificador = 0
	}
	
	return digitoVerificador
}

// ValidarClaveAcceso valida una clave de acceso existente
func ValidarClaveAcceso(claveAcceso string) error {
	if len(claveAcceso) != 49 {
		return fmt.Errorf("clave de acceso debe tener 49 dígitos, tiene %d", len(claveAcceso))
	}
	
	// Extraer partes
	claveSinDV := claveAcceso[:48]
	dv := claveAcceso[48:]
	
	// Calcular dígito verificador esperado
	dvEsperado := calcularDigitoVerificador(claveSinDV)
	dvActual, err := strconv.Atoi(dv)
	if err != nil {
		return fmt.Errorf("dígito verificador inválido: %v", err)
	}
	
	if dvActual != dvEsperado {
		return fmt.Errorf("dígito verificador incorrecto: esperado %d, obtenido %d", dvEsperado, dvActual)
	}
	
	return nil
}