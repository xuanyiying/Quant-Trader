import React from 'react';
import { TrendingUp, Activity, User, LogOut, RefreshCw } from 'lucide-react';
import type { PaperAccount, Subscription } from '../types/data';
import { formatCurrency } from '../utils/formatters';
import { SUBSCRIPTION_TIERS } from '../constants';

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

    return (
        <header className="flex flex-col md:flex-row justify-between items-center gap-6 bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50 backdrop-blur-sm relative overflow-hidden group">
            <div className="absolute inset-0 bg-liner-to-r from-blue-600/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>

            <div className="flex items-center gap-4 relative">
                <div className="p-3 bg-blue-600 rounded-xl shadow-lg shadow-blue-900/20">
                    <TrendingUp size={28} className="text-white" />
                </div>
                <div>
                    <h1 className="text-2xl font-black tracking-tight text-white">
                        QUANT<span className="text-blue-500">TRADER</span>
                    </h1>
                    <div className="flex items-center gap-2">
                        <p className="text-[10px] text-gray-500 uppercase tracking-widest font-black">
                            {subscription ? `${subscription.tier_name} MEMBER` : 'LOADING...'}
                        </p>
                        {isPro ? (
                            <span className="bg-yellow-500/10 text-yellow-500 text-[8px] px-1.5 py-0.5 rounded font-black border border-yellow-500/20 uppercase">
                                Pro Access
                            </span>
                        ) : (
                            <button
                                onClick={onUpgrade}
                                className="bg-blue-600/10 text-blue-400 text-[8px] px-1.5 py-0.5 rounded font-black border border-blue-500/20 hover:bg-blue-600 hover:text-white transition-all uppercase"
                            >
                                Upgrade Now
                            </button>
                        )}
                    </div>
                </div>
            </div>

            <div className="flex flex-wrap justify-center gap-4 relative">
                <div className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px] group relative">
                    <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1 text-center">
                        Paper Balance
                    </span>
                    <span className="text-lg font-black text-up flex items-center gap-2">
                        ${paperAccount ? formatCurrency(paperAccount.balance) : '0.00'}
                    </span>
                    <button
                        onClick={onResetAccount}
                        className="absolute -top-2 -right-2 bg-gray-800 hover:bg-red-900/40 text-gray-500 hover:text-red-400 p-1.5 rounded-lg border border-gray-700 opacity-0 group-hover:opacity-100 transition-all shadow-xl"
                        title="Reset Account"
                    >
                        <RefreshCw size={12} />
                    </button>
                </div>

                <div className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px]">
                    <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1">
                        Status
                    </span>
                    <span className={`text-sm font-black flex items-center gap-2 ${connectionStatus === 'Connected' ? 'text-up' : 'text-down'
                        }`}>
                        <Activity size={16} />
                        {connectionStatus.toUpperCase()}
                    </span>
                </div>

                <div className="bg-gray-900/50 px-5 py-3 rounded-xl border border-gray-800 flex flex-col items-center min-w-[140px]">
                    <span className="text-[10px] text-gray-500 uppercase font-black tracking-tighter mb-1">
                        Market Price
                    </span>
                    <span className="text-lg font-black text-blue-400">
                        {lastTradePrice || '--.--'} <span className="text-[10px] ml-1 text-gray-600 font-bold">USDT</span>
                    </span>
                </div>
            </div>

            <div className="flex items-center gap-3 relative">
                <button className="bg-gray-900 hover:bg-gray-800 text-gray-400 p-2.5 rounded-xl border border-gray-800 transition-all font-bold group/btn">
                    <User size={20} className="group-hover/btn:text-blue-400 transition-colors" />
                </button>
                <button
                    onClick={onLogout}
                    className="bg-gray-900 hover:bg-red-900/20 text-gray-400 hover:text-red-400 p-2.5 rounded-xl border border-gray-800 hover:border-red-900/30 transition-all font-bold group/logout"
                    title="Logout"
                >
                    <LogOut size={20} />
                </button>
            </div>
        </header>
    );
};

export default Header;
