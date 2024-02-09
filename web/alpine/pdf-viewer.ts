import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";
import panzoom from "panzoom";
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

window.alpine.data("pdfViewer", (pdfURL: string) => ({
    currentPage: 1,
    pdfViewer: null,

    init() {
        const viewerContainerElem: HTMLDivElement = this.$root.querySelector("form");

        panzoom(viewerContainerElem, {
            bounds: true,
            minZoom: 0.25,
            maxZoom: 2,
            smoothScroll: false,
            filterKey: () => true,
        });

        this.pdfViewer = new PDFViewer(pdfURL, viewerContainerElem);
        this.$watch("currentPage", (newVal: number) => this.pdfViewer.renderPage(newVal));
    },
}));

function PDFViewer(pdfURL: string, viewerContainerElem: HTMLDivElement) {
    getDocument(pdfURL).promise.then(async (res) => {
        this.doc = res;
        this.pageView = new PDFPageView({
            id: 1,
            container: viewerContainerElem,
            eventBus: new EventBus(),
            defaultViewport: (await res.getPage(1)).getViewport({ scale: 1.0 }),
            scale: 1,
        });

        this.renderPage(1);
    });

    this.renderPage = async (num: number) => {
        await this.pageView.setPdfPage(await this.doc.getPage(num));
        await this.pageView.draw();
    };
}
