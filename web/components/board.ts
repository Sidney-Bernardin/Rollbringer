import Konva from "konva"
import Alpine from "alpinejs"
import { defineComponent } from "."

Alpine.data("board", defineComponent(() => ({
    stage: undefined as Konva.Stage | undefined,
    layer: undefined as Konva.Layer | undefined,
    currentLine: undefined as Konva.Line | undefined,

    init() {
        this.stage = Konva.Node.create(this.$el.dataset.konva, this.$el.querySelector(".konva")) as Konva.Stage
        this.stage.on("mousedown touchstart", (e) => this.mouseDown(e))
        this.stage.on("mouseup touchend", (e) => this.mouseUp(e))
        this.stage.on("mousemove touchmove", (e) => this.mouseMove(e))

        this.layer = this.stage.getLayers()[0]
    },

    mouseDown(e: Konva.KonvaPointerEvent): void {
        const pos = this.stage!.getPointerPosition()!
        this.currentLine = new Konva.Line({
            stroke: "#df4b26",
            strokeWidth: 5,
            globalCompositeOperation: "source-over",
            lineCap: "round",
            lineJoin: "round",
            points: [pos.x, pos.y, pos.x, pos.y],
        })
        this.layer!.add(this.currentLine as Konva.Line)

        window.wsSend(JSON.stringify({
            operation: "CREATE_LINE",
            payload: JSON.parse(this.currentLine.toJSON()),
        }))
    },

    mouseUp(e: Konva.KonvaPointerEvent): void {
        this.currentLine = undefined
    },

    mouseMove(e: Konva.KonvaPointerEvent): void {
        if (!this.currentLine) return
        e.evt.preventDefault() // Prevent scrolling on touch devices.

        const pos = this.stage!.getPointerPosition()!
        this.currentLine.points(this.currentLine.points().concat([pos.x, pos.y]))

        window.wsSend(JSON.stringify({
            operation: "UPDATE_LINE",
            payload: {
                points: this.currentLine.points()
            },
        }))
    },
})))
