export interface PaperAccount {
  balance: string;
}

export interface Position {
  symbol: string;
  qty: string;
  avg_price: string;
}

export interface Strategy {
  id: number;
  name: string;
  description: string;
  price: number;
  author: string;
  is_subscribed: boolean;
  type?: string;
  performance_metrics?: unknown;
  subscriber_count?: number;
}

export interface PortfolioReport {
  total_return: number;
  sharpe_ratio: number;
  max_drawdown: number;
  win_rate: number;
  trades_count: number;
  equity_curve: { time: string; value: number }[];
}

export interface Alert {
  id: number;
  symbol: string;
  condition_type: string;
  target_value: number;
}

export interface Subscription {
  tier_name: string;
  [key: string]: unknown;
}

export interface Order {
  symbol: string;
  side: 'buy' | 'sell';
  type: string;
  qty: number;
}

export interface BackfillParams {
  exchange: string;
  symbol: string;
  start_time: string;
  end_time: string;
}

export interface DataState {
  paperAccount: PaperAccount | null;
  positions: Position[];
  strategies: Strategy[];
  portfolioReport: PortfolioReport | null;
  subscription: Subscription | null;
  alerts: Alert[];
  
  loading: {
    account: boolean;
    positions: boolean;
    strategies: boolean;
    report: boolean;
    subscription: boolean;
    alerts: boolean;
  };
}
