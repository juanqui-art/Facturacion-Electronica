# Deployment Readiness Check

Verify deployment readiness: $ARGUMENTS

## Pre-Deployment Checklist

1. **Code Quality and Testing**
   - Run all tests: `go test ./...`
   - Check test coverage is adequate (>80%)
   - Verify no failing tests or build errors
   - Run linting: `go fmt ./... && go vet ./...`

2. **Configuration Validation**
   - Check config/produccion.json exists and is valid
   - Verify all environment-specific settings
   - Confirm RUC and empresa data for production
   - Validate certificate configuration for production

3. **Security Assessment**
   - Verify no secrets in code or configuration files
   - Check .gitignore includes all sensitive files
   - Confirm HTTPS configuration for production
   - Validate certificate security and permissions

4. **Database Readiness**
   - Verify database migrations are complete
   - Check database backup strategy
   - Test database performance under load
   - Confirm data integrity and consistency

5. **SRI Integration Validation**
   - Test SRI connectivity in production environment
   - Verify digital certificate is valid and not expired
   - Confirm RUC is registered and active with SRI
   - Test factura authorization flow end-to-end

6. **Performance Testing**
   - Load test critical endpoints
   - Verify response times are acceptable
   - Check memory usage patterns
   - Test concurrent user scenarios

7. **Frontend Build Validation**
   - Build frontend for production: `cd web && npm run build`
   - Test frontend production build
   - Verify static assets are optimized
   - Check for JavaScript errors in production mode

## Production Environment Checklist

### Server Configuration
- [ ] Go binary compiled for target architecture
- [ ] All dependencies available on production server
- [ ] Port 8080 (or configured port) is available
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules allow necessary traffic

### SRI Production Requirements
- [ ] Real digital certificate from BCE installed
- [ ] Production RUC configured and validated
- [ ] SRI production endpoints configured
- [ ] Test transactions successful in production

### Database Production Setup
- [ ] Production database initialized
- [ ] Backup procedures configured and tested
- [ ] Database permissions configured correctly
- [ ] Monitoring and alerting configured

### Monitoring and Logging
- [ ] Application logging configured
- [ ] Error monitoring setup
- [ ] Performance monitoring active
- [ ] Health check endpoints configured

## Deployment Commands

```bash
# Build for production
go build -o facturacion-sri main.go test_validaciones.go

# Frontend build
cd web && npm run build && cd ..

# Run production tests
go test ./... -v

# Database backup before deployment
cp demo_facturacion.db backup_pre_deploy_$(date +%Y%m%d_%H%M%S).db
```

## Post-Deployment Verification

1. **Smoke Tests**
   - Health check endpoint responds
   - API endpoints return expected responses
   - Database connections work
   - SRI integration functions correctly

2. **Business Logic Validation**
   - Create test factura in production
   - Verify PDF generation works
   - Test SRI authorization flow
   - Confirm email notifications (if implemented)

3. **Performance Monitoring**
   - Monitor response times
   - Check memory and CPU usage
   - Verify database performance
   - Monitor error rates

## Rollback Plan

- [ ] Previous version backup available
- [ ] Database rollback procedure documented
- [ ] Rollback triggers and procedures defined
- [ ] Communication plan for rollback scenario

## Production Monitoring Setup

```bash
# Health monitoring
curl -f http://localhost:8080/health || echo "Health check failed"

# Log monitoring
tail -f /var/log/facturacion-sri.log

# Database monitoring
ls -la *.db
```

## Common Deployment Issues

- **Certificate Problems**: Ensure production certificates are correctly installed
- **Database Migration**: Verify all schema changes are applied
- **Network Configuration**: Check firewall and port configurations
- **Permission Issues**: Verify file and directory permissions
- **Environment Variables**: Ensure all required environment variables are set

## Success Criteria

- [ ] All tests pass
- [ ] Performance benchmarks met
- [ ] Security checklist completed
- [ ] SRI integration working
- [ ] Monitoring active
- [ ] Rollback plan ready
- [ ] Documentation updated

Deployment is ready when all criteria are met and verified.