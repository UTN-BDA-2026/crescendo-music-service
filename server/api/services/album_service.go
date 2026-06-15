package services

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type AlbumService interface {
	GetAlbumDetails(id int) (models.AlbumDetailed, error)
}

type albumService struct {
	repository      repositories.AlbumRepository
	genreRepository repositories.GenreRepository
	cache           *database.Cache
}

func NewAlbumService(r repositories.AlbumRepository,
	gr repositories.GenreRepository,
	c *database.Cache,
) AlbumService {
	return albumService{
		repository:      r,
		genreRepository: gr,
		cache:           c,
	}
}

func (s albumService) GetAlbumDetails(id int) (models.AlbumDetailed, error) {

	if id <= 0 {
		return models.AlbumDetailed{}, errors.New("invalid id")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("albums:details:%v", id)

		cachedValue, found, err := s.cache.Get(key)
		if err == nil && found {
			var album models.AlbumDetailed
			if err := json.Unmarshal([]byte(cachedValue), &album); err == nil {
				return album, nil
			}
		}
	}

	album, err := s.repository.GetById(id)
	if err != nil {
		log.Printf("fetching album failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	genre, err := s.genreRepository.GetById(album.GenreId)
	if err != nil {
		log.Printf("fetching genre failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	songs, err := s.repository.GetSongsPreviewFromAlbumId(id)
	if err != nil {
		log.Printf("fetching songs failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	result := models.AlbumDetailed{
		Id:            album.Id,
		Title:         album.Title,
		Type:          album.Type,
		CoverImageUrl: album.CoverImageUrl,
		ReleaseDate:   album.ReleaseDate,
		Genre:         genre,
		Songs:         songs,
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("albums:details:%v", id)
		data, err := json.Marshal(result)
		if err == nil {
			_ = s.cache.Set(key, string(data), 10*time.Minute)
		}
	}

	return result, nil
}
