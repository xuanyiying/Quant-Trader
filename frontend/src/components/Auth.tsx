import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Lock, Mail, ArrowRight, Loader2, TrendingUp, User, Eye, EyeOff, Sparkles, Zap, Shield } from 'lucide-react';
import * as authApi from '../api/auth';
import { getErrorMessage } from '../utils/errorHandler';

interface AuthProps {
  onLogin: (token: string) => void;
}

interface Particle {
  id: number;
  x: number;
  y: number;
  vx: number;
  vy: number;
  size: number;
  opacity: number;
  color: string;
}

// 模拟股票数据
const STOCK_TICKERS = [
  { symbol: 'BTC', price: 67234.56, change: 2.34 },
  { symbol: 'ETH', price: 3456.78, change: -1.23 },
  { symbol: 'SOL', price: 178.90, change: 5.67 },
  { symbol: 'AAPL', price: 189.45, change: 0.89 },
  { symbol: 'TSLA', price: 234.56, change: -2.45 },
  { symbol: 'NVDA', price: 876.54, change: 3.21 },
  { symbol: 'MSFT', price: 423.78, change: 1.12 },
  { symbol: 'GOOGL', price: 156.78, change: -0.67 },
];

const Auth: React.FC<AuthProps> = ({ onLogin }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [isFlipping, setIsFlipping] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [username, setUsername] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [focusedField, setFocusedField] = useState<string | null>(null);

  // 粒子动画状态
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const particlesRef = useRef<Particle[]>([]);
  const animationRef = useRef<number | null>(null);
  const mouseRef = useRef({ x: 0, y: 0 });

  // 初始化粒子
  const initParticles = useCallback(() => {
    const particles: Particle[] = [];
    const colors = ['#3b82f6', '#00c087', '#8b5cf6', '#06b6d4', '#f59e0b'];

    for (let i = 0; i < 60; i++) {
      particles.push({
        id: i,
        x: Math.random() * window.innerWidth,
        y: Math.random() * window.innerHeight,
        vx: (Math.random() - 0.5) * 0.8,
        vy: (Math.random() - 0.5) * 0.8,
        size: Math.random() * 3 + 1,
        opacity: Math.random() * 0.5 + 0.2,
        color: colors[Math.floor(Math.random() * colors.length)],
      });
    }
    particlesRef.current = particles;
  }, []);

  // 动画循环
  const animate = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    particlesRef.current.forEach((particle, index) => {
      // 更新位置
      particle.x += particle.vx;
      particle.y += particle.vy;

      // 边界检测
      if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1;
      if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1;

      // 鼠标交互
      const dx = mouseRef.current.x - particle.x;
      const dy = mouseRef.current.y - particle.y;
      const distance = Math.sqrt(dx * dx + dy * dy);

      if (distance < 150) {
        const force = (150 - distance) / 150;
        particle.vx -= (dx / distance) * force * 0.02;
        particle.vy -= (dy / distance) * force * 0.02;
      }

      // 绘制粒子
      ctx.beginPath();
      ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
      ctx.fillStyle = particle.color;
      ctx.globalAlpha = particle.opacity;
      ctx.fill();

      // 绘制连线
      particlesRef.current.slice(index + 1).forEach((other) => {
        const dx2 = particle.x - other.x;
        const dy2 = particle.y - other.y;
        const dist = Math.sqrt(dx2 * dx2 + dy2 * dy2);

        if (dist < 120) {
          ctx.beginPath();
          ctx.moveTo(particle.x, particle.y);
          ctx.lineTo(other.x, other.y);
          ctx.strokeStyle = particle.color;
          ctx.globalAlpha = (1 - dist / 120) * 0.2;
          ctx.lineWidth = 0.5;
          ctx.stroke();
        }
      });
    });

    ctx.globalAlpha = 1;
    animationRef.current = requestAnimationFrame(animate);
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const handleResize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };

    const handleMouseMove = (e: MouseEvent) => {
      mouseRef.current = { x: e.clientX, y: e.clientY };
    };

    handleResize();
    initParticles();
    animate();

    window.addEventListener('resize', handleResize);
    window.addEventListener('mousemove', handleMouseMove);

    return () => {
      window.removeEventListener('resize', handleResize);
      window.removeEventListener('mousemove', handleMouseMove);
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [initParticles, animate]);

  const handleToggleMode = () => {
    setIsFlipping(true);
    setTimeout(() => {
      setIsLogin(!isLogin);
      setError('');
      setSuccess('');
      setEmail('');
      setPassword('');
      setConfirmPassword('');
      setUsername('');
    }, 300);
    setTimeout(() => {
      setIsFlipping(false);
    }, 600);
  };

  const validateForm = () => {
    if (!email || !password) {
      setError('Please fill in all required fields');
      return false;
    }

    if (!isLogin) {
      if (password !== confirmPassword) {
        setError('Passwords do not match');
        return false;
      }
      if (password.length < 6) {
        setError('Password must be at least 6 characters');
        return false;
      }
      if (!username.trim()) {
        setError('Please enter a username');
        return false;
      }
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      if (isLogin) {
        const response = await authApi.login(email, password);
        localStorage.setItem('token', response.token);
        onLogin(response.token);
      } else {
        await authApi.register(email, password);
        setSuccess('Registration successful! Please login.');
        setTimeout(() => {
          setIsLogin(true);
          setPassword('');
          setConfirmPassword('');
        }, 1500);
      }
    } catch (err) {
      setError(getErrorMessage(err) || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0a0f] flex items-center justify-center p-4 relative overflow-hidden">
      {/* 粒子背景画布 */}
      <canvas
        ref={canvasRef}
        className="absolute inset-0 pointer-events-none"
        style={{ zIndex: 0 }}
      />

      {/* 渐变背景 */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-900/20 via-purple-900/10 to-cyan-900/20 pointer-events-none" />

      {/* 股票行情滚动条 */}
      <div className="absolute top-0 left-0 right-0 h-10 bg-black/40 backdrop-blur-sm border-b border-white/5 overflow-hidden z-10">
        <div className="flex animate-marquee whitespace-nowrap">
          {[...STOCK_TICKERS, ...STOCK_TICKERS].map((stock, index) => (
            <div key={index} className="flex items-center mx-6 text-sm">
              <span className="font-bold text-gray-400 mr-2">{stock.symbol}</span>
              <span className="text-white mr-2">${stock.price.toLocaleString()}</span>
              <span className={stock.change >= 0 ? 'text-up' : 'text-down'}>
                {stock.change >= 0 ? '+' : ''}{stock.change}%
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* 装饰性光晕 */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-blue-600/20 rounded-full blur-[120px] pointer-events-none animate-pulse" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-600/20 rounded-full blur-[120px] pointer-events-none animate-pulse" style={{ animationDelay: '1s' }} />
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-cyan-600/10 rounded-full blur-[150px] pointer-events-none" />

      {/* 主卡片容器 */}
      <div className="relative z-10 perspective-1000">
        <div
          className={`relative transition-all duration-600 transform-style-preserve-3d ${isFlipping ? 'rotate-y-180' : ''
            }`}
          style={{
            transformStyle: 'preserve-3d',
            transform: isFlipping ? 'rotateY(180deg)' : 'rotateY(0deg)',
            transition: 'transform 0.6s cubic-bezier(0.4, 0, 0.2, 1)',
          }}
        >
          {/* 正面 - 登录/注册表单 */}
          <div
            className="w-[420px] bg-gradient-to-b from-gray-900/90 to-black/90 backdrop-blur-xl rounded-3xl border border-white/10 shadow-2xl overflow-hidden"
            style={{ backfaceVisibility: 'hidden' }}
          >
            {/* 顶部发光条 */}
            <div className="h-1 w-full bg-gradient-to-r from-blue-500 via-purple-500 to-cyan-500" />

            {/* Logo区域 */}
            <div className="pt-8 pb-6 text-center relative">
              <div className="flex justify-center mb-4">
                <div className="relative">
                  <div className="absolute inset-0 bg-blue-500/50 rounded-2xl blur-xl animate-pulse" />
                  <div className="relative p-4 bg-gradient-to-br from-blue-600 to-purple-600 rounded-2xl shadow-lg">
                    <TrendingUp size={36} className="text-white" />
                  </div>
                  <div className="absolute -top-1 -right-1">
                    <Sparkles size={16} className="text-yellow-400 animate-bounce" />
                  </div>
                </div>
              </div>
              <h1 className="text-3xl font-black tracking-tight">
                <span className="bg-gradient-to-r from-blue-400 via-purple-400 to-cyan-400 bg-clip-text text-transparent">
                  QuantTrader
                </span>
              </h1>
              <p className="mt-2 text-sm text-gray-500 font-medium">
                {isLogin ? 'Welcome back, trader' : 'Join the elite traders'}
              </p>
            </div>

            {/* 表单区域 */}
            <div className="px-8 pb-8">
              {error && (
                <div className="mb-4 p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm text-center animate-shake">
                  {error}
                </div>
              )}

              {success && (
                <div className="mb-4 p-3 rounded-xl bg-green-500/10 border border-green-500/20 text-green-400 text-sm text-center animate-fade-in">
                  {success}
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-4">
                {/* 用户名输入 - 仅注册时显示 */}
                {!isLogin && (
                  <div className="relative group">
                    <div className={`absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-xl blur opacity-0 transition-opacity duration-300 ${focusedField === 'username' ? 'opacity-100' : ''}`} />
                    <div className="relative">
                      <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-400 transition-colors">
                        <User size={18} />
                      </div>
                      <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        onFocus={() => setFocusedField('username')}
                        onBlur={() => setFocusedField(null)}
                        className="block w-full pl-11 pr-4 py-3.5 bg-gray-800/50 border border-gray-700 rounded-xl text-sm text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800/80 transition-all duration-300"
                        placeholder="Username"
                      />
                    </div>
                  </div>
                )}

                {/* 邮箱输入 */}
                <div className="relative group">
                  <div className={`absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-xl blur opacity-0 transition-opacity duration-300 ${focusedField === 'email' ? 'opacity-100' : ''}`} />
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-400 transition-colors">
                      <Mail size={18} />
                    </div>
                    <input
                      type="email"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      onFocus={() => setFocusedField('email')}
                      onBlur={() => setFocusedField(null)}
                      className="block w-full pl-11 pr-4 py-3.5 bg-gray-800/50 border border-gray-700 rounded-xl text-sm text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800/80 transition-all duration-300"
                      placeholder="Email address"
                    />
                  </div>
                </div>

                {/* 密码输入 */}
                <div className="relative group">
                  <div className={`absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-xl blur opacity-0 transition-opacity duration-300 ${focusedField === 'password' ? 'opacity-100' : ''}`} />
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-400 transition-colors">
                      <Lock size={18} />
                    </div>
                    <input
                      type={showPassword ? 'text' : 'password'}
                      required
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      onFocus={() => setFocusedField('password')}
                      onBlur={() => setFocusedField(null)}
                      className="block w-full pl-11 pr-11 py-3.5 bg-gray-800/50 border border-gray-700 rounded-xl text-sm text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800/80 transition-all duration-300"
                      placeholder="Password"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-gray-300 transition-colors"
                    >
                      {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                    </button>
                  </div>
                </div>

                {/* 确认密码 - 仅注册时显示 */}
                {!isLogin && (
                  <div className="relative group">
                    <div className={`absolute inset-0 bg-gradient-to-r from-blue-500/20 to-purple-500/20 rounded-xl blur opacity-0 transition-opacity duration-300 ${focusedField === 'confirmPassword' ? 'opacity-100' : ''}`} />
                    <div className="relative">
                      <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-400 transition-colors">
                        <Shield size={18} />
                      </div>
                      <input
                        type={showConfirmPassword ? 'text' : 'password'}
                        required
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        onFocus={() => setFocusedField('confirmPassword')}
                        onBlur={() => setFocusedField(null)}
                        className="block w-full pl-11 pr-11 py-3.5 bg-gray-800/50 border border-gray-700 rounded-xl text-sm text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800/80 transition-all duration-300"
                        placeholder="Confirm password"
                      />
                      <button
                        type="button"
                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                        className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-gray-300 transition-colors"
                      >
                        {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                      </button>
                    </div>
                  </div>
                )}

                {/* 提交按钮 */}
                <button
                  type="submit"
                  disabled={loading}
                  className="group relative w-full mt-6"
                >
                  <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl blur opacity-50 group-hover:opacity-75 transition-opacity" />
                  <div className="relative flex items-center justify-center py-3.5 px-4 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl text-white font-semibold shadow-lg shadow-blue-900/25 hover:shadow-blue-900/40 transition-all duration-300 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed">
                    {loading ? (
                      <Loader2 className="animate-spin" size={20} />
                    ) : (
                      <>
                        <span>{isLogin ? 'Sign In' : 'Create Account'}</span>
                        <ArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" size={18} />
                      </>
                    )}
                  </div>
                </button>
              </form>

              {/* 切换登录/注册 */}
              <div className="mt-6 text-center">
                <button
                  type="button"
                  onClick={handleToggleMode}
                  className="text-sm text-gray-500 hover:text-white transition-colors"
                >
                  {isLogin ? (
                    <span>
                      Don't have an account?{' '}
                      <span className="text-blue-400 font-semibold hover:text-blue-300">Sign up</span>
                    </span>
                  ) : (
                    <span>
                      Already have an account?{' '}
                      <span className="text-blue-400 font-semibold hover:text-blue-300">Sign in</span>
                    </span>
                  )}
                </button>
              </div>

              {/* 社交登录选项 */}
              <div className="mt-6">
                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <div className="w-full border-t border-gray-800" />
                  </div>
                  <div className="relative flex justify-center text-xs">
                    <span className="px-2 bg-gray-900 text-gray-500">Or continue with</span>
                  </div>
                </div>

                <div className="mt-4 grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    className="flex items-center justify-center py-2.5 px-4 bg-gray-800/50 border border-gray-700 rounded-xl text-gray-400 hover:text-white hover:bg-gray-800 hover:border-gray-600 transition-all duration-300"
                  >
                    <Zap size={18} className="mr-2" />
                    <span className="text-sm">Google</span>
                  </button>
                  <button
                    type="button"
                    className="flex items-center justify-center py-2.5 px-4 bg-gray-800/50 border border-gray-700 rounded-xl text-gray-400 hover:text-white hover:bg-gray-800 hover:border-gray-600 transition-all duration-300"
                  >
                    <svg className="w-[18px] h-[18px] mr-2" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
                    </svg>
                    <span className="text-sm">GitHub</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 底部信息 */}
      <div className="absolute bottom-6 left-0 right-0 text-center z-10">
        <p className="text-xs text-gray-600">
          By continuing, you agree to our Terms of Service and Privacy Policy
        </p>
      </div>
    </div>
  );
};

export default Auth;
