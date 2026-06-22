package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_http_server "github.com/Astek27/todoapp/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, user domain.Task) (domain.Task, error)
}

func NewTasksHTTPHandler(tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.Route{
			Method: http.MethodPost,
			Path: "/tasks",
			Handler: h.CreateTask,
		},
	}
}