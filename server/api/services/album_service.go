package services

import (
	"crescendo-api/models"
	"crescendo-api/repositories"
	"errors"
	"log"
)

type AlbumService interface {
	GetAlbumDetails(id int) (models.AlbumDetailed, error)
}

type albumService struct {
	repository      repositories.AlbumRepository
	genreRepository repositories.GenreRepository
}

func NewAlbumService(repository repositories.AlbumRepository, genreRepo repositories.GenreRepository) AlbumService {
	service := albumService{
		repository:      repository,
		genreRepository: genreRepo,
	}

	return service
}

func (s albumService) GetAlbumDetails(id int) (models.AlbumDetailed, error) {

	if id <= 0 {
		return models.AlbumDetailed{}, errors.New("invalid id")
	}

	album, err := s.repository.GetById(id)

	if err != nil {
		log.Printf("fetching album for album details failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	genre, err := s.genreRepository.GetById(album.GenreId)

	if err != nil {
		log.Printf("fetching genre for album details failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	songs, err := s.repository.GetSongsPreviewFromAlbumId(id)

	if err != nil {
		log.Printf("fetching songs for album details failed: %v", err)
		return models.AlbumDetailed{}, errors.New("something went wrong")
	}

	return models.AlbumDetailed{
		Id:            album.Id,
		Title:         album.Title,
		Type:          album.Type,
		CoverImageUrl: album.CoverImageUrl,
		ReleaseDate:   album.ReleaseDate,
		Genre:         genre,
		Songs:         songs,
	}, nil
}
