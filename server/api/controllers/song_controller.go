package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/security"
	"crescendo-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SongController struct {
	service services.SongService
	encoder security.Encoder
}

func NewSongController(s services.SongService, e security.Encoder) *SongController {
	return &SongController{
		service: s,
		encoder: e,
	}
}

func (sc *SongController) GetSongPlaybackInfo(c *gin.Context) {
	hashID := c.Param("id")
	id, err := sc.encoder.Decode(hashID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid song id",
		})
		return
	}

	playbackData, err := sc.service.GetSongPlaybackInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get playback info",
		})
		return
	}
	playbackDataDTO, err := mapping.PlaybackDataToDTO(sc.encoder, playbackData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get playback info",
		})
		return
	}
	c.JSON(http.StatusOK, playbackDataDTO)
}
