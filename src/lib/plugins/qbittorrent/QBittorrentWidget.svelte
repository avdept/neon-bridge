<script lang="ts">
  import Card from "../../components/core/Card.svelte";
  import Stat from "../../components/core/Stat.svelte";
  import StatsGrid from "../../components/core/StatsGrid.svelte";
  import InlineProgressBar from "../../components/core/InlineProgressBar.svelte";
  import { type Plugin, type PluginAlert } from "../../plugins/types.js";
  import { formatNumber, formatBytes } from "../../utils/formatters.js";

  interface Props {
    config: any;
    data: any;
    plugin: Plugin;
    alerts?: PluginAlert[];
  }

  const { config, data, plugin, alerts }: Props = $props();

  const stats = $derived(data?.data || {});
  const isSuccess = $derived(data?.success || false);
  const error = $derived(data?.error);
  const statusType = $derived(isSuccess ? "online" : "offline");
  const status = $derived(isSuccess ? "Online" : "Offline");

  const title = $derived(config?.title || "qBittorrent");

  const errorAlert = $derived(() => {
    if (!isSuccess) {
      return alerts;
    }

    const errorAlerts = [];

    // TODO: maybe one day add stalled torrents as separate type of alerts
    if (stats.errorTorrents && stats.errorTorrents > 0) {
      errorAlerts.push({
        level: "warning" as const,
        message: `${stats.errorTorrents} torrent${stats.errorTorrents === 1 ? "" : "s"} with errors`,
      });
    }

    return errorAlerts.length > 0 ? errorAlerts : alerts;
  });

  const formatSpeed = (bytesPerSecond: number) => {
    if (!bytesPerSecond || bytesPerSecond === 0) return "0 B/s";
    return `${formatBytes(bytesPerSecond)}/s`;
  };

  const downloadPercentage = $derived(() => {
    if (!config?.maxDownloadSpeed || !stats.downloadSpeed) return null;
    const maxBytes = config.maxDownloadSpeed * 1024; // Convert KB/s to B/s
    return Math.min(100, Math.round((stats.downloadSpeed / maxBytes) * 100));
  });

  const uploadPercentage = $derived(() => {
    if (!config?.maxUploadSpeed || !stats.uploadSpeed) return null;
    const maxBytes = config.maxUploadSpeed * 1024; // Convert KB/s to B/s
    return Math.min(100, Math.round((stats.uploadSpeed / maxBytes) * 100));
  });

  const downloadSpeedDisplay = $derived(() => {
    return formatSpeed(stats.downloadSpeed || 0);
  });

  const uploadSpeedDisplay = $derived(() => {
    return formatSpeed(stats.uploadSpeed || 0);
  });
</script>

<Card
  {title}
  {status}
  {statusType}
  icon={plugin.metadata.icon}
  alerts={errorAlert()}
  href={config.serverUrl}
>
  <div class="qbittorrent-widget">
    {#if !isSuccess}
      <div class="error-state">
        <div class="error-icon">⚠️</div>
        <div class="error-message">{error || 'Connection failed'}</div>
      </div>
    {:else}
      <!-- Main Stats Grid -->
      <StatsGrid columns={3}>
        <Stat
          label="Downloading"
          value={formatNumber(stats.downloadingTorrents || 0)}
        />
        <Stat
          label="Seeding"
          value={formatNumber(stats.seedingTorrents || 0)}
        />
        <Stat label="Total" value={formatNumber(stats.totalTorrents || 0)} />
      </StatsGrid>

      <div class="speeds">
        <InlineProgressBar
          title="DL Speed"
          value={downloadPercentage() ?? 0}
          max={100}
          height="18px"
          status={downloadSpeedDisplay()}
          color="primary"
          showPercentage={false}
        />

        <InlineProgressBar
          title="UL Speed"
          value={uploadPercentage() ?? 0}
          max={100}
          height="18px"
          status={uploadSpeedDisplay()}
          color="primary"
          showPercentage={false}
        />
      </div>
    {/if}
  </div>
</Card>

<style>
  .qbittorrent-widget {
    width: 100%;
  }

  .speeds {
    display: flex;
    flex-direction: row;
    gap: 1rem;
    margin-top: 1rem;
  }

  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem 1rem;
    text-align: center;
    gap: 1rem;
  }

  .error-icon {
    font-size: 2rem;
    opacity: 0.7;
  }

  .error-message {
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.9rem;
    line-height: 1.4;
  }
</style>
