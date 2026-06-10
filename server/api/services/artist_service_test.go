package services_test

import (
	"crescendo-api/models"
	"crescendo-api/services"
	"testing"

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
