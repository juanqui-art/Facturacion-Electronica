# Test API Endpoints

Test API functionality: $ARGUMENTS

## Comprehensive API Testing

1. **Server Health Check**
   - Test health endpoint: `curl -s http://localhost:8080/health | jq '.'`
   - Verify server is responding and healthy
   - Check response time and status

2. **API Documentation**
   - Fetch API documentation: `curl -s http://localhost:8080/ | jq '.'`
   - Verify all endpoints are documented
   - Check endpoint specifications

3. **Facturas CRUD Operations**
   - **GET** all facturas: `curl -s http://localhost:8080/api/facturas | jq '.'`
   - **GET** specific factura: `curl -s http://localhost:8080/api/facturas/{id} | jq '.'`
   - **POST** new factura with test data
   - **PUT** update existing factura
   - **DELETE** factura (test soft delete)

4. **Clientes CRUD Operations**
   - **GET** all clientes: `curl -s http://localhost:8080/api/clientes | jq '.'`
   - **GET** specific cliente: `curl -s http://localhost:8080/api/clientes/{id} | jq '.'`
   - **POST** new cliente with valid cédula
   - **PUT** update cliente information
   - **DELETE** cliente (test business logic)

5. **Database Operations**
   - Test database connection and queries
   - Verify data persistence and retrieval
   - Check foreign key relationships
   - Test backup and restore functionality

6. **PDF Generation**
   - Test PDF generation: `curl -s http://localhost:8080/api/facturas/{id}/pdf`
   - Verify PDF content and format
   - Check file size and download headers

## Test Data Templates

### Valid Factura Test Data
```json
{
  "clienteNombre": "Test Cliente",
  "clienteCedula": "1713175071",
  "productos": [
    {
      "codigo": "TEST001",
      "descripcion": "Producto de prueba",
      "cantidad": 1.0,
      "precio_unitario": 10.00
    }
  ]
}
```

### Valid Cliente Test Data
```json
{
  "nombre": "María González",
  "cedula": "1713175071",
  "email": "maria@example.com",
  "telefono": "0987654321",
  "direccion": "Av. Principal 123, Quito"
}
```

## Expected Results

- All endpoints should return appropriate HTTP status codes
- JSON responses should be well-formed and complete
- Database operations should maintain data integrity
- PDF generation should produce valid files
- Error handling should provide meaningful messages

## Performance Metrics

- Response times should be < 500ms for simple queries
- PDF generation should complete within 2-3 seconds
- Database operations should be efficient
- Memory usage should remain stable

## Error Scenarios to Test

- Invalid input data validation
- Non-existent resource requests (404)
- Malformed JSON requests (400)
- Database connection failures
- Concurrent request handling

## API Testing Script

If test_api.sh exists, run comprehensive API test suite and analyze results for any failures or performance issues.