import React from 'react';
import { motion } from 'framer-motion';
import { cardHover, cardTap, fadeInUp } from '../../config/animations';

interface AnimatedCardProps {
  children: React.ReactNode;
  className?: string;
  delay?: number;
  onClick?: () => void;
  glowColor?: 'blue' | 'purple' | 'green' | 'red';
}

const glowColors = {
  blue: 'hover:shadow-[0_0_30px_rgba(59,130,246,0.4)]',
  purple: 'hover:shadow-[0_0_30px_rgba(139,92,246,0.4)]',
  green: 'hover:shadow-[0_0_30px_rgba(0,192,135,0.4)]',
  red: 'hover:shadow-[0_0_30px_rgba(255,59,48,0.4)]',
};

export const AnimatedCard: React.FC<AnimatedCardProps> = ({
  children,
  className = '',
  delay = 0,
  onClick,
  glowColor = 'blue',
}) => {
  return (
    <motion.div
      variants={fadeInUp}
      initial="hidden"
      animate="visible"
      whileHover={cardHover}
      whileTap={onClick ? cardTap : undefined}
      transition={{ delay }}
      onClick={onClick}
      className={`
        bg-card rounded-2xl border border-gray-800/50
        transition-shadow duration-300
        ${glowColors[glowColor]}
        ${onClick ? 'cursor-pointer' : ''}
        ${className}
      `}
    >
      {children}
    </motion.div>
  );
};

export default AnimatedCard;
