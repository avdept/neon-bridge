<script lang="ts">
  interface Props {
    value?: number;
    max?: number;
    color?: "primary" | "warning" | "danger" | "info";
    height?: string;
    animated?: boolean;
  }

  const {
    value = 0,
    max = 100,
    color = "primary",
    height = "8px",
    animated = true,
  }: Props = $props();

  const percentage = $derived(Math.min(Math.max((value / max) * 100, 0), 100));

  const colorVariants = {
    primary: "linear-gradient(90deg, #34d399, #10b981)",
    warning: "linear-gradient(90deg, #fbbf24, #f59e0b)",
    danger: "linear-gradient(90deg, #f87171, #ef4444)",
    info: "linear-gradient(90deg, #60a5fa, #3b82f6)",
  };
</script>

<div class="progress-bar" style="height: {height}">
  <div
    class="progress-fill"
    class:animated
    style="width: {percentage}%; background: {colorVariants[color]}"
  ></div>
</div>

<style>
  .progress-bar {
    width: 100%;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 3px;
    overflow: hidden;
    margin: 0.75rem 0;
  }

  .progress-fill {
    height: 100%;
    border-radius: 4px;
    transition: width 0.5s ease;
  }

  .progress-fill.animated {
    background-size: 200% 100%;
    animation: shimmer 2s linear infinite;
  }

  @keyframes shimmer {
    0% {
      background-position: -200% 0;
    }
    100% {
      background-position: 200% 0;
    }
  }
</style>
