export type WebSocketResponse = {
    operation: "canvas-node-update",
    payload: CanvasNode,
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
