package models_test

import (
	"crescendo-api/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	user := models.User{
		Id:                6,
		Username:          "TestUsername",
		Email:             "testmail@mail.com",
		PasswordHash:      "0x3556FF",
		RegisterDate:      time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC),
		DateOfBirth:       time.Date(1999, 8, 15, 14, 30, 45, 100, time.UTC),
		ProfilePictureUrl: "files/grrs.png",
	}
	check := assert.New(t)
	check.Equal(6, user.Id)
	check.Equal("TestUsername", user.Username)
	check.Equal("testmail@mail.com", user.Email)
	check.Equal("0x3556FF", user.PasswordHash)
	check.Equal(time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC), user.RegisterDate)
	check.Equal(time.Date(1999, 8, 15, 14, 30, 45, 100, time.UTC), user.DateOfBirth)
	check.Equal("files/grrs.png", user.ProfilePictureUrl)
}
