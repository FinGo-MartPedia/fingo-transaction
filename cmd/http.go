package cmd

import (
	"log"

	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	dependencies := InitDependency()

	r := gin.Default()

	apiV1 := r.Group("/api/v1/transaction")
	apiV1.POST("/create", dependencies.MiddlewareValidateToken, dependencies.TransactionController.CreateTransaction)
	apiV1.PUT("/update-status/:reference", dependencies.MiddlewareValidateToken, dependencies.TransactionController.UpdateStatusTransaction)
	apiV1.GET("/", dependencies.MiddlewareValidateToken, dependencies.TransactionController.GetTransactions)
	apiV1.GET("/:reference", dependencies.MiddlewareValidateToken, dependencies.TransactionController.GetTransactionDetail)
	apiV1.POST("/refund", dependencies.MiddlewareValidateToken, dependencies.TransactionController.RefundTransaction)

	err := r.Run(":" + helpers.GetEnv("PORT", "8082"))
	if err != nil {
		log.Fatal(err)
	}
}
