package users_service

import (
	"context"
	"fmt"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func (s *UsersService) GetUsers (
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return []domain.User{}, fmt.Errorf(
			"limit must be not-negative: %w",
			core_errors.ErrBadRequest,
		)
	}

	if offset != nil && *offset < 0 {
		return []domain.User{}, fmt.Errorf(
			"offset must be not-negative: %w",
			core_errors.ErrBadRequest,
		)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return []domain.User{}, fmt.Errorf("get users from repo: %w", err)
	}

	return users, nil
}