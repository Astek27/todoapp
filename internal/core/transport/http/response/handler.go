package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
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

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.logger.Error("json encode:", zap.Error(err))
		return
	}
	
}

func (h *HTTPResponseHandler) HTMLResponse(html []byte) {
	h.rw.WriteHeader(http.StatusOK)
	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := h.rw.Write(html); err != nil {
		h.logger.Error("write HTML HTTP response:", zap.Error(err))
		return
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		loggerFunc func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrBadRequest):
		statusCode = http.StatusBadRequest
		loggerFunc = h.logger.Warn
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		loggerFunc = h.logger.Debug
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		loggerFunc = h.logger.Warn
	default:
		statusCode = http.StatusInternalServerError
		loggerFunc = h.logger.Error
	}

	loggerFunc(msg, zap.Error(err))
	
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := ErrorResponse{
		Error: err.Error(),
		Message: msg,
	}

	h.JSONResponse(response, statusCode)
}