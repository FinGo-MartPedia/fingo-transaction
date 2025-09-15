package requests

import "github.com/go-playground/validator/v10"

type UpdateStatusTransaction struct {
	Reference         string `json:"reference" valid:"required"`
	TransactionStatus string `json:"transaction_status" valid:"required"`
	AdditionalInfo    string `json:"additional_info"`
}

func (l UpdateStatusTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type RefundTransaction struct {
	Reference      string `json:"reference" valid:"required"`
	Description    string `json:"description" valid:"required"`
	AdditionalInfo string `json:"additional_info"`
}

func (l RefundTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
