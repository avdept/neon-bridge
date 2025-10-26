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
    }
  }

  return {
    subscribe,
    startUpdates: () => {
      fetchSystemStats();

      updateInterval = setInterval(() => {
        fetchSystemStats();
      }, 5000); // this is hardcoded to 5 seconds for now, only for system stats. TODO: Make it configurable
    },
    stopUpdates: () => {
      if (updateInterval) {
        clearInterval(updateInterval);
      }
    },
    refresh: fetchSystemStats
  };
}

function createServicesStore() {
  const { subscribe, set, update } = writable<ServiceStatus[]>([]);

  let pluginSubscription: (() => void) | null = null;
  let pluginTimers: Map<number, number> = new Map();

  const fetchPluginData = async (instance: PluginInstance) => {
    const plugin = pluginRegistry.get(instance.pluginId);
    if (!plugin?.fetchData) {
      return {
        status: 'online',
        stats: {},
        success: true
      };
    }

    try {
      const data = await plugin.fetchData(instance.config, instance.id);

      if (data && typeof data === 'object' && data.success === false) {
        console.warn(`Plugin ${instance.pluginId} returned error:`, data.error);
        return {
          status: 'offline',
          stats: {},
          error: data.error || 'Plugin returned error',
          success: false
        };
      }

      return {
        ...data,
        success: data?.success !== false
      };
    } catch (error) {
      console.warn(`Failed to fetch data for plugin ${instance.pluginId}:`, error);
      return {
        status: 'offline',
        stats: {},
        error: error instanceof Error ? error.message : 'Unknown error',
        success: false
      };
    }
  };

  const pluginInstanceToService = async (instance: PluginInstance): Promise<ServiceStatus | null> => {
    const plugin = pluginRegistry.get(instance.pluginId);
    const pluginData = await fetchPluginData(instance);

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
        update(services => services.filter(s => s.id !== instance.id));
        console.log(`Removed failed plugin ${instance.pluginId} from dashboard`);
      }
    } catch (error) {
      console.error(`Failed to update plugin ${instance.pluginId}:`, error);
    }
  };

  const setupPluginTimer = (instance: PluginInstance) => {
    const existingTimer = pluginTimers.get(instance.id);
    if (existingTimer) {
      clearInterval(existingTimer);
    }

    const intervalSeconds = getPluginRefreshRate(instance.config);
    const intervalMs = intervalSeconds * 1000;

    updateSinglePlugin(instance);

    const timer = setInterval(() => {
      updateSinglePlugin(instance);
    }, intervalMs);

    pluginTimers.set(instance.id, timer);
    console.log(`Set up timer for ${instance.pluginId}: ${intervalSeconds}s interval`);
  };

  const updateServicesFromPlugins = async () => {
    try {
      let currentInstances: PluginInstance[] = [];
      pluginInstancesFromDB.subscribe(instances => {
        currentInstances = instances.filter(instance => instance.enabled);
      })();

      pluginTimers.forEach(timer => clearInterval(timer));
      pluginTimers.clear();

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
      pluginSubscription = pluginInstancesFromDB.subscribe((instances) => {
        updateServicesFromPlugins();
      });
    },
    stopUpdates: () => {
      pluginTimers.forEach(timer => clearInterval(timer));
      pluginTimers.clear();

      if (pluginSubscription) {
        pluginSubscription();
        pluginSubscription = null;
      }
    },
    refresh: () => {
      updateServicesFromPlugins();
    }
  };
}

export const systemStats = createSystemStatsStore();
export const services = createServicesStore();
