# SRI Electronic Invoicing System - Project Summary

## 🏗️ Architecture Overview

**Core System:** Go-based electronic invoicing system for Ecuador's SRI (Servicio de Rentas Internas)
**Frontend:** Modern Astro 5.x with TypeScript and Tailwind CSS
**Database:** SQLite with comprehensive audit trails
**Integration:** Real SRI web services with digital certificates

## 📊 Project Scale
- **Backend:** 53 Go files (18,319 lines)
- **Frontend:** Modern Astro application (5,199 lines source)
- **Total:** ~23,500 lines of custom code
- **Architecture:** Microservices-ready with clean separation

## 🎯 Key Components

### Backend (Go)
```
api/           - REST API handlers and middleware
sri/           - SRI integration (SOAP, certificates, XAdES-BES)
database/      - SQLite operations and CRUD
models/        - Data structures and XML generation
factory/       - Business logic and validation
config/        - Environment configuration
validators/    - Input validation (RUC, cédulas)
pdf/           - Professional PDF generation
```

### Frontend (Astro)
```
web/src/pages/       - Application pages (facturas, clientes, dashboard)
web/src/components/  - Reusable UI components
web/src/layouts/     - Page layouts
web/src/lib/         - API integration and utilities
web/src/styles/      - Design system and global styles
```

## 🚀 Core Features

### ✅ Implemented
- **Complete CRUD:** Facturas and clientes with business logic
- **SRI Integration:** Real connectivity with demo/production modes  
- **PDF Generation:** Professional invoice PDFs
- **Digital Certificates:** BCE certificate support (demo mode ready)
- **Modern Frontend:** Responsive Astro application
- **Database:** Full persistence with audit trails
- **API:** RESTful endpoints with validation

### 🔄 In Progress  
- Authentication and authorization system
- Advanced reporting and analytics
- Multi-company support

## 🛠️ Development Workflow

### Quick Start
```bash
go run main.go test_validaciones.go api    # Start API server
cd web && npm run dev                      # Start frontend
```

### Custom Commands (Token-Optimized)
```bash
/debug:auth "login issue"     # Intelligent auth debugging
/debug:sri "cert problem"     # SRI-specific diagnostics  
/test:api "endpoints"         # API testing with analysis
/setup:cert                   # Interactive certificate setup
/db:query "SQL"              # Database queries with insights
/deploy:check                # Pre-deployment validation
```

## 📁 Critical Files (Quick Access)

### Configuration
- `config/desarrollo.json` - Development environment settings
- `.claude/settings.local.json` - Claude optimization settings

### Core Backend
- `api/database_handlers.go` - Main API endpoints  
- `sri/integration_real.go` - SRI communication
- `database/database.go` - Database operations
- `models/factura.go` - Invoice data structures

### Frontend  
- `web/src/pages/index.astro` - Main dashboard
- `web/src/pages/facturas.astro` - Invoice management
- `web/src/layouts/Layout.astro` - Main layout

## 🔧 Technical Specifications

### SRI Integration
- **Environment:** Certificación (pruebas) / Producción
- **Certificates:** BCE digital certificates (.p12)
- **Protocols:** SOAP web services, XAdES-BES signatures
- **Validation:** RUC validation, clave de acceso generation

### Database Schema
```sql
facturas    - Invoice records with SRI compliance
clientes    - Customer data with validation  
productos   - Product catalog (planned)
audit_log   - Full audit trail
```

### API Endpoints
```
GET/POST /api/facturas     - Invoice CRUD
GET/POST /api/clientes     - Client CRUD  
GET      /api/health       - System health
GET      /api/facturas/{id}/pdf - PDF generation
```

## 🎯 Business Logic

### Invoice Workflow
1. **Create:** Validate client and products
2. **Generate:** Create XML with SRI compliance
3. **Sign:** Apply digital signature (XAdES-BES)  
4. **Submit:** Send to SRI for authorization
5. **Authorize:** Receive SRI authorization number
6. **PDF:** Generate RIDE (printable invoice)

### SRI Compliance
- **RUC Validation:** Módulo-11 algorithm  
- **Clave de Acceso:** 49-digit unique key
- **XML Schema:** Full SRI specification compliance
- **Digital Signature:** XAdES-BES with BCE certificates

## 🔒 Security & Production

### Security Features
- Input validation and sanitization
- Digital certificate management
- Audit trail for all operations
- Secure configuration management

### Production Readiness
- Environment-specific configurations
- Digital certificate integration  
- SRI production endpoints
- Professional PDF generation
- Comprehensive error handling

## 📈 Scalability Considerations

### Current Optimizations
- Custom commands for efficient development
- Context management for large codebase
- Specialized tools for frequent operations
- Token usage optimization

### Future Enhancements
- Multi-tenant architecture
- API rate limiting
- Horizontal scaling support
- Advanced caching strategies

This summary provides context for efficient development without requiring full codebase analysis.