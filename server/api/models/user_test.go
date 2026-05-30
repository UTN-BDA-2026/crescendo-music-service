package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	user := User{
		Id:                6,
		Username:          "TestUsername",
		Email:             "testmail@mail.com",
		PasswordHash:      0x3556FF,
		RegisterDate:      time.Date(2024, 4, 23, 14, 30, 45, 0, time.Local),
		DateOfBirth:       time.Date(1999, 8, 15, 14, 30, 45, 100, time.Local),
		ProfilePictureUrl: "files/grrs.png",
	}
	check := assert.New(t)
	check.Equal(uint(6), user.Id)
	check.Equal("TestUsername", user.Username)
	check.Equal("testmail@mail.com", user.Email)
	check.Equal(0x3556FF, user.PasswordHash)
	check.Equal(time.Date(2024, 4, 23, 14, 30, 45, 0, time.Local), user.RegisterDate)
	check.Equal(time.Date(1999, 8, 15, 14, 30, 45, 100, time.Local), user.DateOfBirth)
	check.Equal("files/grrs.png", user.ProfilePictureUrl)
}
