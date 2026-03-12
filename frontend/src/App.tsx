import React, { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import Chart from './components/Chart';
import AlertsManager from './components/AlertsManager';
import StrategyMarketplace from './components/StrategyMarketplace';
import PortfolioReport from './components/PortfolioReport';
import Auth from './components/Auth';
import Header from './components/Header';
import ChartControls from './components/ChartControls';
import TradingPanel from './components/TradingPanel';
import SignalFeed from './components/SignalFeed';
import LoadingScreen from './components/LoadingScreen';
import { GradientOrbs } from './components/animations';
import { useMarketStore } from './store/useMarketStore';
import { useDataStore } from './store/useDataStore';
import { useWebSocket } from './hooks/useWebSocket';
import axios from './api/axios';
import { UI_CONFIG, PRICE_IDS, API_ENDPOINTS } from './constants';
import { getErrorMessage } from './utils/errorHandler';
import { showSuccess, showError, showConfirm } from './utils/notifications';
import { safeParseFloat } from './utils/formatters';
import './styles/animations.css';

const App: React.FC = () => {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);
  const [loadingProgress, setLoadingProgress] = useState(0);

  const {
    symbol,
    period,
    setSymbol,
    setPeriod,
    setKLines,
    lastTrade,
    connectionStatus,
    signals,
  } = useMarketStore();

  const {
    paperAccount,
    positions,
    strategies,
    portfolioReport,
    subscription,
    alerts,
    loading: dataLoading,
    fetchPaperAccount,
    fetchPositions,
    fetchAlerts,
    resetAccount,
    purchaseStrategy,
    createOrder,
    upgradeSubscription,
    triggerBackfill,
    fetchAll,
  } = useDataStore();

  useWebSocket();

  // Set up axios interceptor to sync token state
  useEffect(() => {
    const interceptor = axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          setToken(null);
        }
        return Promise.reject(error);
      }
    );
    return () => axios.interceptors.response.eject(interceptor);
  }, []);

  const loadHistory = useCallback(async () => {
    try {
      // 请求500条K线数据，让图表显示更密集
      const response = await axios.get(API_ENDPOINTS.MARKET.KLINES(symbol, period, 500));
      const data = response.data;
      if (Array.isArray(data)) {
        const mappedData = data
          .map((k: { symbol: string; exchange: string; period: string; o: string; h: string; l: string; c: string; v: string; t: string }) => ({
            s: k.symbol,
            e: k.exchange,
            p: k.period,
            o: k.o,
            h: k.h,
            l: k.l,
            c: k.c,
            v: k.v,
            t: k.t,
          }))
          .sort((a: { t: string }, b: { t: string }) => new Date(a.t).getTime() - new Date(b.t).getTime());
        setKLines(mappedData);
      }
    } catch (error) {
      console.error('Failed to load history:', error);
    }
  }, [symbol, period, setKLines]);

  // Load initial data with progress simulation
  useEffect(() => {
    const fetchData = async () => {
      if (!token) {
        setLoading(false);
        return;
      }
      
      // Simulate loading progress
      const progressInterval = setInterval(() => {
        setLoadingProgress((prev) => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return 90;
          }
          return prev + Math.random() * 15;
        });
      }, 200);

      try {
        await Promise.all([loadHistory(), fetchAll()]);
        setLoadingProgress(100);
        setTimeout(() => setLoading(false), 500);
      } catch (error) {
        console.error('Fetch data failed:', error);
        setLoading(false);
      } finally {
        clearInterval(progressInterval);
      }
    };

    const timer = setTimeout(() => {
      fetchData();
    }, UI_CONFIG.FETCH_DELAY);

    return () => clearTimeout(timer);
  }, [symbol, period, token, loadHistory, fetchAll]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setToken(null);
  };

  const handleLogin = (newToken: string) => {
    setToken(newToken);
  };

  const handleResetAccount = async () => {
    if (!token) return;

    const confirmed = showConfirm(
      'Are you sure you want to reset your paper account? This will clear all positions and reset balance to $100,000.'
    );

    if (!confirmed) return;

    try {
      await resetAccount();
      showSuccess('Account reset successfully');
    } catch (error) {
      showError(`Failed to reset account: ${getErrorMessage(error)}`);
    }
  };

  const handlePurchaseStrategy = async (id: number) => {
    if (!token) return;
    try {
      await purchaseStrategy(id);
      showSuccess('Strategy purchased successfully!');
    } catch (error) {
      showError(`Purchase failed: ${getErrorMessage(error)}`);
    }
  };

  const handleCreateOrder = async (side: 'buy' | 'sell', qty: number) => {
    if (!token) return;
    try {
      await createOrder({
        symbol,
        side,
        type: 'market',
        qty,
      });
      showSuccess('Order placed successfully');
      fetchPaperAccount();
      fetchPositions();
    } catch (error) {
      showError(`Order failed: ${getErrorMessage(error)}`);
    }
  };

  const handleUpgrade = async () => {
    if (!token) return;
    try {
      await upgradeSubscription(PRICE_IDS.PRO_DEFAULT);
    } catch (error) {
      showError(`Failed to start checkout: ${getErrorMessage(error)}`);
    }
  };

  const handleTriggerBackfill = async () => {
    if (!token) return;
    try {
      await triggerBackfill({
        exchange: 'binance',
        symbol,
        start_time: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
        end_time: new Date().toISOString(),
      });
      showSuccess('Backfill task started');
    } catch (error) {
      console.error('Backfill failed:', error);
      showError('Failed to start backfill');
    }
  };

  if (!token) {
    return <Auth onLogin={handleLogin} />;
  }

  if (loading) {
    return <LoadingScreen progress={loadingProgress} />;
  }

  const lastTradePrice = lastTrade ? safeParseFloat(lastTrade.price).toFixed(2) : null;

  return (
    <div className="min-h-screen bg-background text-gray-200 selection:bg-blue-500/30 relative">
      {/* Background Effects */}
      <GradientOrbs />
      
      {/* Main Content */}
      <div className="relative z-10">
        <Header
          paperAccount={paperAccount}
          subscription={subscription}
          connectionStatus={connectionStatus}
          lastTradePrice={lastTradePrice}
          onLogout={handleLogout}
          onResetAccount={handleResetAccount}
          onUpgrade={handleUpgrade}
        />

        <motion.main 
          className="max-w-[1600px] mx-auto p-4 lg:p-8 space-y-8 pt-32"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.5 }}
        >
          <div className="grid grid-cols-12 gap-8">
            {/* Main Content Area */}
            <motion.div 
              className="col-span-12 lg:col-span-9 space-y-8"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
            >
              {/* Chart Section */}
              <motion.div 
                className="bg-card/80 backdrop-blur-sm p-6 rounded-2xl shadow-xl border border-gray-800/50"
                whileHover={{ borderColor: 'rgba(59, 130, 246, 0.2)' }}
                transition={{ duration: 0.3 }}
              >
                <ChartControls
                  symbol={symbol}
                  period={period}
                  onSymbolChange={setSymbol}
                  onPeriodChange={setPeriod}
                  onRefresh={loadHistory}
                  onBackfill={handleTriggerBackfill}
                />
                <div className="h-[500px] w-full">
                  <Chart />
                </div>
              </motion.div>

              {/* Trading Panel */}
              <TradingPanel
                symbol={symbol}
                positions={positions}
                lastTradePrice={lastTradePrice}
                onCreateOrder={handleCreateOrder}
              />
            </motion.div>

            {/* Sidebar */}
            <motion.div 
              className="col-span-12 lg:col-span-3 space-y-8"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
            >
              <StrategyMarketplace
                strategies={strategies}
                loading={dataLoading.strategies}
                onPurchase={handlePurchaseStrategy}
              />
              <PortfolioReport report={portfolioReport} loading={dataLoading.report} />
              <AlertsManager alerts={alerts} symbol={symbol} onRefresh={fetchAlerts} />
              <SignalFeed signals={signals} />
            </motion.div>
          </div>
        </motion.main>
      </div>
    </div>
  );
};

export default App;
