import "htmx.org"
import "iconify-icon"
import Split from "split-grid"

Split({
    columnGutters: [
        {
            track: 1,
            element: document.querySelector(".gutter.g2"),
        },
        {
            track: 3,
            element: document.querySelector(".gutter.g3"),
        },
    ],
    rowGutters: [
        {
            track: 3,
            element: document.querySelector(".gutter.g4"),
        },
    ],
})
