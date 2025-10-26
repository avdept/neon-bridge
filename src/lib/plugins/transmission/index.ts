import type { Plugin, PluginConfig, PluginData } from '../types.js';
import TransmissionWidget from './TransmissionWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'transmission',
    name: 'Transmission',
    description: 'Monitor Transmission BitTorrent client status and torrents',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    category: 'media',
    icon: 'transmission'
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Card Title',
        type: 'text',
        required: false,
        default: 'Transmission',
        description: 'The title displayed on the card',
        placeholder: 'Enter card title'
      },
      {
        key: 'serverUrl',
        label: 'Server URL',
        type: 'url',
        required: true,
        description: 'URL of your Transmission server (including port)',
        placeholder: 'http://192.168.1.100:9091'
      },
      {
        key: 'username',
        label: 'Username',
        type: 'text',
        required: false,
        description: 'Username for Transmission authentication',
        placeholder: 'transmission'
      },
      {
        key: 'password',
        label: 'Password',
        type: 'password',
        required: false,
        credential: true,
        description: 'Password for Transmission authentication'
      },
      {
        key: 'rpcPath',
        label: 'RPC Path',
        type: 'text',
        required: false,
        default: '/transmission/rpc',
        description: 'RPC endpoint path',
        placeholder: '/transmission/rpc'
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

  component: TransmissionWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number, test?: boolean): Promise<PluginData> {
    try {
      const result = await handleApiCall(async () => {
        if (widgetId === undefined || test) {
          const apiUrl = `http://localhost:8080/api/v1/transmission/test`;
          return fetch(apiUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
          });
        } else {
          const apiUrl = `http://localhost:8080/api/v1/transmission/${widgetId}`;
          return fetch(apiUrl, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
        }
      }, 'Transmission');

      if (!result.success) {
        throw new Error(result.error || 'Unknown error occurred');
      }

      return {
        success: true,
        data: result.data,
        error: null
      };
    } catch (error) {
      console.error('Error fetching Transmission data:', error);
      return {
        success: false,
        data: null,
        error: error instanceof Error ? error.message : 'Unknown error occurred'
      };
    }
  }
};
