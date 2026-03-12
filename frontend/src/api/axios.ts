import axios, { AxiosError } from 'axios';
import type { ErrorResponse } from '../types/errors';
import {
  ValidationError,
  UnauthorizedError,
  ForbiddenError,
  NotFoundError,
  ConflictError,
  InternalError,
} from '../types/errors';

const instance = axios.create({
  baseURL: '/',
});

instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ErrorResponse>) => {
    if (!error.response) {
      return Promise.reject(new InternalError('Network error'));
    }

    const { status, data } = error.response;
    const errorMessage = data?.error || 'An unexpected error occurred';

    switch (status) {
      case 400:
        return Promise.reject(new ValidationError(errorMessage));
      case 401:
        localStorage.removeItem('token');
        return Promise.reject(new UnauthorizedError(errorMessage));
      case 403:
        return Promise.reject(new ForbiddenError(errorMessage));
      case 404:
        return Promise.reject(new NotFoundError(errorMessage));
      case 409:
        return Promise.reject(new ConflictError(errorMessage));
      case 500:
        return Promise.reject(new InternalError(errorMessage));
      default:
        return Promise.reject(new InternalError(errorMessage));
    }
  }
);

export default instance;
