import { LitElement, html } from "lit";
import { customElement, property, state, query } from "lit/decorators.js";

import { PDFDocumentProxy, getDocument } from "pdfjs-dist";
import { EventBus, PDFPageView } from "pdfjs-dist/web/pdf_viewer";

import Styles from "./pdf-viewer.scss";

@customElement("pdf-viewer")
export class PDFViewerElement extends LitElement {
  static styles = Styles;

  @property() pdfURL: string;

  @state() protected zoom: number = 100;

  @query(".wrapper") wrapper: HTMLDivElement;
  @query(".viewer-container") viewerContainer: HTMLDivElement;

  pdfDoc: PDFDocumentProxy;
  pageView: PDFPageView;
  panning: boolean;

  async connectedCallback(): Promise<void> {
    super.connectedCallback();
    await this.updateComplete;
    this.setPageView();
  }

  async setPage(pageNum: number) {
    this.pageView.setPdfPage(await this.pdfDoc.getPage(pageNum));
    this.pageView.draw();
  }

  protected async setPageView() {
    this.pdfDoc = await getDocument(this.pdfURL).promise;
    const firstPage = await this.pdfDoc.getPage(1);

    this.pageView = new PDFPageView({
      id: 1,
      container: this.viewerContainer,
      eventBus: new EventBus(),
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });

    await this.setPage(1);
  }

  protected setZoom(e: WheelEvent) {
    e.preventDefault();
    const newZoom = this.zoom - e.deltaY * 0.1;
    if (newZoom > 5 && newZoom < 300) this.zoom = newZoom;
  }

  protected startPanning(e: PointerEvent) {
    this.panning =
      e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
        ? false
        : true;
  }

  protected pan(e: PointerEvent) {
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
        @pointermove=${this.pan}
        @pointerdown=${this.startPanning}
        @pointerup=${() => (this.panning = false)}
        @pointerleave=${() => (this.panning = false)}
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
