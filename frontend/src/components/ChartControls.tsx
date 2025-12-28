import React from 'react';
import { RefreshCw, Settings, Database } from 'lucide-react';
import { TRADING_CONFIG } from '../constants';

interface ChartControlsProps {
    symbol: string;
    period: string;
    onSymbolChange: (symbol: string) => void;
    onPeriodChange: (period: string) => void;
    onRefresh: () => void;
    onBackfill: () => void;
}

const ChartControls: React.FC<ChartControlsProps> = ({
    symbol,
    period,
    onSymbolChange,
    onPeriodChange,
    onRefresh,
    onBackfill,
}) => {
    const [inputSymbol, setInputSymbol] = React.useState(symbol);

    const handleUpdateSymbol = () => {
        onSymbolChange(inputSymbol.toUpperCase());
    };

    return (
        <div className="flex flex-wrap justify-between items-center gap-4 mb-8">
            <div className="flex items-center gap-2 bg-gray-900 px-2 py-1.5 rounded-xl border border-gray-800">
                <input
                    value={inputSymbol}
                    onChange={(e) => setInputSymbol(e.target.value)}
                    className="bg-transparent border-none outline-none px-3 py-1 text-sm font-black w-28 text-white placeholder-gray-700"
                    placeholder="BTCUSDT"
                />
                <button
                    onClick={handleUpdateSymbol}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-1.5 rounded-lg text-[10px] font-black transition-all active:scale-95 shadow-lg shadow-blue-900/20"
                >
                    SYNC
                </button>
            </div>

            <div className="flex gap-1.5 bg-gray-900 p-1 rounded-xl border border-gray-800">
                {TRADING_CONFIG.PERIODS.map((p) => (
                    <button
                        key={p}
                        onClick={() => onPeriodChange(p)}
                        className={`px-4 py-1.5 rounded-lg text-[10px] font-black transition-all ${period === p
                                ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20'
                                : 'text-gray-500 hover:text-gray-300 hover:bg-gray-800'
                            }`}
                    >
                        {p.toUpperCase()}
                    </button>
                ))}
            </div>

            <div className="flex gap-2">
                <button
                    onClick={onRefresh}
                    className="bg-gray-900 hover:bg-gray-800 text-gray-400 p-2.5 rounded-xl border border-gray-800 transition-all font-bold"
                    title="Refresh"
                >
                    <RefreshCw size={18} />
                </button>
                <button
                    onClick={onBackfill}
                    className="bg-gray-900 hover:bg-gray-800 text-gray-400 p-2.5 rounded-xl border border-gray-800 transition-all font-bold"
                    title="Backfill"
                >
                    <Database size={18} />
                </button>
                <button
                    className="bg-gray-900 hover:bg-gray-800 text-gray-400 p-2.5 rounded-xl border border-gray-800 transition-all font-bold"
                    title="Settings"
                >
                    <Settings size={18} />
                </button>
            </div>
        </div>
    );
};

export default ChartControls;
