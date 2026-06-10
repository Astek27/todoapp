package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var RequestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decoce json: %v, %w",
			err,
			core_errors.ErrBadRequest,
		)
	}
	
	if err := RequestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"validate: %v, %w",
			err,
			core_errors.ErrBadRequest,
		)
	}

	return nil
}
