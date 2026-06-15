package mapping

import (
	"crescendo-api/models"
	"crescendo-api/security"
	"time"
)

type SongCreateDTO struct {
	Title       string    `json:"title"`
	FileId      string    `json:"file_id"`
	GenreId     int       `json:"genre_id"`
	Duration    int       `json:"duration"`
	Bpm         int       `json:"bpm"`
	ReleaseDate time.Time `json:"release_date"`
}

type SongPreviewDTO struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

func SongPreviewListToDTO(encoder security.Encoder, list []models.SongPreview) ([]SongPreviewDTO, error) {
	var songsDTO []SongPreviewDTO
	for _, song := range list {
		hashedId, err := encoder.Encode(song.Id)
		if err != nil {
			return []SongPreviewDTO{}, nil
		}
		songsDTO = append(songsDTO, SongPreviewDTO{
			Id:       hashedId,
			Title:    song.Title,
			Duration: song.Duration,
		})
	}
	return songsDTO, nil
}

type SongSearchResultDTO struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	ArtistNames string `json:"artist_names"`
	AlbumTitles string `json:"album_titles"`
}

func SongSearchResultListToDTO(encoder security.Encoder, list []models.SongSearchResult) ([]SongSearchResultDTO, error) {
	var songsDTO []SongSearchResultDTO
	for _, song := range list {
		hashedId, err := encoder.Encode(song.Id)
		if err != nil {
			return []SongSearchResultDTO{}, nil
		}
		songsDTO = append(songsDTO, SongSearchResultDTO{
			Id:          hashedId,
			Title:       song.Title,
			Duration:    song.Duration,
			ArtistNames: song.ArtistNames,
			AlbumTitles: song.AlbumTitles,
		})
	}
	return songsDTO, nil
}

type ArtistLabelDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlaybackDataDTO struct {
	Id        string           `json:"id"`
	Title     string           `json:"title"`
	Duration  int              `json:"duration"`
	StreamURL string           `json:"stream_url"`
	Artists   []ArtistLabelDTO `json:"artists"`
}

func PlaybackDataToDTO(encoder security.Encoder, data models.PlaybackData) (PlaybackDataDTO, error) {

	hashedId, err := encoder.Encode(data.Id)
	if err != nil {
		return PlaybackDataDTO{}, err
	}

	var artistsDTO []ArtistLabelDTO

	for _, artist := range data.Artists {
		hashedArtistID, err := encoder.Encode(artist.Id)
		if err != nil {
			return PlaybackDataDTO{}, err
		}

		artistsDTO = append(artistsDTO, ArtistLabelDTO{
			Id:   hashedArtistID,
			Name: artist.Name,
		})
	}

	return PlaybackDataDTO{
		Id:        hashedId,
		Title:     data.Title,
		Duration:  data.Duration,
		StreamURL: data.StreamURL,
		Artists:   artistsDTO,
	}, nil
}

type ListedSongDTO struct {
	Id            string           `json:"id"`
	Title         string           `json:"title"`
	Duration      int              `json:"duration"`
	Artists       []ArtistLabelDTO `json:"artists"`
	TrackPosition int              `json:"track_position"`
}
