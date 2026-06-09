package mapping

import (
	"crescendo-api/models"
	"crescendo-api/security"
)

type ArtistLabelDTO struct {
	Id   string
	Name string
}

type PlaybackDataDTO struct {
	Id        string
	Title     string
	Duration  int
	StreamURL string
	Artists   []ArtistLabelDTO
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
