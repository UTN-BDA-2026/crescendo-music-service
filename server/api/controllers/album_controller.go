package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/security"
	"crescendo-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlbumController struct {
	service services.AlbumService
	encoder security.Encoder
}

func NewAlbumController(s services.AlbumService, e security.Encoder) *AlbumController {
	return &AlbumController{
		service: s,
		encoder: e,
	}
}

func (sc *AlbumController) GetAlbumDetails(c *gin.Context) {
	hashID := c.Param("id")
	id, err := sc.encoder.Decode(hashID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid album id",
		})
		return
	}

	albumDetails, err := sc.service.GetAlbumDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get album details",
		})
		return
	}
	albumDetailsDTO, err := mapping.AlbumDetailedToDTO(sc.encoder, albumDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get album details",
		})
		return
	}
	c.JSON(http.StatusOK, albumDetailsDTO)
}
