// Package sri implementa manejo avanzado de errores del SRI Ecuador
package sri

import (
	"fmt"
	"strings"
)

// TipoErrorSRI tipos específicos de errores del SRI
type TipoErrorSRI int

const (
	ErrorConexion TipoErrorSRI = iota + 1
	ErrorValidacion
	ErrorAutenticacion
	ErrorFormato
	ErrorDatos
	ErrorSistema
	ErrorTimeout
	ErrorCertificado
	ErrorFirma
	ErrorClaveAcceso
)

// String implementa Stringer para TipoErrorSRI
func (te TipoErrorSRI) String() string {
	switch te {
	case ErrorConexion:
		return "ERROR_CONEXION"
	case ErrorValidacion:
		return "ERROR_VALIDACION"
	case ErrorAutenticacion:
		return "ERROR_AUTENTICACION"
	case ErrorFormato:
		return "ERROR_FORMATO"
	case ErrorDatos:
		return "ERROR_DATOS"
	case ErrorSistema:
		return "ERROR_SISTEMA"
	case ErrorTimeout:
		return "ERROR_TIMEOUT"
	case ErrorCertificado:
		return "ERROR_CERTIFICADO"
	case ErrorFirma:
		return "ERROR_FIRMA"
	case ErrorClaveAcceso:
		return "ERROR_CLAVE_ACCESO"
	default:
		return "ERROR_DESCONOCIDO"
	}
}

// ErrorSRI estructura para errores específicos del SRI
type ErrorSRI struct {
	Tipo         TipoErrorSRI `json:"tipo"`
	Codigo       string       `json:"codigo"`
	Mensaje      string       `json:"mensaje"`
	Detalle      string       `json:"detalle"`
	Recuperable  bool         `json:"recuperable"`
	SugerenciaFix string      `json:"sugerencia_fix"`
}

// Error implementa la interfaz error
func (e *ErrorSRI) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Tipo, e.Codigo, e.Mensaje)
}

// String implementa la interfaz Stringer para mejor presentación
func (e *ErrorSRI) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ErrorSRI{\n"))
	builder.WriteString(fmt.Sprintf("  Tipo: %s\n", e.Tipo))
	builder.WriteString(fmt.Sprintf("  Código: %s\n", e.Codigo))
	builder.WriteString(fmt.Sprintf("  Mensaje: %s\n", e.Mensaje))
	if e.Detalle != "" {
		builder.WriteString(fmt.Sprintf("  Detalle: %s\n", e.Detalle))
	}
	builder.WriteString(fmt.Sprintf("  Recuperable: %t\n", e.Recuperable))
	if e.SugerenciaFix != "" {
		builder.WriteString(fmt.Sprintf("  Sugerencia: %s\n", e.SugerenciaFix))
	}
	builder.WriteString("}")
	return builder.String()
}

// IsRecuperable indica si el error es recuperable con reintentos
func (e *ErrorSRI) IsRecuperable() bool {
	return e.Recuperable
}

// GetSugerencia obtiene sugerencia para resolver el error
func (e *ErrorSRI) GetSugerencia() string {
	return e.SugerenciaFix
}

// Códigos de error específicos del SRI Ecuador
var codigosErrorSRI = map[string]*ErrorSRI{
	// Errores de conexión
	"CONEXION_TIMEOUT": {
		Tipo:          ErrorTimeout,
		Codigo:        "CONEXION_TIMEOUT",
		Mensaje:       "Timeout conectando con SRI",
		Detalle:       "El servidor del SRI no respondió en el tiempo esperado",
		Recuperable:   true,
		SugerenciaFix: "Reintentar en unos minutos. Verificar conexión a internet.",
	},
	"SERVIDOR_NO_DISPONIBLE": {
		Tipo:          ErrorConexion,
		Codigo:        "SERVIDOR_NO_DISPONIBLE",
		Mensaje:       "Servidor SRI no disponible",
		Detalle:       "El servicio web del SRI está temporalmente inaccesible",
		Recuperable:   true,
		SugerenciaFix: "Reintentar más tarde. Verificar estado de servicios SRI.",
	},

	// Errores de validación XML
	"CLAVE-01": {
		Tipo:          ErrorClaveAcceso,
		Codigo:        "CLAVE-01",
		Mensaje:       "Clave de acceso registrada",
		Detalle:       "La clave de acceso ya está registrada en el SRI",
		Recuperable:   false,
		SugerenciaFix: "Generar nueva clave de acceso con secuencial diferente.",
	},
	"CLAVE-02": {
		Tipo:          ErrorClaveAcceso,
		Codigo:        "CLAVE-02",
		Mensaje:       "Clave de acceso mal formada",
		Detalle:       "La clave de acceso no cumple con el formato de 49 dígitos",
		Recuperable:   false,
		SugerenciaFix: "Verificar algoritmo de generación de clave de acceso.",
	},
	"CLAVE-03": {
		Tipo:          ErrorClaveAcceso,
		Codigo:        "CLAVE-03",
		Mensaje:       "Fecha de emisión incorrecta",
		Detalle:       "La fecha en la clave de acceso no coincide con el XML",
		Recuperable:   false,
		SugerenciaFix: "Sincronizar fechas entre clave de acceso y XML.",
	},

	// Errores de estructura XML
	"ESTRUCTURA-01": {
		Tipo:          ErrorFormato,
		Codigo:        "ESTRUCTURA-01",
		Mensaje:       "Estructura XML incorrecta",
		Detalle:       "El XML no cumple con el esquema XSD del SRI",
		Recuperable:   false,
		SugerenciaFix: "Validar XML contra esquema oficial del SRI.",
	},
	"ESTRUCTURA-02": {
		Tipo:          ErrorFormato,
		Codigo:        "ESTRUCTURA-02",
		Mensaje:       "Codificación incorrecta",
		Detalle:       "El XML debe estar codificado en UTF-8",
		Recuperable:   false,
		SugerenciaFix: "Verificar codificación UTF-8 del XML.",
	},

	// Errores de datos
	"RUC-01": {
		Tipo:          ErrorDatos,
		Codigo:        "RUC-01",
		Mensaje:       "RUC no válido",
		Detalle:       "El RUC del emisor no está registrado en el SRI",
		Recuperable:   false,
		SugerenciaFix: "Verificar RUC en portal SRI. Debe estar activo.",
	},
	"CEDULA-01": {
		Tipo:          ErrorDatos,
		Codigo:        "CEDULA-01",
		Mensaje:       "Cédula no válida",
		Detalle:       "La cédula del comprador no cumple algoritmo de validación",
		Recuperable:   false,
		SugerenciaFix: "Verificar dígito verificador de la cédula ecuatoriana.",
	},

	// Errores de certificado
	"CERT-01": {
		Tipo:          ErrorCertificado,
		Codigo:        "CERT-01",
		Mensaje:       "Certificado expirado",
		Detalle:       "El certificado digital usado para firmar ha expirado",
		Recuperable:   false,
		SugerenciaFix: "Renovar certificado digital con entidad certificadora autorizada.",
	},
	"CERT-02": {
		Tipo:          ErrorCertificado,
		Codigo:        "CERT-02",
		Mensaje:       "Certificado revocado",
		Detalle:       "El certificado digital ha sido revocado",
		Recuperable:   false,
		SugerenciaFix: "Obtener nuevo certificado digital.",
	},
	"CERT-03": {
		Tipo:          ErrorCertificado,
		Codigo:        "CERT-03",
		Mensaje:       "Cadena de certificación inválida",
		Detalle:       "La cadena de certificación no es válida",
		Recuperable:   false,
		SugerenciaFix: "Verificar cadena completa de certificación.",
	},

	// Errores de firma
	"FIRMA-01": {
		Tipo:          ErrorFirma,
		Codigo:        "FIRMA-01",
		Mensaje:       "Firma digital inválida",
		Detalle:       "La firma XAdES-BES no es válida",
		Recuperable:   false,
		SugerenciaFix: "Verificar proceso de firma digital XAdES-BES.",
	},
	"FIRMA-02": {
		Tipo:          ErrorFirma,
		Codigo:        "FIRMA-02",
		Mensaje:       "Algoritmo de firma no soportado",
		Detalle:       "El algoritmo de firma no está soportado por SRI",
		Recuperable:   false,
		SugerenciaFix: "Usar algoritmo RSA-SHA256 para firma.",
	},

	// Errores del sistema SRI
	"SRI-01": {
		Tipo:          ErrorSistema,
		Codigo:        "SRI-01",
		Mensaje:       "Sistema en mantenimiento",
		Detalle:       "El sistema del SRI está en mantenimiento programado",
		Recuperable:   true,
		SugerenciaFix: "Reintentar después del horario de mantenimiento.",
	},
	"SRI-02": {
		Tipo:          ErrorSistema,
		Codigo:        "SRI-02",
		Mensaje:       "Sobrecarga del sistema",
		Detalle:       "El sistema del SRI está experimentando alta carga",
		Recuperable:   true,
		SugerenciaFix: "Reintentar con intervalos exponenciales.",
	},
}

// ParsearErrorSRI parsea mensajes de error del SRI y los clasifica
func ParsearErrorSRI(mensajeError string, codigoHTTP int) *ErrorSRI {
	mensajeLower := strings.ToLower(mensajeError)

	// Buscar por código específico en el mensaje
	for codigo, errorInfo := range codigosErrorSRI {
		if strings.Contains(mensajeLower, strings.ToLower(codigo)) ||
		   strings.Contains(mensajeLower, strings.ToLower(errorInfo.Mensaje)) {
			return &ErrorSRI{
				Tipo:          errorInfo.Tipo,
				Codigo:        errorInfo.Codigo,
				Mensaje:       errorInfo.Mensaje,
				Detalle:       errorInfo.Detalle,
				Recuperable:   errorInfo.Recuperable,
				SugerenciaFix: errorInfo.SugerenciaFix,
			}
		}
	}

	// Clasificar por código HTTP
	switch codigoHTTP {
	case 408, 504:
		return &ErrorSRI{
			Tipo:          ErrorTimeout,
			Codigo:        "HTTP_TIMEOUT",
			Mensaje:       "Timeout en petición HTTP",
			Detalle:       fmt.Sprintf("Código HTTP: %d", codigoHTTP),
			Recuperable:   true,
			SugerenciaFix: "Reintentar petición con timeout mayor.",
		}
	case 500, 502, 503:
		return &ErrorSRI{
			Tipo:          ErrorSistema,
			Codigo:        "HTTP_SERVER_ERROR",
			Mensaje:       "Error interno del servidor SRI",
			Detalle:       fmt.Sprintf("Código HTTP: %d", codigoHTTP),
			Recuperable:   true,
			SugerenciaFix: "El SRI tiene problemas internos. Reintentar más tarde.",
		}
	case 401, 403:
		return &ErrorSRI{
			Tipo:          ErrorAutenticacion,
			Codigo:        "HTTP_AUTH_ERROR",
			Mensaje:       "Error de autenticación",
			Detalle:       fmt.Sprintf("Código HTTP: %d", codigoHTTP),
			Recuperable:   false,
			SugerenciaFix: "Verificar certificado digital y credenciales.",
		}
	case 400:
		return &ErrorSRI{
			Tipo:          ErrorValidacion,
			Codigo:        "HTTP_BAD_REQUEST",
			Mensaje:       "Petición malformada",
			Detalle:       fmt.Sprintf("Código HTTP: %d", codigoHTTP),
			Recuperable:   false,
			SugerenciaFix: "Verificar formato de la petición SOAP.",
		}
	}

	// Clasificar por contenido del mensaje
	if strings.Contains(mensajeLower, "timeout") ||
	   strings.Contains(mensajeLower, "connection") {
		return &ErrorSRI{
			Tipo:          ErrorConexion,
			Codigo:        "CONEXION_GENERAL",
			Mensaje:       "Error de conexión",
			Detalle:       mensajeError,
			Recuperable:   true,
			SugerenciaFix: "Verificar conectividad. Reintentar.",
		}
	}

	if strings.Contains(mensajeLower, "xml") ||
	   strings.Contains(mensajeLower, "schema") ||
	   strings.Contains(mensajeLower, "formato") {
		return &ErrorSRI{
			Tipo:          ErrorFormato,
			Codigo:        "FORMATO_GENERAL",
			Mensaje:       "Error de formato",
			Detalle:       mensajeError,
			Recuperable:   false,
			SugerenciaFix: "Verificar estructura XML contra especificaciones SRI.",
		}
	}

	if strings.Contains(mensajeLower, "certificado") ||
	   strings.Contains(mensajeLower, "firma") {
		return &ErrorSRI{
			Tipo:          ErrorCertificado,
			Codigo:        "CERT_GENERAL",
			Mensaje:       "Error de certificado",
			Detalle:       mensajeError,
			Recuperable:   false,
			SugerenciaFix: "Verificar certificado digital y proceso de firma.",
		}
	}

	// Error genérico
	return &ErrorSRI{
		Tipo:          ErrorSistema,
		Codigo:        "ERROR_GENERAL",
		Mensaje:       "Error no clasificado",
		Detalle:       mensajeError,
		Recuperable:   true,
		SugerenciaFix: "Revisar logs detallados y contactar soporte técnico.",
	}
}

// CrearErrorConexion crea un error específico de conexión
func CrearErrorConexion(detalle string) *ErrorSRI {
	return &ErrorSRI{
		Tipo:          ErrorConexion,
		Codigo:        "CONEXION_FALLO",
		Mensaje:       "Fallo de conexión con SRI",
		Detalle:       detalle,
		Recuperable:   true,
		SugerenciaFix: "Verificar conexión a internet y estado de servicios SRI.",
	}
}

// CrearErrorValidacion crea un error específico de validación
func CrearErrorValidacion(campo, detalle string) *ErrorSRI {
	return &ErrorSRI{
		Tipo:          ErrorValidacion,
		Codigo:        "VALIDACION_" + strings.ToUpper(campo),
		Mensaje:       fmt.Sprintf("Error de validación en campo: %s", campo),
		Detalle:       detalle,
		Recuperable:   false,
		SugerenciaFix: fmt.Sprintf("Corregir el valor del campo %s según especificaciones SRI.", campo),
	}
}

// CrearErrorCertificado crea un error específico de certificado
func CrearErrorCertificado(detalle string) *ErrorSRI {
	return &ErrorSRI{
		Tipo:          ErrorCertificado,
		Codigo:        "CERT_ERROR",
		Mensaje:       "Error con certificado digital",
		Detalle:       detalle,
		Recuperable:   false,
		SugerenciaFix: "Verificar validez y configuración del certificado digital.",
	}
}

// EsErrorRecuperable determina si un error permite reintentos
func EsErrorRecuperable(err error) bool {
	if errorSRI, ok := err.(*ErrorSRI); ok {
		return errorSRI.IsRecuperable()
	}
	return false
}

// ObtenerSugerencia obtiene sugerencia para resolver un error
func ObtenerSugerencia(err error) string {
	if errorSRI, ok := err.(*ErrorSRI); ok {
		return errorSRI.GetSugerencia()
	}
	return "Error no clasificado. Revisar logs y documentación SRI."
}

// MostrarInformacionError muestra información detallada de un error SRI
func MostrarInformacionError(err error) {
	if errorSRI, ok := err.(*ErrorSRI); ok {
		fmt.Printf("\n❌ ERROR SRI DETECTADO\n")
		fmt.Printf("=======================\n")
		fmt.Printf("🔍 Tipo: %s\n", errorSRI.Tipo)
		fmt.Printf("📋 Código: %s\n", errorSRI.Codigo)
		fmt.Printf("💬 Mensaje: %s\n", errorSRI.Mensaje)
		fmt.Printf("📝 Detalle: %s\n", errorSRI.Detalle)
		
		if errorSRI.Recuperable {
			fmt.Printf("🔄 Recuperable: ✅ SÍ\n")
		} else {
			fmt.Printf("🔄 Recuperable: ❌ NO\n")
		}
		
		fmt.Printf("💡 Sugerencia: %s\n", errorSRI.SugerenciaFix)
		fmt.Printf("=======================\n")
	} else {
		fmt.Printf("❌ Error general: %v\n", err)
	}
}