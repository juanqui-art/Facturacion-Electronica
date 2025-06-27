# Component Index - Quick Reference

## 🎯 Quick Navigation for Common Tasks

### Authentication Issues → `api/middleware.go`, `config/desarrollo.json`
### SRI Integration → `sri/integration_real.go`, `sri/soap_client.go`  
### Database Problems → `database/database.go`, `database/additional_functions.go`
### API Endpoints → `api/database_handlers.go`, `api/server.go`
### PDF Generation → `pdf/generator.go`
### Frontend Issues → `web/src/pages/*.astro`, `web/src/lib/api.ts`

---

## 📁 Core Components by Function

### 🔐 Authentication & Security
```
api/middleware.go           - CORS, logging, request validation
config/desarrollo.json      - Environment settings and credentials  
sri/certificado.go         - Digital certificate management
validators/validations.go   - Input validation (RUC, cédulas)
```

### 🇪🇨 SRI Integration
```
sri/integration_real.go     - Main SRI communication logic
sri/soap_client.go          - SOAP web service client
sri/autorizacion.go         - Authorization and clave de acceso
sri/xades_bes.go           - Digital signature implementation
sri/errores_sri.go         - SRI-specific error handling
```

### 🗄️ Database Operations  
```
database/database.go        - Core database operations
database/additional_functions.go - Extended CRUD operations
database/backup.go          - Backup and recovery
models/factura.go          - Data structures and XML generation
```

### 🌐 API Layer
```
api/database_handlers.go    - Main CRUD endpoints
api/handlers.go            - Base API functionality  
api/server.go              - Server setup and routing
api/middleware.go          - Request processing
```

### 📄 Document Generation
```
pdf/generator.go           - Professional PDF creation
models/factura.go          - XML generation for SRI
factory/factura_factory.go - Business logic and validation
```

### 🎨 Frontend (Astro)
```
web/src/pages/index.astro      - Main dashboard
web/src/pages/facturas.astro   - Invoice management
web/src/pages/clientes.astro   - Client management  
web/src/layouts/Layout.astro   - Main layout and navigation
web/src/lib/api.ts            - API integration utilities
```

---

## 🚨 Common Issues & Quick Fixes

### Authentication Failures
**Files:** `api/middleware.go:25`, `config/desarrollo.json:16`
**Common Fix:** Check CORS configuration and certificate settings

### SRI Connection Issues  
**Files:** `sri/soap_client.go:45`, `sri/integration_real.go:80`
**Common Fix:** Verify endpoints and certificate configuration

### Database Errors
**Files:** `database/database.go:30`, `database/additional_functions.go:15`  
**Common Fix:** Check file permissions and SQLite connection

### PDF Generation Problems
**Files:** `pdf/generator.go:40`, `models/factura.go:120`
**Common Fix:** Verify invoice data completeness and formatting

### Frontend API Errors
**Files:** `web/src/lib/api.ts:10`, `web/src/pages/facturas.astro:165`
**Common Fix:** Check API_BASE configuration and CORS settings

---

## 🔧 Configuration Files

### Environment Settings
```
config/desarrollo.json     - Development environment
config/produccion.json     - Production environment (when ready)
.claude/settings.local.json - Claude optimization settings
```

### Project Management
```
docs/CLAUDE.md             - Complete project documentation
.claude/project-summary.md  - This summary for quick context
.claude/component-index.md  - Component navigation (this file)
```

---

## 📊 Key Metrics & Locations

### Performance Critical
- **Database queries:** `database/database.go:150-200`
- **SRI communication:** `sri/soap_client.go:80-120`  
- **PDF generation:** `pdf/generator.go:60-100`

### Business Logic
- **Invoice validation:** `factory/factura_factory.go:25-60`
- **RUC validation:** `validators/validations.go:15-50`
- **Clave de acceso:** `sri/autorizacion.go:30-70`

### API Performance  
- **Main endpoints:** `api/database_handlers.go:20-150`
- **Error handling:** `api/handlers.go:15-40`
- **Middleware:** `api/middleware.go:10-50`

---

## 💡 Development Tips

### For Quick Debugging
1. **Use Task tool** instead of reading full files
2. **Check component index** for specific functionality
3. **Use custom commands** for targeted analysis
4. **Reference project summary** for context

### For New Features
1. **Start with models/** for data structures
2. **Add API endpoints** in `api/database_handlers.go`
3. **Implement frontend** in `web/src/pages/`
4. **Test integration** with custom commands

This index optimizes navigation without requiring full file reads.