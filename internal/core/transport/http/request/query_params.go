package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"Param %s in key=%s not integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrBadRequest,
		)
	}

	return &val, nil
}
