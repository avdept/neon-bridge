<script lang="ts">
  import StatusBadge from "../core/StatusBadge.svelte";
  import Stat from "../core/Stat.svelte";
  import type { SystemStats } from "../../stores/system.ts";

  interface Props {
    title?: string;
    subtitle?: string;
    systemStats: SystemStats;
  }

  const {
    title = "Neon Bridge",
    subtitle = "System Overview & Management",
    systemStats,
  }: Props = $props();

  const overallStatus = $derived(
    (systemStats.cpu.usage > 80 || systemStats.cpu.temperature > 60
      ? "warning"
      : "online") as "online" | "warning"
  );

  const statusMessage = $derived(
    overallStatus === "warning" ? "High Load Detected" : "All Systems Online"
  );
</script>

<header class="header-bar">
  <div class="header-left">
    <h1 class="title">{title}</h1>
  </div>

  <div class="header-stats">
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
