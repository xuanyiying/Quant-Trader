import { create } from 'zustand';
import axios from '../api/axios';
import type { DataState, Order, BackfillParams } from '../types/data';
import { API_ENDPOINTS } from '../constants';
import { logError } from '../utils/errorHandler';

interface DataActions {
  fetchPaperAccount: () => Promise<void>;
  fetchPositions: () => Promise<void>;
  fetchStrategies: () => Promise<void>;
  fetchPortfolioReport: () => Promise<void>;
  fetchSubscription: () => Promise<void>;
  fetchAlerts: () => Promise<void>;
  resetAccount: () => Promise<void>;
  purchaseStrategy: (id: number) => Promise<void>;
  createOrder: (order: Order) => Promise<void>;
  upgradeSubscription: (priceId: string) => Promise<void>;
  triggerBackfill: (params: BackfillParams) => Promise<void>;
  fetchAll: () => Promise<void>;
}

export const useDataStore = create<DataState & DataActions>((set, get) => ({
  paperAccount: null,
  positions: [],
  strategies: [],
  portfolioReport: null,
  subscription: null,
  alerts: [],
  loading: {
    account: false,
    positions: false,
    strategies: false,
    report: false,
    subscription: false,
    alerts: false,
  },

  fetchPaperAccount: async () => {
    set((state) => ({ loading: { ...state.loading, account: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.PAPER.ACCOUNT);
      set({ paperAccount: response.data });
    } catch (error) {
      logError('fetchPaperAccount', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, account: false } }));
    }
  },

  fetchPositions: async () => {
    set((state) => ({ loading: { ...state.loading, positions: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.PAPER.POSITIONS);
      set({ positions: response.data || [] });
    } catch (error) {
      logError('fetchPositions', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, positions: false } }));
    }
  },

  fetchStrategies: async () => {
    set((state) => ({ loading: { ...state.loading, strategies: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.MARKET.STRATEGIES);
      set({ strategies: response.data || [] });
    } catch (error) {
      logError('fetchStrategies', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, strategies: false } }));
    }
  },

  fetchPortfolioReport: async () => {
    set((state) => ({ loading: { ...state.loading, report: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.ANALYTICS.PORTFOLIO);
      set({ portfolioReport: response.data });
    } catch (error) {
      logError('fetchPortfolioReport', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, report: false } }));
    }
  },

  resetAccount: async () => {
    try {
      await axios.post(API_ENDPOINTS.PAPER.RESET);
      await Promise.all([get().fetchPaperAccount(), get().fetchPositions()]);
    } catch (error) {
      logError('resetAccount', error);
      throw error;
    }
  },

  purchaseStrategy: async (id: number) => {
    try {
      await axios.post(`${API_ENDPOINTS.MARKET.STRATEGIES}/${id}/purchase`);
      await get().fetchStrategies();
    } catch (error) {
      logError('purchaseStrategy', error);
      throw error;
    }
  },

  createOrder: async (order: Order) => {
    try {
      await axios.post(API_ENDPOINTS.PAPER.ORDERS, order);
      await Promise.all([get().fetchPaperAccount(), get().fetchPositions()]);
    } catch (error) {
      logError('createOrder', error);
      throw error;
    }
  },

  upgradeSubscription: async (priceId: string) => {
    try {
      const response = await axios.post(API_ENDPOINTS.SUBSCRIPTION.CHECKOUT, {
        price_id: priceId
      });
      if (response.data?.url) {
        window.location.href = response.data.url;
      }
    } catch (error) {
      logError('upgradeSubscription', error);
      throw error;
    }
  },

  triggerBackfill: async (params: BackfillParams) => {
    try {
      await axios.post(API_ENDPOINTS.BACKFILL, params);
    } catch (error) {
      logError('triggerBackfill', error);
      throw error;
    }
  },

  fetchSubscription: async () => {
    set((state) => ({ loading: { ...state.loading, subscription: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.SUBSCRIPTION.INFO);
      set({ subscription: response.data });
    } catch (error) {
      logError('fetchSubscription', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, subscription: false } }));
    }
  },

  fetchAlerts: async () => {
    set((state) => ({ loading: { ...state.loading, alerts: true } }));
    try {
      const response = await axios.get(API_ENDPOINTS.ALERTS.LIST);
      set({ alerts: response.data || [] });
    } catch (error) {
      logError('fetchAlerts', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, alerts: false } }));
    }
  },

  fetchAll: async () => {
    await Promise.all([
      get().fetchPaperAccount(),
      get().fetchPositions(),
      get().fetchStrategies(),
      get().fetchPortfolioReport(),
      get().fetchSubscription(),
      get().fetchAlerts(),
    ]);
  },
}));