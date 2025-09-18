package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/requests"
	"github.com/fingo-martpedia/fingo-transaction/internal/models/responses"
	"github.com/pkg/errors"
)

type WalletExternal struct {
}

func NewWalletExternal() *WalletExternal {
	return &WalletExternal{}
}

func (e *WalletExternal) CreditBalance(ctx context.Context, token string, req requests.UpdateBalance) (*responses.UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marhsal json")
	}

	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_CREDIT", "")

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http request")
	}
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error response from wallet service: %d", resp.StatusCode)
	}

	result := &responses.UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}

func (e *WalletExternal) DebitBalance(ctx context.Context, token string, req requests.UpdateBalance) (*responses.UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marhsal json")
	}

	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_DEBIT", "")

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http request")
	}
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect wallet service")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error response from wallet service: %d", resp.StatusCode)
	}

	result := &responses.UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()

	return result, nil
}
