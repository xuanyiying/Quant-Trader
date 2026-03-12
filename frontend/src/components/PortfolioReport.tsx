import React from 'react';
import { PieChart } from 'lucide-react';
import type { PortfolioReport as PortfolioReportType } from '../types/data';

interface Props {
    report: PortfolioReportType | null;
    loading: boolean;
}

const PortfolioReport: React.FC<Props> = ({ report, loading }) => {
    if (loading) {
        return (
            <div className="bg-card p-6 rounded-2xl border border-gray-800 animate-pulse">
                <div className="h-4 w-32 bg-gray-800 rounded mb-4"></div>
                <div className="grid grid-cols-2 gap-4">
                    <div className="h-16 bg-gray-900/50 rounded-xl"></div>
                    <div className="h-16 bg-gray-900/50 rounded-xl"></div>
                </div>
            </div>
        );
    }

    if (!report) return null;

    const totalReturn = Number(report.total_return) || 0;
    const sharpeRatio = Number(report.sharpe_ratio) || 0;
    const maxDrawdown = Number(report.max_drawdown) || 0;
    const winRate = Number(report.win_rate) || 0;

    return (
        <div className="bg-card p-6 rounded-2xl shadow-xl border border-gray-800/50">
            <div className="flex items-center gap-2 mb-6">
                <PieChart size={18} className="text-blue-400" />
                <h2 className="font-black uppercase tracking-tight">Performance Analytics</h2>
            </div>

            <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-800">
                    <p className="text-[10px] text-gray-500 font-black uppercase mb-1">Total Return</p>
                    <p className={`text-lg font-black ${totalReturn >= 0 ? 'text-up' : 'text-down'}`}>
                        {totalReturn >= 0 ? '+' : ''}{totalReturn.toFixed(2)}%
                    </p>
                </div>
                <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-800">
                    <p className="text-[10px] text-gray-500 font-black uppercase mb-1">Sharpe Ratio</p>
                    <p className="text-lg font-black text-blue-400">{sharpeRatio.toFixed(2)}</p>
                </div>
                <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-800">
                    <p className="text-[10px] text-gray-500 font-black uppercase mb-1">Max Drawdown</p>
                    <p className="text-lg font-black text-down">-{maxDrawdown.toFixed(2)}%</p>
                </div>
                <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-800">
                    <p className="text-[10px] text-gray-500 font-black uppercase mb-1">Win Rate</p>
                    <p className="text-lg font-black text-white">{winRate.toFixed(1)}%</p>
                </div>
            </div>

            <div className="space-y-3">
                <div className="flex justify-between items-center">
                    <span className="text-[10px] text-gray-500 font-black uppercase">Trade Count</span>
                    <span className="text-xs font-mono text-gray-300">{report.total_trades} Executed</span>
                </div>
                <div className="h-1.5 w-full bg-gray-800 rounded-full overflow-hidden">
                    <div
                        className="h-full bg-blue-500 rounded-full"
                        style={{ width: `${Math.min(report.total_trades, 100)}%` }}
                    ></div>
                </div>
            </div>
        </div>
    );
};

export default PortfolioReport;
