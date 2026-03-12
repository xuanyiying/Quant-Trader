import axios from './axios';
import type { Subscription, CheckoutSessionResponse } from '../types/data';

export async function getSubscription(): Promise<Subscription> {
  const response = await axios.get<Subscription>('/api/v1/subscription');
  return response.data;
}

export async function createCheckoutSession(priceId: string): Promise<CheckoutSessionResponse> {
  const response = await axios.post<CheckoutSessionResponse>('/api/v1/subscription/checkout', {
    price_id: priceId,
  });
  return response.data;
}
