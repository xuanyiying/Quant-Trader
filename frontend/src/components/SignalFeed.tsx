import React from 'react';
import { Activity } from 'lucide-react';
import type { StrategySignal } from '../types/market';
import { formatTime, safeParseFloat } from '../utils/formatters';

interface SignalFeedProps {
    signals: StrategySignal[];
}

const SignalFeed: React.FC<SignalFeedProps> = ({ signals }) => {
    return (
        <div className="bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50">
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                    <Activity size={18} className="text-blue-400" />
                    <h2 className="font-black uppercase tracking-tight">Signal Feed</h2>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-1.5 h-1.5 bg-up rounded-full animate-ping"></div>
                    <span className="text-[8px] font-black text-up uppercase">LIVE</span>
                </div>
            </div>

            <div className="space-y-4 overflow-y-auto max-h-[400px] pr-2 custom-scrollbar">
                {signals.length === 0 ? (
                    <div className="text-center py-12 text-gray-700 bg-gray-900/20 border border-dashed border-gray-800 rounded-xl">
                        <p className="text-xs font-bold italic">Listening for market triggers...</p>
                    </div>
                ) : (
                    signals.map((sig, i) => (
                        <div
                            key={i}
                            className="bg-gray-900/50 p-4 rounded-xl border border-gray-800 space-y-3 relative overflow-hidden group"
                        >
                            <div className={`absolute left-0 top-0 bottom-0 w-1 ${sig.action === 'buy' ? 'bg-up' : 'bg-down'
                                }`}></div>
                            <div className="flex justify-between items-start">
                                <span className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">
                                    {sig.strategy}
                                </span>
                                <span className="text-[8px] text-gray-600 font-mono tracking-tighter">
                                    {formatTime(sig.time)}
                                </span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-xs font-black text-white">{sig.symbol}</span>
                                <span className={`px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-widest ${sig.action === 'buy' ? 'bg-up/10 text-up' : 'bg-down/10 text-down'
                                    }`}>
                                    {sig.action}
                                </span>
                            </div>
                            <div className="text-right text-[10px] font-mono text-gray-500 font-bold">
                                PRC: {safeParseFloat(sig.price).toFixed(2)}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};

export default SignalFeed;
