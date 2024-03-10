import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";
import panzoom from "panzoom";
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

window.alpine.data("pdfViewer", (pdfURL: string, pdfID: string) => ({
    pageViewer: null,
    initForm: null,

    init() {
        this.initForm = this.$el.querySelector(".pdf-viewer__init-form");
        const viewerContainer: HTMLDivElement = this.$el.querySelector(
            ".pdf-viewer__viewer-container",
        );

        panzoom(viewerContainer, {
            bounds: true,
            minZoom: 0.25,
            maxZoom: 2,
            smoothScroll: false,
            filterKey: () => true,
        });

        this.pageViewer = new PDFPageViewer(pdfURL, viewerContainer);
    },

    async changePage(e: CustomEvent<number>) {
        if ((e.target as HTMLElement).dataset.pdfId !== pdfID) return;

        await this.pageViewer.renderPage(e.detail);
        // htmx.trigger(this.initForm, "submit", null);
    },
}));

function PDFPageViewer(pdfURL: string, viewerContainer: HTMLDivElement) {
    getDocument(pdfURL).promise.then(async (res) => {
        this.doc = res;
        this.pageView = new PDFPageView({
            id: 1,
            container: viewerContainer,
            eventBus: new EventBus(),
            defaultViewport: (await res.getPage(1)).getViewport({ scale: 1.0 }),
            scale: 1,
        });
    });

    this.renderPage = async (num: number) => {
        await this.pageView.setPdfPage(await this.doc.getPage(num));
        await this.pageView.draw();
    };
}
