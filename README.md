# Sistema de Facturación Electrónica Ecuador - Go

Sistema de facturación electrónica compatible con el SRI (Servicio de Rentas Internas) de Ecuador, desarrollado en Go.

## 🚀 Estado del Proyecto

**Semana 1/12** - Fundamentos completados ✅
- [x] Estructuras básicas de datos (InfoTributaria, InfoFactura, Detalle)
- [x] Generación de XML compatible con SRI
- [x] Factory functions para crear facturas
- [x] Cálculo automático de IVA (15%)
- [x] Numeración secuencial automática

## 📋 Características Actuales

### ✅ Implementado
- Generación de facturas electrónicas en formato XML
- Estructuras de datos conformes a esquemas XSD del SRI
- Cálculo automático de impuestos (IVA 15%)
- Numeración secuencial de documentos
- Validación básica de datos de entrada
- Soporte para ambiente de pruebas

### 🔄 En Desarrollo (Próximas Semanas)
- Generación de clave de acceso de 49 dígitos (Semana 4)
- Firma digital XAdES-BES (Semana 5-6)
- Comunicación SOAP con servicios SRI (Semana 5-6)
- Interfaz web básica (Semana 7-8)
- Múltiples tipos de documentos (Semana 9-12)

## 🛠️ Tecnologías

- **Go 1.24+** - Lenguaje principal
- **encoding/xml** - Generación de XML nativo
- **Certificados PKCS#12** - Para firma digital (próximo)
- **SOAP** - Comunicación con SRI (próximo)

## 📖 Instalación y Uso

### Prerrequisitos
- Go 1.21 o superior
- GoLand (recomendado) o VS Code

### Instalación
```bash
# Clonar o crear el proyecto
mkdir go-facturacion-sri
cd go-facturacion-sri
go mod init facturacion-sri

# Copiar el código main.go
# Ejecutar
go run main.go
```

### Uso Básico
```go
// Crear datos de factura
facturaData := FacturaInput{
    ClienteNombre:       "JUAN CARLOS PEREZ",
    ClienteCedula:       "1234567890",
    ProductoCodigo:      "LAPTOP001",
    ProductoDescripcion: "Laptop Dell Inspiron 15",
    Cantidad:            2.0,
    PrecioUnitario:      450.00,
}

// Generar factura
factura := CrearFactura(facturaData)
xmlFactura, err := factura.GenerarXML()
```

## 📊 Ejemplo de Salida

```
=== FACTURA ELECTRÓNICA ECUATORIANA ===
Secuencial: 000000001
Cliente: JUAN CARLOS PEREZ (1234567890)
Producto: Laptop Dell Inspiron 15
Cantidad: 2 x $450.00 = $900.00
IVA 15%: $135.00
TOTAL: $1035.00
```

## 🏗️ Arquitectura del Sistema

### Estructuras Principales

```go
type InfoTributaria struct {
    Ambiente        string // "1"=pruebas, "2"=producción
    TipoEmision     string // "1"=normal, "2"=contingencia
    RazonSocial     string // Nombre de la empresa
    RUC             string // RUC del emisor (13 dígitos)
    ClaveAcceso     string // 49 dígitos únicos por documento
    CodDoc          string // "01"=factura, "04"=nota crédito
    Establecimiento string // 3 dígitos (001, 002, etc.)
    PuntoEmision    string // 3 dígitos (001, 002, etc.)
    Secuencial      string // 9 dígitos incrementales
}
```

### Flujo de Generación
1. **Input** → Datos simples del cliente y producto
2. **Factory** → `CrearFactura()` construye estructura completa
3. **Cálculos** → IVA y totales automáticos
4. **XML** → Generación conforme a esquemas SRI
5. **Output** → XML listo para firma digital

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