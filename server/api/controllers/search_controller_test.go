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
	"time"

	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/require"
)

func (m *MockArtistService) SearchArtists(filter string) ([]models.Artist, error) {
	return []models.Artist{
		{
			Id:          5,
			Name:        "ABBA",
			Information: "Description of the artist",
			ImageUrl:    "a/dfv/gf.png",
		},
	}, nil
}

func (m *MockSongService) SearchSongs(filter string) ([]models.SongPreviewWithArtists, error) {
	return []models.SongPreviewWithArtists{
		{
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
	}, nil
}

func (m *MockAlbumService) SearchAlbums(filter string) ([]models.AlbumPreview, error) {
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

func TestCanProvideSearchResults(t *testing.T) {
	check := require.New(t)

	artistService := &MockArtistService{}
	songService := &MockSongService{}
	albumService := &MockAlbumService{}

	sq, err := sqids.New(sqids.Options{
		Alphabet:  os.Getenv("SQID_ALPHABET"),
		MinLength: 6,
	})
	check.NoError(err)
	sqEncoder := security.NewSquidEncoder(sq)

	testRouter := router.NewRouter(&app.Container{
		Search: controllers.NewSearchController(artistService, songService, albumService, sqEncoder),
	})

	searchWord := "searchword"
	searchType := "all"

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/search?q=%v&type=%v", searchWord, searchType),
		nil,
	)

	check.NoError(err)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	check.Equal(http.StatusOK, w.Code)
	check.NotEmpty(w.Body.String())
	var response mapping.SearchDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	check.NoError(err)

	check.Len(response.Artists, 1)
	check.Equal([]uint64{uint64(5)}, sq.Decode(response.Artists[0].Id))
	check.Equal("ABBA", response.Artists[0].Name)
	check.Equal("a/dfv/gf.png", response.Artists[0].ImageUrl)

	check.Len(response.Songs, 1)
	check.Equal([]uint64{uint64(4)}, sq.Decode(response.Songs[0].Id))
	check.Equal("Song Title", response.Songs[0].Title)
	check.NotEmpty(response.Songs[0].Artists)
	check.Equal(response.Songs[0].Artists[0].Name, "Artist 1")

	check.Len(response.Albums, 2)
	check.Equal([]uint64{uint64(8)}, sq.Decode(response.Albums[0].Id))
	check.Equal("JJ", response.Albums[0].Title)
	check.Equal("EP", response.Albums[0].Type)
	check.Equal("aaa/f.png", response.Albums[0].CoverImageUrl)
	check.Equal(
		time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
		response.Albums[0].ReleaseDate,
	)

	check.Equal([]uint64{uint64(9)}, sq.Decode(response.Albums[1].Id))
	check.Equal("SS", response.Albums[1].Title)
	check.Equal("Album", response.Albums[1].Type)
	check.Equal("bbb/j.png", response.Albums[1].CoverImageUrl)
	check.Equal(
		time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
		response.Albums[1].ReleaseDate,
	)
}
