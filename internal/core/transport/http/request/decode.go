package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var RequestValidator = validator.New()

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decoce json: %w", err)
	}
	
	return nil
}
