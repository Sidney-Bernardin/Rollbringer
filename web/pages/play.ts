import "iconify-icon"
import Split from "split-grid"

import Alpine from "alpinejs"
import { Alpine as AlpineExtention } from "alpinejs"
import "../components/board.ts"
window.Alpine = Alpine

import { HtmxExtension } from "htmx.org"
window.htmx = require("htmx.org")
import "htmx-ext-ws";


declare global {
    interface Window {
        defaultKonvaBoard: string,
        htmx: HtmxExtension,
        Alpine: AlpineExtention,
        wsSend: (msg: any, fromElem?: HTMLElement) => void,
    }
}


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


window.defaultKonvaBoard = `{"attrs":{"width":578,"height":200},"className":"Stage","children":[{"attrs":{},"className":"Layer","children":[{"attrs":{"x":100,"y":100,"sides":6,"radius":70,"fill":"red","stroke":"black","strokeWidth":4},"className":"RegularPolygon"}]}]}`
window.addEventListener("htmx:wsOpen", (e: CustomEventInit) => window.wsSend = e.detail.socketWrapper.send)

Alpine.start()
