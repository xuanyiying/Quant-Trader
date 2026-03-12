package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Agent AI Agent 接口
type Agent interface {
	GenerateSignal(ctx context.Context, req SignalRequest) (*SignalResponse, error)
	AnalyzeStrategy(ctx context.Context, req StrategyAnalysisRequest) (*StrategyAnalysisResponse, error)
	NaturalLanguageToStrategy(ctx context.Context, req NLStrategyRequest) (*NLStrategyResponse, error)
}

// SignalRequest 信号生成请求
type SignalRequest struct {
	Symbol          string                 `json:"symbol"`
	Timeframe       string                 `json:"timeframe"`
	MarketData      map[string]interface{} `json:"market_data"`
	Indicators      map[string]float64     `json:"indicators"`
	StrategyContext string                 `json:"strategy_context,omitempty"`
}

// SignalResponse 信号生成响应
type SignalResponse struct {
	Action     string  `json:"action"`     // buy, sell, hold
	Confidence float64 `json:"confidence"` // 0-1
	Reasoning  string  `json:"reasoning"`
	EntryPrice float64 `json:"entry_price,omitempty"`
	StopLoss   float64 `json:"stop_loss,omitempty"`
	TakeProfit float64 `json:"take_profit,omitempty"`
	Timestamp  int64   `json:"timestamp"`
}

// StrategyAnalysisRequest 策略分析请求
type StrategyAnalysisRequest struct {
	StrategyCode string                 `json:"strategy_code"`
	BacktestData map[string]interface{} `json:"backtest_data"`
}

// StrategyAnalysisResponse 策略分析响应
type StrategyAnalysisResponse struct {
	Analysis    string   `json:"analysis"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Suggestions []string `json:"suggestions"`
	RiskLevel   string   `json:"risk_level"`
}

// NLStrategyRequest 自然语言策略请求
type NLStrategyRequest struct {
	Description string `json:"description"`
	Language    string `json:"language"` // zh, en
}

// NLStrategyResponse 自然语言策略响应
type NLStrategyResponse struct {
	StrategyName string                 `json:"strategy_name"`
	StrategyType string                 `json:"strategy_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	Conditions   []Condition            `json:"conditions"`
	Explanation  string                 `json:"explanation"`
}

// Condition 策略条件
type Condition struct {
	Type      string  `json:"type"` // entry, exit
	Indicator string  `json:"indicator"`
	Operator  string  `json:"operator"` // >, <, =, cross_above, cross_below
	Value     float64 `json:"value"`
	Logic     string  `json:"logic"` // AND, OR
}

// OpenAIAgent OpenAI 实现
type OpenAIAgent struct {
	apiKey string
	model  string
	logger *zap.Logger
}

// NewOpenAIAgent 创建 OpenAI Agent
func NewOpenAIAgent(apiKey, model string, logger *zap.Logger) *OpenAIAgent {
	if model == "" {
		model = "gpt-4"
	}
	return &OpenAIAgent{
		apiKey: apiKey,
		model:  model,
		logger: logger,
	}
}

// GenerateSignal 生成交易信号
func (a *OpenAIAgent) GenerateSignal(ctx context.Context, req SignalRequest) (*SignalResponse, error) {
	_ = fmt.Sprintf(`
分析以下市场数据并给出交易建议:

交易对: %s
时间周期: %s
市场数据: %v
技术指标: %v

请提供:
1. 建议操作 (buy/sell/hold)
2. 置信度 (0-1)
3. 推理过程
4. 建议入场价、止损价、止盈价 (如适用)
`, req.Symbol, req.Timeframe, req.MarketData, req.Indicators)

	// TODO: 这里应该调用 OpenAI API
	// 简化实现，返回模拟数据
	a.logger.Info("generating signal with OpenAI", zap.String("symbol", req.Symbol))

	return &SignalResponse{
		Action:     "hold",
		Confidence: 0.7,
		Reasoning:  "基于当前市场分析，建议观望。",
		Timestamp:  0,
	}, nil
}

// AnalyzeStrategy 分析策略
func (a *OpenAIAgent) AnalyzeStrategy(ctx context.Context, req StrategyAnalysisRequest) (*StrategyAnalysisResponse, error) {
	a.logger.Info("analyzing strategy with OpenAI")

	return &StrategyAnalysisResponse{
		Analysis:    "策略整体表现良好，但需要注意风险控制。",
		Strengths:   []string{"盈利稳定", "回撤控制较好"},
		Weaknesses:  []string{"交易频率偏低", "对市场波动敏感"},
		Suggestions: []string{"考虑增加止损机制", "优化入场时机"},
		RiskLevel:   "medium",
	}, nil
}

// NaturalLanguageToStrategy 自然语言转策略
func (a *OpenAIAgent) NaturalLanguageToStrategy(ctx context.Context, req NLStrategyRequest) (*NLStrategyResponse, error) {
	a.logger.Info("converting natural language to strategy", zap.String("description", req.Description))

	// 解析自然语言描述
	// 这里简化实现，实际应该调用 LLM API

	response := &NLStrategyResponse{
		StrategyName: "自然语言策略",
		StrategyType: "conditional",
		Parameters: map[string]interface{}{
			"lookback_period": 20,
			"threshold":       0.05,
		},
		Conditions:  []Condition{},
		Explanation: fmt.Sprintf("已解析您的描述: %s", req.Description),
	}

	// 简单的关键词匹配
	if contains(req.Description, "跌") && contains(req.Description, "买") {
		response.StrategyName = "抄底策略"
		response.Conditions = append(response.Conditions, Condition{
			Type:      "entry",
			Indicator: "price_change",
			Operator:  "<",
			Value:     -0.05,
		})
	}

	if contains(req.Description, "涨") && contains(req.Description, "卖") {
		response.Conditions = append(response.Conditions, Condition{
			Type:      "exit",
			Indicator: "price_change",
			Operator:  ">",
			Value:     0.05,
		})
	}

	return response, nil
}

// ClaudeAgent Claude 实现
type ClaudeAgent struct {
	apiKey string
	model  string
	logger *zap.Logger
}

// NewClaudeAgent 创建 Claude Agent
func NewClaudeAgent(apiKey, model string, logger *zap.Logger) *ClaudeAgent {
	if model == "" {
		model = "claude-3-opus-20240229"
	}
	return &ClaudeAgent{
		apiKey: apiKey,
		model:  model,
		logger: logger,
	}
}

// GenerateSignal 生成交易信号
func (a *ClaudeAgent) GenerateSignal(ctx context.Context, req SignalRequest) (*SignalResponse, error) {
	a.logger.Info("generating signal with Claude", zap.String("symbol", req.Symbol))
	// 实现与 OpenAI 类似
	return &SignalResponse{
		Action:     "hold",
		Confidence: 0.75,
		Reasoning:  "基于Claude分析，建议观望。",
		Timestamp:  0,
	}, nil
}

// AnalyzeStrategy 分析策略
func (a *ClaudeAgent) AnalyzeStrategy(ctx context.Context, req StrategyAnalysisRequest) (*StrategyAnalysisResponse, error) {
	a.logger.Info("analyzing strategy with Claude")
	return &StrategyAnalysisResponse{
		Analysis:    "Claude分析：策略具有较好的风险收益比。",
		Strengths:   []string{"逻辑清晰", "执行纪律性强"},
		Weaknesses:  []string{"参数优化空间", "适应性有限"},
		Suggestions: []string{"增加动态参数调整", "考虑多周期验证"},
		RiskLevel:   "low",
	}, nil
}

// NaturalLanguageToStrategy 自然语言转策略
func (a *ClaudeAgent) NaturalLanguageToStrategy(ctx context.Context, req NLStrategyRequest) (*NLStrategyResponse, error) {
	a.logger.Info("converting natural language to strategy with Claude", zap.String("description", req.Description))

	return &NLStrategyResponse{
		StrategyName: "Claude解析策略",
		StrategyType: "ai_generated",
		Parameters: map[string]interface{}{
			"confidence_threshold": 0.8,
		},
		Conditions: []Condition{
			{
				Type:      "entry",
				Indicator: "ai_signal",
				Operator:  ">",
				Value:     0.7,
			},
		},
		Explanation: fmt.Sprintf("Claude理解您的需求: %s", req.Description),
	}, nil
}

// SignalService 信号服务
type SignalService struct {
	agent  Agent
	logger *zap.Logger
}

// NewSignalService 创建信号服务
func NewSignalService(agent Agent, logger *zap.Logger) *SignalService {
	return &SignalService{
		agent:  agent,
		logger: logger,
	}
}

// GenerateTradingSignal 生成交易信号
func (s *SignalService) GenerateTradingSignal(ctx context.Context, symbol, timeframe string, marketData map[string]interface{}) (*SignalResponse, error) {
	req := SignalRequest{
		Symbol:     symbol,
		Timeframe:  timeframe,
		MarketData: marketData,
	}
	return s.agent.GenerateSignal(ctx, req)
}

// StrategyParserService 策略解析服务
type StrategyParserService struct {
	agent  Agent
	logger *zap.Logger
}

// NewStrategyParserService 创建策略解析服务
func NewStrategyParserService(agent Agent, logger *zap.Logger) *StrategyParserService {
	return &StrategyParserService{
		agent:  agent,
		logger: logger,
	}
}

// ParseNaturalLanguage 解析自然语言策略
func (s *StrategyParserService) ParseNaturalLanguage(ctx context.Context, description, language string) (*NLStrategyResponse, error) {
	req := NLStrategyRequest{
		Description: description,
		Language:    language,
	}
	return s.agent.NaturalLanguageToStrategy(ctx, req)
}

// 辅助函数
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ToJSON 转换为 JSON
func ToJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
