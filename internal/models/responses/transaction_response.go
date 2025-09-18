package responses

type CreateTransactionResponse struct {
	Reference         string `json:"reference"`
	TransactionStatus string `json:"transaction_status"`
}

type TransactionResponse struct {
	Amount            float64 `json:"amount"`
	TransactionType   string  `json:"transaction_type"`
	TransactionStatus string  `json:"transaction_status"`
	Reference         string  `json:"reference"`
	Description       string  `json:"description"`
	AdditionalInfo    string  `json:"additional_info"`
}
