import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist"
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer"
import panzoom from "panzoom"
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js"

const pdfDocuments = {}
const pdfViews = {}

window.alpine.data("pdfViewer", (pdfID: string, pdfSchema: string) => ({
    currentPage: 0,
    pdfPageView: null,

    async init() {
        const container = this.$el.querySelector(".viewer")

        panzoom(container, {
            bounds: true,
            minZoom: 0.25,
            maxZoom: 2,
            smoothScroll: false,
            filterKey: () => true,
        })

        pdfDocuments[pdfID] = await getDocument(`/static/${pdfSchema}.pdf`).promise
        pdfViews[pdfID] = new PDFPageView({
            id: 1,
            container: container,
            eventBus: new EventBus(),
            defaultViewport: (await pdfDocuments[pdfID].getPage(1)).getViewport({ scale: 1.0 }),
            scale: 1,
        })

        this.$watch("currentDynamicTab", (newVal) => {
            if (newVal === pdfID && this.currentPage !== 0) {
                htmx.trigger(this.$refs.subscribe_form, "submit", null)
            }
        })
    },

    async openPage(pageNum: number) {
        this.currentPage = pageNum

        await pdfViews[pdfID].setPdfPage(await pdfDocuments[pdfID].getPage(pageNum))
        await pdfViews[pdfID].draw()

        htmx.trigger(this.$refs.subscribe_form, "submit", null)
    }
}))
