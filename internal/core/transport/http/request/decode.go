package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var RequestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decoce json: %w", err)
	}
	
	if err := RequestValidator.Struct(dest); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}
