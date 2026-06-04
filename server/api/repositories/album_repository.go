package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type AlbumRepository struct {
	databaseContext database.DBTX
}

func NewAlbumRepository(db database.DBTX) AlbumRepository {
	repository := AlbumRepository{
		databaseContext: db,
	}

	return repository
}

func (r AlbumRepository) Create(album models.Album) (int, error) {
	var id int

	err := r.databaseContext.QueryRow(`
		INSERT INTO albums (
			title,
			type,
			genre_id,
			cover_image_url,
			release_date
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`,
		album.Title,
		album.Type,
		album.GenreId,
		album.CoverImageUrl,
		album.ReleaseDate,
	).Scan(&id)

	return id, err
}

func (r AlbumRepository) GetById(id int) (models.Album, error) {
	row := r.databaseContext.QueryRow(`
			SELECT 
				id, 
				title,
				type,
				genre_id,
				cover_image_url,
				release_date
			FROM albums
			WHERE id = $1
		`, id)
	var album models.Album
	err := row.Scan(
		&album.Id,
		&album.Title,
		&album.Type,
		&album.GenreId,
		&album.CoverImageUrl,
		&album.ReleaseDate,
	)

	album.ReleaseDate = album.ReleaseDate.UTC()
	return album, err
}

func (r AlbumRepository) Update(album models.Album) (models.Album, error) {
	var updated models.Album

	err := r.databaseContext.QueryRow(`
		UPDATE albums
		SET title = $1,
			type = $2,
			genre_id = $3,
			cover_image_url = $4,
			release_date = $5
		WHERE id = $6
		RETURNING id, title, type, genre_id, cover_image_url, release_date
	`,
		album.Title,
		album.Type,
		album.GenreId,
		album.CoverImageUrl,
		album.ReleaseDate,
		album.Id,
	).Scan(
		&updated.Id,
		&updated.Title,
		&updated.Type,
		&updated.GenreId,
		&updated.CoverImageUrl,
		&updated.ReleaseDate,
	)

	updated.ReleaseDate = updated.ReleaseDate.UTC()

	return updated, err
}

func (r AlbumRepository) Delete(id int) error {
	result, err := r.databaseContext.Exec(`
		DELETE FROM albums
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
		return fmt.Errorf("no album found with id %d", id)
	}

	return nil
}
