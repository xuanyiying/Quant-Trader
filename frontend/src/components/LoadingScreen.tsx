import React from 'react';
import { Activity } from 'lucide-react';
import { UI_CONFIG } from '../constants';

const LoadingScreen: React.FC = () => {
    return (
        <div className="min-h-screen bg-background flex items-center justify-center">
            <div className="flex flex-col items-center gap-4">
                <Activity size={UI_CONFIG.LOADING_SPINNER_SIZE} className="text-blue-600 animate-spin" />
                <p className="text-gray-500 font-bold uppercase tracking-widest animate-pulse">
                    Initializing Terminal...
                </p>
            </div>
        </div>
    );
};

export default LoadingScreen;
