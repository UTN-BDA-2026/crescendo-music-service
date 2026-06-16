CREATE INDEX idx_playlists_user_id
ON playlists(user_id);

CREATE INDEX idx_artists_songs_song_id
ON artists_songs(song_id);

CREATE INDEX idx_artists_songs_artist_id
ON artists_songs(artist_id);

CREATE INDEX idx_artists_albums_artist_id
ON artists_albums(artist_id);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_artists_name_trgm
ON artists
USING gin (name gin_trgm_ops);