package models_test

import (
	"crescendo-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArtist(t *testing.T) {
	artist := models.Artist{
		Id:          5,
		Name:        "ABBA",
		Information: "Description of the artist",
		ImageUrl:    "a/dfv/gf.png",
	}
	check := assert.New(t)

	check.Equal(5, artist.Id)
	check.Equal("ABBA", artist.Name)
	check.Equal("Description of the artist", artist.Information)
	check.Equal("a/dfv/gf.png", artist.ImageUrl)
}
