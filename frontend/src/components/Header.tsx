import React from 'react';
import { motion } from 'framer-motion';
import { TrendingUp, Activity, User, LogOut, RefreshCw } from 'lucide-react';
import type { PaperAccount, Subscription } from '../types/data';
import { safeParseFloat } from '../utils/formatters';
import { SUBSCRIPTION_TIERS } from '../constants';
import { AnimatedNumber } from './animations';
import { useScrollDirection } from '../hooks/useAnimation';

interface HeaderProps {
  paperAccount: PaperAccount | null;
  subscription: Subscription | null;
  connectionStatus: 'Connected' | 'Disconnected' | 'Connecting';
  lastTradePrice: string | null;
  onLogout: () => void;
  onResetAccount: () => void;
  onUpgrade: () => void;
}

const Header: React.FC<HeaderProps> = ({
  paperAccount,
  subscription,
  connectionStatus,
  lastTradePrice,
  onLogout,
  onResetAccount,
  onUpgrade,
}) => {
  const isPro = subscription?.tier_name !== SUBSCRIPTION_TIERS.FREE;
  const { scrollY } = useScrollDirection();
  const isScrolled = scrollY > 50;

  const balance = paperAccount ? safeParseFloat(paperAccount.balance) : 0;
  const prevPriceRef = React.useRef<string | null>(null);
  const currentPrice = lastTradePrice || '0';
  
  const priceDirection = prevPriceRef.current !== null 
    ? parseFloat(currentPrice) > parseFloat(prevPriceRef.current) ? 'up' : 
      parseFloat(currentPrice) < parseFloat(prevPriceRef.current) ? 'down' : null
    : null;
  
  React.useEffect(() => {
    prevPriceRef.current = currentPrice;
  }, [currentPrice]);

  return (
    <motion.header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled ? 'py-2' : 'py-4'
      }`}
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5, ease: 'easeOut' }}
    >
      <div
        className={`
          mx-4 lg:mx-8 rounded-2xl border border-gray-800/50
          transition-all duration-300
          ${isScrolled 
            ? 'bg-card/90 backdrop-blur-xl shadow-2xl' 
            : 'bg-card backdrop-blur-sm shadow-xl'
          }
        `}
      >
        <div className="flex flex-col md:flex-row justify-between items-center gap-6 p-6 relative overflow-hidden">
          {/* Background gradient on hover */}
          <motion.div
            className="absolute inset-0 bg-gradient-to-r from-blue-600/5 to-transparent"
            initial={{ opacity: 0 }}
            whileHover={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          />

          {/* Logo section */}
          <motion.div
            className="flex items-center gap-4 relative"
            whileHover={{ scale: 1.02 }}
            transition={{ duration: 0.2 }}
          >
            <motion.div
              className="p-3 bg-blue-600 rounded-xl shadow-lg"
              whileHover={{ 
                boxShadow: '0 0 30px rgba(59, 130, 246, 0.5)',
              }}
              animate={{
                boxShadow: [
                  '0 10px 40px rgba(59, 130, 246, 0.2)',
                  '0 10px 40px rgba(59, 130, 246, 0.4)',
                  '0 10px 40px rgba(59, 130, 246, 0.2)',
                ],
              }}
              transition={{
                boxShadow: {
                  duration: 2,
                  repeat: Infinity,
                  ease: 'easeInOut',
                },
              }}
            >
              <TrendingUp size={28} className="text-white" />
            </motion.div>
            <div>
              <h1 className="text-2xl font-black tracking-tight text-white">
                QUANT<span className="text-blue-500">TRADER</span>
              </h1>
              <div className="flex items-center gap-2">
                <p className="text-[10px] text-gray-500 uppercase tracking-widest font-black">
                  {subscription ? `${subscription.tier_name} MEMBER` : 'LOADING...'}
                </p>
                {isPro ? (
                  <motion.span
                    className="bg-yellow-500/10 text-yellow-500 text-[8px] px-1.5 py-0.5 rounded font-black border border-yellow-500/20 uppercase"
                    animate={{
                      boxShadow: [
                        '0 0 0 rgba(234, 179, 8, 0)',
                        '0 0 10px rgba(234, 179, 8, 0.3)',
                        '0 0 0 rgba(234, 179, 8, 0)',
                      ],
                    }}
                    transition={{
                      duration: 2,
                      repeat: Infinity,
                      ease: 'easeInOut',
                    }}
                  >
                    Pro Access
                  </motion.span>
                ) : (
                  <motion.button
                    onClick={onUpgrade}
                    className="bg-blue-600/10 text-blue-400 text-[8px] px-1.5 py-0.5 rounded font-black border border-blue-500/20 uppercase"
                    whileHover={{
                      backgroundColor: 'rgba(59, 130, 246, 1)',
                      color: '#fff',
                      scale: 1.05,
                    }}
                    whileTap={{ scale: 0.95 }}
                  >
                    Upgrade Now
                  </motion.button>
                )}
              </div>
            </div>
          </motion.div>

          {/* Stats section */}
          <div className="flex flex-wrap justify-center gap-4 relative">
            {/* Balance card */}
            <motion.div
              className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px] group relative"
              whileHover={{
                scale: 1.02,
                borderColor: 'rgba(59, 130, 246, 0.3)',
                boxShadow: '0 0 20px rgba(59, 130, 246, 0.1)',
              }}
              transition={{ duration: 0.2 }}
            >
              <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1 text-center">
                Paper Balance
              </span>
              <span className="text-lg font-black text-up flex items-center gap-2">
                $<AnimatedNumber value={balance} decimals={2} duration={1500} />
              </span>
              <motion.button
                onClick={onResetAccount}
                className="absolute -top-2 -right-2 bg-gray-800 text-gray-500 p-1.5 rounded-lg border border-gray-700 opacity-0 group-hover:opacity-100 transition-all shadow-xl"
                title="Reset Account"
                whileHover={{ 
                  backgroundColor: 'rgba(239, 68, 68, 0.2)',
                  color: '#ef4444',
                  scale: 1.1,
                }}
                whileTap={{ scale: 0.9 }}
              >
                <RefreshCw size={12} />
              </motion.button>
            </motion.div>

            {/* Status card */}
            <motion.div
              className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px]"
              whileHover={{
                scale: 1.02,
                borderColor: connectionStatus === 'Connected' 
                  ? 'rgba(0, 192, 135, 0.3)' 
                  : 'rgba(239, 68, 68, 0.3)',
              }}
              transition={{ duration: 0.2 }}
            >
              <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1">
                Status
              </span>
              <span className={`text-sm font-black flex items-center gap-2 ${
                connectionStatus === 'Connected' ? 'text-up' : 'text-down'
              }`}>
                <motion.div
                  animate={connectionStatus === 'Connected' ? {
                    scale: [1, 1.2, 1],
                    opacity: [1, 0.7, 1],
                  } : {}}
                  transition={{
                    duration: 1.5,
                    repeat: Infinity,
                    ease: 'easeInOut',
                  }}
                >
                  <Activity size={16} />
                </motion.div>
                {connectionStatus.toUpperCase()}
              </span>
            </motion.div>

            {/* Price card */}
            <motion.div
              className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px]"
              whileHover={{
                scale: 1.02,
                borderColor: 'rgba(59, 130, 246, 0.3)',
              }}
              transition={{ duration: 0.2 }}
            >
              <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1">
                Market Price
              </span>
              <motion.span
                className={`text-lg font-black ${
                  priceDirection === 'up' ? 'text-up' : 
                  priceDirection === 'down' ? 'text-down' : 'text-blue-400'
                }`}
                key={lastTradePrice}
                initial={{ scale: 1.1 }}
                animate={{ scale: 1 }}
                transition={{ duration: 0.2 }}
              >
                {lastTradePrice || '--.--'} 
                <span className="text-[10px] ml-1 text-gray-600 font-bold">USDT</span>
              </motion.span>
            </motion.div>
          </div>

          {/* Actions section */}
          <div className="flex items-center gap-3 relative">
            <motion.button
              className="bg-gray-900 text-gray-400 p-2.5 rounded-xl border border-gray-800 font-bold"
              whileHover={{
                backgroundColor: 'rgba(31, 41, 55, 1)',
                color: '#60a5fa',
                scale: 1.05,
              }}
              whileTap={{ scale: 0.95 }}
            >
              <User size={20} />
            </motion.button>
            <motion.button
              onClick={onLogout}
              className="bg-gray-900 text-gray-400 p-2.5 rounded-xl border border-gray-800 font-bold"
              title="Logout"
              whileHover={{
                backgroundColor: 'rgba(239, 68, 68, 0.2)',
                color: '#ef4444',
                borderColor: 'rgba(239, 68, 68, 0.3)',
                scale: 1.05,
              }}
              whileTap={{ scale: 0.95 }}
            >
              <LogOut size={20} />
            </motion.button>
          </div>
        </div>
      </div>
    </motion.header>
  );
};

export default Header;
