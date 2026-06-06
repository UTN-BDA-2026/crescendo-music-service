package main

import (
	"crescendo-api/config/app"
	"crescendo-api/database"
	"crescendo-api/router"
	"fmt"
)

func main() {
	fmt.Println("Starting API Service")

	db, _ := database.NewConnection()

	container := app.NewContainer(db)
	router := router.NewRouter(container)

	router.Run(":8080")
}
