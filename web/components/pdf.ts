import { LitElement, html } from "lit";
import { customElement, property, query } from "lit/decorators.js";

import { getDocument } from "pdfjs-dist";
import {
  EventBus,
  PDFPageView,
  PDFSinglePageViewer,
} from "pdfjs-dist/web/pdf_viewer";

import Styles from "./pdf.scss";

@customElement("pdf-element")
export class PDF extends LitElement {
  static styles = Styles;

  @property() zoom: number;
  @property() panning: boolean;

  @query(".wrapper") wrapper: HTMLDivElement;
  @query(".viewer-container") viewerContainer: HTMLDivElement;

  constructor() {
    super();
    this.zoom = 100;
  }

  async connectedCallback(): Promise<void> {
    super.connectedCallback();
    await this.updateComplete;

    const doc = await getDocument("static/assets/character_sheet.pdf").promise;
    const firstPage = await doc.getPage(1);

    const eventBus = new EventBus();
    const pageView = new PDFPageView({
      id: 1,
      container: this.viewerContainer,
      eventBus,
      defaultViewport: firstPage.getViewport({ scale: 1.0 }),
      scale: 1,
    });

    pageView.setPdfPage(firstPage);
    pageView.draw();
    // setTimeout(() => {
    //   pageView.update({ scale: 1 });
    //   pageView.draw();
    // }, 1000);
  }

  setZoom(e: WheelEvent) {
    const newZoom = this.zoom - e.deltaY * 0.1;
    if (newZoom > 5 && newZoom < 300) this.zoom = newZoom;
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
        @pointerdown=${() => (this.panning = true)}
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
