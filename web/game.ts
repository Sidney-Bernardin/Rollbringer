import Split from "split-grid";

import "./alpine";
import "./styles/game.scss";

(window as any).htmx = require("htmx.org");

Split({
  columnGutters: [
    {
      track: 1,
      element: document.querySelector(".layout__left-gutter"),
    },
    {
      track: 3,
      element: document.querySelector(".layout__right-gutter"),
    },
  ],
  rowGutters: [
    {
      track: 1,
      element: document.querySelector(".layout__middle-gutter"),
    },
  ],
});
