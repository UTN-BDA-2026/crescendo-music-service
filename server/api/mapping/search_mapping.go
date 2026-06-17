package mapping

type SearchDTO struct {
	Artists []ArtistDTO                 `json:"artists"`
	Songs   []SongPreviewWithArtistsDTO `json:"songs"`
	Albums  []AlbumPreviewDTO           `json:"albums"`
}
