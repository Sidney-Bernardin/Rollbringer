import * as Models from "./models"
import Konva from "konva"


export var board: HTMLDivElement
export var stage: Konva.Stage | undefined
export var layer: Konva.Layer | undefined

export function initStage(): void {
    board = document.querySelector(".board .board") as HTMLDivElement

    stage = new Konva.Stage({
        container: board.querySelector(".stage-container") as HTMLDivElement,
        width: 500,
        height: 500,
    })

    const canvas: Models.Canvas = JSON.parse(board.dataset.canvas!)

    layer = new Konva.Layer({ name: "layer" })!
    stage.add(layer)

    for (const node of canvas.nodes) {
        let shape: Konva.Shape
        switch (node.type) {
            case "rect": shape = new Konva.Rect({})
            case "circle": shape = new Konva.Circle({})
        }
        updateShape(shape, node)
        shape.on("dragmove", onKonvaDragMove)
        layer.add(shape)
    }

    window.wsSend(JSON.stringify({
        operation: "board-subscribe",
        payload: {
            board_id: board.dataset.boardId,
        }
    }))
}

export function updateShape(shape: Konva.Shape, update: Models.CanvasNode) {
    for (let [k, v] of Object.entries(update)) {
        if (k == "color") k = "fill"
        shape.setAttr(k, v)
    }
}

function onKonvaDragMove(e: Konva.KonvaPointerEvent): void {
    window.wsSend(JSON.stringify({
        operation: "canvas-node-update",
        payload: {
            name: e.target.name(),
            x: e.target.x(),
            y: e.target.y(),
        },
    }))
}
