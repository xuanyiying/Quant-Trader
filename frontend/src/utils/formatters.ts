/**
 * Format number as currency
 */
export const formatCurrency = (value: number | string, decimals = 2): string => {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  return num.toLocaleString(undefined, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
};

/**
 * Format number as percentage
 */
export const formatPercentage = (value: number, decimals = 2): string => {
  return `${value.toFixed(decimals)}%`;
};

/**
 * Format date to locale string
 */
export const formatDateTime = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date;
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString()}`;
};

/**
 * Format time only
 */
export const formatTime = (date: Date | string): string => {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.toLocaleTimeString();
};

/**
 * Parse float safely
 */
export const safeParseFloat = (value: string | number, defaultValue = 0): number => {
  const parsed = typeof value === 'string' ? parseFloat(value) : value;
  return isNaN(parsed) ? defaultValue : parsed;
};

/**
 * Calculate price change percentage
 */
export const calculatePriceChange = (current: number, previous: number): number => {
  if (previous === 0) return 0;
  return ((current / previous - 1) * 100);
};
