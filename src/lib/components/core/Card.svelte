<script lang="ts">
  import type { ComponentType } from "svelte";
  import type { PluginAlert } from "../../plugins/types.js";

  interface Props {
    title: string;
    icon?: ComponentType | string;
    status?: string;
    statusType?: "online" | "offline" | "warning";
    href?: string;
    large?: boolean;
    span?: number; // 1 = normal, 2 = 2x width, 3 = 3x width, 4 = full width
    alerts?: PluginAlert[];
    children?: import("svelte").Snippet;
  }

  const {
    title,
    icon = "",
    status = "",
    statusType = "online",
    href = "",
    large = false,
    span = 1,
    alerts = [],
    children,
  }: Props = $props();

  let statusElement: HTMLElement | undefined = $state();

  const effectiveStatus = $derived(status);
  const effectiveStatusType = $derived(statusType);

  const getAlertIcon = (alertLevel: "warning" | "error") => {
    if (alertLevel === "warning") {
      return `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--alert-warning-color)" stroke-width="2" style="opacity: 0.8; transition: opacity 0.2s ease;">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
      </svg>`;
    } else {
      return `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--alert-error-color)" stroke-width="2" style="opacity: 0.8; transition: opacity 0.2s ease;">
        <circle cx="12" cy="12" r="10"/>
        <path d="m15 9-6 6"/>
        <path d="m9 9 6 6"/>
      </svg>`;
    }
  };

  const formatTime = (timestamp?: Date) => {
    if (timestamp) {
      return new Date(timestamp).toLocaleString();
    }
    return "";
  };

  function handleMouseEnter(alertItems: PluginAlert[], element: HTMLElement) {
    const rect = element.getBoundingClientRect();
    const event = new CustomEvent("tooltip-show", {
      detail: {
        alerts: alertItems,
        position: {
          top: rect.top - 60,
          left: rect.right - element.clientWidth * 2,
        },
      },
    });
    document.dispatchEvent(event);
  }

  function handleMouseLeave() {
    const event = new CustomEvent("tooltip-hide");
    document.dispatchEvent(event);
  }

  const iconPath = $derived(`/services/${icon}.svg`);

  const statusClasses = {
    online: "status-online",
    offline: "status-offline",
    warning: "status-warning",
  };
</script>

<div
  class="card"
  class:large-card={large}
  class:card-span-2={span === 2}
  class:card-span-3={span === 3}
  class:card-span-4={span === 4}
>
  <div class="card-header">
    <div class="card-inner">
      <div class="card-icon" style="background-image: url('{iconPath}')"></div>
      <a
        class="card-title"
        class:clickable={!!href}
        {href}
        onkeydown={href
          ? (e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                window.location.href = href;
              }
            }
          : undefined}
        title={href ? `Go to ${href}` : undefined}
      >
        {title}
      </a>
    </div>

    {#if effectiveStatus}
      <div
        bind:this={statusElement}
        class="card-status {statusClasses[effectiveStatusType]}"
        class:has-alert={alert.length > 0}
        role={alerts.length > 0 ? "tooltip" : undefined}
      >
        <span class="status-text">{effectiveStatus}</span>
        {#if alerts.length > 0}
          <div
            onmouseenter={(e) => handleMouseEnter(alerts, e.currentTarget)}
            onmouseleave={handleMouseLeave}
            role="tooltip"
            aria-label={`Alert: ${alerts.map((alert) => alert.message).join(", ")}`}
          >
            <span class="alert-icon">{@html getAlertIcon(alerts[0].level)}</span
            >
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <div class="card-content">
    {@render children?.()}
  </div>
</div>

<style>
  .card {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 16px;
    height: 100%;
    /* padding: 1.25rem; */
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    /* overflow: hidden; */
  }

  .card::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.4),
      transparent
    );
  }

  .card:hover {
    background: rgba(255, 255, 255, 0.15);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-color: rgba(255, 255, 255, 0.3);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  }

  .large-card {
    grid-column: 1 / -1;
    max-width: 100%;
  }

  .card-header {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
  }

  .card-inner {
    padding-top: 1.25rem;
    padding-left: 1.25rem;
    display: flex;
    align-items: center;
  }

  .card-icon {
    width: 36px;
    height: 36px;
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.2),
      rgba(255, 255, 255, 0.1)
    );
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 0.75rem;
    color: rgba(255, 255, 255, 0.9);
  }

  .icon-emoji {
    font-size: 1.2rem;
  }

  .card-title {
    color: rgba(255, 255, 255, 0.9);
    font-size: 1rem;
    font-weight: 600;
  }

  .card-status {
    margin-left: auto;
    margin-top: 1.25rem;
    padding: 0.2rem 0.6rem;
    border-radius: 4rem;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    font-size: 0.7rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 0.3rem;
    cursor: default;
    transition: all 0.2s ease;
  }

  .card-status.has-alert {
    cursor: pointer;
  }

  .card-status.has-alert:hover {
    transform: scale(1.05);
  }

  .alert-icon {
    display: flex;
    align-items: center;
  }

  .status-text {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 120px;
  }

  .tooltip-portal {
    pointer-events: none;
  }

  .tooltip-content {
    background: rgba(17, 24, 39, 0.95);
    color: white;
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 11px;
    max-width: 200px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    backdrop-filter: blur(8px);
    white-space: nowrap;
  }

  .alert-message {
    font-weight: 500;
    margin-bottom: 2px;
  }

  .alert-time {
    font-size: 10px;
    opacity: 0.7;
    font-weight: 400;
  }

  .status-online {
    background: rgba(52, 211, 153, 0.2);
    color: var(--status-online-color);
  }

  .status-offline {
    background: rgba(248, 113, 113, 0.2);
    color: var(--status-offline-color);
  }

  .status-warning {
    background: rgba(251, 191, 36, 0.2);
    color: var(--status-warning-color);
  }

  .card-content {
    padding: 0 1.25rem 1.25rem 1.25rem;
    /* Content styling will be handled by slot content */
  }
</style>
