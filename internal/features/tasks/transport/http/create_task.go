package tasks_transport_http

import (
	"net/http"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_request "github.com/Astek27/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string `json:"title" validate:"required,min=1,max=100" example:"Покушать"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Вкусно покушать"`
	AuthorUserID int    `json:"author_user_id" validate:"required" example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask   godoc
// @Summary     Создать задачу
// @Description Создать новую задачу в системе
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body CreateTaskRequest true "CreateTask тело запроса"
// @Success     201 {object} CreateTaskResponse "Успешно созданная задача"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [post]
func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskDomain := domain.NewTaskUninitilized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}