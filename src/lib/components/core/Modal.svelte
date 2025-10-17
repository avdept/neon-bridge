<script lang="ts">
  interface Props {
    isOpen?: boolean;
    title?: string;
    maxWidth?: string;
    onclose?: () => void;
    children?: import("svelte").Snippet;
  }

  const {
    isOpen = false,
    title = "",
    maxWidth = "500px",
    onclose,
    children,
  }: Props = $props();

  function closeModal() {
    if (onclose) {
      onclose();
    }
  }

  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      closeModal();
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      closeModal();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
  <div
    class="modal-backdrop"
    on:click={handleBackdropClick}
    role="presentation"
  >
    <div class="modal-container" style="max-width: {maxWidth}">
      <div class="modal-header">
        {#if title}
          <h2 class="modal-title">{title}</h2>
        {/if}
        <button
          class="close-button"
          on:click={closeModal}
          aria-label="Close modal"
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M15 5L5 15M5 5L15 15"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
      </div>
      <div class="modal-content">
        {@render children?.()}
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    animation: fadeIn 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  }

  .modal-container {
    background: rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(40px);
    -webkit-backdrop-filter: blur(40px);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 24px;
    width: 100%;
    max-height: 90vh;
    overflow: hidden;
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.3),
      0 2px 8px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
    animation: slideIn 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    margin-bottom: 1.5rem;
  }

  .modal-title {
    color: white;
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    background: linear-gradient(135deg, #ffffff 0%, #e8e8e8 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .close-button {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 50%;
    color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
    padding: 0;
    width: 2.5rem;
    height: 2.5rem;
    transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
  }

  .close-button:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.2);
    color: white;
    transform: scale(1.05);
  }

  .close-button:active {
    transform: scale(0.95);
  }

  .modal-content {
    padding: 0 1.5rem 1.5rem 1.5rem;
    color: rgba(255, 255, 255, 0.9);
    max-height: calc(90vh - 120px);
    overflow-y: auto;
  }

  /* Custom scrollbar for webkit browsers */
  .modal-content::-webkit-scrollbar {
    width: 6px;
  }

  .modal-content::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 3px;
  }

  .modal-content::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
  }

  .modal-content::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.3);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes slideIn {
    from {
      transform: scale(0.8) translateY(20px);
      opacity: 0;
    }
    to {
      transform: scale(1) translateY(0);
      opacity: 1;
    }
  }

  /* Mobile responsive */
  @media (max-width: 768px) {
    .modal-backdrop {
      padding: 0.5rem;
    }

    .modal-container {
      max-height: 95vh;
      border-radius: 20px;
    }

    .modal-header {
      padding: 1rem 1rem 0 1rem;
      margin-bottom: 1rem;
    }

    .modal-content {
      padding: 0 1rem 1rem 1rem;
    }

    .modal-title {
      font-size: 1.25rem;
    }

    .close-button {
      width: 2rem;
      height: 2rem;
    }
  }
</style>
