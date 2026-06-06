package main

import (
	"crescendo-streaming/router"
	"fmt"
	"log"
)

func main() {
	r := router.SetupRouter()

	fmt.Println("Iniciando servicio de streaming en el puerto 8081...")

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
