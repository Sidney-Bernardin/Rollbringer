<script lang="ts">
  import Icon from "@iconify/svelte";
  import PDFViewer from "./PDFViewer.svelte";

  export let title: string;
  let currentPage: number = 1;
</script>

<div class="dnd-character-sheet" {title}>
  <ul>
    <button class:active={currentPage == 1} on:click={() => (currentPage = 1)}>
      <Icon icon="system-uicons:list" />
      main
    </button>

    <button class:active={currentPage == 2} on:click={() => (currentPage = 2)}>
      <Icon icon="material-symbols:history-edu" />
      info
    </button>

    <button class:active={currentPage == 3} on:click={() => (currentPage = 3)}>
      <Icon icon="game-icons:spell-book" />
      spells
    </button>
  </ul>

  <PDFViewer pdfURL="static/assets/character_sheet.pdf" {currentPage} />
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

      button {
        display: flex;
        cursor: pointer;
        gap: 5px;
        border: none;
        color: $txt4;
        background: transparent;
        font-weight: bolder;
        text-transform: capitalize;
        align-items: center;

        &.active {
          color: #b8bb26;
        }

        iconify-icon {
          font-size: 1.5rem;
        }
      }
    }

    :global(.pdf-viewer) {
      border: 2px solid $shadow;
      flex: 5;
    }
  }
</style>
