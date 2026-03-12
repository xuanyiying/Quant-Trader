# 量化交易平台重构计划

## 项目现状分析

### 已完成的 Repository 层
- `user.go` - 用户管理
- `api_key.go` - API密钥管理
- `alert.go` - 告警管理
- `kline.go` - K线数据管理 ✅
- `market.go` - 市场数据查询 ✅
- `marketplace.go` - 策略市场
- `paper_account.go` - 模拟交易账户
- `paper_order.go` - 模拟交易订单
- `paper_position.go` - 模拟交易持仓
- `portfolio.go` - 投资组合管理
- `subscription.go` - 订阅管理

### Handler 现状
- `paper.go` - 已使用 biz.PaperTrading ✅
- `portfolio.go` - 已使用 biz.Portfolio ✅
- `kline.go` - 已使用 biz.Kline ✅，但 `TriggerBackfill` 直接操作数据库 ⚠️

### 需要改进的地方
1. **kline.go TriggerBackfill** - 直接创建 `storage.HistoryBackfiller`，需要改为通过 Repository
2. **缺少单元测试** - 核心逻辑缺乏测试覆盖
3. **配置管理** - 已使用 viper，但可能需要完善

---

## 执行计划

### 阶段一：完善基础设施

#### 任务 1: 重构 kline.go TriggerBackfill（P0）
**问题**: `TriggerBackfill` 直接创建 `storage.HistoryBackfiller` 并操作数据库，违反了分层架构原则。

**解决方案**:
1. 创建 `HistoryBackfiller` Repository 或将其逻辑移至 Biz 层
2. 修改 Handler 通过 Biz 层调用

**具体步骤**:
1. 创建 `internal/biz/backfill.go` - 业务逻辑层
2. 修改 `internal/storage/backfiller.go` - 或将其整合到 Repository
3. 更新 `api/kline.go` - 使用 Biz 层调用

#### 任务 2: 补齐缺失的 Repository（P0）
经检查，以下 Repository 已存在且完整：
- ✅ `kline.go` - 已存在
- ✅ `market.go` - 已存在
- ✅ `strategy.go` - 策略配置通过 marketplace.go 管理

**结论**: Repository 层已基本完整，无需新增。

#### 任务 3: 添加单元测试（P1）
**目标文件**:
1. `repository/kline_test.go` - K线数据操作测试
2. `strategy/ma_test.go` - 策略测试（已存在，检查覆盖率）
3. `engine/backtester_test.go` - 回测引擎测试（已存在，检查覆盖率）

**新增测试**:
1. `repository/paper_account_test.go` - 模拟账户测试
2. `repository/portfolio_test.go` - 投资组合测试
3. `biz/paper_test.go` - 模拟交易业务逻辑测试
4. `biz/kline_test.go` - K线业务逻辑测试

#### 任务 4: 环境变量配置完善（P1）
**现状**: 已使用 viper 管理配置
**改进**:
1. 添加更多配置项的默认值
2. 添加配置验证逻辑
3. 考虑添加配置文件模板

---

### 阶段二：增强功能

#### 任务 5: 回测功能完善（P1）
**目标**:
1. 增加更多回测指标（夏普比率、最大回撤等）
2. 支持多策略对比
3. 优化回测性能

#### 任务 6: 策略市场（P1）
**现状**: `marketplace.go` Repository 已存在
**待完成**:
1. 完善策略发布流程
2. 添加策略购买逻辑
3. 策略版本管理

#### 任务 7: 实时监控面板（P2）
**目标**: Grafana 集成
**步骤**:
1. 添加 Prometheus 指标暴露
2. 创建 Grafana Dashboard 配置

#### 任务 8: 告警通知（P2）
**目标**: 支持微信/钉钉/邮件通知
**步骤**:
1. 扩展 alert Repository
2. 添加通知渠道配置
3. 实现各渠道发送逻辑

---

### 阶段三：AI 集成（长期）

#### 任务 9: AI Agent 接口（P2）
**目标**: 支持 OpenAI/Claude 调用

#### 任务 10: 信号策略（P2）
**目标**: AI 生成交易信号

#### 任务 11: 自然语言策略（P3）
**目标**: "当 BTC 跌 5% 时买入" 这类自然语言策略解析

---

## 立即可执行的任务清单

### 本周任务（Week 1）

- [ ] **Task 1.1**: 创建 `internal/biz/backfill.go` 业务逻辑层
- [ ] **Task 1.2**: 修改 `api/kline.go` 使用 Biz 层调用 backfill
- [ ] **Task 1.3**: 添加 `repository/kline_test.go` 单元测试
- [ ] **Task 1.4**: 添加 `biz/paper_test.go` 单元测试

### 下周任务（Week 2）

- [ ] **Task 2.1**: 添加 `repository/paper_account_test.go` 单元测试
- [ ] **Task 2.2**: 添加 `repository/portfolio_test.go` 单元测试
- [ ] **Task 2.3**: 完善配置验证逻辑
- [ ] **Task 2.4**: 跑通模拟交易 Demo

---

## 技术细节

### Task 1.1: Backfill Biz 层设计

```go
// internal/biz/backfill.go
package biz

type Backfill struct {
    backfiller *storage.HistoryBackfiller
    logger     *zap.Logger
}

func NewBackfill(db *pgxpool.Pool, logger *zap.Logger) *Backfill {
    return &Backfill{
        backfiller: storage.NewHistoryBackfiller(db, logger),
        logger:     logger,
    }
}

func (b *Backfill) BackfillBinance(ctx context.Context, symbol string, startTime, endTime time.Time) error {
    return b.backfiller.BackfillBinance(ctx, symbol, startTime, endTime)
}
```

### Task 1.3: Kline Repository 测试设计

```go
// repository/kline_test.go
func TestKline_GetLatest(t *testing.T) {
    // 使用 testcontainers 或内存数据库
    // 测试获取最新 K 线数据
}

func TestKline_BatchCreate(t *testing.T) {
    // 测试批量创建 K 线数据
}
```

---

## 风险评估

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 重构引入 Bug | 中 | 添加充分测试，小步提交 |
| 测试环境搭建复杂 | 低 | 使用 testcontainers |
| 业务逻辑理解偏差 | 低 | 与团队确认需求 |

---

## 成功标准

1. ✅ 所有 Handler 通过 Biz 层访问数据
2. ✅ 核心 Repository 有单元测试覆盖
3. ✅ 模拟交易 Demo 可正常运行
4. ✅ 代码通过 lint 和 typecheck
