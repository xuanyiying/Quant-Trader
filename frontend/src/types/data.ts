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
  price: string;
  author: string;
  metrics: Record<string, unknown>;
  is_subscribed?: boolean;
}

export interface EquityPoint {
  timestamp: string;
  equity: string;
  drawdown: string;
}

export interface PortfolioReport {
  final_balance: string;
  total_return: string;
  sharpe_ratio: string;
  max_drawdown: string;
  win_rate: string;
  total_trades: number;
  equity_curve: EquityPoint[];
}

export interface Alert {
  id: number;
  symbol: string;
  condition_type: string;
  target_value: string;
  status: string;
  created_at: string;
  triggered_at: string | null;
}

export interface Subscription {
  tier_name: string;
  max_symbols: number;
  status: string;
  expires_at: string | null;
}

export interface Order {
  symbol: string;
  side: 'buy' | 'sell';
  type: string;
  qty: number;
}

export interface OrderResponse {
  id: number;
}

export interface BackfillParams {
  exchange: string;
  symbol: string;
  start_time: string;
  end_time: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface ErrorResponse {
  success: false;
  error: string;
  code?: number;
}

export interface RegisterResponse {
  message: string;
  user_id: number;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    email: string;
  };
}

export interface CreateOrderRequest {
  symbol: string;
  side: 'buy' | 'sell';
  type: string;
  qty: number;
}

export interface CreateAlertRequest {
  symbol: string;
  condition_type: string;
  target_value: string;
}

export interface CheckoutSessionResponse {
  session_url: string;
}

export interface KLine {
  timestamp: string;
  open: string;
  high: string;
  low: string;
  close: string;
  volume: string;
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
