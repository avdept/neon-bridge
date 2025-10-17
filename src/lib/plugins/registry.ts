import type { Plugin, PluginInstance } from './types.js';
import { writable, derived } from 'svelte/store';

// Plugin registry
const registeredPlugins = new Map<string, Plugin>();

// Plugin instances store
export const pluginInstances = writable<PluginInstance[]>([]);

// Available plugins derived store
export const availablePlugins = writable<Plugin[]>([]);

export const pluginRegistry = {
  // Register a plugin
  register(plugin: Plugin) {
    registeredPlugins.set(plugin.metadata.id, plugin);
    availablePlugins.update(plugins => [...plugins.filter(p => p.metadata.id !== plugin.metadata.id), plugin]);
  },

  // Get a plugin by ID
  get(pluginId: string): Plugin | undefined {
    return registeredPlugins.get(pluginId);
  },

  // Get all registered plugins
  getAll(): Plugin[] {
    return Array.from(registeredPlugins.values());
  },

  // Unregister a plugin
  unregister(pluginId: string) {
    registeredPlugins.delete(pluginId);
    availablePlugins.update(plugins => plugins.filter(p => p.metadata.id !== pluginId));
  },

  // Get plugins by category
  getByCategory(category: string): Plugin[] {
    return Array.from(registeredPlugins.values()).filter(p => p.metadata.category === category);
  }
};

// Plugin instance management
export const pluginManager = {
  // Add a plugin instance
  addInstance(instance: PluginInstance) {
    pluginInstances.update(instances => [...instances, instance]);
  },

  // Remove a plugin instance
  removeInstance(instanceId: number) {
    pluginInstances.update(instances => instances.filter(i => i.id !== instanceId));
  },

  // Update a plugin instance
  updateInstance(instanceId: number, updates: Partial<PluginInstance>) {
    pluginInstances.update(instances =>
      instances.map(instance =>
        instance.id === instanceId ? { ...instance, ...updates } : instance
      )
    );
  },

  // Get sorted and enabled plugin instances
  getActiveInstances: derived(pluginInstances, $instances =>
    $instances
      .filter(instance => instance.enabled)
      .sort((a, b) => a.order - b.order)
  ),

  // Generate unique instance ID
  generateInstanceId(): string {
    return `instance_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
};
