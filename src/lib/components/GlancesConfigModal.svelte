<script lang="ts">
  import { dashboardStore } from "../stores/dashboard.js";
  import Modal from "./core/Modal.svelte";

  export let isOpen: boolean = false;
  export let dashboardId: number;
  export let initialConfig: GlancesConfig | null = null;
  export let onClose: (() => void) ;

  interface GlancesConfig {
    url: string;
    username?: string;
    password?: string;
  }

  let config: GlancesConfig = {
    url: "",
    username: "",
    password: "",
  };

  let isSaving = false;
  let configError: string | null = null;

  $: if (isOpen && initialConfig) {
    config = { ...initialConfig };
  } else if (isOpen && !initialConfig) {
    config = {
      url: "",
      username: "",
      password: "",
    };
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    if (!config.url.trim()) {
      configError = "URL is required";
      return;
    }

    isSaving = true;
    configError = null;

    try {
      const success = await dashboardStore.updateDashboard({
        glances_config: JSON.stringify(config),
      });

      if (success) {
        onClose();
      } else {
        throw new Error("Failed to save configuration");
      }
    } catch (error) {
      configError =
        error instanceof Error ? error.message : "Failed to save configuration";
    } finally {
      isSaving = false;
    }
  }

  async function handleRemoveConfig() {
    try {
      const success = await dashboardStore.updateDashboard({
        glances_config: "",
      });

      if (success) {
        onClose();
      } else {
        throw new Error("Failed to remove configuration");
      }
    } catch (error) {
      configError =
        error instanceof Error
          ? error.message
          : "Failed to remove configuration";
    }
  }

  function handleClose() {
    configError = null;
    onClose();
  }
</script>

<Modal
  {isOpen}
  title="Configure Glances Stats"
  maxWidth="600px"
  onclose={handleClose}
>
  <div class="glances-config">
    <div class="header-info">
      <div class="glances-icon">
        <img
          src="/services/glances.svg"
          alt="Glances Icon"
          width="40"
          height="40"
        />
      </div>
      <div>
        <p class="modal-description">
          Configure connection to Glances server for real-time system statistics
        </p>
      </div>
    </div>

    {#if configError}
      <div class="error-message">
        <div class="error-icon">⚠️</div>
        <div class="error-content">
          <h4>Configuration Error</h4>
          <p>{configError}</p>
        </div>
      </div>
    {/if}

    <form onsubmit={handleSubmit} class="config-form">
      <div class="form-group">
        <label class="form-label" for="glances-url">
          Glances Server URL
          <span class="required">*</span>
        </label>
        <p class="form-description">
          The URL to your Glances server (e.g., http://localhost:61208)
        </p>
        <input
          type="url"
          id="glances-url"
          bind:value={config.url}
          placeholder="http://localhost:61208"
          required
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="glances-username">
          Username (Optional)
        </label>
        <p class="form-description">
          Username if your Glances server requires authentication
        </p>
        <input
          type="text"
          id="glances-username"
          bind:value={config.username}
          placeholder="Enter username"
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="glances-password">
          Password (Optional)
        </label>
        <p class="form-description">
          Password if your Glances server requires authentication
        </p>
        <input
          type="password"
          id="glances-password"
          bind:value={config.password}
          placeholder="Enter password"
          class="form-input"
        />
      </div>

      <div class="step-actions">
        <div class="left-actions">
          <button
            type="button"
            class="btn-secondary"
            onclick={handleClose}
            disabled={isSaving}
          >
            Cancel
          </button>
          {#if initialConfig}
            <button
              type="button"
              class="btn-danger"
              onclick={handleRemoveConfig}
              disabled={isSaving}
            >
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
              Remove Config
            </button>
          {/if}
        </div>

        <div class="action-buttons">
          <button type="submit" class="btn-primary" disabled={isSaving}>
            {isSaving ? "Saving..." : "Save Configuration"}
          </button>
        </div>
      </div>
    </form>
  </div>
</Modal>

<style>
  .glances-config {
    width: 100%;
  }

  .header-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    border: 1px solid rgba(255, 255, 255, 0.08);
  }

  .glances-icon {
    flex-shrink: 0;
  }

  .modal-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.875rem;
    margin: 0;
    line-height: 1.4;
  }

  .error-message {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    padding: 1rem;
    background: rgba(248, 113, 113, 0.1);
    border: 1px solid rgba(248, 113, 113, 0.3);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    color: #f87171;
  }

  .error-icon {
    font-size: 1.5rem;
    flex-shrink: 0;
    margin-top: 0.125rem;
  }

  .error-content h4 {
    margin: 0 0 0.5rem 0;
    font-size: 1rem;
    font-weight: 600;
  }

  .error-content p {
    margin: 0;
    font-size: 0.875rem;
    line-height: 1.4;
    opacity: 0.9;
  }

  .config-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-label {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.875rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .required {
    color: #ef4444;
    font-weight: 600;
  }

  .form-description {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.75rem;
    margin: 0;
    line-height: 1.4;
  }

  .form-input {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 8px;
    padding: 0.75rem;
    color: white;
    font-size: 0.875rem;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

  .form-input:focus {
    outline: none;
    border-color: #3b82f6;
    background: rgba(255, 255, 255, 0.08);
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  .form-input::placeholder {
    color: rgba(255, 255, 255, 0.4);
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
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    min-width: 120px;
  }

  .btn-primary:disabled,
  .btn-secondary:disabled,
  .btn-danger:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: rgba(255, 255, 255, 0.8);
  }

  .btn-secondary:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.15);
    color: white;
  }

  .btn-primary {
    background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
    color: white;
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  }

  .btn-primary:hover:not(:disabled) {
    background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
  }

  .btn-danger {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
    color: white;
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
  }

  .btn-danger:hover:not(:disabled) {
    background: linear-gradient(135deg, #dc2626 0%, #b91c1c 100%);
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4);
  }

  @media (max-width: 768px) {
    .header-info {
      gap: 0.75rem;
      padding: 0.75rem;
    }

    .step-actions {
      flex-direction: column;
      gap: 1rem;
    }

    .left-actions,
    .action-buttons {
      width: 100%;
      justify-content: center;
    }

    .left-actions {
      order: 2;
    }

    .action-buttons {
      order: 1;
    }
  }
</style>
