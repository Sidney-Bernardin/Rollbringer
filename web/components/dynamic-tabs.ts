import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import "./dynamic-tab-buttons";
import Styles from "./dynamic-tabs.scss";

@customElement("dynamic-tabs")
export class DynamicTabs extends LitElement {
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

  updated(x) {
    super.updated(x)
    console.log(x)
  }

  setCurrentTab(e: CustomEvent<string>) {
    this.currentTab = e.detail;
  }

  removeTab(e: CustomEvent<string>) {
    this.tabTitles.splice(this.tabTitles.indexOf(e.detail), 1);
    this.requestUpdate("tabTitles")
    this.currentTab = this.tabTitles[0];
    this.removeChild(this.querySelector(`div[slot="${e.detail}"]`));
  }

  render() {
    return html`
      <div class="wrapper">
        <dynamic-tab-buttons
          tabtitles=${JSON.stringify(this.tabTitles)}
          currenttab=${this.currentTab}
          @open=${this.setCurrentTab}
          @close=${this.removeTab}
        >
        </dynamic-tab-buttons>

        <div class="tab-border">
          <slot name=${this.currentTab}></slot>
        </div>
      </div>
    `;
  }
}
