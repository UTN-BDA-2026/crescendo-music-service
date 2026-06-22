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
	AddArtistToSong(artistId int, songId int) error
	GetArtistsForPlaybackBySongId(id int) ([]models.ArtistLabel, error)
	FindByNameLike(name string) ([]models.SongPreviewWithArtists, error)
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

func (r songRepository) AddArtistToSong(artistId int, songId int) error {
	_, err := r.db.Exec(`
        INSERT INTO artists_songs (
            artist_id,
            song_id
        )
        VALUES ($1, $2)
    `,
		artistId,
		songId,
	)

	return err
}

func (r songRepository) GetArtistsForPlaybackBySongId(id int) ([]models.ArtistLabel, error) {
	rows, err := r.db.Query(`
		SELECT
			a.id,
			a.name
		FROM artists a
		JOIN artists_songs rel
			ON a.id = rel.artist_id
		WHERE rel.song_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.ArtistLabel

	for rows.Next() {
		var artist models.ArtistLabel

		err := rows.Scan(
			&artist.Id,
			&artist.Name,
		)
		if err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (r songRepository) FindByNameLike(name string) ([]models.SongPreviewWithArtists, error) {
	rows, err := r.db.Query(`
		SELECT
			s.id,
			s.title,
			s.duration,
			a.id,
			a.name
		FROM songs s
		JOIN artists_songs ars ON ars.song_id = s.id
		JOIN artists a ON a.id = ars.artist_id
		WHERE s.title ILIKE '%' || $1 || '%'
		ORDER BY similarity(s.title, $1) DESC
		LIMIT 20
	`, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songsMap := make(map[int]*models.SongPreviewWithArtists)

	for rows.Next() {
		var (
			songID     int
			title      string
			duration   int
			artistID   int
			artistName string
		)

		err := rows.Scan(
			&songID,
			&title,
			&duration,
			&artistID,
			&artistName,
		)
		if err != nil {
			return nil, err
		}

		song, exists := songsMap[songID]
		if !exists {
			song = &models.SongPreviewWithArtists{
				Id:       songID,
				Title:    title,
				Duration: duration,
				Artists:  []models.ArtistLabel{},
			}
			songsMap[songID] = song
		}

		song.Artists = append(song.Artists, models.ArtistLabel{
			Id:   artistID,
			Name: artistName,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	songs := make([]models.SongPreviewWithArtists, 0, len(songsMap))

	for _, song := range songsMap {
		songs = append(songs, *song)
	}

	return songs, nil
}
