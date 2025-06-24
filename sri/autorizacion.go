// Package sri implementa generación de números de autorización y claves de acceso del SRI
package sri

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// TipoComprobante tipos de comprobantes del SRI
type TipoComprobante int

const (
	Factura TipoComprobante = iota + 1
	NotaCredito
	NotaDebito
	GuiaRemision
	ComprobanteRetencion
	LiquidacionCompra
)

// String implementa Stringer para TipoComprobante
func (tc TipoComprobante) String() string {
	switch tc {
	case Factura:
		return "01"
	case NotaCredito:
		return "04"
	case NotaDebito:
		return "05"
	case GuiaRemision:
		return "06"
	case ComprobanteRetencion:
		return "07"
	case LiquidacionCompra:
		return "03"
	default:
		return "01"
	}
}

// Ambiente ambientes del SRI
type Ambiente int

const (
	Pruebas Ambiente = iota + 1
	Produccion
)

// String implementa Stringer para Ambiente
func (a Ambiente) String() string {
	switch a {
	case Pruebas:
		return "1"
	case Produccion:
		return "2"
	default:
		return "1"
	}
}

// TipoEmision tipos de emisión
type TipoEmision int

const (
	EmisionNormal TipoEmision = iota + 1
	EmisionContingencia
)

// String implementa Stringer para TipoEmision
func (te TipoEmision) String() string {
	switch te {
	case EmisionNormal:
		return "1"
	case EmisionContingencia:
		return "2"
	default:
		return "1"
	}
}

// ClaveAccesoConfig configuración para generar clave de acceso
type ClaveAccesoConfig struct {
	FechaEmision    time.Time
	TipoComprobante TipoComprobante
	RUCEmisor       string
	Ambiente        Ambiente
	Serie           string // Ejemplo: "001001"
	NumeroSecuencial string // Ejemplo: "000000001"
	CodigoNumerico  string // 8 dígitos aleatorios
	TipoEmision     TipoEmision
}

// AutorizacionInfo información de autorización del SRI
type AutorizacionInfo struct {
	ClaveAcceso         string    `json:"claveAcceso"`
	NumeroAutorizacion  string    `json:"numeroAutorizacion"`
	FechaAutorizacion   time.Time `json:"fechaAutorizacion"`
	Estado              string    `json:"estado"`
	Ambiente            string    `json:"ambiente"`
	TipoEmision         string    `json:"tipoEmision"`
}

// GenerarClaveAcceso genera una clave de acceso según especificaciones SRI
func GenerarClaveAcceso(config ClaveAccesoConfig) (string, error) {
	// Validar RUC
	if len(config.RUCEmisor) != 13 {
		return "", fmt.Errorf("RUC debe tener 13 dígitos")
	}

	// Validar serie
	if len(config.Serie) != 6 {
		return "", fmt.Errorf("serie debe tener 6 dígitos (ej: 001001)")
	}

	// Validar número secuencial
	if len(config.NumeroSecuencial) != 9 {
		return "", fmt.Errorf("número secuencial debe tener 9 dígitos")
	}

	// Generar código numérico si no se proporciona
	codigoNumerico := config.CodigoNumerico
	if codigoNumerico == "" {
		codigoNumerico = generarCodigoNumerico()
	}

	if len(codigoNumerico) != 8 {
		return "", fmt.Errorf("código numérico debe tener 8 dígitos")
	}

	// Construir clave de acceso (49 dígitos sin dígito verificador)
	fechaString := config.FechaEmision.Format("02012006") // ddMMyyyy
	
	claveBase := fechaString +                           // 8 dígitos
		config.TipoComprobante.String() +               // 2 dígitos
		config.RUCEmisor +                              // 13 dígitos
		config.Ambiente.String() +                      // 1 dígito
		config.Serie +                                  // 6 dígitos
		config.NumeroSecuencial +                       // 9 dígitos
		codigoNumerico +                                // 8 dígitos
		config.TipoEmision.String()                     // 1 dígito
		// Total: 48 dígitos

	// Calcular dígito verificador (módulo 11)
	digitoVerificador := calcularDigitoVerificador(claveBase)

	// Clave de acceso completa (49 dígitos)
	claveCompleta := claveBase + strconv.Itoa(digitoVerificador)

	return claveCompleta, nil
}

// generarCodigoNumerico genera 8 dígitos aleatorios
func generarCodigoNumerico() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// calcularDigitoVerificador calcula el dígito verificador usando módulo 11
func calcularDigitoVerificador(clave string) int {
	// Algoritmo módulo 11 utilizado por el SRI
	multiplicadores := []int{2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7, 2, 3, 4, 5, 6, 7}
	
	suma := 0
	for i, char := range clave {
		digito, _ := strconv.Atoi(string(char))
		suma += digito * multiplicadores[i]
	}

	residuo := suma % 11
	
	switch residuo {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return 11 - residuo
	}
}

// ValidarClaveAcceso valida el formato y dígito verificador de una clave de acceso
func ValidarClaveAcceso(claveAcceso string) error {
	// Validar longitud
	if len(claveAcceso) != 49 {
		return fmt.Errorf("clave de acceso debe tener 49 dígitos")
	}

	// Verificar que todos sean números
	for _, char := range claveAcceso {
		if char < '0' || char > '9' {
			return fmt.Errorf("clave de acceso debe contener solo números")
		}
	}

	// Extraer dígito verificador
	claveBase := claveAcceso[:48]
	digitoVerificadorRecibido, _ := strconv.Atoi(claveAcceso[48:])

	// Calcular dígito verificador esperado
	digitoVerificadorCalculado := calcularDigitoVerificador(claveBase)

	if digitoVerificadorRecibido != digitoVerificadorCalculado {
		return fmt.Errorf("dígito verificador inválido: esperado %d, recibido %d", 
			digitoVerificadorCalculado, digitoVerificadorRecibido)
	}

	return nil
}

// ParsearClaveAcceso extrae información de una clave de acceso
func ParsearClaveAcceso(claveAcceso string) (ClaveAccesoConfig, error) {
	if err := ValidarClaveAcceso(claveAcceso); err != nil {
		return ClaveAccesoConfig{}, err
	}

	// Extraer componentes
	fechaStr := claveAcceso[0:8]
	tipoCompStr := claveAcceso[8:10]
	rucEmisor := claveAcceso[10:23]
	ambienteStr := claveAcceso[23:24]
	serie := claveAcceso[24:30]
	secuencial := claveAcceso[30:39]
	codigoNumerico := claveAcceso[39:47]
	tipoEmisionStr := claveAcceso[47:48]

	// Parsear fecha
	fecha, err := time.Parse("02012006", fechaStr)
	if err != nil {
		return ClaveAccesoConfig{}, fmt.Errorf("fecha inválida en clave de acceso: %v", err)
	}

	// Parsear tipo de comprobante
	tipoComp, _ := strconv.Atoi(tipoCompStr)
	
	// Parsear ambiente
	amb, _ := strconv.Atoi(ambienteStr)
	
	// Parsear tipo de emisión
	tipoEm, _ := strconv.Atoi(tipoEmisionStr)

	config := ClaveAccesoConfig{
		FechaEmision:     fecha,
		TipoComprobante:  TipoComprobante(tipoComp),
		RUCEmisor:        rucEmisor,
		Ambiente:         Ambiente(amb),
		Serie:            serie,
		NumeroSecuencial: secuencial,
		CodigoNumerico:   codigoNumerico,
		TipoEmision:      TipoEmision(tipoEm),
	}

	return config, nil
}

// FormatearClaveAcceso formatea una clave de acceso para visualización
func FormatearClaveAcceso(claveAcceso string) string {
	if len(claveAcceso) != 49 {
		return claveAcceso
	}

	// Insertar guiones para legibilidad
	parts := []string{
		claveAcceso[0:8],   // Fecha
		claveAcceso[8:10],  // Tipo comprobante
		claveAcceso[10:23], // RUC
		claveAcceso[23:24], // Ambiente
		claveAcceso[24:30], // Serie
		claveAcceso[30:39], // Secuencial
		claveAcceso[39:47], // Código numérico
		claveAcceso[47:48], // Tipo emisión
		claveAcceso[48:49], // Dígito verificador
	}

	return strings.Join(parts, "-")
}

// GenerarNumeroAutorizacion simula la generación de un número de autorización del SRI
func GenerarNumeroAutorizacion(claveAcceso string) string {
	// En producción, este sería devuelto por el SRI
	// Por ahora retornamos la misma clave de acceso como número de autorización
	// (en ambiente de pruebas)
	return claveAcceso
}

// SimularAutorizacionSRI simula el proceso de autorización del SRI
func SimularAutorizacionSRI(claveAcceso string, ambiente Ambiente) AutorizacionInfo {
	autorizacion := AutorizacionInfo{
		ClaveAcceso:        claveAcceso,
		NumeroAutorizacion: GenerarNumeroAutorizacion(claveAcceso),
		FechaAutorizacion:  time.Now(),
		Estado:             "AUTORIZADO",
		Ambiente:           ambiente.String(),
		TipoEmision:        "1", // Normal
	}

	return autorizacion
}

// MostrarInformacionClaveAcceso muestra información detallada de una clave de acceso
func MostrarInformacionClaveAcceso(claveAcceso string) {
	fmt.Println("\n🔑 INFORMACIÓN DE CLAVE DE ACCESO")
	fmt.Println("=================================")
	fmt.Printf("🎯 Clave de Acceso: %s\n", FormatearClaveAcceso(claveAcceso))
	
	config, err := ParsearClaveAcceso(claveAcceso)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("📅 Fecha de Emisión: %s\n", config.FechaEmision.Format("02/01/2006"))
	fmt.Printf("📋 Tipo Comprobante: %s\n", obtenerNombreTipoComprobante(config.TipoComprobante))
	fmt.Printf("🏢 RUC Emisor: %s\n", config.RUCEmisor)
	fmt.Printf("🌍 Ambiente: %s\n", obtenerNombreAmbiente(config.Ambiente))
	fmt.Printf("📍 Serie: %s\n", config.Serie)
	fmt.Printf("🔢 Número Secuencial: %s\n", config.NumeroSecuencial)
	fmt.Printf("🎲 Código Numérico: %s\n", config.CodigoNumerico)
	fmt.Printf("📤 Tipo Emisión: %s\n", obtenerNombreTipoEmision(config.TipoEmision))
	
	// Validar
	if err := ValidarClaveAcceso(claveAcceso); err != nil {
		fmt.Printf("❌ Validación: %v\n", err)
	} else {
		fmt.Printf("✅ Validación: Clave de acceso válida\n")
	}
}

// Funciones auxiliares para nombres descriptivos
func obtenerNombreTipoComprobante(tc TipoComprobante) string {
	switch tc {
	case Factura:
		return "Factura (01)"
	case NotaCredito:
		return "Nota de Crédito (04)"
	case NotaDebito:
		return "Nota de Débito (05)"
	case GuiaRemision:
		return "Guía de Remisión (06)"
	case ComprobanteRetencion:
		return "Comprobante de Retención (07)"
	case LiquidacionCompra:
		return "Liquidación de Compra (03)"
	default:
		return "Desconocido"
	}
}

func obtenerNombreAmbiente(a Ambiente) string {
	switch a {
	case Pruebas:
		return "Pruebas (1)"
	case Produccion:
		return "Producción (2)"
	default:
		return "Desconocido"
	}
}

func obtenerNombreTipoEmision(te TipoEmision) string {
	switch te {
	case EmisionNormal:
		return "Normal (1)"
	case EmisionContingencia:
		return "Contingencia (2)"
	default:
		return "Desconocido"
	}
}