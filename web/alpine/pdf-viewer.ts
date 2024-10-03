import Alpine from "alpinejs";
import panzoom from "panzoom"
import htmx from "htmx.org";

import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist"
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer"


GlobalWorkerOptions.workerSrc = "static/pdf.worker.js"

const pdfDocuments = {}
const pdfViews = {}

Alpine.data("pdfViewer", (pdfID: string, pdfSchema: string) => ({
    currentPage: 0,

    async init(): Promise<void> {

        // Initialize Panzoom for the viewer container.
        panzoom(this.$refs.viewer, {
            bounds: true,
            minZoom: 0.25,
            maxZoom: 2,
            smoothScroll: false,
            filterKey: () => true,
        })

        // Create a PDF view in the viewer container.
        pdfDocuments[pdfID] = await getDocument(`/static/${pdfSchema}.pdf`).promise
        pdfViews[pdfID] = new PDFPageView({
            id: 1,
            container: this.$refs.viewer,
            eventBus: new EventBus(),
            defaultViewport: (await pdfDocuments[pdfID].getPage(1)).getViewport({ scale: 1.0 }),
            scale: 1,
        })

        // Subscribes when currentDynamicTab is changed.
        this.$watch("currentDynamicTab", (newVal) => {
            if (newVal === pdfID && this.currentPage !== 0) {
                this.subscribe()
            }
        })
    },

    // Submits the subscribe form.
    subscribe(): void {
        htmx.trigger(this.$refs.subscribe_form, "submit", null)
    },

    // Renders a new PDF page and prepares it's fields for live updating.
    async openPage(pageNum: number): Promise<void> {
        this.currentPage = pageNum

        // Render the new PDF page.
        await pdfViews[pdfID].setPdfPage(await pdfDocuments[pdfID].getPage(pageNum))
        await pdfViews[pdfID].draw()

        // Prepare each field for live updating.
        Array
            .from(this.$root.querySelectorAll(".annotationLayer input, .annotationLayer textarea"))
            .forEach(this.prepareField)

        this.subscribe()
    },

    // Prepares a field for live updating.
    prepareField(field: HTMLInputElement, hash: number): void {
        // <temporary
        const prefix: string = field.tagName === "TEXTAREA" ? "textarea" : field.type;
        field.name = `${prefix}__${field.name}`;
        // temporary>

        field.id = field.name.replace(/\s/g, "");
        field.name = "field_value_" + hash;
        field.removeAttribute("style");
        field.setAttribute("ws-send", "");
        field.setAttribute("hx-trigger", "change");

        htmx.process(field);
    },
}))
