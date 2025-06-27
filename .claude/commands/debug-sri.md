# Debug SRI Integration Issues

Analyze SRI integration error: $ARGUMENTS

## Diagnostic Steps

1. **Check Certificate Configuration**
   - Verify certificate files exist in certificados/ directory
   - Check certificate configuration in config/desarrollo.json
   - Test certificate loading and validation
   - Confirm password and file path are correct

2. **Test SRI Connectivity**
   - Test SRI endpoints connectivity:
     - Certificación: https://celcer.sri.gob.ec/comprobantes-electronicos-ws/RecepcionComprobantesOffline
     - Autorización: https://celcer.sri.gob.ec/comprobantes-electronicos-ws/AutorizacionComprobantesOffline
   - Check internet connectivity and firewall restrictions
   - Verify SSL/TLS handshake with SRI servers

3. **Validate XML Generation**
   - Test factura XML generation with sample data
   - Verify XML schema compliance with SRI specifications
   - Check clave de acceso generation (49-digit validation)
   - Validate RUC format and módulo-11 calculation

4. **Test SOAP Client**
   - Execute SOAP client demo: `go run main.go test_validaciones.go soap`
   - Check circuit breaker status and retry logic
   - Verify SOAP message format and headers
   - Test both recepción and autorización services

5. **Verify Business Logic**
   - Check RUC validation: empresa RUC must be valid
   - Verify establecimiento and punto de emisión configuration
   - Test secuencial numbering generation
   - Validate producto and cliente data formats

## Common SRI Integration Issues

- **Certificate Problems**: Invalid, expired, or missing digital certificates
- **XML Schema Errors**: Non-compliant XML structure or missing required fields
- **Network Issues**: Firewall blocking SRI endpoints or connectivity problems
- **RUC Validation**: Invalid RUC format or módulo-11 calculation errors
- **Clave de Acceso**: Incorrect 49-digit key generation or validation
- **SOAP Format**: Malformed SOAP messages or incorrect headers
- **Environment Config**: Wrong SRI endpoints for ambiente (pruebas vs producción)

## SRI-Specific Commands to Run

```bash
# Test SRI integration demo
go run main.go test_validaciones.go sri

# Test real SRI communication
go run main.go test_validaciones.go test-sri

# Test certificate loading
go run main.go test_validaciones.go certificacion

# Test SOAP client
go run main.go test_validaciones.go soap
```

## Suggested Fixes

Based on the specific SRI error, provide targeted solutions:
- Certificate configuration fixes
- XML schema corrections
- Network configuration adjustments
- Business logic validation fixes
- Environment-specific configuration changes