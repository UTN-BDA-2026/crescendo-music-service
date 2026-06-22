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
            ${album.cover_image_url ? `<img src="${album.cover_image_url}" alt="${album.title}" class="detail-img">` : ''}
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
            const artistNames = song.artists?.length
                ? song.artists.map(a => a.name).join(', ')
                : 'Desconocido';

            div.innerHTML = `
                <div>
                    <span class="track-number">${song.track_position}. </span>
                    <strong class="song-title">${song.title}</strong>
                    <span class="song-artist">- ${artistNames}</span>
                    <div class="song-duration">Duración: ${m}:${s}</div>
                </div>
            `;

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

export function initViews() {
    const backBtn = document.getElementById('back-to-search-btn');
    if (backBtn) {
        backBtn.addEventListener('click', () => {
            document.getElementById('detail-view').classList.add('hidden');
            document.getElementById('search-results').classList.remove('hidden');
        });
    }
}
