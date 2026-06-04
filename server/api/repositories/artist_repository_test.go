package repositories_test

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArtistCRUD(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	artist := models.Artist{
		Id:          5,
		Name:        "ABBA",
		Information: "Description of the artist",
		ImageUrl:    "a/dfv/gf.png",
	}

	t.Run("Create", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewArtistRepository(transaction)
		id, err := repository.Create(artist)
		check.NoError(err)
		check.NotZero(id)

		expected := artist
		expected.Id = id

		created, err := repository.GetById(expected.Id)
		check.NoError(err)
		check.NotEmpty(created)

		check.Equal(expected, created)
	})
	t.Run("Read", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewArtistRepository(transaction)
		reference := artist
		reference.Id, err = repository.Create(reference)
		check.NoError(err)
		check.NotEqual(0, reference.Id)
		fetched, err := repository.GetById(reference.Id)

		check.NoError(err)
		check.Equal(reference, fetched)
	})

	t.Run("Update", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewArtistRepository(transaction)
		reference := artist
		reference.Id, err = repository.Create(reference)
		check.NoError(err)
		check.NotEqual(0, reference.Id)

		changed := reference
		changed.Name = "Galahad"

		updated, err := repository.Update(changed)

		check.NoError(err)

		read, err := repository.GetById(reference.Id)

		check.NoError(err)
		check.NotEmpty(read)

		check.Equal(updated, read)

		check.NotEqual(reference.Name, updated.Name)
	})
	t.Run("Delete", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		repository := repositories.NewArtistRepository(transaction)

		id, err := repository.Create(artist)

		check.NoError(err)
		check.NotZero(id)

		err = repository.Delete(id)

		check.NoError(err)

		deleted, err := repository.GetById(id)

		check.Error(err)
		check.Zero(deleted.Id)
	})
}
