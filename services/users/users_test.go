package users

import (
	"challenge/model"
	"testing"

	"challenge/services/passwords"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsersRequestDTO(t *testing.T) {

	emptyDTO := model.UserRequestDTO{}

	validDTO1 := emptyDTO
	validDTO1.Username = "user1"
	validDTO1.Password = "pass1"
	validDTO2 := emptyDTO
	validDTO2.Username = "user2"
	validDTO2.Password = "pass2"

	hugePasswordString := string(make([]byte, model.MAX_PASSWORD_FIELD_SIZE*10))
	hugeUsernameString := string(make([]byte, model.MAX_USERNAME_FIELD_SIZE*10))

	invalidPasswordDTO1 := validDTO1
	invalidPasswordDTO1.Password = ""
	invalidPasswordDTO2 := validDTO1
	invalidPasswordDTO2.Password = hugePasswordString

	invalidUserDTO1 := validDTO1
	invalidUserDTO1.Username = ""
	invalidUserDTO2 := validDTO1
	invalidUserDTO2.Username = hugeUsernameString

	cases := []struct {
		name     string
		message  model.UserRequestDTO
		expected bool
	}{
		{"testValidateUsersRequestDTOEmptyBody", emptyDTO, false},
		{"testValidateUsersRequestDTOValidDTO1", validDTO1, true},
		{"testValidateUsersRequestDTOValidDTO2", validDTO2, true},
		{"testValidateUsersRequestDTOInvalidPasswordDTO1", invalidPasswordDTO1, false},
		{"testValidateUsersRequestDTOInvalidPasswordDTO2", invalidPasswordDTO1, false},
		{"testValidateUsersRequestDTOInvalidUserDTO1", invalidUserDTO1, false},
		{"testValidateUsersRequestDTOInvalidUserDTO2", invalidUserDTO2, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			actual := ValidateUsersRequestDTO(c.message)
			assert.Equal(tt, c.expected, actual)

		})
	}
}

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testCreateUserEmptyUserAndEmptyPassword", "", ""},
		{"testCreateUserEmptyUserAndNonEmptyPassword", "", "pass"},
		{"testCreateUserNonEmptyUserAndEmptyPassword", "user", ""},
		{"testCreateUserNonEmptyUserAndNonEmptyPassword", "user", "pass"},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			user, err := CreateUser(c.username, c.password)
			assert.Nil(tt, err)
			assert.Equal(tt, int64(0), user.Userid)
			assert.Equal(tt, c.username, user.Username)
			assert.Equal(tt, passwords.GetEncryptedPassword(c.password, user.Salt), user.Password)
			assert.True(tt, len(user.Salt) > 0)
		})
	}
}
