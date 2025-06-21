// Package validators contiene todas las funciones de validación
package validators

import (
	"errors"
	"fmt"
	"strconv"
	
	"go-facturacion-sri/models" // Import local!
)

// ValidarCedula - Valida que una cédula ecuatoriana sea correcta
// Devuelve un error si la cédula no es válida
func ValidarCedula(cedula string) error {
	// Verificar longitud
	if len(cedula) != 10 {
		return errors.New("la cédula debe tener exactamente 10 dígitos")
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

// ValidarProducto - Valida un producto individual
func ValidarProducto(producto models.ProductoInput) error {
	// Validar código de producto
	if producto.Codigo == "" {
		return errors.New("el código del producto no puede estar vacío")
	}
	
	// Validar descripción
	if producto.Descripcion == "" {
		return errors.New("la descripción del producto no puede estar vacía")
	}
	
	// Validar cantidad
	if producto.Cantidad <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	
	// Validar precio
	if producto.PrecioUnitario <= 0 {
		return errors.New("el precio unitario debe ser mayor a cero")
	}
	
	return nil
}

// ValidarFacturaInput - Valida todos los datos de entrada
func ValidarFacturaInput(input models.FacturaInput) error {
	// Validar nombre del cliente
	if input.ClienteNombre == "" {
		return errors.New("el nombre del cliente no puede estar vacío")
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