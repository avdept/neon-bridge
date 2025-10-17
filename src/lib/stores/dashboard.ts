import { writable, derived, get } from 'svelte/store';
import { dashboardAPI, type Dashboard, type DashboardWidget } from '../api/dashboard.js';
import type { PluginInstance } from '../plugins/types.js';

// Current dashboard store
export const currentDashboard = writable<Dashboard | null>(null);

// Widgets store
export const dashboardWidgets = writable<DashboardWidget[]>([]);

// Loading states
export const isDashboardLoading = writable(false);
export const isWidgetSaving = writable(false);

// Error state
export const dashboardError = writable<string | null>(null);

// Edit mode state
export const isDashboardEditMode = writable(false);

// Convert PluginInstance to DashboardWidget format
function pluginInstanceToWidget(instance: PluginInstance, dashboardId: number): Omit<DashboardWidget, 'id' | 'created_at' | 'updated_at'> {
  return {
    dashboard_id: dashboardId,
    name: instance.config.title || instance.pluginId,
    type: instance.pluginId,
    position: instance.order,
    config: instance.config,
    is_enabled: instance.enabled,
  };
}

// Convert DashboardWidget to PluginInstance format
function widgetToPluginInstance(widget: DashboardWidget): PluginInstance {
  return {
    id: widget.id,
    pluginId: widget.type,
    config: widget.config,
    span: widget.config.span || 1,
    order: widget.position,
    enabled: widget.is_enabled,
  };
}

class DashboardStore {
  // Initialize dashboard - load default or create one
  async initializeDashboard(): Promise<void> {
    try {
      isDashboardLoading.set(true);
      dashboardError.set(null);

      const dashboards = await dashboardAPI.getDashboards();

      let dashboard: Dashboard;
      if (dashboards.length > 0) {
        // Use the first dashboard (or could implement dashboard selection later)
        dashboard = dashboards[0];
      } else {
        // Create default dashboard if none exists
        dashboard = await dashboardAPI.createDashboard({
          name: 'Main Dashboard',
          description: 'Primary monitoring dashboard'
        });
      }

      // Load dashboard with widgets
      const fullDashboard = await dashboardAPI.getDashboard(dashboard.id!);

      currentDashboard.set(fullDashboard);
      dashboardWidgets.set(fullDashboard.widgets || []);

      console.log('Dashboard initialized:', fullDashboard);
    } catch (error) {
      console.error('Failed to initialize dashboard:', error);
      dashboardError.set(error instanceof Error ? error.message : 'Failed to load dashboard');
    } finally {
      isDashboardLoading.set(false);
    }
  }

  // Add a new widget to the dashboard
  async addWidget(pluginInstance: PluginInstance): Promise<DashboardWidget | null> {
    try {
      isWidgetSaving.set(true);
      dashboardError.set(null);

      const dashboard = await this.getCurrentDashboard();
      if (!dashboard) {
        throw new Error('No active dashboard');
      }

      const widgetData = pluginInstanceToWidget(pluginInstance, dashboard.id!);

      const savedWidget = await dashboardAPI.createWidget(dashboard.id!, widgetData);

      // Update local store
      dashboardWidgets.update(widgets => [...widgets, savedWidget]);

      console.log('Widget saved to database:', savedWidget);
      return savedWidget;
    } catch (error) {
      console.error('Failed to save widget:', error);
      dashboardError.set(error instanceof Error ? error.message : 'Failed to save widget');
      return null;
    } finally {
      isWidgetSaving.set(false);
    }
  }

  // Remove a widget from the dashboard
  async removeWidget(widgetId: number): Promise<boolean> {
    try {
      await dashboardAPI.deleteWidget(widgetId);

      // Update local store
      dashboardWidgets.update(widgets => widgets.filter(w => w.id !== widgetId));

      console.log('Widget removed from database:', widgetId);
      return true;
    } catch (error) {
      console.error('Failed to remove widget:', error);
      dashboardError.set(error instanceof Error ? error.message : 'Failed to remove widget');
      return false;
    }
  }

  // Update widget configuration
  async updateWidget(widgetId: number, updates: Partial<DashboardWidget>): Promise<boolean> {
    try {
      const updatedWidget = await dashboardAPI.updateWidget(widgetId, updates);

      // Update local store
      dashboardWidgets.update(widgets =>
        widgets.map(w => w.id === widgetId ? updatedWidget : w)
      );

      console.log('Widget updated in database:', updatedWidget);
      return true;
    } catch (error) {
      console.error('Failed to update widget:', error);
      dashboardError.set(error instanceof Error ? error.message : 'Failed to update widget');
      return false;
    }
  }

  // Update widget state (runtime data)
  async updateWidgetState(widgetId: number, state: Record<string, any>): Promise<boolean> {
    try {
      await dashboardAPI.updateWidgetState(widgetId, state);

      // Update local store
      dashboardWidgets.update(widgets =>
        widgets.map(w => w.id === widgetId ? { ...w, last_state: state } : w)
      );

      return true;
    } catch (error) {
      console.error('Failed to update widget state:', error);
      return false;
    }
  }

  // Get current dashboard or ensure one exists
  private async getCurrentDashboard(): Promise<Dashboard | null> {
    const current = get(currentDashboard);

    if (!current) {
      await this.initializeDashboard();
      return get(currentDashboard);
    }

    return current;
  }

  // Get plugin instances from widgets (for compatibility with existing code)
  getPluginInstances(): PluginInstance[] {
    const widgets = get(dashboardWidgets);

    return widgets
      .filter(widget => widget.is_enabled)
      .map(widgetToPluginInstance)
      .sort((a, b) => a.order - b.order);
  }

  // Reload dashboard data from server
  async refreshDashboard(): Promise<void> {
    const dashboard = await this.getCurrentDashboard();
    if (dashboard?.id) {
      const refreshedDashboard = await dashboardAPI.getDashboard(dashboard.id);
      currentDashboard.set(refreshedDashboard);
      dashboardWidgets.set(refreshedDashboard.widgets || []);
    }
  }
}

export const dashboardStore = new DashboardStore();

// Derived store for plugin instances (for backward compatibility)
export const pluginInstancesFromDB = derived<typeof dashboardWidgets, PluginInstance[]>(
  dashboardWidgets,
  ($widgets) => {
    try {
      if (!$widgets || !Array.isArray($widgets)) {
        return [];
      }

      return $widgets
        .filter(widget => widget.is_enabled)
        .map(widgetToPluginInstance)
        .sort((a, b) => a.order - b.order);
    } catch (error) {
      console.error('Error in pluginInstancesFromDB derived store:', error);
      return [];
    }
  }
);
