package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx    context.Context,
	userID *int,
	from   *time.Time,
	to     *time.Time,
) (domain.Statistics, error) {
	if to != nil && from != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from': %w",
				core_errors.ErrBadRequest,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get task from repository: %w", err)
	}

	statistics := calculateStatistics(tasks)
	return statistics, nil
}

func calculateStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var tasksCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}
		complitionDuration := task.CompletionDuration()
		if complitionDuration != nil {
			tasksCompletedDuration += *complitionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && tasksCompletedDuration != 0 {
		avg := tasksCompletedDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)
}