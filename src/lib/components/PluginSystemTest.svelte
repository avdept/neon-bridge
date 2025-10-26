<script lang="ts">
  import { onMount } from "svelte";
  import { pluginRegistry } from "../plugins/registry.js";
  import type { Plugin } from "../plugins/types.js";

  let plugins: Plugin[] = [];
  let loading = true;

  onMount(async () => {
    await new Promise((resolve) => setTimeout(resolve, 100));

    plugins = pluginRegistry.getAll();
    loading = false;
  });
</script>

{#if loading}
  <div class="test-status loading">
    <h3>🔄 Checking Plugin Registry...</h3>
  </div>
{:else if plugins.length === 0}
  <div class="test-status error">
    <h3>⚠️ No Plugins Found</h3>
    <p>Plugin system may not be initialized yet</p>
  </div>
{:else}
  <div class="test-status success">
    <h3>✅ Plugin System Active</h3>
    <p>Registry contains {plugins.length} plugins:</p>
    <ul>
      {#each plugins as plugin}
        <li>
          <strong>{plugin.metadata.name}</strong>
          <small>({plugin.metadata.id})</small>
          <span class="category">{plugin.metadata.category}</span>
        </li>
      {/each}
    </ul>
  </div>
{/if}

<style>
  .test-status {
    background: rgba(0, 0, 0, 0.1);
    border-radius: 8px;
    padding: 1rem;
    margin: 1rem 0;
    backdrop-filter: blur(10px);
  }

  .test-status.loading {
    border: 1px solid rgba(255, 165, 0, 0.3);
    color: #ffa500;
  }

  .test-status.error {
    border: 1px solid rgba(248, 113, 113, 0.3);
    color: #f87171;
  }

  .test-status.success {
    border: 1px solid rgba(34, 197, 94, 0.3);
    color: #22c55e;
  }

  .test-status ul {
    margin: 0.5rem 0;
    list-style: none;
    padding: 0;
  }

  .test-status li {
    padding: 0.25rem 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .test-status li small {
    opacity: 0.7;
  }

  .category {
    background: rgba(255, 255, 255, 0.1);
    padding: 0.1rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
    text-transform: uppercase;
    margin-left: auto;
  }

  .test-status h3 {
    margin: 0 0 0.5rem 0;
  }

  .test-status p {
    margin: 0.5rem 0;
  }
</style>
