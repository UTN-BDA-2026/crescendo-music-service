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

	repo := mockSongRepository{
		findByNameLikeFunc: func(name string) ([]models.SongPreviewWithArtists, error) {
			return []models.SongPreviewWithArtists{
				{
					Id:       5,
					Title:    "Song Title",
					Duration: 125,
					Artists: []models.ArtistLabel{
						{
							Id:   1,
							Name: "Artist 1",
						},
						{
							Id:   2,
							Name: "Artist 2",
						},
					},
				},
			}, nil
		},
	}
	service := services.NewSongService(repo)

	songs, err := service.SearchSongs("Song")

	check.NoError(err)
	check.NotEmpty(songs)
	check.Len(songs, 1)
	check.Equal(5, songs[0].Id)
	check.Equal("Song Title", songs[0].Title)
	check.Equal(125, songs[0].Duration)
	check.NotEmpty(songs[0].Artists)
	check.Equal(models.ArtistLabel{Id: 1, Name: "Artist 1"}, songs[0].Artists[0])
	check.Equal(models.ArtistLabel{Id: 2, Name: "Artist 2"}, songs[0].Artists[1])
}
