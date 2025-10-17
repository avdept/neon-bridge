<script lang="ts">
  interface Props {
    icon?: string;
    label?: string;
    position?: "bottom-right" | "bottom-left" | "top-right" | "top-left";
    onclick?: () => void;
    class?: string;
  }

  const {
    icon = "plus",
    label = "Add",
    position = "bottom-right",
    onclick,
    class: customClass = "",
  }: Props = $props();

  function handleClick() {
    if (onclick) {
      onclick();
    }
  }

  // Dynamic import for the icon - using public directory
  const iconSrc = $derived(`/icons/${icon}.svg`);
</script>

<button
  class="fab {position} {customClass}"
  onclick={handleClick}
  aria-label={label}
  title={label}
>
  {#if icon}
    <div class="icon" style="background-image: url('{iconSrc}')"></div>
  {/if}
</button>

<style>
  .fab {
    position: fixed;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: none;
    color: rgba(255, 255, 255, 0.9);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease;
    z-index: 100;
  }

  .fab:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: scale(1.1);
  }

  .fab:active {
    transform: scale(0.95);
  }

  .icon {
    width: 24px;
    height: 24px;
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    filter: brightness(0) invert(1);
    opacity: 0.85;
    transition: opacity 0.3s ease;
  }

  .fab:hover .icon {
    opacity: 1;
  }

  /* Positioning classes */
  .fab.bottom-right {
    bottom: 2rem;
    right: 2rem;
  }

  .fab.bottom-left {
    bottom: 2rem;
    left: 2rem;
  }

  .fab.top-right {
    top: 2rem;
    right: 2rem;
  }

  .fab.top-left {
    top: 2rem;
    left: 2rem;
  }

  /* Mobile responsive */
  @media (max-width: 768px) {
    .fab {
      width: 44px;
      height: 44px;
    }

    .fab.bottom-right,
    .fab.bottom-left {
      bottom: 1rem;
    }

    .fab.bottom-right {
      right: 1rem;
    }

    .fab.bottom-left {
      left: 1rem;
    }

    .fab.top-right,
    .fab.top-left {
      top: 1rem;
    }

    .fab.top-right {
      right: 1rem;
    }

    .fab.top-left {
      left: 1rem;
    }
  }

  @media (max-width: 480px) {
    .fab {
      width: 40px;
      height: 40px;
    }

    .fab.bottom-right,
    .fab.bottom-left,
    .fab.top-right,
    .fab.top-left {
      bottom: 0.75rem;
      right: 0.75rem;
    }

    .fab.bottom-left {
      left: 0.75rem;
      right: auto;
    }

    .fab.top-right,
    .fab.top-left {
      top: 0.75rem;
      bottom: auto;
    }

    .fab.top-left {
      left: 0.75rem;
      right: auto;
    }
  }
</style>
