# UI动画优化检查清单

## 依赖与配置
- [x] framer-motion 已安装并可用
- [x] gsap 和 @gsap/react 已安装并可用
- [x] canvas-confetti 已安装并可用
- [x] src/config/animations.ts 已创建并导出配置
- [x] src/styles/animations.css 已创建并在应用中引入
- [x] src/hooks/useAnimation.ts 已创建并可用

## 核心动画组件
- [x] AnimatedCard 组件实现悬停发光和缩放效果
- [x] AnimatedButton 组件实现涟漪和光晕效果
- [x] AnimatedNumber 组件实现数字滚动动画
- [x] FadeInView 组件实现滚动触发渐显
- [x] StaggerContainer 组件实现子元素交错动画

## 背景效果
- [x] AnimatedGrid 组件显示动态网格背景
- [x] GradientOrbs 组件显示流动渐变光晕
- [x] ParticleField 组件优化粒子效果性能

## 组件动画升级
- [x] LoadingScreen 显示品牌Logo动画和进度条
- [x] Header 组件有毛玻璃效果和滚动收缩
- [x] Header 状态指示器有脉冲呼吸动画
- [x] Header 余额使用数字滚动动画
- [x] Header 市场价格有涨跌闪烁效果
- [x] Chart 组件K线有绘制动画
- [x] Chart 工具提示有淡入和毛玻璃效果
- [x] TradingPanel 买卖按钮有光晕脉冲
- [x] TradingPanel 订单提交有加载动画
- [x] TradingPanel 成功订单有confetti效果
- [x] Auth 组件表单切换有3D翻转动画
- [x] Auth 输入验证失败有抖动反馈
- [x] SignalFeed 新信号有滑入动画
- [x] SignalFeed 高优先级信号有脉冲警示
- [x] StrategyMarketplace 策略卡片有悬停效果
- [x] PortfolioReport 数字有滚动动画
- [x] AlertsManager 警报项有进入/退出动画

## 全局效果
- [x] 页面切换有过渡动画
- [x] 组件加载有渐显动画
- [x] Skeleton 骨架屏有流动光泽效果
- [x] 支持 prefers-reduced-motion 媒体查询

## 性能与质量
- [x] 动画使用 will-change 优化
- [x] 移动端粒子效果性能良好
- [x] ESLint 检查无错误
- [x] TypeScript 类型检查无错误
- [x] 构建成功无警告
