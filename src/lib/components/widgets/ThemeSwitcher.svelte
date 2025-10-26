<script lang="ts">
  import { currentTheme, themes, type Theme } from "../../stores/theme.js";
  import { onMount } from "svelte";

  let isOpen = $state(false);
  let switcherElement: HTMLElement = $state()!;

  function toggleSwitcher() {
    isOpen = !isOpen;
  }

  function setTheme(themeName: Theme) {
    currentTheme.set(themeName);
    setTimeout(() => {
      isOpen = false;
    }, 500);
  }

  function handleClickOutside(event: MouseEvent) {
    if (switcherElement && !switcherElement.contains(event.target as Node)) {
      isOpen = false;
    }
  }

  onMount(() => {
    currentTheme.init();
    document.addEventListener("click", handleClickOutside);

    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  });
</script>

<div class="theme-switcher-container" bind:this={switcherElement}>
  <button
    class="theme-toggle"
    onclick={toggleSwitcher}
    aria-label="Toggle theme switcher"
  >
    <svg
      class="theme-icon"
      version="1.1"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      viewBox="0 0 21.8524 20.1726"
    >
      <g>
        <rect height="20.1726" opacity="0" width="21.8524" x="0" y="0" />
        <path
          d="M0.301179 17.8255C2.11759 20.2767 5.92618 20.9408 8.28946 19.183C10.3109 17.6888 10.7406 15.3353 9.24649 13.3627C7.85001 11.4877 5.77969 11.0287 4.1586 12.2494C2.63516 13.392 3.17227 15.0326 2.32266 15.6673C1.58048 16.224 0.887117 15.9115 0.350007 16.3216C-0.0406176 16.6341-0.167571 17.2005 0.301179 17.8255ZM2.29337 17.6693C2.16641 17.5326 2.17618 17.3959 2.26407 17.308C2.40079 17.181 2.94766 17.142 3.36759 16.722C4.31485 15.7845 3.99259 14.4076 5.08634 13.558C5.99454 12.8353 7.2543 13.099 8.10391 14.1927C9.06094 15.433 8.72891 16.9466 7.41055 17.9623C5.93594 19.1244 3.66055 19.1146 2.29337 17.6693ZM9.85196 15.5013C10.9066 15.3841 11.7953 14.8568 12.7914 13.8607C15.9262 10.7357 20.7113 3.73375 21.1313 3.14781C22.3617 1.43883 20.2524-0.651017 18.5336 0.550155C17.9574 0.960311 10.9457 5.7357 7.82071 8.89C6.84415 9.87633 6.30704 10.7552 6.18009 11.7904L7.71329 12.2103C7.73282 11.4974 8.11368 10.8041 8.92423 10.0033C12.0199 6.94664 18.8656 2.24937 19.2856 1.93687C19.6176 1.69273 19.9887 2.0443 19.7348 2.40562C19.4809 2.77672 14.725 9.71031 11.6781 12.7572C10.8774 13.558 10.2621 13.8802 9.58829 13.9584ZM11.5805 13.0502L12.8305 12.7279C12.6938 10.8627 10.8676 9.00719 9.03165 8.85094L8.61173 10.1009C9.86173 10.0716 11.5805 11.7904 11.5805 13.0502Z"
          fill="currentColor"
          fill-opacity="0.85"
        />
      </g>
    </svg>
  </button>

  <div class="theme-switcher" class:open={isOpen}>
    <h3>Choose Theme</h3>
    <div class="theme-options">
      {#each Object.entries(themes) as [key, theme]}
        <button
          class="theme-option"
          class:active={$currentTheme === key}
          onclick={() => setTheme(key as Theme)}
        >
          <div class="theme-preview {key}"></div>
          <span class="theme-name">{theme.displayName}</span>
        </button>
      {/each}
    </div>
  </div>
</div>

<style>
  .theme-switcher-container {
    position: fixed;
    top: 2rem;
    right: 2rem;
    z-index: 1000;
  }

  .theme-toggle {
    position: relative;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    border: none;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(20px);
    color: rgba(255, 255, 255, 0.9);
    font-size: 1.5rem;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .theme-toggle:hover {
    background: rgba(255, 255, 255, 0.2);
    transform: scale(1.1);
  }

  .theme-icon {
    width: 24px;
    height: 24px;
    color: rgba(255, 255, 255, 0.85);
    transition: color 0.3s ease;
  }

  .theme-toggle:hover .theme-icon {
    color: rgba(255, 255, 255, 1);
  }

  .theme-switcher {
    position: absolute;
    top: 60px;
    right: 0;
    width: 280px;
    background: rgba(0, 0, 0, 0.9);
    backdrop-filter: blur(20px);
    border-radius: 20px;
    padding: 1.5rem;
    transform: translateY(-20px);
    opacity: 0;
    visibility: hidden;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .theme-switcher.open {
    transform: translateY(0);
    opacity: 1;
    visibility: visible;
  }

  .theme-switcher h3 {
    color: rgba(255, 255, 255, 0.9);
    font-size: 1.1rem;
    margin: 0 0 1rem 0;
    font-weight: 600;
  }

  .theme-options {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .theme-option {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.9rem;
  }

  .theme-option:hover {
    background: rgba(255, 255, 255, 0.1);
    transform: translateX(5px);
  }

  .theme-option.active {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
    color: rgba(255, 255, 255, 1);
  }

  .theme-preview {
    width: 32px;
    height: 32px;
    border-radius: 8px;
    border: 2px solid rgba(255, 255, 255, 0.2);
  }

  .theme-preview.default {
    background: linear-gradient(135deg, #1a1a2e 0%, #533483 50%, #7209b7 100%);
  }

  .theme-preview.cyberpunk {
    background: linear-gradient(135deg, #0a0a0a 0%, #2d1b69 50%, #38ef7d 100%);
  }

  .theme-preview.ocean {
    background: linear-gradient(135deg, #0f172a 0%, #1e40af 50%, #06b6d4 100%);
  }

  .theme-preview.sunset {
    background: linear-gradient(135deg, #1a0404 0%, #dc2626 50%, #fbbf24 100%);
  }

  .theme-preview.darkgold {
    background: linear-gradient(135deg, #000000 0%, #b8860b 50%, #ffd700 100%);
  }

  .theme-preview.image {
    background: linear-gradient(135deg, #1a1a2e 0%, #2d1b69 50%, #38ef7d 100%);
    position: relative;
  }

  .theme-preview.image::after {
    content: "🖼️";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 12px;
  }

  .theme-preview.video {
    background: linear-gradient(135deg, #1a1a2e 0%, #0f172a 50%, #2d1b69 100%);
    position: relative;
  }

  .theme-preview.video::after {
    content: "🎬";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 12px;
  }

  .theme-name {
    font-weight: 500;
  }

  @media (max-width: 768px) {
    .theme-switcher-container {
      top: 1rem;
      right: 1rem;
    }

    .theme-switcher {
      width: 250px;
      right: -50px;
    }
  }
</style>
