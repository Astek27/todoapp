package tasks_transport_http

import (
	"time"

	"github.com/Astek27/todoapp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"          example:"1"`
	Version      int        `json:"version"     example:"3"`
	Title        string     `json:"title"       example:"Погулять"`
	Description  *string    `json:"description" example:"Погулять с бобиком"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id" example:"1"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,        
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromDomains(tasksDomain []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasksDomain))

	for i, task := range tasksDomain {
		dtos[i] = taskDTOFromDomain(task)
	}

	return  dtos
} 