package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Astek27/todoapp/internal/core/domain"
	core_logger "github.com/Astek27/todoapp/internal/core/logger"
	core_http_request "github.com/Astek27/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Astek27/todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Astek27/todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Max Jon"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can not be null")
		}

		lenFullName := len([]rune(*r.FullName.Value))
		if lenFullName < 3 || lenFullName > 100 {
			return fmt.Errorf("length FullName min=3 and max=100")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			lenPhoneNumber := len([]rune(*r.PhoneNumber.Value))
			if lenPhoneNumber < 10 || lenPhoneNumber > 15 {
				return fmt.Errorf("length PhoneNumber min=10 and max=15")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must startwith +")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser    godoc
// @Summary     Изменение пользователе
// @Description Изменение информации о зарегистрированном пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number": "+79998887766"` - устанавливается новый номер телефона.
// @Description 3. **Передан null**: `"phone_number": null` - очищает поле в БД (set NULL)
// @Description Ограничения: `full_name` не может быть выставлен как null
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true                     "ID пользователя"
// @Param       request body PatchUserRequest true   "PatchUser тело запроса"
// @Success     200 {object} PatchUserResponse                "Успешно измененный пользователь"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User not found"
// @Failure     409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get path value")
		return
	}

	var patchUserRequest PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &patchUserRequest); err != nil {
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
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}