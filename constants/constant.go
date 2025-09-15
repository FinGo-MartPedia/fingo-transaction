package constants

import "time"

const (
	SuccessMessage        = "Success"
	ErrFailedBadRequest   = "Bad Request"
	ErrFailedServerError  = "Internal Server Error"
	ErrFailedUnauthorized = "Unauthorized"
)

const (
	TransactionStatusPending  = "PENDING"
	TransactionStatusSuccess  = "SUCCESS"
	TransactionStatusFailed   = "FAILED"
	TransactionStatusReversed = "REVERSED"
)

const (
	TransactionTypeTopup    = "TOPUP"
	TransactionTypePurchase = "PURCHASE"
	TransactionTypeRefund   = "REFUND"
)

var MapTransactionType = map[string]bool{
	TransactionTypeTopup:    true,
	TransactionTypePurchase: true,
	TransactionTypeRefund:   true,
}

var MapTransactionStatusFlow = map[string][]string{
	TransactionStatusPending: {TransactionStatusSuccess, TransactionStatusFailed},
	TransactionStatusSuccess: {TransactionStatusReversed},
	TransactionStatusFailed:  {TransactionStatusSuccess},
}

const (
	MaximumReversalDuration = time.Hour * 24
)
