package mapping

import (
	"crescendo-api/models"
	"crescendo-api/security"
	"time"
)

type AlbumPreviewDTO struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	CoverImageUrl string    `json:"cover_image_url"`
	ReleaseDate   time.Time `json:"release_date"`
}

type AlbumDetailedDTO struct {
	Id            string          `json:"id"`
	Title         string          `json:"title"`
	Type          string          `json:"type"`
	Genre         GenreDTO        `json:"genre"`
	CoverImageUrl string          `json:"cover_image_url"`
	ReleaseDate   time.Time       `json:"release_date"`
	Songs         []ListedSongDTO `json:"songs"`
}

func AlbumPreviewListToDTO(encoder security.Encoder, list []models.AlbumPreview) ([]AlbumPreviewDTO, error) {
	var albumsDTO []AlbumPreviewDTO
	for _, album := range list {
		hashedId, err := encoder.Encode(album.Id)
		if err != nil {
			return []AlbumPreviewDTO{}, nil
		}
		albumsDTO = append(albumsDTO, AlbumPreviewDTO{
			Id:            hashedId,
			Title:         album.Title,
			Type:          album.Type,
			CoverImageUrl: album.CoverImageUrl,
			ReleaseDate:   album.ReleaseDate,
		})
	}
	return albumsDTO, nil
}

func AlbumDetailedToDTO(encoder security.Encoder, data models.AlbumDetailed) (AlbumDetailedDTO, error) {
	hashedId, err := encoder.Encode(data.Id)
	if err != nil {
		return AlbumDetailedDTO{}, err
	}

	hashedGenreId, err := encoder.Encode(data.Genre.Id)

	genreDTO := GenreDTO{
		Id:   hashedGenreId,
		Name: data.Genre.Name,
	}

	var songsDTO []ListedSongDTO

	for _, song := range data.Songs {
		hashedSongId, err := encoder.Encode(song.Id)
		if err != nil {
			return AlbumDetailedDTO{}, nil
		}
		var artistsDTO []ArtistLabelDTO

		for _, artist := range song.Artists {
			hashedArtistID, err := encoder.Encode(artist.Id)
			if err != nil {
				return AlbumDetailedDTO{}, err
			}

			artistsDTO = append(artistsDTO, ArtistLabelDTO{
				Id:   hashedArtistID,
				Name: artist.Name,
			})
		}
		songsDTO = append(songsDTO, ListedSongDTO{
			Id:            hashedSongId,
			Title:         song.Title,
			Duration:      song.Duration,
			TrackPosition: song.TrackPosition,
			Artists:       artistsDTO,
		})
	}

	return AlbumDetailedDTO{
		Id:            hashedId,
		Title:         data.Title,
		Type:          data.Type,
		CoverImageUrl: data.CoverImageUrl,
		ReleaseDate:   data.ReleaseDate,
		Genre:         genreDTO,
		Songs:         songsDTO,
	}, nil
}
