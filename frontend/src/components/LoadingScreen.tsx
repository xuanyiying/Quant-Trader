import React from 'react';
import { motion } from 'framer-motion';
import { TrendingUp } from 'lucide-react';
import { GradientOrbs } from './animations';

interface LoadingScreenProps {
  progress?: number;
}

const LoadingScreen: React.FC<LoadingScreenProps> = ({ progress = 0 }) => {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center relative overflow-hidden">
      <GradientOrbs />
      
      <div className="relative z-10 flex flex-col items-center gap-8">
        {/* Logo Animation */}
        <motion.div
          className="relative"
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.5, ease: 'easeOut' }}
        >
          {/* Glow effect */}
          <motion.div
            className="absolute inset-0 bg-blue-500/30 rounded-3xl blur-3xl"
            animate={{
              scale: [1, 1.2, 1],
              opacity: [0.3, 0.5, 0.3],
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: 'easeInOut',
            }}
          />
          
          {/* Logo container */}
          <motion.div
            className="relative p-6 bg-gradient-to-br from-blue-600 to-purple-600 rounded-3xl shadow-2xl"
            animate={{
              boxShadow: [
                '0 0 40px rgba(59, 130, 246, 0.3)',
                '0 0 60px rgba(59, 130, 246, 0.5)',
                '0 0 40px rgba(59, 130, 246, 0.3)',
              ],
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: 'easeInOut',
            }}
          >
            <TrendingUp size={48} className="text-white" />
          </motion.div>
        </motion.div>

        {/* Brand name */}
        <motion.div
          className="text-center"
          initial={{ y: 20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
        >
          <h1 className="text-3xl font-black tracking-tight">
            <span className="text-white">QUANT</span>
            <span className="bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              TRADER
            </span>
          </h1>
        </motion.div>

        {/* Progress bar */}
        <motion.div
          className="w-64 h-1 bg-gray-800 rounded-full overflow-hidden"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4 }}
        >
          <motion.div
            className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full"
            initial={{ width: 0 }}
            animate={{ width: `${progress}%` }}
            transition={{ duration: 0.3 }}
          />
        </motion.div>

        {/* Loading text */}
        <motion.p
          className="text-gray-500 font-bold uppercase tracking-widest text-sm"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{
            duration: 1.5,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        >
          Initializing Terminal...
        </motion.p>
      </div>
    </div>
  );
};

export default LoadingScreen;
