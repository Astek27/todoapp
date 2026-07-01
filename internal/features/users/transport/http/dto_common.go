package users_transport_http

import "github.com/Astek27/todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int      `json:"id"           example:"10"`
	Version     int      `json:"version"      example:"3"`
	FullName    string   `json:"full_name"    example:"Ivan Ivanov"`
	PhoneNumber *string  `json:"phone_number" example:"+79998887766"`
}

func userDTOFromDomain(userDomain domain.User) UserDTOResponse {
	return  UserDTOResponse{
		ID:          userDomain.ID,
		Version:     userDomain.Version,
		FullName:    userDomain.FullName,
		PhoneNumber: userDomain.PhoneNumber,
	}
}

func usersDTOFromDomain(userDomains []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(userDomains))
	
	for i, userDomain := range userDomains {
		usersDTO[i] = userDTOFromDomain(userDomain)
	}

	return usersDTO
}