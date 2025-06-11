import Split from "split-grid"
import "iconify-icon"

import Alpine from "alpinejs"
import { Alpine as AlpineExtention } from "alpinejs"
window.Alpine = Alpine

import { HtmxExtension } from "htmx.org"
window.HTMX = require("htmx.org")
import "htmx-ext-ws";



document.addEventListener("htmx:beforeSwap", (e: CustomEvent<{
    isError: boolean,
}>) => {
    if (e.detail.isError) alert(e.detail.serverResponse)
})

Alpine.start()

Split({
    columnGutters: [
        {
            track: 1,
            element: document.querySelector(".gutter.g2")!,
        },
        {
            track: 3,
            element: document.querySelector(".gutter.g3")!,
        },
    ],
    rowGutters: [
        {
            track: 3,
            element: document.querySelector(".gutter.g4")!,
        },
    ],
})
