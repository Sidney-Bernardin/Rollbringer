<script>
  import Icon from "@iconify/svelte";
  import PDFViewer from "./PDFViewer";
</script>

<div class="dnd-character-sheet" {...$$props}>
  <ul>
    <li>
      <Icon icon="system-uicons:list" />
      main
    </li>

    <li>
      <Icon icon="material-symbols:history-edu" />
      info
    </li>

    <li>
      <Icon icon="game-icons:spell-book" />
      spells
    </li>
  </ul>
</div>

<style lang="scss">
  @import "../global.scss";

  .dnd-character-sheet {
    display: flex;
    height: 100%;
    flex-direction: column;

    ul {
      display: flex;
      gap: 10px;
      width: 100%;
      background: $bg2;
      margin: 0;
      padding: 0 0 5px 0;
      list-style-type: none;
      justify-content: space-around;
      align-self: center;

      li {
        display: flex;
        cursor: pointer;
        gap: 5px;
        color: $txt4;
        font-weight: bolder;
        text-transform: capitalize;
        align-items: center;

        &.active {
          color: $txt1;
        }

        iconify-icon {
          font-size: 1.5rem;
        }
      }
    }

    .pdf-viewer {
      border: 2px solid $shadow;
      flex: 5;
    }
  }
</style>
