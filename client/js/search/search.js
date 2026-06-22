import { API_BASE_URL } from '../config/api.js';
import { playQueue } from '../player/player.js';
import { showAlbumView } from '../ui/views.js';

export function initSearch() {
    document.getElementById('global-search-form')?.addEventListener('submit', handleSearch);
}

async function handleSearch(e) {
    e.preventDefault();

    const query = document.getElementById('search-query').value.trim();
    if (!query) return;

    try {
        const data = await fetchSearch(query);

        renderArtists(data.artists || []);
        renderAlbums(data.albums || []);
        renderSongs(data.songs || []);

        renderEmptyState(query, data);

    } catch (error) {
        alert(error.message || 'Error de conexión al realizar la búsqueda.');
    }
}

async function fetchSearch(query) {
    const token = localStorage.getItem('crescendo_token');

    const response = await fetch(
        `${API_BASE_URL}/search?q=${encodeURIComponent(query)}&type=all`,
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        }
    );

    if (!response.ok) {
        throw new Error('Error al realizar la búsqueda');
    }

    return await response.json();
}

function renderArtists(artists) {
    const container = document.getElementById('artist-results');
    container.innerHTML = '';

    if (!artists.length) return;

    container.innerHTML = `
        <h3 class="section-title">
            Artistas
        </h3>
    `;

    artists.forEach(artist => {
        const div = document.createElement('div');
        div.className = 'glass-card';

        div.innerHTML = `
            <div class="card-content">
                ${artist.image_url ? `
                    <img src="${artist.image_url}" 
                         alt="${artist.name}"
                         class="artist-img">
                ` : ''}
                
                <div>
                    <h4 class="card-title">
                        ${artist.name}
                    </h4>
                    <p class="card-subtitle">
                        ${artist.information || 'Sin información adicional'}
                    </p>
                </div>
            </div>
        `;

        container.appendChild(div);
    });
}

function renderAlbums(albums) {
    const container = document.getElementById('album-results');
    container.innerHTML = '';

    if (!albums.length) return;

    container.innerHTML = `
        <h3 class="section-title">
            Álbumes
        </h3>
    `;

    albums.forEach(album => {
        const div = document.createElement('div');
        div.className = 'glass-card';

        div.innerHTML = `
            <div class="card-content">
                ${album.cover_image_url ? `
                    <img src="${album.cover_image_url}"
                         alt="${album.title}"
                         class="album-img">
                ` : ''}

                <div>
                    <h4 class="card-title">
                        ${album.title}
                        <span class="badge">
                            ${album.type}
                        </span>
                    </h4>

                    <p class="card-subtitle">
                        Lanzamiento: ${new Date(album.release_date).toLocaleDateString()}
                    </p>
                </div>
            </div>
        `;
        
        div.style.cursor = 'pointer';
        div.addEventListener('click', () => showAlbumView(album.id));

        container.appendChild(div);
    });
}

function renderSongs(songs) {
    const container = document.getElementById('song-results');
    container.innerHTML = '';

    if (!songs.length) return;

    container.innerHTML = `
        <h3 class="section-title">
            Canciones
        </h3>
    `;

    songs.forEach((song, idx) => {
        const div = document.createElement('div');
        div.className = 'song-item glass-card';

        const m = Math.floor(song.duration / 60);
        const s = (song.duration % 60).toString().padStart(2, '0');

        const artistNames = song.artists?.length
            ? song.artists.map(a => a.name).join(', ')
            : 'Desconocido';

        div.innerHTML = `
            <div>
                <strong class="song-title">
                    ${song.title}
                </strong>

                <span class="song-artist">
                    - ${artistNames}
                </span>

                <div class="song-duration">
                    Duración: ${m}:${s}
                </div>
            </div>
        `;

        const btn = document.createElement('button');
        btn.textContent = '▶ Play';
        btn.className = 'play-btn';

        btn.addEventListener('click', () => {
            playQueue(songs, idx);
        });

        div.appendChild(btn);
        container.appendChild(div);
    });
}

function renderEmptyState(query, data) {
    const artistData = data.artists || [];
    const albumData = data.albums || [];
    const songData = data.songs || [];

    if (
        artistData.length === 0 &&
        albumData.length === 0 &&
        songData.length === 0
    ) {
        document.getElementById('artist-results').innerHTML =
            `<p class="empty-state-text">
                No se encontraron resultados para "${query}".
            </p>`;
    }
}