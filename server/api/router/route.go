package router

import (
	"crescendo-api/config/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(container *app.Container) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Service Active",
		})
	})

	router.POST("/users", container.User.Register)
	router.POST("/users/login", container.User.Login)

	router.GET("/songs/:id/playback", container.Song.GetSongPlaybackInfo)

	router.GET("/albums/:id", container.Album.GetAlbumDetails)

	return router
}
