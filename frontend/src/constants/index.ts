// API Endpoints
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    REGISTER: '/api/v1/auth/register',
  },
  PAPER: {
    ACCOUNT: '/api/v1/paper/account',
    POSITIONS: '/api/v1/paper/positions',
    ORDERS: '/api/v1/paper/orders',
    RESET: '/api/v1/paper/account/reset',
  },
  MARKET: {
    STRATEGIES: '/api/v1/market/strategies',
    KLINES: (symbol: string, period: string) => `/api/v1/klines/${symbol}?period=${period}`,
  },
  MARKETPLACE: {
    STRATEGIES: '/api/v1/marketplace',
    PURCHASE: (id: number) => `/api/v1/marketplace/${id}/purchase`,
  },
  KLINES: {
    LATEST: '/api/v1/klines/latest',
    BACKFILL: '/api/v1/klines/backfill',
  },
  ANALYTICS: {
    PORTFOLIO: '/api/v1/analytics/portfolio',
  },
  SUBSCRIPTION: {
    INFO: '/api/v1/subscription',
    CHECKOUT: '/api/v1/subscription/checkout',
  },
  ALERTS: {
    LIST: '/api/v1/alerts',
    DELETE: (id: number) => `/api/v1/alerts/${id}`,
  },
  BACKFILL: '/api/v1/klines/backfill',
} as const;

// WebSocket Configuration
export const WS_CONFIG = {
  RECONNECT_DELAY: 5000,
  DEV_HOST: 'localhost:8080',
} as const;

// Chart Configuration
export const CHART_CONFIG = {
  UP_COLOR: '#00c087',
  DOWN_COLOR: '#ff3b30',
  BACKGROUND_COLOR: '#1e1e1e',
  MAX_KLINES: 1000,
  MAX_SIGNALS: 100,
} as const;

// Trading Configuration
export const TRADING_CONFIG = {
  DEFAULT_SYMBOL: 'BTCUSDT',
  DEFAULT_PERIOD: '1m',
  PERIODS: ['1m', '5m', '15m', '1h', '4h', '1d'] as const,
  DEFAULT_ORDER_QTY: '1',
  DEFAULT_BALANCE: 100000,
} as const;

// UI Configuration
export const UI_CONFIG = {
  FETCH_DELAY: 300,
  LOADING_SPINNER_SIZE: 48,
} as const;

// Subscription Tiers
export const SUBSCRIPTION_TIERS = {
  FREE: 'Free',
  PRO: 'Pro',
} as const;

// Price IDs
export const PRICE_IDS = {
  PRO_DEFAULT: 'price_pro_default',
} as const;
