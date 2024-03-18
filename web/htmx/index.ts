window.htmx = require("htmx.org");
import "htmx.org/dist/ext/ws.js";

// TODO: Improve error handling.
document.body.addEventListener("htmx:responseError", (e: CustomEvent) => {
    alert(e.detail.xhr.response);
});
