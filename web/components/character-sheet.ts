import { LitElement, html } from "lit";
import { customElement } from "lit/decorators.js";

import "./pdf";
import Styles from "./character-sheet.scss";

@customElement("character-sheet")
export class CharacterSheet extends LitElement {
  static styles = Styles;

  render() {
    return html`
      <div class="wrapper">
        <ul>
          <li>main</li>
          <li>info</li>
          <li>spells</li>
        </ul>

        <pdf-element></pdf-element>
      </div>
    `;
  }
}
