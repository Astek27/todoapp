package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Astek27/todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.pool.Query(ctx, sqlQuery, limit, offset)
	if err != nil {
		return []domain.User{}, fmt.Errorf("select users from database: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel
		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return []domain.User{}, fmt.Errorf("rows scan: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if rows.Err() != nil {
		return []domain.User{}, fmt.Errorf("next rows: %w", err)
	}
	
	userDomains := UserDomainsFromModels(userModels)
	return userDomains, nil
}