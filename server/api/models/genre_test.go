package models_test

import (
	"crescendo-api/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGenre(t *testing.T) {
	genre := models.Genre{
		Id:   4,
		Name: "Rock",
	}

	check := assert.New(t)
	check.Equal(4, genre.Id)
	check.Equal("Rock", genre.Name)
}
