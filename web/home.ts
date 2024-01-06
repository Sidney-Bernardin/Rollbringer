import Alpine from "alpinejs";
import Split from "split-grid";

import "./alpine";
import "./htmx/htmx";

import "./styles/home.scss";
require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);

Split({
  columnGutters: [
    {
      track: 1,
      element: document.querySelector(".layout__gutter"),
    },
  ],
});
