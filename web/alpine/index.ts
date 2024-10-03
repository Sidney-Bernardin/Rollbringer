import Alpine from "alpinejs";
window.alpine = Alpine;

import "./pdf-viewer";
import "./modal";

Alpine.store('game', {
    id: document.body.dataset.gameId,
})

Alpine.start();
