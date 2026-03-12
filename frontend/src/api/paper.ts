import axios from './axios';
import type { PaperAccount, Position, CreateOrderRequest, OrderResponse } from '../types/data';

export async function getAccount(): Promise<PaperAccount> {
  const response = await axios.get<PaperAccount>('/api/v1/paper/account');
  return response.data;
}

export async function resetAccount(): Promise<PaperAccount> {
  const response = await axios.post<PaperAccount>('/api/v1/paper/account/reset');
  return response.data;
}

export async function createOrder(order: CreateOrderRequest): Promise<OrderResponse> {
  const response = await axios.post<OrderResponse>('/api/v1/paper/orders', order);
  return response.data;
}

export async function getPositions(): Promise<Position[]> {
  const response = await axios.get<Position[]>('/api/v1/paper/positions');
  return response.data;
}
