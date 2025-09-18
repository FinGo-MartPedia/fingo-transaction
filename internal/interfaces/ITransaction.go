package interfaces

import (
	"context"

	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/requests"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
	"github.com/gin-gonic/gin"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error)
	UpdateStatusTransaction(ctx context.Context, reference string, status string, additional_info string) error
	GetTransactions(ctx context.Context, userID int) ([]models.Transaction, error)
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (responses.CreateTransactionResponse, error)
	UpdateStatusTransaction(ctx context.Context, tokenData string, req *requests.UpdateStatusTransaction) error
	GetTransactions(ctx context.Context, userID int) ([]responses.TransactionResponse, error)
	GetTransactionDetail(ctx context.Context, reference string) (responses.TransactionResponse, error)
	RefundTransaction(ctx context.Context, token string, userId int, req *requests.RefundTransaction) (responses.CreateTransactionResponse, error)
}

type ITransactionController interface {
	CreateTransaction(c *gin.Context)
	UpdateStatusTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
	GetTransactionDetail(c *gin.Context)
	RefundTransaction(c *gin.Context)
}
