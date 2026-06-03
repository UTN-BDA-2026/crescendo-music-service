package models_test

import (
	"crescendo-api/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlaylist(t *testing.T) {
	playlist := models.Playlist{
		Id:           6,
		Title:        "TestTitle",
		Description:  "Texto de descripción",
		CreationDate: time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC),
		UserId:       3,
	}

	check := assert.New(t)

	check.Equal(6, playlist.Id)
	check.Equal("TestTitle", playlist.Title)
	check.Equal("Texto de descripción", playlist.Description)
	check.Equal(time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC), playlist.CreationDate)
	check.Equal(3, playlist.UserId)
}
