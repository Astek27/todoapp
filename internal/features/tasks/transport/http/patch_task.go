package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_request "github.com/Astek27/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Astek27/todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Поспать"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Долго поспать"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can't be NULL")
		}

		lenTitle := len([]rune(*r.Title.Value))
		if lenTitle < 1 || lenTitle > 100 {
			return fmt.Errorf("length of title must be beetwen 1 and 100")
		}
	}

	if r.Description.Set {
		lenDescription := len([]rune(*r.Description.Value))
		if lenDescription < 1 || lenDescription > 1000 {
			return fmt.Errorf("length description must be beetwen 1 and 1000")
		}
	}

	if r.Completed.Set && r.Completed.Value == nil {
		return fmt.Errorf("completed can't be NULL")
	}

	return nil
}

type PatchedTaskResponse TaskDTOResponse

// PatchTask    godoc
// @Summary     Изменение задачи
// @Description Изменение информации о задаче
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description": "Tasty"` - устанавливается новое описание.
// @Description 3. **Передан null**: `"description": null` - очищает поле в БД (set NULL)
// @Description Ограничения: `title` и `completed` не может быть выставлен как null
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id path int true                     "ID задачи"
// @Param       request body PatchTaskRequest true   "PatchTask тело запроса"
// @Success     200 {object} PatchedTaskResponse                "Успешно измененная задача"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure     409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID from path value")
		return
	}

	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate")
		return
	}

	taskPatch := taskPatchFromRequest(request)
	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchedTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}