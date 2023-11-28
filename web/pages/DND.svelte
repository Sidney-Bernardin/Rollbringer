<script lang="ts">
  import { onMount } from "svelte";
  import Split from "split-grid";
  import Tabs from "../components/Tabs.svelte";
  import DNDCharacterSheet from "../components/DNDCharacterSheet.svelte";

  onMount(() => {
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
  });
</script>

<div class="grid">
  <Tabs class="left">
    <DNDCharacterSheet title="hoid" />
    <DNDCharacterSheet title="lee" />
  </Tabs>
  <Tabs class="right" />
  <Tabs class="top" />
  <Tabs class="bottom" />

  <span class="left" />
  <span class="middle" />
  <span class="right" />
</div>

<style lang="scss">
  .grid {
    display: grid;
    box-sizing: border-box;
    width: 100%;
    height: 100%;
    grid-template-rows: 2fr 10px 1fr;
    grid-template-columns: 3fr 10px 2fr 10px 1fr;
    grid-template-areas:
      "l gl t  gr r"
      "l gl gm gr r"
      "l gl b  gr r";
    padding: 10px;

    :global(.tabs.left) {
      grid-area: l;
    }

    :global(.tabs.right) {
      grid-area: r;
    }

    :global(.tabs.top) {
      grid-area: t;
    }

    :global(.tabs.bottom) {
      grid-area: b;
    }

    span {
      &.left:hover,
      &.right:hover {
        cursor: col-resize;
      }

      &.middle:hover {
        cursor: row-resize;
      }

      &.left {
        grid-area: gl;
      }

      &.right {
        grid-area: gr;
      }

      &.middle {
        grid-area: gm;
      }
    }
  }
</style>
