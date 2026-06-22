import { API_BASE_URL } from '../config/api.js';
import { showArtistView } from '../ui/views.js';

let queue = [];
let index = -1;

let isRandom = false;
let loopMode = 'none';

export function playQueue(songsArray, startIndex = 0) {
    if (!songsArray || songsArray.length === 0) return;

    queue = [...songsArray];

    if (isRandom) {
        index = Math.floor(Math.random() * queue.length);
    } else {
        index = startIndex;
    }

    playCurrentSong();
}

async function playCurrentSong() {
    if (index < 0 || index >= queue.length) return;

    const song = queue[index];
    const token = localStorage.getItem('crescendo_token');

    try {
        const res = await fetch(
            `${API_BASE_URL}/songs/${song.id}/playback`,
            {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }
        );

        if (!res.ok) throw new Error("No se pudo obtener información de reproducción");

        const data = await res.json();

        const titleEl = document.getElementById('now-playing-title');
        const audioEl = document.getElementById('audio-player');

        const artistLinksHTML = data.artists?.length
            ? data.artists.map(a => `<a href="#" class="player-artist-link" data-id="${a.id}" style="color: inherit; text-decoration: underline; cursor: pointer;">${a.name}</a>`).join(', ')
            : 'Desconocido';

        titleEl.innerHTML = `${data.title} - ${artistLinksHTML}`;

        const artistLinks = titleEl.querySelectorAll('.player-artist-link');
        artistLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                showArtistView(link.getAttribute('data-id'));
            });
        });

        audioEl.src = data.stream_url;
        audioEl.play();

    } catch (error) {
        alert(error.message || "Error al reproducir la canción");
    }
}

function playNext() {
    if (queue.length === 0) return;

    if (isRandom) {
        index = Math.floor(Math.random() * queue.length);
    } else {
        index++;

        if (index >= queue.length) {
            if (loopMode === 'all') {
                index = 0;
            } else {
                index = queue.length - 1;
                return;
            }
        }
    }

    playCurrentSong();
}

function playPrev() {
    if (queue.length === 0) return;

    if (isRandom) {
        index = Math.floor(Math.random() * queue.length);
    } else {
        index--;

        if (index < 0) {
            if (loopMode === 'all') {
                index = queue.length - 1;
            } else {
                index = 0;
            }
        }
    }

    playCurrentSong();
}

export function initPlayerControls() {

    document.getElementById('btn-prev')?.addEventListener('click', playPrev);
    document.getElementById('btn-next')?.addEventListener('click', playNext);

    document.getElementById('btn-random')?.addEventListener('click', (e) => {
        isRandom = !isRandom;

        e.target.innerText = isRandom ? '🔀 (On)' : '🔀 (Off)';
        e.target.title = `Aleatorio (${isRandom ? 'Activado' : 'Desactivado'})`;
        e.target.className = `player-btn ${isRandom ? 'btn-active-random' : 'btn-inactive'}`;
    });

    document.getElementById('btn-loop')?.addEventListener('click', (e) => {
        const modes = ['none', 'all', 'one'];
        const labels = ['🔁 (None)', '🔁 (All)', '🔂 (One)'];
        const classes = ['btn-loop-none', 'btn-loop-all', 'btn-loop-one'];

        let idx = modes.indexOf(loopMode);
        idx = (idx + 1) % modes.length;

        loopMode = modes[idx];

        e.target.innerText = labels[idx];
        e.target.title = `Repetir (${loopMode})`;
        e.target.className = `player-btn ${classes[idx]}`;

        const audioEl = document.getElementById('audio-player');
        audioEl.loop = loopMode === 'one';
    });

    document.getElementById('audio-player')?.addEventListener('ended', () => {
        if (loopMode === 'one') return;
        playNext();
    });
}