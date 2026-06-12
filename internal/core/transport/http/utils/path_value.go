package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("not path value key=%s: %w", core_errors.ErrBadRequest)
	}

	pathValueInt, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path_value=%s by key=%s not integer: %s: %w",
			pathValue,
			key,
			err,
			core_errors.ErrBadRequest,
		)
	}

	return pathValueInt, nil
}