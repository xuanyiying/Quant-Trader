import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Activity, AlertTriangle } from 'lucide-react';
import type { StrategySignal } from '../types/market';
import { formatTime, safeParseFloat } from '../utils/formatters';

interface SignalFeedProps {
    signals: StrategySignal[];
}

const SignalFeed: React.FC<SignalFeedProps> = ({ signals }) => {
    return (
        <motion.div 
            className="bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
        >
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-2">
                    <motion.div
                        animate={{
                            rotate: [0, 10, -10, 0],
                        }}
                        transition={{
                            duration: 2,
                            repeat: Infinity,
                            ease: 'easeInOut',
                            repeatDelay: 3,
                        }}
                    >
                        <Activity size={18} className="text-blue-400" />
                    </motion.div>
                    <h2 className="font-black uppercase tracking-tight">Signal Feed</h2>
                </div>
                <div className="flex items-center gap-1.5">
                    <motion.div
                        className="w-1.5 h-1.5 bg-up rounded-full"
                        animate={{
                            scale: [1, 1.5, 1],
                            opacity: [1, 0.5, 1],
                        }}
                        transition={{
                            duration: 1.5,
                            repeat: Infinity,
                            ease: 'easeInOut',
                        }}
                    />
                    <span className="text-[8px] font-black text-up uppercase">LIVE</span>
                </div>
            </div>

            <div className="space-y-4 overflow-y-auto max-h-[400px] pr-2 custom-scrollbar">
                <AnimatePresence mode="popLayout">
                    {signals.length === 0 ? (
                        <motion.div 
                            className="text-center py-12 text-gray-700 bg-gray-900/20 border border-dashed border-gray-800 rounded-xl"
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                        >
                            <motion.div
                                animate={{
                                    y: [0, -5, 0],
                                }}
                                transition={{
                                    duration: 2,
                                    repeat: Infinity,
                                    ease: 'easeInOut',
                                }}
                            >
                                <Activity size={32} className="mx-auto mb-3 text-gray-600" />
                            </motion.div>
                            <p className="text-xs font-bold italic">Listening for market triggers...</p>
                        </motion.div>
                    ) : (
                        signals.map((sig, i) => {
                            const isHighPriority = sig.action === 'sell' && safeParseFloat(sig.price) > 50000;
                            
                            return (
                                <motion.div
                                    key={`${sig.symbol}-${sig.time}-${i}`}
                                    layout
                                    initial={{ opacity: 0, x: 50, scale: 0.9 }}
                                    animate={{ opacity: 1, x: 0, scale: 1 }}
                                    exit={{ opacity: 0, x: -50, scale: 0.9 }}
                                    transition={{ 
                                        duration: 0.4,
                                        delay: i * 0.05,
                                        type: 'spring',
                                        stiffness: 100,
                                    }}
                                    className={`
                                        bg-gray-900/50 p-4 rounded-xl border border-gray-800 
                                        space-y-3 relative overflow-hidden group
                                        ${isHighPriority ? 'border-yellow-500/30' : ''}
                                    `}
                                    whileHover={{
                                        scale: 1.02,
                                        borderColor: sig.action === 'buy' ? 'rgba(0, 192, 135, 0.3)' : 'rgba(239, 68, 68, 0.3)',
                                    }}
                                >
                                    {/* Priority indicator */}
                                    {isHighPriority && (
                                        <motion.div
                                            className="absolute top-2 right-2 text-yellow-500"
                                            animate={{
                                                scale: [1, 1.2, 1],
                                                opacity: [0.7, 1, 0.7],
                                            }}
                                            transition={{
                                                duration: 1,
                                                repeat: Infinity,
                                                ease: 'easeInOut',
                                            }}
                                        >
                                            <AlertTriangle size={14} />
                                        </motion.div>
                                    )}

                                    {/* Side indicator bar */}
                                    <motion.div 
                                        className={`absolute left-0 top-0 bottom-0 w-1 ${sig.action === 'buy' ? 'bg-up' : 'bg-down'}`}
                                        initial={{ scaleY: 0 }}
                                        animate={{ scaleY: 1 }}
                                        transition={{ duration: 0.3, delay: i * 0.05 }}
                                    />
                                    
                                    <div className="flex justify-between items-start">
                                        <span className="text-[10px] font-black text-gray-400 uppercase tracking-tighter">
                                            {sig.strategy}
                                        </span>
                                        <span className="text-[8px] text-gray-600 font-mono tracking-tighter">
                                            {formatTime(sig.time)}
                                        </span>
                                    </div>
                                    <div className="flex justify-between items-center">
                                        <motion.span 
                                            className="text-xs font-black text-white"
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            transition={{ delay: 0.1 + i * 0.05 }}
                                        >
                                            {sig.symbol}
                                        </motion.span>
                                        <motion.span 
                                            className={`
                                                px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-widest
                                                ${sig.action === 'buy' ? 'bg-up/10 text-up' : 'bg-down/10 text-down'}
                                            `}
                                            initial={{ scale: 0 }}
                                            animate={{ scale: 1 }}
                                            transition={{ 
                                                delay: 0.15 + i * 0.05,
                                                type: 'spring',
                                                stiffness: 200,
                                            }}
                                        >
                                            {sig.action}
                                        </motion.span>
                                    </div>
                                    <motion.div 
                                        className="text-right text-[10px] font-mono text-gray-500 font-bold"
                                        initial={{ opacity: 0, x: 10 }}
                                        animate={{ opacity: 1, x: 0 }}
                                        transition={{ delay: 0.2 + i * 0.05 }}
                                    >
                                        PRC: {safeParseFloat(sig.price).toFixed(2)}
                                    </motion.div>
                                </motion.div>
                            );
                        })
                    )}
                </AnimatePresence>
            </div>
        </motion.div>
    );
};

export default SignalFeed;
