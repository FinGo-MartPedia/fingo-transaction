package responses

type UpdateBalanceResponse struct {
	Message string `json:"message"`
	Data    struct {
		Balance float64 `json:"balance"`
	} `json:"data"`
}
