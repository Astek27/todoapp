package domain

import (
	"fmt"
	"time"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

type Task struct {
	ID            int
	Version       int

	Title         string
	Description  *string
	Completed     bool
	CreatedAt    time.Time
	CompletedAt *time.Time

	AuthorUserID  int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID: id,
		Version: version,
		Title: title,
		Description: description,
		Completed: completed,
		CreatedAt: createdAt,
		CompletedAt: completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitilized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninializedID,
		UninializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLen,
			core_errors.ErrBadRequest,
		)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 100 {
			return fmt.Errorf(
				"invalid `Title` len: %d: %w",
				descriptionLen,
				core_errors.ErrBadRequest,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`ComletedAt` can not be `nil` if `Completed` == `true`: %w",
				core_errors.ErrBadRequest,
			)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"CreatedAt must be <= CompletedAt: %w",
				core_errors.ErrBadRequest,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"CompletedAt must be nil if task not completed: %w",
				core_errors.ErrBadRequest,
			)
		}
	}

	return nil
}