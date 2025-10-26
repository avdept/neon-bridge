import type { Plugin, PluginConfig, PluginData } from '../types.js';
import LidarrWidget from './LidarrWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'lidarr',
    name: 'Lidarr',
    description: 'Monitor music downloads, queue status, and system health from Lidarr',
    version: '1.0.0',
    author: 'Homepage',
    icon: 'lidarr',
    category: 'media',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'Lidarr',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'Lidarr Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:8686',
        description: 'The base URL of your Lidarr instance'
      },
      {
        key: 'apiKey',
        label: 'API Key',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'your-lidarr-api-key',
        description: 'Lidarr API key (found in Settings > General > Security)'
      },
      {
        key: 'refreshRate',
        label: 'Refresh Rate (seconds)',
        type: 'number',
        required: false,
        credential: false,
        default: 30,
        description: 'How often to refresh the statistics (10-300 seconds)'
      },
      {
        key: 'showSpaceUsage',
        label: 'Show Space Usage',
        type: 'boolean',
        required: false,
        credential: false,
        default: true,
        description: 'Display storage usage progress bar'
      },
      {
        key: 'showUsageThreshold',
        label: 'Show if usage more than %',
        type: 'number',
        required: false,
        credential: false,
        placeholder: 'e.g. 80',
        description: 'Show storage bar only if usage is above this percentage (leave empty to always show)'
      }
    ]
  },

  component: LidarrWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    try {
      const data = await handleApiCall(async () => {
        if (widgetId === undefined || test) {
          const apiUrl = `http://localhost:8080/api/v1/lidarr/test`;
          return fetch(apiUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
          });
        } else {
          const apiUrl = `http://localhost:8080/api/v1/lidarr/${widgetId}`;
          return fetch(apiUrl, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
        }
      }, 'Lidarr');

      return {
        success: true,
        data: data,
        lastUpdated: new Date().toISOString()
      };

    } catch (error) {
      console.error('Lidarr plugin error:', error);
      throw error;
    }
  }
};
