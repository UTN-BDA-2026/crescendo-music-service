package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type ArtistRepository interface {
	Create(artist models.Artist) (int, error)
	GetById(id int) (models.Artist, error)
	Update(artist models.Artist) (models.Artist, error)
	Delete(id int) error
	AddAlbumToArtist(albumId, artistId int) error
	GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error)
	GetArtistSongPreviews(id int) ([]models.SongPreview, error)
}

type artistRepository struct {
	databaseContext database.DBTX
}

func NewArtistRepository(db database.DBTX) ArtistRepository {
	repository := artistRepository{
		databaseContext: db,
	}

	return repository
}

func (r artistRepository) Create(artist models.Artist) (int, error) {
	var id int

	err := r.databaseContext.QueryRow(`
		INSERT INTO artists (
			name,
			information,
			image_url
		)
		VALUES ($1,$2,$3)
		RETURNING id
	`,
		artist.Name,
		artist.Information,
		artist.ImageUrl,
	).Scan(&id)

	return id, err
}

func (r artistRepository) GetById(id int) (models.Artist, error) {
	row := r.databaseContext.QueryRow(`
			SELECT id, name, information, image_url
			FROM artists
			WHERE id = $1
		`, id)
	var artist models.Artist
	err := row.Scan(
		&artist.Id,
		&artist.Name,
		&artist.Information,
		&artist.ImageUrl,
	)
	return artist, err
}

func (r artistRepository) Update(artist models.Artist) (models.Artist, error) {
	var updated models.Artist

	err := r.databaseContext.QueryRow(`
		UPDATE artists
		SET name = $1,
			information = $2,
			image_url = $3
		WHERE id = $4
		RETURNING id, name, information, image_url
	`,
		artist.Name,
		artist.Information,
		artist.ImageUrl,
		artist.Id,
	).Scan(
		&updated.Id,
		&updated.Name,
		&updated.Information,
		&updated.ImageUrl,
	)
	return updated, err
}

func (r artistRepository) Delete(id int) error {
	result, err := r.databaseContext.Exec(`
		DELETE FROM artists
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no artist found with id %d", id)
	}

	return nil
}

func (r artistRepository) AddAlbumToArtist(albumId, artistId int) error {
	_, err := r.databaseContext.Exec(`
        INSERT INTO artists_albums (
            artist_id,
            album_id
        )
        VALUES ($1, $2)
    `,
		artistId,
		albumId,
	)

	return err
}

func (r artistRepository) GetArtistAlbumPreviews(id int) ([]models.AlbumPreview, error) {
	rows, err := r.databaseContext.Query(`
		SELECT
			a.id,
			a.title,
			a.type,
			a.cover_image_url,
			a.release_date
		FROM albums a
		JOIN artists_albums rel
			ON a.id = rel.album_id
		WHERE rel.artist_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.AlbumPreview

	for rows.Next() {
		var album models.AlbumPreview

		err := rows.Scan(
			&album.Id,
			&album.Title,
			&album.Type,
			&album.CoverImageUrl,
			&album.ReleaseDate,
		)
		if err != nil {
			return nil, err
		}

		album.ReleaseDate = album.ReleaseDate.UTC()
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (r artistRepository) GetArtistSongPreviews(id int) ([]models.SongPreview, error) {
	rows, err := r.databaseContext.Query(`
		SELECT
			s.id,
			s.title,
			s.duration
		FROM songs s
		JOIN artists_songs rel
			ON s.id = rel.song_id
		WHERE rel.artist_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.SongPreview

	for rows.Next() {
		var song models.SongPreview

		err := rows.Scan(
			&song.Id,
			&song.Title,
			&song.Duration,
		)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}
