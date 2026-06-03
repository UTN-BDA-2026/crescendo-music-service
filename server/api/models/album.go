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
