package app

import (
	"crescendo-api/controllers"
	"crescendo-api/database"
	"crescendo-api/repositories"
	"crescendo-api/services"
)

type Container struct {
	User *controllers.UserController
}

func NewContainer(db database.DBTX) *Container {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	return &Container{
		User: userController,
	}
}
