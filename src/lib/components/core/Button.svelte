<script lang="ts">
  interface Props {
    variant?: "primary" | "secondary" | "danger";
    size?: "small" | "medium" | "large";
    disabled?: boolean;
    loading?: boolean;
    onclick?: () => void;
    children?: import("svelte").Snippet;
  }

  const {
    variant = "secondary",
    size = "medium",
    disabled = false,
    loading = false,
    onclick,
    children,
  }: Props = $props();

  function handleClick() {
    if (!disabled && !loading && onclick) {
      onclick();
    }
  }
</script>

<button
  class="btn btn-{variant} btn-{size}"
  class:disabled
  class:loading
  on:click={handleClick}
  {disabled}
>
  {#if loading}
    <div class="spinner"></div>
  {/if}
  {@render children?.()}
</button>

<style>
  .btn {
    padding: 0.5rem 0.75rem;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 10px;
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.8rem;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.4rem;
    font-family: inherit;
    font-weight: 500;
  }

  .btn:hover:not(.disabled) {
    background: rgba(255, 255, 255, 0.2);
    transform: translateY(-2px);
  }

  .btn:active:not(.disabled) {
    transform: translateY(-1px) scale(0.98);
  }

  .btn-primary {
    background: rgba(52, 211, 153, 0.2);
    border-color: rgba(52, 211, 153, 0.3);
    color: #34d399;
  }

  .btn-primary:hover:not(.disabled) {
    background: rgba(52, 211, 153, 0.3);
  }

  .btn-danger {
    background: rgba(248, 113, 113, 0.2);
    border-color: rgba(248, 113, 113, 0.3);
    color: #f87171;
  }

  .btn-danger:hover:not(.disabled) {
    background: rgba(248, 113, 113, 0.3);
  }

  .btn-small {
    padding: 0.5rem 0.75rem;
    font-size: 0.8rem;
  }

  .btn-large {
    padding: 1rem 1.5rem;
    font-size: 1rem;
  }

  .disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
  }

  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid transparent;
    border-top: 2px solid currentColor;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>
