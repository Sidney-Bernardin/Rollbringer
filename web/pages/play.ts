import "../main";

import "htmx.org/dist/ext/ws.js";
import Split from "split-grid";


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
