package models_test

import (
	"crescendo-api/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateSong(t *testing.T) {

	song := models.Song{
		Id:          1,
		Title:       "Song title",
		FileId:      "0xq2454",
		GenreId:     4,
		Duration:    253,
		Bpm:         110,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	check := assert.New(t)

	check.Equal(1, song.Id)
	check.Equal("Song Title", song.Title)
	check.Equal("0xq2454", song.FileId)
	check.Equal(4, song.GenreId)
	check.Equal(253, song.Duration)
	check.Equal(110, song.Bpm)
	check.Equal(time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC), song.ReleaseDate)
}
