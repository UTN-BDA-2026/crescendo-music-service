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

func TestAlbumCRUD(t *testing.T) {

	db, err := database.NewConnection()

	assert.NoError(t, err)

	t.Cleanup(func() {
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	album := models.Album{
		Id:            8,
		Title:         "JJ",
		Type:          "EP",
		GenreId:       4,
		CoverImageUrl: "aaa/f.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	genre := models.Genre{
		Name: "Rock",
	}

	t.Run("Create", func(t *testing.T) {
		check := require.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewAlbumRepository(transaction)

		expected := album

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
		check := require.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewAlbumRepository(transaction)

		reference := album

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
		check := require.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewAlbumRepository(transaction)

		reference := album

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
		check := require.New(t)

		transaction, err := db.Begin()

		t.Cleanup(func() {
			transaction.Rollback()
		})

		check.NoError(err)

		genreRepository := repositories.NewGenreRepository(transaction)
		repository := repositories.NewAlbumRepository(transaction)

		reference := album

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

func TestGetSongsPreviewFromAlbumId(t *testing.T) {
	check := require.New(t)

	db, err := database.NewConnection()

	check.NoError(err)

	transaction, err := db.Begin()

	t.Cleanup(func() {
		transaction.Rollback()
		db.Close() // Cerramos la conexion al terminar el test (exitoso o no)
	})

	check.NoError(err)

	album := models.Album{
		Title:         "JJ",
		Type:          "EP",
		CoverImageUrl: "aaa/f.png",
		ReleaseDate:   time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}
	genre := models.Genre{
		Name: "Rock",
	}

	song := models.Song{
		Title:       "Song Title",
		FileId:      "0xq2454",
		Duration:    253,
		Bpm:         110,
		ReleaseDate: time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC),
	}

	artist := models.Artist{
		Name: "Artist 1",
	}

	listedSong := models.ListedSong{
		SongPreviewWithArtists: models.SongPreviewWithArtists{
			Title:    song.Title,
			Duration: song.Duration,
			Artists: []models.ArtistLabel{
				{
					Name: artist.Name,
				},
			},
		},
		TrackPosition: 1,
	}

	albumRepo := repositories.NewAlbumRepository(transaction)
	artistRepo := repositories.NewArtistRepository(transaction)
	genreRepo := repositories.NewGenreRepository(transaction)
	songRepo := repositories.NewSongRepository(transaction)

	genre.Id, err = genreRepo.Create(genre)
	check.NoError(err)

	song.GenreId = genre.Id
	album.GenreId = genre.Id

	song.Id, err = songRepo.Create(song)
	check.NoError(err)
	listedSong.Id = song.Id

	album.Id, err = albumRepo.Create(album)
	check.NoError(err)

	artist.Id, err = artistRepo.Create(artist)
	check.NoError(err)
	listedSong.Artists[0].Id = artist.Id

	err = albumRepo.AddSongToAlbum(song.Id, album.Id, listedSong.TrackPosition)
	check.NoError(err)

	err = songRepo.AddArtistToSong(artist.Id, song.Id)

	songList, err := albumRepo.GetSongsPreviewFromAlbumId(album.Id)
	check.NoError(err)

	check.Len(songList, 1)
	check.Equal(listedSong, songList[0])
}

func TestSearchAlbumsByTitle(t *testing.T) {
	check := require.New(t)

	db, err := database.NewConnection()
	check.NoError(err)

	transaction, err := db.Begin()
	t.Cleanup(func() {
		transaction.Rollback()
		db.Close()
	})
	check.NoError(err)

	genreRepo := repositories.NewGenreRepository(transaction)
	repository := repositories.NewAlbumRepository(transaction)

	genre := models.Genre{
		Name: "Pop",
	}
	genre.Id, err = genreRepo.Create(genre)
	check.NoError(err)

	album := models.Album{
		Title:         "After Hours",
		Type:          "Album",
		GenreId:       genre.Id,
		CoverImageUrl: "afterhours.png",
		ReleaseDate:   time.Date(2020, 3, 20, 0, 0, 0, 0, time.UTC),
	}

	album.Id, err = repository.Create(album)
	check.NoError(err)

	results, err := repository.SearchByTitle("After Hours")
	check.NoError(err)
	check.NotEmpty(results)
	check.Equal(album.Title, results[0].Title)

	results, err = repository.SearchByTitle("hours")
	check.NoError(err)
	check.NotEmpty(results)
	check.Equal(album.Title, results[0].Title)

	results, err = repository.SearchByTitle("Starboy")
	check.NoError(err)
	check.Empty(results)
}
