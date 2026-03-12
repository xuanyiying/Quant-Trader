import axios from './axios';
import type { Strategy } from '../types/data';

export async function listStrategies(): Promise<Strategy[]> {
  const response = await axios.get<Strategy[]>('/api/v1/marketplace/strategies');
  return response.data;
}

export async function purchaseStrategy(id: number): Promise<void> {
  await axios.post(`/api/v1/marketplace/strategies/${id}/purchase`);
}
