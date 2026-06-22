BEGIN;

-- =====================================================
-- GENRES
-- =====================================================

INSERT INTO genres (name)
VALUES
('Rock'),
('Pop'),
('Hip Hop'),
('Reggaeton'),
('Electronic'),
('Jazz'),
('Blues'),
('Classical'),
('Indie'),
('Metal'),
('Trap'),
('House'),
('Techno'),
('R&B'),
('Soul'),
('Funk'),
('Latin'),
('K-Pop'),
('Country'),
('Ambient');

-- =====================================================
-- USERS (5000)
-- =====================================================

INSERT INTO users (
    username,
    email,
    password_hash,
    date_of_birth,
    profile_image_url
)
SELECT
    lower(
        first_name || '_' ||
        last_name || '_' || g
    ),
    lower(
        first_name || '.' ||
        last_name || g || '@mail.com'
    ),
    md5('pass_' || g),
    DATE '1970-01-01' + ((random()*18000)::int),
    'https://picsum.photos/200?user=' || g
FROM generate_series(1,5000) g
CROSS JOIN LATERAL (
    SELECT
        (ARRAY[
            'Lucas','Mateo','Sofia','Emma','Valentina',
            'Martina','Joaquin','Tomas','Camila','Nicolas',
            'Julieta','Agustin','Thiago','Olivia','Benjamin',
            'Lautaro','Paula','Mia','Renata','Bruno'
        ])[((g % 20)+1)] AS first_name,
        (ARRAY[
            'Gomez','Perez','Lopez','Rodriguez','Fernandez',
            'Diaz','Torres','Sanchez','Acosta','Ruiz',
            'Molina','Castro','Silva','Vega','Herrera',
            'Suarez','Romero','Medina','Navarro','Rojas'
        ])[(((g*7)%20)+1)] AS last_name
) n;

-- =====================================================
-- ARTISTS (1000)
-- =====================================================

INSERT INTO artists (
    name,
    information,
    image_url
)
SELECT
    CASE (g % 4)
        WHEN 0 THEN solo_name
        WHEN 1 THEN first_name || ' ' || last_name
        WHEN 2 THEN 'The ' || adjective || ' ' || noun
        ELSE adjective || ' ' || noun
    END AS name,

    genre_desc || '. ' || style_desc || '. ' || career_desc || '.' AS information,

    'https://picsum.photos/300?artist=' || g
FROM generate_series(1,1000) g
CROSS JOIN LATERAL (
    SELECT
        solo_names[((g - 1) % array_length(solo_names, 1)) + 1] AS solo_name,

        first_names[((g - 1) % array_length(first_names, 1)) + 1] AS first_name,

        last_names[(((g - 1) / array_length(first_names, 1))
            % array_length(last_names, 1)) + 1] AS last_name,

        adjectives[((g - 1) % array_length(adjectives, 1)) + 1] AS adjective,

        nouns[(((g - 1) / array_length(adjectives, 1))
            % array_length(nouns, 1)) + 1] AS noun,

        genres[((g * 3) % array_length(genres, 1)) + 1] AS genre_desc,

        styles[((g * 7) % array_length(styles, 1)) + 1] AS style_desc,

        careers[((g * 11) % array_length(careers, 1)) + 1] AS career_desc

    FROM (
        SELECT
            ARRAY[
                'Aurora','Nova','Eclipse','Luna','Solstice',
                'Zenith','Ember','Atlas','Echo','Phoenix',
                'Sierra','Orion','Lyric','Vega','Aria',
                'Rogue','Indigo','Halo','Cascade','Vertex'
            ] AS solo_names,

            ARRAY[
                'Alex','Luna','Mia','Leo','Nora',
                'Theo','Ruby','Kai','Ivy','Milo',
                'Aria','Noah','Zoe','Liam','Jade',
                'Ethan','Ava','Mason','Willow','Ezra'
            ] AS first_names,

            ARRAY[
                'Nova','Rivers','Stone','Vale','Knight',
                'Woods','Skye','Blake','Hart','Cross',
                'Lane','Fox','Reed','Brooks','West',
                'Ray','Cole','Storm','Frost','Banks'
            ] AS last_names,

            ARRAY[
                'Midnight','Golden','Silent','Electric','Scarlet',
                'Crystal','Urban','Wild','Burning','Velvet',
                'Hidden','Parallel','Northern','Solar','Neon',
                'Cosmic','Ancient','Radiant','Magnetic','Endless'
            ] AS adjectives,

            ARRAY[
                'Waves','Horizon','Bloom','Frequency','Mirage',
                'Lights','Storm','Signal','Orbit','Motion',
                'Galaxy','Forest','Journey','Rhythm','Vision',
                'Whisper','Beacon','Empire','Ocean','Theory'
            ] AS nouns,

            ARRAY[
                'Electronic music producer',
                'Alternative rock band',
                'Indie pop artist',
                'Hip-hop performer',
                'Jazz fusion musician',
                'Ambient composer',
                'Folk singer-songwriter',
                'Progressive metal project',
                'Synthwave producer',
                'Experimental electronic act'
            ] AS genres,

            ARRAY[
                'Known for atmospheric soundscapes',
                'Blending modern and classic influences',
                'Recognized for energetic live performances',
                'Focused on emotional storytelling',
                'Combining digital and acoustic elements',
                'Exploring cinematic arrangements',
                'Creating immersive musical experiences',
                'Inspired by urban culture and travel',
                'Mixing nostalgic and futuristic sounds',
                'Pushing the boundaries of genre conventions'
            ] AS styles,

            ARRAY[
                'Has built a loyal international fanbase',
                'Regularly collaborates with emerging artists',
                'Has released multiple acclaimed records',
                'Frequently appears at music festivals',
                'Continues to evolve with each release',
                'Gained popularity through streaming platforms',
                'Maintains an active touring schedule',
                'Has earned recognition from critics worldwide',
                'Is known for a distinctive artistic identity',
                'Consistently attracts new listeners globally'
            ] AS careers
    ) t
) a;

-- =====================================================
-- ALBUMS (2000)
-- =====================================================

INSERT INTO albums (
    title,
    type,
    genre_id,
    cover_image_url,
    release_date
)
SELECT
    adjective || ' ' || noun || ' Vol. ' || ((g % 8)+1),
    (ARRAY['LP','EP','Single'])[((g % 3)+1)],
    ((g % 20)+1),
    'https://picsum.photos/400?album=' || g,
    DATE '1995-01-01' + ((random()*11000)::int)
FROM generate_series(1,2000) g
CROSS JOIN LATERAL (
    SELECT
        (ARRAY[
            'Lost','Hidden','Golden','Broken',
            'Electric','Silent','Infinite','Dark',
            'Parallel','Crimson','Neon','Fading',
            'Burning','Silver','Midnight','Urban',
            'Velvet','Solar','Cold','Endless'
        ])[((g % 20)+1)] AS adjective,
        (ARRAY[
            'Horizons','Memories','Skies','Echoes',
            'Roads','Stories','Dreams','Lights',
            'Voices','Nights','Reflections','Waves',
            'Signals','Visions','Patterns','Shapes',
            'Moments','Distances','Hearts','Frequencies'
        ])[(((g*11)%20)+1)] AS noun
) x;

-- =====================================================
-- ARTISTS -> ALBUMS
-- =====================================================

INSERT INTO artists_albums (
    artist_id,
    album_id
)
SELECT
    ((id - 1) % 1000) + 1,
    id
FROM albums;

-- =====================================================
-- SONGS (20000)
-- =====================================================

INSERT INTO songs (
    title,
    file_id,
    genre_id,
    duration,
    bpm,
    release_date
)
SELECT
    adjective || ' ' || noun,
    md5('song_' || g),
    ((g % 20) + 1),
    120 + (random() * 240)::int,
    70 + (random() * 100)::int,
    DATE '2000-01-01' + ((random() * 9000)::int)
FROM generate_series(1,20000) g
CROSS JOIN LATERAL (
    SELECT
        adjectives[((g - 1) % array_length(adjectives, 1)) + 1] AS adjective,
        nouns[(((g - 1) / array_length(adjectives, 1))
               % array_length(nouns, 1)) + 1] AS noun
    FROM (
        SELECT
            ARRAY[
                'Midnight','Golden','Broken','Silent','Electric','Fading',
                'Parallel','Burning','Digital','Cold','Higher','Hidden',
                'Lost','Dark','Crimson','Blue','Neon','Infinite','Urban','Velvet',
                'Silver','Amber','Crystal','Scarlet','Cosmic','Wild',
                'Ancient','Bright','Frozen','Lonely','Rapid','Gentle',
                'Restless','Secret','Radiant','Magnetic','Electric','Endless',
                'Dreaming','Fallen'
            ] AS adjectives,
            ARRAY[
                'Heart','Dream','Memory','Sky','Light','Shadow','River','Signal',
                'Horizon','Wave','Road','Voice','Fire','Motion','Storm','Love',
                'Night','Distance','Frequency','Echo',
                'Ocean','Thunder','Machine','Galaxy','Forest','Sunrise',
                'Sunset','Rain','Mirror','Whisper','Pulse','Dust',
                'Flame','Silence','Vision','Journey','Rhythm','Star',
                'Moon','Beacon'
            ] AS nouns
    ) t
) s;

-- =====================================================
-- ALBUMS -> SONGS
-- 10 canciones por álbum
-- =====================================================

INSERT INTO albums_songs (
    track_position,
    album_id,
    song_id
)
SELECT
    ((song_id - 1) % 10) + 1,
    ((song_id - 1) / 10) + 1,
    song_id
FROM generate_series(1,20000) song_id;

-- =====================================================
-- ARTISTA PRINCIPAL DE CADA CANCIÓN
-- =====================================================

INSERT INTO artists_songs (
    artist_id,
    song_id,
    role
)
SELECT
    aa.artist_id,
    als.song_id,
    'main'
FROM albums_songs als
JOIN artists_albums aa
    ON aa.album_id = als.album_id;

-- =====================================================
-- FEATURINGS
-- =====================================================

INSERT INTO artists_songs (
    artist_id,
    song_id,
    role
)
SELECT
    ((song_id * 37) % 1000) + 1,
    song_id,
    'feat'
FROM generate_series(1,20000) song_id
WHERE song_id % 4 = 0
ON CONFLICT DO NOTHING;

-- =====================================================
-- PLAYLISTS (50000)
-- =====================================================

INSERT INTO playlists (
    title,
    description,
    user_id
)
SELECT
    mood || ' Mix #' || g,
    'Curated playlist for ' || mood,
    ((g - 1) % 5000) + 1
FROM generate_series(1,50000) g
CROSS JOIN LATERAL (
    SELECT
        (ARRAY[
            'Chill',
            'Workout',
            'Focus',
            'Party',
            'Driving',
            'Coding',
            'Relax',
            'Summer',
            'Night',
            'Morning',
            'Electronic',
            'Rock',
            'Latin',
            'Indie',
            'Deep House',
            'Jazz'
        ])[((g % 16)+1)] AS mood
) p;

COMMIT;

-- =====================================================
-- RESET SEQUENCES
-- =====================================================

SELECT setval(
    pg_get_serial_sequence('users','id'),
    COALESCE((SELECT MAX(id) FROM users),1)
);

SELECT setval(
    pg_get_serial_sequence('artists','id'),
    COALESCE((SELECT MAX(id) FROM artists),1)
);

SELECT setval(
    pg_get_serial_sequence('albums','id'),
    COALESCE((SELECT MAX(id) FROM albums),1)
);

SELECT setval(
    pg_get_serial_sequence('songs','id'),
    COALESCE((SELECT MAX(id) FROM songs),1)
);

SELECT setval(
    pg_get_serial_sequence('genres','id'),
    COALESCE((SELECT MAX(id) FROM genres),1)
);

SELECT setval(
    pg_get_serial_sequence('playlists','id'),
    COALESCE((SELECT MAX(id) FROM playlists),1)
);