package repo

import (
	"context"
	"testing"
	"time"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type KlineRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *Kline
	ctx  context.Context
}

func (s *KlineRepoTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)

	// 创建表
	err = s.db.AutoMigrate(&model.KLine{})
	s.Require().NoError(err)

	s.repo = NewKline(s.db)
	s.ctx = context.Background()
}

func (s *KlineRepoTestSuite) TearDownTest() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

func (s *KlineRepoTestSuite) TestCreate() {
	kline := &model.KLine{
		Symbol:    "BTCUSDT",
		Exchange:  model.ExchangeBinance,
		Period:    model.Period1m,
		Open:      decimal.NewFromFloat(50000.0),
		High:      decimal.NewFromFloat(51000.0),
		Low:       decimal.NewFromFloat(49000.0),
		Close:     decimal.NewFromFloat(50500.0),
		Volume:    decimal.NewFromFloat(100.5),
		Timestamp: time.Now(),
	}

	err := s.repo.Create(s.ctx, kline)
	s.NoError(err)
	s.NotZero(kline.ID)
}

func (s *KlineRepoTestSuite) TestBatchCreate() {
	now := time.Now()
	klines := []model.KLine{
		{
			Symbol:    "BTCUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(50000.0),
			High:      decimal.NewFromFloat(51000.0),
			Low:       decimal.NewFromFloat(49000.0),
			Close:     decimal.NewFromFloat(50500.0),
			Volume:    decimal.NewFromFloat(100.5),
			Timestamp: now,
		},
		{
			Symbol:    "BTCUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(50500.0),
			High:      decimal.NewFromFloat(51500.0),
			Low:       decimal.NewFromFloat(49500.0),
			Close:     decimal.NewFromFloat(51000.0),
			Volume:    decimal.NewFromFloat(150.0),
			Timestamp: now.Add(time.Minute),
		},
	}

	err := s.repo.BatchCreate(s.ctx, klines)
	s.NoError(err)

	// 验证数据已插入
	var count int64
	s.db.Model(&model.KLine{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *KlineRepoTestSuite) TestGetLatest() {
	now := time.Now()

	// 插入测试数据
	klines := []model.KLine{
		{
			Symbol:    "BTCUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(50000.0),
			High:      decimal.NewFromFloat(51000.0),
			Low:       decimal.NewFromFloat(49000.0),
			Close:     decimal.NewFromFloat(50500.0),
			Volume:    decimal.NewFromFloat(100.5),
			Timestamp: now.Add(-2 * time.Minute),
		},
		{
			Symbol:    "BTCUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(50500.0),
			High:      decimal.NewFromFloat(51500.0),
			Low:       decimal.NewFromFloat(49500.0),
			Close:     decimal.NewFromFloat(51000.0),
			Volume:    decimal.NewFromFloat(150.0),
			Timestamp: now.Add(-time.Minute),
		},
		{
			Symbol:    "BTCUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(51000.0),
			High:      decimal.NewFromFloat(52000.0),
			Low:       decimal.NewFromFloat(50000.0),
			Close:     decimal.NewFromFloat(51500.0),
			Volume:    decimal.NewFromFloat(200.0),
			Timestamp: now,
		},
	}
	err := s.repo.BatchCreate(s.ctx, klines)
	s.NoError(err)

	// 测试获取最新数据
	result, err := s.repo.GetLatest(s.ctx, "BTCUSDT", model.Period1m, 2)
	s.NoError(err)
	s.Len(result, 2)
	s.True(result[0].Timestamp.After(result[1].Timestamp))
}

func (s *KlineRepoTestSuite) TestGetBySymbol() {
	now := time.Now()
	startTime := now.Add(-1 * time.Hour)
	endTime := now

	// 插入测试数据
	klines := []model.KLine{
		{
			Symbol:    "ETHUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(3000.0),
			High:      decimal.NewFromFloat(3100.0),
			Low:       decimal.NewFromFloat(2900.0),
			Close:     decimal.NewFromFloat(3050.0),
			Volume:    decimal.NewFromFloat(500.0),
			Timestamp: startTime.Add(10 * time.Minute),
		},
		{
			Symbol:    "ETHUSDT",
			Exchange:  model.ExchangeBinance,
			Period:    model.Period1m,
			Open:      decimal.NewFromFloat(3050.0),
			High:      decimal.NewFromFloat(3150.0),
			Low:       decimal.NewFromFloat(2950.0),
			Close:     decimal.NewFromFloat(3100.0),
			Volume:    decimal.NewFromFloat(600.0),
			Timestamp: startTime.Add(20 * time.Minute),
		},
	}
	err := s.repo.BatchCreate(s.ctx, klines)
	s.NoError(err)

	// 测试按时间范围查询
	result, err := s.repo.GetBySymbol(s.ctx, "ETHUSDT", model.Period1m, startTime, endTime, 0)
	s.NoError(err)
	s.Len(result, 2)
}

func TestKlineRepoTestSuite(t *testing.T) {
	suite.Run(t, new(KlineRepoTestSuite))
}

// 独立的边界测试

func TestKline_GetLatest_EmptyResult(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.KLine{})
	repo := NewKline(db)
	ctx := context.Background()

	result, err := repo.GetLatest(ctx, "NONEXISTENT", model.Period1m, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestKline_GetBySymbol_NoDataInRange(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.KLine{})
	repo := NewKline(db)
	ctx := context.Background()

	now := time.Now()
	result, err := repo.GetBySymbol(ctx, "BTCUSDT", model.Period1m, now.Add(-2*time.Hour), now.Add(-1*time.Hour), 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
}
