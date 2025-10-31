import type { Plugin, PluginConfig, PluginData } from '../types.js';
import ImmichWidget from './ImmichWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'immich',
    name: 'Immich',
    description: 'Monitor photos, videos, users, and storage from your Immich instance',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    icon: 'immich',
    category: 'media',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'Immich',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'Immich Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:2283',
        description: 'The base URL of your Immich instance'
      },
      {
        key: 'apiKey',
        label: 'API Key',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'your-immich-api-key',
        description: 'Immich API key (found in Account Settings > API Keys)'
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
        key: 'showStorage',
        label: 'Show Storage Usage',
        type: 'boolean',
        required: false,
        credential: false,
        default: true,
        description: 'Display storage usage progress bar'
      },
      {
        key: 'showStorageThreshold',
        label: 'Show storage if usage more than %',
        type: 'number',
        required: false,
        credential: false,
        placeholder: 'e.g. 80',
        description: 'Show storage bar only if usage is above this percentage (leave empty to always show)'
      }
    ]
  },

  component: ImmichWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    const data = await handleApiCall(async () => {
      if (widgetId === undefined || test) {
        const apiUrl = `http://localhost:8080/api/v1/immich/test`;
        return fetch(apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(config)
        });
      } else {
        const apiUrl = `http://localhost:8080/api/v1/immich/${widgetId}`;
        return fetch(apiUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          }
        });
      }
    }, 'Immich');

    return {
      success: true,
      data: data,
      lastUpdated: new Date().toISOString()
    };
  }
};
