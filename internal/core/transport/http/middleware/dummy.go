package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
)

func Dummy(s string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			log.Debug(fmt.Sprintf("-> before: %s", s))
			h.ServeHTTP(w, r)
			log.Debug(fmt.Sprintf("-> after: %s", s))
		})
	}
}