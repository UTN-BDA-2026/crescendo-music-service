package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type PlaylistRepository struct {
	databaseContext database.DBTX
}

func NewPlaylistRepository(db database.DBTX) PlaylistRepository {
	repository := PlaylistRepository{
		databaseContext: db,
	}

	return repository
}

func (r PlaylistRepository) Create(playlist models.Playlist) (int, error) {
	var id int

	err := r.databaseContext.QueryRow(`
		INSERT INTO playlists (
			title,
			description,
			creation_date,
			user_id
		)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`,
		playlist.Title,
		playlist.Description,
		playlist.CreationDate,
		playlist.UserId,
	).Scan(&id)

	return id, err
}

func (r PlaylistRepository) GetById(id int) (models.Playlist, error) {
	row := r.databaseContext.QueryRow(`
			SELECT id, title, description, creation_date, user_id
			FROM playlists
			WHERE id = $1
		`, id)
	var playlist models.Playlist
	err := row.Scan(
		&playlist.Id,
		&playlist.Title,
		&playlist.Description,
		&playlist.CreationDate,
		&playlist.UserId,
	)
	playlist.CreationDate = playlist.CreationDate.UTC() //Correccion de conversión de fechas de Postgres a Go
	return playlist, err
}

func (r PlaylistRepository) Update(playlist models.Playlist) (models.Playlist, error) {
	var updated models.Playlist

	err := r.databaseContext.QueryRow(`
		UPDATE playlists
		SET title = $1,
			description = $2,
			creation_date = $3,
			user_id = $4
		WHERE id = $5
		RETURNING id, title, description, creation_date, user_id
	`,
		playlist.Title,
		playlist.Description,
		playlist.CreationDate,
		playlist.UserId,
		playlist.Id,
	).Scan(
		&updated.Id,
		&updated.Title,
		&updated.Description,
		&updated.CreationDate,
		&updated.UserId,
	)
	updated.CreationDate = updated.CreationDate.UTC() //Correccion de conversión de fechas de Postgres a Go
	return updated, err
}

func (r PlaylistRepository) Delete(id int) error {
	result, err := r.databaseContext.Exec(`
		DELETE FROM playlists
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
		return fmt.Errorf("no playlist found with id %d", id)
	}

	return nil
}
