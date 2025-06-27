# CLAUDE.md

Este archivo proporciona orientación a Claude Code (claude.ai/code) al trabajar con código en este repositorio.

## Descripción del Proyecto

Sistema de facturación electrónica basado en Go para el SRI de Ecuador (Servicio de Rentas Internas). Sistema avanzado con integración completa al SRI, certificados digitales, firmas XAdES-BES, y API REST. Actualmente soporta procesamiento completo de facturas, generación de XML, y creación de documentos compatibles con el SRI.

## Comandos de Desarrollo

### Custom Slash Commands (Claude Code)

**Sistema de comandos inteligentes** para debugging y gestión rápida:

```bash
# Debugging y Diagnóstico
/debug:auth "login failing with 401"          # Análisis inteligente de autenticación
/debug:sri "certificate expired"              # Debugging específico SRI
/test:api "all endpoints"                      # Testing completo API

# Configuración y Setup
/setup:cert                                    # Guía interactiva certificados BCE
/db:query "SELECT COUNT(*) FROM facturas"     # Consultas con análisis automático
/deploy:check                                  # Verificación pre-deployment

# Ejemplos de uso específicos
/debug:auth "CORS error in frontend"          # → Analiza problemas CORS
/debug:sri "XML validation failed"            # → Debug validación XML SRI
/setup:cert "production environment"          # → Guía certificados producción
/db:query "performance analysis"              # → Análisis rendimiento BD
```

### Comandos Tradicionales de Desarrollo

```bash
# Ejecutar modos de aplicación
go run main.go test_validaciones.go           # Modo demo con ejemplos
go run main.go test_validaciones.go api       # Iniciar servidor API REST (puerto 8080)
go run main.go test_validaciones.go api 3000  # Iniciar API en puerto personalizado
go run main.go test_validaciones.go sri       # Demo de integración SRI
go run main.go test_validaciones.go soap      # Demo cliente SOAP
go run main.go test_validaciones.go database  # Demo base de datos

# Compilar ejecutable
go build -o facturacion-sri main.go test_validaciones.go

# Pruebas
go test ./...                    # Ejecutar todas las pruebas
go test -v ./...                # Salida verbosa de pruebas
go test -cover ./...            # Pruebas con cobertura
go test -bench=. ./...          # Ejecutar benchmarks
go test ./sri -v               # Probar paquete específico

# Pruebas de API
./test_api.sh                   # Ejecutar pruebas de integración API (requiere jq)
curl http://localhost:8080/health  # Verificación rápida de salud

# Calidad de código
go fmt ./...                    # Formatear código
go vet ./...                   # Análisis estático
go mod tidy                    # Limpiar dependencias
go mod verify                  # Verificar dependencias

# Análisis de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Arquitectura

### Estructura de Paquetes
```
go-facturacion-sri/
├── api/          # Servidor API REST y manejadores
├── config/       # Configuración JSON externa (desarrollo.json, produccion.json)
├── database/     # Integración base de datos SQLite
├── factory/      # Creación de facturas con lógica de negocio
├── models/       # Estructuras de datos principales y generación XML
├── sri/          # Integración SRI (certificados, firmas, autorización)
└── validators/   # Validaciones de negocio (cédulas, RUC, productos)
```

### Flujo Principal de Datos
1. **Entrada**: JSON vía API o programáticamente `models.FacturaInput`
2. **Validación**: Validación de cédulas, productos (`validators/`)
3. **Creación**: Patrón Factory con reglas de negocio (`factory/`)
4. **Procesamiento SRI**: Generación de clave de acceso, firmas digitales (`sri/`)  
5. **Salida**: Generación XML, respuestas API (`models/`, `api/`)

### Estructuras de Datos Clave

**Estructura Principal de Factura** (`models/factura.go`):
- `InfoTributaria`: Información tributaria (RUC, establecimiento, secuencias)
- `InfoFactura`: Metadatos de factura (fechas, cliente, totales)
- `Detalle`: Líneas de productos (productos con cantidades, precios, impuestos)
- `Factura`: Estructura completa de factura compatible con SRI

**Integración SRI** (paquete `sri/`):
- `ClaveAccesoConfig`: Generación de clave de acceso de 49 dígitos con validación módulo-11
- `CertificadoDigital`: Manejo de certificados PKCS#12 para firmas XAdES-BES
- `AutorizacionSRI`: Simulación de autorización SRI y manejo de respuestas

### Requisitos Específicos del SRI

**Códigos y Formatos de Documentos**:
- **Ambiente**: "1" (pruebas), "2" (producción)
- **ClaveAcceso**: Clave única de 49 dígitos con fecha, RUC, secuencia y dígito verificador
- **CodDoc**: "01" (factura), "04" (nota de crédito), "05" (nota de débito), "06" (guía de remisión)
- **Secuencial**: Numeración secuencial de 9 dígitos (000000001, 000000002, etc.)
- **Formato de Fecha**: DD/MM/YYYY como requieren los esquemas SRI

**Cálculos de Impuestos**:
- Cálculo automático del 15% de IVA para Ecuador
- Manejo adecuado de códigos de impuestos (tarifa "2" para IVA 15%)
- Cálculos de subtotal y total con manejo de precisión

### Sistema de Configuración

**Configuración JSON Externa** (paquete `config/`):
- `desarrollo.json`: Configuraciones de ambiente de pruebas
- `produccion.json`: Configuraciones de ambiente de producción  
- Carga automática con fallback a valores por defecto
- Códigos de RUC y establecimiento específicos por ambiente

### Arquitectura de API

**Endpoints REST** (paquete `api/`):
- `POST /api/facturas`: Crear factura con entrada JSON
- `GET /api/facturas`: Listar todas las facturas creadas
- `GET /api/facturas/{id}`: Obtener factura específica con XML opcional
- `GET /health`: Verificación de salud del sistema
- Middleware CORS y logging de peticiones

### Puntos de Integración SRI

**Endpoints**:
- Certificación: `https://celcer.sri.gob.ec/comprobantes-electronicos-ws/`
- Producción: `https://cel.sri.gob.ec/comprobantes-electronicos-ws/`

**Flujo de Firma Digital**:
1. Cargar certificado PKCS#12 del BCE/Security Data/ANF
2. Generar clave de acceso de 49 dígitos con validación módulo-11
3. Aplicar firma XAdES-BES al documento XML
4. Enviar al SRI para autorización

## Patrones Clave de Desarrollo

### Implementación del Patrón Factory
El sistema usa funciones factory para creación de objetos con validación adecuada:
```go
// factory/factura_factory.go
func CrearFactura(input models.FacturaInput) (*models.Factura, error) {
    // Validación -> Lógica de negocio -> Creación de objeto
}
```

### Convenciones de Manejo de Errores
- Siempre retornar `error` como último valor de retorno
- Usar mensajes de error descriptivos con contexto
- Validar entradas temprano y retornar errores inmediatamente
- Registrar errores pero no exponer detalles internos a consumidores de API

### Estructura de Pruebas
- Archivos `*_test.go` junto a archivos de implementación
- Funciones de prueba con prefijo `Test`
- Usar pruebas basadas en tablas para múltiples escenarios
- Pruebas de integración en `sri/integration_test.go`
- Reportes de cobertura con `coverage.out`

### Detalles de Implementación Específicos del SRI

**Generación de Clave de Acceso** (`sri/autorizacion.go`):
- Clave de 49 dígitos: `ddMMyyyyTTrrrrrrrrrrrraeeeeeeNNNNNNNNNccccccccee`
- Cálculo de dígito verificador módulo-11
- Implementación real con formato de fecha adecuado

**Generación de XML** (`models/factura.go`):
- Etiquetas struct para marshaling XML: `xml:"campo,attr"`
- Manejo adecuado de namespaces para esquemas SRI
- Formato de fecha en DD/MM/YYYY para cumplimiento SRI

**Lógica de Validación** (`validators/validations.go`):
- Validación de cédula ecuatoriana con algoritmo 10
- Validación de código y descripción de productos
- Validación de RUC para identificación empresarial

## Gestión de Configuración

El sistema carga configuración desde archivos JSON externos:
- `config/desarrollo.json`: Ambiente de desarrollo/pruebas
- `config/produccion.json`: Ambiente de producción
- Fallback automático a valores por defecto si falla la carga JSON
- Códigos de RUC, establecimiento y endpoints SRI específicos por ambiente

## Principios de Diseño de API

- Endpoints RESTful con métodos HTTP adecuados
- Entrada/salida JSON con inclusión opcional de XML
- Formato consistente de respuesta de errores
- Middleware CORS para peticiones cross-origin
- Logging de peticiones para debugging y auditoría

## Estrategia de Pruebas

- **Pruebas Unitarias**: Pruebas de funciones individuales (21 pruebas actualmente)
- **Pruebas de Integración**: Flujos de trabajo específicos del SRI
- **Pruebas de API**: Pruebas end-to-end HTTP request/response vía `test_api.sh`
- **Pruebas de Benchmark**: Medición de rendimiento para rutas críticas
- **Meta de Cobertura**: 45.5% actual, objetivo 80%+ para producción

## Flujo de Trabajo de Desarrollo

### Desarrollo Rápido con Custom Commands

1. **Análisis de Problemas**: `/debug:auth "problema específico"` o `/debug:sri "error SRI"`
2. **Testing Dirigido**: `/test:api "funcionalidad específica"`
3. **Configuración**: `/setup:cert` para certificados o `/db:query "SQL"`
4. **Verificación**: `/deploy:check` antes de deployment

### Flujo Tradicional

1. **Modificar Código**: Hacer cambios a archivos `.go`
2. **Ejecutar Pruebas**: `go test ./...` para asegurar que nada se rompa
3. **Formatear Código**: `go fmt ./...` para estilo consistente
4. **Probar API**: `./test_api.sh` para pruebas de integración
5. **Verificar Cobertura**: `go test -cover ./...` para monitorear cobertura de pruebas

## Custom Commands disponibles

### Debugging Inteligente
- `/debug:auth "mensaje"` - Análisis completo de problemas de autenticación
- `/debug:sri "error"` - Debugging específico de integración SRI

### Testing y Validación  
- `/test:api "scope"` - Testing dirigido de endpoints API
- `/deploy:check` - Verificación completa pre-deployment

### Configuración y Setup
- `/setup:cert` - Guía interactiva para certificados digitales BCE
- `/db:query "SQL"` - Ejecución y análisis de consultas de base de datos

Los custom commands proporcionan análisis contextual, sugerencias específicas y debugging automático adaptado al stack tecnológico del proyecto (Go + SRI + Astro).