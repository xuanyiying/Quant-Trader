package repo

import (
	"context"
	"errors"

	"quant-trader/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaperAccount struct {
	db *gorm.DB
}

func NewPaperAccount(db *gorm.DB) *PaperAccount {
	return &PaperAccount{db: db}
}

func (r *PaperAccount) GetByUserID(ctx context.Context, userID int64) (*model.PaperAccount, error) {
	var account model.PaperAccount
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *PaperAccount) GetOrCreate(ctx context.Context, userID int64, initialBalance decimal.Decimal) (*model.PaperAccount, error) {
	account, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return account, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	account = &model.PaperAccount{
		UserID:         userID,
		Balance:        initialBalance,
		InitialBalance: initialBalance,
	}
	err = r.db.WithContext(ctx).Create(account).Error
	if err != nil {
		return nil, err
	}

	return r.GetByUserID(ctx, userID)
}

func (r *PaperAccount) UpdateBalance(ctx context.Context, userID int64, balance decimal.Decimal) error {
	return r.db.WithContext(ctx).Model(&model.PaperAccount{}).
		Where("user_id = ?", userID).
		Update("balance", balance).Error
}

func (r *PaperAccount) Reset(ctx context.Context, userID int64, balance decimal.Decimal) error {
	return r.db.WithContext(ctx).Model(&model.PaperAccount{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"balance":         balance,
			"initial_balance": balance,
		}).Error
}

// WithTx 返回一个使用指定事务的 PaperAccount 实例
func (r *PaperAccount) WithTx(tx *gorm.DB) *PaperAccount {
	return &PaperAccount{db: tx}
}

// TransferBalance 在事务中转移余额（原子操作）
// 使用悲观锁确保并发安全
func (r *PaperAccount) TransferBalance(ctx context.Context, buyerID, authorID int64, amount decimal.Decimal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 获取买家账户并加锁（FOR UPDATE）
		var buyerAccount model.PaperAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", buyerID).
			First(&buyerAccount).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("buyer account not found")
			}
			return err
		}

		// 检查余额是否足够
		if buyerAccount.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// 获取作者账户并加锁（FOR UPDATE）
		var authorAccount model.PaperAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", authorID).
			First(&authorAccount).Error; err != nil {
			// 作者账户不存在，创建一个
			if errors.Is(err, gorm.ErrRecordNotFound) {
				authorAccount = model.PaperAccount{
					UserID:         authorID,
					Balance:        decimal.Zero,
					InitialBalance: decimal.Zero,
				}
				if err := tx.Create(&authorAccount).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// 扣除买家余额
		if err := tx.Model(&model.PaperAccount{}).
			Where("user_id = ?", buyerID).
			Update("balance", buyerAccount.Balance.Sub(amount)).Error; err != nil {
			return err
		}

		// 增加作者余额
		if err := tx.Model(&model.PaperAccount{}).
			Where("user_id = ?", authorID).
			Update("balance", authorAccount.Balance.Add(amount)).Error; err != nil {
			return err
		}

		return nil
	})
}
