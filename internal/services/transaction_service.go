package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fingo-martpedia/fingo-transaction/constants"
	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/fingo-martpedia/fingo-transaction/internal/interfaces"
	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/requests"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepository interfaces.ITransactionRepository
	WalletExternal        interfaces.IWalletExternal
}

func NewTransactionService(transactionRepository interfaces.ITransactionRepository, walletExternal interfaces.IWalletExternal) *TransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
		WalletExternal:        walletExternal,
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

func (s *TransactionService) UpdateStatusTransaction(ctx context.Context, tokenData string, req *requests.UpdateStatusTransaction) error {
	// get transaction by reference
	trx, err := s.TransactionRepository.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction")
	}

	// validate transaction status flow
	statusValid := false
	mapStatusFlow := constants.MapTransactionStatusFlow[trx.TransactionStatus]
	for i := range mapStatusFlow {
		if mapStatusFlow[i] == req.TransactionStatus {
			statusValid = true
		}
	}
	if !statusValid {
		return fmt.Errorf("transaction status flow invalid. request status = %s", req.TransactionStatus)
	}

	// request update balance to ewallet-wallet
	reqUpdateBalance := requests.UpdateBalance{
		Amount:    trx.Amount,
		Reference: req.Reference,
	}
	if req.TransactionStatus == constants.TransactionStatusReversed {
		reqUpdateBalance.Reference = "REVERSED-" + req.Reference

		now := time.Now()
		expiredReversalTime := trx.CreatedAt.Add(constants.MaximumReversalDuration)
		if now.After(expiredReversalTime) {
			return errors.New("reversal duration is already expired")
		}
	}
	var (
		errUpdateBalance error
	)
	switch trx.TransactionType {
	case constants.TransactionTypeTopup:
		if req.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.WalletExternal.CreditBalance(ctx, tokenData, reqUpdateBalance)
		} else if req.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.WalletExternal.DebitBalance(ctx, tokenData, reqUpdateBalance)
		}
	case constants.TransactionTypePurchase:
		if req.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.WalletExternal.DebitBalance(ctx, tokenData, reqUpdateBalance)
		} else if req.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.WalletExternal.CreditBalance(ctx, tokenData, reqUpdateBalance)
		}
	}
	if errUpdateBalance != nil {
		return errors.Wrap(errUpdateBalance, "failed to update balance")
	}

	// update additional info
	var (
		newAdditionalInfo     = map[string]interface{}{}
		currentAdditionalInfo = map[string]interface{}{}
	)

	if trx.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(trx.AdditionalInfo), &currentAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}

	if req.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(req.AdditionalInfo), &newAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal new additional info")
		}
	}

	for key, val := range newAdditionalInfo {
		currentAdditionalInfo[key] = val
	}

	byteAdditionalInfo, err := json.Marshal(currentAdditionalInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal merged additional info")
	}

	// update status in DB
	err = s.TransactionRepository.UpdateStatusTransaction(ctx, req.Reference, req.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		return errors.Wrap(err, "failed to update status transaction")
	}

	trx.TransactionStatus = req.TransactionStatus
	// s.sendNotification(ctx, tokenData, trx)

	return nil
}
