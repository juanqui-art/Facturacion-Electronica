package sri

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCargarCertificado tests certificate loading functionality
func TestCargarCertificado(t *testing.T) {
	tests := []struct {
		name       string
		rutaP12    string
		password   string
		expectErr  bool
		setupFile  bool
		fileContent string
	}{
		{
			name:      "Certificate file does not exist",
			rutaP12:   "/nonexistent/path/cert.p12",
			password:  "password",
			expectErr: true,
			setupFile: false,
		},
		{
			name:      "Empty certificate path",
			rutaP12:   "",
			password:  "password",
			expectErr: true,
			setupFile: false,
		},
		{
			name:        "Invalid certificate file content",
			rutaP12:     "test_invalid_cert.p12",
			password:    "password",
			expectErr:   true,
			setupFile:   true,
			fileContent: "invalid certificate content",
		},
		{
			name:      "Empty password",
			rutaP12:   "test_cert.p12",
			password:  "",
			expectErr: true,
			setupFile: true,
			fileContent: "mock certificate content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file if needed
			if tt.setupFile {
				err := os.WriteFile(tt.rutaP12, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer os.Remove(tt.rutaP12)
			}

			config := CertificadoConfig{
				RutaArchivo: tt.rutaP12,
				Password:    tt.password,
			}
			cert, err := CargarCertificado(config)

			if tt.expectErr {
				if err == nil {
					t.Error("CargarCertificado() expected error, got nil")
				}
				if cert != nil {
					t.Error("CargarCertificado() expected nil certificate on error")
				}
			} else {
				if err != nil {
					t.Errorf("CargarCertificado() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestCargarCertificadoWithValidStructure tests certificate loading with valid file structure
func TestCargarCertificadoWithValidStructure(t *testing.T) {
	// Create a temporary directory for test certificates
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with various file extensions
	extensions := []string{".p12", ".pfx", ".P12", ".PFX"}
	
	for _, ext := range extensions {
		t.Run("Extension_"+ext, func(t *testing.T) {
			certPath := filepath.Join(tempDir, "test_cert"+ext)
			
			// Create a dummy file (not a real certificate)
			err := os.WriteFile(certPath, []byte("dummy certificate content for testing"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test certificate file: %v", err)
			}

			// Test loading - should fail because it's not a real P12 file
			config := CertificadoConfig{
				RutaArchivo: certPath,
				Password:    "test_password",
			}
			cert, err := CargarCertificado(config)
			
			// We expect an error because it's not a real P12 certificate
			if err == nil {
				t.Error("Expected error for dummy certificate file")
			}
			
			if cert != nil {
				t.Error("Certificate should be nil when loading fails")
			}
		})
	}
}

// TestCargarCertificadoPathValidation tests path validation
func TestCargarCertificadoPathValidation(t *testing.T) {
	invalidPaths := []string{
		"",
		" ",
		"\n",
		"\t",
		".",
		"..",
		"/",
		"~/nonexistent.p12",
		"relative/path/cert.p12",
	}

	for _, path := range invalidPaths {
		t.Run("InvalidPath_"+path, func(t *testing.T) {
			config := CertificadoConfig{
				RutaArchivo: path,
				Password:    "password",
			}
			cert, err := CargarCertificado(config)
			
			if err == nil {
				t.Errorf("Expected error for invalid path: %s", path)
			}
			
			if cert != nil {
				t.Errorf("Certificate should be nil for invalid path: %s", path)
			}
		})
	}
}

// TestCargarCertificadoPasswordValidation tests password validation
func TestCargarCertificadoPasswordValidation(t *testing.T) {
	// Create a temporary test file
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "test.p12")
	err = os.WriteFile(certPath, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	invalidPasswords := []string{
		"", // Empty password
		// Note: nil password would be handled by Go's type system
	}

	for i, password := range invalidPasswords {
		t.Run("InvalidPassword_"+string(rune(i+48)), func(t *testing.T) {
			config := CertificadoConfig{
				RutaArchivo: certPath,
				Password:    password,
			}
			cert, err := CargarCertificado(config)
			
			// Should get an error for invalid passwords
			if err == nil {
				t.Error("Expected error for invalid password")
			}
			
			if cert != nil {
				t.Error("Certificate should be nil for invalid password")
			}
		})
	}
}

// TestCargarCertificadoFilePermissions tests file permission handling
func TestCargarCertificadoFilePermissions(t *testing.T) {
	// Skip this test on Windows as it doesn't handle Unix permissions the same way
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping file permission test on Windows")
	}

	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "no_read_cert.p12")
	
	// Create file and remove read permissions
	err = os.WriteFile(certPath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remove read permissions
	err = os.Chmod(certPath, 0000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	// Restore permissions for cleanup
	defer os.Chmod(certPath, 0644)

	config := CertificadoConfig{
		RutaArchivo: certPath,
		Password:    "password",
	}
	cert, err := CargarCertificado(config)
	
	if err == nil {
		t.Error("Expected error for unreadable file")
	}
	
	if cert != nil {
		t.Error("Certificate should be nil for unreadable file")
	}
}

// TestCargarCertificadoConcurrency tests concurrent certificate loading
func TestCargarCertificadoConcurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cert_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "concurrent_test.p12")
	err = os.WriteFile(certPath, []byte("dummy content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test concurrent access to the same file
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			config := CertificadoConfig{
				RutaArchivo: certPath,
				Password:    "password",
			}
			_, err := CargarCertificado(config)
			// We expect an error since it's not a real P12 file
			if err == nil {
				t.Error("Expected error for dummy certificate")
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent operations")
		}
	}
}

// BenchmarkCargarCertificado benchmarks certificate loading
func BenchmarkCargarCertificado(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "cert_bench")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	certPath := filepath.Join(tempDir, "benchmark.p12")
	err = os.WriteFile(certPath, []byte("benchmark certificate content"), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Note: This will fail because it's not a real P12 file,
		// but we're measuring the performance of the loading attempt
		config := CertificadoConfig{
			RutaArchivo: certPath,
			Password:    "password",
		}
		_, _ = CargarCertificado(config)
	}
}

// TestCargarCertificadoEdgeCases tests edge cases
func TestCargarCertificadoEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (string, string, func())
		wantErr  bool
	}{
		{
			name: "Very long file path",
			setup: func() (string, string, func()) {
				longPath := string(make([]byte, 1000)) + "cert.p12"
				return longPath, "password", func() {}
			},
			wantErr: true,
		},
		{
			name: "Very long password",
			setup: func() (string, string, func()) {
				tempDir, _ := os.MkdirTemp("", "cert_test")
				certPath := filepath.Join(tempDir, "test.p12")
				os.WriteFile(certPath, []byte("content"), 0644)
				longPassword := string(make([]byte, 10000))
				return certPath, longPassword, func() { os.RemoveAll(tempDir) }
			},
			wantErr: true,
		},
		{
			name: "File with special characters in name",
			setup: func() (string, string, func()) {
				tempDir, _ := os.MkdirTemp("", "cert_test")
				certPath := filepath.Join(tempDir, "test-cert_file (1).p12")
				os.WriteFile(certPath, []byte("content"), 0644)
				return certPath, "password", func() { os.RemoveAll(tempDir) }
			},
			wantErr: true, // Will fail because it's not a real P12 file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, password, cleanup := tt.setup()
			defer cleanup()

			config := CertificadoConfig{
				RutaArchivo: path,
				Password:    password,
			}
			cert, err := CargarCertificado(config)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if cert != nil && tt.wantErr {
				t.Error("Expected nil certificate on error")
			}
		})
	}
}