package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_request "github.com/Astek27/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Astek27/todoapp/internal/core/transport/http/types"
	core_http_utils "github.com/Astek27/todoapp/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get path value")
		return
	}

	var patchUserRequest PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, patchUserRequest); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(patchUserRequest)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest:\nFullName:'%s'\nPhoneNumber:'%s'",
			patchUserRequest.FullName,
			patchUserRequest.PhoneNumber,
		),
	)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}