import Split from "split-grid";

require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);
import "./styles/pages/play/index.scss";

import "./window";
import "./alpine";
import "./htmx";

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
