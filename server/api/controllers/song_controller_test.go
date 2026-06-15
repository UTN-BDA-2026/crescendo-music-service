package controllers_test

import (
	"crescendo-api/config/app"
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

	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/require"
)

type MockSongService struct{}

func (m *MockUserService) GetSongPlaybackInfo(song_id int) (models.PlaybackData, error) {
	return models.PlaybackData{
		SongPreviewWithArtists: models.SongPreviewWithArtists{
			Id:       4,
			Title:    "Song Title",
			Duration: 125,
			Artists: []models.ArtistLabel{
				{
					Id:   6,
					Name: "Artist 1",
				},
			},
		},
		StreamURL: "https://example.com/stream",
	}, nil
}

func TestCanProvidSongPlaybackInfo(t *testing.T) {
	check := require.New(t)
	service := &MockUserService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	sqEncoder := security.NewSquidEncoder(sq)
	testRouter := router.NewRouter(&app.Container{
		Song: controllers.NewSongController(service, sqEncoder),
	})
	hashedId, err := sq.Encode([]uint64{uint64(4)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/songs/%s/playback", hashedId), nil)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response mapping.PlaybackDataDTO

	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Equal([]uint64{uint64(4)}, sq.Decode(response.Id))
	check.Equal("Song Title", response.Title)
	check.NotEmpty(response.Artists)
	check.Equal(response.Artists[0].Name, "Artist 1")
}

func (m *MockUserService) SearchSongs(title string) ([]models.SongSearchResult, error) {
	if title == "Blinding Lights" {
		return []models.SongSearchResult{
			{
				Id:          1,
				Title:       "Blinding Lights",
				ArtistNames: "The Weeknd",
				AlbumTitles: "After Hours",
			},
		}, nil
	}
	return []models.SongSearchResult{}, nil
}

func TestSearchSongsByTitle(t *testing.T) {
	check := require.New(t)
	service := &MockUserService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	sqEncoder := security.NewSquidEncoder(sq)
	testRouter := router.NewRouter(&app.Container{
		Song: controllers.NewSongController(service, sqEncoder),
	})

	req, err := http.NewRequest(http.MethodGet, "/songs/search?title=Blinding Lights", nil)
	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)

	var response []mapping.SongSearchResultDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Len(response, 1)
	check.Equal([]uint64{uint64(1)}, sq.Decode(response[0].Id))
	check.Equal("Blinding Lights", response[0].Title)
	check.Equal("The Weeknd", response[0].ArtistNames)
	check.Equal("After Hours", response[0].AlbumTitles)
}
