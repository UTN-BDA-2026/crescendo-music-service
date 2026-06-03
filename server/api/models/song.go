package models

import "time"

type Song struct {
	Id          int
	Title       string
	FileId      string
	GenreId     int
	Duration    int
	Bpm         int
	ReleaseDate time.Time
}
