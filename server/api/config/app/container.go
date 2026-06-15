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
	User   *controllers.UserController
	Song   *controllers.SongController
	Album  *controllers.AlbumController
	Artist *controllers.ArtistController
}

func NewContainer(db database.DBTX) *Container {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	songRepo := repositories.NewSongRepository(db)
	songService := services.NewSongService(songRepo)
	sq, _ := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	idEncoder := security.NewSquidEncoder(sq)
	songController := controllers.NewSongController(songService, idEncoder)

	cache := database.NewCacheConnection()

	albumRepo := repositories.NewAlbumRepository(db)
	genreRepo := repositories.NewGenreRepository(db)
	albumService := services.NewAlbumService(albumRepo, genreRepo, cache)
	albumController := controllers.NewAlbumController(albumService, idEncoder)

	artistRepo := repositories.NewArtistRepository(db)
	artistService := services.NewArtistService(artistRepo)
	artistController := controllers.NewArtistController(artistService, idEncoder)

	return &Container{
		User:   userController,
		Song:   songController,
		Album:  albumController,
		Artist: artistController,
	}
}
