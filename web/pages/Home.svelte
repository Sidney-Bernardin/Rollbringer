<script lang="ts">
  import { onMount } from "svelte";
  import htmx from "htmx.org";

  onMount(() => htmx.process(document.querySelector("body")));
</script>

<a href="/users/oauth-login"> LOGIN WITH GOOGLE </a>

<style lang="scss">
  button {
    color: #000;
  }
</style>
