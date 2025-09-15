package external

import (
	"context"
	"fmt"

	"github.com/fingo-martpedia/fingo-transaction/cmd/proto/tokenvalidation"
	"github.com/fingo-martpedia/fingo-transaction/constants"
	"github.com/fingo-martpedia/fingo-transaction/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func ValidateToken(ctx context.Context, token string) (models.User, error) {
	var user models.User

	conn, err := grpc.Dial("localhost:7000", grpc.WithInsecure())
	if err != nil {
		return user, errors.Wrap(err, "failed connect to server grpc ums")
	}
	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)

	req := &tokenvalidation.TokenRequest{Token: token}
	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		return user, errors.Wrap(err, "failed to validate token")
	}

	if response.Message != constants.SuccessMessage {
		return user, fmt.Errorf("failed to validate token: %s", response.Message)
	}

	user.ID = response.Data.UserId
	user.Email = response.Data.Email
	user.Username = response.Data.Username
	user.FullName = response.Data.FullName

	return user, nil
}
