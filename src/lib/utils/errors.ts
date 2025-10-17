/**
 * Handles HTTP error responses with status-specific error messages
 * @param response - The failed HTTP response
 * @param serviceType - The type of service (e.g., 'Sonarr', 'Transmission', etc.)
 * @returns Promise<Error> - A promise that resolves to an Error with a descriptive message
 */
export async function handleApiError(response: Response, serviceType: string): Promise<Error> {
  let errorMessage = `Failed to fetch ${serviceType} stats: ${response.status}`;

  try {
    const errorData = await response.json();
    errorMessage = errorData.error || errorMessage;
  } catch {
    // If we can't parse error JSON, use status-based messages
    switch (response.status) {
      case 400:
        errorMessage = 'Invalid widget ID or configuration';
        break;
      case 401:
        errorMessage = 'Authentication failed - check API key';
        break;
      case 403:
        errorMessage = 'Access forbidden - API key may be invalid';
        break;
      case 404:
        errorMessage = `Widget not found or ${serviceType} API not found`;
        break;
      case 502:
        errorMessage = `Connection refused - ${serviceType} may be offline`;
        break;
      case 408:
        errorMessage = `Connection timeout - ${serviceType} is slow to respond`;
        break;
      case 500:
        errorMessage = `Internal server error in ${serviceType}`;
        break;
      case 503:
        errorMessage = `${serviceType} service unavailable`;
        break;
      default:
        errorMessage = `HTTP ${response.status}: ${response.statusText}`;
    }
  }

  return new Error(errorMessage);
}

/**
 * Generic wrapper for API calls that handles errors consistently
 * @param apiCall - Function that returns a Promise<Response>
 * @param serviceType - The type of service for error messages
 * @returns Promise<any> - The parsed JSON response data
 */
export async function handleApiCall(apiCall: () => Promise<Response>, serviceType: string): Promise<any> {
  try {
    const response = await apiCall();

    if (!response.ok) {
      throw await handleApiError(response, serviceType);
    }

    return await response.json();
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error(`Failed to fetch ${serviceType} statistics`);
  }
}
