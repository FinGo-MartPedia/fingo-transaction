package repository

import (
	"context"

	"github.com/fingo-martpedia/fingo-transaction/constants"
	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(trx).Error
}

func (r *TransactionRepository) GetTransactionByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error) {
	var (
		resp models.Transaction
	)
	sql := r.DB.Where("reference=?", reference)
	if !includeRefund {
		sql = sql.Where("transaction_type != ?", constants.TransactionTypeRefund)
	}
	err := sql.Last(&resp).Error
	return resp, err
}

func (r *TransactionRepository) UpdateStatusTransaction(ctx context.Context, reference string, status string, additional_info string) error {
	return r.DB.Exec("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?", status, additional_info, reference).Error
}

func (r *TransactionRepository) GetTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	var (
		resp []models.Transaction
	)
	err := r.DB.Where("user_id = ?", userID).Find(&resp).Order("id DESC").Error
	return resp, err
}
