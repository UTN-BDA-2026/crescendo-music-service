package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type SongRepository interface {
	Create(user models.Song) (int, error)
	GetById(id int) (models.Song, error)
	Update(user models.Song) (models.Song, error)
	Delete(id int) error
}

type songRepository struct {
	db database.DBTX
}

func NewSongRepository(db database.DBTX) SongRepository {
	return &songRepository{
		db: db,
	}
}

func (r songRepository) Create(song models.Song) (int, error) {
	var id int

	err := r.db.QueryRow(`
		INSERT INTO songs (
			title,
			file_id,
			genre_id,
			duration,
			bpm,
			release_date
		)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id
	`,
		song.Title,
		song.FileId,
		song.GenreId,
		song.Duration,
		song.Bpm,
		song.ReleaseDate,
	).Scan(&id)

	return id, err
}

func (r songRepository) GetById(id int) (models.Song, error) {
	row := r.db.QueryRow(`
			SELECT 
				id, 
				title,
				file_id,
				genre_id,
				duration,
				bpm,
				release_date
			FROM songs
			WHERE id = $1
		`, id)
	var song models.Song
	err := row.Scan(
		&song.Id,
		&song.Title,
		&song.FileId,
		&song.GenreId,
		&song.Duration,
		&song.Bpm,
		&song.ReleaseDate,
	)

	song.ReleaseDate = song.ReleaseDate.UTC()
	return song, err
}

func (r songRepository) Update(song models.Song) (models.Song, error) {
	var updated models.Song

	err := r.db.QueryRow(`
		UPDATE songs
		SET title = $1,
			file_id = $2,
			genre_id = $3,
			duration = $4,
			bpm = $5,
			release_date = $6
		WHERE id = $7
		RETURNING id, title, file_id, genre_id, duration, bpm, release_date
	`,
		song.Title,
		song.FileId,
		song.GenreId,
		song.Duration,
		song.Bpm,
		song.ReleaseDate,
		song.Id,
	).Scan(
		&updated.Id,
		&updated.Title,
		&updated.FileId,
		&updated.GenreId,
		&updated.Duration,
		&updated.Bpm,
		&updated.ReleaseDate,
	)

	updated.ReleaseDate = updated.ReleaseDate.UTC()

	return updated, err
}

func (r songRepository) Delete(id int) error {
	result, err := r.db.Exec(`
		DELETE FROM songs
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
		return fmt.Errorf("no song found with id %d", id)
	}

	return nil
}
