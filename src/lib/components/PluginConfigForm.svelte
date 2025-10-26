<script lang="ts">
  import { pluginRegistry } from "../plugins/registry.js";
  import type { Plugin, PluginConfigField } from "../plugins/types.js";


  export let selectedPluginId: string;
  export let errorMessage: string | null = null;
  export let initialConfig: Record<string, any> | undefined = undefined;
  export let onSubmit: (pluginId: string, config: Record<string, any>) => void;

  let plugin: Plugin | undefined;
  let config: Record<string, any> = {};


  $: {
    plugin = pluginRegistry.get(selectedPluginId);
    if (plugin) {
      config = {};
      plugin.configTemplate.fields.forEach((field) => {
        if (initialConfig && initialConfig[field.key] !== undefined) {
          config[field.key] = initialConfig[field.key];
        } else {
          config[field.key] =
            field.default || getDefaultValueForType(field.type);
        }
      });
    }
  }

  function getDefaultValueForType(type: string): any {
    switch (type) {
      case "boolean":
        return false;
      case "number":
        return 0;
      default:
        return "";
    }
  }

  function handleSubmit() {
    onSubmit(selectedPluginId, config);
  }

  function renderField(field: PluginConfigField) {
    switch (field.type) {
      case "text":
      case "url":
      case "email":
        return "text";
      case "password":
        return "password";
      case "number":
        return "number";
      default:
        return "text";
    }
  }
</script>

<div class="plugin-config">
  {#if plugin}
    {@const iconPath = `/services/${plugin.metadata.icon}.svg`}
    <div class="plugin-header">
      <div
        class="plugin-icon"
        style="background-image: url('{iconPath}')"
      ></div>
      <div class="plugin-info">
        <h3 class="plugin-name">{plugin.metadata.name}</h3>
        <p class="plugin-description">{plugin.metadata.description}</p>
      </div>
    </div>

    {#if errorMessage}
      <div class="error-message">
        <div class="error-icon">⚠️</div>
        <div class="error-content">
          <h4>Configuration Test Failed</h4>
          <p>{errorMessage}</p>
        </div>
      </div>
    {/if}

    <form
      class="config-form"
      id="plugin-config-form"
      on:submit|preventDefault={handleSubmit}
    >
      {#each plugin.configTemplate.fields as field (field.key)}
        <div class="form-group">
          <label class="form-label" for={field.key}>
            {field.label}
            {#if field.required}
              <span class="required">*</span>
            {/if}
          </label>

          {#if field.description}
            <p class="form-description">{field.description}</p>
          {/if}

          {#if field.type === "boolean"}
            <label class="checkbox-wrapper">
              <input
                type="checkbox"
                id={field.key}
                bind:checked={config[field.key]}
                class="checkbox"
              />
              <span class="checkbox-label"
                >Enable {field.label.toLowerCase()}</span
              >
            </label>
          {:else if field.type === "select" && field.options}
            <select
              id={field.key}
              bind:value={config[field.key]}
              class="form-select"
              required={field.required}
            >
              <option value="">Select an option...</option>
              {#each field.options as option}
                <option value={option.value}>{option.label}</option>
              {/each}
            </select>
          {:else}
            <input
              type={renderField(field)}
              id={field.key}
              bind:value={config[field.key]}
              placeholder={field.placeholder}
              required={field.required}
              class="form-input"
            />
          {/if}
        </div>
      {/each}
    </form>
  {:else}
    <div class="error">
      <p>Plugin not found: {selectedPluginId}</p>
    </div>
  {/if}
</div>

<style>
  .plugin-config {
    width: 100%;
  }

  .plugin-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 12px;
    margin-bottom: 1.5rem;
    border: 1px solid rgba(255, 255, 255, 0.08);
  }

  .plugin-icon {
    font-size: 2.5rem;
    width: 3.5rem;
    height: 3.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    flex-shrink: 0;
  }

  .plugin-info {
    flex: 1;
  }

  .plugin-name {
    color: white;
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0 0 0.25rem 0;
    background: linear-gradient(135deg, #ffffff 0%, #e8e8e8 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .plugin-description {
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

  .form-input,
  .form-select {
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

  .form-input:focus,
  .form-select:focus {
    outline: none;
    border-color: #3b82f6;
    background: rgba(255, 255, 255, 0.08);
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  .form-input::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .checkbox-wrapper {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    padding: 0.75rem;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  }

  .checkbox-wrapper:hover {
    background: rgba(255, 255, 255, 0.06);
    border-color: rgba(255, 255, 255, 0.15);
  }

  .checkbox {
    width: 1.25rem;
    height: 1.25rem;
    accent-color: #3b82f6;
    cursor: pointer;
  }

  .checkbox-label {
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.875rem;
    cursor: pointer;
  }

  .error {
    padding: 2rem;
    text-align: center;
    color: rgba(255, 255, 255, 0.7);
  }

  @media (max-width: 768px) {
    .plugin-header {
      padding: 0.75rem;
      gap: 0.75rem;
    }

    .plugin-icon {
      font-size: 2rem;
      width: 3rem;
      height: 3rem;
    }

    .plugin-name {
      font-size: 1.125rem;
    }

    .form-input,
    .form-select {
      padding: 0.625rem;
    }
  }
</style>
