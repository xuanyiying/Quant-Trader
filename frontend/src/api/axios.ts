import axios from 'axios';

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
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      // We don't want to use window.location.reload() here as it might cause loops
      // The App.tsx component will handle the state change if we trigger it correctly
    }
    return Promise.reject(error);
  }
);

export default instance;
