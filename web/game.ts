import Split from "split-grid";
import { GlobalWorkerOptions } from "pdfjs-dist";

import "./styles/game.scss";
require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);

import "./alpine";

(window as any).htmx = require("htmx.org");
GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

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
