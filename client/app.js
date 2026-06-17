const API_BASE_URL = 'http://localhost:8080';

// ESTADO DEL REPRODUCTOR
let playerQueue = [];
let playerCurrentIndex = -1;
let playerLoopMode = 'none'; // 'none', 'all', 'one'
let playerIsRandom = false;

// UTILIDADES
function checkAuth() {
    const token = localStorage.getItem('crescendo_token');
    const authSection = document.getElementById('auth-section');
    const appSection = document.getElementById('app-section');
    const playerSection = document.getElementById('player-section');

    if (token) {
        authSection.classList.add('hidden');
        appSection.classList.remove('hidden');
        playerSection.classList.remove('hidden');
    } else {
        authSection.classList.remove('hidden');
        appSection.classList.add('hidden');
        playerSection.classList.add('hidden');
    }
}

// Inicializar estado
checkAuth();

// LOGOUT
document.getElementById('logout-btn')?.addEventListener('click', () => {
    localStorage.removeItem('crescendo_token');
    document.getElementById('artist-results').innerHTML = '';
    document.getElementById('album-results').innerHTML = '';
    document.getElementById('song-results').innerHTML = '';
    document.getElementById('audio-player').pause();
    document.getElementById('audio-player').src = '';
    document.getElementById('now-playing-title').innerText = 'Selecciona una canción...';
    checkAuth();
});

// REGISTRO
document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = {
        username: document.getElementById('reg-username').value.trim(),
        email: document.getElementById('reg-email').value.trim(),
        password: document.getElementById('reg-password').value,
        date_of_birth: document.getElementById('reg-dob').value + "T00:00:00Z" 
    };

    try {
        const res = await fetch(`${API_BASE_URL}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const msgEl = document.getElementById('reg-msg');
        if (res.ok) {
            msgEl.innerText = "¡Usuario creado! Ya puedes hacer login a la izquierda.";
            msgEl.style.color = "green";
            e.target.reset();
        } else {
            const errData = await res.json();
            msgEl.innerText = `Error: ${errData.error}`;
            msgEl.style.color = "red";
        }
    } catch (error) {
        alert("Error de conexión con la API");
    }
});

// LOGIN
document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const input = document.getElementById('log-username').value.trim();
    const isEmail = input.includes('@');
    
    const payload = {
        username: isEmail ? "" : input,
        email: isEmail ? input : "",
        password: document.getElementById('log-password').value
    };

    try {
        const res = await fetch(`${API_BASE_URL}/users/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const msgEl = document.getElementById('log-msg');
        if (res.ok) {
            const data = await res.json();
            localStorage.setItem('crescendo_token', data.token);
            e.target.reset();
            msgEl.innerText = "";
            checkAuth();
        } else {
            const errData = await res.json();
            msgEl.innerText = `Error: ${errData.error}`;
            msgEl.style.color = "red";
        }
    } catch (error) {
        document.getElementById('log-msg').innerText = "Error de conexión con la API.";
    }
});

// BÚSQUEDA GLOBAL
document.getElementById('global-search-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();

    const query = document.getElementById('search-query').value.trim();
    if (!query) return;

    const token = localStorage.getItem('crescendo_token');

    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };

    try {
        const response = await fetch(
            `${API_BASE_URL}/search?q=${encodeURIComponent(query)}&type=all`,
            { headers }
        );

        if (!response.ok) {
            throw new Error('Error al realizar la búsqueda');
        }

        const data = await response.json();

        const artistData = data.artists || [];
        const albumData = data.albums || [];
        const songData = data.songs || [];

        // ==========================
        // ARTISTAS
        // ==========================
        const artistContainer = document.getElementById('artist-results');
        artistContainer.innerHTML = '';

        if (artistData.length > 0) {
            artistContainer.innerHTML =
                `<h3 style="margin-bottom:10px;color:#111827;">
                    Artistas Encontrados (${artistData.length})
                </h3>`;

            artistData.forEach(artist => {
                const div = document.createElement('div');
                div.className = 'glass-card';

                div.innerHTML = `
                    <div style="display:flex;align-items:center;gap:20px;">
                        ${
                            artist.image_url
                                ? `<img src="${artist.image_url}" 
                                       alt="${artist.name}"
                                       style="width:60px;height:60px;border-radius:50%;object-fit:cover;">`
                                : ''
                        }
                        <div>
                            <h4 style="margin:0;font-size:1.2em;color:#111827;">
                                ${artist.name}
                            </h4>
                            <p style="color:#6b7280;margin:5px 0 0 0;font-size:.9em;">
                                ${artist.information || 'Sin información adicional'}
                            </p>
                        </div>
                    </div>
                `;

                artistContainer.appendChild(div);
            });
        }

        // ==========================
        // ÁLBUMES
        // ==========================
        const albumContainer = document.getElementById('album-results');
        albumContainer.innerHTML = '';

        if (albumData.length > 0) {
            albumContainer.innerHTML =
                `<h3 style="margin-bottom:10px;color:#111827;">
                    Álbumes Encontrados (${albumData.length})
                </h3>`;

            albumData.forEach(album => {
                const div = document.createElement('div');
                div.className = 'glass-card';

                div.innerHTML = `
                    <div style="display:flex;align-items:center;gap:20px;">
                        ${
                            album.cover_image_url
                                ? `<img src="${album.cover_image_url}"
                                       alt="${album.title}"
                                       style="width:60px;height:60px;border-radius:4px;object-fit:cover;">`
                                : ''
                        }
                        <div>
                            <h4 style="margin:0;font-size:1.2em;color:#111827;">
                                ${album.title}
                                <span style="
                                    font-size:.7em;
                                    background:#e5e7eb;
                                    padding:2px 6px;
                                    border-radius:4px;
                                ">
                                    ${album.type}
                                </span>
                            </h4>
                            <p style="color:#6b7280;margin:5px 0 0 0;font-size:.9em;">
                                Lanzamiento:
                                ${new Date(album.release_date).toLocaleDateString()}
                            </p>
                        </div>
                    </div>
                `;

                albumContainer.appendChild(div);
            });
        }

        // ==========================
        // CANCIONES
        // ==========================
        const songContainer = document.getElementById('song-results');
        songContainer.innerHTML = '';

        if (songData.length > 0) {
            songContainer.innerHTML =
                `<h3 style="margin-bottom:10px;color:#111827;">
                    Canciones Encontradas (${songData.length})
                </h3>`;

            const escapedSongs = JSON.stringify(songData)
                .replace(/'/g, "&apos;")
                .replace(/"/g, "&quot;");

            songData.forEach((song, idx) => {
                const div = document.createElement('div');
                div.className = 'song-item glass-card';

                const m = Math.floor(song.duration / 60);
                const s = (song.duration % 60).toString().padStart(2, '0');

                const artistNames = song.artists?.length
                    ? song.artists.map(a => a.name).join(', ')
                    : 'Desconocido';

                div.innerHTML = `
                    <div>
                        <strong style="color:#1f2937;font-size:1.1em;">
                            ${song.title}
                        </strong>

                        <span style="color:#4b5563;font-size:.9em;margin-left:5px;">
                            - ${artistNames}
                        </span>

                        <div style="color:#6b7280;font-size:.85em;margin-top:4px;">
                            Duración: ${m}:${s}
                        </div>
                    </div>

                    <button
                        onclick="playQueue(${escapedSongs}, ${idx})"
                        class="play-btn"
                        style="padding:4px 10px;border-radius:4px;border:none;cursor:pointer;color:white;"
                    >
                        ▶ Play
                    </button>
                `;

                songContainer.appendChild(div);
            });
        }

        // ==========================
        // SIN RESULTADOS
        // ==========================
        if (
            artistData.length === 0 &&
            albumData.length === 0 &&
            songData.length === 0
        ) {
            artistContainer.innerHTML =
                `<p style="color:#6b7280;">
                    No se encontraron resultados para "${query}".
                </p>`;
        }

    } catch (error) {
        alert(error.message || 'Error de conexión al realizar la búsqueda.');
    }
});

// REPRODUCCIÓN Y COLA
function playQueue(songsArray, startIndex = 0) {
    if (!songsArray || songsArray.length === 0) return;
    
    playerQueue = [...songsArray];
    
    if (playerIsRandom) {
        playerCurrentIndex = Math.floor(Math.random() * playerQueue.length);
    } else {
        playerCurrentIndex = startIndex;
    }
    
    playCurrentSong();
}

function playSong(songId, title = '') {
    playQueue([{ id: songId, title: title }], 0);
}

async function playCurrentSong() {
    if (playerCurrentIndex < 0 || playerCurrentIndex >= playerQueue.length) return;
    
    const song = playerQueue[playerCurrentIndex];
    const token = localStorage.getItem('crescendo_token');
    const headers = { 'Authorization': `Bearer ${token}` };

    try {
        const res = await fetch(`${API_BASE_URL}/songs/${song.id}/playback`, { headers });
        if (!res.ok) throw new Error("No se pudo obtener información de reproducción");
        
        const data = await res.json();
        const titleEl = document.getElementById('now-playing-title');
        const audioEl = document.getElementById('audio-player');
        
        const artistNames = (data.artists && data.artists.length > 0) ? data.artists.map(a => a.name).join(', ') : 'Desconocido';
        titleEl.innerText = `${data.title} - ${artistNames}`;
        
        audioEl.src = data.stream_url;
        audioEl.play();
    } catch (error) {
        alert(error.message || "Error al intentar reproducir la canción");
    }
}

function playNext() {
    if (playerQueue.length === 0) return;

    if (playerIsRandom) {
        playerCurrentIndex = Math.floor(Math.random() * playerQueue.length);
    } else {
        playerCurrentIndex++;
        if (playerCurrentIndex >= playerQueue.length) {
            if (playerLoopMode === 'all') {
                playerCurrentIndex = 0;
            } else {
                playerCurrentIndex = playerQueue.length - 1;
                return; // don't play, end of queue
            }
        }
    }
    playCurrentSong();
}

function playPrev() {
    if (playerQueue.length === 0) return;

    if (playerIsRandom) {
        playerCurrentIndex = Math.floor(Math.random() * playerQueue.length);
    } else {
        playerCurrentIndex--;
        if (playerCurrentIndex < 0) {
            if (playerLoopMode === 'all') {
                playerCurrentIndex = playerQueue.length - 1;
            } else {
                playerCurrentIndex = 0;
            }
        }
    }
    playCurrentSong();
}

// CONTROLES DEL REPRODUCTOR
document.getElementById('btn-prev')?.addEventListener('click', playPrev);
document.getElementById('btn-next')?.addEventListener('click', playNext);

document.getElementById('btn-random')?.addEventListener('click', (e) => {
    playerIsRandom = !playerIsRandom;
    e.target.innerText = playerIsRandom ? '🔀 (On)' : '🔀 (Off)';
    e.target.title = `Aleatorio (${playerIsRandom ? 'Activado' : 'Desactivado'})`;
    e.target.style.color = playerIsRandom ? '#10b981' : 'white';
});

document.getElementById('btn-loop')?.addEventListener('click', (e) => {
    const modes = ['none', 'all', 'one'];
    const labels = ['🔁 (None)', '🔁 (All)', '🔂 (One)'];
    const colors = ['white', '#10b981', '#3b82f6'];
    
    let idx = modes.indexOf(playerLoopMode);
    idx = (idx + 1) % modes.length;
    playerLoopMode = modes[idx];
    
    e.target.innerText = labels[idx];
    e.target.title = `Repetir (${playerLoopMode})`;
    e.target.style.color = colors[idx];
    
    const audioEl = document.getElementById('audio-player');
    if (playerLoopMode === 'one') {
        audioEl.loop = true;
    } else {
        audioEl.loop = false;
    }
});

// Evento ended para siguiente canción
document.getElementById('audio-player')?.addEventListener('ended', () => {
    if (playerLoopMode === 'one') return; // audio loops automatically via audioEl.loop = true
    playNext();
});