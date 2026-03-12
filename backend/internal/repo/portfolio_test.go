package repo

import (
	"context"
	"testing"
	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PortfolioRepoTestSuite struct {
	suite.Suite
	db        *gorm.DB
	repo      *Portfolio
	assetRepo *PortfolioAsset
	ctx       context.Context
}

func (s *PortfolioRepoTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)

	err = s.db.AutoMigrate(&model.Portfolio{}, &model.PortfolioAsset{})
	s.Require().NoError(err)

	s.repo = NewPortfolio(s.db)
	s.assetRepo = NewPortfolioAsset(s.db)
	s.ctx = context.Background()
}

func (s *PortfolioRepoTestSuite) TearDownTest() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

func (s *PortfolioRepoTestSuite) TestCreate() {
	portfolio := &model.Portfolio{
		UserID: 1,
		Name:   "Test Portfolio",
	}

	err := s.repo.Create(s.ctx, portfolio)
	s.NoError(err)
	s.NotZero(portfolio.ID)
	s.NotZero(portfolio.CreatedAt)
}

func (s *PortfolioRepoTestSuite) TestGetByID() {
	// 创建投资组合
	portfolio := &model.Portfolio{
		UserID: 1,
		Name:   "Test Portfolio",
	}
	err := s.repo.Create(s.ctx, portfolio)
	s.NoError(err)

	// 查询
	result, err := s.repo.GetByID(s.ctx, portfolio.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(portfolio.Name, result.Name)
	s.Equal(portfolio.UserID, result.UserID)
}

func (s *PortfolioRepoTestSuite) TestGetByUserID() {
	userID := int64(1)

	// 创建多个投资组合
	portfolios := []*model.Portfolio{
		{UserID: userID, Name: "Portfolio 1"},
		{UserID: userID, Name: "Portfolio 2"},
		{UserID: 2, Name: "Other User Portfolio"},
	}

	for _, p := range portfolios {
		err := s.repo.Create(s.ctx, p)
		s.NoError(err)
	}

	// 查询用户1的投资组合
	result, err := s.repo.GetByUserID(s.ctx, userID)
	s.NoError(err)
	s.Len(result, 2)
}

func (s *PortfolioRepoTestSuite) TestUpdate() {
	// 创建
	portfolio := &model.Portfolio{
		UserID: 1,
		Name:   "Original Name",
	}
	err := s.repo.Create(s.ctx, portfolio)
	s.NoError(err)

	// 更新
	portfolio.Name = "Updated Name"
	err = s.repo.Update(s.ctx, portfolio)
	s.NoError(err)

	// 验证
	result, err := s.repo.GetByID(s.ctx, portfolio.ID)
	s.NoError(err)
	s.Equal("Updated Name", result.Name)
}

func (s *PortfolioRepoTestSuite) TestDelete() {
	// 创建
	portfolio := &model.Portfolio{
		UserID: 1,
		Name:   "To Delete",
	}
	err := s.repo.Create(s.ctx, portfolio)
	s.NoError(err)

	// 删除
	err = s.repo.Delete(s.ctx, portfolio.ID)
	s.NoError(err)

	// 验证已删除（软删除）
	var count int64
	s.db.Model(&model.Portfolio{}).Where("id = ?", portfolio.ID).Count(&count)
	s.Equal(int64(0), count)
}

// PortfolioAsset 测试

func (s *PortfolioRepoTestSuite) TestPortfolioAsset_Upsert_Create() {
	asset := &model.PortfolioAsset{
		PortfolioID: 1,
		Symbol:      "BTCUSDT",
		Side:        "long",
		Qty:         decimal.NewFromFloat(1.5),
		EntryPrice:  decimal.NewFromFloat(50000),
	}

	err := s.assetRepo.Upsert(s.ctx, asset)
	s.NoError(err)
	s.NotZero(asset.ID)
}

func (s *PortfolioRepoTestSuite) TestPortfolioAsset_Upsert_Update() {
	// 创建资产
	asset := &model.PortfolioAsset{
		PortfolioID: 1,
		Symbol:      "BTCUSDT",
		Side:        "long",
		Qty:         decimal.NewFromFloat(1.5),
		EntryPrice:  decimal.NewFromFloat(50000),
	}
	err := s.assetRepo.Upsert(s.ctx, asset)
	s.NoError(err)

	// 更新资产
	asset.Qty = decimal.NewFromFloat(2.0)
	asset.EntryPrice = decimal.NewFromFloat(51000)
	err = s.assetRepo.Upsert(s.ctx, asset)
	s.NoError(err)

	// 验证
	assets, err := s.assetRepo.GetByPortfolioID(s.ctx, 1)
	s.NoError(err)
	s.Len(assets, 1)
	s.Equal(decimal.NewFromFloat(2.0), assets[0].Qty)
}

func (s *PortfolioRepoTestSuite) TestPortfolioAsset_GetByPortfolioID() {
	// 创建多个资产
	assets := []*model.PortfolioAsset{
		{PortfolioID: 1, Symbol: "BTCUSDT", Side: "long", Qty: decimal.NewFromFloat(1.0), EntryPrice: decimal.NewFromFloat(50000)},
		{PortfolioID: 1, Symbol: "ETHUSDT", Side: "long", Qty: decimal.NewFromFloat(10.0), EntryPrice: decimal.NewFromFloat(3000)},
		{PortfolioID: 2, Symbol: "BTCUSDT", Side: "short", Qty: decimal.NewFromFloat(0.5), EntryPrice: decimal.NewFromFloat(51000)},
	}

	for _, a := range assets {
		err := s.assetRepo.Upsert(s.ctx, a)
		s.NoError(err)
	}

	// 查询 Portfolio 1 的资产
	result, err := s.assetRepo.GetByPortfolioID(s.ctx, 1)
	s.NoError(err)
	s.Len(result, 2)
}

func (s *PortfolioRepoTestSuite) TestPortfolioAsset_Delete() {
	// 创建资产
	asset := &model.PortfolioAsset{
		PortfolioID: 1,
		Symbol:      "BTCUSDT",
		Side:        "long",
		Qty:         decimal.NewFromFloat(1.5),
		EntryPrice:  decimal.NewFromFloat(50000),
	}
	err := s.assetRepo.Upsert(s.ctx, asset)
	s.NoError(err)

	// 删除
	err = s.assetRepo.Delete(s.ctx, asset.ID)
	s.NoError(err)

	// 验证
	assets, err := s.assetRepo.GetByPortfolioID(s.ctx, 1)
	s.NoError(err)
	s.Empty(assets)
}

func TestPortfolioRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PortfolioRepoTestSuite))
}

// 边界测试

func TestPortfolio_GetByID_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.Portfolio{})
	repo := NewPortfolio(db)
	ctx := context.Background()

	result, err := repo.GetByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPortfolio_GetByUserID_Empty(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.Portfolio{})
	repo := NewPortfolio(db)
	ctx := context.Background()

	result, err := repo.GetByUserID(ctx, 999)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestPortfolioAsset_GetByPortfolioID_Empty(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.PortfolioAsset{})
	repo := NewPortfolioAsset(db)
	ctx := context.Background()

	result, err := repo.GetByPortfolioID(ctx, 999)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
