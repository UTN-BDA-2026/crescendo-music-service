package services_test

import (
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/services"
	"testing"
	"time"

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

func TestCanCreateSong_EmptyTitle(t *testing.T) {
	check := require.New(t)

	repo := mockSongRepository{
		createFunc: func(song models.Song) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	service := services.NewSongService(repo)

	song, err := service.Create(mapping.SongCreateDTO{
		Title:       "",
		FileId:      "507f1f77bcf86cd799439011",
		GenreId:     1,
		Duration:    200,
		Bpm:         120,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.Song{}, song)
}

func TestCanCreateSong_InvalidFileId(t *testing.T) {
	check := require.New(t)

	repo := mockSongRepository{
		createFunc: func(song models.Song) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	service := services.NewSongService(repo)

	song, err := service.Create(mapping.SongCreateDTO{
		Title:       "Title",
		FileId:      "idFalsa",
		GenreId:     1,
		Duration:    200,
		Bpm:         120,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.Song{}, song)
}

func TestCanCreateSong_InvalidGenreId(t *testing.T) {
	check := require.New(t)

	repo := mockSongRepository{
		createFunc: func(song models.Song) (int, error) {
			t.Fatal("repository should not be called")
			return 0, nil
		},
	}

	service := services.NewSongService(repo)

	song, err := service.Create(mapping.SongCreateDTO{
		Title:       "Title",
		FileId:      "507f1f77bcf86cd799439011",
		GenreId:     0,
		Duration:    200,
		Bpm:         120,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	})

	check.Error(err)
	check.Equal(models.Song{}, song)
}
