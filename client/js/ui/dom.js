export function checkAuth() {
    const token = localStorage.getItem('crescendo_token');

    document.getElementById('auth-section').classList.toggle('hidden', !!token);
    document.getElementById('app-section').classList.toggle('hidden', !token);
    document.getElementById('player-section').classList.toggle('hidden', !token);
}