package router

import (
	"crescendo-api/config/app"
	"crescendo-api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(container *app.Container) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Service Active",
		})
	})

	router.POST("/users", container.User.Register)
	router.POST("/users/login", container.User.Login)

	api := router.Group("/")
	api.Use(middleware.Authentication())
	api.GET("/songs/:id/playback", container.Song.GetSongPlaybackInfo)

	api.GET("/albums/:id", container.Album.GetAlbumDetails)

	api.GET("/artists", container.Artist.GetAllArtist)
	api.GET("/artists/:id", container.Artist.GetArtist)
	api.GET("/artists/:id/albums", container.Artist.GetArtistAlbumPreviews)
	api.GET("/artists/:id/songs", container.Artist.GetArtistSongPreviews)

	api.GET("/search", container.Search.SearchByName)

	return router
}
