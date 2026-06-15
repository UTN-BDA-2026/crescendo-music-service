package services_test

import (
	"crescendo-api/models"
	"crescendo-api/services"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetAlbumDetails(t *testing.T) {
	check := require.New(t)

	expectedSongs := []models.ListedSong{
		{
			SongPreviewWithArtists: models.SongPreviewWithArtists{
				Id:       7,
				Title:    "Song 1",
				Duration: 125,
				Artists: []models.ArtistLabel{
					{Id: 1, Name: "Artist 1"},
				},
			},
			TrackPosition: 1,
		},
		{
			SongPreviewWithArtists: models.SongPreviewWithArtists{
				Id:       7,
				Title:    "Song 2",
				Duration: 653,
				Artists: []models.ArtistLabel{
					{Id: 2, Name: "Artist 2"},
				},
			},
			TrackPosition: 1,
		},
		{
			SongPreviewWithArtists: models.SongPreviewWithArtists{
				Id:       2,
				Title:    "Song 3",
				Duration: 222,
				Artists: []models.ArtistLabel{
					{Id: 1, Name: "Artist 1"},
				},
			},
			TrackPosition: 1,
		},
	}
	expectedAlbum := models.AlbumDetailed{
		Id:    3,
		Title: "Album Title",
		Genre: models.Genre{
			Id:   1,
			Name: "Rock",
		},
		CoverImageUrl: "http://host:8080/albumcover.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
		Songs:         expectedSongs,
	}
	albumId := 2
	albumRepository := mockAlbumRepository{
		getSongsPreviewFromAlbumIdFunc: func(albumId int) ([]models.ListedSong, error) {
			return expectedSongs, nil
		},
		getByIdFunc: func(id int) (models.Album, error) {
			return models.Album{
				Id:            3,
				Title:         "Album Title",
				CoverImageUrl: "http://host:8080/albumcover.png",
				ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
				GenreId:       1,
			}, nil
		},
	}
	genreRepository := mockGenreRepository{
		getByIdFunc: func(id int) (models.Genre, error) {
			return models.Genre{
				Id:   1,
				Name: "Rock",
			}, nil
		},
	}
	service := services.NewAlbumService(albumRepository, genreRepository, nil)

	album, err := service.GetAlbumDetails(albumId)

	check.NoError(err)
	check.Equal(expectedAlbum, album)
}
