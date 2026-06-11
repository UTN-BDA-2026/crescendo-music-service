package repositories_test

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"crescendo-api/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		check := require.New(t)

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
		check := require.New(t)

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
		check := require.New(t)

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
		check := require.New(t)

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

func TestArtistAlbumRelation(t *testing.T) {
	check := require.New(t)

	db, err := database.NewConnection()

	check.NoError(err)

	transaction, err := db.Begin()

	t.Cleanup(func() {
		transaction.Rollback()
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	check.NoError(err)

	artist := models.Artist{
		Name:        "ABBA",
		Information: "Description of the artist",
		ImageUrl:    "a/dfv/gf.png",
	}
	album := models.Album{
		Title:         "JJ",
		Type:          "EP",
		CoverImageUrl: "aaa/f.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}
	genre := models.Genre{
		Name: "Rock",
	}

	repository := repositories.NewArtistRepository(transaction)
	albumRepo := repositories.NewAlbumRepository(transaction)
	genreRepo := repositories.NewGenreRepository(transaction)

	artist.Id, err = repository.Create(artist)
	check.NoError(err)

	genre.Id, err = genreRepo.Create(genre)
	check.NoError(err)

	album.GenreId = genre.Id
	album.Id, err = albumRepo.Create(album)
	check.NoError(err)

	err = repository.AddAlbumToArtist(album.Id, artist.Id)
	check.NoError(err)

	expected := []models.AlbumPreview{
		{
			Id:            album.Id,
			Title:         album.Title,
			Type:          album.Type,
			CoverImageUrl: album.CoverImageUrl,
			ReleaseDate:   album.ReleaseDate,
		},
	}

	fetchedAlbums, err := repository.GetArtistAlbumPreviews(artist.Id)

	check.NoError(err)
	check.Len(fetchedAlbums, len(expected))
	check.Equal(expected[0], fetchedAlbums[0])
}
