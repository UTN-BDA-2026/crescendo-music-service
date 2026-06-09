package controllers_test

import (
	"crescendo-api/controllers"
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/security"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/require"
)

type MockSongService struct{}

func (m *MockUserService) GetSongPlaybackInfo(song_id int) (models.PlaybackData, error) {
	return models.PlaybackData{
		SongPreview: models.SongPreview{
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
	controller := controllers.NewSongController(service, sqEncoder)
	router := gin.Default()
	router.GET("/songs/:id/playback", controller.GetSongPlaybackInfo)
	hashedId, err := sq.Encode([]uint64{uint64(4)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/songs/%s/playback", hashedId), nil)
	check.NoError(err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response mapping.PlaybackDataDTO

	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Equal([]uint64{uint64(4)}, sq.Decode(response.Id))
	check.Equal("Song Title", response.Title)
	check.NotEmpty(response.Artists)
	check.Equal(response.Artists[0].Name, "Artist 1")
}
