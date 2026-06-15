package mapping

import (
	"crescendo-api/models"
	"crescendo-api/security"
)

type ArtistDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Information string `json:"information"`
	ImageUrl    string `json:"image_url"`
}

func ArtistListToDTO(encoder security.Encoder, list []models.Artist) ([]ArtistDTO, error) {
	var artistsDTO []ArtistDTO
	for _, artist := range list {
		hashedId, err := encoder.Encode(artist.Id)
		if err != nil {
			return []ArtistDTO{}, nil
		}
		artistsDTO = append(artistsDTO, ArtistDTO{
			Id:          hashedId,
			Name:        artist.Name,
			Information: artist.Information,
			ImageUrl:    artist.ImageUrl,
		})
	}
	return artistsDTO, nil
}
