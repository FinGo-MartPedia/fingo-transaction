package requests

type UpdateBalance struct {
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}
