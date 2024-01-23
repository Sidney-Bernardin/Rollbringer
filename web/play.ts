import Alpine from "alpinejs";
import Split from "split-grid";

require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);
import "./styles/layouts/play.scss";

window.alpine = Alpine;
window.htmx = require("htmx.org");

import "./window";
import "./components/pdf-viewer";

document.body.addEventListener("htmx:responseError", (e: CustomEvent) =>
  alert(e.detail.xhr.response),
);

Alpine.start();

Split({
  columnGutters: [
    {
      track: 1,
      element: document.querySelector(".layout__gutter-2"),
    },
    {
      track: 3,
      element: document.querySelector(".layout__gutter-3"),
    },
  ],
  rowGutters: [
    {
      track: 3,
      element: document.querySelector(".layout__gutter-4"),
    },
  ],
});
