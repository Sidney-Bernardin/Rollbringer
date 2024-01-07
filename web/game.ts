import Alpine from "alpinejs";
import Split from "split-grid";

import "./alpine";
import "./htmx/htmx";

import "./styles/layouts/game.scss";
require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);

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
