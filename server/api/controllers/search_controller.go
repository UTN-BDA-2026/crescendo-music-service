package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/security"
	"crescendo-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	artistService services.ArtistService
	songService   services.SongService
	albumService  services.AlbumService
	encoder       security.Encoder
}

func NewSearchController(
	ar services.ArtistService,
	s services.SongService,
	al services.AlbumService,
	e security.Encoder,
) *SearchController {
	return &SearchController{
		artistService: ar,
		songService:   s,
		albumService:  al,
		encoder:       e,
	}
}

func (sc *SearchController) SearchByName(c *gin.Context) {
	q := c.Query("q")
	searchType := c.Query("type")

	validSearchTypes := map[string]struct{}{
		"all":    {},
		"artist": {},
		"song":   {},
		"album":  {},
	}

	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid search parameter",
		})
		return
	}

	if _, ok := validSearchTypes[searchType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid search type",
		})
		return
	}
	var searchResults mapping.SearchDTO

	if searchType == "all" || searchType == "artist" {
		artists, err := sc.artistService.SearchArtists(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		searchResults.Artists, err = mapping.ArtistListToDTO(sc.encoder, artists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
	}
	if searchType == "all" || searchType == "song" {
		songs, err := sc.songService.SearchSongs(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		searchResults.Songs, err = mapping.SongPreviewListWithArtistsToDTO(sc.encoder, songs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
	}
	if searchType == "all" || searchType == "album" {
		albums, err := sc.albumService.SearchAlbums(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		searchResults.Albums, err = mapping.AlbumPreviewListToDTO(sc.encoder, albums)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
	}
	c.JSON(http.StatusOK, searchResults)
}
