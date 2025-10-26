import type { Plugin, PluginInstance } from './types.js';
import { writable, derived } from 'svelte/store';

const registeredPlugins = new Map<string, Plugin>();

export const pluginInstances = writable<PluginInstance[]>([]);

export const availablePlugins = writable<Plugin[]>([]);

export const pluginRegistry = {
  register(plugin: Plugin) {
    registeredPlugins.set(plugin.metadata.id, plugin);
    availablePlugins.update(plugins => [...plugins.filter(p => p.metadata.id !== plugin.metadata.id), plugin]);
  },

  get(pluginId: string): Plugin | undefined {
    return registeredPlugins.get(pluginId);
  },

  getAll(): Plugin[] {
    return Array.from(registeredPlugins.values());
  },

  unregister(pluginId: string) {
    registeredPlugins.delete(pluginId);
    availablePlugins.update(plugins => plugins.filter(p => p.metadata.id !== pluginId));
  },

  getByCategory(category: string): Plugin[] {
    return Array.from(registeredPlugins.values()).filter(p => p.metadata.category === category);
  }
};

export const pluginManager = {
  addInstance(instance: PluginInstance) {
    pluginInstances.update(instances => [...instances, instance]);
  },

  removeInstance(instanceId: number) {
    pluginInstances.update(instances => instances.filter(i => i.id !== instanceId));
  },

  updateInstance(instanceId: number, updates: Partial<PluginInstance>) {
    pluginInstances.update(instances =>
      instances.map(instance =>
        instance.id === instanceId ? { ...instance, ...updates } : instance
      )
    );
  },

  getActiveInstances: derived(pluginInstances, $instances =>
    $instances
      .filter(instance => instance.enabled)
      .sort((a, b) => a.order - b.order)
  ),

  generateInstanceId(): string {
    return `instance_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
};
