// Package validators contiene todas las funciones de validación
package validators

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
	
	"go-facturacion-sri/models" // Import local!
)

// ValidarRUC - Valida que un RUC ecuatoriano sea correcto
// RUC puede ser de persona natural (13 dígitos) o jurídica (13 dígitos)
// Devuelve un error si el RUC no es válido
func ValidarRUC(ruc string) error {
	// Sanitizar input
	ruc = strings.TrimSpace(ruc)
	
	// Verificar longitud
	if len(ruc) != 13 {
		return errors.New("el RUC debe tener exactamente 13 dígitos")
	}
	
	// Verificar que todos sean números
	for _, char := range ruc {
		if char < '0' || char > '9' {
			return errors.New("el RUC solo puede contener números")
		}
	}
	
	// Verificar provincia (primeros 2 dígitos)
	provincia, err := strconv.Atoi(ruc[:2])
	if err != nil {
		return errors.New("error al procesar los primeros dos dígitos del RUC")
	}
	
	if provincia < 1 || provincia > 24 {
		return errors.New("los dos primeros dígitos del RUC deben estar entre 01 y 24")
	}
	
	// Verificar tercer dígito (tipo de RUC)
	tercerDigito, _ := strconv.Atoi(string(ruc[2]))
	
	// RUC de persona natural: tercer dígito debe ser menor a 6
	// RUC de empresa privada: tercer dígito debe ser 9
	// RUC de institución pública: tercer dígito debe ser 6
	if tercerDigito >= 6 && tercerDigito != 6 && tercerDigito != 9 {
		return errors.New("el tercer dígito del RUC debe ser menor a 6, o igual a 6 (sector público) o 9 (empresa privada)")
	}
	
	// Validar dígito verificador según tipo de RUC
	if tercerDigito < 6 {
		// RUC persona natural - usar algoritmo de cédula
		return validarRUCPersonaNatural(ruc)
	} else if tercerDigito == 6 {
		// RUC sector público
		return validarRUCSectorPublico(ruc)
	} else if tercerDigito == 9 {
		// RUC empresa privada
		return validarRUCEmpresaPrivada(ruc)
	}
	
	return nil
}

// validarRUCPersonaNatural valida RUC de persona natural (similar a cédula)
func validarRUCPersonaNatural(ruc string) error {
	// Los primeros 10 dígitos deben cumplir algoritmo de cédula
	cedulaParte := ruc[:10]
	if err := ValidarCedula(cedulaParte); err != nil {
		return fmt.Errorf("RUC persona natural inválido: %v", err)
	}
	
	// Los últimos 3 dígitos deben ser "001"
	if ruc[10:] != "001" {
		return errors.New("RUC persona natural debe terminar en 001")
	}
	
	return nil
}

// validarRUCSectorPublico valida RUC de sector público
func validarRUCSectorPublico(ruc string) error {
	coeficientes := []int{3, 2, 7, 6, 5, 4, 3, 2}
	suma := 0
	
	for i := 0; i < 8; i++ {
		digito, _ := strconv.Atoi(string(ruc[i]))
		suma += digito * coeficientes[i]
	}
	
	resto := suma % 11
	digitoVerificador := 11 - resto
	
	if resto == 0 {
		digitoVerificador = 0
	} else if resto == 1 {
		return errors.New("RUC sector público inválido por algoritmo")
	}
	
	nonoDigito, _ := strconv.Atoi(string(ruc[8]))
	if digitoVerificador != nonoDigito {
		return errors.New("el dígito verificador del RUC sector público no es válido")
	}
	
	// Los últimos 4 dígitos deben ser "0001"
	if ruc[9:] != "0001" {
		return errors.New("RUC sector público debe terminar en 0001")
	}
	
	return nil
}

// validarRUCEmpresaPrivada valida RUC de empresa privada
func validarRUCEmpresaPrivada(ruc string) error {
	coeficientes := []int{4, 3, 2, 7, 6, 5, 4, 3, 2}
	suma := 0
	
	for i := 0; i < 9; i++ {
		digito, _ := strconv.Atoi(string(ruc[i]))
		suma += digito * coeficientes[i]
	}
	
	resto := suma % 11
	digitoVerificador := 11 - resto
	
	if resto == 0 {
		digitoVerificador = 0
	} else if resto == 1 {
		return errors.New("RUC empresa privada inválido por algoritmo")
	}
	
	decimoDigito, _ := strconv.Atoi(string(ruc[9]))
	if digitoVerificador != decimoDigito {
		return errors.New("el dígito verificador del RUC empresa privada no es válido")
	}
	
	// Los últimos 3 dígitos deben ser "001"
	if ruc[10:] != "001" {
		return errors.New("RUC empresa privada debe terminar en 001")
	}
	
	return nil
}

// ValidarCedula - Valida que una cédula ecuatoriana sea correcta
// Devuelve un error si la cédula no es válida
func ValidarCedula(cedula string) error {
	// Sanitizar input
	cedula = strings.TrimSpace(cedula)
	// Verificar longitud
	if len(cedula) != 10 {
		return errors.New("la cédula debe tener exactamente 10 dígitos")
	}
	
	// Verificar que no esté vacía después de sanitizar
	if cedula == "" {
		return errors.New("la cédula no puede estar vacía")
	}
	
	// Verificar que todos sean números
	for _, char := range cedula {
		if char < '0' || char > '9' {
			return errors.New("la cédula solo puede contener números")
		}
	}
	
	// Verificar que los dos primeros dígitos sean válidos (01-24)
	provincia, err := strconv.Atoi(cedula[:2])
	if err != nil {
		return errors.New("error al procesar los primeros dos dígitos de la cédula")
	}
	
	if provincia < 1 || provincia > 24 {
		return errors.New("los dos primeros dígitos de la cédula deben estar entre 01 y 24")
	}
	
	// Algoritmo de validación del dígito verificador
	coeficientes := []int{2, 1, 2, 1, 2, 1, 2, 1, 2}
	suma := 0
	
	for i := 0; i < 9; i++ {
		digito, _ := strconv.Atoi(string(cedula[i]))
		resultado := digito * coeficientes[i]
		
		if resultado >= 10 {
			resultado = resultado - 9
		}
		
		suma += resultado
	}
	
	digitoVerificador := suma % 10
	if digitoVerificador != 0 {
		digitoVerificador = 10 - digitoVerificador
	}
	
	ultimoDigito, _ := strconv.Atoi(string(cedula[9]))
	
	if digitoVerificador != ultimoDigito {
		return errors.New("el dígito verificador de la cédula no es válido")
	}
	
	return nil // nil significa "no hay error"
}

// SanitizarTexto - Sanitiza texto para prevenir inyección XML/XSS
func SanitizarTexto(texto string) string {
	// Remover espacios al inicio y final
	texto = strings.TrimSpace(texto)
	
	// Escapar HTML para prevenir XSS
	texto = html.EscapeString(texto)
	
	// Remover caracteres de control peligrosos
	controlChars := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
	texto = controlChars.ReplaceAllString(texto, "")
	
	// Limitar longitud máxima para prevenir ataques DoS
	if len(texto) > 1000 {
		texto = texto[:1000]
	}
	
	return texto
}

// ValidarLimitesExtremos - Valida límites extremos para prevenir ataques
func ValidarLimitesExtremos(cantidad, precio float64) error {
	// Validar cantidad extrema
	if cantidad > 999999 {
		return errors.New("cantidad excede el límite máximo permitido (999,999)")
	}
	
	// Validar precio extremo
	if precio > 9999999.99 {
		return errors.New("precio excede el límite máximo permitido ($9,999,999.99)")
	}
	
	// Validar precisión decimal
	if fmt.Sprintf("%.2f", precio) != fmt.Sprintf("%.2f", float64(int(precio*100))/100) {
		return errors.New("precio tiene más de 2 decimales")
	}
	
	return nil
}

// ValidarFecha - Valida que una fecha esté en rango aceptable
func ValidarFecha(fecha string) error {
	// Parsear fecha en formato DD/MM/YYYY
	t, err := time.Parse("02/01/2006", fecha)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido, use DD/MM/YYYY: %v", err)
	}
	
	// Validar que no sea muy antigua (más de 5 años)
	limiteAntiguo := time.Now().AddDate(-5, 0, 0)
	if t.Before(limiteAntiguo) {
		return errors.New("la fecha no puede ser mayor a 5 años en el pasado")
	}
	
	// Validar que no sea futura (más de 1 día)
	limiteFuturo := time.Now().AddDate(0, 0, 1)
	if t.After(limiteFuturo) {
		return errors.New("la fecha no puede ser futura")
	}
	
	return nil
}

// ValidarProducto - Valida un producto individual con sanitización
func ValidarProducto(producto models.ProductoInput) error {
	// Sanitizar y validar código de producto
	codigoSanitizado := SanitizarTexto(producto.Codigo)
	if codigoSanitizado == "" {
		return errors.New("el código del producto no puede estar vacío")
	}
	if len(codigoSanitizado) > 25 {
		return errors.New("el código del producto no puede exceder 25 caracteres")
	}
	
	// Sanitizar y validar descripción
	descripcionSanitizada := SanitizarTexto(producto.Descripcion)
	if descripcionSanitizada == "" {
		return errors.New("la descripción del producto no puede estar vacía")
	}
	if len(descripcionSanitizada) > 300 {
		return errors.New("la descripción del producto no puede exceder 300 caracteres")
	}
	
	// Validar cantidad
	if producto.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	
	// Validar precio
	if producto.PrecioUnitario <= 0 {
		return errors.New("el precio unitario debe ser mayor a cero")
	}
	
	// Validar límites extremos
	if err := ValidarLimitesExtremos(producto.Cantidad, producto.PrecioUnitario); err != nil {
		return err
	}
	
	return nil
}

// ValidarFacturaInput - Valida todos los datos de entrada con sanitización
func ValidarFacturaInput(input models.FacturaInput) error {
	// Sanitizar y validar nombre del cliente
	nombreSanitizado := SanitizarTexto(input.ClienteNombre)
	if nombreSanitizado == "" {
		return errors.New("el nombre del cliente no puede estar vacío")
	}
	if len(nombreSanitizado) > 300 {
		return errors.New("el nombre del cliente no puede exceder 300 caracteres")
	}
	
	// Validar que el nombre no contenga solo números
	soloNumeros := regexp.MustCompile(`^[0-9]+$`)
	if soloNumeros.MatchString(nombreSanitizado) {
		return errors.New("el nombre del cliente no puede contener solo números")
	}
	
	// Validar cédula
	if err := ValidarCedula(input.ClienteCedula); err != nil {
		return fmt.Errorf("cédula inválida: %v", err)
	}
	
	// Validar que tenga al menos un producto
	if len(input.Productos) == 0 {
		return errors.New("debe incluir al menos un producto")
	}
	
	// Validar cada producto usando un loop
	for i, producto := range input.Productos {
		if err := ValidarProducto(producto); err != nil {
			return fmt.Errorf("producto %d inválido: %v", i+1, err)
		}
	}
	
	return nil
}