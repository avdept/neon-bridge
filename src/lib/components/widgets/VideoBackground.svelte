<script lang="ts">
  import { currentTheme, themes } from "../../stores/theme.js";
  import { onMount } from "svelte";

  let videoElement: HTMLVideoElement = $state()!;
  let isVisible = $state(false);

  const showVideo = $derived($currentTheme === "video");
  const videoUrl = $derived(themes[$currentTheme].backgroundVideo);

  onMount(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        isVisible = entry.isIntersecting;
        if (videoElement) {
          if (isVisible && showVideo) {
            videoElement.play().catch(console.error);
          } else {
            videoElement.pause();
          }
        }
      },
      { threshold: 0.1 }
    );

    if (videoElement) {
      observer.observe(videoElement);
    }

    return () => {
      observer.disconnect();
    };
  });

  $effect(() => {
    if (videoElement && showVideo) {
      if (isVisible) {
        videoElement.play().catch(console.error);
      }
    } else if (videoElement) {
      videoElement.pause();
    }
  });
</script>

{#if showVideo}
  <video
    bind:this={videoElement}
    class="background-video"
    autoplay
    muted
    loop
    playsinline
    preload="metadata"
  >
    <source src={videoUrl} type="video/mp4" />
    <!-- Fallback for browsers that don't support video -->
    <div class="video-fallback"></div>
  </video>
{/if}

<style>
  .background-video {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    object-fit: cover;
    opacity: 0.8;
  }

  .video-fallback {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: var(--bg-fallback, #1a1a2e);
    z-index: -2;
  }
</style>
