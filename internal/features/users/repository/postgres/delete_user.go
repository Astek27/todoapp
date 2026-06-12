package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1
	`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete user id=%d: %w", id, core_errors.ErrNotFound)
	}
	return nil
}