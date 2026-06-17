package models

type Artist struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Information string `json:"information"`
	ImageUrl    string `json:"image_url"`
}

type ArtistLabel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
