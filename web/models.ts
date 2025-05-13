export type WebSocketResponse = {
    operation: "error",
    payload: Error,
} | {
    operation: "update-canvas-node"
    payload: CanvasNode
}

export type Error = {
    type: string,
    msg: string,
}

export type Canvas = {
    nodes: CanvasNode[],
}

export type CanvasNode = {
    name: string,

    x: number,
    y: number,

    width: number,
    height: number,

    color: string,
    border_color: string,

    image: string,
} & ({
    type: "rect",
} | {
    type: "circle",
    radius: number,
})
