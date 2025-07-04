package mocks

import "asniki/snippetbox/internal/models"

// UserModel mocks models.UserModel
type UserModel struct{}

// Insert mocks models.UserModel.Insert
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate mocks models.UserModel.Authenticate
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "bob@example.com" && password == "validPa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

// Exists mocks models.UserModel.Exists
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
