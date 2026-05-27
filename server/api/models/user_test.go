package models

import (
	"testing"
	"time"
)

func assertComparison(t *testing.T, reference any, result any) {
	if result != reference {
		t.Errorf("Expected= %v; obtained %v", reference, result)
	}
}

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
	assertComparison(t, uint(6), user.Id)
	assertComparison(t, "TestUsername", user.Username)
	assertComparison(t, "testmail@mail.com", user.Email)
	assertComparison(t, 0x3556FF, user.PasswordHash)
	assertComparison(t, time.Date(2024, 4, 23, 14, 30, 45, 0, time.Local), user.RegisterDate)
	assertComparison(t, time.Date(1999, 8, 15, 14, 30, 45, 100, time.Local), user.DateOfBirth)
	assertComparison(t, "files/grrs.png", user.ProfilePictureUrl)
}
