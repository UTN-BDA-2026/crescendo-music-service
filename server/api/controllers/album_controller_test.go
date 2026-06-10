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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/require"
)

type MockAlbumService struct{}

func (m *MockAlbumService) GetAlbumDetails(id int) (models.AlbumDetailed, error) {
	return models.AlbumDetailed{
		Id:    5,
		Title: "Album Title",
		Type:  "EP",
		Genre: models.Genre{
			Id:   3,
			Name: "Pop",
		},
		CoverImageUrl: "aaa/f.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
		Songs: []models.ListedSong{
			{
				SongPreview: models.SongPreview{
					Id:       1,
					Title:    "Song Title",
					Duration: 344,
					Artists: []models.ArtistLabel{
						{
							Id:   8,
							Name: "Artist 1",
						},
					},
				},
				TrackPosition: 1},
		},
	}, nil
}

func TestCanProvideAlbumDetails(t *testing.T) {
	check := require.New(t)
	service := &MockAlbumService{}
	sq, err := sqids.New(sqids.Options{
		Alphabet: os.Getenv("SQID_ALPHABET"),
	})
	sqEncoder := security.NewSquidEncoder(sq)
	controller := controllers.NewAlbumController(service, sqEncoder)
	router := gin.Default()
	router.GET("/albums/:id", controller.GetAlbumDetails)
	hashedId, err := sq.Encode([]uint64{uint64(5)})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/albums/%s", hashedId), nil)
	check.NoError(err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	var response mapping.AlbumDetailedDTO

	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Equal([]uint64{uint64(5)}, sq.Decode(response.Id))
	check.Equal("Album Title", response.Title)
	check.Len(response.Songs, 1)
	check.Equal("EP", response.Type)
	check.Equal("aaa/f.png", response.CoverImageUrl)

	check.Equal([]uint64{3}, sq.Decode(response.Genre.Id))

	check.Equal("Pop", response.Genre.Name)

	check.Len(response.Songs, 1)

	song := response.Songs[0]
	check.Equal([]uint64{1}, sq.Decode(song.Id))

	check.Equal("Song Title", song.Title)
	check.Equal(344, song.Duration)
	check.Equal(1, song.TrackPosition)

	check.Len(song.Artists, 1)

	artist := song.Artists[0]

	check.Equal([]uint64{8}, sq.Decode(artist.Id))

	check.Equal("Artist 1", artist.Name)
}
