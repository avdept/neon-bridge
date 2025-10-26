<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import StatusBadge from "../core/StatusBadge.svelte";
  import Stat from "../core/Stat.svelte";
  import type { SystemStats } from "../../stores/system.ts";

  const dispatch = createEventDispatcher();

  interface Props {
    title?: string;
    subtitle?: string;
    systemStats: SystemStats;
    editMode?: boolean;
  }

  const {
    title = "Neon Bridge",
    subtitle = "System Overview & Management",
    systemStats,
    editMode = false,
  }: Props = $props();

  const overallStatus = $derived(
    (systemStats.cpu.usage > 80 || systemStats.cpu.temperature > 60
      ? "warning"
      : "online") as "online" | "warning"
  );

  const statusMessage = $derived(
    overallStatus === "warning" ? "High Load Detected" : "All Systems Online"
  );

  function handleStatsClick() {
    if (editMode) {
      dispatch("configureStats");
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (editMode && (event.key === "Enter" || event.key === " ")) {
      handleStatsClick();
    }
  }
</script>

<header class="header-bar">
  <div class="header-left">
    <h1 class="title">{title}</h1>
  </div>

  <div
    class="header-stats {editMode ? 'edit-mode' : ''}"
    onclick={editMode ? handleStatsClick : undefined}
    onkeydown={editMode ? handleKeydown : undefined}
    role={editMode ? "button" : undefined}
    title={editMode ? "Click to configure stats source" : undefined}
  >
    <Stat
      value="{Math.round(systemStats.cpu.usage)}%"
      label="CPU Usage"
      withBackground={false}
      size="small"
    />

    <div class="header-divider"></div>

    <Stat
      value="{systemStats.memory.used.toFixed(1)}GB"
      label="RAM Used"
      withBackground={false}
      size="small"
    />

    <div class="header-divider"></div>

    <Stat
      value="{Math.round(systemStats.cpu.temperature)}°C"
      label="CPU Temp"
      withBackground={false}
      size="small"
    />

    <div class="header-divider"></div>

    <Stat
      value={systemStats.uptime.display}
      label="Uptime"
      size="small"
      withBackground={false}
    />

    {#if editMode}
      <div class="edit-overlay">
        <button
          class="edit-btn"
          onclick={handleStatsClick}
          title="Configure Stats Source"
          aria-label="Configure Stats Source"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
            ></path>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
            ></path>
          </svg>
        </button>
      </div>
    {/if}
  </div>

  <div class="header-status">
    <StatusBadge status={overallStatus}>
      {statusMessage}
    </StatusBadge>
  </div>
</header>

<style>
  .header-bar {
    background: rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 24px;
    padding: 1.5rem 2rem;
    margin-bottom: 3rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 2rem;
    position: relative;
    overflow: hidden;
  }

  .header-bar::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent
    );
  }

  .header-left {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .title {
    font-size: 2.5rem;
    font-weight: 700;
    background: linear-gradient(135deg, #ffffff 0%, #f0f0f0 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 2px 20px rgba(255, 255, 255, 0.1);
    margin: 0;
  }

  .header-stats {
    display: flex;
    gap: 2rem;
    align-items: center;
    position: relative;
  }

  .header-stats.edit-mode {
    cursor: pointer;
    border: 1px solid transparent;
    border-radius: 12px;
    transition: all 0.2s ease;
    padding: 0.5rem;
    margin: -0.5rem;
    background: none;
    border: none;
    color: inherit;
    font: inherit;
    position: relative;
  }

  .header-stats.edit-mode:hover {
    border-color: rgba(245, 158, 11, 0.5);
    transform: translateY(-2px);
    filter: brightness(1.1);
  }

  .header-stats.edit-mode:focus {
    outline: 2px solid rgba(245, 158, 11, 0.5);
    outline-offset: 2px;
  }

  .header-divider {
    width: 1px;
    height: 40px;
    background: linear-gradient(
      180deg,
      transparent,
      rgba(255, 255, 255, 0.2),
      transparent
    );
  }

  .header-status {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .edit-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: all 0.2s ease;
    backdrop-filter: blur(2px);
  }

  .header-stats.edit-mode:hover .edit-overlay {
    opacity: 1;
  }
  .header-stats.edit-mode {
    border: 2px dashed rgba(245, 158, 11, 0.5);
    border-radius: 12px;
  }

  .edit-btn {
    background: rgba(245, 158, 11, 0.9);
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .edit-btn:hover {
    background: rgba(245, 158, 11, 1);
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
  }

  .edit-btn svg {
    width: 16px;
    height: 16px;
  }

  @media (max-width: 768px) {
    .header-bar {
      flex-direction: column;
      text-align: center;
      gap: 1rem;
      padding: 1rem;
    }

    .title {
      font-size: 1.5rem;
    }

    .header-stats {
      justify-content: center;
      gap: 1rem;
    }

    .header-status {
      justify-content: center;
    }
  }
</style>
