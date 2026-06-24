package tasks_service

import (
	"context"
	"fmt"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx    context.Context,
	userID *int,
	limit  *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be not-negative: %w",
			core_errors.ErrBadRequest,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be not-negative: %w",
			core_errors.ErrBadRequest,
		)
	}

	users, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"get users from repository: %w",
			err,
		)
	}

	return users, nil
}