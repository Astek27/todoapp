package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	logger *core_logger.Logger
	rw http.ResponseWriter
}

func NewHTTPResponseHandler(
	logger *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		logger: logger,
		rw: rw,
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statusCode)

	response := map[string]string {
		"msg": msg,
		"err": err.Error(),
	}

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.logger.Error("json encode error:", zap.Error(err))
	}

}