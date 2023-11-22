<script lang="ts">
  import { onMount } from "svelte";
  import Icon from "@iconify/svelte";

  let clazz: string;
  export { clazz as class };

  let wrapper: HTMLDivElement;
  let tabTitles: string[] = [];
  let currentTab: string;

  // Whenever currentTab changes, displays it's corresponding element with
  // the .active class.
  $: currentTab,
    wrapper?.querySelectorAll(`[title]`).forEach((tabElem) => {
      const force = tabElem.attributes["title"].value == currentTab;
      tabElem.classList.toggle("active", force);
    });

  onMount(() => {
    // Create a MutationObserver that resets tab data whenever a tab is
    // added/removed from wrapper.
    new MutationObserver(() => {
      tabTitles = [...wrapper.querySelectorAll("[title]")].map(
        (tt) => tt.attributes["title"].value
      );
      currentTab = tabTitles[0];
    }).observe(wrapper, { childList: true });

    // Trigger the observer.
    wrapper.appendChild(document.createElement("foo"));
    wrapper.removeChild(wrapper.lastChild);
  });

  const removeTab = (tabTitle: string): void =>
    wrapper.querySelector(`[title=${tabTitle}]`).remove();
</script>

<div class={`tabs ${clazz}`}>
  <ul>
    {#each tabTitles as tt}
      <button
        class:active={currentTab == tt}
        on:click={() => (currentTab = tt)}
      >
        {tt}
        <button on:click={() => removeTab(tt)}>
          <Icon icon="material-symbols:close" />
        </button>
      </button>
    {/each}
  </ul>

  <div class="wrapper" bind:this={wrapper}>
    <slot />
  </div>
</div>

<style lang="scss">
  @import "../global.scss";

  .tabs {
    display: flex;
    flex-direction: column;
    height: 100%;

    .wrapper {
      border: 5px solid $bg2;
      outline: 2px solid $shadow;
      flex: 5;

      :global([title]) {
        display: none;
      }

      :global(.active[title]) {
        display: block;
      }
    }

    ul {
      z-index: 1;
      display: flex;
      gap: 2px;
      margin: 0;
      padding: 0;
      list-style-type: none;
      translate: -2px 0;

      button {
        display: flex;
        cursor: pointer;
        gap: 5px;
        border: 2px solid $bg1;
        border-bottom: none;
        opacity: 0.5;
        background: $bg2;
        padding: 5px;
        text-transform: capitalize;
        align-items: center;

        &.active {
          border: 2px solid $shadow;
          border-bottom: none;
          opacity: 1;
        }

        button {
          cursor: pointer;
          border: none;
          border-radius: 50%;
          background: none;
          opacity: 1;
          padding: 3px;

          &:hover {
            background: $bg3;
          }
        }
      }
    }
  }
</style>
