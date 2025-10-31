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

  const stats = $derived(data?.data || {});
  const isSuccess = $derived(data?.success || false);
  const error = $derived(data?.error);
  const statusType = $derived(isSuccess ? "online" : "offline");
  const status = $derived(isSuccess ? "Online" : "Offline");

  const title = $derived(config?.title || "Prowlarr");

  const successRate = $derived(() => {
    const total = stats.totalQueries || 0;
    const failed = stats.totalFailedQueries || 0;
    const successful = total - failed;

    if (total === 0) return 0;
    return Math.round((successful / total) * 100);
  });
</script>

<Card
  {title}
  {status}
  {statusType}
  icon={plugin.metadata.icon}
  alerts={stats.alerts as PluginAlert[]}
  href={config.serverUrl}
>
  <div class="prowlarr-widget">
    {#if !isSuccess || error}
      <div class="error-state">
        <div class="error-icon">⚠️</div>
        <div class="error-message">{error}</div>
      </div>
    {:else}
      <StatsGrid columns={3}>
        <Stat
          label="Failed Queries"
          value={formatNumber(stats.totalFailedQueries || 0)}
        />
        <Stat
          label="Total Queries"
          value={formatNumber(stats.totalQueries || 0)}
        />
        <Stat label="Total Grabs" value={formatNumber(stats.totalGrabs || 0)} />
      </StatsGrid>
    {/if}
  </div>
</Card>

<style>
  .prowlarr-widget {
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
</style>
