package mapping

import (
	"crescendo-api/models"
	"crescendo-api/security"
)

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
