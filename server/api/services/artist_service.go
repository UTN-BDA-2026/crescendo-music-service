package services

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type ArtistService interface {
	GetArtist(id int) (models.Artist, error)
	GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error)
	GetArtistSongPreviews(id int) ([]models.SongPreview, error)
	GetAllArtist() ([]models.Artist, error)
}

type artistService struct {
	repository repositories.ArtistRepository
	cache      *database.Cache
}

func NewArtistService(repository repositories.ArtistRepository, cache *database.Cache) ArtistService {
	service := artistService{
		repository: repository,
		cache:      cache,
	}

	return service
}

func (s artistService) GetArtist(id int) (models.Artist, error) {

	if id <= 0 {
		return models.Artist{}, errors.New("invalid id")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v", id)

		cachedValue, found, err := s.cache.Get(key)
		if err == nil && found {
			var artist models.Artist
			if err := json.Unmarshal([]byte(cachedValue), &artist); err == nil {
				return artist, nil
			}
		}
	}

	artist, err := s.repository.GetById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Artist{}, errors.New("invalid id")
		}

		log.Printf("fetching artists failed: %v", err)
		return models.Artist{}, errors.New("something went wrong")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v", id)
		data, err := json.Marshal(artist)
		if err == nil {
			_ = s.cache.Set(key, string(data), 30*time.Minute)
		}
	}

	return artist, nil
}

func (s artistService) GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error) {
	if id <= 0 {
		return []models.AlbumPreview{}, errors.New("invalid id")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v:albums", id)

		cachedValue, found, err := s.cache.Get(key)
		if err == nil && found {
			var albums []models.AlbumPreview
			if err := json.Unmarshal([]byte(cachedValue), &albums); err == nil {
				return albums, nil
			}
		}
	}

	albums, err := s.repository.GetArtistAlbumPreviews(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.AlbumPreview{}, nil
		}

		log.Printf("fetching albums from artist %v failed: %v", id, err)
		return []models.AlbumPreview{}, errors.New("something went wrong")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v:albums", id)
		data, err := json.Marshal(albums)
		if err == nil {
			_ = s.cache.Set(key, string(data), 30*time.Minute)
		}
	}
	return albums, nil
}

func (s artistService) GetArtistSongPreviews(id int) ([]models.SongPreview, error) {
	if id <= 0 {
		return []models.SongPreview{}, errors.New("invalid id")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v:songs", id)

		cachedValue, found, err := s.cache.Get(key)
		if err == nil && found {
			var songs []models.SongPreview
			if err := json.Unmarshal([]byte(cachedValue), &songs); err == nil {
				return songs, nil
			}
		}
	}

	songs, err := s.repository.GetArtistSongPreviews(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.SongPreview{}, nil
		}

		log.Printf("fetching songs from artist %v failed: %v", id, err)
		return []models.SongPreview{}, errors.New("something went wrong")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := fmt.Sprintf("artist:%v:songs", id)
		data, err := json.Marshal(songs)
		if err == nil {
			_ = s.cache.Set(key, string(data), 30*time.Minute)
		}
	}

	return songs, nil
}

func (s artistService) GetAllArtist() ([]models.Artist, error) {

	if s.cache != nil && s.cache.IsReady() {
		key := "artists"

		cachedValue, found, err := s.cache.Get(key)
		if err == nil && found {
			var artists []models.Artist
			if err := json.Unmarshal([]byte(cachedValue), &artists); err == nil {
				return artists, nil
			}
		}
	}
	artists, err := s.repository.GetAll()

	if err != nil {
		log.Printf("fetching artists list failed: %v", err)
		return []models.Artist{}, errors.New("something went wrong")
	}

	if s.cache != nil && s.cache.IsReady() {
		key := "artists"
		data, err := json.Marshal(artists)
		if err == nil {
			_ = s.cache.Set(key, string(data), 30*time.Minute)
		}
	}
	return artists, nil
}
