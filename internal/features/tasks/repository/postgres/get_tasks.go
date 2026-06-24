package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Astek27/todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx    context.Context,
	userID *int,
	limit  *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	args := []any{limit, offset}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer rows.Close()

	var tasksModel []TaskModel
	for rows.Next() {
		var task TaskModel
		if err := rows.Scan(
			&task.ID,
			&task.Version,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasksModel = append(tasksModel, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	tasksDomain := tasksDomainFromModels(tasksModel)
	return tasksDomain, nil
}