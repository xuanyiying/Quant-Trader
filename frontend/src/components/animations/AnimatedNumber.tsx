import React, { useEffect, useState } from 'react';
import { motion, useReducedMotion } from 'framer-motion';
import { useAnimatedNumber } from '../../hooks/useAnimation';

interface AnimatedNumberProps {
  value: number;
  prefix?: string;
  suffix?: string;
  decimals?: number;
  duration?: number;
  className?: string;
  flashOnChange?: boolean;
  flashColor?: 'green' | 'red' | 'none';
}

export const AnimatedNumber: React.FC<AnimatedNumberProps> = ({
  value,
  prefix = '',
  suffix = '',
  decimals = 2,
  duration = 1000,
  className = '',
  flashOnChange = true,
  flashColor = 'none',
}) => {
  const displayValue = useAnimatedNumber(value, { duration, decimals });
  const [isFlashing, setIsFlashing] = useState(false);
  const [flashDirection, setFlashDirection] = useState<'green' | 'red' | null>(null);
  const [prevValue, setPrevValue] = useState(value);
  const prefersReducedMotion = useReducedMotion();

  useEffect(() => {
    if (!flashOnChange || prefersReducedMotion) return;
    
    if (value !== prevValue) {
      const newFlashColor = flashColor === 'none' 
        ? (value > prevValue ? 'green' : 'red')
        : flashColor;
      
      if (newFlashColor === 'green' || newFlashColor === 'red') {
        setFlashDirection(newFlashColor);
        setIsFlashing(true);
        const timer = setTimeout(() => setIsFlashing(false), 500);
        setPrevValue(value);
        return () => clearTimeout(timer);
      }
      setPrevValue(value);
    }
  }, [value, flashOnChange, flashColor, prevValue, prefersReducedMotion]);

  const flashClass = isFlashing
    ? flashDirection === 'green'
      ? 'animate-flash-green'
      : 'animate-flash-red'
    : '';

  return (
    <motion.span
      className={`inline-block ${flashClass} ${className}`}
      initial={false}
      animate={isFlashing ? { scale: [1, 1.05, 1] } : {}}
      transition={{ duration: 0.3 }}
    >
      {prefix}{displayValue}{suffix}
    </motion.span>
  );
};

export default AnimatedNumber;
