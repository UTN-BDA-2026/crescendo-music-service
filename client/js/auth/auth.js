import { API_BASE_URL } from '../config/api.js';
import { checkAuth } from '../ui/dom.js';

export function initAuth(updateAuthUI) {
    const logoutBtn = document.getElementById('logout-btn');

    logoutBtn?.addEventListener('click', () => {
        localStorage.removeItem('crescendo_token');

        document.getElementById('artist-results').innerHTML = '';
        document.getElementById('album-results').innerHTML = '';
        document.getElementById('song-results').innerHTML = '';

        document.getElementById('audio-player').pause();
        document.getElementById('audio-player').src = '';

        document.getElementById('now-playing-title').innerText =
            'Selecciona una canción...';

        updateAuthUI();
    });

    // REGISTER
    document.getElementById('register-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();

        const payload = {
            username: reg('reg-username'),
            email: reg('reg-email'),
            password: reg('reg-password'),
            date_of_birth: reg('reg-dob') + "T00:00:00Z"
        };

        const msg = document.getElementById('reg-msg');

        try {
            const res = await fetch(`${API_BASE_URL}/users`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                msg.textContent = "Usuario creado";
                msg.className = "msg-success";
                e.target.reset();
            } else {
                const err = await res.json();
                msg.textContent = err.error || "Error en el registro";
                msg.className = "msg-error";
            }
        } catch (networkErr) {
            console.error('Register fetch error:', networkErr);
            msg.textContent = "Error de conexión con el servidor";
            msg.className = "msg-error";
        }
    });

    // LOGIN
    document.getElementById('login-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();

        const input = val('log-username');
        const isEmail = input.includes('@');

        const payload = {
            username: isEmail ? "" : input,
            email: isEmail ? input : "",
            password: val('log-password')
        };

        const msg = document.getElementById('log-msg');

        try {
            const res = await fetch(`${API_BASE_URL}/users/login`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                const data = await res.json();
                localStorage.setItem('crescendo_token', data.token);
                msg.textContent = "";
                updateAuthUI();
            } else {
                const err = await res.json();
                msg.textContent = err.error || "Error en el login";
                msg.className = "msg-error";
            }
        } catch (networkErr) {
            console.error('Login fetch error:', networkErr);
            msg.textContent = "Error de conexión con el servidor";
            msg.className = "msg-error";
        }
    });
}

const val = (id) => document.getElementById(id).value.trim();
const reg = val;