package router

import (
	"crescendo-api/config/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewRouter(container *app.Container) *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Service Active",
		})
	})

	router.POST("/users", container.User.Register)
	router.POST("/users/login", container.User.Login)

	router.GET("/songs/search", container.Song.SearchSongs)
	router.GET("/songs/:id/playback", container.Song.GetSongPlaybackInfo)

	router.GET("/albums/search", container.Album.SearchAlbums)
	router.GET("/albums/:id", container.Album.GetAlbumDetails)

	router.GET("/artists/search", container.Artist.SearchArtists)
	router.GET("/artists/:id", container.Artist.GetArtist)
	router.GET("/artists/:id/albums", container.Artist.GetArtistAlbumPreviews)
	router.GET("/artists/:id/songs", container.Artist.GetArtistSongPreviews)

	return router
}
