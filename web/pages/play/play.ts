import Alpine from "alpinejs"
import Split from "split-grid"
import "htmx.org"
import "iconify-icon"

/////

document.addEventListener("htmx:beforeSwap", (e: CustomEvent<{
    isError: boolean,
}>) => {
    if (e.detail.isError) alert(e.detail.serverResponse)
})

/////

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
