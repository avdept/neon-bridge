<script lang="ts">
  import Card from "../../components/core/Card.svelte";
  import Stat from "../../components/core/Stat.svelte";
  import StatsGrid from "../../components/core/StatsGrid.svelte";
  import InlineProgressBar from "../../components/core/InlineProgressBar.svelte";
  import { type Plugin, type PluginAlert } from "../../plugins/types.js";
  import { formatNumber, getStorageColor } from "../../utils/formatters.js";

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

  const title = $derived(config?.title || "Immich");

  const storagePercentage = $derived(() => {
    if (!stats.storage?.diskUsagePercentage) {
      return 0;
    }
    return stats.storage.diskUsagePercentage;
  });

  const shouldShowStorage = $derived(() => {
    if (config?.showStorage === false) {
      return false;
    }

    if (!stats.storage?.diskUsagePercentage) {
      return false;
    }

    const threshold = config?.showStorageThreshold;
    if (threshold !== undefined && threshold !== null && threshold !== "") {
      const currentUsage = storagePercentage();
      return currentUsage >= Number(threshold);
    }

    return true;
  });

  const storageColor = $derived(() => {
    return getStorageColor(storagePercentage());
  });

  const formattedStorage = $derived(() => {
    return {
      used: stats.storage.diskUse,
      total: stats.storage.diskSize,
    };
  });
</script>

<Card
  href={config.serverUrl}
  {title}
  {status}
  {statusType}
  icon={plugin.metadata.icon}
  alerts={stats.alerts as PluginAlert[]}
>
  <div class="immich-widget">
    {#if isSuccess && stats}
      <StatsGrid columns={3}>
        <Stat label="Users" value={formatNumber(stats.users || 0)} />
        <Stat
          label="Photos"
          value={formatNumber(stats.serverStats?.photos || 0)}
        />
        <Stat
          label="Videos"
          value={formatNumber(stats.serverStats?.videos || 0)}
        />
      </StatsGrid>

      {#if shouldShowStorage()}
        <InlineProgressBar
          title="Storage"
          height="18px"
          value={storagePercentage()}
          max={100}
          status="{formattedStorage().used} / {formattedStorage().total}"
          color={storageColor()}
          showPercentage={false}
        />
      {/if}
    {:else if error}
      <div class="error-message">
        <p>Error: {error}</p>
      </div>
    {:else}
      <div class="loading-message">
        <p>Loading Immich data...</p>
      </div>
    {/if}
  </div>
</Card>

<style>
  .immich-widget {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    display: flex;
    justify-content: center;
    min-height: 80px;
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.875rem;
  }

  .error-message p {
    color: #ef4444;
  }
</style>
