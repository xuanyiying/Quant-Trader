/**
 * Notification utility functions
 * TODO: Replace with a proper toast notification library like react-hot-toast
 */

export const showSuccess = (message: string): void => {
  alert(message);
};

export const showError = (message: string): void => {
  alert(message);
};

export const showInfo = (message: string): void => {
  alert(message);
};

export const showConfirm = (message: string): boolean => {
  return window.confirm(message);
};
