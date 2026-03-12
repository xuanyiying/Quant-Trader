import React, { useEffect, useRef } from 'react';

interface AnimatedGridProps {
  className?: string;
  gridSize?: number;
  lineColor?: string;
  glowColor?: string;
}

export const AnimatedGrid: React.FC<AnimatedGridProps> = ({
  className = '',
  gridSize = 50,
  lineColor = 'rgba(59, 130, 246, 0.1)',
  glowColor = 'rgba(59, 130, 246, 0.3)',
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationRef = useRef<number | undefined>(undefined);
  const timeRef = useRef(0);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const resize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };

    resize();
    window.addEventListener('resize', resize);

    const drawGrid = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      const cols = Math.ceil(canvas.width / gridSize);
      const rows = Math.ceil(canvas.height / gridSize);

      // Draw vertical lines
      for (let i = 0; i <= cols; i++) {
        const x = i * gridSize;
        const wave = Math.sin(timeRef.current * 0.001 + i * 0.1) * 2;
        
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x + wave, canvas.height);
        ctx.strokeStyle = lineColor;
        ctx.lineWidth = 1;
        ctx.stroke();

        // Draw glow at intersections
        for (let j = 0; j <= rows; j += 3) {
          const y = j * gridSize;
          const glowIntensity = (Math.sin(timeRef.current * 0.002 + i * 0.2 + j * 0.1) + 1) / 2;
          
          ctx.beginPath();
          ctx.arc(x + wave, y, 2 + glowIntensity * 2, 0, Math.PI * 2);
          ctx.fillStyle = glowColor.replace('0.3', (glowIntensity * 0.5).toString());
          ctx.fill();
        }
      }

      // Draw horizontal lines
      for (let j = 0; j <= rows; j++) {
        const y = j * gridSize;
        const wave = Math.cos(timeRef.current * 0.001 + j * 0.1) * 2;
        
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(canvas.width, y + wave);
        ctx.strokeStyle = lineColor;
        ctx.lineWidth = 1;
        ctx.stroke();
      }

      timeRef.current += 16;
      animationRef.current = requestAnimationFrame(drawGrid);
    };

    drawGrid();

    return () => {
      window.removeEventListener('resize', resize);
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [gridSize, lineColor, glowColor]);

  return (
    <canvas
      ref={canvasRef}
      className={`fixed inset-0 pointer-events-none z-0 ${className}`}
    />
  );
};

export default AnimatedGrid;
