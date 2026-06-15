package services_test

import (
	"crescendo-api/models"
	"os"
)

// Simulate .env variables for testing

func init() {
	os.Setenv("JWT_SECRET", "secret-key")
	os.Setenv("BASE_URL", "http://localhost:8080")
}

// User Repository Mock

type mockUserRepository struct {
	createFunc           func(user models.User) (int, error)
	getByIdFunc          func(id int) (models.User, error)
	updateFunc           func(user models.User) (models.User, error)
	deleteFunc           func(id int) error
	getByUsernameOrEmail func(username string, email string) (models.User, error)
}

func (m mockUserRepository) Create(user models.User) (int, error) {
	return m.createFunc(user)
}

func (m mockUserRepository) GetById(id int) (models.User, error) {
	return m.getByIdFunc(id)
}

func (m mockUserRepository) Update(user models.User) (models.User, error) {
	return m.updateFunc(user)
}

func (m mockUserRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

func (m mockUserRepository) GetByUsernameOrEmail(username string, email string) (models.User, error) {
	return m.getByUsernameOrEmail(username, email)
}

// Song Repository Mock

type mockSongRepository struct {
	createFunc                        func(user models.Song) (int, error)
	getByIdFunc                       func(id int) (models.Song, error)
	updateFunc                        func(user models.Song) (models.Song, error)
	deleteFunc                        func(id int) error
	getArtistsForPlaybackBySongIdFunc func(id int) ([]models.ArtistLabel, error)
	addArtistToSongFunc               func(artistId int, songId int) error
}

func (m mockSongRepository) Create(song models.Song) (int, error) {
	return m.createFunc(song)
}

func (m mockSongRepository) GetById(id int) (models.Song, error) {
	return m.getByIdFunc(id)
}

func (m mockSongRepository) Update(song models.Song) (models.Song, error) {
	return m.updateFunc(song)
}

func (m mockSongRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

func (m mockSongRepository) GetArtistsForPlaybackBySongId(id int) ([]models.ArtistLabel, error) {
	return m.getArtistsForPlaybackBySongIdFunc(id)
}

func (m mockSongRepository) AddArtistToSong(artistId int, songId int) error {
	return m.addArtistToSongFunc(artistId, songId)
}

// Album Repository Mock

type mockAlbumRepository struct {
	createFunc                     func(user models.Album) (int, error)
	getByIdFunc                    func(id int) (models.Album, error)
	updateFunc                     func(user models.Album) (models.Album, error)
	deleteFunc                     func(id int) error
	addSongToAlbumFunc             func(songId, albumId, trackPosition int) error
	getSongsPreviewFromAlbumIdFunc func(id int) ([]models.ListedSong, error)
}

func (m mockAlbumRepository) Create(album models.Album) (int, error) {
	return m.createFunc(album)
}

func (m mockAlbumRepository) GetById(id int) (models.Album, error) {
	return m.getByIdFunc(id)
}

func (m mockAlbumRepository) Update(album models.Album) (models.Album, error) {
	return m.updateFunc(album)
}

func (m mockAlbumRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

func (m mockAlbumRepository) AddSongToAlbum(songId, albumId, trackPosition int) error {
	return m.addSongToAlbumFunc(songId, albumId, trackPosition)
}

func (m mockAlbumRepository) GetSongsPreviewFromAlbumId(id int) ([]models.ListedSong, error) {
	return m.getSongsPreviewFromAlbumIdFunc(id)
}

// Genre Repository Mock

type mockGenreRepository struct {
	createFunc  func(genre models.Genre) (int, error)
	getByIdFunc func(id int) (models.Genre, error)
	updateFunc  func(genre models.Genre) (models.Genre, error)
	deleteFunc  func(id int) error
}

func (m mockGenreRepository) Create(genre models.Genre) (int, error) {
	return m.createFunc(genre)
}

func (m mockGenreRepository) GetById(id int) (models.Genre, error) {
	return m.getByIdFunc(id)
}

func (m mockGenreRepository) Update(genre models.Genre) (models.Genre, error) {
	return m.updateFunc(genre)
}

func (m mockGenreRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

// Artist Repository Mock

type mockArtistRepository struct {
	createFunc                     func(artist models.Artist) (int, error)
	getByIdFunc                    func(id int) (models.Artist, error)
	getAllFunc                     func() ([]models.Artist, error)
	updateFunc                     func(artist models.Artist) (models.Artist, error)
	deleteFunc                     func(id int) error
	addAlbumToArtistFunc           func(albumId, artistId int) error
	getAlbumPreviewsByArtistIdFunc func(id int) ([]models.AlbumPreview, error)
	getArtistSongPreviewsFunc      func(id int) ([]models.SongPreview, error)
}

func (m mockArtistRepository) Create(artist models.Artist) (int, error) {
	return m.createFunc(artist)
}

func (m mockArtistRepository) GetById(id int) (models.Artist, error) {
	return m.getByIdFunc(id)
}

func (m mockArtistRepository) Update(artist models.Artist) (models.Artist, error) {
	return m.updateFunc(artist)
}

func (m mockArtistRepository) Delete(id int) error {
	return m.deleteFunc(id)
}

func (m mockArtistRepository) AddAlbumToArtist(albumId, artistId int) error {
	return m.addAlbumToArtistFunc(albumId, artistId)
}

func (m mockArtistRepository) GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error) {
	return m.getAlbumPreviewsByArtistIdFunc(id)
}

func (m mockArtistRepository) GetArtistSongPreviews(id int) ([]models.SongPreview, error) {
	return m.getArtistSongPreviewsFunc(id)
}

func (m mockArtistRepository) GetAll() ([]models.Artist, error) {
	return m.getAllFunc()
}
