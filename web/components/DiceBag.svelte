<script lang="ts">
  import Icon from "@iconify/svelte";

  export let dice: number[] = [20, 12, 10, 8, 6, 4]
</script>

<div class="dice-bag">
  <div class="top">
    <div class="tower">
    </div>

    <button>
      <Icon icon="fa-solid:dice-d20" />
    </button>
  </div>

  <div class="dice">
    {#each dice as d}
      <button>
        <Icon icon={`mdi:dice-d${d}-outline`} />
      </button>
    {/each}
  </div>

  <div class="stats">
    <button>foo bar</button>
    <button>foo bar</button>
    <button>foo bar</button>
    <button>foo bar</button>
    <button>foo bar</button>
  </div>
</div>

<style lang="scss">
  @import "../global.scss";

  .dice-bag {
    display: flex;
    gap: 10px;
    border: 2px solid $shadow;
    background: $bg2;
    padding: 10px;
    flex-direction: column;

    .top {
      display: flex;
      gap: 10px;

      button {
        color: $purple;
        border: none;
        background: transparent;
        font-size: 2.5rem;
        flex: 1;
      }

      .tower {
        width: 100%;
        border: 2px solid $shadow;
        background: $bg1;
        flex: 5;
      }
    }

    .dice {
      display: flex;

      button {
        border: none;
        color: $purple;
        background: transparent;
        font-size: 2rem;
        padding: 0;

        &.active {
          color: $bg2;
          background: $purple;
        }
      }
    }

    .stats {
      display: flex;
      flex-direction: column;

      button {
        border: none;
        color: $blue;
        background: transparent;
        font-weight: bolder;

        &.active {
          color: $bg2;
          background: $blue;
        }
      }
    }
  }
</style>
