package repo

import (
	"context"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Marketplace struct {
	db *gorm.DB
}

func NewMarketplace(db *gorm.DB) *Marketplace {
	return &Marketplace{db: db}
}

type MarketItem struct {
	ID          int64
	Name        string
	Description string
	Price       decimal.Decimal
	AuthorEmail string
	Metrics     string
}

func (r *Marketplace) ListPublic(ctx context.Context) ([]MarketItem, error) {
	var items []MarketItem
	err := r.db.WithContext(ctx).
		Table("strategy_market m").
		Select("m.id, m.price, m.description, m.performance_metrics as metrics, s.name, u.email as author_email").
		Joins("JOIN strategies s ON m.strategy_id = s.id").
		Joins("JOIN users u ON m.owner_id = u.id").
		Where("m.is_public = ?", true).
		Scan(&items).Error
	return items, err
}

func (r *Marketplace) GetPrice(ctx context.Context, id int64) (decimal.Decimal, error) {
	var item model.StrategyMarket
	err := r.db.WithContext(ctx).
		Select("price").
		Where("id = ? AND is_published = ?", id, true).
		First(&item).Error
	return item.Price, err
}

func (r *Marketplace) Purchase(ctx context.Context, userID, strategyID int64, price decimal.Decimal) error {
	purchase := &model.StrategyPurchase{
		UserID:     userID,
		StrategyID: strategyID,
		Price:      price,
	}
	return r.db.WithContext(ctx).Create(purchase).Error
}

// Create 创建策略市场条目
func (r *Marketplace) Create(ctx context.Context, item *model.StrategyMarket) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetByID 根据ID获取策略
func (r *Marketplace) GetByID(ctx context.Context, id int64) (*model.StrategyMarket, error) {
	var item model.StrategyMarket
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新策略
func (r *Marketplace) Update(ctx context.Context, item *model.StrategyMarket) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete 删除策略
func (r *Marketplace) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.StrategyMarket{}, id).Error
}

// HasPurchased 检查用户是否已购买策略
func (r *Marketplace) HasPurchased(ctx context.Context, userID, strategyID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.StrategyPurchase{}).
		Where("user_id = ? AND strategy_id = ?", userID, strategyID).
		Count(&count).Error
	return count > 0, err
}

// GetByAuthorID 获取作者发布的所有策略
func (r *Marketplace) GetByAuthorID(ctx context.Context, userID int64) ([]model.StrategyMarket, error) {
	var items []model.StrategyMarket
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&items).Error
	return items, err
}

// GetPurchasedByUserID 获取用户购买的所有策略
func (r *Marketplace) GetPurchasedByUserID(ctx context.Context, userID int64) ([]model.StrategyPurchase, error) {
	var purchases []model.StrategyPurchase
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&purchases).Error
	return purchases, err
}

// CreatePurchase 创建购买记录
func (r *Marketplace) CreatePurchase(ctx context.Context, purchase *model.StrategyPurchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

// Search 搜索策略
func (r *Marketplace) Search(ctx context.Context, query string, minPrice, maxPrice decimal.Decimal, sortBy string) ([]model.StrategyMarket, error) {
	var items []model.StrategyMarket
	db := r.db.WithContext(ctx).
		Where("is_published = ?", true)

	if query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if !minPrice.IsZero() {
		db = db.Where("price >= ?", minPrice)
	}

	if !maxPrice.IsZero() {
		db = db.Where("price <= ?", maxPrice)
	}

	// 排序
	switch sortBy {
	case "price_asc":
		db = db.Order("price ASC")
	case "price_desc":
		db = db.Order("price DESC")
	case "newest":
		db = db.Order("created_at DESC")
	default:
		db = db.Order("created_at DESC")
	}

	err := db.Find(&items).Error
	return items, err
}

// GetTopBySales 获取热门策略（按销量排序）
func (r *Marketplace) GetTopBySales(ctx context.Context, limit int) ([]model.StrategyMarket, error) {
	var items []model.StrategyMarket
	err := r.db.WithContext(ctx).
		Table("strategy_market sm").
		Select("sm.*, COUNT(sp.id) as sales_count").
		Joins("LEFT JOIN strategy_purchases sp ON sm.id = sp.strategy_id").
		Where("sm.is_published = ?", true).
		Group("sm.id").
		Order("sales_count DESC").
		Limit(limit).
		Scan(&items).Error
	return items, err
}
