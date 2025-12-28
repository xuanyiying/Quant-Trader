import React, { useState } from 'react';
import { Play } from 'lucide-react';
import type { Position } from '../types/data';
import { formatCurrency, safeParseFloat, calculatePriceChange } from '../utils/formatters';
import { TRADING_CONFIG } from '../constants';

interface TradingPanelProps {
    symbol: string;
    positions: Position[];
    lastTradePrice: string | null;
    onCreateOrder: (side: 'buy' | 'sell', qty: number) => void;
}

const TradingPanel: React.FC<TradingPanelProps> = ({
    symbol,
    positions,
    lastTradePrice,
    onCreateOrder,
}) => {
    const [orderQty, setOrderQty] = useState<string>(TRADING_CONFIG.DEFAULT_ORDER_QTY);

    const handleOrder = (side: 'buy' | 'sell') => {
        const qty = safeParseFloat(orderQty);
        if (qty > 0) {
            onCreateOrder(side, qty);
        }
    };

    const currentPrice = lastTradePrice ? safeParseFloat(lastTradePrice) : 0;

    return (
        <div className="bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50">
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                    <div className="p-2 bg-blue-500/10 rounded-lg">
                        <Play size={20} className="text-blue-500" />
                    </div>
                    <h2 className="text-lg font-black uppercase tracking-tight">Simulator: {symbol}</h2>
                </div>
                <span className="text-[10px] font-mono text-gray-500 bg-gray-900 px-3 py-1 rounded-lg border border-gray-800">
                    MARKET PRICE: {lastTradePrice ? safeParseFloat(lastTradePrice).toFixed(2) : '0.00'}
                </span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-12 gap-8">
                <div className="md:col-span-4 space-y-6 pr-6 border-r border-gray-800/50">
                    <div className="space-y-2">
                        <label className="text-[10px] text-gray-500 uppercase font-black tracking-widest">
                            Order Amount ({symbol.replace('USDT', '')})
                        </label>
                        <input
                            type="number"
                            value={orderQty}
                            onChange={(e) => setOrderQty(e.target.value)}
                            className="w-full bg-gray-900 border border-gray-800 rounded-xl px-4 py-3 text-sm font-black outline-none focus:ring-2 focus:ring-blue-600/50 transition-all"
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        <button
                            onClick={() => handleOrder('buy')}
                            className="bg-up hover:bg-green-600 text-white font-black py-4 rounded-xl transition-all shadow-xl shadow-green-900/20 active:scale-95 text-xs tracking-widest"
                        >
                            BUY / LONG
                        </button>
                        <button
                            onClick={() => handleOrder('sell')}
                            className="bg-down hover:bg-red-600 text-white font-black py-4 rounded-xl transition-all shadow-xl shadow-red-900/20 active:scale-95 text-xs tracking-widest"
                        >
                            SELL / SHORT
                        </button>
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
                        <span className="text-[10px] bg-blue-600/10 text-blue-400 px-2 py-0.5 rounded-full font-bold">
                            {positions.length} Open
                        </span>
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-h-[160px] overflow-y-auto pr-2 custom-scrollbar">
                        {positions.length === 0 ? (
                            <div className="col-span-2 py-10 text-center bg-gray-900/20 border border-dashed border-gray-800 rounded-xl">
                                <span className="text-xs text-gray-600 font-bold italic">
                                    No active positions for this asset
                                </span>
                            </div>
                        ) : (
                            positions.map((pos, i) => {
                                const avgPrice = safeParseFloat(pos.avg_price);
                                const priceChange = calculatePriceChange(currentPrice, avgPrice);
                                const isProfit = priceChange > 0;

                                return (
                                    <div
                                        key={i}
                                        className="bg-gray-900/40 p-4 rounded-xl border border-gray-800 flex justify-between items-center group hover:border-blue-500/30 transition-all"
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
                                            <div className={`text-[10px] font-black mt-1 ${isProfit ? 'text-up' : 'text-down'}`}>
                                                {currentPrice ? `${priceChange.toFixed(2)}%` : '-%'}
                                            </div>
                                        </div>
                                    </div>
                                );
                            })
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TradingPanel;
