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
		Alphabet: os.Getenv("SQID_ALPHABET"),
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
