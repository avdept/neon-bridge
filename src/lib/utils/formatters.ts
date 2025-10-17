/**
 * Common utility functions for formatting data in widgets
 */

/**
 * Format large numbers with unit suffixes (K, M)
 * @param num - The number to format
 * @returns Formatted string with unit suffix
 */
export function formatNumber(num: number | undefined | null): string {
  if (num == null || isNaN(num)) {
    return "0";
  }
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + "M";
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + "K";
  }
  return num.toString();
}

/**
 * Format a number as a percentage with one decimal place
 * @param value - The percentage value to format
 * @returns Formatted percentage string
 */
export function formatPercentage(value: number | undefined | null): string {
  if (value == null || isNaN(value)) {
    return "0";
  }
  return value.toFixed(1);
}

/**
 * Format bytes to human readable format (B, KB, MB, GB, TB)
 * @param bytes - The number of bytes to format
 * @returns Formatted string with unit
 */
export function formatBytes(bytes: number | undefined | null): string {
  if (!bytes || bytes === 0) return "0 B";

  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return (bytes / Math.pow(1024, i)).toFixed(1) + " " + sizes[i];
}

/**
 * Calculate storage percentage from total and free storage
 * @param totalStorage - Total storage in bytes
 * @param freeStorage - Free storage in bytes
 * @returns Percentage used (0-100)
 */
export function calculateStoragePercentage(totalStorage: number | undefined | null, freeStorage: number | undefined | null): number {
  if (totalStorage && freeStorage) {
    const used = totalStorage - freeStorage;
    return (used / totalStorage) * 100;
  }
  return 0;
}

/**
 * Get the appropriate color for storage usage
 * @param percentage - Storage usage percentage (0-100)
 * @returns Color string for UI components
 */
export function getStorageColor(percentage: number): "primary" | "warning" | "danger" {
  if (percentage > 90) return "danger";
  if (percentage > 75) return "warning";
  return "primary";
}
