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

	jsonAdditionalInfo := map[string]interface{}{}
	if req.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(req.AdditionalInfo), &jsonAdditionalInfo)
		if err != nil {
			return resp, errors.Wrap(err, "additional info type is invalid")
		}
	}

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

	// request update balance to wallet service
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

func (s *TransactionService) GetTransactions(ctx context.Context, userID int) ([]responses.TransactionResponse, error) {
	transactions, err := s.TransactionRepository.GetTransactions(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transactions")
	}

	transactionResponses := make([]responses.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = responses.TransactionResponse{
			Amount:            transaction.Amount,
			TransactionType:   transaction.TransactionType,
			TransactionStatus: transaction.TransactionStatus,
			Reference:         transaction.Reference,
			Description:       transaction.Description,
			AdditionalInfo:    transaction.AdditionalInfo,
		}
	}

	return transactionResponses, nil
}

func (s *TransactionService) GetTransactionDetail(ctx context.Context, reference string) (responses.TransactionResponse, error) {
	var resp responses.TransactionResponse

	transaction, err := s.TransactionRepository.GetTransactionByReference(ctx, reference, true)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get transaction detail")
	}

	resp = responses.TransactionResponse{
		Amount:            transaction.Amount,
		TransactionType:   transaction.TransactionType,
		TransactionStatus: transaction.TransactionStatus,
		Reference:         transaction.Reference,
		Description:       transaction.Description,
		AdditionalInfo:    transaction.AdditionalInfo,
	}
	return resp, nil
}

func (s *TransactionService) RefundTransaction(ctx context.Context, token string, userId int, req *requests.RefundTransaction) (responses.CreateTransactionResponse, error) {
	var (
		resp responses.CreateTransactionResponse
	)

	trx, err := s.TransactionRepository.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get transaction")
	}

	if trx.TransactionStatus != constants.TransactionStatusSuccess || trx.TransactionType != constants.TransactionTypePurchase {
		return resp, errors.New("current transaction status is not success or transaction type is not purchase")
	}

	refundReference := "REFUND-" + req.Reference
	reqCreditBalance := requests.UpdateBalance{
		Reference: refundReference,
		Amount:    trx.Amount,
	}
	_, err = s.WalletExternal.CreditBalance(ctx, token, reqCreditBalance)
	if err != nil {
		return resp, errors.Wrap(err, "failed to credit balance")
	}

	transaction := models.Transaction{
		UserID:            userId,
		Amount:            trx.Amount,
		TransactionType:   constants.TransactionTypeRefund,
		TransactionStatus: constants.TransactionStatusSuccess,
		Reference:         refundReference,
		Description:       req.Description,
		AdditionalInfo:    req.AdditionalInfo,
	}
	err = s.TransactionRepository.CreateTransaction(ctx, &transaction)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert new transaction refund")
	}

	resp.Reference = refundReference
	resp.TransactionStatus = transaction.TransactionStatus

	return resp, nil
}
