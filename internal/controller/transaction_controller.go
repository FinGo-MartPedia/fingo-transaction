package controller

import (
	"net/http"

	"github.com/fingo-martpedia/fingo-transaction/constants"
	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/fingo-martpedia/fingo-transaction/internal/interfaces"
	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService interfaces.ITransactionService
}

func NewTransactionController(transactionService interfaces.ITransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

func (api *TransactionController) CreateTransaction(c *gin.Context) {
	var (
		log    = helpers.Logger
		req    models.Transaction
		errMsg string
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request: ", err)
		errMsg = err.Error()
		helpers.SendResponse(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil, &errMsg)
		return
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request: ", err)
		errMsg = err.Error()
		helpers.SendResponse(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil, &errMsg)
		return
	}

	u, exists := c.Get("user")
	if !exists {
		log.Error("Failed to get user from context")
		errMsg = "failed to get user from context"
		helpers.SendResponse(c, http.StatusUnauthorized, constants.ErrFailedUnauthorized, nil, &errMsg)
		return
	}

	user, ok := u.(models.User)
	if !ok {
		log.Error("Invalid user type in context")
		errMsg = "invalid user type in context"
		helpers.SendResponse(c, http.StatusUnauthorized, constants.ErrFailedUnauthorized, nil, &errMsg)
		return
	}

	if !constants.MapTransactionType[req.TransactionType] {
		log.Error("invalid transaction type")
		errMsg = "invalid transaction type"
		helpers.SendResponse(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil, &errMsg)
		return
	}

	req.UserID = int(user.ID)

	resp, err := api.TransactionService.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		log.Error("failed to create transaction: ", err)
		errMsg = err.Error()
		helpers.SendResponse(c, http.StatusInternalServerError, constants.ErrFailedServerError, nil, &errMsg)
		return
	}

	helpers.SendResponse(c, http.StatusOK, constants.SuccessMessage, resp, nil)
}
