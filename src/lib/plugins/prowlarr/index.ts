import type { Plugin, PluginConfig, PluginData } from '../types.js';
import ProwlarrWidget from './ProwlarrWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'prowlarr',
    name: 'Prowlarr',
    description: 'Monitor indexer statistics, query performance, and system health from Prowlarr',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    icon: 'prowlarr',
    category: 'media',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'Prowlarr',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'Prowlarr Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:9696',
        description: 'The base URL of your Prowlarr instance'
      },
      {
        key: 'apiKey',
        label: 'API Key',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'your-prowlarr-api-key',
        description: 'Prowlarr API key (found in Settings > General > Security)'
      },
      {
        key: 'refreshRate',
        label: 'Refresh Rate (seconds)',
        type: 'number',
        required: false,
        credential: false,
        default: 30,
        description: 'How often to refresh the statistics (10-300 seconds)'
      }
    ]
  },

  component: ProwlarrWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    const data = await handleApiCall(async () => {
      if (widgetId === undefined || test) {
        const apiUrl = `http://localhost:8080/api/v1/prowlarr/test`;
        return fetch(apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(config)
        });
      } else {
        const apiUrl = `http://localhost:8080/api/v1/prowlarr/${widgetId}`;
        return fetch(apiUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });
      }
    }, 'Prowlarr');

    return {
      success: true,
      data: data,
      lastUpdated: new Date().toISOString()
    };
  }
};
