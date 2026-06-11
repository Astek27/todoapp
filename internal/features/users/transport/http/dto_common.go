package users_transport_http

import "github.com/Astek27/todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int      `json:"id"`
	Version     int      `json:"version"`
	FullName    string   `json:"full_name"`
	PhoneNumber *string  `json:"phone_number"`
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