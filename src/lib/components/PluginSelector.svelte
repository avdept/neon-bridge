<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { pluginRegistry } from "../plugins/registry.js";
  import type { Plugin } from "../plugins/types.js";

  const dispatch = createEventDispatcher();

  export let selectedPluginId: string | null = null;
  export let onSelect: (pluginId: string | null) => void;

  let plugins: Plugin[] = [];
  let iconPath: string;
  $: {
    plugins = pluginRegistry.getAll();
  }

  function selectPlugin(pluginId: string) {
    if (selectedPluginId === pluginId) {
      selectedPluginId = null;
      onSelect(null);
    } else {
      selectedPluginId = pluginId;
      onSelect(pluginId);
    }
  }
</script>

<div class="plugin-selector">
  <h3 class="selector-title">Choose a Plugin</h3>
  <p class="selector-description">Select a plugin to add to your dashboard</p>

  <div class="plugin-grid">
    {#each plugins as plugin (plugin.metadata.id)}
      {@const iconPath = `/services/${plugin.metadata.icon}.svg`}
      <button
        class="plugin-card {selectedPluginId === plugin.metadata.id
          ? 'selected'
          : ''}"
        on:click={() => selectPlugin(plugin.metadata.id)}
      >
        <div
          class="plugin-icon"
          style="background-image: url('{iconPath}')"
        ></div>
        <div class="plugin-info">
          <h4 class="plugin-name">{plugin.metadata.name}</h4>
          <p class="plugin-description">{plugin.metadata.description}</p>
          <span class="plugin-category">{plugin.metadata.category}</span>
        </div>
        {#if selectedPluginId === plugin.metadata.id}
          <div class="selected-indicator">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
              <path
                d="M16.667 5L7.5 14.167 3.333 10"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
        {/if}
      </button>
    {/each}
  </div>
</div>

<style>
  .plugin-selector {
    width: 100%;
  }

  .selector-title {
    color: white;
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 0.5rem 0;
    background: linear-gradient(135deg, #ffffff 0%, #e8e8e8 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .selector-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.875rem;
    margin: 0 0 1.5rem 0;
  }

  .plugin-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.75rem;
    max-height: 400px;
    overflow-y: auto;
    padding-right: 0.5rem;
    padding-top: 0.5rem;
  }

  .plugin-card {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    padding: 1rem;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    display: flex;
    align-items: center;
    gap: 1rem;
    text-align: left;
    width: 100%;
    position: relative;
  }

  .plugin-card:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.2);
    transform: translateY(-1px);
  }

  .plugin-card.selected {
    background: rgba(59, 130, 246, 0.15);
    border-color: rgba(59, 130, 246, 0.4);
    box-shadow: 0 0 0 1px rgba(59, 130, 246, 0.2);
  }

  .plugin-icon {
    font-size: 2rem;
    width: 3rem;
    height: 3rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    flex-shrink: 0;
  }

  .plugin-info {
    flex: 1;
    min-width: 0;
  }

  .plugin-name {
    color: white;
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 0.25rem 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .plugin-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.75rem;
    margin: 0 0 0.5rem 0;
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .plugin-category {
    color: rgba(255, 255, 255, 0.5);
    font-size: 0.625rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-weight: 500;
    background: rgba(255, 255, 255, 0.1);
    padding: 0.125rem 0.5rem;
    border-radius: 4px;
    display: inline-block;
  }

  .selected-indicator {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    color: #3b82f6;
    background: rgba(59, 130, 246, 0.2);
    border-radius: 50%;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Custom scrollbar */
  .plugin-grid::-webkit-scrollbar {
    width: 6px;
  }

  .plugin-grid::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 3px;
  }

  .plugin-grid::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }

  .plugin-grid::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  @media (max-width: 768px) {
    .plugin-card {
      padding: 0.75rem;
    }

    .plugin-icon {
      font-size: 1.5rem;
      width: 2.5rem;
      height: 2.5rem;
    }

    .plugin-name {
      font-size: 0.875rem;
    }

    .plugin-description {
      font-size: 0.7rem;
    }
  }
</style>
