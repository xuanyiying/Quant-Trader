import axios from './axios';
import type { PortfolioReport } from '../types/data';

export async function getReport(): Promise<PortfolioReport> {
  const response = await axios.get<PortfolioReport>('/api/v1/analytics/portfolio');
  return response.data;
}
