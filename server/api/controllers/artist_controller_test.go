package controllers_test

import (
	"crescendo-api/config/app"
	"crescendo-api/config/env"
	"crescendo-api/controllers"
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/router"
	"crescendo-api/security"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/require"
)

func init() {
	env.Load()
}

type MockArtistService struct{}

func (m *MockArtistService) GetArtist(id int) (models.Artist, error) {
	return models.Artist{
		Id:          5,
		Name:        "ABBA",
		Information: "Description of the artist",
		ImageUrl:    "a/dfv/gf.png",
	}, nil
}

func (m *MockArtistService) GetAllArtist() ([]models.Artist, error) {
	return []models.Artist{
		{
			Id:          5,
			Name:        "ABBA",
			Information: "Description of the artist",
			ImageUrl:    "a/dfv/gf.png",
		},
	}, nil
}

func (m *MockArtistService) GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error) {
	return []models.AlbumPreview{
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
	}, nil
}

func (m *MockArtistService) GetArtistSongPreviews(id int) ([]models.SongPreview, error) {
	return []models.SongPreview{
		{
			Id:       12,
			Title:    "Song 1",
			Duration: 124,
		},
		{
			Id:       43,
			Title:    "Song 2",
			Duration: 332,
		},
	}, nil
}

func TestGetArtist(t *testing.T) {
	check := require.New(t)

	service := &MockArtistService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	check.NoError(err)
	sqEncoder := security.NewSquidEncoder(sq)

	testRouter := router.NewRouter(&app.Container{
		Artist: controllers.NewArtistController(service, sqEncoder),
	})
	token, err := security.GenerateLoginToken(1, "testuser")
	check.NoError(err)
	hashedId, err := sq.Encode([]uint64{uint64(5)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/artists/%s", hashedId), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response mapping.ArtistDTO

	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Equal([]uint64{uint64(5)}, sq.Decode(response.Id))
	check.Equal("ABBA", response.Name)
	check.Equal("Description of the artist", response.Information)
	check.Equal("a/dfv/gf.png", response.ImageUrl)
}

func TestGetArtistAlbumPreviews(t *testing.T) {
	check := require.New(t)

	service := &MockArtistService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	check.NoError(err)
	sqEncoder := security.NewSquidEncoder(sq)

	testRouter := router.NewRouter(&app.Container{
		Artist: controllers.NewArtistController(service, sqEncoder),
	})
	token, err := security.GenerateLoginToken(1, "testuser")
	check.NoError(err)
	hashedId, err := sq.Encode([]uint64{uint64(5)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/artists/%s/albums", hashedId), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response []mapping.AlbumPreviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Len(response, 2)
	check.Equal([]uint64{uint64(8)}, sq.Decode(response[0].Id))
	check.Equal("JJ", response[0].Title)
	check.Equal("EP", response[0].Type)
	check.Equal("aaa/f.png", response[0].CoverImageUrl)
	check.Equal(
		time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
		response[0].ReleaseDate,
	)

	check.Equal([]uint64{uint64(9)}, sq.Decode(response[1].Id))
	check.Equal("SS", response[1].Title)
	check.Equal("Album", response[1].Type)
	check.Equal("bbb/j.png", response[1].CoverImageUrl)
	check.Equal(
		time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
		response[1].ReleaseDate,
	)

}

func TestGetArtistSongPreview(t *testing.T) {
	check := require.New(t)

	service := &MockArtistService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	check.NoError(err)
	sqEncoder := security.NewSquidEncoder(sq)

	testRouter := router.NewRouter(&app.Container{
		Artist: controllers.NewArtistController(service, sqEncoder),
	})
	token, err := security.GenerateLoginToken(1, "testuser")
	check.NoError(err)
	hashedId, err := sq.Encode([]uint64{uint64(5)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/artists/%s/songs", hashedId), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response []mapping.SongPreviewDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)

	check.NoError(err)
	check.Len(response, 2)

	check.Equal([]uint64{uint64(12)}, sq.Decode(response[0].Id))
	check.Equal("Song 1", response[0].Title)
	check.Equal(124, response[0].Duration)

	check.Equal([]uint64{uint64(43)}, sq.Decode(response[1].Id))
	check.Equal("Song 2", response[1].Title)
	check.Equal(332, response[1].Duration)
}

func TestGetAllArtist(t *testing.T) {
	check := require.New(t)

	service := &MockArtistService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	check.NoError(err)
	sqEncoder := security.NewSquidEncoder(sq)

	testRouter := router.NewRouter(&app.Container{
		Artist: controllers.NewArtistController(service, sqEncoder),
	})
	token, err := security.GenerateLoginToken(1, "testuser")
	check.NoError(err)
	req, err := http.NewRequest(http.MethodGet, "/artists", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response []mapping.ArtistDTO

	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Len(response, 1)
	check.Equal([]uint64{uint64(5)}, sq.Decode(response[0].Id))
	check.Equal("ABBA", response[0].Name)
	check.Equal("Description of the artist", response[0].Information)
	check.Equal("a/dfv/gf.png", response[0].ImageUrl)
}
