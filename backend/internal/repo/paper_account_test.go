package repo

import (
	"context"
	"testing"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PaperAccountRepoTestSuite struct {
	suite.Suite
	db     *gorm.DB
	repo   *PaperAccount
	ctx    context.Context
	logger *zap.Logger
}

func (s *PaperAccountRepoTestSuite) SetupTest() {
	var err error
	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	s.Require().NoError(err)

	err = s.db.AutoMigrate(&model.PaperAccount{})
	s.Require().NoError(err)

	s.logger = zap.NewNop()
	s.repo = NewPaperAccount(s.db, s.logger)
	s.ctx = context.Background()
}

func (s *PaperAccountRepoTestSuite) TearDownTest() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

func (s *PaperAccountRepoTestSuite) TestGetByUserID_NotFound() {
	account, err := s.repo.GetByUserID(s.ctx, 999)
	s.ErrorIs(err, ErrNotFound)
	s.Nil(account)
}

func (s *PaperAccountRepoTestSuite) TestGetByUserID_Exists() {
	// 先创建账户
	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)
	_, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)

	// 查询账户
	account, err := s.repo.GetByUserID(s.ctx, userID)
	s.NoError(err)
	s.NotNil(account)
	s.Equal(userID, account.UserID)
	s.Equal(initialBalance, account.Balance)
}

func (s *PaperAccountRepoTestSuite) TestGetOrCreate_NewAccount() {
	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)

	account, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)
	s.NotNil(account)
	s.Equal(userID, account.UserID)
	s.Equal(initialBalance, account.Balance)
	s.Equal(initialBalance, account.InitialBalance)
	s.NotZero(account.ID)
}

func (s *PaperAccountRepoTestSuite) TestGetOrCreate_ExistingAccount() {
	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)

	// 第一次创建
	account1, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)

	// 第二次获取应该返回同一个账户
	account2, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)
	s.Equal(account1.ID, account2.ID)
}

func (s *PaperAccountRepoTestSuite) TestUpdateBalance() {
	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)

	// 创建账户
	_, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)

	// 更新余额
	newBalance := decimal.NewFromFloat(95000)
	err = s.repo.UpdateBalance(s.ctx, userID, newBalance)
	s.NoError(err)

	// 验证更新
	account, err := s.repo.GetByUserID(s.ctx, userID)
	s.NoError(err)
	s.Equal(newBalance, account.Balance)
}

func (s *PaperAccountRepoTestSuite) TestReset() {
	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)

	// 创建账户并修改余额
	_, err := s.repo.GetOrCreate(s.ctx, userID, initialBalance)
	s.NoError(err)

	err = s.repo.UpdateBalance(s.ctx, userID, decimal.NewFromFloat(50000))
	s.NoError(err)

	// 重置账户
	resetBalance := decimal.NewFromFloat(100000)
	err = s.repo.Reset(s.ctx, userID, resetBalance)
	s.NoError(err)

	// 验证重置
	account, err := s.repo.GetByUserID(s.ctx, userID)
	s.NoError(err)
	s.Equal(resetBalance, account.Balance)
	s.Equal(resetBalance, account.InitialBalance)
}

func TestPaperAccountRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PaperAccountRepoTestSuite))
}

// 独立的边界测试

func TestPaperAccount_UpdateBalance_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.PaperAccount{})
	logger := zap.NewNop()
	repo := NewPaperAccount(db, logger)
	ctx := context.Background()

	// 更新不存在的用户余额应该不返回错误（GORM 的 Update 行为）
	err := repo.UpdateBalance(ctx, 999, decimal.NewFromFloat(50000))
	assert.NoError(t, err)
}

func TestPaperAccount_Reset_NotFound(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.PaperAccount{})
	logger := zap.NewNop()
	repo := NewPaperAccount(db, logger)
	ctx := context.Background()

	// 重置不存在的用户应该不返回错误
	err := repo.Reset(ctx, 999, decimal.NewFromFloat(100000))
	assert.NoError(t, err)
}

// 并发测试

func TestPaperAccount_GetOrCreate_Concurrent(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.PaperAccount{})
	logger := zap.NewNop()
	repo := NewPaperAccount(db, logger)
	ctx := context.Background()

	userID := int64(1)
	initialBalance := decimal.NewFromFloat(100000)

	// 并发创建同一个用户的账户
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := repo.GetOrCreate(ctx, userID, initialBalance)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证只创建了一个账户
	var count int64
	db.Model(&model.PaperAccount{}).Where("user_id = ?", userID).Count(&count)
	assert.Equal(t, int64(1), count)
}
