import { LitElement, html } from "lit";
import { customElement, property, query } from "lit/decorators.js";
import "iconify-icon";

import Styles from "./character-sheet.scss";
import "./pdf";
import { PDF } from "./pdf";

@customElement("character-sheet")
export class CharacterSheet extends LitElement {
  static styles = Styles;

  @property() currentPage: number = 1;

  @query("pdf-element") pdfElem: PDF;

  setPage(pageNum: number) {
    this.pdfElem.setPage(pageNum);
    this.currentPage = pageNum;
  }

  render() {
    return html`
      <div class="wrapper">
        <ul>
          <li
            class=${this.currentPage == 1 ? "active" : ""}
            @click=${() => this.setPage(1)}
          >
            <iconify-icon icon="system-uicons:list"></iconify-icon>
            main
          </li>

          <li
            class=${this.currentPage == 2 ? "active" : ""}
            @click=${() => this.setPage(2)}
          >
            <iconify-icon icon="material-symbols:history-edu"></iconify-icon>
            info
          </li>

          <li
            class=${this.currentPage == 3 ? "active" : ""}
            @click=${() => this.setPage(3)}
          >
            <iconify-icon icon="game-icons:spell-book"></iconify-icon>
            spells
          </li>
        </ul>

        <pdf-element pdfurl="static/assets/character_sheet.pdf"></pdf-element>
      </div>
    `;
  }
}
