package repositories

import (
	"crescendo-api/database"
	"crescendo-api/models"
	"fmt"
)

type AlbumRepository interface {
	Create(user models.Album) (int, error)
	GetById(id int) (models.Album, error)
	Update(user models.Album) (models.Album, error)
	Delete(id int) error
	AddSongToAlbum(songId int, albumId int, trackPosition int) error
	GetSongsPreviewFromAlbumId(id int) ([]models.ListedSong, error)
}

type albumRepository struct {
	db database.DBTX
}

func NewAlbumRepository(db database.DBTX) AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

func (r albumRepository) Create(album models.Album) (int, error) {
	var id int

	err := r.db.QueryRow(`
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

func (r albumRepository) GetById(id int) (models.Album, error) {
	row := r.db.QueryRow(`
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

func (r albumRepository) Update(album models.Album) (models.Album, error) {
	var updated models.Album

	err := r.db.QueryRow(`
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

func (r albumRepository) Delete(id int) error {
	result, err := r.db.Exec(`
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

func (r albumRepository) AddSongToAlbum(songId int, albumId int, trackPosition int) error {
	_, err := r.db.Exec(`
        INSERT INTO albums_songs (
            album_id,
            song_id,
			track_position
        )
        VALUES ($1, $2,$3)
    `,
		albumId,
		songId,
		trackPosition,
	)

	return err
}
func (r albumRepository) GetSongsPreviewFromAlbumId(albumId int) ([]models.ListedSong, error) {
	rows, err := r.db.Query(`
		SELECT
			s.id,
			s.title,
			s.duration,
			rel.track_position,
			a.id,
			a.name
		FROM albums_songs rel
		JOIN songs s
			ON s.id = rel.song_id
		JOIN artists_songs ars
			ON ars.song_id = s.id
		JOIN artists a
			ON a.id = ars.artist_id
		WHERE rel.album_id = $1
		ORDER BY rel.track_position, a.name
	`, albumId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songsByID := make(map[int]*models.ListedSong)
	var orderedSongs []models.ListedSong

	for rows.Next() {
		var (
			songID        int
			title         string
			duration      int
			trackPosition int
			artistID      int
			artistName    string
		)

		if err := rows.Scan(
			&songID,
			&title,
			&duration,
			&trackPosition,
			&artistID,
			&artistName,
		); err != nil {
			return nil, err
		}

		song, exists := songsByID[songID]
		if !exists {
			orderedSongs = append(orderedSongs, models.ListedSong{
				SongPreview: models.SongPreview{
					Id:       songID,
					Title:    title,
					Duration: duration,
				},
				TrackPosition: trackPosition,
			})

			song = &orderedSongs[len(orderedSongs)-1]
			songsByID[songID] = song
		}

		song.Artists = append(song.Artists, models.ArtistLabel{
			Id:   artistID,
			Name: artistName,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderedSongs, nil
}
