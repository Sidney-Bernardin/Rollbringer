import { GlobalWorkerOptions, PDFDocumentProxy, PDFPageProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";
import panzoom from "panzoom";
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

window.alpine.data("pdfViewer", (pdfURL: string, pdfID: string) => ({
    pageViewer: null,

    init() {
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

        Array.from(this.$el.querySelectorAll(".annotationLayer input")).forEach(
            (elem: HTMLInputElement) => {
                elem.setAttribute("ws-send", "");
                elem.setAttribute("hx-trigger", "change");
                elem.setAttribute(
                    "hx-include",
                    `.dynamic-tab-container [data-pdf-id="${pdfID}"] input.update-pdf-field-param`,
                );

                htmx.process(elem);
            },
        );

        htmx.ajax("GET", `/play-materials/pdfs/${pdfID}/${e.detail}`, { swap: "none" });
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
