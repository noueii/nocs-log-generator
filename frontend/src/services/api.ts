/**
 * API client configuration and utilities
 * Configures axios client with interceptors, error handling, and base configuration
 */

import axios, { AxiosError } from 'axios';
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

/**
 * API response wrapper interface
 */
export interface IApiResponse<T = any> {
  data: T;
  message?: string;
  success: boolean;
  timestamp: string;
}

/**
 * API error interface
 */
export interface IApiError {
  message: string;
  code?: string;
  status: number;
  details?: any;
  timestamp: string;
}

/**
 * API client configuration
 */
interface IApiConfig {
  baseURL: string;
  timeout: number;
  retries: number;
  retryDelay: number;
}

/**
 * Default API configuration
 */
const defaultConfig: IApiConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: parseInt(import.meta.env.VITE_API_TIMEOUT) || 30000,
  retries: parseInt(import.meta.env.VITE_API_RETRY_ATTEMPTS) || 3,
  retryDelay: parseInt(import.meta.env.VITE_API_RETRY_DELAY) || 1000,
};

/**
 * Create axios instance with base configuration
 */
const createApiClient = (config: IApiConfig): AxiosInstance => {
  const client = axios.create({
    baseURL: config.baseURL,
    timeout: config.timeout,
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  });

  // Request interceptor
  client.interceptors.request.use(
    (requestConfig: InternalAxiosRequestConfig) => {
      // Add timestamp to requests
      const timestamp = new Date().toISOString();
      
      // Add request ID for tracking
      const requestId = `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      
      if (requestConfig.headers) {
        requestConfig.headers['X-Request-ID'] = requestId;
        requestConfig.headers['X-Request-Timestamp'] = timestamp;
      }

      // Log request in development or when debug logging is enabled
      if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
        console.log(`🚀 API Request [${requestId}]:`, {
          method: requestConfig.method?.toUpperCase(),
          url: requestConfig.url,
          baseURL: requestConfig.baseURL,
          data: requestConfig.data,
          params: requestConfig.params,
        });
      }

      return requestConfig;
    },
    (error: AxiosError) => {
      console.error('❌ Request interceptor error:', error);
      return Promise.reject(error);
    }
  );

  // Response interceptor
  client.interceptors.response.use(
    (response: AxiosResponse) => {
      const requestId = response.config.headers?.['X-Request-ID'];
      
      // Log response in development or when debug logging is enabled
      if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
        console.log(`✅ API Response [${requestId}]:`, {
          status: response.status,
          statusText: response.statusText,
          data: response.data,
          headers: response.headers,
        });
      }

      // Transform response to standard format if needed
      if (response.data && typeof response.data === 'object') {
        // If response already has standard format, return as is
        if ('data' in response.data && 'success' in response.data) {
          return response;
        }
        
        // Otherwise, wrap in standard format
        response.data = {
          data: response.data,
          success: true,
          timestamp: new Date().toISOString(),
        };
      }

      return response;
    },
    async (error: AxiosError) => {
      const requestId = error.config?.headers?.['X-Request-ID'];
      
      // Log error in development or when debug logging is enabled
      if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
        console.error(`❌ API Error [${requestId}]:`, {
          message: error.message,
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data,
          config: {
            method: error.config?.method,
            url: error.config?.url,
            baseURL: error.config?.baseURL,
          },
        });
      }

      // Create standardized error
      const apiError: IApiError = {
        message: error.message || 'An unexpected error occurred',
        status: error.response?.status || 0,
        timestamp: new Date().toISOString(),
      };

      // Extract error details from response
      if (error.response?.data) {
        if (typeof error.response.data === 'string') {
          apiError.message = error.response.data;
        } else if (typeof error.response.data === 'object') {
          const errorData = error.response.data as any;
          apiError.message = errorData.message || errorData.error || apiError.message;
          apiError.code = errorData.code;
          apiError.details = errorData.details;
        }
      }

      // Handle specific error cases
      switch (error.response?.status) {
        case 400:
          apiError.message = apiError.message || 'Bad Request - Please check your input';
          break;
        case 401:
          apiError.message = 'Unauthorized - Please check your credentials';
          break;
        case 403:
          apiError.message = 'Forbidden - You do not have permission to access this resource';
          break;
        case 404:
          apiError.message = 'Not Found - The requested resource could not be found';
          break;
        case 422:
          apiError.message = apiError.message || 'Validation Error - Please check your input';
          break;
        case 429:
          apiError.message = 'Too Many Requests - Please try again later';
          break;
        case 500:
          apiError.message = 'Internal Server Error - Please try again later';
          break;
        case 502:
          apiError.message = 'Bad Gateway - Service is temporarily unavailable';
          break;
        case 503:
          apiError.message = 'Service Unavailable - Please try again later';
          break;
        case 504:
          apiError.message = 'Gateway Timeout - Request timed out';
          break;
        default:
          if (error.code === 'ECONNABORTED') {
            apiError.message = 'Request timeout - Please try again';
          } else if (error.code === 'ERR_NETWORK') {
            apiError.message = 'Network error - Please check your connection';
          }
      }

      // Retry logic for specific errors
      if (shouldRetry(error, config)) {
        return retryRequest(error, config, client);
      }

      return Promise.reject(apiError);
    }
  );

  return client;
};

/**
 * Determine if request should be retried
 */
const shouldRetry = (error: AxiosError, config: IApiConfig): boolean => {
  // Don't retry if we've already exceeded max retries
  const retryCount = (error.config as any)?._retryCount || 0;
  if (retryCount >= config.retries) {
    return false;
  }

  // Only retry on specific error conditions
  const retryableStatuses = [408, 429, 500, 502, 503, 504];
  const retryableCodes = ['ECONNABORTED', 'ERR_NETWORK', 'ENOTFOUND', 'ECONNRESET'];
  
  return (
    (error.response && retryableStatuses.includes(error.response.status)) ||
    (error.code && retryableCodes.includes(error.code)) ||
    error.message.includes('timeout')
  );
};

/**
 * Retry failed request with exponential backoff
 */
const retryRequest = async (
  error: AxiosError, 
  config: IApiConfig,
  client: AxiosInstance
): Promise<AxiosResponse> => {
  const retryCount = ((error.config as any)?._retryCount || 0) + 1;
  const delay = config.retryDelay * Math.pow(2, retryCount - 1); // Exponential backoff
  
  if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
    console.log(`🔄 Retrying request (attempt ${retryCount}/${config.retries}) in ${delay}ms`);
  }

  // Wait before retry
  await new Promise(resolve => setTimeout(resolve, delay));

  // Update retry count
  if (error.config) {
    (error.config as any)._retryCount = retryCount;
    return client.request(error.config);
  }

  throw error;
};

/**
 * Create and export the main API client instance
 */
export const apiClient = createApiClient(defaultConfig);

/**
 * Health check function
 */
export const checkApiHealth = async (): Promise<{ status: string; timestamp: string }> => {
  try {
    await apiClient.get('/health');
    return {
      status: 'healthy',
      timestamp: new Date().toISOString(),
    };
  } catch (error) {
    console.error('API health check failed:', error);
    return {
      status: 'unhealthy',
      timestamp: new Date().toISOString(),
    };
  }
};

/**
 * Get API status information
 */
export const getApiStatus = async (): Promise<any> => {
  const response = await apiClient.get('/api/v1/status');
  return response.data;
};

/**
 * Utility function to handle API calls with consistent error handling
 */
export const handleApiCall = async <T>(
  apiCall: () => Promise<AxiosResponse<IApiResponse<T>>>
): Promise<T> => {
  try {
    const response = await apiCall();
    
    // Check if response has standard format
    if (response.data && 'data' in response.data) {
      return response.data.data;
    }
    
    // Fallback to direct response data
    return response.data as any;
  } catch (error) {
    // Re-throw the standardized error
    throw error;
  }
};

/**
 * Create a custom API client with different configuration
 */
export const createCustomApiClient = (customConfig: Partial<IApiConfig>): AxiosInstance => {
  const config = { ...defaultConfig, ...customConfig };
  return createApiClient(config);
};

/**
 * Request timeout utility
 */
export const withTimeout = <T>(
  promise: Promise<T>,
  timeoutMs: number,
  timeoutMessage = 'Request timed out'
): Promise<T> => {
  const timeoutPromise = new Promise<never>((_, reject) =>
    setTimeout(() => reject(new Error(timeoutMessage)), timeoutMs)
  );

  return Promise.race([promise, timeoutPromise]);
};

/**
 * Generic GET request utility
 */
export const get = async <T>(url: string, params?: any): Promise<T> => {
  return handleApiCall<T>(() => apiClient.get(url, { params }));
};

/**
 * Generic POST request utility
 */
export const post = async <T>(url: string, data?: any): Promise<T> => {
  return handleApiCall<T>(() => apiClient.post(url, data));
};

/**
 * Generic PUT request utility
 */
export const put = async <T>(url: string, data?: any): Promise<T> => {
  return handleApiCall<T>(() => apiClient.put(url, data));
};

/**
 * Generic DELETE request utility
 */
export const del = async <T>(url: string): Promise<T> => {
  return handleApiCall<T>(() => apiClient.delete(url));
};

/**
 * Toast notification interface (to be implemented by the app)
 */
export interface IToastNotification {
  showToast: (message: string, type: 'success' | 'error' | 'info') => void;
}

/**
 * Global toast notification handler (set by the app)
 */
let toastHandler: IToastNotification | null = null;

/**
 * Set global toast handler
 */
export const setToastHandler = (handler: IToastNotification) => {
  toastHandler = handler;
};

/**
 * Show error toast if handler is set
 */
const showErrorToast = (message: string) => {
  if (toastHandler) {
    toastHandler.showToast(message, 'error');
  }
};

/**
 * Show success toast if handler is set
 */
const showSuccessToast = (message: string) => {
  if (toastHandler) {
    toastHandler.showToast(message, 'success');
  }
};

/**
 * Offline detection
 */
export const isOnline = () => navigator.onLine;

/**
 * Network status event listener
 */
let networkStatusListeners: Array<(online: boolean) => void> = [];

window.addEventListener('online', () => {
  networkStatusListeners.forEach(listener => listener(true));
});

window.addEventListener('offline', () => {
  networkStatusListeners.forEach(listener => listener(false));
});

/**
 * Add network status listener
 */
export const addNetworkStatusListener = (listener: (online: boolean) => void) => {
  networkStatusListeners.push(listener);
  return () => {
    networkStatusListeners = networkStatusListeners.filter(l => l !== listener);
  };
};

/**
 * Create API call with offline handling
 */
export const apiCallWithOfflineHandling = async <T>(
  apiCall: () => Promise<T>,
  fallback?: () => Promise<T>
): Promise<T> => {
  try {
    // Check if offline mode is enabled and we're offline
    if (import.meta.env.VITE_ENABLE_OFFLINE_MODE === 'true' && !isOnline()) {
      if (fallback) {
        console.log('📱 Using offline fallback');
        return await fallback();
      }
      throw new Error('Network unavailable and no offline fallback provided');
    }
    
    return await apiCall();
  } catch (error) {
    // If network error and offline mode enabled, try fallback
    if (import.meta.env.VITE_ENABLE_OFFLINE_MODE === 'true' && 
        (error as any)?.code === 'ERR_NETWORK' && fallback) {
      console.log('📱 Network error, using offline fallback');
      return await fallback();
    }
    
    throw error;
  }
};

/**
 * Export default client
 */
export default apiClient;