package repository

import (
	"database/sql"
	"errors"
	"go-finance-tracker/internal/models"
	"gorm.io/gorm"
)

var (
	ErrFinanceHistoryNotFound = errors.New("finance history not found")
)

type UserFinanceRepository struct {
	db *gorm.DB
}

func NewUserFinanceRepository(db *gorm.DB) *UserFinanceRepository {
	return &UserFinanceRepository{db: db}
}

func (r *UserFinanceRepository) GetAll(userID int) (*[]models.FinanceRecord, error) {
	var records []models.FinanceRecord
	if err := r.db.Where("user_id = ?", userID).Preload("TransactionType").Preload("Category").Find(&records).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFinanceHistoryNotFound
		}
		return nil, err
	}
	return &records, nil
}

func (r *UserFinanceRepository) Create(record *models.FinanceRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return err
	}
	return nil
}
