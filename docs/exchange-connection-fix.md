# 交易所连接问题解决方案

## 问题描述

所有交易所 (Binance, Kraken, Bybit, Coinbase, OKX) 的 WebSocket 连接出现 `i/o timeout` 错误:

```
{"level":"error","msg":"failed to connect to binance","error":"dial tcp 45.114.11.25:9443: i/o timeout"}
{"level":"error","msg":"failed to connect to Kraken","error":"dial tcp 104.23.124.189:443: i/o timeout"}
```

## 根本原因

网络防火墙或地域限制导致无法直接访问加密货币交易所的 API 服务器。这种情况常见于:
- 中国大陆地区 (GFW 防火墙)
- 公司/学校网络限制
- ISP 服务商限制

## 解决方案

### 方案 1: 配置代理 (推荐)

已实现代理支持，现在所有连接器都支持以下三种代理配置方式:

#### 1.1 使用环境变量

在 `.env` 文件中添加:

```bash
# HTTP/HTTPS 代理
HTTPS_PROXY=http://127.0.0.1:7890

# 或 SOCKS5 代理
SOCKS5_PROXY=socks5://127.0.0.1:7891

# 或简写格式 (默认使用 SOCKS5)
SOCKS5_PROXY=127.0.0.1:7891
```

#### 1.2 在终端中设置

```bash
# 临时设置 (当前终端会话有效)
export HTTPS_PROXY=http://127.0.0.1:7890

# 或
export ALL_PROXY=socks5://127.0.0.1:7891

# 永久设置 (添加到 ~/.zshrc 或 ~/.bashrc)
echo 'export HTTPS_PROXY=http://127.0.0.1:7890' >> ~/.zshrc
source ~/.zshrc
```

#### 1.3 常见代理工具默认端口

| 代理工具 | 默认端口 | 协议 |
|---------|---------|------|
| Clash X | 7890 | HTTP/HTTPS |
| Clash X | 7891 | SOCKS5 |
| Shadowsocks | 1080 | SOCKS5 |
| V2Ray | 10808 | SOCKS5 |

### 方案 2: 禁用部分交易所 (备选)

如果某些交易所不需要，可以修改 `internal/app/symbols.go` 只保留可访问的交易所:

```go
func (a *App) ingestionTargets() []ingestionTarget {
	return []ingestionTarget{
		// 只保留可以访问的交易所
		{Exchange: "okx", Symbol: "BTC-USDT"},  // OKX 在中国大陆通常可访问
		// {Exchange: "binance", Symbol: "btcusdt"},  // 注释掉无法访问的
		// {Exchange: "bybit", Symbol: "BTCUSDT"},
		// {Exchange: "coinbase", Symbol: "BTC-USD"},
		// {Exchange: "kraken", Symbol: "XBT/USD"},
	}
}
```

### 方案 3: 使用云服务器

在不受限制的地区 (如新加坡、美国) 部署云服务器，通过远程服务器连接交易所。

## 验证连接

配置代理后，可以通过以下方式验证:

### 2.1 测试网络连通性

```bash
# 测试 Binance WebSocket
curl -v https://stream.binance.com:9443

# 或使用 ping (部分服务器禁 ping)
ping stream.binance.com
```

### 2.2 查看应用日志

启动应用后，观察日志:

```bash
cd backend
go run .

# 成功连接应该看到:
{"level":"info","msg":"connected to binance websocket"}
{"level":"info","msg":"connected to Kraken websocket"}
```

## 其他相关问题

### 数据库外键错误

```
ERROR: insert or update on table "paper_accounts" violates foreign key constraint "paper_accounts_user_id_fkey"
```

这是因为尝试为不存在的用户创建模拟账户。解决方法:

1. 确保用户已创建:
```sql
-- 检查用户是否存在
SELECT * FROM users WHERE id = 2;

-- 如果不存在，创建测试用户
INSERT INTO users (email, password_hash, name) 
VALUES ('test@example.com', 'hashed_password', 'Test User');
```

2. 或者先注册新用户，再访问模拟账户接口

### NATS 流不匹配错误

```
failed to subscribe to NATS","error":"nats: no stream matches subject"
```

确保 NATS JetStream 已正确初始化并创建了所需的 Stream:

```bash
# 检查 NATS Stream
nats stream ls

# 应该看到 MARKET stream
# 如果没有，重启应用会自动创建
```

## 技术实现细节

### 修改的文件

1. `internal/connector/proxy.go` - 新增，共享代理配置读取逻辑
2. `internal/connector/binance.go` - 添加代理支持
3. `internal/connector/kraken.go` - 添加代理支持
4. `internal/connector/bybit.go` - 添加代理支持
5. `internal/connector/coinbase.go` - 添加代理支持
6. `internal/connector/okx.go` - 添加代理支持
7. `backend/.env.example` - 新增，配置模板

### 实现原理

所有连接器现在都会:
1. 启动时读取环境变量 `HTTPS_PROXY`, `HTTP_PROXY`, `SOCKS5_PROXY`
2. 自动解析代理地址
3. 在 WebSocket Dialer 中配置代理
4. 通过代理服务器建立 WebSocket 连接

## 常见问题

### Q: 我没有代理怎么办？

A: 可以考虑:
1. 使用 OKX 交易所 (通常在中国大陆可访问)
2. 使用免费的代理测试 (不推荐用于生产)
3. 在本地运行测试时使用 mock 数据

### Q: 代理配置后仍然无法连接？

A: 检查:
1. 代理服务器是否正常运行
2. 代理地址和端口是否正确
3. 代理是否需要认证
4. 尝试更换代理协议 (HTTP vs SOCKS5)

### Q: 如何测试代理是否生效？

A: 使用 curl 测试:
```bash
# 不使用代理
curl -I https://stream.binance.com

# 使用代理
curl -I -x http://127.0.0.1:7890 https://stream.binance.com
```

## 下一步

配置代理后，重新启动应用:

```bash
cd backend
go run .
```

观察日志确认所有交易所连接成功。
