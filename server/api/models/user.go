package models

import "time"

type User struct {
	Id                uint
	Username          string
	Email             string
	PasswordHash      int
	RegisterDate      time.Time
	DateOfBirth       time.Time
	ProfilePictureUrl string
}
