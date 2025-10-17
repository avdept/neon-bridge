import type { Plugin, PluginConfig, PluginData } from '../types.js';
import SonarrWidget from './SonarrWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'sonarr',
    name: 'Sonarr',
    description: 'Monitor TV series downloads, queue status, and system health from Sonarr',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    icon: 'sonarr',
    category: 'media',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'Sonarr',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'Sonarr Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:8989',
        description: 'The base URL of your Sonarr instance'
      },
      {
        key: 'apiKey',
        label: 'API Key',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'your-sonarr-api-key',
        description: 'Sonarr API key (found in Settings > General > Security)'
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
        default: 0,
        description: 'Show storage bar only if usage percentage is above this threshold (0-100)'
      }
    ]
  },

  component: SonarrWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number): Promise<PluginData> {
    try {
      const data = await handleApiCall(async () => {
        if (widgetId === undefined) {
          const apiUrl = `http://localhost:8080/api/v1/sonarr/test`;
          return fetch(apiUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
          });
        } else {
          const apiUrl = `http://localhost:8080/api/v1/sonarr/${widgetId}`;
          return fetch(apiUrl, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
        }
      }, 'Sonarr');

      return {
        success: true,
        data: data,
        lastUpdated: new Date().toISOString()
      };

    } catch (error) {
      console.error('Sonarr plugin error:', error);
      throw error;
    }
  }
};
