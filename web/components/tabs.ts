import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import "iconify-icon";
import Styles from "./tabs.scss";

@customElement("tabs-element")
export class Tabs extends LitElement {
  static styles = Styles;

  @property() tabTitles: string[];
  @property() currentTab: string;

  constructor() {
    super();

    this.tabTitles = [...this.querySelectorAll("div[slot]")].map(
      (tabElem) => tabElem.attributes.getNamedItem("slot").value,
    );
    this.currentTab = this.tabTitles[0];
  }

  removeTab(tabTitle: string) {
    this.tabTitles.splice(this.tabTitles.indexOf(tabTitle), 1);
    this.requestUpdate("tabTitles");
    this.currentTab = this.tabTitles[0];
    this.removeChild(this.querySelector(`div[slot="${tabTitle}"]`));
  }

  render() {
    return html`
      <div class="wrapper">
        <ul>
          ${this.tabTitles.map(
            (tt) => html`
              <li
                class=${this.currentTab == tt ? "active" : ""}
                @click=${() => (this.currentTab = tt)}
              >
                ${tt}
                <button @click=${() => this.removeTab(tt)}>
                  <iconify-icon icon="material-symbols:close"></iconify-icon>
                </button>
              </li>
            `,
          )}
        </ul>

        <div class="tab-border">
          <slot name=${this.currentTab}></slot>
        </div>
      </div>
    `;
  }
}
