<script lang="ts">
  export let message: string;
  export let statusCode: string;
  export let statusText: string;
</script>

<h1>{message}</h1>
<p>{statusCode} - {statusText}</p>

<style lang="scss">
  button {
    color: #000;
  }
</style>
