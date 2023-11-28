<script>
  import Icon from "@iconify/svelte";
</script>

<label>
  <input type="checkbox" />
  <span>
    <Icon icon="mdi:dice-d20-outline" />
  </span>
</label>

<style lang="scss">
  @import "../global.scss";

  label {
    input {
      opacity: 0;
      width: 0;
      height: 0;
    }
  }
</style>
