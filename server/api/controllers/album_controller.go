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

func (ac *AlbumController) SearchAlbums(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title query parameter is required"})
		return
	}

	albums, err := ac.service.SearchAlbums(title)

	if err == nil {
		var responseAlbums []mapping.AlbumPreviewDTO
		for _, album := range albums {
			encodedId, _ := ac.encoder.Encode(album.Id)
			responseAlbums = append(responseAlbums, mapping.AlbumPreviewDTO{
				Id:            encodedId,
				Title:         album.Title,
				Type:          album.Type,
				CoverImageUrl: album.CoverImageUrl,
				ReleaseDate:   album.ReleaseDate,
			})
		}
		if responseAlbums == nil {
			responseAlbums = []mapping.AlbumPreviewDTO{}
		}

		c.JSON(http.StatusOK, responseAlbums)

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
