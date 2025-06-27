# Setup Digital Certificates

Guide for digital certificate setup: $ARGUMENTS

## Interactive Certificate Setup

1. **Check Current Certificate Status**
   - List files in certificados/ directory
   - Check config/desarrollo.json certificate configuration
   - Verify current certificate mode (demo vs real)

2. **Certificate Requirements Analysis**
   - Identify if this is for development or production
   - Determine certificate type needed (archivo .p12 vs token)
   - Check RUC and empresa requirements

3. **BCE Certificate Acquisition Guide**
   - **Cost**: $24.64 USD for archivo certificate (recommended for development)
   - **Time**: 30 minutes for online process
   - **URL**: https://www.eci.bce.ec/
   - **Process**: 100% online with biometric validation

4. **Step-by-Step Instructions**
   ```
   1. Visit: https://www.eci.bce.ec/
   2. Select "Solicitud de Certificado"
   3. Choose provider: Latinus S.A. or Sodig S.A.
   4. Complete personal/business data
   5. Biometric validation (photo + ID)
   6. Payment: $24.64 USD
   7. Download .p12 file immediately
   ```

5. **Certificate Installation**
   - Save .p12 file to certificados/ directory
   - Update config/desarrollo.json with:
     - Certificate file path
     - Password
     - Enable certificate mode
   - Test certificate loading

6. **Configuration Update**
   ```json
   {
     "certificado": {
       "rutaArchivo": "./certificados/mi-certificado.p12",
       "password": "mi_contraseña_segura",
       "habilitado": true,
       "modoDemo": false
     }
   }
   ```

7. **Certificate Validation**
   - Test certificate loading: `go run main.go test_validaciones.go sri`
   - Verify certificate information and validity
   - Test SRI communication with real certificate
   - Confirm XAdES-BES signature capability

## Security Checklist

- [ ] Certificate file has restricted permissions (600)
- [ ] Password is stored securely (consider environment variables)
- [ ] Certificate backup created in secure location
- [ ] .gitignore prevents certificate files from being committed
- [ ] Certificate expiration monitoring configured

## Production Considerations

- Certificate renewal strategy (before 2-year expiration)
- Backup and disaster recovery procedures
- Multi-environment certificate management
- Monitoring and alerting for certificate issues

## Troubleshooting Common Issues

- **Invalid password**: Verify password exactly as provided by BCE
- **File not found**: Check file path and permissions
- **Expired certificate**: Verify certificate validity dates
- **Wrong format**: Ensure file is PKCS#12 (.p12) format
- **SRI communication fails**: Test with certificado mode first

## Migration from Demo to Production

1. Obtain real certificate from BCE
2. Update configuration files
3. Test in certificación environment first
4. Switch to production environment
5. Monitor first real transactions

## Next Steps After Setup

- Test factura generation with real certificate
- Configure production environment
- Set up monitoring and logging
- Train users on the system
- Establish backup procedures