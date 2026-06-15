package services_test

import (
	"crescendo-api/models"
	"crescendo-api/services"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetArtist(t *testing.T) {
	check := require.New(t)

	artist := models.Artist{
		Id:          5,
		Name:        "Artist 1",
		Information: "Artist description",
		ImageUrl:    "fvfu/erui.png",
	}
	repo := mockArtistRepository{
		getByIdFunc: func(id int) (models.Artist, error) {
			return artist, nil
		},
	}

	service := services.NewArtistService(repo)

	fetchedArtist, err := service.GetArtist(artist.Id)

	check.NoError(err)
	check.Equal(artist, fetchedArtist)
}

func TestGetArtistAlbums(t *testing.T) {
	check := require.New(t)

	artistId := 5
	referenceAlbums := []models.AlbumPreview{
		{
			Id:            8,
			Title:         "JJ",
			Type:          "EP",
			CoverImageUrl: "aaa/f.png",
			ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:            9,
			Title:         "SS",
			Type:          "Album",
			CoverImageUrl: "bbb/j.png",
			ReleaseDate:   time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
		},
	}
	repo := mockArtistRepository{
		getAlbumPreviewsByArtistIdFunc: func(id int) ([]models.AlbumPreview, error) {
			return referenceAlbums, nil
		},
	}

	service := services.NewArtistService(repo)

	fetchedAlbums, err := service.GetArtistAlbumPreviews(artistId)

	check.NoError(err)
	check.Equal(referenceAlbums, fetchedAlbums)

}

func TestGetArtistSongPreviews(t *testing.T) {
	check := require.New(t)
	artistId := 5

	referenceSongs := []models.SongPreview{
		{
			Id:       4,
			Title:    "Song 1",
			Duration: 123,
		},
		{
			Id:       8,
			Title:    "Song 2",
			Duration: 345,
		},
	}

	repo := mockArtistRepository{
		getArtistSongPreviewsFunc: func(id int) ([]models.SongPreview, error) {
			return referenceSongs, nil
		},
	}

	service := services.NewArtistService(repo)

	fetchedSongs, err := service.GetArtistSongPreviews(artistId)

	check.NoError(err)
	check.Len(fetchedSongs, len(referenceSongs))
	check.Equal(referenceSongs, fetchedSongs)

}

func TestGetAllArtist(t *testing.T) {
	check := require.New(t)

	artists := []models.Artist{
		{
			Id:          5,
			Name:        "Artist 1",
			Information: "Artist description",
			ImageUrl:    "fvfu/erui.png",
		},
		{
			Id:          7,
			Name:        "Artist 2",
			Information: "Artist description 4",
			ImageUrl:    "fvfu/erui.png",
		},
	}
	repo := mockArtistRepository{
		getAllFunc: func() ([]models.Artist, error) {
			return artists, nil
		},
	}

	service := services.NewArtistService(repo)

	fetchedArtists, err := service.GetAllArtist()

	check.NoError(err)
	check.Equal(artists, fetchedArtists)
}
