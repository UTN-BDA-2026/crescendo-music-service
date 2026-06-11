package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/security"
	"crescendo-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArtistController struct {
	service services.ArtistService
	encoder security.Encoder
}

func NewArtistController(s services.ArtistService, e security.Encoder) *ArtistController {
	return &ArtistController{
		service: s,
		encoder: e,
	}
}

func (ac *ArtistController) GetArtist(c *gin.Context) {
	hashID := c.Param("id")
	id, err := ac.encoder.Decode(hashID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid artist id",
		})
		return
	}

	artist, err := ac.service.GetArtist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist details",
		})
		return
	}
	hashedId, err := ac.encoder.Encode(artist.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist details",
		})
		return
	}
	artistDTO := mapping.ArtistDTO{
		Id:          hashedId,
		Name:        artist.Name,
		Information: artist.Information,
		ImageUrl:    artist.ImageUrl,
	}
	c.JSON(http.StatusOK, artistDTO)
}

func (ac *ArtistController) GetArtistAlbumPreviews(c *gin.Context) {
	hashID := c.Param("id")
	id, err := ac.encoder.Decode(hashID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid artist id",
		})
		return
	}
	albumPreviews, err := ac.service.GetArtistAlbumPreviews(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist albums",
		})
		return
	}
	albumPreviewsDTOs, err := mapping.AlbumPreviewListToDTO(ac.encoder, albumPreviews)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist albums",
		})
		return
	}
	c.JSON(http.StatusOK, albumPreviewsDTOs)
}

func (ac *ArtistController) GetArtistSongPreviews(c *gin.Context) {
	hashID := c.Param("id")
	id, err := ac.encoder.Decode(hashID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid artist id",
		})
		return
	}
	songPreviews, err := ac.service.GetArtistSongPreviews(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist songs",
		})
		return
	}
	songPreviewsDTOs, err := mapping.SongPreviewListToDTO(ac.encoder, songPreviews)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get artist songs",
		})
		return
	}
	c.JSON(http.StatusOK, songPreviewsDTOs)
}
