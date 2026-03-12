import axios from './axios';
import type { RegisterResponse, LoginResponse } from '../types/data';

export async function register(email: string, password: string): Promise<RegisterResponse> {
  const response = await axios.post<RegisterResponse>('/api/v1/auth/register', {
    email,
    password,
  });
  return response.data;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await axios.post<LoginResponse>('/api/v1/auth/login', {
    email,
    password,
  });
  return response.data;
}
