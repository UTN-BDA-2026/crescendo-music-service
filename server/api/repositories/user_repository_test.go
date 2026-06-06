package repositories_test

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

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

		repository := repositories.NewUserRepository(transaction)
		userId, err := repository.Create(user)
		check.NoError(err)
		check.NotZero(userId)

		expectedUser := user
		expectedUser.Id = userId

		createdUser, err := repository.GetById(expectedUser.Id)
		check.NoError(err)
		check.NotEmpty(createdUser)

		check.Equal(expectedUser, createdUser)
	})
	t.Run("Read", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewUserRepository(transaction)
		refUser := user
		refUser.Id, err = repository.Create(refUser)
		check.NoError(err)
		check.NotEqual(0, refUser.Id)
		fetchedUser, err := repository.GetById(refUser.Id)

		check.NoError(err)
		check.Equal(refUser, fetchedUser)
	})

	t.Run("Update", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewUserRepository(transaction)
		refUser := user
		refUser.Id, err = repository.Create(refUser)
		check.NoError(err)
		check.NotEqual(0, refUser.Id)

		changedUser := refUser
		changedUser.Username = "Galahad"

		updatedUser, err := repository.Update(changedUser)

		check.NoError(err)

		readUser, err := repository.GetById(refUser.Id)

		check.NoError(err)
		check.NotEmpty(readUser)

		check.Equal(updatedUser, readUser)

		check.NotEqual(refUser.Username, updatedUser.Username)
	})
	t.Run("Delete", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewUserRepository(transaction)

		userId, err := repository.Create(user)

		check.NoError(err)
		check.NotZero(userId)

		err = repository.Delete(userId)

		check.NoError(err)

		deletedUser, err := repository.GetById(userId)

		check.Error(err)
		check.Zero(deletedUser.Id)
	})
}

func TestUserGetByUsernameOrEmail(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	user := models.User{
		Username:          "TestUsername",
		Email:             "testmail@mail.com",
		PasswordHash:      "0x3556FF",
		RegisterDate:      time.Date(2024, 4, 23, 14, 30, 45, 0, time.UTC),
		DateOfBirth:       time.Date(1999, 8, 15, 0, 0, 0, 0, time.UTC),
		ProfilePictureUrl: "files/grrs.png",
	}
	t.Run("Username", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewUserRepository(transaction)
		refUser := user
		refUser.Id, err = repository.Create(refUser)
		check.NoError(err)
		check.NotEqual(0, refUser.Id)
		fetchedUser, err := repository.GetByUsernameOrEmail(refUser.Username, "")

		check.NoError(err)
		check.Equal(refUser, fetchedUser)
	})

	t.Run("Email", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewUserRepository(transaction)
		refUser := user
		refUser.Id, err = repository.Create(refUser)
		check.NoError(err)
		check.NotEqual(0, refUser.Id)
		fetchedUser, err := repository.GetByUsernameOrEmail("", refUser.Email)

		check.NoError(err)
		check.Equal(refUser, fetchedUser)
	})
}
