import "iconify-icon"
import Split from "split-grid"

import Alpine from "alpinejs"
import { Alpine as AlpineExtention } from "alpinejs"
window.Alpine = Alpine
//import "../components/board.ts"

import { HtmxExtension } from "htmx.org"
window.htmx = require("htmx.org")
import "htmx-ext-ws";



declare global {
    interface Window {
        htmx: HtmxExtension,
        Alpine: AlpineExtention,
        wsSend: (msg: any, fromElem?: HTMLElement) => void,
    }
}


window.addEventListener("htmx:wsOpen", (e: CustomEventInit) => {
    window.wsSend = e.detail.socketWrapper.send
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
