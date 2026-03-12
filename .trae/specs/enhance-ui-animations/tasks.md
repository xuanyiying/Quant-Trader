# UI动画优化任务列表

## 依赖安装与基础配置
- [x] Task 1: 安装动画库依赖
  - [x] SubTask 1.1: 安装 framer-motion
  - [x] SubTask 1.2: 安装 gsap 和 @gsap/react
  - [x] SubTask 1.3: 安装 canvas-confetti (用于庆祝效果)
  - [x] SubTask 1.4: 更新 package.json 确认依赖版本

- [x] Task 2: 创建动画配置文件
  - [x] SubTask 2.1: 创建 src/config/animations.ts (动画配置常量)
  - [x] SubTask 2.2: 创建 src/styles/animations.css (全局CSS动画)
  - [x] SubTask 2.3: 创建 src/hooks/useAnimation.ts (动画Hooks)

## 核心动画组件开发
- [x] Task 3: 创建动画组件库
  - [x] SubTask 3.1: 创建 AnimatedCard 组件 (带悬停效果的卡片)
  - [x] SubTask 3.2: 创建 AnimatedButton 组件 (带涟漪和光晕效果)
  - [x] SubTask 3.3: 创建 AnimatedNumber 组件 (数字滚动动画)
  - [x] SubTask 3.4: 创建 FadeInView 组件 (滚动触发渐显)
  - [x] SubTask 3.5: 创建 StaggerContainer 组件 (子元素交错动画)

- [x] Task 4: 背景效果组件
  - [x] SubTask 4.1: 创建 AnimatedGrid 组件 (动态网格背景)
  - [x] SubTask 4.2: 创建 GradientOrbs 组件 (流动渐变光晕)
  - [x] SubTask 4.3: 创建 ParticleField 组件 (优化版粒子效果)

## 现有组件动画升级
- [x] Task 5: LoadingScreen 组件升级
  - [x] SubTask 5.1: 添加品牌Logo动画
  - [x] SubTask 5.2: 添加进度条流动效果
  - [x] SubTask 5.3: 添加背景渐变动画

- [x] Task 6: Header 组件动画升级
  - [x] SubTask 6.1: 添加毛玻璃效果 + 滚动收缩
  - [x] SubTask 6.2: 状态指示器添加脉冲呼吸动画
  - [x] SubTask 6.3: 余额数字改用 AnimatedNumber 组件
  - [x] SubTask 6.4: 市场价格的涨跌闪烁动画

- [x] Task 7: Chart 组件动画升级
  - [x] SubTask 7.1: K线数据加载时的绘制动画
  - [x] SubTask 7.2: 工具提示淡入 + 毛玻璃背景
  - [x] SubTask 7.3: 数据更新时的平滑过渡

- [x] Task 8: TradingPanel 组件动画升级
  - [x] SubTask 8.1: 买卖按钮添加光晕脉冲效果
  - [x] SubTask 8.2: 订单提交时的加载动画
  - [x] SubTask 8.3: 持仓列表项进入动画
  - [x] SubTask 8.4: 成功订单的 confetti 庆祝效果

- [x] Task 9: Auth 组件动画优化
  - [x] SubTask 9.1: 优化粒子效果性能
  - [x] SubTask 9.2: 表单切换3D翻转动画
  - [x] SubTask 9.3: 输入验证失败时的抖动反馈

- [x] Task 10: SignalFeed 组件动画
  - [x] SubTask 10.1: 新信号滑入动画
  - [x] SubTask 10.2: 高优先级信号脉冲警示
  - [x] SubTask 10.3: 信号列表 stagger 动画

- [x] Task 11: StrategyMarketplace 组件动画
  - [x] SubTask 11.1: 策略卡片悬停效果
  - [x] SubTask 11.2: 购买成功的庆祝动画
  - [x] SubTask 11.3: 列表加载 stagger 效果

- [x] Task 12: PortfolioReport 组件动画
  - [x] SubTask 12.1: 数据卡片数字滚动动画
  - [x] SubTask 12.2: 图表绘制动画
  - [x] SubTask 12.3: 盈亏数据的正负反馈动画

- [x] Task 13: AlertsManager 组件动画
  - [x] SubTask 13.1: 警报项进入/退出动画
  - [x] SubTask 13.2: 触发状态变化动画
  - [x] SubTask 13.3: 添加警报按钮涟漪效果

## 全局动画效果
- [x] Task 14: 页面过渡动画
  - [x] SubTask 14.1: 路由切换时的页面过渡
  - [x] SubTask 14.2: 组件加载时的渐显动画

- [x] Task 15: 骨架屏动画
  - [x] SubTask 15.1: 创建 Skeleton 组件
  - [x] SubTask 15.2: 添加流动光泽效果
  - [x] SubTask 15.3: 应用到各数据加载场景

## 性能优化与测试
- [x] Task 16: 动画性能优化
  - [x] SubTask 16.1: 使用 will-change 优化渲染
  - [x] SubTask 16.2: 实现 reduced-motion 媒体查询支持
  - [x] SubTask 16.3: 优化粒子效果在移动端的性能

- [x] Task 17: 代码质量检查
  - [x] SubTask 17.1: 运行 ESLint 检查
  - [x] SubTask 17.2: 运行 TypeScript 类型检查
  - [x] SubTask 17.3: 构建测试

# Task Dependencies
- Task 3 依赖 Task 1, Task 2
- Task 4 依赖 Task 1, Task 2
- Task 5-15 依赖 Task 3, Task 4
- Task 16 依赖 Task 5-15
- Task 17 依赖所有实现任务
