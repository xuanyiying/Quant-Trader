import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Play, Check } from 'lucide-react';
import type { Position } from '../types/data';
import { formatCurrency, safeParseFloat, calculatePriceChange } from '../utils/formatters';
import { TRADING_CONFIG } from '../constants';
import { useConfetti } from '../hooks/useAnimation';

interface TradingPanelProps {
    symbol: string;
    positions: Position[];
    lastTradePrice: string | null;
    onCreateOrder: (side: 'buy' | 'sell', qty: number) => Promise<void>;
}

const TradingPanel: React.FC<TradingPanelProps> = ({
    symbol,
    positions,
    lastTradePrice,
    onCreateOrder,
}) => {
    const [orderQty, setOrderQty] = useState<string>(TRADING_CONFIG.DEFAULT_ORDER_QTY);
    const [loading, setLoading] = useState<'buy' | 'sell' | null>(null);
    const [success, setSuccess] = useState<'buy' | 'sell' | null>(null);
    const triggerConfetti = useConfetti();

    const handleOrder = async (side: 'buy' | 'sell') => {
        const qty = safeParseFloat(orderQty);
        if (qty <= 0) return;

        setLoading(side);
        try {
            await onCreateOrder(side, qty);
            setSuccess(side);
            triggerConfetti({
                particleCount: 60,
                spread: 50,
                origin: { y: 0.7 },
                colors: side === 'buy' 
                    ? ['#00c087', '#10b981', '#34d399'] 
                    : ['#ef4444', '#f87171', '#fca5a5'],
            });
            setTimeout(() => setSuccess(null), 2000);
        } finally {
            setLoading(null);
        }
    };

    const currentPrice = lastTradePrice ? safeParseFloat(lastTradePrice) : 0;

    return (
        <motion.div 
            className="bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
        >
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                    <motion.div 
                        className="p-2 bg-blue-500/10 rounded-lg"
                        whileHover={{ scale: 1.1, rotate: 5 }}
                    >
                        <Play size={20} className="text-blue-500" />
                    </motion.div>
                    <h2 className="text-lg font-black uppercase tracking-tight">Simulator: {symbol}</h2>
                </div>
                <motion.span 
                    className="text-[10px] font-mono text-gray-500 bg-gray-900 px-3 py-1 rounded-lg border border-gray-800"
                    key={lastTradePrice}
                    initial={{ scale: 1.05 }}
                    animate={{ scale: 1 }}
                    transition={{ duration: 0.2 }}
                >
                    MARKET PRICE: {lastTradePrice ? safeParseFloat(lastTradePrice).toFixed(2) : '0.00'}
                </motion.span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
                <div className="md:col-span-4 space-y-6 pr-6 border-r border-gray-800/50">
                    <div className="space-y-2">
                        <label className="text-[10px] text-gray-500 uppercase font-black tracking-widest">
                            Order Amount ({symbol.replace('USDT', '')})
                        </label>
                        <motion.input
                            type="number"
                            value={orderQty}
                            onChange={(e) => setOrderQty(e.target.value)}
                            className="w-full bg-gray-900 border border-gray-800 rounded-xl px-4 py-3 text-sm font-black outline-none focus:ring-2 focus:ring-blue-600/50 transition-all"
                            whileFocus={{ scale: 1.02 }}
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        {/* Buy Button */}
                        <motion.button
                            onClick={() => handleOrder('buy')}
                            disabled={loading !== null}
                            className="relative bg-up text-white font-black py-4 rounded-xl transition-all shadow-xl shadow-green-900/20 text-xs tracking-widest overflow-hidden"
                            whileHover={{ 
                                scale: 1.02,
                                boxShadow: '0 0 30px rgba(0, 192, 135, 0.4)',
                            }}
                            whileTap={{ scale: 0.98 }}
                            animate={success === 'buy' ? {
                                boxShadow: [
                                    '0 0 30px rgba(0, 192, 135, 0.4)',
                                    '0 0 50px rgba(0, 192, 135, 0.6)',
                                    '0 0 30px rgba(0, 192, 135, 0.4)',
                                ],
                            } : {}}
                        >
                            <AnimatePresence mode="wait">
                                {loading === 'buy' ? (
                                    <motion.span
                                        key="loading"
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        exit={{ opacity: 0 }}
                                        className="flex items-center justify-center gap-2"
                                    >
                                        <motion.span
                                            animate={{ rotate: 360 }}
                                            transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                                            className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full"
                                        />
                                    </motion.span>
                                ) : success === 'buy' ? (
                                    <motion.span
                                        key="success"
                                        initial={{ scale: 0 }}
                                        animate={{ scale: 1 }}
                                        exit={{ scale: 0 }}
                                        className="flex items-center justify-center gap-2"
                                    >
                                        <Check size={16} />
                                        DONE
                                    </motion.span>
                                ) : (
                                    <motion.span
                                        key="default"
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        exit={{ opacity: 0 }}
                                    >
                                        BUY / LONG
                                    </motion.span>
                                )}
                            </AnimatePresence>
                        </motion.button>

                        {/* Sell Button */}
                        <motion.button
                            onClick={() => handleOrder('sell')}
                            disabled={loading !== null}
                            className="relative bg-down text-white font-black py-4 rounded-xl transition-all shadow-xl shadow-red-900/20 text-xs tracking-widest overflow-hidden"
                            whileHover={{ 
                                scale: 1.02,
                                boxShadow: '0 0 30px rgba(239, 68, 68, 0.4)',
                            }}
                            whileTap={{ scale: 0.98 }}
                            animate={success === 'sell' ? {
                                boxShadow: [
                                    '0 0 30px rgba(239, 68, 68, 0.4)',
                                    '0 0 50px rgba(239, 68, 68, 0.6)',
                                    '0 0 30px rgba(239, 68, 68, 0.4)',
                                ],
                            } : {}}
                        >
                            <AnimatePresence mode="wait">
                                {loading === 'sell' ? (
                                    <motion.span
                                        key="loading"
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        exit={{ opacity: 0 }}
                                        className="flex items-center justify-center gap-2"
                                    >
                                        <motion.span
                                            animate={{ rotate: 360 }}
                                            transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                                            className="inline-block w-4 h-4 border-2 border-white/30 border-t-white rounded-full"
                                        />
                                    </motion.span>
                                ) : success === 'sell' ? (
                                    <motion.span
                                        key="success"
                                        initial={{ scale: 0 }}
                                        animate={{ scale: 1 }}
                                        exit={{ scale: 0 }}
                                        className="flex items-center justify-center gap-2"
                                    >
                                        <Check size={16} />
                                        DONE
                                    </motion.span>
                                ) : (
                                    <motion.span
                                        key="default"
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        exit={{ opacity: 0 }}
                                    >
                                        SELL / SHORT
                                    </motion.span>
                                )}
                            </AnimatePresence>
                        </motion.button>
                    </div>
                    <div className="bg-gray-900/30 p-3 rounded-xl border border-dashed border-gray-800">
                        <p className="text-[10px] text-gray-600 italic leading-relaxed text-center">
                            Market orders are subject to platform risk engine validation.
                        </p>
                    </div>
                </div>

                <div className="md:col-span-8 space-y-4">
                    <div className="flex justify-between items-center">
                        <label className="text-[10px] text-gray-500 uppercase font-black tracking-widest">
                            Active Inventory
                        </label>
                        <motion.span 
                            className="text-[10px] bg-blue-600/10 text-blue-400 px-2 py-0.5 rounded-full font-bold"
                            key={positions.length}
                            initial={{ scale: 1.2 }}
                            animate={{ scale: 1 }}
                        >
                            {positions.length} Open
                        </motion.span>
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-h-[160px] overflow-y-auto pr-2 custom-scrollbar">
                        <AnimatePresence mode="popLayout">
                            {positions.length === 0 ? (
                                <motion.div 
                                    className="col-span-2 py-10 text-center bg-gray-900/20 border border-dashed border-gray-800 rounded-xl"
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    exit={{ opacity: 0 }}
                                >
                                    <span className="text-xs text-gray-600 font-bold italic">
                                        No active positions for this asset
                                    </span>
                                </motion.div>
                            ) : (
                                positions.map((pos, i) => {
                                    const avgPrice = safeParseFloat(pos.avg_price);
                                    const priceChange = calculatePriceChange(currentPrice, avgPrice);
                                    const isProfit = priceChange > 0;

                                    return (
                                        <motion.div
                                            key={`${pos.symbol}-${i}`}
                                            layout
                                            initial={{ opacity: 0, scale: 0.9, y: 20 }}
                                            animate={{ opacity: 1, scale: 1, y: 0 }}
                                            exit={{ opacity: 0, scale: 0.9 }}
                                            transition={{ delay: i * 0.05 }}
                                            className="bg-gray-900/40 p-4 rounded-xl border border-gray-800 flex justify-between items-center group"
                                            whileHover={{
                                                borderColor: isProfit ? 'rgba(0, 192, 135, 0.3)' : 'rgba(239, 68, 68, 0.3)',
                                                boxShadow: isProfit 
                                                    ? '0 0 20px rgba(0, 192, 135, 0.1)' 
                                                    : '0 0 20px rgba(239, 68, 68, 0.1)',
                                            }}
                                        >
                                            <div>
                                                <span className="text-xs font-black text-gray-200 uppercase">{pos.symbol}</span>
                                                <div className="text-[10px] text-gray-500 font-bold mt-1">
                                                    VOL: {safeParseFloat(pos.qty).toFixed(4)}
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-xs font-mono text-blue-400">
                                                    @{formatCurrency(avgPrice)}
                                                </div>
                                                <motion.div 
                                                    className={`text-[10px] font-black mt-1 ${isProfit ? 'text-up' : 'text-down'}`}
                                                    key={priceChange}
                                                    initial={{ scale: 1.2 }}
                                                    animate={{ scale: 1 }}
                                                >
                                                    {currentPrice ? `${priceChange.toFixed(2)}%` : '-%'}
                                                </motion.div>
                                            </div>
                                        </motion.div>
                                    );
                                })
                            )}
                        </AnimatePresence>
                    </div>
                </div>
            </div>
        </motion.div>
    );
};

export default TradingPanel;
