package models

import "time"

type Playlist struct {
	Id           int
	Title        string
	Description  string
	CreationDate time.Time
	UserId       int
}
