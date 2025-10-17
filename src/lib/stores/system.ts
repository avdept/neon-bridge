import { writable, derived } from 'svelte/store';
import { pluginRegistry } from '../plugins/registry.js';
import { pluginInstancesFromDB } from './dashboard.js';
import type { PluginInstance } from '../plugins/types.js';
import { getPluginRefreshRate } from '../plugins/types.js';

export interface SystemStats {
  cpu: {
    usage: number;
    temperature: number;
  };
  memory: {
    used: number;
    total: number;
    percentage: number;
  };
  uptime: {
    days: number;
    display: string;
  };
  loadAverage: number;
  processes: number;
}

export interface ServiceStatus {
  id: number;
  name: string;
  iconId?: string; // Icon identifier for SVG loading (e.g., 'sonarr', 'radarr')
  status: 'online' | 'offline' | 'warning';
  ping?: number;
  lastUpdate: number;
  stats: Record<string, any>;
  pluginData?: any; // Store the original plugin data for widget consumption
}

function createSystemStatsStore() {
  const { subscribe, set, update } = writable<SystemStats>({
    cpu: {
      usage: 0,
      temperature: 0
    },
    memory: {
      used: 0,
      total: 0,
      percentage: 0
    },
    uptime: {
      days: 0,
      display: '0d'
    },
    loadAverage: 0,
    processes: 0
  });

  let updateInterval: number;

  // Fetch system stats from server
  async function fetchSystemStats(): Promise<void> {
    try {
      const response = await fetch('http://localhost:8080/api/v1/system/stats');
      if (response.ok) {
        const result = await response.json();
        if (result.success && result.data) {
          set(result.data);
        } else {
          console.error('Server returned error:', result);
        }
      } else {
        console.error('Failed to fetch system stats:', response.status, response.statusText);
      }
    } catch (error) {
      console.error('Failed to fetch system stats:', error);
      // Keep the current values instead of resetting to avoid flickering
    }
  }

  return {
    subscribe,
    startUpdates: () => {
      // Fetch immediately
      fetchSystemStats();

      // Then update every 5 seconds
      updateInterval = setInterval(() => {
        fetchSystemStats();
      }, 5000);
    },
    stopUpdates: () => {
      if (updateInterval) {
        clearInterval(updateInterval);
      }
    },
    // Manual refresh method
    refresh: fetchSystemStats
  };
}

function createServicesStore() {
  const { subscribe, set, update } = writable<ServiceStatus[]>([]);

  let pluginSubscription: (() => void) | null = null;
  let pluginTimers: Map<number, number> = new Map(); // Store individual plugin timers

  // Function to fetch data from a plugin instance
  const fetchPluginData = async (instance: PluginInstance) => {
    const plugin = pluginRegistry.get(instance.pluginId);
    if (!plugin?.fetchData) {
      // Return mock data for plugins without fetchData
      return {
        status: 'online',
        stats: {},
        success: true
      };
    }

    try {
      const data = await plugin.fetchData(instance.config, instance.id);

      // Check if the plugin explicitly returned a success indicator
      if (data && typeof data === 'object' && data.success === false) {
        console.warn(`Plugin ${instance.pluginId} returned error:`, data.error);
        return {
          status: 'offline',
          stats: {},
          error: data.error || 'Plugin returned error',
          success: false
        };
      }

      // If data is returned and no explicit success false, consider it successful
      return {
        ...data,
        success: data?.success !== false // Preserve existing success if defined, default to true
      };
    } catch (error) {
      debugger
      console.warn(`Failed to fetch data for plugin ${instance.pluginId}:`, error);
      return {
        status: 'offline',
        stats: {},
        error: error instanceof Error ? error.message : 'Unknown error',
        success: false
      };
    }
  };

  // Function to convert plugin instance to service status
  const pluginInstanceToService = async (instance: PluginInstance): Promise<ServiceStatus | null> => {
    const plugin = pluginRegistry.get(instance.pluginId);
    const pluginData = await fetchPluginData(instance);

    // Only return service if plugin data fetch was successful
    if (pluginData.success === false) {
      console.log(`Excluding plugin ${instance.pluginId} from dashboard due to failed data fetch`);
      return null;
    }

    return {
      id: instance.id,
      name: instance.config.title || plugin?.metadata.name || instance.pluginId,
      status: pluginData.status === 'offline' ? 'offline' :
        pluginData.status === 'warning' ? 'warning' : 'online',
      ping: (pluginData as any).ping || 0, // Plugins may not provide ping data
      lastUpdate: Date.now(),
      stats: pluginData.stats || {},
      pluginData: pluginData // Preserve original plugin data for widget consumption
    };
  };

  // Function to update a single plugin's data
  const updateSinglePlugin = async (instance: PluginInstance) => {
    try {
      const serviceData = await pluginInstanceToService(instance);
      if (serviceData) {
        update(services => {
          const existingIndex = services.findIndex(s => s.id === instance.id);
          if (existingIndex >= 0) {
            services[existingIndex] = serviceData;
          } else {
            services.push(serviceData);
          }
          return [...services];
        });
        console.log(`Updated plugin ${instance.pluginId} (interval: ${getPluginRefreshRate(instance.config)}s)`);
      } else {
        // Remove failed plugin from services
        update(services => services.filter(s => s.id !== instance.id));
        console.log(`Removed failed plugin ${instance.pluginId} from dashboard`);
      }
    } catch (error) {
      console.error(`Failed to update plugin ${instance.pluginId}:`, error);
    }
  };

  // Function to set up timer for a single plugin
  const setupPluginTimer = (instance: PluginInstance) => {
    // Clear existing timer if any
    const existingTimer = pluginTimers.get(instance.id);
    if (existingTimer) {
      clearInterval(existingTimer);
    }

    // Get standardized refresh rate from plugin config
    const intervalSeconds = getPluginRefreshRate(instance.config);
    const intervalMs = intervalSeconds * 1000;

    // Initial update
    updateSinglePlugin(instance);

    // Set up recurring timer
    const timer = setInterval(() => {
      updateSinglePlugin(instance);
    }, intervalMs);

    pluginTimers.set(instance.id, timer);
    console.log(`Set up timer for ${instance.pluginId}: ${intervalSeconds}s interval`);
  };

  // Function to update services from plugin instances
  const updateServicesFromPlugins = async () => {
    try {
      // Get current plugin instances from database
      let currentInstances: PluginInstance[] = [];
      pluginInstancesFromDB.subscribe(instances => {
        currentInstances = instances.filter(instance => instance.enabled);
      })();

      // Clear all existing timers
      pluginTimers.forEach(timer => clearInterval(timer));
      pluginTimers.clear();

      // Set up individual timers for each plugin
      currentInstances.forEach(instance => {
        setupPluginTimer(instance);
      });

      console.log(`Set up individual timers for ${currentInstances.length} plugins`);
    } catch (error) {
      console.error('Failed to update services from plugins:', error);
    }
  };

  return {
    subscribe,
    updateService: (id: number, updates: Partial<ServiceStatus>) => {
      update(services =>
        services.map(service =>
          service.id === id
            ? { ...service, ...updates, lastUpdate: Date.now() }
            : service
        )
      );
    },
    startUpdates: () => {
      // Listen for plugin instance changes and set up timers
      pluginSubscription = pluginInstancesFromDB.subscribe((instances) => {
        // Always update when instances change (including empty -> populated transitions)
        updateServicesFromPlugins();
      });
    },
    stopUpdates: () => {
      // Clear all plugin timers
      pluginTimers.forEach(timer => clearInterval(timer));
      pluginTimers.clear();

      if (pluginSubscription) {
        pluginSubscription();
        pluginSubscription = null;
      }
    },
    // Manual refresh function
    refresh: () => {
      updateServicesFromPlugins();
    }
  };
}

export const systemStats = createSystemStatsStore();
export const services = createServicesStore();
