package repositories_test

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSongCRUD(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	song := models.Song{
		Id:          1,
		Title:       "Song Title",
		FileId:      "0xq2454",
		GenreId:     4,
		Duration:    253,
		Bpm:         110,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	genre := models.Genre{
		Name: "Rock",
	}

	t.Run("Create", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewSongRepository(transaction)

		expected := song

		genreId, err := genreRepository.Create(genre)

		check.NoError(err)
		check.NotZero(genreId)

		expected.GenreId = genreId

		id, err := repository.Create(expected)
		check.NoError(err)
		check.NotZero(id)

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

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewSongRepository(transaction)

		reference := song

		genreId, err := genreRepository.Create(genre)

		check.NoError(err)
		check.NotZero(genreId)

		reference.GenreId = genreId

		reference.Id, err = repository.Create(reference)
		check.NoError(err)
		check.NotZero(reference.Id)

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

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewSongRepository(transaction)

		reference := song

		genreId, err := genreRepository.Create(genre)

		check.NoError(err)
		check.NotZero(genreId)

		reference.GenreId = genreId
		reference.Id, err = repository.Create(reference)

		check.NoError(err)
		check.NotEqual(0, reference.Id)

		changed := reference
		changed.Title = "Galahad"

		updated, err := repository.Update(changed)

		check.NoError(err)

		read, err := repository.GetById(reference.Id)

		check.NoError(err)
		check.NotEmpty(read)

		check.Equal(updated, read)

		check.NotEqual(reference.Title, updated.Title)
	})
	t.Run("Delete", func(t *testing.T) {
		check := assert.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewSongRepository(transaction)

		reference := song

		genreId, err := genreRepository.Create(genre)

		check.NoError(err)
		check.NotZero(genreId)

		reference.GenreId = genreId
		id, err := repository.Create(reference)

		check.NoError(err)
		check.NotEqual(0, reference.Id)

		check.NoError(err)
		check.NotZero(id)

		err = repository.Delete(id)

		check.NoError(err)

		deleted, err := repository.GetById(id)

		check.Error(err)
		check.Zero(deleted.Id)
	})
}
