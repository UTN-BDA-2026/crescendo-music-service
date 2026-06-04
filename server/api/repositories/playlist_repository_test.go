package repositories_test

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlaylistCRUD(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	playlist := models.Playlist{
		Id:           6,
		Title:        "TestTitle",
		Description:  "Texto de descripción",
		CreationDate: time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC),
	}

	user := models.User{
		Username:          "TestUsername",
		Email:             "testmail@mail.com",
		PasswordHash:      "0x3556FF",
		RegisterDate:      time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC),
		DateOfBirth:       time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC),
		ProfilePictureUrl: "files/grrs.png",
	}

	t.Run("Create", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		userRepository := repositories.NewUserRepository(transaction)
		repository := repositories.NewPlaylistRepository(transaction)

		userId, err := userRepository.Create(user)
		check.NoError(err)
		check.NotZero(userId)

		expectedPlaylist := playlist
		expectedPlaylist.UserId = userId

		id, err := repository.Create(expectedPlaylist)
		check.NoError(err)
		check.NotZero(id)

		expectedPlaylist.Id = id

		createdPlaylist, err := repository.GetById(expectedPlaylist.Id)
		check.NoError(err)
		check.NotEmpty(createdPlaylist)

		check.Equal(expectedPlaylist, createdPlaylist)
	})

	t.Run("Read", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		userRepository := repositories.NewUserRepository(transaction)
		repository := repositories.NewPlaylistRepository(transaction)

		userId, err := userRepository.Create(user)
		check.NoError(err)
		check.NotZero(userId)

		refPlaylist := playlist
		refPlaylist.UserId = userId
		refPlaylist.Id, err = repository.Create(refPlaylist)
		check.NoError(err)
		check.NotEqual(0, refPlaylist.Id)
		fetchedPlaylist, err := repository.GetById(refPlaylist.Id)

		check.NoError(err)
		check.Equal(refPlaylist, fetchedPlaylist)
	})

	t.Run("Update", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		userRepository := repositories.NewUserRepository(transaction)
		repository := repositories.NewPlaylistRepository(transaction)

		userId, err := userRepository.Create(user)
		check.NoError(err)
		check.NotZero(userId)

		refPlaylist := playlist
		refPlaylist.UserId = userId
		refPlaylist.Id, err = repository.Create(refPlaylist)

		check.NoError(err)
		check.NotEqual(0, refPlaylist.Id)

		changedPlaylist := refPlaylist
		changedPlaylist.Title = "Galahad"

		updatedPlaylist, err := repository.Update(changedPlaylist)

		check.NoError(err)

		readPlaylist, err := repository.GetById(refPlaylist.Id)

		check.NoError(err)
		check.NotEmpty(readPlaylist)

		check.Equal(updatedPlaylist, readPlaylist)

		check.NotEqual(refPlaylist.Title, updatedPlaylist.Title)
	})
	t.Run("Delete", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		userRepository := repositories.NewUserRepository(transaction)
		repository := repositories.NewPlaylistRepository(transaction)

		userId, err := userRepository.Create(user)
		check.NoError(err)
		check.NotZero(userId)

		refPlaylist := playlist
		refPlaylist.UserId = userId
		refPlaylist.Id, err = repository.Create(refPlaylist)
		id, err := repository.Create(refPlaylist)

		check.NoError(err)
		check.NotZero(id)

		err = repository.Delete(id)

		check.NoError(err)

		deletedPlaylist, err := repository.GetById(id)

		check.Error(err)
		check.Zero(deletedPlaylist.Id)
	})
}
