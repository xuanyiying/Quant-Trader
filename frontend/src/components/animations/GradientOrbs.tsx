import React from 'react';
import { motion } from 'framer-motion';

interface GradientOrbsProps {
  className?: string;
}

export const GradientOrbs: React.FC<GradientOrbsProps> = ({ className = '' }) => {
  return (
    <div className={`fixed inset-0 overflow-hidden pointer-events-none z-0 ${className}`}>
      {/* Blue orb */}
      <motion.div
        className="absolute w-[600px] h-[600px] rounded-full"
        style={{
          background: 'radial-gradient(circle, rgba(59,130,246,0.15) 0%, transparent 70%)',
          filter: 'blur(80px)',
        }}
        animate={{
          x: ['-20%', '30%', '-20%'],
          y: ['-10%', '40%', '-10%'],
          scale: [1, 1.2, 1],
        }}
        transition={{
          duration: 20,
          repeat: Infinity,
          ease: 'easeInOut',
        }}
      />

      {/* Purple orb */}
      <motion.div
        className="absolute w-[500px] h-[500px] rounded-full"
        style={{
          background: 'radial-gradient(circle, rgba(139,92,246,0.12) 0%, transparent 70%)',
          filter: 'blur(80px)',
          right: '-10%',
          top: '20%',
        }}
        animate={{
          x: ['10%', '-20%', '10%'],
          y: ['0%', '30%', '0%'],
          scale: [1, 1.3, 1],
        }}
        transition={{
          duration: 25,
          repeat: Infinity,
          ease: 'easeInOut',
          delay: 5,
        }}
      />

      {/* Cyan orb */}
      <motion.div
        className="absolute w-[400px] h-[400px] rounded-full"
        style={{
          background: 'radial-gradient(circle, rgba(6,182,212,0.1) 0%, transparent 70%)',
          filter: 'blur(60px)',
          left: '30%',
          bottom: '10%',
        }}
        animate={{
          x: ['-10%', '20%', '-10%'],
          y: ['10%', '-20%', '10%'],
          scale: [1, 1.1, 1],
        }}
        transition={{
          duration: 18,
          repeat: Infinity,
          ease: 'easeInOut',
          delay: 3,
        }}
      />

      {/* Green orb - subtle */}
      <motion.div
        className="absolute w-[300px] h-[300px] rounded-full"
        style={{
          background: 'radial-gradient(circle, rgba(0,192,135,0.08) 0%, transparent 70%)',
          filter: 'blur(60px)',
          right: '20%',
          bottom: '30%',
        }}
        animate={{
          x: ['5%', '-15%', '5%'],
          y: ['-5%', '15%', '-5%'],
          scale: [1, 1.2, 1],
        }}
        transition={{
          duration: 22,
          repeat: Infinity,
          ease: 'easeInOut',
          delay: 8,
        }}
      />
    </div>
  );
};

export default GradientOrbs;
