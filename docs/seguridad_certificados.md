# Seguridad de Certificados Digitales: Mejores Prácticas

## 🔐 Visión General

La gestión segura de certificados digitales es **crítica** en un sistema SaaS de facturación electrónica. Cada certificado representa la **identidad legal** de una empresa y su **autorización fiscal**, por lo que debe protegerse con el más alto nivel de seguridad.

---

## 🎯 Principios Fundamentales de Seguridad

### **Analogía: Caja Fuerte del Banco**

```
🏦 BANCO TRADICIONAL:
├── 🏛️ Edificio blindado (infraestructura segura)
├── 🚪 Múltiples puertas con códigos (autenticación por capas)
├── 📹 Cámaras 24/7 (monitoreo y auditoría)
├── 🔐 Cajas individuales por cliente (aislamiento)
├── 👮 Guardias especializados (personal de seguridad)
└── 🚨 Alarmas automáticas (detección de intrusiones)

💾 TU SISTEMA DE CERTIFICADOS:
├── 🛡️ Servidores endurecidos (infraestructura segura)
├── 🔑 Autenticación multi-factor (acceso por capas)
├── 📊 Logs completos (monitoreo y auditoría)
├── 🗂️ Encriptación por empresa (aislamiento)
├── 👨‍💻 Equipo especializado (personal de seguridad)
└── 🚨 Alertas automáticas (detección de anomalías)
```

### **Los 7 Pilares de Seguridad**

1. **🔐 Encriptación en Reposo** - Certificados nunca almacenados en texto plano
2. **🚛 Encriptación en Tránsito** - TLS 1.3 para todas las comunicaciones
3. **🏢 Aislamiento por Empresa** - Cada tenant completamente separado
4. **🔑 Gestión de Claves** - Claves de encriptación rotadas y seguras
5. **📊 Auditoría Completa** - Todo acceso registrado y monitoreado
6. **⏰ Validación Continua** - Verificación automática de integridad
7. **🚨 Respuesta a Incidentes** - Procedimientos de emergencia definidos

---

## 🔒 Arquitectura de Seguridad

### **Almacenamiento Multi-Capa**

```
📁 ESTRUCTURA DE ALMACENAMIENTO SEGURO:

/data/certificados/
├── empresa_1792146739001/          # RUC como identificador
│   ├── cert_active.p12.enc         # Certificado actual encriptado
│   ├── cert_backup.p12.enc         # Backup del certificado
│   ├── .metadata                   # Metadatos encriptados
│   ├── .integrity_hash             # Hash de verificación
│   └── .access_log                 # Log de accesos
│
├── empresa_0992345678001/
│   ├── cert_active.p12.enc
│   ├── cert_backup.p12.enc
│   ├── .metadata
│   ├── .integrity_hash
│   └── .access_log
│
└── .master_keys/                   # Directorio protegido
    ├── current.key                 # Clave maestra actual
    ├── previous.key                # Clave anterior (para rotación)
    └── salt.dat                    # Salt para derivación de claves

🔐 PERMISOS UNIX:
├── /data/certificados/: 700 (rwx------)
├── Subdirectorios empresa: 700 (rwx------)
├── Archivos .p12.enc: 600 (rw-------)
└── .master_keys/: 700 (rwx------)
```

### **Encriptación Avanzada**

```go
// Implementación de encriptación por capas
type CertificateSecurityManager struct {
    masterKey        []byte           // Clave maestra del sistema
    keyDerivationSalt []byte          // Salt para derivación
    encryptionMethod string          // AES-256-GCM
    compressionLevel int             // Nivel de compresión
    auditLogger      *AuditLogger    // Logger de auditoría
}

// Función de encriptación con múltiples capas de seguridad
func (csm *CertificateSecurityManager) EncryptCertificate(
    empresaID string, 
    ruc string, 
    certificateData []byte, 
    userID string) (*EncryptedCertificate, error) {

    // 1. Validar integridad del certificado original
    originalHash := sha512.Sum512(certificateData)
    
    // 2. Comprimir datos para reducir superficie de ataque
    compressedData, err := compressData(certificateData)
    if err != nil {
        return nil, fmt.Errorf("error comprimiendo certificado: %v", err)
    }
    
    // 3. Derivar clave específica de empresa
    empresaKey := csm.deriveEnterpriseKey(empresaID, ruc)
    
    // 4. Generar IV único para esta operación
    iv := make([]byte, 12) // GCM recommended IV size
    if _, err := rand.Read(iv); err != nil {
        return nil, fmt.Errorf("error generando IV: %v", err)
    }
    
    // 5. Encriptar con AES-256-GCM (autenticación incluida)
    block, err := aes.NewCipher(empresaKey)
    if err != nil {
        return nil, fmt.Errorf("error creando cipher: %v", err)
    }
    
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("error creando GCM: %v", err)
    }
    
    // 6. Datos adicionales para autenticación (AAD)
    aad := []byte(fmt.Sprintf("empresa:%s:ruc:%s:user:%s:timestamp:%d", 
        empresaID, ruc, userID, time.Now().Unix()))
    
    encryptedData := aesgcm.Seal(nil, iv, compressedData, aad)
    
    // 7. Crear estructura final con metadatos
    encryptedCert := &EncryptedCertificate{
        EmpresaID:     empresaID,
        RUC:           ruc,
        EncryptedData: encryptedData,
        IV:            iv,
        AAD:           aad,
        Algorithm:     "AES-256-GCM",
        CompressionLevel: csm.compressionLevel,
        OriginalHash:  hex.EncodeToString(originalHash[:]),
        CreatedAt:     time.Now(),
        CreatedBy:     userID,
        Version:       "1.0",
    }
    
    // 8. Registrar en auditoría
    csm.auditLogger.LogCertificateEncryption(empresaID, userID, "SUCCESS")
    
    return encryptedCert, nil
}

// Derivación de clave específica por empresa
func (csm *CertificateSecurityManager) deriveEnterpriseKey(empresaID, ruc string) []byte {
    // Usar PBKDF2 con salt específico de empresa
    enterpriseSalt := append(csm.keyDerivationSalt, []byte(empresaID+":"+ruc)...)
    
    // 100,000 iteraciones para mayor seguridad
    return pbkdf2.Key(csm.masterKey, enterpriseSalt, 100000, 32, sha256.New)
}

// Desencriptación con verificación de integridad
func (csm *CertificateSecurityManager) DecryptCertificate(
    encryptedCert *EncryptedCertificate, 
    userID string) ([]byte, error) {

    // 1. Verificar permisos de acceso
    if !csm.hasAccessPermission(encryptedCert.EmpresaID, userID) {
        csm.auditLogger.LogCertificateAccess(encryptedCert.EmpresaID, userID, "DENIED")
        return nil, fmt.Errorf("acceso denegado")
    }
    
    // 2. Derivar clave de empresa
    empresaKey := csm.deriveEnterpriseKey(encryptedCert.EmpresaID, encryptedCert.RUC)
    
    // 3. Configurar desencriptación AES-GCM
    block, err := aes.NewCipher(empresaKey)
    if err != nil {
        return nil, fmt.Errorf("error creando cipher: %v", err)
    }
    
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("error creando GCM: %v", err)
    }
    
    // 4. Desencriptar y verificar autenticidad
    compressedData, err := aesgcm.Open(nil, encryptedCert.IV, 
        encryptedCert.EncryptedData, encryptedCert.AAD)
    if err != nil {
        csm.auditLogger.LogCertificateAccess(encryptedCert.EmpresaID, userID, "FAILED_DECRYPT")
        return nil, fmt.Errorf("error desencriptando: certificado corrupto o clave incorrecta")
    }
    
    // 5. Descomprimir datos
    originalData, err := decompressData(compressedData)
    if err != nil {
        return nil, fmt.Errorf("error descomprimiendo datos: %v", err)
    }
    
    // 6. Verificar integridad
    currentHash := sha512.Sum512(originalData)
    expectedHash, err := hex.DecodeString(encryptedCert.OriginalHash)
    if err != nil || !bytes.Equal(currentHash[:], expectedHash) {
        csm.auditLogger.LogCertificateAccess(encryptedCert.EmpresaID, userID, "INTEGRITY_FAILED")
        return nil, fmt.Errorf("verificación de integridad fallida")
    }
    
    // 7. Registrar acceso exitoso
    csm.auditLogger.LogCertificateAccess(encryptedCert.EmpresaID, userID, "SUCCESS")
    
    return originalData, nil
}
```

---

## 🔑 Gestión de Claves

### **Rotación Automática de Claves**

```go
// Sistema de rotación de claves maestras
type KeyRotationManager struct {
    currentKey    []byte
    previousKey   []byte
    rotationSchedule time.Duration
    notificationWindow time.Duration
    db            *database.Connection
}

func (krm *KeyRotationManager) RotateKeys() error {
    // 1. Generar nueva clave maestra
    newKey := make([]byte, 32) // 256 bits
    if _, err := rand.Read(newKey); err != nil {
        return fmt.Errorf("error generando nueva clave: %v", err)
    }
    
    // 2. Backup de clave actual
    krm.previousKey = krm.currentKey
    
    // 3. Activar nueva clave
    krm.currentKey = newKey
    
    // 4. Guardar claves de forma segura
    err := krm.saveKeysSecurely()
    if err != nil {
        // Rollback en caso de error
        krm.currentKey = krm.previousKey
        return fmt.Errorf("error guardando claves: %v", err)
    }
    
    // 5. Programar re-encriptación de certificados
    go krm.reencryptAllCertificates()
    
    // 6. Notificar a administradores
    krm.notifyKeyRotation()
    
    // 7. Programar eliminación de clave anterior (después de re-encriptación)
    time.AfterFunc(24*time.Hour, func() {
        krm.cleanupPreviousKey()
    })
    
    return nil
}

// Re-encriptación gradual de todos los certificados
func (krm *KeyRotationManager) reencryptAllCertificates() {
    // Obtener lista de todas las empresas
    empresas, err := krm.db.GetAllActiveEmpresas()
    if err != nil {
        log.Printf("Error obteniendo lista de empresas: %v", err)
        return
    }
    
    // Re-encriptar de forma gradual (no bloquear sistema)
    for _, empresa := range empresas {
        // Pausa entre re-encriptaciones para no sobrecargar
        time.Sleep(100 * time.Millisecond)
        
        err := krm.reencryptEmpresaCertificate(empresa.ID)
        if err != nil {
            log.Printf("Error re-encriptando certificado empresa %s: %v", 
                empresa.ID, err)
            // Continuar con la siguiente empresa
        }
    }
    
    log.Printf("Re-encriptación completa de %d certificados", len(empresas))
}
```

### **Hardware Security Module (HSM) - Nivel Enterprise**

```go
// Para clientes enterprise con máxima seguridad
type HSMManager struct {
    hsmClient     *pkcs11.Client
    slotID        uint
    keyLabel      string
    pin           string
}

func (hsm *HSMManager) GenerateAndStoreKey(keyID string) error {
    // 1. Conectar al HSM
    session, err := hsm.hsmClient.OpenSession(hsm.slotID, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
    if err != nil {
        return fmt.Errorf("error abriendo sesión HSM: %v", err)
    }
    defer hsm.hsmClient.CloseSession(session)
    
    // 2. Autenticar con PIN
    err = hsm.hsmClient.Login(session, pkcs11.CKU_USER, hsm.pin)
    if err != nil {
        return fmt.Errorf("error autenticando HSM: %v", err)
    }
    
    // 3. Generar clave AES-256 dentro del HSM
    keyTemplate := []*pkcs11.Attribute{
        pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
        pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
        pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, 32), // 256 bits
        pkcs11.NewAttribute(pkcs11.CKA_LABEL, keyID),
        pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
        pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
        pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false), // Clave no extraíble
        pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
    }
    
    keyHandle, err := hsm.hsmClient.GenerateKey(session, 
        []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)}, 
        keyTemplate)
    if err != nil {
        return fmt.Errorf("error generando clave en HSM: %v", err)
    }
    
    log.Printf("Clave generada en HSM con handle: %d", keyHandle)
    return nil
}

// Encriptación usando HSM (máxima seguridad)
func (hsm *HSMManager) EncryptWithHSM(data []byte, keyID string) ([]byte, error) {
    session, err := hsm.hsmClient.OpenSession(hsm.slotID, pkcs11.CKF_SERIAL_SESSION)
    if err != nil {
        return nil, fmt.Errorf("error abriendo sesión HSM: %v", err)
    }
    defer hsm.hsmClient.CloseSession(session)
    
    // Encontrar clave por etiqueta
    keyHandle, err := hsm.findKeyByLabel(session, keyID)
    if err != nil {
        return nil, fmt.Errorf("clave no encontrada en HSM: %v", err)
    }
    
    // Encriptar usando el HSM
    err = hsm.hsmClient.EncryptInit(session, 
        []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_GCM, nil)}, 
        keyHandle)
    if err != nil {
        return nil, fmt.Errorf("error inicializando encriptación HSM: %v", err)
    }
    
    encryptedData, err := hsm.hsmClient.Encrypt(session, data)
    if err != nil {
        return nil, fmt.Errorf("error encriptando con HSM: %v", err)
    }
    
    return encryptedData, nil
}
```

---

## 📊 Auditoría y Monitoreo

### **Sistema de Auditoría Completa**

```go
// Estructura para logs de auditoría de certificados
type CertificateAuditLog struct {
    ID           string    `json:"id"`
    Timestamp    time.Time `json:"timestamp"`
    EmpresaID    string    `json:"empresa_id"`
    RUC          string    `json:"ruc"`
    UserID       string    `json:"user_id"`
    Action       string    `json:"action"`       // UPLOAD, ACCESS, DELETE, ROTATE, etc.
    Result       string    `json:"result"`       // SUCCESS, FAILURE, DENIED
    IPAddress    string    `json:"ip_address"`
    UserAgent    string    `json:"user_agent"`
    Details      string    `json:"details"`
    RiskLevel    string    `json:"risk_level"`   // LOW, MEDIUM, HIGH, CRITICAL
    
    // Información técnica
    CertificateHash string `json:"certificate_hash"`
    KeyVersion      string `json:"key_version"`
    AccessMethod    string `json:"access_method"` // WEB, API, SYSTEM
    
    // Contexto de seguridad
    SessionID       string `json:"session_id"`
    RequestID       string `json:"request_id"`
    GeolocationData string `json:"geolocation"`
}

type CertificateAuditor struct {
    db          *database.Connection
    logFile     *os.File
    alertSystem *SecurityAlertSystem
    
    // Configuración de alertas
    maxFailedAttempts int
    suspiciousIPList  []string
    geoAnomalyDetection bool
}

func (ca *CertificateAuditor) LogCertificateEvent(event CertificateAuditLog) {
    // 1. Enriquecer el evento con información adicional
    event.ID = generateUUID()
    event.Timestamp = time.Now()
    
    // 2. Determinar nivel de riesgo
    event.RiskLevel = ca.calculateRiskLevel(event)
    
    // 3. Guardar en base de datos
    err := ca.saveToDB(event)
    if err != nil {
        log.Printf("Error guardando log de auditoría: %v", err)
    }
    
    // 4. Escribir a archivo de log
    ca.writeToLogFile(event)
    
    // 5. Verificar patrones sospechosos
    ca.checkSuspiciousPatterns(event)
    
    // 6. Enviar alertas si es necesario
    if event.RiskLevel == "HIGH" || event.RiskLevel == "CRITICAL" {
        ca.sendSecurityAlert(event)
    }
}

func (ca *CertificateAuditor) calculateRiskLevel(event CertificateAuditLog) string {
    riskScore := 0
    
    // Factores de riesgo
    if event.Result == "FAILURE" {
        riskScore += 2
    }
    
    if event.Action == "DELETE" || event.Action == "ROTATE" {
        riskScore += 3
    }
    
    if ca.isFromSuspiciousIP(event.IPAddress) {
        riskScore += 4
    }
    
    if ca.isOutsideBusinessHours(event.Timestamp) {
        riskScore += 1
    }
    
    if ca.isGeoAnomalous(event.IPAddress, event.EmpresaID) {
        riskScore += 3
    }
    
    // Determinar nivel
    switch {
    case riskScore >= 7:
        return "CRITICAL"
    case riskScore >= 4:
        return "HIGH"
    case riskScore >= 2:
        return "MEDIUM"
    default:
        return "LOW"
    }
}

// Detección de patrones sospechosos
func (ca *CertificateAuditor) checkSuspiciousPatterns(event CertificateAuditLog) {
    // 1. Múltiples intentos fallidos
    recentFailures := ca.countRecentFailures(event.EmpresaID, event.UserID, 15*time.Minute)
    if recentFailures >= ca.maxFailedAttempts {
        ca.alertSystem.SendAlert("MULTIPLE_FAILED_ATTEMPTS", event)
    }
    
    // 2. Acceso desde múltiples ubicaciones geográficas
    if ca.geoAnomalyDetection {
        recentLocations := ca.getRecentAccessLocations(event.EmpresaID, 1*time.Hour)
        if len(recentLocations) > 3 {
            ca.alertSystem.SendAlert("GEOGRAPHIC_ANOMALY", event)
        }
    }
    
    // 3. Acceso fuera de horario laboral
    if ca.isOutsideBusinessHours(event.Timestamp) && event.Action != "SYSTEM" {
        ca.alertSystem.SendAlert("AFTER_HOURS_ACCESS", event)
    }
    
    // 4. Cambio de certificado no programado
    if event.Action == "UPLOAD" && !ca.isScheduledUpdate(event.EmpresaID) {
        ca.alertSystem.SendAlert("UNSCHEDULED_CERTIFICATE_CHANGE", event)
    }
}
```

### **Dashboard de Seguridad en Tiempo Real**

```go
// API para dashboard de seguridad
func (ca *CertificateAuditor) GetSecurityDashboard() *SecurityDashboard {
    dashboard := &SecurityDashboard{
        LastUpdated: time.Now(),
    }
    
    // Estadísticas de las últimas 24 horas
    since24h := time.Now().Add(-24 * time.Hour)
    
    // Eventos por nivel de riesgo
    dashboard.EventsByRisk = ca.getEventCountsByRisk(since24h)
    
    // Top empresas con más actividad
    dashboard.TopActiveEmpresas = ca.getTopActiveEmpresas(since24h, 10)
    
    // Alertas activas
    dashboard.ActiveAlerts = ca.getActiveAlerts()
    
    // Certificados por vencer
    dashboard.ExpiringCertificates = ca.getCertificatesExpiringInDays(30)
    
    // Métricas de acceso
    dashboard.AccessMetrics = ca.getAccessMetrics(since24h)
    
    // Ubicaciones de acceso
    dashboard.AccessLocations = ca.getAccessLocationStats(since24h)
    
    return dashboard
}

type SecurityDashboard struct {
    LastUpdated    time.Time `json:"last_updated"`
    
    EventsByRisk   map[string]int `json:"events_by_risk"`
    // {"LOW": 150, "MEDIUM": 25, "HIGH": 5, "CRITICAL": 1}
    
    TopActiveEmpresas []struct {
        EmpresaID    string `json:"empresa_id"`
        RazonSocial  string `json:"razon_social"`
        EventCount   int    `json:"event_count"`
        LastAccess   time.Time `json:"last_access"`
    } `json:"top_active_empresas"`
    
    ActiveAlerts []struct {
        ID          string    `json:"id"`
        Type        string    `json:"type"`
        EmpresaID   string    `json:"empresa_id"`
        Severity    string    `json:"severity"`
        CreatedAt   time.Time `json:"created_at"`
        Description string    `json:"description"`
    } `json:"active_alerts"`
    
    ExpiringCertificates []struct {
        EmpresaID    string    `json:"empresa_id"`
        RazonSocial  string    `json:"razon_social"`
        ExpiryDate   time.Time `json:"expiry_date"`
        DaysRemaining int      `json:"days_remaining"`
    } `json:"expiring_certificates"`
    
    AccessMetrics struct {
        TotalAccesses    int `json:"total_accesses"`
        SuccessfulAccesses int `json:"successful_accesses"`
        FailedAccesses   int `json:"failed_accesses"`
        UniqueUsers      int `json:"unique_users"`
        UniqueIPs        int `json:"unique_ips"`
    } `json:"access_metrics"`
    
    AccessLocations []struct {
        Country     string  `json:"country"`
        City        string  `json:"city"`
        AccessCount int     `json:"access_count"`
        Percentage  float64 `json:"percentage"`
    } `json:"access_locations"`
}
```

---

## 🚨 Respuesta a Incidentes

### **Procedimientos de Emergencia**

```go
// Sistema de respuesta automática a incidentes
type IncidentResponseSystem struct {
    alertThresholds map[string]int
    responseActions map[string][]ResponseAction
    notificationChannels []NotificationChannel
    escalationMatrix []EscalationLevel
}

type SecurityIncident struct {
    ID           string
    Type         string    // BREACH_ATTEMPT, MULTIPLE_FAILURES, GEO_ANOMALY, etc.
    Severity     string    // LOW, MEDIUM, HIGH, CRITICAL
    EmpresaID    string
    UserID       string
    IPAddress    string
    DetectedAt   time.Time
    Description  string
    Evidence     []AuditLogEntry
    Status       string    // DETECTED, INVESTIGATING, CONTAINED, RESOLVED
    AssignedTo   string
}

func (irs *IncidentResponseSystem) HandleSecurityIncident(incident SecurityIncident) {
    // 1. Clasificar y priorizar el incidente
    incident.Severity = irs.classifyIncident(incident)
    
    // 2. Ejecutar respuesta automática inmediata
    irs.executeImmediateResponse(incident)
    
    // 3. Notificar a equipos apropiados
    irs.notifySecurityTeam(incident)
    
    // 4. Documentar el incidente
    irs.documentIncident(incident)
    
    // 5. Iniciar investigación si es necesario
    if incident.Severity == "HIGH" || incident.Severity == "CRITICAL" {
        irs.initiateInvestigation(incident)
    }
}

func (irs *IncidentResponseSystem) executeImmediateResponse(incident SecurityIncident) {
    switch incident.Type {
    case "MULTIPLE_FAILED_ATTEMPTS":
        // Bloquear temporalmente la cuenta
        irs.temporaryAccountLockout(incident.UserID, 30*time.Minute)
        
        // Bloquear IP si múltiples cuentas afectadas
        if irs.countAffectedAccounts(incident.IPAddress) > 3 {
            irs.blockIPAddress(incident.IPAddress, 24*time.Hour)
        }
        
    case "CERTIFICATE_BREACH_ATTEMPT":
        // Suspender acceso inmediato al certificado
        irs.suspendCertificateAccess(incident.EmpresaID)
        
        // Notificar a la empresa afectada
        irs.notifyAffectedEnterprise(incident.EmpresaID, incident)
        
        // Forzar rotación de claves
        irs.forceKeyRotation(incident.EmpresaID)
        
    case "GEOGRAPHIC_ANOMALY":
        // Requerir autenticación adicional
        irs.requireAdditionalAuth(incident.UserID, 24*time.Hour)
        
        // Notificar al usuario por canales alternativos
        irs.sendOutOfBandNotification(incident.UserID, incident)
        
    case "INTEGRITY_VIOLATION":
        // Aislar certificado comprometido inmediatamente
        irs.quarantineCertificate(incident.EmpresaID)
        
        // Activar protocolo de recuperación
        irs.activateRecoveryProtocol(incident.EmpresaID)
        
        // Escalar a CRITICAL automáticamente
        incident.Severity = "CRITICAL"
    }
}

// Protocolo de recuperación para certificados comprometidos
func (irs *IncidentResponseSystem) activateRecoveryProtocol(empresaID string) {
    // 1. Crear backup de emergencia del estado actual
    backupID := irs.createEmergencyBackup(empresaID)
    
    // 2. Revocar acceso inmediato
    irs.revokeAllAccess(empresaID)
    
    // 3. Notificar automáticamente a SRI (si está configurado)
    irs.notifySRIOfCompromise(empresaID)
    
    // 4. Preparar certificado de emergencia temporal
    tempCert := irs.generateTemporaryCertificate(empresaID)
    
    // 5. Crear plan de recuperación
    recoveryPlan := RecoveryPlan{
        EmpresaID:           empresaID,
        BackupID:            backupID,
        TemporaryCertID:     tempCert.ID,
        EstimatedDowntime:   "2-6 horas",
        RequiredActions:     irs.generateRecoveryActions(empresaID),
        ApprovalRequired:    true,
        CreatedAt:          time.Now(),
    }
    
    irs.saveRecoveryPlan(recoveryPlan)
    
    // 6. Notificar a todos los stakeholders
    irs.notifyAllStakeholders(empresaID, recoveryPlan)
}
```

### **Comunicación de Crisis**

```go
// Sistema de comunicación durante incidentes de seguridad
type CrisisCommunicationManager struct {
    templates         map[string]MessageTemplate
    notificationMatrix map[string][]NotificationChannel
    legalRequirements LegalNotificationRequirements
}

func (ccm *CrisisCommunicationManager) NotifySecurityBreach(incident SecurityIncident) {
    // 1. Determinar audiencias que deben ser notificadas
    audiences := ccm.determineNotificationAudiences(incident)
    
    // 2. Preparar mensajes personalizados por audiencia
    for _, audience := range audiences {
        message := ccm.prepareMessage(incident, audience)
        channels := ccm.notificationMatrix[audience]
        
        // Enviar por todos los canales apropiados
        for _, channel := range channels {
            ccm.sendNotification(message, channel)
        }
    }
    
    // 3. Cumplir requisitos legales de notificación
    ccm.handleLegalNotifications(incident)
}

type MessageTemplate struct {
    Subject     string
    BodyHTML    string
    BodyText    string
    Urgency     string
    Language    string
    Audience    string
    Channel     string
}

// Plantilla para notificación a empresa afectada
var EnterpriseBreachNotificationTemplate = MessageTemplate{
    Subject: "URGENTE: Actividad sospechosa detectada en su certificado digital",
    BodyHTML: `
    <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; margin: 10px 0;">
        <h2 style="color: #856404;">⚠️ ALERTA DE SEGURIDAD</h2>
        <p><strong>Estimado cliente,</strong></p>
        
        <p>Hemos detectado actividad sospechosa relacionada con el certificado digital de su empresa 
        <strong>{{.RazonSocial}}</strong> (RUC: {{.RUC}}).</p>
        
        <h3>📋 Detalles del Incidente:</h3>
        <ul>
            <li><strong>Fecha/Hora:</strong> {{.DetectedAt}}</li>
            <li><strong>Tipo:</strong> {{.IncidentType}}</li>
            <li><strong>Severidad:</strong> {{.Severity}}</li>
            <li><strong>IP de Origen:</strong> {{.IPAddress}}</li>
        </ul>
        
        <h3>🔒 Medidas Tomadas Automáticamente:</h3>
        <ul>
            <li>✅ Acceso al certificado suspendido temporalmente</li>
            <li>✅ Notificación a nuestro equipo de seguridad</li>
            <li>✅ Inicio de investigación del incidente</li>
            <li>✅ Activación de protocolos de seguridad</li>
        </ul>
        
        <h3>📞 Qué Debe Hacer INMEDIATAMENTE:</h3>
        <ol>
            <li><strong>Llamar a nuestra línea de emergencia:</strong> +593 2 XXX-XXXX</li>
            <li><strong>Verificar identidad:</strong> Prepare su RUC y datos de contacto</li>
            <li><strong>No intente acceder al sistema</strong> hasta que se resuelva el incidente</li>
            <li><strong>Revisar actividad reciente</strong> en su empresa relacionada con facturación</li>
        </ol>
        
        <div style="background-color: #f8d7da; border: 1px solid #f5c6cb; padding: 10px; margin: 15px 0;">
            <p><strong>⏰ TIEMPO CRÍTICO:</strong> Para minimizar el impacto, debe contactarnos 
            dentro de las próximas <strong>2 horas</strong>.</p>
        </div>
        
        <p>Este incidente será resuelto con la máxima prioridad. Su seguridad y continuidad operativa 
        son nuestra principal preocupación.</p>
        
        <p><strong>Equipo de Seguridad</strong><br>
        Tu Empresa SaaS<br>
        📧 seguridad@tu-dominio.com<br>
        📞 +593 2 XXX-XXXX (24/7)</p>
    </div>
    `,
    Urgency: "HIGH",
    Audience: "AFFECTED_ENTERPRISE",
    Channel: "EMAIL",
}
```

---

## 🔍 Compliance y Regulaciones

### **Cumplimiento Legal Ecuador**

```go
// Gestor de cumplimiento legal para Ecuador
type ComplianceManager struct {
    regulationsDB map[string]Regulation
    auditTrail    []ComplianceEvent
    certifications []Certification
}

type Regulation struct {
    ID           string
    Name         string
    Authority    string    // SRI, SUPERBANCOS, etc.
    Type         string    // ENCRYPTION, AUDIT, RETENTION, etc.
    Requirements []string
    Penalties    []Penalty
    LastUpdated  time.Time
}

// Regulaciones aplicables en Ecuador
func (cm *ComplianceManager) InitializeEcuadorianRegulations() {
    // Ley de Comercio Electrónico, Firmas Electrónicas y Mensajes de Datos
    cm.regulationsDB["LEY_COMERCIO_ELECTRONICO"] = Regulation{
        ID:        "LEY_COMERCIO_ELECTRONICO",
        Name:      "Ley de Comercio Electrónico, Firmas Electrónicas y Mensajes de Datos",
        Authority: "MINTEL",
        Type:      "DIGITAL_SIGNATURE",
        Requirements: []string{
            "Certificados digitales deben ser emitidos por entidades autorizadas",
            "Algoritmos de encriptación deben cumplir estándares internacionales",
            "Logs de auditoría deben mantenerse por mínimo 7 años",
            "Integridad de documentos debe ser verificable",
            "No repudio debe estar garantizado",
        },
        LastUpdated: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
    }
    
    // Regulaciones SRI
    cm.regulationsDB["SRI_FACTURACION_ELECTRONICA"] = Regulation{
        ID:        "SRI_FACTURACION_ELECTRONICA",
        Name:      "Resolución NAC-DGERCGC18-00000434",
        Authority: "SRI",
        Type:      "TAX_COMPLIANCE",
        Requirements: []string{
            "Certificados digitales de entidades autorizadas por BCE",
            "Conservación de comprobantes por 7 años",
            "Disponibilidad de documentos para auditoría SRI",
            "Backup y recuperación de información tributaria",
            "Medidas de seguridad para integridad fiscal",
        },
        LastUpdated: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
    }
    
    // Protección de Datos Personales
    cm.regulationsDB["PROTECCION_DATOS"] = Regulation{
        ID:        "PROTECCION_DATOS",
        Name:      "Ley Orgánica de Protección de Datos Personales",
        Authority: "DINARDAP",
        Type:      "DATA_PROTECTION",
        Requirements: []string{
            "Consentimiento explícito para procesamiento de datos",
            "Encriptación de datos personales en reposo y tránsito",
            "Notificación de brechas en máximo 72 horas",
            "Derecho al olvido y portabilidad de datos",
            "Evaluaciones de impacto de privacidad",
        },
        LastUpdated: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
    }
}

// Verificación automática de cumplimiento
func (cm *ComplianceManager) VerifyCompliance() *ComplianceReport {
    report := &ComplianceReport{
        GeneratedAt: time.Now(),
        Status:      "COMPLIANT",
        Issues:      []ComplianceIssue{},
        Score:       100,
    }
    
    // Verificar cada regulación
    for _, regulation := range cm.regulationsDB {
        for _, requirement := range regulation.Requirements {
            compliant := cm.checkRequirement(regulation.ID, requirement)
            if !compliant {
                issue := ComplianceIssue{
                    RegulationID: regulation.ID,
                    Requirement:  requirement,
                    Severity:     cm.getSeverity(regulation.ID, requirement),
                    Description:  cm.getIssueDescription(regulation.ID, requirement),
                    Remediation:  cm.getRemediationSteps(regulation.ID, requirement),
                }
                report.Issues = append(report.Issues, issue)
                report.Score -= cm.getScoreImpact(issue.Severity)
            }
        }
    }
    
    if len(report.Issues) > 0 {
        report.Status = "NON_COMPLIANT"
    }
    
    return report
}
```

### **Auditorías Automáticas**

```go
// Sistema de auditorías automáticas
type AutomaticAuditSystem struct {
    auditSchedules []AuditSchedule
    auditTemplates map[string]AuditTemplate
    compliance     *ComplianceManager
    reporting      *AuditReportingSystem
}

type AuditSchedule struct {
    ID          string
    Name        string
    Type        string        // SECURITY, COMPLIANCE, OPERATIONAL
    Frequency   time.Duration // DAILY, WEEKLY, MONTHLY, QUARTERLY
    LastRun     time.Time
    NextRun     time.Time
    Enabled     bool
    Template    string
}

func (aas *AutomaticAuditSystem) InitializeAuditSchedules() {
    aas.auditSchedules = []AuditSchedule{
        {
            ID:        "DAILY_SECURITY_AUDIT",
            Name:      "Auditoría Diaria de Seguridad",
            Type:      "SECURITY",
            Frequency: 24 * time.Hour,
            NextRun:   time.Now().Add(24 * time.Hour),
            Enabled:   true,
            Template:  "SECURITY_TEMPLATE",
        },
        {
            ID:        "WEEKLY_COMPLIANCE_AUDIT",
            Name:      "Auditoría Semanal de Cumplimiento",
            Type:      "COMPLIANCE",
            Frequency: 7 * 24 * time.Hour,
            NextRun:   time.Now().Add(7 * 24 * time.Hour),
            Enabled:   true,
            Template:  "COMPLIANCE_TEMPLATE",
        },
        {
            ID:        "MONTHLY_CERTIFICATE_AUDIT",
            Name:      "Auditoría Mensual de Certificados",
            Type:      "OPERATIONAL",
            Frequency: 30 * 24 * time.Hour,
            NextRun:   time.Now().Add(30 * 24 * time.Hour),
            Enabled:   true,
            Template:  "CERTIFICATE_TEMPLATE",
        },
    }
}

func (aas *AutomaticAuditSystem) RunScheduledAudits() {
    for _, schedule := range aas.auditSchedules {
        if schedule.Enabled && time.Now().After(schedule.NextRun) {
            go aas.executeAudit(schedule)
        }
    }
}

func (aas *AutomaticAuditSystem) executeAudit(schedule AuditSchedule) {
    auditRun := AuditRun{
        ID:          generateUUID(),
        ScheduleID:  schedule.ID,
        StartTime:   time.Now(),
        Status:      "RUNNING",
        Type:        schedule.Type,
    }
    
    // Ejecutar checks específicos según el tipo de auditoría
    switch schedule.Type {
    case "SECURITY":
        auditRun.Results = aas.runSecurityAudit()
    case "COMPLIANCE":
        auditRun.Results = aas.runComplianceAudit()
    case "OPERATIONAL":
        auditRun.Results = aas.runOperationalAudit()
    }
    
    auditRun.EndTime = time.Now()
    auditRun.Duration = auditRun.EndTime.Sub(auditRun.StartTime)
    auditRun.Status = "COMPLETED"
    
    // Generar reporte
    report := aas.reporting.GenerateReport(auditRun)
    
    // Enviar notificaciones si hay problemas
    if auditRun.HasCriticalIssues() {
        aas.sendCriticalIssueAlert(auditRun, report)
    }
    
    // Actualizar próxima ejecución
    aas.updateNextRun(schedule.ID)
    
    log.Printf("Auditoría %s completada en %v", schedule.Name, auditRun.Duration)
}
```

---

## 🎯 Conclusión

Este sistema de seguridad para certificados digitales proporciona:

### ✅ **Protección de Clase Empresarial**
- **Encriptación AES-256-GCM** con claves derivadas por empresa
- **Aislamiento completo** entre tenants
- **Hardware Security Module** para clientes enterprise
- **Rotación automática** de claves maestras

### ✅ **Auditoría y Cumplimiento**
- **Logs completos** de todos los accesos
- **Alertas en tiempo real** para actividad sospechosa
- **Cumplimiento regulaciones** ecuatorianas e internacionales
- **Auditorías automáticas** programadas

### ✅ **Respuesta a Incidentes**
- **Detección automática** de patrones sospechosos
- **Respuesta inmediata** a amenazas
- **Protocolos de recuperación** bien definidos
- **Comunicación de crisis** automatizada

### ✅ **Monitoreo Continuo**
- **Dashboard de seguridad** en tiempo real
- **Métricas de acceso** y comportamiento
- **Alertas proactivas** de certificados por vencer
- **Análisis geográfico** de accesos

### 🛡️ **Resultado Final**

Tu sistema SaaS maneja certificados digitales con el mismo nivel de seguridad que los bancos más grandes del mundo, garantizando:

- **Confianza total** de tus clientes empresariales
- **Cumplimiento legal** en Ecuador y regulaciones internacionales  
- **Protección contra amenazas** conocidas y emergentes
- **Escalabilidad segura** para miles de empresas

**¡Tus clientes pueden confiar completamente en la seguridad de sus activos digitales más valiosos! 🔐**