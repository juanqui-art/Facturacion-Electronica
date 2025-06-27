# Directorio de Certificados Digitales

## 🔐 Certificados para Integración SRI

Este directorio contiene los certificados digitales necesarios para la integración real con el SRI de Ecuador.

### Estado Actual
- ⚠️ **Sin certificado digital real**
- ✅ Sistema preparado para certificado real
- ✅ Modo demo funcional sin certificado

### Archivos Esperados

#### Certificado de Producción
```
certificado-produccion.p12  # Certificado real del BCE
password-produccion.txt     # Contraseña del certificado (NUNCA en Git)
```

#### Certificado de Desarrollo (Opcional)
```
certificado-desarrollo.p12  # Certificado de pruebas/desarrollo
password-desarrollo.txt     # Contraseña del certificado de desarrollo
```

## 🚀 Obtener Certificado Digital Real

### Proceso BCE Online (Recomendado)
- **Costo:** $24.64 USD (certificado de archivo)
- **Tiempo:** 30 minutos (100% online)
- **Vigencia:** 2 años
- **URL:** https://www.eci.bce.ec/

### Pasos Detallados
1. **Acceder al portal BCE:** https://www.eci.bce.ec/
2. **Seleccionar proveedor:** Latinus S.A. o Sodig S.A.
3. **Completar datos:** Cédula, datos personales, email
4. **Validación biométrica:** Fotografía de cédula + rostro
5. **Pago online:** $24.64 USD con tarjeta
6. **Descarga inmediata:** Archivo .p12 con contraseña

### Requisitos
- ✅ Cédula ecuatoriana vigente
- ✅ Dispositivo con cámara funcional
- ✅ Conexión a internet estable
- ✅ Tarjeta de crédito/débito

## 🔧 Configuración Posterior

### Una vez obtenido el certificado:

1. **Guardar archivo .p12**
   ```bash
   # Copiar certificado a este directorio
   cp /Downloads/certificado.p12 ./certificados/certificado-produccion.p12
   ```

2. **Actualizar configuración**
   ```json
   // config/desarrollo.json o config/produccion.json
   {
     "certificado": {
       "rutaArchivo": "./certificados/certificado-produccion.p12",
       "password": "tu_contraseña_aquí"
     }
   }
   ```

3. **Probar integración**
   ```bash
   # Testing con certificado real
   go run main.go test_validaciones.go test-sri
   ```

4. **Activar producción**
   ```json
   // config/produccion.json
   {
     "ambiente": {
       "codigo": "2",
       "descripcion": "Ambiente de Producción"
     }
   }
   ```

## 🛡️ Seguridad

### ⚠️ IMPORTANTE - Nunca versionar:
- ❌ Archivos .p12 (certificados)
- ❌ Contraseñas en texto plano
- ❌ Claves privadas
- ❌ Datos de configuración sensibles

### ✅ Incluido en .gitignore:
```gitignore
# Certificados digitales y llaves privadas
*.p12
*.pfx
*.key
*.pem
*.crt
*.cer
certificados/
```

## 📋 Checklist Post-Certificado

### Cuando obtengas el certificado real:
- [ ] Certificado .p12 guardado en este directorio
- [ ] Contraseña configurada en JSON (temporal) o variable de entorno
- [ ] Testing básico completado (`test-sri`)
- [ ] Comunicación SRI verificada
- [ ] Firma digital funcionando
- [ ] Autorización de facturas exitosa
- [ ] Backup seguro del certificado creado
- [ ] Documentación de empresa actualizada

## 🔄 Migración Automática

El sistema está preparado para activación automática:

```bash
# El sistema detecta automáticamente la presencia del certificado
# y cambia de modo demo a integración real
go run main.go test_validaciones.go sri
```

## 📞 Soporte

### Problemas con certificado:
- **BCE:** 1700 2233733 / eci@bce.fin.ec
- **Proveedores:** Latinus S.A. o Sodig S.A.

### Problemas de integración:
- Verificar configuración en `config/`
- Ejecutar tests: `go test ./sri -v`
- Revisar logs de conexión SRI

---

**🎯 Objetivo:** Una vez obtenido el certificado, el sistema pasa de demo a producción en minutos.