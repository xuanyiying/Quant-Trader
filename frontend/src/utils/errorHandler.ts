import { AxiosError } from 'axios';

export interface ApiError {
  message: string;
  status?: number;
  code?: string;
}

/**
 * Extract error message from various error types
 */
export const getErrorMessage = (error: unknown): string => {
  if (error instanceof Error) {
    return error.message;
  }
  
  if (typeof error === 'string') {
    return error;
  }
  
  const axiosError = error as AxiosError<{ error?: string }>;
  if (axiosError.response?.data?.error) {
    return axiosError.response.data.error;
  }
  
  return 'An unexpected error occurred';
};

/**
 * Parse API error response
 */
export const parseApiError = (error: unknown): ApiError => {
  const axiosError = error as AxiosError<{ error?: string }>;
  
  return {
    message: getErrorMessage(error),
    status: axiosError.response?.status,
    code: axiosError.code,
  };
};

/**
 * Check if error is authentication error
 */
export const isAuthError = (error: unknown): boolean => {
  const axiosError = error as AxiosError;
  return axiosError.response?.status === 401;
};

/**
 * Log error to console with context
 */
export const logError = (context: string, error: unknown): void => {
  console.error(`[${context}]`, error);
};
