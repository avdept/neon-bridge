<script lang="ts">
  import Card from "../../components/core/Card.svelte";
  import Stat from "../../components/core/Stat.svelte";
  import StatsGrid from "../../components/core/StatsGrid.svelte";
  import InlineProgressBar from "../../components/core/InlineProgressBar.svelte";
  import { type Plugin, type PluginAlert } from "../../plugins/types.js";
  import {
    formatNumber,
    formatPercentage,
    formatBytes,
    calculateStoragePercentage,
    getStorageColor,
  } from "../../utils/formatters.js";

  interface Props {
    config: any;
    data: any;
    plugin: Plugin;
    alerts?: PluginAlert[];
  }

  const { config, data, plugin, alerts }: Props = $props();

  // Extract the stats from the plugin data
  const stats = $derived(data?.data || {});
  const isSuccess = $derived(data?.success || false);
  const error = $derived(data?.error);
  const statusType = $derived(isSuccess ? "online" : "offline");
  const status = $derived(isSuccess ? "Online" : "Offline");

  // Get the title from config or default
  const title = $derived(config?.title || "Sonarr");

  // Calculate storage percentage if available
  const storagePercentage = $derived(() => {
    return calculateStoragePercentage(stats.totalStorage, stats.freeStorage);
  });

  // Determine if storage bar should be shown
  const shouldShowStorage = $derived(() => {
    // First check if space usage display is enabled
    if (config?.showSpaceUsage === false) {
      return false;
    }

    // Check if we have storage data
    if (!stats.totalStorage || !stats.freeStorage) {
      return false;
    }

    // Check threshold if configured
    const threshold = config?.showUsageThreshold;
    if (threshold !== undefined && threshold !== null && threshold !== "") {
      const currentUsage = storagePercentage();
      return currentUsage >= Number(threshold);
    }

    // Default: show if showSpaceUsage is not explicitly false
    return true;
  });

  // Process health alerts from backend
  const healthAlerts = $derived(() => {
    if (!stats.healthAlerts || !Array.isArray(stats.healthAlerts) || stats.healthAlerts.length === 0) {
      return [];
    }

    // Convert all health alerts to PluginAlert format
    debugger
    const convertedAlerts = stats.healthAlerts
      .filter((h: any) => h.type === 'error' || h.type === 'warning')
      .map((h: any) => ({
        level: h.type as 'error' | 'warning',
        message: h.message as string,
      }));

    return convertedAlerts;
  });
</script>

<Card
  {title}
  {status}
  statusType={statusType}
  icon={plugin.metadata.icon}
  alerts={healthAlerts()}
  href={config.serverUrl}
>
  <div class="sonarr-widget">
    {#if !isSuccess && error}
      <div class="error-state">
        <div class="error-icon">⚠️</div>
        <div class="error-message">{error}</div>
      </div>
    {:else}
      <!-- Main Stats Grid -->
      <StatsGrid columns={3}>
        <Stat
          label="Queued Downloads"
          value={formatNumber(stats.queuedItems || 0)}
        />
        <Stat
          label="Download Progress"
          value={formatPercentage(stats.downloadProgress || 0) + "%"}
        />
        <Stat
          label="Missing"
          value={formatNumber(stats.missingEpisodes || 0)}
        />
      </StatsGrid>

      <!-- Storage Information -->
      {#if shouldShowStorage()}
        <InlineProgressBar
          title="Storage"
          height="18px"
          value={storagePercentage()}
          max={100}
          status="{formatBytes(
            stats.totalStorage - stats.freeStorage
          )} / {formatBytes(stats.totalStorage)}"
          color={getStorageColor(storagePercentage())}
          showPercentage={false}
        />
      {/if}
    {/if}
  </div>
</Card>

<style>
  .sonarr-widget {
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
