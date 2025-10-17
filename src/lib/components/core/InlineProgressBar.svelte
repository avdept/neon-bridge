<script lang="ts">
  interface Props {
    title: string;
    value?: number;
    max?: number;
    color?: "primary" | "warning" | "danger" | "info";
    height?: string;
    status?: string;
    showPercentage?: boolean;
  }

  const {
    title,
    value = 0,
    max = 100,
    color = "primary",
    height = "24px",
    status,
    showPercentage = true,
  }: Props = $props();

  const percentage = $derived(Math.min(Math.max((value / max) * 100, 0), 100));

  const colorVariants = {
    primary: "linear-gradient(90deg, #34d399, #10b981)",
    warning: "linear-gradient(90deg, #fbbf24, #f59e0b)",
    danger: "linear-gradient(90deg, #f87171, #ef4444)",
    info: "linear-gradient(90deg, #60a5fa, #3b82f6)",
  };

  // Format percentage to display
  const displayPercentage = $derived(percentage.toFixed(1) + "%");

  // Determine text color based on progress for better contrast
  const textColor = $derived(
    percentage > 50 ? "rgba(255, 255, 255, 0.95)" : "rgba(255, 255, 255, 0.8)"
  );
</script>

<div class="inline-progress-bar" style="height: {height}">
  <div
    class="progress-fill"
    style="width: {percentage}%; background: {colorVariants[color]}"
  ></div>

  <div class="progress-content" style="color: {textColor}">
    <div class="progress-title">{title}</div>
    <div class="progress-status">
      {#if status}
        {status}
      {:else if showPercentage}
        {displayPercentage}
      {/if}
    </div>
  </div>
</div>

<style>
  .inline-progress-bar {
    position: relative;
    width: 100%;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
    overflow: hidden;
    margin: 0.5rem 0;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .progress-fill {
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    border-radius: 2px;
    transition: width 0.5s ease;
    z-index: 1;
  }

  .progress-content {
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 100%;
    padding: 0 10px;
    z-index: 2;
    font-size: 12px;
    font-weight: 500;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  }

  .progress-title {
    font-weight: 600;
    letter-spacing: 0.3px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
    margin-right: 8px;
  }

  .progress-status {
    font-family: monospace;
    font-size: 11px;
    font-weight: 500;
    white-space: nowrap;
    opacity: 0.95;
  }

  /* Responsive adjustments */
  @media (max-width: 600px) {
    .progress-content {
      padding: 0 8px;
      font-size: 11px;
    }

    .progress-status {
      font-size: 10px;
    }
  }
</style>
