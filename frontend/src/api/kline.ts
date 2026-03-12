import axios from './axios';
import type { KLine, BackfillParams } from '../types/data';

export async function getLatest(symbol: string, period: string, limit: number): Promise<KLine[]> {
  const response = await axios.get<KLine[]>('/api/v1/kline/latest', {
    params: {
      symbol,
      period,
      limit,
    },
  });
  return response.data;
}

export async function triggerBackfill(params: BackfillParams): Promise<void> {
  await axios.post('/api/v1/kline/backfill', params);
}
