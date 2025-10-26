<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import HeaderBar from "./lib/components/widgets/HeaderBar.svelte";
  import ThemeSwitcher from "./lib/components/widgets/ThemeSwitcher.svelte";
  import FloatingParticles from "./lib/components/widgets/FloatingParticles.svelte";
  import VideoBackground from "./lib/components/widgets/VideoBackground.svelte";
  import PluginDashboard from "./lib/components/PluginDashboard.svelte";
  import Modal from "./lib/components/core/Modal.svelte";
  import GlancesConfigModal from "./lib/components/GlancesConfigModal.svelte";
  import FloatingActionButton from "./lib/components/core/FloatingActionButton.svelte";
  import AddPluginForm from "./lib/components/AddPluginForm.svelte";
  import { systemStats, services } from "./lib/stores/system.js";
  import { pluginRegistry } from "./lib/plugins/registry.js";
  import {
    dashboardStore,
    isDashboardEditMode,
    currentDashboard,
  } from "./lib/stores/dashboard.js";
  import type { PluginInstance } from "./lib/plugins/types.js";


  let isModalOpen = $state(false);
  let isGlancesConfigOpen = $state(false);
  let editingWidget = $state<{
    instance: PluginInstance;
    plugin: any;
    config: any;
  } | null>(null);
  let newWidgetPosition = $state<number | undefined>(undefined);

  const isEditMode = $derived($isDashboardEditMode);
  const dashboard = $derived($currentDashboard);
  let modalTitle = $derived(
    editingWidget
      ? `Edit ${editingWidget.plugin?.metadata.name || "Widget"}`
      : "Add New Plugin"
  );

  let glancesConfig = $derived(() => {
    if (dashboard?.glances_config) {
      try {
        return JSON.parse(dashboard.glances_config);
      } catch {
        return null;
      }
    }
    return null;
  });

  onMount(async () => {
    systemStats.startUpdates();

    try {
      await dashboardStore.initializeDashboard();
      services.startUpdates();
    } catch (error) {
      console.error("Failed to initialize dashboard store:", error);
      services.startUpdates();
    }
  });

  onDestroy(() => {
    systemStats.stopUpdates();
    services.stopUpdates();
  });

  function openModal(position?: number) {
    editingWidget = null;
    newWidgetPosition = position;
    console.log("Opening modal with position:", position);
    isModalOpen = true;
  }

  function toggleEditMode() {
    isDashboardEditMode.update((mode) => !mode);
  }

  function handleBellClick() {
    window.openUpdatify();
  }

  function handleStatsConfig() {
    isGlancesConfigOpen = true;
  }

  function closeGlancesConfig() {
    isGlancesConfigOpen = false;
  }

  function closeModal() {
    isModalOpen = false;
    editingWidget = null;
    newWidgetPosition = undefined;
  }

  function handleWidgetEdit(event: CustomEvent) {
    const { instance, plugin, config } = event.detail;
    editingWidget = { instance, plugin, config };
    isModalOpen = true;
  }
</script>

<main>
  <VideoBackground />
  <div class="background-overlay"></div>
  <FloatingParticles />
  <ThemeSwitcher />

  <div class="container">
    <HeaderBar
      systemStats={$systemStats}
      editMode={isEditMode}
      on:configureStats={handleStatsConfig}
    />

    <PluginDashboard
      onedit={handleWidgetEdit}
      editMode={isEditMode}
      onaddwidget={openModal}
    />
  </div>

  <div class="fab-group">
    <FloatingActionButton
      icon="edit"
      label={isEditMode ? "Exit Edit Mode" : "Edit Dashboard"}
      onclick={toggleEditMode}
      position="bottom-right"
      class="edit-fab {isEditMode ? 'active' : ''}"
    />
    <FloatingActionButton
      icon="bell"
      label="Updates & Notifications"
      onclick={handleBellClick}
      position="bottom-right"
      class="bell-fab releasenotes"
    />
  </div>

  <Modal
    isOpen={isModalOpen}
    title={modalTitle}
    onclose={closeModal}
    maxWidth="600px"
  >
    <AddPluginForm
      editMode={!!editingWidget}
      initialPluginId={editingWidget?.instance.pluginId}
      initialConfig={editingWidget?.config}
      editingWidget={editingWidget?.instance}
      {closeModal}
      position={newWidgetPosition}
    />
  </Modal>

  <!-- Glances Configuration Modal -->
  <GlancesConfigModal
    isOpen={isGlancesConfigOpen}
    dashboardId={dashboard?.id || 1}
    initialConfig={glancesConfig()}
    onClose={closeGlancesConfig}
  />
</main>

<style>
  :root {
    --bg-gradient: linear-gradient(
      135deg,
      #1a1a2e 0%,
      #16213e 25%,
      #0f3460 50%,
      #533483 75%,
      #7209b7 100%
    );
    --bg-overlay: radial-gradient(
        circle at 85% 15%,
        rgba(255, 193, 7, 0.15) 0%,
        transparent 45%
      ),
      radial-gradient(
        circle at 15% 85%,
        rgba(255, 235, 59, 0.1) 0%,
        transparent 50%
      ),
      radial-gradient(
        circle at 50% 50%,
        rgba(139, 69, 19, 0.05) 0%,
        transparent 60%
      );
  }

  :global([data-theme="cyberpunk"]) {
    --bg-gradient: linear-gradient(
      135deg,
      #0a0a0a 0%,
      #1a0033 25%,
      #2d1b69 50%,
      #11998e 75%,
      #38ef7d 100%
    );
    --bg-overlay: radial-gradient(
        circle at 85% 15%,
        rgba(0, 255, 127, 0.15) 0%,
        transparent 45%
      ),
      radial-gradient(
        circle at 15% 85%,
        rgba(255, 20, 147, 0.1) 0%,
        transparent 50%
      ),
      radial-gradient(
        circle at 50% 50%,
        rgba(0, 255, 255, 0.05) 0%,
        transparent 60%
      );
  }

  :global([data-theme="ocean"]) {
    --bg-gradient: linear-gradient(
      135deg,
      #0f172a 0%,
      #1e3a8a 25%,
      #1e40af 50%,
      #3b82f6 75%,
      #06b6d4 100%
    );
    --bg-overlay: radial-gradient(
        circle at 85% 15%,
        rgba(6, 182, 212, 0.15) 0%,
        transparent 45%
      ),
      radial-gradient(
        circle at 15% 85%,
        rgba(59, 130, 246, 0.1) 0%,
        transparent 50%
      ),
      radial-gradient(
        circle at 50% 50%,
        rgba(30, 64, 175, 0.05) 0%,
        transparent 60%
      );
  }

  :global([data-theme="sunset"]) {
    --bg-gradient: linear-gradient(
      135deg,
      #1a0404 0%,
      #4c1d95 25%,
      #dc2626 50%,
      #f59e0b 75%,
      #fbbf24 100%
    );
    --bg-overlay: radial-gradient(
        circle at 85% 15%,
        rgba(251, 191, 36, 0.15) 0%,
        transparent 45%
      ),
      radial-gradient(
        circle at 15% 85%,
        rgba(239, 68, 68, 0.1) 0%,
        transparent 50%
      ),
      radial-gradient(
        circle at 50% 50%,
        rgba(220, 38, 38, 0.05) 0%,
        transparent 60%
      );
  }

  :global([data-theme="darkgold"]) {
    --bg-gradient: linear-gradient(
      135deg,
      #000000 0%,
      #1a1000 15%,
      #2d1f00 25%,
      #4a3300 35%,
      #b8860b 45%,
      #ffd700 50%,
      #b8860b 55%,
      #4a3300 65%,
      #2d1f00 75%,
      #1a1000 85%,
      #000000 100%
    );
    --bg-overlay: radial-gradient(
        circle at 50% 50%,
        rgba(255, 215, 0, 0.25) 0%,
        transparent 60%
      ),
      radial-gradient(
        circle at 80% 20%,
        rgba(184, 134, 11, 0.15) 0%,
        transparent 50%
      ),
      radial-gradient(
        circle at 20% 80%,
        rgba(218, 165, 32, 0.1) 0%,
        transparent 50%
      );
  }

  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  main {
    position: relative;
    width: 100%;
    min-height: 100vh;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      sans-serif;
    background: var(--bg-fallback, #1a1a2e);
    background-image: var(--bg-image, none), var(--bg-gradient);
    background-size: cover, auto;
    background-position: center, center;
    background-attachment: fixed, scroll;
    min-height: 100vh;
    overflow-x: hidden;
  }

  .background-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: var(--bg-overlay);
    pointer-events: none;
    z-index: 0;
  }

  .container {
    position: relative;
    z-index: 1;
    width: 100%;
    max-width: 1920px;
    margin: 0 auto;
    padding: 2rem;
    min-height: 100vh;
  }

  @media (max-width: 1440px) {
    .container {
      max-width: 1400px;
    }
  }

  @media (max-width: 1366px) {
    .container {
      max-width: 1200px;
    }
  }

  @media (max-width: 1200px) {
    .container {
      max-width: 1100px;
      padding: 1.5rem;
    }
  }

  @media (max-width: 1024px) {
    .container {
      max-width: 960px;
      padding: 1.25rem;
    }
  }

  @media (max-width: 900px) {
    .container {
      max-width: 800px;
      padding: 1rem;
    }
  }

  @media (max-width: 768px) {
    .container {
      max-width: 100%;
      padding: 1rem;
    }
  }

  @media (max-width: 480px) {
    .container {
      padding: 0.75rem;
    }
  }

  /* Ultra-wide and 4K displays */
  @media (min-width: 1921px) {
    .container {
      max-width: 2400px;
    }
  }

  /* Large desktop (1440p) */
  @media (max-width: 1440px) and (min-width: 1367px) {
    .container {
      max-width: 1400px;
    }
  }

  /* Edit Button Styling */
  :global(.edit-fab) {
    background: linear-gradient(135deg, #6b7280, #4b5563) !important;
  }

  :global(.edit-fab:hover) {
    background: linear-gradient(135deg, #4b5563, #374151) !important;
    transform: scale(1.1);
  }

  :global(.edit-fab.active) {
    background: linear-gradient(135deg, #fbbf24, #f59e0b) !important;
    box-shadow: 0 4px 16px rgba(251, 191, 36, 0.3) !important;
  }

  :global(.edit-fab.active:hover) {
    background: linear-gradient(135deg, #f59e0b, #d97706) !important;
    box-shadow: 0 6px 24px rgba(251, 191, 36, 0.4) !important;
  }

  /* FAB Group Styling */
  .fab-group {
    position: fixed;
    bottom: 24px;
    right: 24px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    z-index: 1000;
  }

  /* Override FAB positioning within the group */
  .fab-group :global(.fab) {
    position: relative !important;
    bottom: auto !important;
    right: auto !important;
    top: auto !important;
    left: auto !important;
  }

  /* Bell Button Styling */
  :global(.bell-fab) {
    background: linear-gradient(135deg, #3b82f6, #2563eb) !important;
  }

  :global(.bell-fab:hover) {
    background: linear-gradient(135deg, #2563eb, #1d4ed8) !important;
    transform: scale(1.1);
  }

  /* Mobile responsiveness for FAB group */
  @media (max-width: 768px) {
    .fab-group {
      bottom: 20px;
      right: 20px;
      gap: 10px;
    }
  }
</style>
