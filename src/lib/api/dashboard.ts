// API client for dashboard backend
const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface DashboardWidget {
  id: number;
  dashboard_id: number;
  name: string;
  type: string;
  position: number;
  config: Record<string, any>;
  last_state?: Record<string, any>;
  is_enabled: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface Dashboard {
  id: number;
  name: string;
  description?: string;
  created_at?: string;
  updated_at?: string;
  widgets?: DashboardWidget[];
}

class DashboardAPI {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
    }

    const data = await response.json();
    return data.data || data;
  }

  // Dashboard methods
  async getDashboards(): Promise<Dashboard[]> {
    return this.request<Dashboard[]>('/dashboards');
  }

  async getDashboard(id: number): Promise<Dashboard> {
    return this.request<Dashboard>(`/dashboards/${id}`);
  }

  async createDashboard(dashboard: Omit<Dashboard, 'id' | 'created_at' | 'updated_at'>): Promise<Dashboard> {
    return this.request<Dashboard>('/dashboards', {
      method: 'POST',
      body: JSON.stringify(dashboard),
    });
  }

  async updateDashboard(id: number, dashboard: Partial<Dashboard>): Promise<Dashboard> {
    return this.request<Dashboard>(`/dashboards/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dashboard),
    });
  }

  async deleteDashboard(id: number): Promise<void> {
    await this.request<void>(`/dashboards/${id}`, {
      method: 'DELETE',
    });
  }

  // Widget methods
  async getWidgets(dashboardId: number): Promise<DashboardWidget[]> {
    return this.request<DashboardWidget[]>(`/dashboards/${dashboardId}/widgets`);
  }

  async getWidget(id: number): Promise<DashboardWidget> {
    return this.request<DashboardWidget>(`/widgets/${id}`);
  }

  async createWidget(dashboardId: number, widget: Omit<DashboardWidget, 'id' | 'dashboard_id' | 'created_at' | 'updated_at'>): Promise<DashboardWidget> {
    return this.request<DashboardWidget>(`/dashboards/${dashboardId}/widgets`, {
      method: 'POST',
      body: JSON.stringify({
        ...widget,
        dashboard_id: dashboardId,
      }),
    });
  }

  async updateWidget(id: number, widget: Partial<DashboardWidget>): Promise<DashboardWidget> {
    return this.request<DashboardWidget>(`/widgets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(widget),
    });
  }

  async updateWidgetState(id: number, state: Record<string, any>): Promise<DashboardWidget> {
    return this.request<DashboardWidget>(`/widgets/${id}/state`, {
      method: 'PUT',
      body: JSON.stringify({ last_state: state }),
    });
  }

  async deleteWidget(id: number): Promise<void> {
    await this.request<void>(`/widgets/${id}`, {
      method: 'DELETE',
    });
  }
}

export const dashboardAPI = new DashboardAPI();
