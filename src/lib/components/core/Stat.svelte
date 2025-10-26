<script lang="ts">
  interface Props {
    value: string | number;
    label: string;
    withBackground?: boolean;
    size?: "small" | "medium" | "large";
  }

  const {
    value,
    label,
    withBackground = true,
    size = "small",
  }: Props = $props();

  const sizeClasses = {
    small: "stat-small",
    medium: "stat-medium",
    large: "stat-large",
  };
</script>

<div
  class="stat-item {sizeClasses[size]} {withBackground
    ? 'stat-background'
    : ''}"
>
  <div class="stat-value">{value}</div>
  <div class="stat-label">{label}</div>
</div>

<style>
  .stat-item {
    text-align: center;
    padding: 0.5rem 0 0.5rem 0;
  }

  .stat-background {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 8px;
  }

  .stat-value {
    font-weight: 700;
    color: rgba(255, 255, 255, 0.9);
    margin-bottom: 0.25rem;
  }

  .stat-label {
    color: rgba(255, 255, 255, 0.6);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stat-small .stat-value {
    font-size: 0.75rem;
  }

  .stat-small .stat-label {
    font-size: 0.65rem;
  }

  .stat-medium .stat-value {
    font-size: 1.4rem;
  }

  .stat-medium .stat-label {
    font-size: 0.75rem;
  }

  .stat-large .stat-value {
    font-size: 2.4rem;
  }

  .stat-large .stat-label {
    font-size: 1rem;
  }
</style>
