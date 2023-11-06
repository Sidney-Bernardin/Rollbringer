import "../components/dynamic-tabs";
import Split from "split-grid";
import Styles from "./game.scss";

document.adoptedStyleSheets.push(Styles.styleSheet);

const grid = document.querySelector(".grid") as HTMLElement;
const computedStyles = grid.computedStyleMap();
grid.style.gridTemplateRows = computedStyles
  .get("grid-template-rows")
  .toString();
grid.style.gridTemplateColumns = computedStyles
  .get("grid-template-columns")
  .toString();

Split({
  columnGutters: [
    {
      track: 1,
      element: document.querySelector("span.left"),
    },
    {
      track: 3,
      element: document.querySelector("span.right"),
    },
  ],
  rowGutters: [
    {
      track: 1,
      element: document.querySelector("span.middle"),
    },
  ],
});
