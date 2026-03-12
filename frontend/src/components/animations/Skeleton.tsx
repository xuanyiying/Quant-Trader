import React from 'react';
import { motion } from 'framer-motion';

interface SkeletonProps {
  className?: string;
  width?: string | number;
  height?: string | number;
  circle?: boolean;
  shimmer?: boolean;
}

export const Skeleton: React.FC<SkeletonProps> = ({
  className = '',
  width = '100%',
  height = '20px',
  circle = false,
  shimmer = true,
}) => {
  return (
    <motion.div
      className={`
        bg-gray-800/50 overflow-hidden
        ${circle ? 'rounded-full' : 'rounded-lg'}
        ${shimmer ? 'relative' : ''}
        ${className}
      `}
      style={{ width, height }}
      initial={{ opacity: 0.5 }}
      animate={{ opacity: [0.5, 0.8, 0.5] }}
      transition={{
        duration: 1.5,
        repeat: Infinity,
        ease: 'easeInOut',
      }}
    >
      {shimmer && (
        <motion.div
          className="absolute inset-0"
          style={{
            background: 'linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.05) 50%, transparent 100%)',
            backgroundSize: '200% 100%',
          }}
          animate={{
            backgroundPosition: ['200% 0', '-200% 0'],
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
            ease: 'linear',
          }}
        />
      )}
    </motion.div>
  );
};

interface SkeletonCardProps {
  className?: string;
  lines?: number;
  hasImage?: boolean;
}

export const SkeletonCard: React.FC<SkeletonCardProps> = ({
  className = '',
  lines = 3,
  hasImage = false,
}) => {
  return (
    <div className={`bg-card rounded-2xl p-6 border border-gray-800/50 ${className}`}>
      {hasImage && (
        <Skeleton width="100%" height="120px" className="mb-4 rounded-xl" />
      )}
      <Skeleton width="60%" height="24px" className="mb-4" />
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton
          key={i}
          width={i === lines - 1 ? '80%' : '100%'}
          height="16px"
          className="mb-3"
        />
      ))}
    </div>
  );
};

export default Skeleton;
