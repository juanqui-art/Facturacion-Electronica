# Guía Completa: Certificado Digital Online BCE Ecuador 2024

## 🎉 Proceso 100% Virtual Disponible

El Banco Central del Ecuador (BCE) ahora ofrece certificados digitales **completamente online** a través de empresas autorizadas, sin necesidad de citas presenciales.

## 💰 Costos y Opciones

### Certificado de Archivo (Recomendado para desarrollo)
- **Costo:** $22.00 + IVA = **$24.64 USD**
- **Vigencia:** 2 años
- **Proceso:** 100% virtual
- **Validación:** Biométrica facial
- **Tiempo:** Inmediato

### Certificado en Token (Tradicional)
- **Costo:** $49.00 + IVA = **$54.88 USD**
- **Vigencia:** 2 años
- **Proceso:** Mixto (online + presencial)
- **Tiempo:** 3-5 días + cita

## 🚀 Proceso Paso a Paso - Certificado Online

### Paso 1: Acceder al Portal BCE
```
🌐 URL: https://www.eci.bce.ec/
📋 Ir a: "Solicitud de Certificado"
```

### Paso 2: Seleccionar Proveedor
Elegir uno de los proveedores autorizados:
- **Latinus S.A.**
- **Sodig S.A.**

### Paso 3: Completar Datos
- Número de cédula
- Datos personales
- Email y teléfono
- Tipo de procedimiento: "Emisión"

### Paso 4: Validación Biométrica
- Permitir acceso a cámara
- Fotografía de cédula (anverso y reverso)
- Fotografía de rostro para validación facial
- Sistema compara biométricamente

### Paso 5: Pago Online
- **Monto:** $24.64 USD
- **Métodos:** Tarjeta de crédito/débito
- **Confirmación:** Inmediata

### Paso 6: Descarga del Certificado
- **Formato:** Archivo .p12
- **Contraseña:** La que defines durante el proceso
- **Descarga:** Inmediata tras confirmación de pago

## 🔧 Integración con Sistema de Facturación

### Configuración del Certificado
```go
// En tu sistema Go
config := ConfigTestSRI{
    RutaCertificado: "/ruta/al/certificado.p12",
    PasswordCertificado: "tu_password_elegido",
    ValidarCertificado: true,
    UsarAmbientePruebas: true, // Empezar en pruebas
}
```

### Activar Comunicación Real con SRI
```bash
# Ejecutar tests con certificado real
go run main.go test_validaciones.go test-sri

# Si todo funciona, cambiar a producción
config.UsarAmbientePruebas = false
```

## 📋 Requisitos para el Proceso Online

### Documentos Necesarios
- ✅ Cédula de identidad ecuatoriana vigente
- ✅ Dispositivo con cámara (PC, laptop, móvil)
- ✅ Conexión a internet estable
- ✅ Tarjeta de crédito/débito para pago

### Datos que Necesitarás
- Número de cédula
- Nombre completo (como aparece en cédula)
- Email válido
- Número de teléfono
- Dirección actual

## ⚠️ Consideraciones Importantes

### Validez del Certificado
- ✅ **Válido para SRI:** Ambiente de pruebas y producción
- ✅ **Válido para:** Facturación electrónica, firma de documentos
- ✅ **Reconocido por:** Todas las instituciones públicas ecuatorianas

### Limitaciones
- 🔒 **Solo para personas naturales** con cédula ecuatoriana
- 📱 **Requiere biometría facial** (cámara funcional)
- 🌐 **Internet estable** durante el proceso

### Diferencias vs. Token
| Característica | Archivo (.p12) | Token Físico |
|---------------|----------------|--------------|
| Costo | $24.64 | $54.88 |
| Proceso | 100% online | Mixto |
| Tiempo | Inmediato | 3-5 días |
| Portabilidad | Software | Hardware |
| Backup | Fácil | Limitado |
| Pérdida | Recuperable | Problemático |

## 🎯 Para Desarrolladores

### Ventajas del Certificado de Archivo
1. **Desarrollo ágil:** Obtienes certificado en minutos
2. **Costo bajo:** Inversión mínima para pruebas
3. **Flexibilidad:** Fácil de integrar y respaldar
4. **Testing:** Perfecto para ambiente de desarrollo

### Integración Inmediata
```bash
# Una vez obtenido el certificado:
# 1. Guardar archivo .p12 en proyecto
# 2. Configurar ruta y contraseña
# 3. Ejecutar tests reales
go run main.go test_validaciones.go test-sri

# 4. ¡Sistema funcionando con SRI real!
```

## 🏆 Resultado Final

Con esta inversión de **$24.64 USD** y **30 minutos de tiempo**, tu sistema de facturación pasa de:

- ✅ **Simulación** → 🚀 **Comunicación real con SRI**
- ✅ **Demo técnico** → 💼 **Sistema empresarial funcional**
- ✅ **Facturas de prueba** → 📋 **Facturas con validez fiscal**

## 📞 Soporte y Contactos

### BCE - Entidad de Certificación
- **Web:** https://www.eci.bce.ec/
- **Teléfono:** 1700 2233733
- **Email:** eci@bce.fin.ec

### Proveedores Autorizados
- **Latinus S.A.**
- **Sodig S.A.**

### En caso de problemas
1. Verificar que cédula esté al día
2. Limpiar caché del navegador
3. Usar conexión estable
4. Contactar soporte del proveedor elegido

---

**💡 Tip:** Este proceso es ideal para desarrolladores que quieren probar su sistema con SRI real sin complicaciones burocráticas.