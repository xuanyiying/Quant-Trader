import React, { useState } from 'react';
import { Lock, Mail, ArrowRight, Loader2, TrendingUp } from 'lucide-react';
import * as authApi from '../api/auth';
import { getErrorMessage } from '../utils/errorHandler';

interface AuthProps {
  onLogin: (token: string) => void;
}

const Auth: React.FC<AuthProps> = ({ onLogin }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      if (isLogin) {
        const response = await authApi.login(email, password);
        localStorage.setItem('token', response.token);
        onLogin(response.token);
      } else {
        await authApi.register(email, password);
        setIsLogin(true);
        setError('Registration successful! Please login.');
      }
    } catch (err) {
      setError(getErrorMessage(err) || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="max-w-md w-full space-y-8 bg-card p-8 rounded-2xl border border-gray-800 shadow-2xl relative overflow-hidden">
        {/* Background Decorative Element */}
        <div className="absolute -top-24 -right-24 w-48 h-48 bg-blue-600/10 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-24 -left-24 w-48 h-48 bg-up/10 rounded-full blur-3xl"></div>

        <div className="text-center relative">
          <div className="flex justify-center mb-4">
            <div className="p-3 bg-blue-600 rounded-xl shadow-lg shadow-blue-900/20">
              <TrendingUp size={32} className="text-white" />
            </div>
          </div>
          <h2 className="text-3xl font-black tracking-tight text-white uppercase">
            Quant<span className="text-blue-500">Trader</span>
          </h2>
          <p className="mt-2 text-sm text-gray-500 font-bold tracking-widest uppercase">
            {isLogin ? 'Welcome Back Terminal' : 'Create Operator Account'}
          </p>
        </div>

        <form className="mt-8 space-y-6 relative" onSubmit={handleSubmit}>
          {error && (
            <div className={`p-3 rounded-lg text-xs font-bold text-center border ${error.includes('successful') ? 'bg-up/10 border-up/20 text-up' : 'bg-down/10 border-down/20 text-down'
              }`}>
              {error}
            </div>
          )}

          <div className="space-y-4">
            <div className="relative group">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-500 transition-colors">
                <Mail size={18} />
              </div>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="block w-full pl-10 pr-3 py-3 bg-gray-900 border border-gray-800 rounded-xl text-sm font-bold text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all"
                placeholder="Email Address"
              />
            </div>

            <div className="relative group">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-500 group-focus-within:text-blue-500 transition-colors">
                <Lock size={18} />
              </div>
              <input
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="block w-full pl-10 pr-3 py-3 bg-gray-900 border border-gray-800 rounded-xl text-sm font-bold text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all"
                placeholder="Password"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={loading}
              className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-black rounded-xl text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-blue-900/20"
            >
              {loading ? (
                <Loader2 className="animate-spin" size={20} />
              ) : (
                <>
                  <span className="uppercase tracking-widest">{isLogin ? 'Sign In' : 'Register'}</span>
                  <ArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" size={18} />
                </>
              )}
            </button>
          </div>
        </form>

        <div className="text-center relative">
          <button
            onClick={() => setIsLogin(!isLogin)}
            className="text-xs font-black text-gray-500 hover:text-blue-400 uppercase tracking-widest transition-colors"
          >
            {isLogin ? "Don't have an account? Register" : "Already have an account? Sign In"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default Auth;
