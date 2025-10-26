<script lang="ts">
  import { services } from "../stores/system.js";
  import {
    pluginInstancesFromDB,
    isDashboardLoading,
    dashboardError,
  } from "../stores/dashboard.js";
  import PluginRenderer from "./PluginRenderer.svelte";
  import type { PluginInstance } from "../plugins/types.js";
  import type { PluginAlert } from "../plugins/types.js";

  interface Props {
    onedit?: (event: CustomEvent) => void;
    editMode?: boolean;
    onaddwidget?: (position?: number) => void;
  }

  const { onedit, editMode = false, onaddwidget }: Props = $props();

  let showTooltip = $state(false);
  let tooltipData = $state<{
    alerts: PluginAlert[];
    position: { top: number; left: number };
  } | null>(null);

  function handleTooltipShow(event: Event) {
    const customEvent = event as CustomEvent;
    tooltipData = customEvent.detail;
    showTooltip = true;
  }

  function handleTooltipHide() {
    showTooltip = false;
    tooltipData = null;
  }

  $effect(() => {
    document.addEventListener("tooltip-show", handleTooltipShow);
    document.addEventListener("tooltip-hide", handleTooltipHide);

    return () => {
      document.removeEventListener("tooltip-show", handleTooltipShow);
      document.removeEventListener("tooltip-hide", handleTooltipHide);
    };
  });

  const instances = $derived(() => {
    const originalInstances = $pluginInstancesFromDB || [];

    return originalInstances;
  });

  const instancesArray = $derived(instances());

  let renderItems: Array<{
    type: "widget" | "gap";
    instance?: PluginInstance;
    gapPosition?: number;
  }> = $state([]);

  $effect(() => {
    if (instancesArray.length === 0) {
      renderItems = [];
      return;
    }

    const sortedInstances = [...instancesArray].sort(
      (a, b) => (a.order || 0) - (b.order || 0)
    );
    const items: Array<{
      type: "widget" | "gap";
      instance?: PluginInstance;
      gapPosition?: number;
    }> = [];

    let expectedPosition = 1;

    for (const instance of sortedInstances) {
      const currentPosition = instance.order || 0;
      while (expectedPosition < currentPosition) {
        items.push({ type: "gap", gapPosition: expectedPosition });
        expectedPosition++;
      }

      items.push({ type: "widget", instance });
      expectedPosition = currentPosition + (instance.span || 1);
    }

    renderItems = items;
  });

  const servicesData = $derived($services || []);
  const isLoading = $derived($isDashboardLoading || false);
  const error = $derived($dashboardError || null);

  let placeholders: Array<{ id: string; type: string; position: number }> =
    $state([]);

  $effect(() => {
    if (instancesArray.length === 0) {
      placeholders = [];
      return;
    }

    const gridColumns = 4; // Default to 4 columns for desktop

    let totalPositionsUsed = 0;
    renderItems.forEach((item) => {
      if (item.type === "widget" && item.instance) {
        totalPositionsUsed += Math.min(item.instance.span || 1, gridColumns);
      } else if (item.type === "gap") {
        totalPositionsUsed += 1;
      }
    });

    const currentRowUsed = totalPositionsUsed % gridColumns;
    const placeholdersInCurrentRow =
      currentRowUsed === 0 ? 0 : gridColumns - currentRowUsed;
    const placeholdersInNewRow = gridColumns;

    const maxPosition = Math.max(0, ...instancesArray.map((i) => i.order || 0));
    const maxSpan = Math.max(0, ...instancesArray.map((i) => i.span || 1));
    let nextPosition = maxPosition + maxSpan;

    const placeholderList: Array<{
      id: string;
      type: string;
      position: number;
    }> = [];

    for (let i = 0; i < placeholdersInCurrentRow; i++) {
      placeholderList.push({
        id: `current-${i}`,
        type: "current-row",
        position: nextPosition + i,
      });
    }

    for (let i = 0; i < placeholdersInNewRow; i++) {
      placeholderList.push({
        id: `new-${i}`,
        type: "new-row",
        position: nextPosition + placeholdersInCurrentRow + i,
      });
    }

    placeholders = placeholderList;
    console.log("Generated placeholders:", placeholderList);
  });
  function handleAddWidget(position?: number) {
    if (onaddwidget) {
      onaddwidget(position);
    }
  }
</script>

<div class="dashboard-container">
  {#if isLoading}
    <div class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading dashboard...</p>
    </div>
  {:else if error}
    <div class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Dashboard Error</h3>
      <p>{error}</p>
      <button onclick={() => window.location.reload()}>Retry</button>
    </div>
  {:else if renderItems.length > 0}
    <div class="dashboard-grid">
      {#each renderItems as item, index (item.type === "widget" ? item.instance?.id : `gap-${item.gapPosition}`)}
        {#if item.type === "widget" && item.instance}
          {@const serviceData = servicesData.find(
            (s) => s.id === item.instance?.id
          )}
          {@const pluginData = serviceData?.pluginData || {}}
          <div
            class="grid-item {editMode ? 'edit-mode' : ''}"
            class:span-2={item.instance.span === 2}
            class:span-3={item.instance.span === 3}
            class:span-4={item.instance.span === 4}
          >
            <PluginRenderer
              instance={item.instance}
              data={pluginData}
              {onedit}
              {editMode}
            />
          </div>
        {:else if item.type === "gap"}
          <div class="grid-item gap-item">
            <div class="gap-card">
              {#if editMode}
                <div
                  class="gap-placeholder"
                  role="button"
                  tabindex="0"
                  onclick={() => handleAddWidget(item.gapPosition)}
                  onkeydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      e.preventDefault();
                      handleAddWidget(item.gapPosition);
                    }
                  }}
                >
                  <div class="placeholder-icon">
                    <svg
                      width="24"
                      height="24"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <line x1="12" y1="5" x2="12" y2="19"></line>
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                    </svg>
                  </div>
                  <p class="placeholder-text">Fill Gap</p>
                  <p class="placeholder-position">
                    Position: {item.gapPosition}
                  </p>
                </div>
              {:else}
                <div class="empty-gap"></div>
              {/if}
            </div>
          </div>
        {/if}
      {/each}

      <!-- Placeholder Cards - Only shown in edit mode -->
      {#if editMode}
        {#each placeholders as placeholder (placeholder.id)}
          <div class="grid-item placeholder-item">
            <div
              class="placeholder-card"
              role="button"
              tabindex="0"
              onclick={() => handleAddWidget(placeholder.position)}
              onkeydown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  handleAddWidget(placeholder.position);
                }
              }}
            >
              <div class="placeholder-icon">
                <svg
                  width="32"
                  height="32"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <line x1="12" y1="5" x2="12" y2="19"></line>
                  <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
              </div>
              <p class="placeholder-text">Add Widget</p>
              <p class="placeholder-position">
                Position: {placeholder.position}
              </p>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  {:else if editMode}
    <div class="dashboard-grid">
      <!-- Show placeholder cards when no widgets exist and in edit mode -->
      {#each Array(4) as _, i}
        <div class="grid-item placeholder-item">
          <div
            class="placeholder-card empty-placeholder"
            role="button"
            tabindex="0"
            onclick={() => handleAddWidget(i + 1)}
            onkeydown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                handleAddWidget(i + 1);
              }
            }}
          >
            <div class="placeholder-icon">
              <svg
                width="32"
                height="32"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <line x1="12" y1="5" x2="12" y2="19"></line>
                <line x1="5" y1="12" x2="19" y2="12"></line>
              </svg>
            </div>
            <p class="placeholder-text">
              {i === 0 ? "Add Your First Widget" : "Add Widget"}
            </p>
            {#if i === 0}
              <p class="placeholder-hint">Click here to get started</p>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty-state">
      <div class="empty-icon">📊</div>
      <h3 class="empty-title">No Plugins Added Yet</h3>
      <p class="empty-description">
        Click the <strong>pencil</strong> button in the bottom-right corner to enter
        edit mode, then click on placeholder cards to add your first plugin to the
        dashboard.
      </p>
    </div>
  {/if}
</div>

<!-- Global tooltip portal - renders outside all grid constraints -->
{#if showTooltip && tooltipData}
  <div class="global-tooltip-portal">
    <div
      class="tooltip"
      style="
      position: fixed;
      top: {tooltipData.position.top}px;
      left: {tooltipData.position.left}px;
      z-index: 10000;
      pointer-events: none;
    "
    >
      <div class="tooltip-content">
        {#each tooltipData.alerts as alertItem, index (index)}
          <div class="alert-message">{alertItem.message}</div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .dashboard-container {
    width: 100%;
    padding: 20px;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1.5rem;
    margin-top: 2rem;
  }

  .grid-item {
    grid-column: span 1;
  }

  .grid-item.span-2 {
    grid-column: span 2;
  }

  .grid-item.span-3 {
    grid-column: span 3;
  }

  .grid-item.span-4 {
    grid-column: span 4;
  }

  .grid-item.edit-mode {
    position: relative;
    border: 2px dashed rgba(245, 158, 11, 0.5);
    border-radius: 12px;
    transition: all 0.2s ease;
  }

  .grid-item.edit-mode:hover {
    border-color: rgba(245, 158, 11, 0.8);
    background: rgba(245, 158, 11, 0.05);
  }

  /* Placeholder Card Styles */
  .placeholder-item {
    grid-column: span 1;
    min-height: 200px;
  }

  .placeholder-card {
    width: 100%;
    height: 100%;
    min-height: 200px;
    background: rgba(255, 255, 255, 0.02);
    border: 2px dashed rgba(255, 255, 255, 0.2);
    border-radius: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
  }

  .placeholder-card:hover {
    background: rgba(59, 130, 246, 0.1);
    border-color: rgba(59, 130, 246, 0.4);
    transform: translateY(-2px);
    box-shadow: 0 8px 32px rgba(59, 130, 246, 0.2);
  }

  .placeholder-card:active {
    transform: translateY(0);
  }

  .placeholder-card:focus {
    outline: 2px solid rgba(59, 130, 246, 0.5);
    outline-offset: 2px;
  }

  .placeholder-icon {
    color: rgba(255, 255, 255, 0.5);
    margin-bottom: 8px;
    transition: all 0.3s ease;
  }

  .placeholder-card:hover .placeholder-icon {
    color: rgba(59, 130, 246, 0.8);
    transform: scale(1.1);
  }

  .placeholder-text {
    color: rgba(255, 255, 255, 0.6);
    font-size: 14px;
    font-weight: 500;
    margin: 0;
    transition: all 0.3s ease;
  }

  .placeholder-card:hover .placeholder-text {
    color: rgba(59, 130, 246, 0.9);
  }

  /* Empty state */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    text-align: center;
    max-width: 100%;
    margin: 80px auto;
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.1),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
  }

  .empty-icon {
    font-size: 48px;
    margin-bottom: 20px;
    opacity: 0.7;
  }

  .empty-title {
    margin: 0 0 12px 0;
    font-size: 24px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
  }

  .empty-description {
    margin: 0;
    font-size: 16px;
    color: rgba(255, 255, 255, 0.7);
    line-height: 1.5;
  }

  .empty-description strong {
    color: rgba(255, 255, 255, 0.9);
    font-weight: 600;
    background: rgba(255, 255, 255, 0.1);
    padding: 2px 6px;
    border-radius: 6px;
  }

  /* Empty state placeholder cards (when in edit mode) */
  .empty-placeholder {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.3);
  }

  .empty-placeholder:hover {
    background: rgba(59, 130, 246, 0.15);
    border-color: rgba(59, 130, 246, 0.5);
  }

  .placeholder-hint {
    color: rgba(255, 255, 255, 0.4);
    font-size: 12px;
    margin: 4px 0 0 0;
    transition: all 0.3s ease;
  }

  .placeholder-card:hover .placeholder-hint {
    color: rgba(59, 130, 246, 0.7);
  }

  .placeholder-position {
    color: rgba(255, 255, 255, 0.3);
    font-size: 10px;
    margin: 2px 0 0 0;
    font-family: monospace;
  }

  .placeholder-card:hover .placeholder-position {
    color: rgba(59, 130, 246, 0.5);
  }

  /* Gap Cards */
  .gap-item {
    grid-column: span 1;
  }

  .gap-card {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .gap-placeholder {
    width: 100%;
    height: 100%;
    background: rgba(255, 193, 7, 0.05);
    border: 2px dashed rgba(255, 193, 7, 0.3);
    border-radius: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .gap-placeholder:hover {
    background: rgba(255, 193, 7, 0.1);
    border-color: rgba(255, 193, 7, 0.5);
    transform: translateY(-2px);
    box-shadow: 0 8px 32px rgba(255, 193, 7, 0.2);
  }

  .gap-placeholder .placeholder-icon {
    color: rgba(255, 193, 7, 0.6);
  }

  .gap-placeholder:hover .placeholder-icon {
    color: rgba(255, 193, 7, 0.8);
    transform: scale(1.1);
  }

  .gap-placeholder .placeholder-text {
    color: rgba(255, 193, 7, 0.7);
  }

  .gap-placeholder:hover .placeholder-text {
    color: rgba(255, 193, 7, 0.9);
  }

  .gap-placeholder .placeholder-position {
    color: rgba(255, 193, 7, 0.4);
  }

  .gap-placeholder:hover .placeholder-position {
    color: rgba(255, 193, 7, 0.6);
  }

  .empty-gap {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  @media (max-width: 1440px) {
    .dashboard-grid {
      grid-template-columns: repeat(3, 1fr);
    }

    .grid-item.span-4 {
      grid-column: span 3;
    }
  }

  @media (max-width: 1280px) {
    .dashboard-grid {
      grid-template-columns: repeat(2, 1fr);
    }

    .grid-item.span-3,
    .grid-item.span-4 {
      grid-column: span 2;
    }
  }

  @media (max-width: 768px) {
    .dashboard-container {
      padding: 10px;
    }

    .empty-state {
      margin: 40px auto;
      padding: 40px 20px;
    }

    .empty-icon {
      font-size: 36px;
    }

    .empty-title {
      font-size: 20px;
    }

    .empty-description {
      font-size: 14px;
    }
  }

  .loading-state,
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    text-align: center;
    max-width: 100%;
    margin: 80px auto;
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(255, 255, 255, 0.2);
    border-top: 3px solid #fff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 20px;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .error-icon {
    font-size: 48px;
    margin-bottom: 20px;
  }

  .error-state h3 {
    color: #ff6b6b;
    margin-bottom: 10px;
    font-size: 24px;
  }

  .error-state p {
    color: rgba(255, 255, 255, 0.8);
    margin-bottom: 20px;
    font-size: 16px;
  }

  .error-state button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .error-state button:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }

  @media (max-width: 600px) {
    .dashboard-grid {
      grid-template-columns: 1fr;
    }

    .grid-item,
    .grid-item.span-2,
    .grid-item.span-3,
    .grid-item.span-4,
    .placeholder-item,
    .gap-item {
      grid-column: span 1;
    }

    .placeholder-card,
    .gap-placeholder,
    .empty-gap {
      min-height: 150px;
    }
  }

  /* Global tooltip styles */
  .global-tooltip-portal {
    pointer-events: none;
  }

  .tooltip-content {
    background: rgba(17, 24, 39, 0.95);
    color: white;
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 11px;

    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    backdrop-filter: blur(8px);
    white-space: nowrap;
  }

  .alert-message {
    font-weight: 500;
    margin-bottom: 2px;
  }

</style>
