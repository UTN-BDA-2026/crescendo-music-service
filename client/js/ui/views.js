import { API_BASE_URL } from '../config/api.js';
import { playQueue } from '../player/player.js';

export async function showAlbumView(albumId) {
    const searchResults = document.getElementById('search-results');
    const detailView = document.getElementById('detail-view');
    const detailContent = document.getElementById('detail-content');


    searchResults.classList.add('hidden');
    detailView.classList.remove('hidden');

    detailContent.innerHTML = '<p>Cargando...</p>';

    try {
        const token = localStorage.getItem('crescendo_token');
        const response = await fetch(`${API_BASE_URL}/albums/${albumId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('No se pudo cargar el álbum');
        }

        const album = await response.json();
        renderAlbumDetails(album, detailContent);
    } catch (error) {
        detailContent.innerHTML = `<p class="error">${error.message}</p>`;
    }
}

function renderAlbumDetails(album, container) {
    const releaseDate = new Date(album.release_date).toLocaleDateString();

    let html = `
        <div class="detail-header">
            ${album.cover_image_url ? `<img src="${album.cover_image_url}" referrerpolicy="no-referrer" alt="${album.title}" class="detail-img">` : ''}
            <div class="detail-info">
                <h2>${album.title}</h2>
                <p><strong>Tipo:</strong> ${album.type}</p>
                <p><strong>Género:</strong> ${album.genre ? album.genre.name : 'Desconocido'}</p>
                <p><strong>Lanzamiento:</strong> ${releaseDate}</p>
            </div>
        </div>
        <div class="detail-tracks">
            <h3>Canciones</h3>
            <div class="song-list">
    `;

    if (!album.songs || album.songs.length === 0) {
        html += `<p>No hay canciones en este álbum.</p>`;
    } else {
        album.songs.forEach((song, idx) => {
            const m = Math.floor(song.duration / 60);
            const s = (song.duration % 60).toString().padStart(2, '0');
            const artistNames = song.artists?.length
                ? song.artists.map(a => a.name).join(', ')
                : 'Desconocido';


        });
    }

    html += `
            </div>
        </div>
    `;

    container.innerHTML = html;

    if (album.songs && album.songs.length > 0) {
        const songList = container.querySelector('.song-list');
        album.songs.forEach((song, idx) => {
            const div = document.createElement('div');
            div.className = 'song-item glass-card';

            const m = Math.floor(song.duration / 60);
            const s = (song.duration % 60).toString().padStart(2, '0');
            const artistLinksHTML = song.artists?.length
                ? song.artists.map(a => `<a href="#" class="album-artist-link" data-id="${a.id}" style="color: inherit; text-decoration: underline; cursor: pointer;">${a.name}</a>`).join(', ')
                : 'Desconocido';

            div.innerHTML = `
                <div>
                    <span class="track-number">${song.track_position}. </span>
                    <strong class="song-title">${song.title}</strong>
                    <span class="song-artist">- ${artistLinksHTML}</span>
                    <div class="song-duration">Duración: ${m}:${s}</div>
                </div>
            `;

            const artistLinks = div.querySelectorAll('.album-artist-link');
            artistLinks.forEach(link => {
                link.addEventListener('click', (e) => {
                    e.preventDefault();
                    showArtistView(link.getAttribute('data-id'));
                });
            });

            const btn = document.createElement('button');
            btn.textContent = '▶ Play';
            btn.className = 'play-btn';


            btn.addEventListener('click', () => {
                playQueue(album.songs, idx);
            });

            div.appendChild(btn);
            songList.appendChild(div);
        });
    }
}

export async function showArtistView(artistId) {
    const searchResults = document.getElementById('search-results');
    const detailView = document.getElementById('detail-view');
    const detailContent = document.getElementById('detail-content');

    searchResults.classList.add('hidden');
    detailView.classList.remove('hidden');

    detailContent.innerHTML = '<p>Cargando artista...</p>';

    try {
        const token = localStorage.getItem('crescendo_token');
        const headers = { 'Authorization': `Bearer ${token}` };

        const [artistRes, albumsRes, songsRes] = await Promise.all([
            fetch(`${API_BASE_URL}/artists/${artistId}`, { headers }),
            fetch(`${API_BASE_URL}/artists/${artistId}/albums`, { headers }),
            fetch(`${API_BASE_URL}/artists/${artistId}/songs`, { headers })
        ]);

        if (!artistRes.ok) throw new Error('No se pudo cargar el artista');

        const artist = await artistRes.json();
        const albums = albumsRes.ok ? await albumsRes.json() : [];
        const songs = songsRes.ok ? await songsRes.json() : [];

        renderArtistDetails(artist, albums, songs, detailContent);
    } catch (error) {
        detailContent.innerHTML = `<p class="error">${error.message}</p>`;
    }
}

function renderArtistDetails(artist, albums, songs, container) {
    let html = `
        <div class="detail-header">
            ${artist.image_url ? `<img src="${artist.image_url}" referrerpolicy="no-referrer" alt="${artist.name}" class="detail-img" style="border-radius: 50%;">` : ''}
            <div class="detail-info">
                <h2>${artist.name}</h2>
                <p>${artist.information || 'Sin información adicional'}</p>
            </div>
        </div>
    `;

    if (songs && songs.length > 0) {
        html += `
            <div class="detail-tracks" style="margin-top: 20px;">
                <h3>Canciones Populares</h3>
                <div class="song-list"></div>
            </div>
        `;
    }

    if (albums && albums.length > 0) {
        html += `
            <div class="detail-albums" style="margin-top: 20px;">
                <h3>Álbumes</h3>
                <div class="album-grid" style="display: flex; gap: 15px; flex-wrap: wrap;"></div>
            </div>
        `;
    }

    container.innerHTML = html;

    if (songs && songs.length > 0) {
        const songList = container.querySelector('.song-list');
        songs.forEach((song, idx) => {
            const div = document.createElement('div');
            div.className = 'song-item glass-card';

            const m = Math.floor(song.duration / 60);
            const s = (song.duration % 60).toString().padStart(2, '0');

            div.innerHTML = `
                <div>
                    <strong class="song-title">${song.title}</strong>
                    <div class="song-duration">Duración: ${m}:${s}</div>
                </div>
            `;

            const btn = document.createElement('button');
            btn.textContent = '▶ Play';
            btn.className = 'play-btn';

            btn.addEventListener('click', () => {
                playQueue(songs, idx);
            });

            div.appendChild(btn);
            songList.appendChild(div);
        });
    }


    if (albums && albums.length > 0) {
        const albumGrid = container.querySelector('.album-grid');
        albums.forEach(album => {
            const div = document.createElement('div');
            div.className = 'glass-card';
            div.style.width = '180px';
            div.style.cursor = 'pointer';

            div.innerHTML = `
                ${album.cover_image_url ? `<img src="${album.cover_image_url}" referrerpolicy="no-referrer" alt="${album.title}" class="album-img" style="width: 100%; height: 180px;">` : ''}
                <div style="padding-top: 10px;">
                    <h4 style="margin: 0;">${album.title}</h4>
                    <p style="margin: 5px 0 0 0; font-size: 0.85em; color: #6b7280;">${new Date(album.release_date).getFullYear()} • ${album.type}</p>
                </div>
            `;

            div.addEventListener('click', () => {
                showAlbumView(album.id);
            });

            albumGrid.appendChild(div);
        });
    }
}

export function initViews() {
    const backBtn = document.getElementById('back-to-search-btn');
    if (backBtn) {
        backBtn.addEventListener('click', () => {
            document.getElementById('detail-view').classList.add('hidden');
            document.getElementById('search-results').classList.remove('hidden');
        });
    }
}
