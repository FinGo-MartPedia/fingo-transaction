package services

import (
	"context"

	"github.com/fingo-martpedia/fingo-transaction/constants"
	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/fingo-martpedia/fingo-transaction/internal/interfaces"
	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepository interfaces.ITransactionRepository
}

func NewTransactionService(transactionRepository interfaces.ITransactionRepository) *TransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, req *models.Transaction) (responses.CreateTransactionResponse, error) {
	var resp responses.CreateTransactionResponse

	req.TransactionStatus = constants.TransactionStatusPending
	req.Reference = helpers.GenerateReference()

	// jsonAdditionalInfo := map[string]interface{}{}
	// if req.AdditionalInfo != "" {
	// 	err := json.Unmarshal([]byte(req.AdditionalInfo), &jsonAdditionalInfo)
	// 	if err != nil {
	// 		return resp, errors.Wrap(err, "additional info type is invalid")
	// 	}
	// }

	err := s.TransactionRepository.CreateTransaction(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert create transaction")
	}

	resp.Reference = req.Reference
	resp.TransactionStatus = req.TransactionStatus
	return resp, nil

}
