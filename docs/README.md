# Sistema de Facturación Electrónica Ecuador - Go

Sistema de facturación electrónica compatible con el SRI (Servicio de Rentas Internas) de Ecuador, desarrollado en Go.

## 🚀 Estado del Proyecto

**Avanzado** - Sistema empresarial con integración SRI ✅
- [x] Estructuras de datos y XML compatible con SRI
- [x] API REST completa con endpoints JSON
- [x] Sistema de paquetes organizados profesionalmente
- [x] 🆕 **Integración completa con SRI Ecuador**
- [x] 🆕 **Generación de claves de acceso de 49 dígitos**
- [x] 🆕 **Sistema de certificados digitales PKCS#12**
- [x] 🆕 **Firma digital XAdES-BES**
- [x] 🆕 **Simulación de autorización SRI**

## 📋 Características Actuales

### ✅ Core del Sistema
- **API REST completa** con endpoints JSON (POST/GET facturas, health check)
- **Generación de XML** conforme a especificaciones SRI v2.31
- **Cálculo automático** de impuestos (IVA 15%) y totales
- **Validación robusta** de cédulas ecuatorianas (algoritmo oficial)
- **Sistema de paquetes** profesional (models/, factory/, validators/, config/, api/, sri/)

### ✅ Integración SRI Ecuador
- **Claves de acceso de 49 dígitos** con algoritmo módulo 11 oficial
- **Certificados digitales PKCS#12** del Banco Central del Ecuador
- **Firma digital XAdES-BES** según especificaciones técnicas SRI
- **Tipos de comprobantes** completos (Facturas, Notas, Guías, Retenciones)
- **Ambientes duales** (Pruebas y Producción)
- **Autorización simulada** para testing

### ✅ Testing y Calidad
- **21 tests unitarios** con cobertura del 45.5%
- **Tests de integración** SRI específicos
- **Benchmarks** de performance
- **Validation** automática de documentos
- **Demo interactivo** con ejemplos reales

### 🔄 Próximas Mejoras
- Cliente SOAP para comunicación real con SRI
- Interfaz web para gestión visual
- Base de datos para persistencia
- Reportes y consultas avanzadas

## 🛠️ Tecnologías

- **Go 1.24+** - Lenguaje principal
- **encoding/xml** - Generación de XML nativo
- **Certificados PKCS#12** - Para firma digital (próximo)
- **SOAP** - Comunicación con SRI (próximo)

## 📖 Instalación y Uso

### Prerrequisitos
- **Go 1.21+** - [Descargar aquí](https://golang.org/dl/)
- **Git** - Para clonar el repositorio
- **Editor** - GoLand, VS Code, o tu preferido

### Instalación Rápida
```bash
# Clonar el proyecto
git clone <tu-repositorio>/go-facturacion-sri
cd go-facturacion-sri

# Instalar dependencias
go mod tidy

# Ejecutar en modo demo
go run main.go test_validaciones.go
```

### Modos de Ejecución

```bash
# 🧪 Modo Demo - Ejemplos básicos de facturación
go run main.go test_validaciones.go

# 🇪🇨 Modo SRI - Demo completo de integración SRI
go run main.go test_validaciones.go sri

# 🌐 Modo API - Servidor REST en puerto 8080
go run main.go test_validaciones.go api

# 🌐 API en puerto personalizado
go run main.go test_validaciones.go api 3000
```

### Usar la API REST

```bash
# Health check
curl http://localhost:8080/health

# Crear factura
curl -X POST http://localhost:8080/api/facturas \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "Juan Perez",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "LAPTOP001",
        "descripcion": "Laptop HP",
        "cantidad": 1,
        "precioUnitario": 800.00
      }
    ],
    "includeXML": true
  }'

# Listar facturas
curl http://localhost:8080/api/facturas

# Obtener factura específica con XML
curl "http://localhost:8080/api/facturas/FAC-000001?includeXML=true"
```

### Uso Programático

```go
package main

import (
    "go-facturacion-sri/factory"
    "go-facturacion-sri/models"
    "go-facturacion-sri/sri"
)

func main() {
    // Crear factura
    facturaData := models.FacturaInput{
        ClienteNombre: "JUAN CARLOS PEREZ",
        ClienteCedula: "1713175071", // Cédula válida
        Productos: []models.ProductoInput{
            {
                Codigo:         "LAPTOP001",
                Descripcion:    "Laptop Dell Inspiron 15",
                Cantidad:       1.0,
                PrecioUnitario: 450.00,
            },
        },
    }

    // Generar factura usando factory
    factura, err := factory.CrearFactura(facturaData)
    if err != nil {
        panic(err)
    }

    // Generar XML
    xmlData, _ := factura.GenerarXML()
    
    // Generar clave de acceso SRI
    claveConfig := sri.ClaveAccesoConfig{
        FechaEmision:     time.Now(),
        TipoComprobante:  sri.Factura,
        RUCEmisor:        "1792146739001",
        Ambiente:         sri.Pruebas,
        Serie:            "001001",
        NumeroSecuencial: "000000001",
        TipoEmision:      sri.EmisionNormal,
    }
    
    claveAcceso, _ := sri.GenerarClaveAcceso(claveConfig)
    
    // Mostrar resultados
    factura.MostrarResumen()
    sri.MostrarInformacionClaveAcceso(claveAcceso)
}
```

## 📊 Ejemplo de Salida

### Demo Básico
```
🚀 GENERANDO FACTURA PRINCIPAL
==================================================
✅ Factura creada exitosamente
📋 Cliente: JUAN CARLOS PEREZ (1713175071)
💰 Subtotal: $1010.00
🧮 IVA 15%: $151.50
💵 TOTAL: $1161.50
📑 Productos: 3 items
```

### Demo SRI Completo
```
🇪🇨 DEMO SISTEMA INTEGRACIÓN SRI ECUADOR
============================================================

1️⃣ GENERACIÓN DE CLAVE DE ACCESO
✅ Clave de acceso generada: 23062025-01-1792146739001-1-001001-000000001-91527284-1-1

🔑 INFORMACIÓN DE CLAVE DE ACCESO
=================================
🎯 Clave de Acceso: 23062025-01-1792146739001-1-001001-000000001-91527284-1-1
📅 Fecha de Emisión: 23/06/2025
📋 Tipo Comprobante: Factura (01)
🏢 RUC Emisor: 1792146739001
🌍 Ambiente: Pruebas (1)
✅ Validación: Clave de acceso válida

2️⃣ SIMULACIÓN AUTORIZACIÓN SRI
📝 Número de Autorización: 2306202501179214673900110010010000000019152728411
📅 Fecha de Autorización: 23/06/2025 14:14:33
✅ Estado: AUTORIZADO
```

### API Response
```json
{
  "id": "FAC-000001",
  "status": "created",
  "factura": {
    "InfoFactura": {
      "ImporteTotal": "920.00"
    }
  },
  "xml": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<factura>...</factura>",
  "createdAt": "2025-06-23T14:15:33Z"
}
```

## 🏗️ Arquitectura del Sistema

### Estructura de Paquetes
```
go-facturacion-sri/
├── 📁 api/          # Servidor HTTP REST
│   ├── server.go    # Configuración del servidor
│   ├── handlers.go  # Manejadores de endpoints
│   └── middleware.go # CORS, logging, etc.
│
├── 📁 models/       # Estructuras de datos
│   └── factura.go   # Tipos y métodos XML
│
├── 📁 factory/      # Lógica de creación
│   └── factura_factory.go # Factory functions
│
├── 📁 validators/   # Validaciones de negocio
│   └── validations.go # Cédulas, productos, etc.
│
├── 📁 config/       # Configuración externa
│   ├── config.go    # Estructuras de config
│   ├── loader.go    # Carga desde JSON
│   ├── desarrollo.json # Ambiente de pruebas
│   └── produccion.json # Ambiente productivo
│
├── 📁 sri/          # 🆕 Integración SRI Ecuador
│   ├── certificado.go    # Certificados PKCS#12
│   ├── xades_bes.go     # Firma digital XAdES-BES
│   ├── autorizacion.go  # Claves de acceso y autorización
│   ├── demo.go          # Demostraciones interactivas
│   └── integration_test.go # Tests de integración
│
└── 📄 main.go       # Punto de entrada
```

### Flujo Completo de Facturación Electrónica
```
1. 📝 INPUT
   ├── Cliente (nombre, cédula)
   └── Productos (código, descripción, cantidad, precio)

2. 🏭 FACTORY
   ├── Validar datos (cédula ecuatoriana, productos)
   ├── Calcular impuestos (IVA 15%)
   ├── Generar secuencial
   └── Construir estructura XML

3. 🇪🇨 SRI INTEGRATION
   ├── Generar clave de acceso (49 dígitos)
   ├── Firmar con certificado PKCS#12
   ├── Aplicar XAdES-BES
   └── Preparar para envío SRI

4. 🌐 API REST
   ├── Recibir vía HTTP POST
   ├── Procesar en background
   ├── Almacenar temporalmente
   └── Retornar JSON + XML

5. 📋 OUTPUT
   ├── XML firmado SRI-compliant
   ├── Clave de acceso válida
   ├── JSON response con metadata
   └── Logs de auditoría
```

### Conceptos Clave SRI

```go
// Clave de Acceso (49 dígitos)
type ClaveAccesoConfig struct {
    FechaEmision     time.Time      // ddMMyyyy
    TipoComprobante  TipoComprobante // 01=Factura, 04=NotaCredito
    RUCEmisor        string         // 13 dígitos
    Ambiente         Ambiente       // 1=Pruebas, 2=Producción
    Serie            string         // EstabPuntoEmi (6 dígitos)
    NumeroSecuencial string         // 9 dígitos incrementales
    CodigoNumerico   string         // 8 dígitos aleatorios
    TipoEmision      TipoEmision    // 1=Normal, 2=Contingencia
    // + 1 dígito verificador (módulo 11)
}

// Certificado Digital
type CertificadoDigital struct {
    Archivo    string           // Ruta al .p12
    Password   string           // Contraseña
    PrivateKey interface{}      // Clave privada RSA
    Cert       *x509.Certificate // Certificado X.509
    CACerts    []*x509.Certificate // Cadena de CA
}
```

## 🎯 Roadmap de Desarrollo

### Semana 1-2: Fundamentos ✅
- [x] Go básico y estructuras XML
- [x] Conceptos de facturación electrónica

### Semana 3-4: Core del Sistema
- [ ] Validaciones avanzadas (RUC, secuencias)
- [ ] Cálculos tributarios completos
- [ ] Generación de clave de acceso real

### Semana 5-6: Firma Digital y SRI
- [ ] Certificados PKCS#12
- [ ] Firma XAdES-BES
- [ ] Cliente SOAP para comunicación SRI

### Semana 7-8: Interfaz y Deployment
- [ ] Frontend web básico
- [ ] CRUD de clientes y productos
- [ ] Binario único para distribución

### Semana 9-12: Funcionalidades Avanzadas
- [ ] Notas de crédito/débito
- [ ] Guías de remisión
- [ ] Comprobantes de retención
- [ ] Reportes y consultas

## 🔧 Configuración del SRI

### Endpoints
- **Recepción**: `https://celcer.sri.gob.ec/comprobantes-electronicos-ws/` (Certificación)
- **Autorización**: `https://celcer.sri.gob.ec/comprobantes-electronicos-ws/` (Certificación)
- **Producción**: `https://cel.sri.gob.ec/comprobantes-electronicos-ws/`

### Requisitos Técnicos
- Certificado digital vigente (BCE, Security Data, ANF)
- Formato XML con codificación UTF-8
- Firma digital XAdES-BES
- Numeración secuencial obligatoria

## 📈 Ventajas de Go para Facturación Electrónica

- **Performance**: 3-40x más rápido que Java/.NET
- **Concurrencia**: Procesamiento paralelo de documentos
- **Deployment**: Binario único sin dependencias
- **XML nativo**: Sin librerías externas necesarias
- **Memoria**: 60% menos recursos que alternativas

## 🤝 Contribuir

1. Fork el proyecto
2. Crear rama para feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE.md](LICENSE.md) para detalles.

## 📞 Contacto

- **Desarrollador**: Juanqui
- **Email**: [tu-email@ejemplo.com]
- **LinkedIn**: [tu-linkedin]

## 🙏 Reconocimientos

- **SRI Ecuador** - Especificaciones técnicas oficiales
- **Go Team** - Por el excelente ecosistema de desarrollo
- **Comunidad Go Ecuador** - Soporte y feedback

---

**Nota**: Este sistema está en desarrollo activo. Para uso en producción, completar las fases de firma digital y certificación SRI.