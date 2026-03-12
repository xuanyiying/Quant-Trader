import { create } from 'zustand';
import type { DataState, Order, BackfillParams } from '../types/data';
import { logError } from '../utils/errorHandler';
import {
  getAccount,
  resetAccount as resetAccountApi,
  createOrder as createOrderApi,
  getPositions,
} from '../api/paper';
import { listStrategies, purchaseStrategy as purchaseStrategyApi } from '../api/marketplace';
import { getReport } from '../api/portfolio';
import { getSubscription, createCheckoutSession } from '../api/subscription';
import { getAlerts } from '../api/alert';
import { triggerBackfill as triggerBackfillApi } from '../api/kline';

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
      const data = await getAccount();
      set({ paperAccount: data });
    } catch (error) {
      logError('fetchPaperAccount', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, account: false } }));
    }
  },

  fetchPositions: async () => {
    set((state) => ({ loading: { ...state.loading, positions: true } }));
    try {
      const data = await getPositions();
      set({ positions: data || [] });
    } catch (error) {
      logError('fetchPositions', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, positions: false } }));
    }
  },

  fetchStrategies: async () => {
    set((state) => ({ loading: { ...state.loading, strategies: true } }));
    try {
      const data = await listStrategies();
      set({ strategies: data || [] });
    } catch (error) {
      logError('fetchStrategies', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, strategies: false } }));
    }
  },

  fetchPortfolioReport: async () => {
    set((state) => ({ loading: { ...state.loading, report: true } }));
    try {
      const data = await getReport();
      set({ portfolioReport: data });
    } catch (error) {
      logError('fetchPortfolioReport', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, report: false } }));
    }
  },

  resetAccount: async () => {
    try {
      await resetAccountApi();
      await Promise.all([get().fetchPaperAccount(), get().fetchPositions()]);
    } catch (error) {
      logError('resetAccount', error);
      throw error;
    }
  },

  purchaseStrategy: async (id: number) => {
    try {
      await purchaseStrategyApi(id);
      await get().fetchStrategies();
    } catch (error) {
      logError('purchaseStrategy', error);
      throw error;
    }
  },

  createOrder: async (order: Order) => {
    try {
      await createOrderApi(order);
      await Promise.all([get().fetchPaperAccount(), get().fetchPositions()]);
    } catch (error) {
      logError('createOrder', error);
      throw error;
    }
  },

  upgradeSubscription: async (priceId: string) => {
    try {
      const data = await createCheckoutSession(priceId);
      if (data?.session_url) {
        window.location.href = data.session_url;
      }
    } catch (error) {
      logError('upgradeSubscription', error);
      throw error;
    }
  },

  triggerBackfill: async (params: BackfillParams) => {
    try {
      await triggerBackfillApi(params);
    } catch (error) {
      logError('triggerBackfill', error);
      throw error;
    }
  },

  fetchSubscription: async () => {
    set((state) => ({ loading: { ...state.loading, subscription: true } }));
    try {
      const data = await getSubscription();
      set({ subscription: data });
    } catch (error) {
      logError('fetchSubscription', error);
    } finally {
      set((state) => ({ loading: { ...state.loading, subscription: false } }));
    }
  },

  fetchAlerts: async () => {
    set((state) => ({ loading: { ...state.loading, alerts: true } }));
    try {
      const data = await getAlerts();
      set({ alerts: data || [] });
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