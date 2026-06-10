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
	Id            int
	Title         string
	Type          string
	Genre         Genre
	CoverImageUrl string
	ReleaseDate   time.Time
	Songs         []ListedSong
}

type AlbumPreview struct {
	Id            int
	Title         string
	Type          string
	CoverImageUrl string
	ReleaseDate   time.Time
}
