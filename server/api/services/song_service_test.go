package services_test

import (
	"crescendo-api/models"
	"crescendo-api/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaySong_Success(t *testing.T) {
	check := require.New(t)

	repo := mockSongRepository{
		getByIdFunc: func(id int) (models.Song, error) {
			return models.Song{
				Id:       5,
				Title:    "Song Title",
				FileId:   "gridfs:abc123",
				Duration: 125,
			}, nil
		},
		addArtistToSongFunc: func(artistId int, songId int) error {
			return nil
		},
		getArtistsForPlaybackBySongIdFunc: func(id int) ([]models.ArtistLabel, error) {
			return []models.ArtistLabel{
				{
					Id:   1,
					Name: "Artist 1",
				},
				{
					Id:   2,
					Name: "Artist 2",
				},
			}, nil
		},
	}

	service := services.NewSongService(repo)
	songId := 5

	playbackData, err := service.GetSongPlaybackInfo(songId)

	check.NoError(err)
	check.NotEqual(models.PlaybackData{}, playbackData)
	check.Equal(5, playbackData.Id)
	check.Equal("Song Title", playbackData.Title)
	check.Equal(125, playbackData.Duration)
	check.NotEmpty(playbackData.Artists)
	check.NotEmpty(playbackData.StreamURL)
	check.Equal(models.ArtistLabel{Id: 1, Name: "Artist 1"}, playbackData.Artists[0])
	check.Equal(models.ArtistLabel{Id: 2, Name: "Artist 2"}, playbackData.Artists[1])
}

func TestSearchSongs(t *testing.T) {
	check := require.New(t)

	referenceSongs := []models.SongSearchResult{
		{
			Id:          1,
			Title:       "Blinding Lights",
			ArtistNames: "The Weeknd",
			AlbumTitles: "After Hours",
		},
	}

	repo := mockSongRepository{
		searchByTitleFunc: func(title string) ([]models.SongSearchResult, error) {
			if title == "Blinding Lights" {
				return referenceSongs, nil
			}
			return []models.SongSearchResult{}, nil
		},
	}

	service := services.NewSongService(repo)

	fetchedSongs, err := service.SearchSongs("Blinding Lights")
	check.NoError(err)
	check.Equal(referenceSongs, fetchedSongs)

	_, err = service.SearchSongs("")
	check.Error(err)
}
