import { LitElement, TemplateResult, html } from "lit";
import { customElement, state, query } from "lit/decorators.js";
import "iconify-icon";

import Styles from "./character-sheet.scss";

import "./pdf-viewer";
import { PDFViewerElement } from "./pdf-viewer";

@customElement("character-sheet")
export class CharacterSheet extends LitElement {
  static styles = Styles;

  @state() protected currentPage: number = 1;

  @query("pdf-viewer") pdfElem: PDFViewerElement;

  setPage(pageNum: number) {
    this.pdfElem.setPage(pageNum);
    this.currentPage = pageNum;
  }

  protected isCurrentPage(pageNum: number): "active" | "" {
    return this.currentPage == pageNum ? "active" : "";
  }

  render() {
    return html`
      <div class="wrapper">
        <ul>
          <li class=${this.isCurrentPage(1)} @click=${() => this.setPage(1)}>
            <iconify-icon icon="system-uicons:list"></iconify-icon>
            main
          </li>

          <li class=${this.isCurrentPage(2)} @click=${() => this.setPage(2)}>
            <iconify-icon icon="material-symbols:history-edu"></iconify-icon>
            info
          </li>

          <li class=${this.isCurrentPage(3)} @click=${() => this.setPage(3)}>
            <iconify-icon icon="game-icons:spell-book"></iconify-icon>
            spells
          </li>
        </ul>

        <pdf-viewer pdfurl="static/assets/character_sheet.pdf"></pdf-viewer>
      </div>
    `;
  }
}
