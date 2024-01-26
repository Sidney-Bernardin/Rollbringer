window.htmx = require("htmx.org");
import "/static/ws.js";

document.body.addEventListener("htmx:responseError", (e: CustomEvent) =>
  alert(e.detail.xhr.response),
);
