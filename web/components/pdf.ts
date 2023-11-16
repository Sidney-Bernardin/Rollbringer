import { LitElement, html } from "lit";
import { customElement, property, query } from "lit/decorators.js";

import { PDFDocumentProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";

import Styles from "./pdf.scss";

@customElement("pdf-element")
export class PDF extends LitElement {
  static styles = Styles;

  @property() pdfURL: string;
  @property() zoom: number;
  @property() panning: boolean;

  @query(".wrapper") wrapper: HTMLDivElement;
  @query(".viewer-container") viewerContainer: HTMLDivElement;

  pdfDoc: PDFDocumentProxy;
  pageView: PDFPageView;

  constructor() {
    super();
    this.zoom = 100;
  }

  async connectedCallback(): Promise<void> {
    super.connectedCallback();
    await this.updateComplete;

    this.pdfDoc = await getDocument(this.pdfURL).promise;
    const firstPage = await this.pdfDoc.getPage(1);

    const eventBus = new EventBus();
    this.pageView = new PDFPageView({
      id: 1,
      container: this.viewerContainer,
      eventBus,
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });

    this.pageView.setPdfPage(firstPage);
    this.pageView.draw();
  }

  async setPage(pageNum: number) {
    this.pageView.setPdfPage(await this.pdfDoc.getPage(pageNum));
    this.pageView.draw();
  }

  setZoom(e: WheelEvent) {
    e.preventDefault();
    const newZoom = this.zoom - e.deltaY * 0.1;
    if (newZoom > 5 && newZoom < 300) this.zoom = newZoom;
  }

  startPanning(e: PointerEvent) {
    this.panning =
      e.target instanceof HTMLInputElement ||
      e.target instanceof HTMLTextAreaElement
        ? false
        : true;
  }

  pan(e: PointerEvent) {
    if (!this.panning) return;
    this.wrapper.scrollTop -= e.movementY;
    this.wrapper.scrollLeft -= e.movementX;
  }

  render() {
    return html`
      <div
        class="wrapper"
        style="zoom: ${this.zoom}%"
        @wheel=${this.setZoom}
        @pointerdown=${this.startPanning}
        @pointerup=${() => (this.panning = false)}
        @pointerleave=${() => (this.panning = false)}
        @pointermove=${this.pan}
      >
        <div class="content">
          <div class="viewer-container">
            <div class="viewer"></div>
          </div>
        </div>
      </div>
    `;
  }
}
