import type { Component } from 'svelte';

export interface PluginConfig {
  [key: string]: any;
}

export interface PluginMetadata {
  id: string;
  name: string;
  description: string;
  version: string;
  author: string;
  category: 'system' | 'service' | 'monitoring' | 'media' | 'network' | 'storage' | 'custom';
  icon: string;
}

export interface PluginConfigField {
  key: string;
  label: string;
  type: 'text' | 'number' | 'boolean' | 'select' | 'password' | 'url' | 'email';
  required: boolean;
  default?: any;
  options?: Array<{ value: any; label: string }>;
  description?: string;
  credential?: boolean; // Indicates if the field is a credential (e.g., password, token) and will not be shown or sent to frontend
  placeholder?: string;
}

export interface PluginConfigTemplate {
  fields: PluginConfigField[];
}

export interface PluginData {
  [key: string]: any;
}

export interface PluginAlert {
  level: 'warning' | 'error';
  message: string;
  timestamp?: Date;
}

export interface PluginProps {
  config: PluginConfig;
  data?: PluginData;
  span?: number;
  alert?: PluginAlert;
}

export interface Plugin {
  metadata: PluginMetadata;
  configTemplate: PluginConfigTemplate;
  component: Component<any>;
  fetchData?: (config: PluginConfig, widgetId?: string | number, test?: boolean) => Promise<PluginData>;
  validateConfig?: (config: PluginConfig) => boolean;
}

export interface PluginInstance {
  id: number;
  pluginId: string;
  config: PluginConfig;
  span: number;
  order: number;
  enabled: boolean;
  alert?: PluginAlert;
}

export function getPluginRefreshRate(config: PluginConfig): number {
  const rate = config.refreshRate || 30;

  return Math.max(10, Math.min(300, rate)); // rate is within reasonable bounds (10-300 seconds)
}
