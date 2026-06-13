package router

import (
	"crescendo-streaming/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *controllers.StreamingController) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if ctrl != nil {
		r.POST("/upload", ctrl.UploadAudio)
		r.GET("/stream/:file_id", ctrl.StreamAudio)
	}

	return r
}
