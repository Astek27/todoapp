package tasks_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec sql query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user with id=%d: %w", id, core_errors.ErrNotFound)
	}

	return nil
}