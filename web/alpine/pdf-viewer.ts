import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";
import panzoom from "panzoom";
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

window.alpine.data("pdfViewer", (pdfURL: string, pdfID: string) => ({
    currentPage: 1,
    pdfViewer: null,

    init() {
        const initFormElem: HTMLFormElement = this.$root.querySelector("form.init-pdf-page");
        const viewerContainerElem: HTMLDivElement =
            this.$root.querySelector("form.viewer-container");

        panzoom(viewerContainerElem, {
            bounds: true,
            minZoom: 0.25,
            maxZoom: 2,
            smoothScroll: false,
            filterKey: () => true,
        });

        this.pdfViewer = new PDFViewer(pdfURL, viewerContainerElem, async () => {
            await this.pdfViewer.renderPage(1);
            htmx.trigger(initFormElem, "submit", null);
        });

        this.$watch("currentPage", async (newVal: number) => {
            await this.pdfViewer.renderPage(newVal);
            htmx.trigger(initFormElem, "submit", null);
        });

        this.$watch("currentDynamicTab", (newVal: string) => {
            if (newVal === pdfID) htmx.trigger(initFormElem, "submit", null);
        });
    },
}));

function PDFViewer(pdfURL: string, viewerContainerElem: HTMLDivElement, cb?: () => void) {
    getDocument(pdfURL).promise.then(async (res) => {
        this.doc = res;
        this.pageView = new PDFPageView({
            id: 1,
            container: viewerContainerElem,
            eventBus: new EventBus(),
            defaultViewport: (await res.getPage(1)).getViewport({ scale: 1.0 }),
            scale: 1,
        });

        cb();
    });

    this.renderPage = async (num: number) => {
        await this.pageView.setPdfPage(await this.doc.getPage(num));
        await this.pageView.draw();
    };
}
