<script lang="ts">
  import type { PluginInstance } from "../plugins/types.js";
  import { pluginRegistry } from "../plugins/registry.js";

  interface Props {
    instance: PluginInstance;
    data?: any;
    onedit?: (event: CustomEvent) => void;
    editMode?: boolean;
  }

  const { instance, data = {}, onedit, editMode = false }: Props = $props();

  // Get plugin and its component dynamically from the registry
  const plugin = $derived(pluginRegistry.get(instance.pluginId));
  const component = $derived(plugin?.component);

  function handleEdit() {
    if (onedit) {
      const editEvent = new CustomEvent("edit", {
        detail: {
          instance,
          plugin,
          config: instance.config,
        },
      });
      onedit(editEvent);
    }
  }
</script>

{#if component && plugin}
  <div
    class="widget-wrapper {editMode ? 'edit-mode' : ''}"
    role="button"
    tabindex="0"
    onclick={editMode ? handleEdit : undefined}
    onkeydown={(e) => {
      if (editMode && (e.key === "Enter" || e.key === " ")) {
        e.preventDefault();
        handleEdit();
      }
    }}
    title={editMode ? "Click to edit widget configuration" : undefined}
  >
    {#if component}
      {@const Component = component}
      <Component
        {plugin}
        config={instance.config}
        {data}
        span={instance.span}
        alert={instance.alert}
      />
    {/if}

    {#if editMode}
      <div class="edit-overlay">
        <button
          class="edit-btn"
          onclick={handleEdit}
          title="Edit Widget"
          aria-label="Edit Widget"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
            ></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
            ></path>
          </svg>
        </button>
      </div>
    {/if}
  </div>
{:else}
  <div class="error-card">
    <h3>Plugin Error</h3>
    <p>Plugin "{instance.pluginId}" not found</p>
  </div>
{/if}

<style>
  .widget-wrapper {
    position: relative;
    width: 100%;
    height: 100%;
    cursor: pointer;
    background: rgba(255, 193, 7, 0.05);
    -webkit-backdrop-filter: blur(20px);
    backdrop-filter: blur(20px);
    transition: all 0.2s ease;
    border-radius: 16px;
  }

  .widget-wrapper:hover {
    transform: translateY(-2px);
    filter: brightness(1.05);
  }

  .widget-wrapper:active {
    transform: translateY(0);
    filter: brightness(0.95);
  }

  .widget-wrapper:focus {
    outline: 0px solid rgba(255, 255, 255, 0.5);
    outline-offset: 2px;
  }

  .widget-wrapper.edit-mode {
    cursor: pointer;
    border: 1px solid transparent;
    transition: all 0.2s ease;
  }

  .widget-wrapper.edit-mode:hover {
    border-color: rgba(245, 158, 11, 0.5);
    transform: translateY(-2px);
    filter: brightness(1.1);
  }

  .edit-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.2s ease;
    border-radius: 16px;
    backdrop-filter: blur(2px);
  }

  .widget-wrapper.edit-mode:hover .edit-overlay {
    opacity: 1;
  }

  .edit-btn {
    background: rgba(245, 158, 11, 0.9);
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .edit-btn:hover {
    background: rgba(245, 158, 11, 1);
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
  }

  .edit-btn svg {
    width: 16px;
    height: 16px;
  }
  .error-card {
    background: rgba(248, 113, 113, 0.1);
    border: 1px solid rgba(248, 113, 113, 0.3);
    border-radius: 12px;
    padding: 1rem;
    color: #f87171;
    text-align: center;
  }

  .error-card h3 {
    margin: 0 0 0.5rem 0;
    font-size: 1rem;
  }

  .error-card p {
    margin: 0;
    font-size: 0.875rem;
    opacity: 0.8;
  }
</style>
