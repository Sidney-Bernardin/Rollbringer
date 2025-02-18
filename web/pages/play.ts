import "iconify-icon"
import "htmx.org"
import Konva from "konva"
import Alpine from "alpinejs"
import Split from "split-grid"

globalThis.Alpine = Alpine
Alpine.start()

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

const stage = new Konva.Stage({
    container: document.querySelector(".board"),
    width: 500,
    height: 500,
})

const layer = new Konva.Layer()

const img = new Image()
const image = new Konva.Image({
    id: "123abc",
    name: "foobarbaz",
    x: 250,
    y: 250,
    image: img,
    width: 100,
    height: 100,
    draggable: true,
})
img.src = "/static/favicon.png"
layer.add(image)

stage.add(layer)
layer.draw()

const json = stage.toJSON()
console.log(json)
