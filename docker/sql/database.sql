CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    register_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_of_birth DATE NOT NULL,
    profile_image_url TEXT
);

CREATE TABLE playlists (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL REFERENCES users(id)
);

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
	file_id TEXT NOT NULL,
	genre_id INT NOT NULL REFERENCES genres(id),
	duration INT,
	bpm INT,
	release_date DATE
);

CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
	type TEXT NOT NULL,
	genre_id INT NOT NULL REFERENCES genres(id),
    cover_image_url TEXT,
	release_date DATE
);

CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    information TEXT,
    image_url TEXT
);

CREATE TABLE artists_songs (
    artist_id INT NOT NULL REFERENCES artists(id),
    song_id INT NOT NULL REFERENCES songs(id),
    role TEXT,
    PRIMARY KEY(artist_id,song_id)
);

CREATE TABLE albums_songs (
    track_position INT,
    album_id INT NOT NULL REFERENCES albums(id),
    song_id INT NOT NULL REFERENCES songs(id),
    PRIMARY KEY (track_position, album_id)
);

CREATE TABLE artists_albums (
    artist_id INT NOT NULL REFERENCES artists(id),
    album_id INT NOT NULL REFERENCES albums(id),
    PRIMARY KEY(artist_id,album_id)
);