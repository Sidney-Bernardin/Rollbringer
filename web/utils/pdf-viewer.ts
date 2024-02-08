import {
  GlobalWorkerOptions,
  PDFDocumentProxy,
  PDFPageProxy,
  getDocument,
} from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";
import panzoom from "panzoom";
import htmx from "htmx.org";

GlobalWorkerOptions.workerSrc = "static/pdf.worker.js";

export default function(pdfURL: string, viewerContainer: HTMLDivElement) {
  panzoom(viewerContainer, {
    bounds: true,
    minZoom: 0.25,
    maxZoom: 2,
    smoothScroll: false,
    filterKey: () => true,
  });

  getDocument(pdfURL).promise.then(async (res) => {
    this.doc = res;
    this.pageView = new PDFPageView({
      id: 1,
      container: viewerContainer,
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
