package models

import "time"

type User struct {
	Id                int
	Username          string
	Email             string
	PasswordHash      string
	RegisterDate      time.Time
	DateOfBirth       time.Time
	ProfilePictureUrl string
}
