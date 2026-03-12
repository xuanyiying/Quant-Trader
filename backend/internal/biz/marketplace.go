package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrStrategyNotFound    = errors.New("strategy not found")
	ErrAlreadyPurchased    = errors.New("strategy already purchased")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInvalidStrategyData = errors.New("invalid strategy data")
)

// Marketplace 策略市场业务逻辑
type Marketplace struct {
	repo         *repo.Marketplace
	paperAccount *repo.PaperAccount
	logger       *zap.Logger
}

// NewMarketplace 创建策略市场业务逻辑
func NewMarketplace(repo *repo.Marketplace, paperAccount *repo.PaperAccount, logger *zap.Logger) *Marketplace {
	return &Marketplace{
		repo:         repo,
		paperAccount: paperAccount,
		logger:       logger,
	}
}

// ListPublic 列出公开策略
func (b *Marketplace) ListPublic(ctx context.Context) ([]repo.MarketItem, error) {
	return b.repo.ListPublic(ctx)
}

// Purchase 购买策略（基础版本，不含余额检查）
func (b *Marketplace) Purchase(ctx context.Context, userID, strategyID int64) error {
	price, err := b.repo.GetPrice(ctx, strategyID)
	if err != nil {
		return err
	}

	if price.GreaterThan(decimal.Zero) {
		b.logger.Info("purchasing paid strategy",
			zap.Int64("user_id", userID),
			zap.Int64("strategy_id", strategyID),
			zap.String("price", price.String()))
	}

	return b.repo.Purchase(ctx, userID, strategyID, price)
}

// PublishStrategyInput 发布策略输入参数
type PublishStrategyInput struct {
	AuthorID     int64
	Name         string
	Description  string
	Price        decimal.Decimal
	StrategyCode string
	Metrics      StrategyMetrics
}

// StrategyMetrics 策略指标
type StrategyMetrics struct {
	TotalReturn    float64 `json:"total_return"`
	SharpeRatio    float64 `json:"sharpe_ratio"`
	MaxDrawdown    float64 `json:"max_drawdown"`
	WinRate        float64 `json:"win_rate"`
	TotalTrades    int     `json:"total_trades"`
	BacktestPeriod string  `json:"backtest_period"`
}

// PublishStrategy 发布策略到市场
func (b *Marketplace) PublishStrategy(ctx context.Context, input PublishStrategyInput) (*model.StrategyMarket, error) {
	// 验证输入
	if input.Name == "" {
		return nil, ErrInvalidStrategyData
	}

	if input.Price.IsNegative() {
		return nil, ErrInvalidStrategyData
	}

	// 创建市场条目
	item := &model.StrategyMarket{
		UserID:      input.AuthorID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		IsPublished: false, // 需要审核，初始未发布
		CreatedAt:   time.Now(),
	}

	// 保存到数据库
	if err := b.repo.Create(ctx, item); err != nil {
		b.logger.Error("failed to create strategy market item", zap.Error(err))
		return nil, err
	}

	b.logger.Info("strategy published",
		zap.Int64("author_id", input.AuthorID),
		zap.String("name", input.Name),
		zap.String("price", input.Price.String()),
	)

	return item, nil
}

// ApproveStrategy 审核通过策略
func (b *Marketplace) ApproveStrategy(ctx context.Context, strategyID int64, reviewerID int64) error {
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return err
	}

	if item.IsPublished {
		return errors.New("strategy is already approved")
	}

	item.IsPublished = true
	item.UpdatedAt = time.Now()

	if err := b.repo.Update(ctx, item); err != nil {
		return err
	}

	b.logger.Info("strategy approved",
		zap.Int64("strategy_id", strategyID),
		zap.Int64("reviewer_id", reviewerID),
	)

	return nil
}

// RejectStrategy 拒绝策略
func (b *Marketplace) RejectStrategy(ctx context.Context, strategyID int64, reviewerID int64, reason string) error {
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return err
	}

	item.Description = fmt.Sprintf("%s\n\nRejection reason: %s", item.Description, reason)
	item.UpdatedAt = time.Now()

	if err := b.repo.Update(ctx, item); err != nil {
		return err
	}

	b.logger.Info("strategy rejected",
		zap.Int64("strategy_id", strategyID),
		zap.Int64("reviewer_id", reviewerID),
		zap.String("reason", reason),
	)

	return nil
}

// GetStrategyDetail 获取策略详情
func (b *Marketplace) GetStrategyDetail(ctx context.Context, strategyID int64, userID int64) (*StrategyDetail, error) {
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return nil, err
	}

	// 检查用户是否已购买
	hasPurchased, err := b.repo.HasPurchased(ctx, userID, strategyID)
	if err != nil {
		return nil, err
	}

	detail := &StrategyDetail{
		StrategyMarket: *item,
		HasPurchased:   hasPurchased,
		IsAuthor:       item.UserID == userID,
	}

	return detail, nil
}

// StrategyDetail 策略详情
type StrategyDetail struct {
	model.StrategyMarket
	HasPurchased bool `json:"has_purchased"`
	IsAuthor     bool `json:"is_author"`
}

// PurchaseStrategy 购买策略（含余额检查）
func (b *Marketplace) PurchaseStrategy(ctx context.Context, userID int64, strategyID int64) error {
	// 检查策略是否存在
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return ErrStrategyNotFound
	}

	if !item.IsPublished {
		return errors.New("strategy is not available for purchase")
	}

	// 检查是否已购买
	hasPurchased, err := b.repo.HasPurchased(ctx, userID, strategyID)
	if err != nil {
		return err
	}
	if hasPurchased {
		return ErrAlreadyPurchased
	}

	// 检查余额（如果是付费策略）
	if item.Price.GreaterThan(decimal.Zero) {
		if b.paperAccount == nil {
			return errors.New("paper account service not available")
		}

		// 使用事务转移余额（原子操作）
		if err := b.paperAccount.TransferBalance(ctx, userID, item.UserID, item.Price); err != nil {
			b.logger.Error("failed to transfer balance",
				zap.Error(err),
				zap.Int64("buyer_id", userID),
				zap.Int64("author_id", item.UserID),
				zap.String("amount", item.Price.String()),
			)
			return fmt.Errorf("failed to transfer balance: %w", err)
		}
	}

	// 创建购买记录
	purchase := &model.StrategyPurchase{
		UserID:     userID,
		StrategyID: strategyID,
		Price:      item.Price,
		CreatedAt:  time.Now(),
	}

	if err := b.repo.CreatePurchase(ctx, purchase); err != nil {
		// 购买记录创建失败，需要回滚余额
		if item.Price.GreaterThan(decimal.Zero) {
			b.logger.Error("failed to create purchase record, rolling back balances",
				zap.Error(err),
				zap.Int64("user_id", userID),
				zap.Int64("strategy_id", strategyID),
			)
			// 回滚买家余额
			account, _ := b.paperAccount.GetByUserID(ctx, userID)
			if account != nil {
				originalBalance := account.Balance.Add(item.Price)
				if rollbackErr := b.paperAccount.UpdateBalance(ctx, userID, originalBalance); rollbackErr != nil {
					b.logger.Error("failed to rollback buyer balance after purchase failure",
						zap.Error(rollbackErr),
						zap.Int64("user_id", userID),
					)
				}
			}
			// 回滚作者余额
			authorAccount, _ := b.paperAccount.GetByUserID(ctx, item.UserID)
			if authorAccount != nil {
				authorOriginalBalance := authorAccount.Balance.Sub(item.Price)
				if rollbackErr := b.paperAccount.UpdateBalance(ctx, item.UserID, authorOriginalBalance); rollbackErr != nil {
					b.logger.Error("failed to rollback author balance after purchase failure",
						zap.Error(rollbackErr),
						zap.Int64("author_id", item.UserID),
					)
				}
			}
		}
		return fmt.Errorf("failed to create purchase record: %w", err)
	}

	b.logger.Info("strategy purchased",
		zap.Int64("user_id", userID),
		zap.Int64("strategy_id", strategyID),
		zap.String("price", item.Price.String()),
	)

	return nil
}

// GetMyStrategies 获取我的策略（发布的和购买的）
func (b *Marketplace) GetMyStrategies(ctx context.Context, userID int64) (*MyStrategies, error) {
	// 获取发布的策略
	published, err := b.repo.GetByAuthorID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取购买的策略
	purchased, err := b.repo.GetPurchasedByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &MyStrategies{
		Published: published,
		Purchased: purchased,
	}, nil
}

// MyStrategies 我的策略
type MyStrategies struct {
	Published []model.StrategyMarket   `json:"published"`
	Purchased []model.StrategyPurchase `json:"purchased"`
}

// UpdateStrategy 更新策略信息（仅作者）
func (b *Marketplace) UpdateStrategy(ctx context.Context, strategyID int64, userID int64, updates map[string]interface{}) error {
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return err
	}

	if item.UserID != userID {
		return ErrUnauthorized
	}

	// 应用更新
	if name, ok := updates["name"].(string); ok {
		item.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		item.Description = desc
	}
	if price, ok := updates["price"].(decimal.Decimal); ok && !price.IsNegative() {
		item.Price = price
	}

	item.UpdatedAt = time.Now()

	return b.repo.Update(ctx, item)
}

// DeleteStrategy 删除策略（仅作者或管理员）
func (b *Marketplace) DeleteStrategy(ctx context.Context, strategyID int64, userID int64, isAdmin bool) error {
	item, err := b.repo.GetByID(ctx, strategyID)
	if err != nil {
		return err
	}

	if item.UserID != userID && !isAdmin {
		return ErrUnauthorized
	}

	return b.repo.Delete(ctx, strategyID)
}

// SearchStrategies 搜索策略
func (b *Marketplace) SearchStrategies(ctx context.Context, query string, minPrice, maxPrice decimal.Decimal, sortBy string) ([]model.StrategyMarket, error) {
	return b.repo.Search(ctx, query, minPrice, maxPrice, sortBy)
}

// GetTopStrategies 获取热门策略
func (b *Marketplace) GetTopStrategies(ctx context.Context, limit int) ([]model.StrategyMarket, error) {
	return b.repo.GetTopBySales(ctx, limit)
}
