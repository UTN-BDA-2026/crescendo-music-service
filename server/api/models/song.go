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

type SongPreviewWithArtists struct {
	Id       int
	Title    string
	Duration int
	Artists  []ArtistLabel
}

type PlaybackData struct {
	SongPreviewWithArtists
	StreamURL string
}

type ListedSong struct {
	SongPreviewWithArtists
	TrackPosition int
}
