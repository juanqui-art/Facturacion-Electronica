# Database Query and Analysis

Execute database query: $ARGUMENTS

## Database Query Execution

1. **Connection Verification**
   - Check if database file exists and is accessible
   - Verify database file permissions and size
   - Test basic connection to SQLite database

2. **Query Execution**
   - Execute the provided SQL query safely
   - Handle potential SQL injection concerns
   - Format and display results clearly

3. **Result Analysis**
   - Analyze query results for data integrity
   - Check for expected data patterns
   - Identify any anomalies or issues

4. **Common Database Queries for SRI System**

   **Facturas Analysis:**
   ```sql
   -- Total facturas by status
   SELECT estado, COUNT(*) as total FROM facturas GROUP BY estado;
   
   -- Recent facturas
   SELECT id, numero_factura, cliente_nombre, total, fecha_emision 
   FROM facturas ORDER BY fecha_emision DESC LIMIT 10;
   
   -- Facturas summary
   SELECT COUNT(*) as total_facturas, SUM(total) as total_monto 
   FROM facturas WHERE estado = 'AUTORIZADA';
   ```

   **Clientes Analysis:**
   ```sql
   -- Active clients
   SELECT id, nombre, cedula, email, activo FROM clientes WHERE activo = 1;
   
   -- Clients with most invoices
   SELECT c.nombre, COUNT(f.id) as num_facturas, SUM(f.total) as total_ventas
   FROM clientes c LEFT JOIN facturas f ON c.cedula = f.cliente_cedula
   GROUP BY c.id ORDER BY num_facturas DESC LIMIT 10;
   ```

   **System Health:**
   ```sql
   -- Database statistics
   SELECT name, sql FROM sqlite_master WHERE type='table';
   
   -- Record counts
   SELECT 'facturas' as tabla, COUNT(*) as registros FROM facturas
   UNION ALL
   SELECT 'clientes' as tabla, COUNT(*) as registros FROM clientes;
   ```

5. **Database Maintenance Queries**
   ```sql
   -- Check for orphaned records
   SELECT * FROM facturas WHERE cliente_cedula NOT IN (SELECT cedula FROM clientes);
   
   -- Database integrity check
   PRAGMA integrity_check;
   
   -- Database size and optimization
   PRAGMA page_count;
   PRAGMA page_size;
   ```

## Query Safety and Best Practices

- Always use READ-ONLY queries unless explicitly specified
- Avoid DELETE or UPDATE operations without explicit confirmation
- Use LIMIT clauses for potentially large result sets
- Validate query syntax before execution

## Result Interpretation

- Explain what the query results mean in business context
- Identify any data quality issues or inconsistencies
- Suggest follow-up queries if needed
- Recommend data cleanup or fixes if problems found

## Database Schema Information

```sql
-- Show table structure
.schema facturas
.schema clientes

-- Show indexes
.indexes

-- Show database info
.dbinfo
```

## Performance Analysis

- Check query execution time
- Identify slow queries and suggest optimizations
- Recommend indexes if needed
- Monitor database file size growth

## Common Issues and Solutions

- **Lock issues**: Check for long-running transactions
- **Corruption**: Run integrity checks and suggest repair
- **Performance**: Analyze query plans and suggest optimizations
- **Data consistency**: Check foreign key relationships and constraints

## Backup and Recovery

- Verify recent backups exist
- Test backup integrity if needed
- Suggest backup strategy improvements
- Document recovery procedures