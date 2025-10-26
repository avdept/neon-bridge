import { writable } from 'svelte/store';

export type Theme = 'default' | 'darkgold' | 'darkocean' | 'video';

export interface ThemeConfig {
  name: string;
  displayName: string;
  gradient?: string;
  backgroundImage?: string;
  backgroundVideo?: string;
  overlay: string;
  fallbackColor?: string;
}

export const themes: Record<Theme, ThemeConfig> = {
  default: {
    name: 'default',
    displayName: 'Default',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    overlay: `radial-gradient(circle at 30% 70%, rgba(255, 255, 255, 0.2) 0%, transparent 50%),
              radial-gradient(circle at 70% 30%, rgba(255, 255, 255, 0.1) 0%, transparent 50%)`,
    fallbackColor: '#667eea'
  },
  // cyberpunk: {
  //   name: 'cyberpunk',
  //   displayName: 'Cyberpunk',
  //   gradient: 'linear-gradient(135deg, #0f0f0f 0%, #1a0033 25%, #330066 50%, #ff00ff 75%, #00ffff 100%)',
  //   overlay: `radial-gradient(circle at 20% 80%, rgba(255, 0, 255, 0.3) 0%, transparent 50%),
  //             radial-gradient(circle at 80% 20%, rgba(0, 255, 255, 0.2) 0%, transparent 50%)`,
  //   fallbackColor: '#1a0033'
  // },
  // ocean: {
  //   name: 'ocean',
  //   displayName: 'Ocean',
  //   gradient: 'linear-gradient(135deg, #0c4a6e 0%, #0369a1 25%, #0ea5e9 50%, #38bdf8 75%, #7dd3fc 100%)',
  //   overlay: `radial-gradient(circle at 50% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 60%),
  //             radial-gradient(circle at 20% 20%, rgba(56, 189, 248, 0.2) 0%, transparent 50%)`,
  //   fallbackColor: '#0c4a6e'
  // },
  // sunset: {
  //   name: 'sunset',
  //   displayName: 'Sunset',
  //   gradient: 'linear-gradient(135deg, #1e1b4b 0%, #7c2d12 25%, #ea580c 50%, #f97316 75%, #fbbf24 100%)',
  //   overlay: `radial-gradient(circle at 80% 20%, rgba(251, 191, 36, 0.3) 0%, transparent 50%),
  //             radial-gradient(circle at 20% 80%, rgba(249, 115, 22, 0.2) 0%, transparent 60%)`,
  //   fallbackColor: '#1e1b4b'
  // },
  darkgold: {
    name: 'darkgold',
    displayName: 'Dark Gold',
    backgroundImage: 'url("https://4kwallpapers.com/images/wallpapers/macos-ventura-macos-13-macos-2022-stock-dark-mode-5k-retina-3840x2160-8133.jpg")',
    overlay: `radial-gradient(circle at 50% 50%, rgba(255, 215, 0, 0.25) 0%, transparent 60%),
              radial-gradient(circle at 80% 20%, rgba(184, 134, 11, 0.15) 0%, transparent 50%),
              radial-gradient(circle at 20% 80%, rgba(218, 165, 32, 0.1) 0%, transparent 50%)`,
    fallbackColor: '#1a1000'
  },
  darkocean: {
    name: 'darkocean',
    displayName: 'Dark ocean',
    backgroundImage: 'url("https://images.unsplash.com/photo-1620121692029-d088224ddc74?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2072&q=80")',
    overlay: `linear-gradient(135deg, rgba(0, 0, 0, 0.4) 0%, rgba(0, 0, 0, 0.2) 50%, rgba(0, 0, 0, 0.6) 100%),
              radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 70%)`,
    fallbackColor: '#1a1a2e'
  },
  video: {
    name: 'video',
    displayName: 'Video Background',
    backgroundVideo: 'https://cdn.pixabay.com/video/2024/03/14/204117-923594068.mp4?download',
    overlay: `linear-gradient(135deg, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.1) 50%, rgba(0, 0, 0, 0.4) 100%),
              radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.05) 0%, transparent 70%)`,
    fallbackColor: '#1a1a2e'
  }
};

function createThemeStore() {
  const { subscribe, set, update } = writable<Theme>('default');

  return {
    subscribe,
    set: (theme: Theme) => {
      set(theme);
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('selectedTheme', theme);
      }
      applyTheme(theme);
    },
    init: () => {
      if (typeof localStorage !== 'undefined') {
        const saved = localStorage.getItem('selectedTheme') as Theme;
        if (saved && themes[saved]) {
          set(saved);
          applyTheme(saved);
        }
      }
    }
  };
}

function applyTheme(theme: Theme) {
  if (typeof document !== 'undefined') {
    const root = document.documentElement;
    const config = themes[theme];

    if (config.backgroundVideo) {
      // Video theme
      root.style.setProperty('--bg-gradient', 'none');
      root.style.setProperty('--bg-image', 'none');
      root.style.setProperty('--bg-video', config.backgroundVideo);
      root.style.setProperty('--bg-fallback', config.fallbackColor || '');
    } else if (config.backgroundImage) {
      // Image theme
      root.style.setProperty('--bg-gradient', 'none');
      root.style.setProperty('--bg-image', config.backgroundImage);
      root.style.setProperty('--bg-video', 'none');
      root.style.setProperty('--bg-fallback', config.fallbackColor || '');
    } else {
      // Gradient theme
      root.style.setProperty('--bg-gradient', config.gradient || '');
      root.style.setProperty('--bg-image', 'none');
      root.style.setProperty('--bg-video', 'none');
      root.style.setProperty('--bg-fallback', config.fallbackColor || '');
    }

    root.style.setProperty('--bg-overlay', config.overlay);

    if (theme === 'default') {
      root.removeAttribute('data-theme');
    } else {
      root.setAttribute('data-theme', theme);
    }
  }
} export const currentTheme = createThemeStore();
