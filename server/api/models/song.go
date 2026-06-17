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
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type SongPreviewWithArtists struct {
	Id       int           `json:"id"`
	Title    string        `json:"title"`
	Duration int           `json:"duration"`
	Artists  []ArtistLabel `json:"artists"`
}

type PlaybackData struct {
	SongPreviewWithArtists
	StreamURL string
}

type ListedSong struct {
	SongPreviewWithArtists
	TrackPosition int `json:"track_position"`
}
