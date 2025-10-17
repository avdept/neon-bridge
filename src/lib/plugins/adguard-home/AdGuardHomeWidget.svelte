<script lang="ts">
  import Card from "../../components/core/Card.svelte";
  import Stat from "../../components/core/Stat.svelte";
  import StatsGrid from "../../components/core/StatsGrid.svelte";
  import { type Plugin, type PluginAlert } from "../../plugins/types.js";
  import { formatNumber } from "../../utils/formatters.js";

  interface Props {
    config: any;
    data: any;
    plugin: Plugin;
    alert?: PluginAlert;
  }

  const { config, data, plugin, alert }: Props = $props();

  // Extract the stats from the plugin data
  const stats = $derived(data?.data || {});
  const isSuccess = $derived(data?.success || false);
  const error = $derived(data?.error);
  const statusType = $derived(isSuccess ? "online" : "offline");
  const status = $derived(isSuccess ? "Online" : "Offline");

  // Format processing time to milliseconds
  const avgProcessingTimeMs = $derived(
    stats.avgProcessingTime
      ? (stats.avgProcessingTime * 1000).toFixed(1)
      : "0.0"
  );

  // Get health status from announcement
  const alerts = $derived(() => {
    if (stats.health) {
      return [
        {
          level: "warning" as "error" | "warning",
          message: stats.health,
        },
      ];
    }
    return [];
  });

  // Get the title from config or default
  const title = $derived(config?.title || "AdGuard Home");
</script>

<Card
  {title}
  {status}
  {statusType}
  icon={plugin.metadata.icon}
  alerts={alerts()}
  href={config.serverUrl}
>
  <div class="adguard-home-widget">
    {#if !isSuccess && error}
      <div class="error-state">
        <div class="error-icon">⚠️</div>
        <div class="error-message">{error}</div>
      </div>
    {:else}
      <!-- Main Stats Grid -->
      <StatsGrid columns={3}>
        <Stat
          label="DNS Queries ({stats.timeUnit || 'total'})"
          value={formatNumber(stats.totalQueries || 0)}
        />
        <Stat
          label="Blocked ({stats.blockingPercentage || 0}%)"
          value={formatNumber(stats.blockedQueries || 0)}
        />
        <Stat label="Avg Response (ms)" value={avgProcessingTimeMs} />
      </StatsGrid>
    {/if}
  </div>
</Card>

<style>
  .adguard-home-widget {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px 20px;
    text-align: center;
  }

  .error-icon {
    font-size: 48px;
    margin-bottom: 12px;
    opacity: 0.7;
  }

  .error-message {
    color: rgba(255, 255, 255, 0.8);
    font-size: 14px;
    line-height: 1.4;
  }

  .health-status {
    background: rgba(0, 255, 127, 0.1);
    border: 1px solid rgba(0, 255, 127, 0.2);
    border-radius: 8px;
    padding: 12px 16px;
  }

  .health-label {
    font-size: 12px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.7);
    margin-bottom: 4px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .health-message {
    font-size: 14px;
    color: rgba(255, 255, 255, 0.9);
    line-height: 1.4;
  }
</style>
