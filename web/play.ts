import Split from "split-grid";

require.context("./assets", false, /\.(png|jpg|gif|pdf)$/);
import "./styles/play.scss";

import "./window";
import "./alpine";
import "./htmx";

Split({
    columnGutters: [
        {
            track: 1,
            element: document.querySelector(".play-layout__gutter-2"),
        },
        {
            track: 3,
            element: document.querySelector(".play-layout__gutter-3"),
        },
    ],
    rowGutters: [
        {
            track: 3,
            element: document.querySelector(".play-layout__gutter-4"),
        },
    ],
});
