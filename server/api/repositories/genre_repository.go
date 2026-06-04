package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type GenreRepository struct {
	databaseContext database.DBTX
}

func NewGenreRepository(db database.DBTX) GenreRepository {
	repository := GenreRepository{
		databaseContext: db,
	}

	return repository
}

func (r GenreRepository) Create(genre models.Genre) (int, error) {
	var id int

	err := r.databaseContext.QueryRow(`
		INSERT INTO genres (
			name
		)
		VALUES ($1)
		RETURNING id
	`,
		genre.Name,
	).Scan(&id)

	return id, err
}

func (r GenreRepository) GetById(id int) (models.Genre, error) {
	row := r.databaseContext.QueryRow(`
			SELECT id, name
			FROM genres
			WHERE id = $1
		`, id)
	var genre models.Genre
	err := row.Scan(
		&genre.Id,
		&genre.Name,
	)
	return genre, err
}

func (r GenreRepository) Update(genre models.Genre) (models.Genre, error) {
	var updated models.Genre

	err := r.databaseContext.QueryRow(`
		UPDATE genres
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`,
		genre.Name,
		genre.Id,
	).Scan(
		&updated.Id,
		&updated.Name,
	)
	return updated, err
}

func (r GenreRepository) Delete(id int) error {
	result, err := r.databaseContext.Exec(`
		DELETE FROM genres
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
		return fmt.Errorf("no genre found with id %d", id)
	}

	return nil
}
