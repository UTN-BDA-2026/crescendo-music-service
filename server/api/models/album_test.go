package models_test

import (
	"crescendo-api/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateAlbum(t *testing.T) {
	album := models.Album{
		Id:            8,
		Title:         "JJ",
		Type:          "EP",
		GenreId:       4,
		CoverImageUrl: "aaa/f.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	check := require.New(t)
	check.Equal(8, album.Id)
	check.Equal("JJ", album.Title)
	check.Equal("EP", album.Type)
	check.Equal(4, album.GenreId)
	check.Equal("aaa/f.png", album.CoverImageUrl)
	check.Equal(time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC), album.ReleaseDate)
}
