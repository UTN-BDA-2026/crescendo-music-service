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

type PlaybackData struct {
	Id        int
	Title     string
	Duration  int
	StreamURL string
	Artists   []ArtistLabel
}
