package users_postgres_repository

import "github.com/Astek27/todoapp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func UserDomainsFromModels(userModels []UserModel) []domain.User {
	userDomains := make([]domain.User, len(userModels))
	for i, userModel := range userModels {
		userDomain := domain.NewUser(
			userModel.ID,
			userModel.Version,
			userModel.FullName,
			userModel.PhoneNumber,
		)
		userDomains[i] = userDomain
	}
	return userDomains
}