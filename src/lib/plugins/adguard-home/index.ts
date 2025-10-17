import type { Plugin, PluginConfig, PluginData } from '../types.js';
import AdGuardHomeWidget from './AdGuardHomeWidget.svelte';
import { handleApiCall } from '../../utils/errors.js';

export const plugin: Plugin = {
  metadata: {
    id: 'adguard-home',
    name: 'AdGuard Home',
    description: 'Monitor DNS queries, blocked requests, and processing times from AdGuard Home',
    version: '1.0.0',
    author: 'Alex <https://x.com/_avdept>',
    icon: 'adguard-home',
    category: 'network',
  },

  configTemplate: {
    fields: [
      {
        key: 'title',
        label: 'Widget Title',
        type: 'text',
        required: false,
        default: 'AdGuard Home',
        credential: false,
        placeholder: 'Custom title for this widget'
      },
      {
        key: 'serverUrl',
        label: 'AdGuard Home Server URL',
        type: 'url',
        required: true,
        credential: false,
        placeholder: 'http://192.168.1.100:3000',
        description: 'The base URL of your AdGuard Home instance'
      },
      {
        key: 'username',
        label: 'Username',
        type: 'text',
        required: true,
        placeholder: 'admin',
        credential: true,
        description: 'AdGuard Home admin username'
      },
      {
        key: 'password',
        label: 'Password',
        type: 'password',
        required: true,
        credential: true,
        placeholder: 'password',
        description: 'AdGuard Home admin password'
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

  component: AdGuardHomeWidget,

  async fetchData(config: PluginConfig, widgetId?: string | number): Promise<PluginData> {
    const { serverUrl, username, password } = config;

    if (!serverUrl || !username) {
      throw new Error('Server URL, username, and password are required');
    }

    try {
      const data = await handleApiCall(async () => {
        if (widgetId === undefined) {
          const apiUrl = `http://localhost:8080/api/v1/adguard/test`;
          return fetch(apiUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
          });
        } else {
          const apiUrl = `http://localhost:8080/api/v1/adguard/${widgetId}`;
          return fetch(apiUrl, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
        }
      }, 'AdGuard Home');

      const stats = {
        totalQueries: data.num_dns_queries || 0,
        blockedQueries: data.num_blocked_filtering || 0,
        avgProcessingTime: data.avg_processing_time || 0,
        health: data.health
      };

      return {
        success: true,
        data: stats,
        lastUpdated: new Date().toISOString()
      };

    } catch (error) {
      console.error('AdGuard Home plugin error:', error);
      throw error;
    }
  }
};
