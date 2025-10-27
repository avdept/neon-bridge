import type { Plugin, PluginConfig, PluginData } from '../types.js';
import QBittorrentWidget from './QBittorrentWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'qbittorrent',
    name: 'qBittorrent',
    description: 'Monitor qBittorrent client status and torrents',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    category: 'media',
    icon: 'qbittorrent'
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Card Title',
        type: 'text',
        required: false,
        default: 'qBittorrent',
        description: 'The title displayed on the card',
        placeholder: 'Enter card title'
      },
      {
        key: 'serverUrl',
        label: 'Server URL',
        type: 'url',
        required: true,
        description: 'URL of your qBittorrent WebUI (including port)',
        placeholder: 'http://192.168.1.100:8080'
      },
      {
        key: 'username',
        label: 'Username',
        type: 'text',
        required: false,
        description: 'Username for qBittorrent authentication',
        placeholder: 'admin'
      },
      {
        key: 'password',
        label: 'Password',
        type: 'password',
        required: false,
        credential: true,
        description: 'Password for qBittorrent authentication'
      },
      {
        key: 'maxDownloadSpeed',
        label: 'Max Download Speed (KB/s)',
        type: 'number',
        required: false,
        description: 'Optional: Maximum download speed for reference (used for percentage calculations)',
        placeholder: '10000'
      },
      {
        key: 'maxUploadSpeed',
        label: 'Max Upload Speed (KB/s)',
        type: 'number',
        required: false,
        description: 'Optional: Maximum upload speed for reference (used for percentage calculations)',
        placeholder: '1000'
      },
      {
        key: 'refreshRate',
        label: 'Refresh Rate (seconds)',
        type: 'number',
        required: false,
        default: 30,
        description: 'How often to refresh the data (10-300 seconds)'
      }
    ]
  },

  component: QBittorrentWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    const result = await handleApiCall(async () => {
      if (widgetId === undefined || test) {
        const apiUrl = `http://localhost:8080/api/v1/qbittorrent/test`;
        return fetch(apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(config)
        });
      } else {
        const apiUrl = `http://localhost:8080/api/v1/qbittorrent/${widgetId}`;
        return fetch(apiUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });
      }
    }, 'qBittorrent');

    if (!result.success) {
      throw new Error(result.error || 'Unknown error occurred');
    }

    return {
      success: true,
      data: result.data,
      error: null
    };
  }
};
