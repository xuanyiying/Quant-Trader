import { create } from 'zustand';
import type { KLine, Trade, StrategySignal } from '../types/market';
import { CHART_CONFIG, TRADING_CONFIG } from '../constants';

interface MarketStore {
  symbol: string;
  period: string;
  klines: KLine[];
  lastTrade: Trade | null;
  signals: StrategySignal[];
  connectionStatus: 'Connected' | 'Disconnected' | 'Connecting';

  setSymbol: (symbol: string) => void;
  setPeriod: (period: string) => void;
  setKLines: (klines: KLine[]) => void;
  updateKLine: (kline: KLine) => void;
  setLastTrade: (trade: Trade) => void;
  addSignal: (signal: StrategySignal) => void;
  setConnectionStatus: (status: 'Connected' | 'Disconnected' | 'Connecting') => void;
}

export const useMarketStore = create<MarketStore>((set) => ({
  symbol: TRADING_CONFIG.DEFAULT_SYMBOL,
  period: TRADING_CONFIG.DEFAULT_PERIOD,
  klines: [],
  lastTrade: null,
  signals: [],
  connectionStatus: 'Disconnected',

  setSymbol: (symbol) => set({ symbol }),
  setPeriod: (period) => set({ period }),
  setKLines: (klines) => set({ klines }),
  updateKLine: (kline) => set((state) => {
    if (kline.s.toUpperCase() !== state.symbol.toUpperCase()) return state;
    if (kline.p !== state.period) return state;

    const newKLines = [...state.klines];
    const lastIdx = newKLines.length - 1;

    if (lastIdx >= 0 && new Date(newKLines[lastIdx].t).getTime() === new Date(kline.t).getTime()) {
      newKLines[lastIdx] = kline;
    } else {
      newKLines.push(kline);
    }

    if (newKLines.length > CHART_CONFIG.MAX_KLINES) newKLines.shift();
    return { klines: newKLines };
  }),
  setLastTrade: (trade) => set((state) => {
    if (trade.symbol.toUpperCase() !== state.symbol.toUpperCase()) return state;
    return { lastTrade: trade };
  }),
  addSignal: (signal) => set((state) => ({
    signals: [signal, ...state.signals].slice(0, CHART_CONFIG.MAX_SIGNALS)
  })),
  setConnectionStatus: (status) => set({ connectionStatus: status }),
}));
