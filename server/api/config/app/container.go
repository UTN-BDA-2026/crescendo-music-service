package app

import (
	"crescendo-api/controllers"
	"crescendo-api/database"
	"crescendo-api/repositories"
	"crescendo-api/security"
	"crescendo-api/services"
	"os"

	"github.com/sqids/sqids-go"
)

type Container struct {
	User *controllers.UserController
	Song *controllers.SongController
}

func NewContainer(db database.DBTX) *Container {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	songRepo := repositories.NewSongRepository(db)
	songService := services.NewSongService(songRepo)
	sq, _ := sqids.New(sqids.Options{
		Alphabet: os.Getenv("SQID_ALPHABET"),
	})
	songEncoder := security.NewSquidEncoder(sq)
	songController := controllers.NewSongController(songService, songEncoder)
	return &Container{
		User: userController,
		Song: songController,
	}
}
