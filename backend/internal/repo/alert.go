package repo

import (
	"context"

	"quant-trader/internal/model"

	"gorm.io/gorm"
)

type Alert struct {
	db *gorm.DB
}

func NewAlert(db *gorm.DB) *Alert {
	return &Alert{db: db}
}

func (r *Alert) Create(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *Alert) GetByUserID(ctx context.Context, userID int64) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

func (r *Alert) GetByID(ctx context.Context, id int64) (*model.Alert, error) {
	var alert model.Alert
	err := r.db.WithContext(ctx).First(&alert, id).Error
	return &alert, err
}

func (r *Alert) GetActive(ctx context.Context) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).Where("is_triggered = ?", false).Find(&alerts).Error
	return alerts, err
}

func (r *Alert) UpdateTriggered(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_triggered": true,
	}).Error
}

func (r *Alert) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Alert{}, id).Error
}
