// Package sri implementa cliente SOAP para comunicación con SRI Ecuador
package sri

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Endpoints oficiales del SRI Ecuador
const (
	// Ambiente de Certificación (Pruebas)
	EndpointRecepcionCertificacion   = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline"
	EndpointAutorizacionCertificacion = "https://celcer.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline"
	
	// Ambiente de Producción
	EndpointRecepcionProduccion   = "https://cel.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline"
	EndpointAutorizacionProduccion = "https://cel.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline"
)

// SOAPClient cliente SOAP para SRI Ecuador
type SOAPClient struct {
	Ambiente        Ambiente
	TimeoutSegundos int
	httpClient      *http.Client
}

// RespuestaSolicitud respuesta del servicio de recepción SRI
type RespuestaSolicitud struct {
	XMLName xml.Name `xml:"respuestaSolicitud"`
	Estado  string   `xml:"estado"`
	Comprobantes []ComprobanteRecepcion `xml:"comprobantes>comprobante"`
}

// ComprobanteRecepcion información de recepción de comprobante
type ComprobanteRecepcion struct {
	XMLName     xml.Name `xml:"comprobante"`
	ClaveAcceso string   `xml:"claveAcceso"`
	Mensajes    []MensajeSRI `xml:"mensajes>mensaje"`
}

// MensajeSRI mensaje del SRI (errores, advertencias, etc.)
type MensajeSRI struct {
	XMLName        xml.Name `xml:"mensaje"`
	Identificador  string   `xml:"identificador"`
	Mensaje        string   `xml:"mensaje"`
	InformacionAdicional string `xml:"informacionAdicional"`
	Tipo           string   `xml:"tipo"`
}

// RespuestaComprobante respuesta del servicio de autorización SRI
type RespuestaComprobante struct {
	XMLName           xml.Name `xml:"respuestaComprobante"`
	ClaveAccesoConsultada string `xml:"claveAccesoConsultada"`
	NumeroComprobantes    string `xml:"numeroComprobantes"`
	Autorizaciones        []AutorizacionSRI `xml:"autorizaciones>autorizacion"`
}

// AutorizacionSRI información de autorización del SRI
type AutorizacionSRI struct {
	XMLName           xml.Name `xml:"autorizacion"`
	Estado            string   `xml:"estado"`
	NumeroAutorizacion string  `xml:"numeroAutorizacion"`
	FechaAutorizacion string   `xml:"fechaAutorizacion"`
	Ambiente          string   `xml:"ambiente"`
	Comprobante       string   `xml:"comprobante"` // XML autorizado en CDATA
	Mensajes          []MensajeSRI `xml:"mensajes>mensaje"`
}

// SolicitudRecepcion estructura para envío a SRI
type SolicitudRecepcion struct {
	XMLName     xml.Name `xml:"soap:Envelope"`
	SoapNS      string   `xml:"xmlns:soap,attr"`
	SriNS       string   `xml:"xmlns:sri,attr"`
	Body        BodyRecepcion `xml:"soap:Body"`
}

// BodyRecepcion cuerpo del SOAP para recepción
type BodyRecepcion struct {
	XMLName         xml.Name `xml:"soap:Body"`
	ValidarComprobante ValidarComprobante `xml:"sri:validarComprobante"`
}

// ValidarComprobante operación de validación
type ValidarComprobante struct {
	XMLName xml.Name `xml:"sri:validarComprobante"`
	XML     string   `xml:"xml"` // XML en base64
}

// SolicitudAutorizacion estructura para consulta de autorización
type SolicitudAutorizacion struct {
	XMLName     xml.Name `xml:"soap:Envelope"`
	SoapNS      string   `xml:"xmlns:soap,attr"`
	SriNS       string   `xml:"xmlns:sri,attr"`
	Body        BodyAutorizacion `xml:"soap:Body"`
}

// BodyAutorizacion cuerpo del SOAP para autorización
type BodyAutorizacion struct {
	XMLName              xml.Name `xml:"soap:Body"`
	AutorizarComprobante AutorizarComprobante `xml:"sri:autorizacionComprobante"`
}

// AutorizarComprobante operación de autorización
type AutorizarComprobante struct {
	XMLName     xml.Name `xml:"sri:autorizacionComprobante"`
	ClaveAcceso string   `xml:"claveAccesoComprobante"`
}

// NewSOAPClient crea un nuevo cliente SOAP para SRI
func NewSOAPClient(ambiente Ambiente) *SOAPClient {
	// Configurar cliente HTTP con timeout y TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // En producción debe ser false
		},
	}
	
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
	
	return &SOAPClient{
		Ambiente:        ambiente,
		TimeoutSegundos: 30,
		httpClient:      client,
	}
}

// EnviarComprobante envía un comprobante XML al SRI para validación
func (c *SOAPClient) EnviarComprobante(xmlComprobante []byte) (*RespuestaSolicitud, error) {
	// Codificar XML en base64
	xmlBase64 := base64.StdEncoding.EncodeToString(xmlComprobante)
	
	// Crear solicitud SOAP
	solicitud := SolicitudRecepcion{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		SriNS:  "http://ec.gob.sri.ws.recepcion",
		Body: BodyRecepcion{
			ValidarComprobante: ValidarComprobante{
				XML: xmlBase64,
			},
		},
	}
	
	// Serializar a XML
	soapXML, err := xml.MarshalIndent(solicitud, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando solicitud SOAP: %v", err)
	}
	
	// Agregar header XML
	soapRequest := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(soapXML))
	
	// Determinar endpoint según ambiente
	endpoint := EndpointRecepcionCertificacion
	if c.Ambiente == Produccion {
		endpoint = EndpointRecepcionProduccion
	}
	
	// Crear petición HTTP
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %v", err)
	}
	
	// Headers requeridos por SRI
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(soapRequest)))
	
	// Enviar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando petición al SRI: %v", err)
	}
	defer resp.Body.Close()
	
	// Leer respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta del SRI: %v", err)
	}
	
	// Verificar código de respuesta HTTP
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("SRI respondió con código %d: %s", resp.StatusCode, string(respBody))
	}
	
	// Parsear respuesta SOAP
	return c.parsearRespuestaRecepcion(respBody)
}

// ConsultarAutorizacion consulta el estado de autorización de un comprobante
func (c *SOAPClient) ConsultarAutorizacion(claveAcceso string) (*RespuestaComprobante, error) {
	// Crear solicitud SOAP
	solicitud := SolicitudAutorizacion{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		SriNS:  "http://ec.gob.sri.ws.autorizacion",
		Body: BodyAutorizacion{
			AutorizarComprobante: AutorizarComprobante{
				ClaveAcceso: claveAcceso,
			},
		},
	}
	
	// Serializar a XML
	soapXML, err := xml.MarshalIndent(solicitud, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando solicitud SOAP: %v", err)
	}
	
	// Agregar header XML
	soapRequest := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(soapXML))
	
	// Determinar endpoint según ambiente
	endpoint := EndpointAutorizacionCertificacion
	if c.Ambiente == Produccion {
		endpoint = EndpointAutorizacionProduccion
	}
	
	// Crear petición HTTP
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(soapRequest))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %v", err)
	}
	
	// Headers requeridos por SRI
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(soapRequest)))
	
	// Enviar petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando petición al SRI: %v", err)
	}
	defer resp.Body.Close()
	
	// Leer respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta del SRI: %v", err)
	}
	
	// Verificar código de respuesta HTTP
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("SRI respondió con código %d: %s", resp.StatusCode, string(respBody))
	}
	
	// Parsear respuesta SOAP
	return c.parsearRespuestaAutorizacion(respBody)
}

// parsearRespuestaRecepcion parsea la respuesta SOAP de recepción
func (c *SOAPClient) parsearRespuestaRecepcion(respXML []byte) (*RespuestaSolicitud, error) {
	// En una implementación real, necesitaríamos parsear el SOAP envelope completo
	// Por simplicidad, buscamos el contenido entre las tags de respuesta
	
	respStr := string(respXML)
	
	// Buscar el contenido de respuestaSolicitud
	inicio := strings.Index(respStr, "<ns2:respuestaSolicitud")
	if inicio == -1 {
		inicio = strings.Index(respStr, "<respuestaSolicitud")
	}
	if inicio == -1 {
		return nil, fmt.Errorf("no se encontró respuestaSolicitud en la respuesta del SRI")
	}
	
	fin := strings.Index(respStr[inicio:], "</ns2:respuestaSolicitud>")
	if fin == -1 {
		fin = strings.Index(respStr[inicio:], "</respuestaSolicitud>")
	}
	if fin == -1 {
		return nil, fmt.Errorf("respuesta del SRI mal formada")
	}
	
	// Extraer el XML de respuesta limpio
	respuestaXML := respStr[inicio : inicio+fin] + "</respuestaSolicitud>"
	
	// Limpiar namespaces para simplificar el parsing
	respuestaXML = strings.ReplaceAll(respuestaXML, "ns2:", "")
	respuestaXML = strings.ReplaceAll(respuestaXML, "ns3:", "")
	
	// Parsear XML
	var respuesta RespuestaSolicitud
	err := xml.Unmarshal([]byte(respuestaXML), &respuesta)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta del SRI: %v", err)
	}
	
	return &respuesta, nil
}

// parsearRespuestaAutorizacion parsea la respuesta SOAP de autorización
func (c *SOAPClient) parsearRespuestaAutorizacion(respXML []byte) (*RespuestaComprobante, error) {
	// Similar al método anterior pero para respuestas de autorización
	
	respStr := string(respXML)
	
	// Buscar el contenido de respuestaComprobante
	inicio := strings.Index(respStr, "<ns2:respuestaComprobante")
	if inicio == -1 {
		inicio = strings.Index(respStr, "<respuestaComprobante")
	}
	if inicio == -1 {
		return nil, fmt.Errorf("no se encontró respuestaComprobante en la respuesta del SRI")
	}
	
	fin := strings.Index(respStr[inicio:], "</ns2:respuestaComprobante>")
	if fin == -1 {
		fin = strings.Index(respStr[inicio:], "</respuestaComprobante>")
	}
	if fin == -1 {
		return nil, fmt.Errorf("respuesta del SRI mal formada")
	}
	
	// Extraer el XML de respuesta limpio
	respuestaXML := respStr[inicio : inicio+fin] + "</respuestaComprobante>"
	
	// Limpiar namespaces
	respuestaXML = strings.ReplaceAll(respuestaXML, "ns2:", "")
	respuestaXML = strings.ReplaceAll(respuestaXML, "ns3:", "")
	
	// Parsear XML
	var respuesta RespuestaComprobante
	err := xml.Unmarshal([]byte(respuestaXML), &respuesta)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta del SRI: %v", err)
	}
	
	return &respuesta, nil
}

// ProcesarComprobanteCompleto procesa un comprobante de forma completa: envío + autorización
func (c *SOAPClient) ProcesarComprobanteCompleto(xmlComprobante []byte, claveAcceso string) (*AutorizacionSRI, error) {
	fmt.Println("📤 Enviando comprobante al SRI...")
	
	// Paso 1: Enviar comprobante para validación
	respRecepcion, err := c.EnviarComprobante(xmlComprobante)
	if err != nil {
		return nil, fmt.Errorf("error en recepción: %v", err)
	}
	
	// Verificar estado de recepción
	if respRecepcion.Estado != "RECIBIDA" {
		return nil, fmt.Errorf("comprobante no fue recibido por SRI. Estado: %s", respRecepcion.Estado)
	}
	
	fmt.Println("✅ Comprobante recibido por SRI")
	
	// Paso 2: Esperar un momento antes de consultar autorización
	fmt.Println("⏳ Esperando procesamiento del SRI...")
	time.Sleep(3 * time.Second)
	
	// Paso 3: Consultar autorización con reintentos
	var respAutorizacion *RespuestaComprobante
	maxReintentos := 5
	
	for intento := 1; intento <= maxReintentos; intento++ {
		fmt.Printf("🔍 Consultando autorización (intento %d/%d)...\n", intento, maxReintentos)
		
		respAutorizacion, err = c.ConsultarAutorizacion(claveAcceso)
		if err != nil {
			if intento == maxReintentos {
				return nil, fmt.Errorf("error consultando autorización después de %d intentos: %v", maxReintentos, err)
			}
			fmt.Printf("⚠️  Error en intento %d, reintentando...\n", intento)
			time.Sleep(2 * time.Second)
			continue
		}
		
		// Si hay autorizaciones, verificar estado
		if len(respAutorizacion.Autorizaciones) > 0 {
			autorizacion := respAutorizacion.Autorizaciones[0]
			
			if autorizacion.Estado == "AUTORIZADO" {
				fmt.Println("🎉 Comprobante AUTORIZADO por el SRI!")
				return &autorizacion, nil
			} else if autorizacion.Estado == "NO_AUTORIZADO" {
				return nil, fmt.Errorf("comprobante NO AUTORIZADO por el SRI")
			}
		}
		
		// Si aún no está procesado, esperar y reintentar
		if intento < maxReintentos {
			fmt.Println("⏳ Comprobante aún en procesamiento, esperando...")
			time.Sleep(3 * time.Second)
		}
	}
	
	return nil, fmt.Errorf("comprobante no fue autorizado después de %d intentos", maxReintentos)
}

// MostrarRespuestaRecepcion muestra información de respuesta de recepción
func MostrarRespuestaRecepcion(respuesta *RespuestaSolicitud) {
	fmt.Println("\n📥 RESPUESTA DE RECEPCIÓN SRI")
	fmt.Println("=============================")
	fmt.Printf("📊 Estado: %s\n", respuesta.Estado)
	
	for i, comp := range respuesta.Comprobantes {
		fmt.Printf("\n📋 Comprobante %d:\n", i+1)
		fmt.Printf("🔑 Clave Acceso: %s\n", comp.ClaveAcceso)
		
		if len(comp.Mensajes) > 0 {
			fmt.Println("💬 Mensajes:")
			for _, msg := range comp.Mensajes {
				fmt.Printf("  - %s: %s\n", msg.Tipo, msg.Mensaje)
				if msg.InformacionAdicional != "" {
					fmt.Printf("    Info adicional: %s\n", msg.InformacionAdicional)
				}
			}
		}
	}
}

// MostrarRespuestaAutorizacion muestra información de respuesta de autorización
func MostrarRespuestaAutorizacion(respuesta *RespuestaComprobante) {
	fmt.Println("\n🔐 RESPUESTA DE AUTORIZACIÓN SRI")
	fmt.Println("================================")
	fmt.Printf("🔑 Clave Consultada: %s\n", respuesta.ClaveAccesoConsultada)
	fmt.Printf("📊 Número de Comprobantes: %s\n", respuesta.NumeroComprobantes)
	
	for i, auth := range respuesta.Autorizaciones {
		fmt.Printf("\n📋 Autorización %d:\n", i+1)
		fmt.Printf("✅ Estado: %s\n", auth.Estado)
		fmt.Printf("🔢 Número de Autorización: %s\n", auth.NumeroAutorizacion)
		fmt.Printf("📅 Fecha de Autorización: %s\n", auth.FechaAutorizacion)
		fmt.Printf("🌍 Ambiente: %s\n", auth.Ambiente)
		
		if len(auth.Mensajes) > 0 {
			fmt.Println("💬 Mensajes:")
			for _, msg := range auth.Mensajes {
				fmt.Printf("  - %s: %s\n", msg.Tipo, msg.Mensaje)
			}
		}
	}
}