package sri

import (
	"testing"
	"time"
)

func TestNewSOAPClient(t *testing.T) {
	tests := []struct {
		name     string
		ambiente Ambiente
	}{
		{
			name:     "Cliente para ambiente de pruebas",
			ambiente: Pruebas,
		},
		{
			name:     "Cliente para ambiente de producción",
			ambiente: Produccion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewSOAPClient(tt.ambiente)
			
			if client == nil {
				t.Error("NewSOAPClient() devolvió nil")
				return
			}
			
			if client.Ambiente != tt.ambiente {
				t.Errorf("NewSOAPClient() ambiente = %v, esperado %v", client.Ambiente, tt.ambiente)
			}
			
			if client.TimeoutSegundos != 30 {
				t.Errorf("NewSOAPClient() timeout = %v, esperado 30", client.TimeoutSegundos)
			}
			
			if client.httpClient == nil {
				t.Error("NewSOAPClient() httpClient es nil")
			}
		})
	}
}

func TestSOAPClientEndpoints(t *testing.T) {
	tests := []struct {
		name                 string
		ambiente            Ambiente
		expectedRecepcion   string
		expectedAutorizacion string
	}{
		{
			name:                 "Endpoints de certificación",
			ambiente:            Pruebas,
			expectedRecepcion:   EndpointRecepcionCertificacion,
			expectedAutorizacion: EndpointAutorizacionCertificacion,
		},
		{
			name:                 "Endpoints de producción",
			ambiente:            Produccion,
			expectedRecepcion:   EndpointRecepcionProduccion,
			expectedAutorizacion: EndpointAutorizacionProduccion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verificar que las constantes están definidas correctamente
			if tt.ambiente == Pruebas {
				if EndpointRecepcionCertificacion == "" {
					t.Error("EndpointRecepcionCertificacion no está definido")
				}
				if EndpointAutorizacionCertificacion == "" {
					t.Error("EndpointAutorizacionCertificacion no está definido")
				}
			} else {
				if EndpointRecepcionProduccion == "" {
					t.Error("EndpointRecepcionProduccion no está definido")
				}
				if EndpointAutorizacionProduccion == "" {
					t.Error("EndpointAutorizacionProduccion no está definido")
				}
			}
		})
	}
}

// TestParsearRespuestaRecepcionMock simula el parsing de una respuesta de recepción
func TestParsearRespuestaRecepcionMock(t *testing.T) {
	client := NewSOAPClient(Pruebas)
	
	// XML de respuesta simulada del SRI
	respuestaMock := `<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
			<ns2:respuestaSolicitud xmlns:ns2="http://ec.gob.sri.ws.recepcion">
				<estado>RECIBIDA</estado>
				<comprobantes>
					<comprobante>
						<claveAcceso>2306202501179214673900110010010000000019152728411</claveAcceso>
						<mensajes>
							<mensaje>
								<identificador>CLAVE-01</identificador>
								<mensaje>CLAVE DE ACCESO REGISTRADA</mensaje>
								<informacionAdicional></informacionAdicional>
								<tipo>INFORMATIVO</tipo>
							</mensaje>
						</mensajes>
					</comprobante>
				</comprobantes>
			</ns2:respuestaSolicitud>
		</soap:Body>
	</soap:Envelope>`
	
	respuesta, err := client.parsearRespuestaRecepcion([]byte(respuestaMock))
	if err != nil {
		t.Fatalf("Error parseando respuesta mock: %v", err)
	}
	
	if respuesta.Estado != "RECIBIDA" {
		t.Errorf("Estado esperado: RECIBIDA, obtenido: %s", respuesta.Estado)
	}
	
	if len(respuesta.Comprobantes) != 1 {
		t.Errorf("Se esperaba 1 comprobante, se obtuvo: %d", len(respuesta.Comprobantes))
	}
	
	comp := respuesta.Comprobantes[0]
	expectedClave := "2306202501179214673900110010010000000019152728411"
	if comp.ClaveAcceso != expectedClave {
		t.Errorf("Clave de acceso esperada: %s, obtenida: %s", expectedClave, comp.ClaveAcceso)
	}
	
	if len(comp.Mensajes) != 1 {
		t.Errorf("Se esperaba 1 mensaje, se obtuvo: %d", len(comp.Mensajes))
	}
}

// TestParsearRespuestaAutorizacionMock simula el parsing de una respuesta de autorización
func TestParsearRespuestaAutorizacionMock(t *testing.T) {
	client := NewSOAPClient(Pruebas)
	
	// XML de respuesta simulada del SRI
	respuestaMock := `<?xml version="1.0" encoding="UTF-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
			<ns2:respuestaComprobante xmlns:ns2="http://ec.gob.sri.ws.autorizacion">
				<claveAccesoConsultada>2306202501179214673900110010010000000019152728411</claveAccesoConsultada>
				<numeroComprobantes>1</numeroComprobantes>
				<autorizaciones>
					<autorizacion>
						<estado>AUTORIZADO</estado>
						<numeroAutorizacion>2306202501179214673900110010010000000019152728411</numeroAutorizacion>
						<fechaAutorizacion>2025-06-23T14:30:00.000-05:00</fechaAutorizacion>
						<ambiente>PRUEBAS</ambiente>
						<comprobante><![CDATA[<?xml version="1.0" encoding="UTF-8"?><factura>...</factura>]]></comprobante>
						<mensajes></mensajes>
					</autorizacion>
				</autorizaciones>
			</ns2:respuestaComprobante>
		</soap:Body>
	</soap:Envelope>`
	
	respuesta, err := client.parsearRespuestaAutorizacion([]byte(respuestaMock))
	if err != nil {
		t.Fatalf("Error parseando respuesta mock: %v", err)
	}
	
	expectedClave := "2306202501179214673900110010010000000019152728411"
	if respuesta.ClaveAccesoConsultada != expectedClave {
		t.Errorf("Clave consultada esperada: %s, obtenida: %s", expectedClave, respuesta.ClaveAccesoConsultada)
	}
	
	if respuesta.NumeroComprobantes != "1" {
		t.Errorf("Número de comprobantes esperado: 1, obtenido: %s", respuesta.NumeroComprobantes)
	}
	
	if len(respuesta.Autorizaciones) != 1 {
		t.Errorf("Se esperaba 1 autorización, se obtuvo: %d", len(respuesta.Autorizaciones))
	}
	
	auth := respuesta.Autorizaciones[0]
	if auth.Estado != "AUTORIZADO" {
		t.Errorf("Estado esperado: AUTORIZADO, obtenido: %s", auth.Estado)
	}
	
	if auth.NumeroAutorizacion != expectedClave {
		t.Errorf("Número de autorización esperado: %s, obtenido: %s", expectedClave, auth.NumeroAutorizacion)
	}
}

// BenchmarkNewSOAPClient mide el performance de creación de cliente
func BenchmarkNewSOAPClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client := NewSOAPClient(Pruebas)
		_ = client
	}
}

// TestSOAPClientTimeout verifica la configuración de timeout
func TestSOAPClientTimeout(t *testing.T) {
	client := NewSOAPClient(Pruebas)
	
	expectedTimeout := 30 * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Timeout esperado: %v, obtenido: %v", expectedTimeout, client.httpClient.Timeout)
	}
}

// TestConstantesEndpoints verifica que los endpoints estén definidos correctamente
func TestConstantesEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		shouldContain string
	}{
		{
			name:          "Endpoint recepción certificación",
			endpoint:      EndpointRecepcionCertificacion,
			shouldContain: "celcer.sri.gob.ec",
		},
		{
			name:          "Endpoint autorización certificación",
			endpoint:      EndpointAutorizacionCertificacion,
			shouldContain: "celcer.sri.gob.ec",
		},
		{
			name:          "Endpoint recepción producción",
			endpoint:      EndpointRecepcionProduccion,
			shouldContain: "cel.sri.gob.ec",
		},
		{
			name:          "Endpoint autorización producción",
			endpoint:      EndpointAutorizacionProduccion,
			shouldContain: "cel.sri.gob.ec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.endpoint == "" {
				t.Errorf("Endpoint %s está vacío", tt.name)
			}
			
			if !contains(tt.endpoint, tt.shouldContain) {
				t.Errorf("Endpoint %s no contiene %s", tt.endpoint, tt.shouldContain)
			}
			
			if !contains(tt.endpoint, "https://") {
				t.Errorf("Endpoint %s no usa HTTPS", tt.endpoint)
			}
		})
	}
}

// contains verifica si una cadena contiene una subcadena
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      containsInner(s, substr))))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}