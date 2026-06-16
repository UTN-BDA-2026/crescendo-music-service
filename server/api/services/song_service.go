package services

import (
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
	SearchSongs(name string) ([]models.SongPreviewWithArtists, error)
}

type songService struct {
	repository repositories.SongRepository
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

func (s songService) SearchSongs(name string) ([]models.SongPreviewWithArtists, error) {
	if name == "" {
		return nil, errors.New("invalid search string")
	}

	songs, err := s.repository.FindByNameLike(name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.SongPreviewWithArtists{}, nil
		}

		log.Printf("fetching songs for search %q failed: %v", name, err)
		return []models.SongPreviewWithArtists{}, errors.New("something went wrong")
	}
	return songs, nil
}
