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
