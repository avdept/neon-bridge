<script lang="ts">
  interface Props {
    columns?: number;
    gap?: string;
    children?: import("svelte").Snippet;
  }

  const { columns = 2, gap = "1rem", children }: Props = $props();
</script>

<div
  class="stats-grid"
  style="grid-template-columns: repeat({columns}, 1fr); gap: {gap}"
>
  {@render children?.()}
</div>

<style>
  .stats-grid {
    display: grid;
  }
</style>
