import axios from './axios';
import type { Alert, CreateAlertRequest } from '../types/data';

export async function getAlerts(): Promise<Alert[]> {
  const response = await axios.get<Alert[]>('/api/v1/alerts');
  return response.data;
}

export async function createAlert(alert: CreateAlertRequest): Promise<Alert> {
  const response = await axios.post<Alert>('/api/v1/alerts', alert);
  return response.data;
}

export async function deleteAlert(id: number): Promise<void> {
  await axios.delete(`/api/v1/alerts/${id}`);
}
