package user

import (
	"context"

	"github.com/stasdashkevitch/rest-api/cmd/pkg/logging"
)

type Service struct {
	storage Storage
	logger  logging.Logger
}

func (s *Service) Create(ctx *context.Context, dto CreateUserDTO) (user User, err error) {
	return
}
