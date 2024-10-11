require.context("../assets", false, /\.(png|jpg|gif|pdf)$/);
import "../styles/pages/play/play.scss";

import "../window";
import "../htmx";
import "../alpine/modal";
import "../alpine/pdf-viewer";

import Alpine from "alpinejs";
import Split from "split-grid";


window.alpine = Alpine
Alpine.start()

Split({
    columnGutters: [
        {
            track: 1,
            element: document.querySelector(".gutter-2"),
        },
        {
            track: 3,
            element: document.querySelector(".gutter-3"),
        },
    ],
    rowGutters: [
        {
            track: 3,
            element: document.querySelector(".gutter-4"),
        },
    ],
});
