declare global {
    interface Window {
        HTMX: HtmxExtension,
        Alpine: AlpineExtention,
        Board: typeof Board,
        wsSend: (msg: any, fromElem?: HTMLElement) => void,
    }
}



import * as Models from "../models"
import * as Board from "../board"
window.Board = Board

import "iconify-icon"
import Split from "split-grid"
import Konva from "konva"

import Alpine from "alpinejs"
import { Alpine as AlpineExtention } from "alpinejs"
window.Alpine = Alpine


import { HtmxExtension } from "htmx.org"
window.HTMX = require("htmx.org")
import "htmx-ext-ws";



window.addEventListener("htmx:wsOpen", (e: CustomEventInit) => {
    window.wsSend = e.detail.socketWrapper.send
})

window.addEventListener("htmx:wsBeforeMessage", (e: CustomEventInit) => {
    try {
        var msg: Models.WebSocketResponse = JSON.parse(e.detail.message)
    } catch { return }

    switch (msg.operation) {
        case "update-canvas-node":
            const shape: Konva.Shape = Board.layer?.findOne(`.${msg.payload.name}`)!
            Board.updateShape(shape, msg.payload)
            break
    }
})

window.defaultBoardCanvas = {
    nodes: [
        { type: "circle", name: "foo", radius: 50, color: "red" }
    ]
} as Models.Canvas

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
