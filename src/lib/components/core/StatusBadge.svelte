<script lang="ts">
  interface Props {
    status?: "online" | "offline" | "warning" | "loading";
    pulse?: boolean;
    showDot?: boolean;
    children?: import("svelte").Snippet;
  }

  const {
    status = "online",
    pulse = true,
    showDot = true,
    children,
  }: Props = $props();

  const statusConfig = {
    online: {
      color: "#34d399",
      bg: "rgba(52, 211, 153, 0.15)",
      border: "rgba(52, 211, 153, 0.3)",
    },
    offline: {
      color: "#f87171",
      bg: "rgba(248, 113, 113, 0.15)",
      border: "rgba(248, 113, 113, 0.3)",
    },
    warning: {
      color: "#fbbf24",
      bg: "rgba(251, 191, 36, 0.15)",
      border: "rgba(251, 191, 36, 0.3)",
    },
    loading: {
      color: "#60a5fa",
      bg: "rgba(96, 165, 250, 0.15)",
      border: "rgba(96, 165, 250, 0.3)",
    },
  };
</script>

<div
  class="status-badge"
  style="
    background: {statusConfig[status].bg};
    border-color: {statusConfig[status].border};
    color: {statusConfig[status].color}
  "
>
  {#if showDot}
    <div
      class="status-dot"
      class:pulse
      style="background: {statusConfig[status].color}"
    ></div>
  {/if}
  {@render children?.()}
</div>

<style>
  .status-badge {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border: 1px solid;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 500;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  .status-dot.pulse {
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.7;
      transform: scale(1.2);
    }
  }
</style>
