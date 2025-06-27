# Error Analysis Template

## Error Information
- **Error Message**: {ERROR_MESSAGE}
- **Error Code**: {ERROR_CODE}
- **Timestamp**: {TIMESTAMP}
- **Context**: {CONTEXT}

## Diagnostic Steps

1. **Reproduce the Error**
   - Identify exact steps to reproduce
   - Note environmental conditions
   - Check if error is consistent or intermittent

2. **Check System State**
   - Verify server is running and healthy
   - Check database connectivity
   - Confirm configuration is correct

3. **Analyze Logs**
   - Search for related error messages
   - Check timing of related events
   - Look for patterns or trends

4. **Test Related Components**
   - Test individual components
   - Verify dependencies are working
   - Check external service connectivity

## Common Error Categories

### Authentication Errors
- Invalid credentials or certificates
- Session timeout or token expiration
- Permission denied or insufficient privileges

### SRI Integration Errors
- Certificate problems (expired, invalid, missing)
- Network connectivity issues with SRI
- XML schema validation failures
- SOAP communication errors

### Database Errors
- Connection failures or timeouts
- Constraint violations or data integrity issues
- Lock conflicts or transaction failures
- Disk space or permission problems

### API Errors
- Invalid request format or parameters
- Missing required fields or data
- Rate limiting or throttling
- Endpoint not found or method not allowed

## Resolution Strategy

1. **Immediate Actions**
   - Address any critical system failures
   - Implement temporary workarounds if needed
   - Ensure system stability

2. **Root Cause Analysis**
   - Identify underlying cause of error
   - Determine if this is a new issue or regression
   - Check for related configuration changes

3. **Fix Implementation**
   - Implement targeted fix for root cause
   - Test fix in isolated environment
   - Deploy fix with appropriate validation

4. **Prevention Measures**
   - Update error handling and validation
   - Improve monitoring and alerting
   - Document lessons learned

## Follow-up Actions

- Monitor for recurrence of the error
- Update documentation with new insights
- Consider additional preventive measures
- Review and update error handling code