<script lang="ts">
  import PluginSelector from "./PluginSelector.svelte";
  import PluginConfigForm from "./PluginConfigForm.svelte";
  import { dashboardStore } from "../stores/dashboard.js";
  import type { PluginInstance } from "../plugins/types.js";
  import { pluginRegistry } from "../plugins/registry.js";

  export let currentStep = 1;
  export let selectedPluginId: string | null = null;
  export let editMode: boolean = false;
  export let initialPluginId: string | undefined = undefined;
  export let initialConfig: any = undefined;
  export let closeModal: () => void;
  export let position: number | null = null;
  export let editingWidget: PluginInstance | null = null;

  $: if (editMode && initialPluginId) {
    selectedPluginId = initialPluginId;
    currentStep = 2;
  }

  let isTestingConfig = false;
  let errorMessage: string | null = null;

  function handlePluginSelect(pluginId: string | null) {
    selectedPluginId = pluginId;
  }

  function nextStep() {
    if (currentStep === 1 && selectedPluginId) {
      currentStep = 2;
    }
  }

  function previousStep() {
    if (currentStep > 1) {
      currentStep = currentStep - 1;
    }
  }



  async function handleConfigSubmit(pluginId: string, config: Record<string, any>) {
    errorMessage = null;
    try {
      const plugin = pluginRegistry.get(pluginId);
      if (plugin && plugin.fetchData) {
        await plugin.fetchData(config, editingWidget?.id, true);
      }

      if (editingWidget) {
        const updatedWidget = await dashboardStore.updateWidget(
          editingWidget.id,
          {
            name: config.title || pluginId,
            config: {
              title: config.title || pluginId,
              ...config,
            },
            position: editingWidget.order,
          }
        );

        if (updatedWidget) {
          console.log("Widget updated successfully:", updatedWidget);
          // todo: refactor this to use props passing instead of events
          isTestingConfig = false;

          closeModal();
        } else {
          throw new Error("Failed to update widget in database");
        }
      } else {
        const newInstance: PluginInstance = {
          id: Date.now(),
          pluginId,
          config: {
            title: config.title || pluginId,
            ...config,
          },
          span: config.span || 1,
          order: position || Date.now(),
          enabled: true,
        };

        const savedWidget = await dashboardStore.addWidget(newInstance);

        if (savedWidget) {
          // todo: refactor this to use props passing instead of events
          isTestingConfig = false;
          closeModal();
        } else {
          throw new Error("Failed to save widget to database");
        }
      }
    } catch (error) {
      console.error("Plugin configuration test failed:", error);

      // todo: refactor this to use props passing instead of events
      errorMessage =
        error instanceof Error ? error.message : "Configuration test failed";
    }
  }


  async function handleRemove() {
    if (
      confirm(
        "Are you sure you want to remove this widget? This action cannot be undone."
      )
    ) {
      if (!editMode) return;

    try {
      const success = await dashboardStore.removeWidget(
        editingWidget!.id
      );

      if (success) {
        closeModal();
      } else {
        throw new Error("Failed to remove widget from database");
      }
    } catch (error) {
      console.error("Failed to remove widget:", error);
      // todo: add some sort of notification to show results
      alert(error instanceof Error ? error.message : "Failed to remove widget");
    }
    }
  }
</script>

<div class="add-plugin-form">
  <!-- Step Indicator -->
  <div class="step-indicator">
    <div
      class="step {currentStep >= 1 ? 'active' : ''} {currentStep > 1
        ? 'completed'
        : ''}"
    >
      <div class="step-number">
        {#if currentStep > 1}
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path
              d="M13.333 4L6 11.333 2.667 8"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        {:else}
          1
        {/if}
      </div>
      <span class="step-label">Select Plugin</span>
    </div>

    <div class="step-connector"></div>

    <div class="step {currentStep >= 2 ? 'active' : ''}">
      <div class="step-number">2</div>
      <span class="step-label">Configure</span>
    </div>
  </div>

  <!-- Step Content -->
  <div class="step-content">
    {#if currentStep === 1}
      <PluginSelector bind:selectedPluginId onSelect={handlePluginSelect} />
    {:else if currentStep === 2 && selectedPluginId}
      <PluginConfigForm
        {selectedPluginId}
        {initialConfig}
        onSubmit={handleConfigSubmit}
        errorMessage={errorMessage}
      />
    {/if}
  </div>

  <!-- Step Actions -->
  <div class="step-actions">
    <div class="left-actions">
      <button class="btn-secondary" on:click={closeModal}> Cancel </button>
      {#if editMode}
        <button class="btn-danger" on:click={handleRemove}>
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polyline points="3,6 5,6 21,6"></polyline>
            <path
              d="m19,6 v14c0,1 -1,2 -2,2H7c-1,0 -2,-1 -2,-2V6m3,0V4c0,-1 1,-2 2,-2h4c1,0 2,1 2,2v2"
            ></path>
            <line x1="10" y1="11" x2="10" y2="17"></line>
            <line x1="14" y1="11" x2="14" y2="17"></line>
          </svg>
          Remove Widget
        </button>
      {/if}
    </div>

    <div class="action-buttons">
      {#if currentStep > 1 && !editMode}
        <button class="btn-secondary" on:click={previousStep}> Back </button>
      {/if}

      {#if currentStep === 1}
        <button
          class="btn-primary"
          disabled={!selectedPluginId}
          on:click={nextStep}
        >
          Continue
        </button>
      {:else if currentStep === 2}
        <button
          class="btn-primary"
          type="submit"
          form="plugin-config-form"
          disabled={isTestingConfig}
        >
          {#if isTestingConfig}
            <span class="loading-spinner"></span>
            Testing Configuration...
          {:else}
            {editMode ? "Save Changes" : "Add Plugin"}
          {/if}
        </button>
      {/if}
    </div>
  </div>
</div>

<style>
  .add-plugin-form {
    width: 100%;
  }

  .step-indicator {
    display: flex;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .step {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: rgba(255, 255, 255, 0.5);
    transition: color 0.2s ease;
  }

  .step.active {
    color: rgba(255, 255, 255, 0.9);
  }

  .step.completed {
    color: #3b82f6;
  }

  .step-number {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    border: 2px solid rgba(255, 255, 255, 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    font-weight: 600;
    transition: all 0.2s ease;
  }

  .step.active .step-number {
    background: rgba(59, 130, 246, 0.2);
    border-color: #3b82f6;
    color: #3b82f6;
  }

  .step.completed .step-number {
    background: #3b82f6;
    border-color: #3b82f6;
    color: white;
  }

  .step-label {
    font-size: 0.875rem;
    font-weight: 500;
  }

  .step-connector {
    flex: 1;
    height: 2px;
    background: rgba(255, 255, 255, 0.1);
    margin: 0 1rem;
  }

  .step-content {
    min-height: 300px;
    margin-bottom: 2rem;
  }

  .step-2-placeholder {
    text-align: center;
    padding: 2rem;
    color: rgba(255, 255, 255, 0.7);
  }

  .step-2-placeholder h3 {
    color: white;
    margin-bottom: 1rem;
  }

  .step-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 1rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }

  .left-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .action-buttons {
    display: flex;
    gap: 0.75rem;
  }

  .btn-primary,
  .btn-secondary,
  .btn-danger {
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-weight: 500;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    border: none;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-primary {
    background: #3b82f6;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #2563eb;
    transform: translateY(-1px);
  }

  .btn-primary:disabled {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.4);
    cursor: not-allowed;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.12);
  }

  .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.12);
    color: white;
    border-color: rgba(255, 255, 255, 0.2);
  }

  .btn-danger {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.3);
  }

  .btn-danger:hover {
    background: rgba(239, 68, 68, 0.2);
    color: #dc2626;
    border-color: rgba(239, 68, 68, 0.5);
    transform: translateY(-1px);
  }

  .btn-danger svg {
    width: 16px;
    height: 16px;
  }

  .loading-spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-right: 0.5rem;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  @media (max-width: 768px) {
    .step-actions {
      flex-direction: column;
      gap: 1rem;
    }

    .action-buttons {
      width: 100%;
      justify-content: center;
    }
  }
</style>
