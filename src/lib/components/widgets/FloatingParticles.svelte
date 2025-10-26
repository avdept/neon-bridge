<script lang="ts">
  import { onMount } from "svelte";

  interface Props {
    count?: number;
    color?: string;
  }

  const { count = 20, color = "rgba(255, 223, 0, 0.2)" }: Props = $props();

  let particlesContainer: HTMLElement = $state()!;

  function createParticles() {
    if (!particlesContainer) return;

    // Clear existing particles
    particlesContainer.innerHTML = "";

    for (let i = 0; i < count; i++) {
      const particle = document.createElement("div");
      particle.className = "particle";
      particle.style.left = Math.random() * 100 + "%";
      particle.style.animationDelay = Math.random() * 20 + "s";
      particle.style.animationDuration = 15 + Math.random() * 10 + "s";
      particle.style.background = color;
      particlesContainer.appendChild(particle);
    }
  }

  onMount(() => {
    createParticles();
  });
</script>

<div class="floating-particles" bind:this={particlesContainer}></div>

<style>
  .floating-particles {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    overflow: hidden;
    z-index: 0;
  }

  :global(.particle) {
    position: absolute;
    width: 4px;
    height: 4px;
    border-radius: 50%;
    animation: float 20s infinite linear;
  }

  @keyframes float {
    0% {
      transform: translateY(100vh) translateX(-50px);
      opacity: 0;
    }
    10% {
      opacity: 1;
    }
    90% {
      opacity: 1;
    }
    100% {
      transform: translateY(-100vh) translateX(50px);
      opacity: 0;
    }
  }
</style>
