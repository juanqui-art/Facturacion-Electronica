package main

import (
	"go-facturacion-sri/config"
	"go-facturacion-sri/sri"
)

func main() {
	// Load development configuration
	config.CargarConfiguracion("config/desarrollo.json")
	
	// Show endpoints configuration
	sri.MostrarEndpointsSRI()
}