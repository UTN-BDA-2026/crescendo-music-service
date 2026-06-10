package services

import (
	"crescendo-api/models"
	"crescendo-api/repositories"
	"database/sql"
	"errors"
	"log"
)

type ArtistService interface {
	GetArtist(id int) (models.Artist, error)
	GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error)
}

type artistService struct {
	repository repositories.ArtistRepository
}

func NewArtistService(repository repositories.ArtistRepository) ArtistService {
	service := artistService{
		repository: repository,
	}

	return service
}

func (s artistService) GetArtist(id int) (models.Artist, error) {

	if id <= 0 {
		return models.Artist{}, errors.New("invalid id")
	}

	artist, err := s.repository.GetById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Artist{}, errors.New("invalid id")
		}

		log.Printf("fetching artists failed: %v", err)
		return models.Artist{}, errors.New("something went wrong")
	}
	return artist, nil
}

func (s artistService) GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error) {
	if id <= 0 {
		return []models.AlbumPreview{}, errors.New("invalid id")
	}
	albums, err := s.repository.GetArtistAlbumPreviews(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.AlbumPreview{}, nil
		}

		log.Printf("fetching albums from artist %v failed: %v", id, err)
		return []models.AlbumPreview{}, errors.New("something went wrong")
	}
	return albums, nil
}
