import * as Panzoom from "panzoom";
import { PDFDocumentProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";

(window as any).Alpine.data("pdfViewer", (pdfURL: string) => ({
  pdfDoc: null,
  pageView: null,

  currentPage: 0,

  async init() {
    const pdfDoc = await getDocument(pdfURL).promise;
    const firstPage = await pdfDoc.getPage(1);

    const pageView = new PDFPageView({
      id: 1,
      container: this.$refs.viewerContainerElem,
      eventBus: new EventBus(),
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });

    Panzoom(this.$refs.viewerContainerElem, {
      bounds: true,
      minZoom: 0.25,
      maxZoom: 2,
      smoothScroll: false,
      filterKey: () => true,
    });

    pageView.setPdfPage(firstPage);
    pageView.draw();
  },
}));
