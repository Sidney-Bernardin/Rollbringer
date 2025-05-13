import { KonvaEventObject } from "konva/lib/Node"
import * as Models from "./models"
import Konva from "konva"


export var board: HTMLDivElement
export var stage: Konva.Stage | undefined
export var layer: Konva.Layer | undefined

export function initStage(): void {
    board = document.querySelector(".boards .board") as HTMLDivElement

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
        shape.setAttr("draggable", true)
        shape.on("dragmove", onShapeDragMove)
        layer.add(shape)
    }

    window.wsSend(JSON.stringify({ operation: "subscribe-to-canvas" }))
}

export function updateShape(shape: Konva.Shape, update: Models.CanvasNode) {
    for (let [k, v] of Object.entries(update)) {
        if (k == "new_name") k = "name"
        if (k == "color") k = "fill"
        shape.setAttr(k, v)
    }
}

function onShapeDragMove(e: Konva.KonvaPointerEvent): void {
    window.wsSend(JSON.stringify({
        operation: "update-canvas-node",
        name: e.target.name(),
        x: e.target.x(),
        y: e.target.y(),
    }))
}
