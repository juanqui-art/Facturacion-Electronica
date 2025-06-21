# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-based electronic invoicing system for Ecuador's SRI (Servicio de Rentas Internas). Currently in Week 1/12 of development, focusing on XML generation compliant with SRI schemas. The system will eventually support digital signatures (XAdES-BES), SOAP communication with SRI services, and a web interface.

## Development Commands

```bash
# Run the application
go run main.go

# Build executable
go build -o facturacion-sri main.go

# Clean dependencies
go mod tidy

# Format code
go fmt ./...
```

## Architecture

### Core Data Structures

The system uses Go's `encoding/xml` package with struct tags for SRI-compliant XML generation:

- **InfoTributaria**: Tax information (RUC, establishment codes, document sequences)
- **InfoFactura**: Invoice details (dates, customer info, totals)  
- **Detalle**: Line items (products/services with quantities and prices)
- **Factura**: Complete invoice structure combining all components

### SRI-Specific Requirements

- **Ambiente**: "1" for testing, "2" for production environment
- **ClaveAcceso**: 49-digit unique access key (currently placeholder, real algorithm planned for Week 4)
- **CodDoc**: Document type codes ("01" for invoices, "04" for credit notes)
- **Secuencial**: 9-digit sequential invoice numbering (000000001, 000000002, etc.)
- **Tax Calculations**: Automatic 15% IVA calculation for Ecuador
- **Date Format**: DD/MM/YYYY format required by SRI

### Development Roadmap Context

- **Weeks 1-2**: Basic XML generation (current phase)
- **Weeks 3-4**: Real access key generation and advanced validations
- **Weeks 5-6**: PKCS#12 certificates and XAdES-BES digital signatures
- **Weeks 7-8**: Web interface and CRUD operations
- **Weeks 9-12**: Additional document types (credit/debit notes, delivery guides)

### SRI Endpoints

- **Certification**: `https://celcer.sri.gob.ec/comprobantes-electronicos-ws/`
- **Production**: `https://cel.sri.gob.ec/comprobantes-electronicos-ws/`

## Key Patterns

- Factory functions for creating invoices with proper defaults
- Automatic calculation of taxes and totals
- XML marshaling with proper SRI namespace requirements
- Sequential document numbering system
- Environment-aware configuration (test vs production)

## Go Learning Context

**IMPORTANT**: The user is learning Go while developing this project. Always provide educational explanations when:

- Introducing new Go concepts
- Writing or modifying code
- Explaining architectural decisions
- Debugging or troubleshooting

### Learning Progression by Project Phase

**Weeks 1-2 (Current)**: Go Fundamentals
- Structs and struct tags
- Basic types (string, float64, int)
- Slices and arrays
- XML marshaling/unmarshaling
- Error handling patterns
- Package imports and basic functions

**Weeks 3-4**: Intermediate Go
- Interfaces and methods
- Pointers and memory management
- File I/O operations
- String manipulation and validation
- Time and date handling
- Testing with Go's testing package

**Weeks 5-6**: Advanced Go
- Goroutines and concurrency
- Channels for communication
- HTTP client implementations
- Cryptographic operations
- Certificate handling
- Third-party package integration

**Weeks 7-8**: Web Development
- HTTP server setup
- Template engines
- Middleware patterns
- Database integration
- JSON handling
- RESTful API design

**Weeks 9-12**: Production Go
- Build and deployment
- Configuration management
- Logging and monitoring
- Performance optimization
- Error recovery
- Production best practices

### Teaching Approach

1. **Explain before implementing**: Always explain what Go concept you're about to use
2. **Show alternatives**: When there are multiple ways to do something in Go, mention them
3. **Reference learning resources**: Point to Go documentation when introducing new concepts
4. **Progressive complexity**: Build on previously learned concepts
5. **Real-world context**: Connect Go features to the SRI invoicing requirements

## Current Limitations

- Access key generation is placeholder (real algorithm pending)
- No digital signature implementation yet
- Single invoice type support only
- No persistence layer (in-memory only)
- No SRI communication (XML generation only)