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

  @query(".viewer-container")
  viewerContainer: HTMLDivElement;

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

  render() {
    return html`
      <div class="wrapper">
        <div class="viewer-container">
          <div class="viewer"></div>
        </div>
      </div>
    `;
  }
}
