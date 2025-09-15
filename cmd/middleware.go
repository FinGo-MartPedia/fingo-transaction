package cmd

import (
	"log"
	"net/http"

	"github.com/fingo-martpedia/fingo-transaction/external"
	"github.com/fingo-martpedia/fingo-transaction/helpers"
	"github.com/gin-gonic/gin"
)

func (d *Dependency) MiddlewareValidateToken(ctx *gin.Context) {
	var errMsg string
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("authorization empty")
		errMsg = "authorization empty"
		helpers.SendResponse(ctx, http.StatusUnauthorized, "unauthorized", nil, &errMsg)
		ctx.Abort()
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		log.Println("invalid authorization format")
		errMsg = "invalid authorization format"
		helpers.SendResponse(ctx, http.StatusUnauthorized, "unauthorized", nil, &errMsg)
		ctx.Abort()
		return
	}

	accessToken := authHeader[len(bearerPrefix):]

	userData, err := external.ValidateToken(ctx.Request.Context(), accessToken)
	if err != nil {
		log.Println(err)
		errMsg = err.Error()
		helpers.SendResponse(ctx, http.StatusUnauthorized, "unauthorized", nil, &errMsg)
		ctx.Abort()
		return
	}

	ctx.Set("user", userData)

	ctx.Next()
}
