import "iconify-icon"
import "htmx.org"
import Konva from "konva"
import Alpine from "alpinejs"
import Split from "split-grid"
import htmx from "htmx.org"

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


globalThis.defaultKonvaBoard = `{"attrs":{"width":578,"height":200},"className":"Stage","children":[{"attrs":{},"className":"Layer","children":[{"attrs":{"x":100,"y":100,"sides":6,"radius":70,"fill":"red","stroke":"black","strokeWidth":4},"className":"RegularPolygon"}]}]}`

Alpine.data("board", () => ({
    stage: null,

    init() {
        this.stage = Konva.Node.create(this.$el.dataset.konva, this.$el)
    },
}))

globalThis.Alpine = Alpine
Alpine.start()
