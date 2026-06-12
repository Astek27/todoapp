package users_transport_http

import (
	"net/http"

	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/Astek27/todoapp/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	id, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid id")
		return
	}

	if err := h.usersService.DeleteUser(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, "not delete user")
		return
	}

	responseHandler.NoContentResponse()
}