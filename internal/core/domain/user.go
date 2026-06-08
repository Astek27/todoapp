package domain

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