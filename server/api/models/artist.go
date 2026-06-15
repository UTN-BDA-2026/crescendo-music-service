package models

type Artist struct {
	Id          int
	Name        string
	Information string
	ImageUrl    string
}

type ArtistLabel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
