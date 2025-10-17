// Plugin component initializer - handles dynamic component loading
import { pluginRegistry } from './registry.js';
import type { Plugin } from './types.js';

// Dynamically import all plugin index.ts files from subfolders
const pluginModules = import.meta.glob('./*/index.ts', { eager: true });

export async function initializePlugins() {
  console.log('🚀 Initializing plugin system...');

  const loadedPlugins: Plugin[] = [];

  for (const [path, module] of Object.entries(pluginModules)) {
    try {
      // Extract plugin from module
      const pluginModule = module as { plugin: Plugin };

      if (pluginModule.plugin) {
        const plugin = pluginModule.plugin;

        // Register the plugin
        pluginRegistry.register(plugin);
        loadedPlugins.push(plugin);

        console.log(`✅ Loaded plugin: ${plugin.metadata.name} (${plugin.metadata.id})`);
      } else {
        console.warn(`⚠️  Plugin at ${path} doesn't export a 'plugin' object`);
      }
    } catch (error) {
      console.error(`❌ Failed to load plugin at ${path}:`, error);
    }
  }

  console.log(`✅ Plugin system initialized with ${loadedPlugins.length} plugins`);

  return pluginRegistry.getAll();
}
