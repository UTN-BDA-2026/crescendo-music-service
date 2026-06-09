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

type SongPreview struct {
	Id       int
	Title    string
	Duration int
	Artists  []ArtistLabel
}

type PlaybackData struct {
	SongPreview
	StreamURL string
}
