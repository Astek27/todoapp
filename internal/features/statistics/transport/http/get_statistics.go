package statistics_http_transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_request "github.com/Astek27/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
)

type StatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"                 example:"50"`
	TasksCompleted             int      `json:"tasks_completed"               example:"30"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"          example:"60"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

// GetStatistics godoc
// @Summary     Получение статистики
// @Description Получение статистики с опциональной фильтрацией по user_id и/или промежутку времени
// @Tags        statistics
// @Produce     json
// @Param       user_id query int    false                     "Фильтрация статистики по конкретному пользователю"
// @Param       from query    string false                     "Начало промежутка рассмотрения статистики (включительно), формат: ГГГГ-ММ-ДД"
// @Param       to query      string false                     "Конец промежутка рассмотрения статистики (не включительно), формат: ГГГГ-ММ-ДД"
// @Success     200 {object}  StatisticsResponse               "Успешное получение статистики"
// @Failure     400 {object}  core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object}  core_http_response.ErrorResponse "Internal server error"
// @Router      /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "get statistics error")
		return
	}

	response := toDTOFromDomain(statistics)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) StatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	
	return StatisticsResponse{
		TasksCreated: statistics.TasksCreated,
		TasksCompleted: statistics.TasksCompleted,
		TasksCompletedRate: statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	var (
		userIDQueryParamKey = "user_id"
		toQueryParamKey     = "to"
		fromQueryParamKey = "from"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	to, err := core_http_request.GetTimeQueryParams(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	from, err := core_http_request.GetTimeQueryParams(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	return userID, from, to, nil
}