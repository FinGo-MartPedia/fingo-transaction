package interfaces

import (
	"context"

	"github.com/fingo-martpedia/fingo-transaction/internal/models/requests"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
)

type IWalletExternal interface {
	CreditBalance(ctx context.Context, token string, req requests.UpdateBalance) (*responses.UpdateBalanceResponse, error)
	DebitBalance(ctx context.Context, token string, req requests.UpdateBalance) (*responses.UpdateBalanceResponse, error)
}
