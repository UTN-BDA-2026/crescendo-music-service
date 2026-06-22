import { initAuth } from './auth/auth.js';
import { initSearch } from './search/search.js';
import { initPlayerControls } from './player/player.js';
import { checkAuth } from './ui/dom.js';

checkAuth();

initAuth(checkAuth);
initSearch();
initPlayerControls();