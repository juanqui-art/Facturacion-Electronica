# Debug Authentication Issues

Analyze authentication error: $ARGUMENTS

## Diagnostic Steps

1. **Check API Server Status**
   - Verify Go API server is running: `curl -s http://localhost:8080/health`
   - Check server logs for authentication errors
   - Confirm API is responding to requests

2. **Verify Configuration**
   - Read config/desarrollo.json for authentication settings
   - Check if certificate configuration is correct
   - Validate RUC and empresa settings

3. **Test API Endpoints**
   - Test health endpoint: `curl -s http://localhost:8080/health`
   - Test facturas endpoint: `curl -s http://localhost:8080/api/facturas`
   - Check for 401, 403, or other auth-related status codes

4. **Frontend Authentication**
   - If using frontend, check browser console for errors
   - Verify API_BASE configuration in Astro pages
   - Check for CORS issues between frontend and backend

5. **Database Connection**
   - Verify database file exists and is accessible
   - Check database connection in Go application
   - Test database queries work correctly

## Common Authentication Issues in SRI System

- **Certificate Issues**: Missing or invalid digital certificates
- **RUC Configuration**: Invalid RUC in empresa configuration
- **CORS Problems**: Frontend and backend on different ports
- **Database Access**: File permissions or SQLite connection issues
- **API Endpoint Conflicts**: Route conflicts or middleware issues

## Suggested Fixes

Based on the analysis, provide specific code fixes with file:line references and minimal code changes needed to resolve the authentication issue.