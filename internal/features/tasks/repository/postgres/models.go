package tasks_postgres_repository

import (
	"time"

	"github.com/Astek27/todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainFromModel(task TaskModel) domain.Task {
	return domain.NewTask(
			task.ID,
			task.Version,
			task.Title,
			task.Description,
			task.Completed,
			task.CreatedAt,
			task.CompletedAt,
			task.AuthorUserID,
		)
}

func tasksDomainFromModels(tasksModel []TaskModel) []domain.Task {
	tasksDomain := make([]domain.Task, len(tasksModel))

	for i, task := range tasksModel {
		tasksDomain[i] = taskDomainFromModel(task)
	}
	return tasksDomain
}