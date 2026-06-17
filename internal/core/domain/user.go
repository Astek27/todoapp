package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Astek27/todoapp/internal/core/errors"
)

type User struct {
	ID int
	Version int

	FullName string
	PhoneNumber *string
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninializedID,
		UninializedVersion,
		fullName,
		phoneNumber,
	)
}

func NewUser(
		id          int,
		version     int,
		fullName    string,
		phoneNumber *string,
	) User {
	return User{
		ID: id,
		Version: version,

		FullName: fullName,
		PhoneNumber: phoneNumber,
	}
}

func (d *User) Validate() error {
	lengthFullName := len([]rune(d.FullName))

	if lengthFullName < 3 || lengthFullName > 100 {
		return fmt.Errorf(
			"invalid length full name %d: %w",
			lengthFullName,
			core_errors.ErrBadRequest,
		)
	}

	if d.PhoneNumber != nil {
		lengthPhoneNumber := len([]rune(*d.PhoneNumber))
		if lengthPhoneNumber < 10 || lengthPhoneNumber > 15 {
			return fmt.Errorf(
				"invalid length phone number: %d. %w",
				lengthPhoneNumber,
				core_errors.ErrBadRequest,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*d.PhoneNumber) {
			return fmt.Errorf("invalid phone number format: %w", core_errors.ErrBadRequest)
		}
	}
	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("FullName not be null: %w", core_errors.ErrBadRequest)
	}
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("No valide patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}