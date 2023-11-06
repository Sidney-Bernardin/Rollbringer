import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import "iconify-icon";
import Styles from "./dynamic-tab-buttons.scss";

@customElement("dynamic-tab-buttons")
export class DynamicTabButtons extends LitElement {
  static styles = Styles;

  @property({ type: Array }) tabTitles: string[];
  @property({ type: String }) currentTab: string;

  dispatchTabBtnEvent(type: "open" | "close", tabTitle: string) {
    return (e: Event) => {
      e.stopPropagation();
      this.dispatchEvent(new CustomEvent(type, { detail: tabTitle }));
    };
  }

  render() {
    return html`
      <ul>
        ${this.tabTitles.map(
          (tt) => html`
            <li
              class=${this.currentTab == tt ? "active" : ""}
              @click=${this.dispatchTabBtnEvent("open", tt)}
            >
              ${tt}
              <button @click=${this.dispatchTabBtnEvent("close", tt)}>
                <iconify-icon icon="material-symbols:close"></iconify-icon>
              </button>
            </li>
          `,
        )}
      </ul>
    `;
  }
}
