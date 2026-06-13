package services

import (
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"crescendo-api/security"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

type SongService interface {
	GetSongPlaybackInfo(id int) (models.PlaybackData, error)
	Create(mapping.SongCreateDTO) (models.Song, error)
}

type songService struct {
	repository repositories.SongRepository
}

func isValidFileId(fileId string) bool {
	return len(fileId) == 24
}

func NewSongService(repository repositories.SongRepository) SongService {
	service := songService{
		repository: repository,
	}

	return service
}

func (s songService) GetSongPlaybackInfo(id int) (models.PlaybackData, error) {

	if id <= 0 {
		return models.PlaybackData{}, errors.New("invalid id")
	}
	song, err := s.repository.GetById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PlaybackData{}, errors.New("invalid id")
		}

		log.Printf("fetching song for playback failed: %v", err)
		return models.PlaybackData{}, errors.New("something went wrong")
	}

	artists, err := s.repository.GetArtistsForPlaybackBySongId(id)

	if err != nil {
		log.Printf("fetching artist for playback failed: %v", err)
		return models.PlaybackData{}, errors.New("something went wrong")
	}

	token, err := security.GenerateStreamToken(song.Id, song.FileId)

	if err != nil {
		log.Printf("generating token for playback failed: %v", err)
		return models.PlaybackData{}, errors.New("something went wrong")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Printf("generating playback url failed: base url not configured")
		return models.PlaybackData{}, errors.New("something went wrong")
	}
	streamURL := fmt.Sprintf(
		"%s/stream?token=%s",
		baseURL,
		token,
	)

	return models.PlaybackData{
		SongPreviewWithArtists: models.SongPreviewWithArtists{
			Id:       song.Id,
			Title:    song.Title,
			Duration: song.Duration,
			Artists:  artists,
		},

		StreamURL: streamURL,
	}, nil
}

func (s songService) Create(requestDTO mapping.SongCreateDTO) (models.Song, error) {
	if requestDTO.Title == "" {
		return models.Song{}, errors.New("invalid song title")
	}

	if !isValidFileId(requestDTO.FileId) {
		return models.Song{}, errors.New("invalid file id")
	}

	if requestDTO.GenreId <= 0 {
		return models.Song{}, errors.New("invalid genre id")
	}

	return models.Song{}, nil
}
