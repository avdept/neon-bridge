import type { Plugin, PluginConfig, PluginData } from '../types.js';
import RadarrWidget from './RadarrWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'radarr',
    name: 'Radarr',
    description: 'Monitor movie downloads, queue status, and system health from Radarr',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    icon: 'radarr',
    category: 'media',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'Radarr',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'Radarr Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:7878',
        description: 'The base URL of your Radarr instance'
      },
      {
        key: 'apiKey',
        label: 'API Key',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'your-radarr-api-key',
        description: 'Radarr API key (found in Settings > General > Security)'
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

  component: RadarrWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    const data = await handleApiCall(async () => {
      if (widgetId === undefined || test) {
        const apiUrl = `http://localhost:8080/api/v1/radarr/test`;
        return fetch(apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(config)
        });
      } else {
        const apiUrl = `http://localhost:8080/api/v1/radarr/${widgetId}`;
        return fetch(apiUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });
      }
    }, 'Radarr');

    return {
      success: true,
      data: data,
      lastUpdated: new Date().toISOString()
    };


  }
};
