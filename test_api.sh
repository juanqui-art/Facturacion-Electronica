#!/bin/bash

# Script para probar la API de facturación electrónica

API_URL="http://localhost:8080"

echo "🧪 PROBANDO API DE FACTURACIÓN ELECTRÓNICA"
echo "=========================================="

# Test 1: Health check
echo ""
echo "1. 🏥 Probando health check..."
curl -s "$API_URL/health" | jq '.'

# Test 2: Documentación de la API
echo ""
echo "2. 📚 Obteniendo documentación de la API..."
curl -s "$API_URL/" | jq '.endpoints'

# Test 3: Crear factura
echo ""
echo "3. 📝 Creando nueva factura..."
curl -s -X POST "$API_URL/api/facturas" \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "Ana García",
    "clienteCedula": "1713175071",
    "productos": [
      {
        "codigo": "LAPTOP001",
        "descripcion": "Laptop HP Pavilion",
        "cantidad": 1,
        "precioUnitario": 800.00
      },
      {
        "codigo": "MOUSE001",
        "descripcion": "Mouse Inalámbrico",
        "cantidad": 2,
        "precioUnitario": 25.00
      }
    ],
    "includeXML": true
  }' | jq '.'

# Test 4: Listar facturas
echo ""
echo "4. 📋 Listando facturas creadas..."
curl -s "$API_URL/api/facturas" | jq '.'

# Test 5: Obtener factura específica con XML
echo ""
echo "5. 🔍 Obteniendo factura específica con XML..."
curl -s "$API_URL/api/facturas/FAC-000001?includeXML=true" | jq '.id, .status, .factura.InfoFactura.ImporteTotal'

# Test 6: Error handling - datos inválidos
echo ""
echo "6. ❌ Probando manejo de errores (cédula inválida)..."
curl -s -X POST "$API_URL/api/facturas" \
  -H "Content-Type: application/json" \
  -d '{
    "clienteNombre": "Cliente Test",
    "clienteCedula": "123",
    "productos": [
      {
        "codigo": "TEST001",
        "descripcion": "Producto test",
        "cantidad": 1,
        "precioUnitario": 100.00
      }
    ]
  }' | jq '.'

echo ""
echo "✅ Pruebas de API completadas!"
echo "💡 Para ver logs del servidor, revisa la terminal donde está corriendo."