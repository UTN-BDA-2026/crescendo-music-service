package models

import "time"

type Album struct {
	Id            int
	Title         string
	Type          string
	GenreId       int
	CoverImageUrl string
	ReleaseDate   time.Time
}

type AlbumDetailed struct {
	Id            int          `json:"id"`
	Title         string       `json:"title"`
	Type          string       `json:"type"`
	Genre         Genre        `json:"genre"`
	CoverImageUrl string       `json:"cover_image_url"`
	ReleaseDate   time.Time    `json:"release_date"`
	Songs         []ListedSong `json:"songs"`
}

type AlbumPreview struct {
	Id            int
	Title         string
	Type          string
	CoverImageUrl string
	ReleaseDate   time.Time
}
