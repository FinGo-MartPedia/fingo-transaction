package interfaces

import (
	"context"

	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
	"github.com/gin-gonic/gin"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (responses.CreateTransactionResponse, error)
}

type ITransactionController interface {
	CreateTransaction(c *gin.Context)
}
