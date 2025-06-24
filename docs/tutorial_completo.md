# Tutorial Completo: Sistema de Facturación Electrónica SRI Ecuador

## 🎯 Guía Paso a Paso para Probar Todo el Sistema

Este tutorial te guía a través de **todas las funcionalidades** del sistema de facturación electrónica, desde lo más básico hasta la integración completa con SRI.

## 🚀 Inicio Rápido

### Prerrequisitos
```bash
# Verificar que Go está instalado
go version

# Debe mostrar Go 1.21+ 
# Si no tienes Go: https://golang.org/dl/
```

### Clonar y Preparar
```bash
# En tu directorio de proyectos
cd /ruta/a/tus/proyectos
git clone <este-repositorio>
cd go-facturacion-sri

# Instalar dependencias
go mod tidy
```

## 📚 Demos Disponibles

### 1️⃣ Demo Básico - Primer Contacto
```bash
# Ejecutar demo básico
go run main.go test_validaciones.go

# ¿Qué hace?
✅ Crea facturas automáticamente
✅ Valida cédulas ecuatorianas
✅ Calcula IVA (15%)
✅ Genera XML compatible SRI
✅ Muestra resumen en pantalla
```

**Salida esperada:**
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

### 2️⃣ Demo SRI - Integración Básica
```bash
# Demo de integración SRI
go run main.go test_validaciones.go sri

# ¿Qué hace?
✅ Genera claves de acceso de 49 dígitos
✅ Muestra información SRI detallada
✅ Simula proceso de autorización
✅ Explica tipos de comprobantes
```

**Funcionalidades mostradas:**
- Generación de claves de acceso válidas
- Información detallada de cada componente
- Simulación de autorización SRI
- Tipos de comprobantes soportados

### 3️⃣ Demo SOAP - Cliente Avanzado
```bash
# Demo del cliente SOAP
go run main.go test_validaciones.go soap

# ¿Qué hace?
✅ Muestra endpoints reales del SRI
✅ Simula flujo completo de comunicación
✅ Explica proceso paso a paso
✅ Demuestra estructura SOAP
```

**Puntos clave:**
- Endpoints oficiales de certificación y producción
- Simulación del flujo completo
- Estructura SOAP preparada para SRI real

### 4️⃣ Demo Base de Datos - Persistencia
```bash
# Demo del sistema de base de datos
go run main.go test_validaciones.go database

# ¿Qué hace?
✅ Crea facturas y las guarda en SQLite
✅ Asigna números secuenciales (FAC-000001)
✅ Gestiona estados (BORRADOR → AUTORIZADA)
✅ Muestra estadísticas del sistema
✅ Demuestra gestión de clientes
```

**Funcionalidades clave:**
- Persistencia completa de facturas
- Numeración automática secuencial
- Estados de factura manejados
- Base de datos SQLite integrada

### 5️⃣ Tests de Integración SRI
```bash
# Tests completos de integración
go run main.go test_validaciones.go test-sri

# ¿Qué hace?
✅ Verifica conectividad con endpoints SRI
✅ Valida estructura XML generada
✅ Prueba claves de acceso
✅ Simula proceso completo sin certificado
```

### 6️⃣ Guía de Certificación
```bash
# Guía completa de certificación
go run main.go test_validaciones.go certificacion

# ¿Qué muestra?
✅ Entidades certificadoras autorizadas
✅ Requisitos para personas naturales/jurídicas
✅ Costos y tiempos de cada opción
✅ Checklist para producción
```

## 🌐 API REST - Servidor Web

### Iniciar el Servidor
```bash
# Terminal 1: Iniciar servidor API
go run main.go test_validaciones.go api

# Salida esperada:
🚀 Servidor iniciado en http://localhost:8080
📋 Health check: http://localhost:8080/health
📊 API docs: http://localhost:8080/
```

### Probar Endpoints (Terminal 2)

#### Health Check
```bash
curl http://localhost:8080/health

# Respuesta:
{
  "status": "OK",
  "timestamp": "2025-06-23T15:30:00Z",
  "version": "1.0.0"
}
```

#### Crear Factura (En Memoria)
```bash
curl -X POST http://localhost:8080/api/facturas \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "EMPRESA DEMO S.A.",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "DEMO001",
        "descripcion": "Producto Demo API",
        "cantidad": 2,
        "precioUnitario": 150.00
      }
    ],
    "includeXML": true
  }'
```

#### Crear Factura (Base de Datos)
```bash
curl -X POST http://localhost:8080/api/facturas/db \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "CLIENTE BD DEMO",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "BD001",
        "descripcion": "Producto Base Datos",
        "cantidad": 1,
        "precioUnitario": 200.00
      }
    ]
  }'
```

#### Listar Facturas
```bash
# Facturas en base de datos (paginadas)
curl "http://localhost:8080/api/facturas/db/list?limit=5&offset=0"

# Facturas en memoria
curl http://localhost:8080/api/facturas
```

#### Obtener Factura Específica
```bash
# Obtener factura por ID con XML incluido
curl "http://localhost:8080/api/facturas/db/1?includeXML=true"
```

#### Actualizar Estado de Factura
```bash
curl -X PUT http://localhost:8080/api/facturas/db/1/estado \
  -H "Content-Type: application/json" \
  -d '{
    "estado": "AUTORIZADA",
    "numero_autorizacion": "1234567890123456789",
    "observaciones_sri": "Autorizada correctamente"
  }'
```

#### Ver Estadísticas
```bash
curl http://localhost:8080/api/estadisticas

# Respuesta:
{
  "success": true,
  "data": {
    "total_facturas": 5,
    "total_facturado": 1250.50,
    "por_estado": {
      "BORRADOR": 2,
      "AUTORIZADA": 3
    }
  }
}
```

#### Gestionar Clientes
```bash
# Crear cliente
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cedula": "1713175071",
    "nombre": "JUAN PEREZ DEMO",
    "direccion": "Av. Principal 123",
    "telefono": "0987654321",
    "email": "juan@demo.com",
    "tipoCliente": "PERSONA_NATURAL"
  }'

# Buscar cliente
curl "http://localhost:8080/api/clientes/buscar?cedula=1713175071"
```

## 🧪 Ejecutar Tests

### Tests Completos
```bash
# Ejecutar todos los tests
go test ./... -v

# Solo tests del módulo SRI
go test ./sri -v

# Solo tests de base de datos
go test ./database -v

# Tests con coverage
go test ./... -cover
```

### Benchmarks
```bash
# Benchmarks de performance
go test ./sri -bench=.
go test ./database -bench=.
```

## 📊 Flujo Completo del Sistema

### Proceso de Facturación Completo
```
1. USUARIO INGRESA DATOS
   ↓
2. VALIDACIÓN AUTOMÁTICA
   - Cédula ecuatoriana válida
   - Productos con datos completos
   ↓
3. CÁLCULOS AUTOMÁTICOS
   - Subtotales por producto
   - IVA 15% automático
   - Total general
   ↓
4. GENERACIÓN XML SRI
   - Estructura compatible SRI v2.31
   - Codificación UTF-8
   - Todos los campos requeridos
   ↓
5. CLAVE DE ACCESO
   - 49 dígitos únicos
   - Algoritmo módulo 11
   - Fecha, RUC, secuencial incluidos
   ↓
6. PERSISTENCIA (OPCIONAL)
   - Guarda en base de datos SQLite
   - Numeración secuencial automática
   - Estados manejados
   ↓
7. COMUNICACIÓN SRI (CON CERTIFICADO)
   - Cliente SOAP preparado
   - Reintentos automáticos
   - Manejo avanzado de errores
   ↓
8. AUTORIZACIÓN SRI
   - Número de autorización
   - XML autorizado
   - Validez fiscal
```

## 🎯 Casos de Uso Prácticos

### Caso 1: Desarrollador Probando
```bash
# 1. Ver demo básico
go run main.go test_validaciones.go

# 2. Entender SRI
go run main.go test_validaciones.go sri

# 3. Probar API
go run main.go test_validaciones.go api
# En otra terminal: curl tests...

# 4. Ver persistencia
go run main.go test_validaciones.go database
```

### Caso 2: Empresa Evaluando
```bash
# 1. Ver capacidades completas
go run main.go test_validaciones.go soap

# 2. Verificar base de datos
go run main.go test_validaciones.go database

# 3. Revisar certificación
go run main.go test_validaciones.go certificacion

# 4. API para integración
go run main.go test_validaciones.go api
```

### Caso 3: Implementación Real
```bash
# 1. Tests de integración
go run main.go test_validaciones.go test-sri

# 2. Obtener certificado ($24.64)
# Seguir guía: docs/certificado_online_guide.md

# 3. Configurar certificado real
# Editar ConfigTestSRI con datos reales

# 4. Tests reales con SRI
go run main.go test_validaciones.go test-sri
```

## 📁 Estructura de Archivos Generados

Después de ejecutar los demos, encontrarás:

```
proyecto/
├── database/
│   ├── demo_facturacion.db     # Base de datos demo
│   └── facturacion.db          # Base de datos API
├── docs/
│   ├── certificado_online_guide.md
│   └── tutorial_completo.md
└── logs/                       # Logs del sistema (si aplica)
```

## 🔧 Configuración Avanzada

### Variables de Ambiente
```bash
# Puerto personalizado para API
go run main.go test_validaciones.go api 3000

# Base de datos personalizada
export DB_PATH="/ruta/personalizada/facturas.db"
```

### Configuración SRI Real
```go
// En sri/test_integration_real.go
config := ConfigTestSRI{
    RutaCertificado: "/ruta/real/certificado.p12",
    PasswordCertificado: "password_real",
    ValidarCertificado: true,
    UsarAmbientePruebas: true, // false para producción
}
```

## 🚨 Solución de Problemas

### Error: "go: command not found"
```bash
# Instalar Go desde https://golang.org/dl/
# Verificar PATH incluye Go
export PATH=$PATH:/usr/local/go/bin
```

### Error: "port already in use"
```bash
# Usar puerto diferente
go run main.go test_validaciones.go api 8081
```

### Error: "database locked"
```bash
# Cerrar conexiones existentes o reiniciar
rm database/*.db
go run main.go test_validaciones.go database
```

## 🎉 ¡Listo para Producción!

Una vez que hayas probado todo:

1. ✅ **Sistema funcionando** localmente
2. ✅ **API REST** respondiendo
3. ✅ **Base de datos** operativa
4. ✅ **Certificado digital** obtenido ($24.64)
5. ✅ **Tests SRI** exitosos

**¡Tu sistema está listo para facturación electrónica real en Ecuador!**

---

**💡 Tip:** Guarda este tutorial para referencia futura y compártelo con tu equipo de desarrollo.