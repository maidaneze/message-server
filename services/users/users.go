package users

import (
	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/services/passwords"
)

//Validates if the usersRequestDTO is valid
//Returns true if it is, false if it isn't
//Maximum length for username and passwords fields is 32 characters
//Username and Password can't be empty

func ValidateUsersRequestDTO(dto model.UserRequestDTO) bool {
	if dto.Username == "" || dto.Password == "" || len(dto.Username) > model.MAX_USERNAME_FIELD_SIZE ||
		len(dto.Password) > model.MAX_PASSWORD_FIELD_SIZE {
		return false
	}
	return true
}

//Creates a new user with the given username and password
//The resulting user.UserId is 0, it must be set after inserting into the database
//The resulting user.Username is the given username
//The resulting user.Salt is a randomly insertd 32 byte string
//The resulting user.password is the given password encoded with SHA256 using the insertd salt

func CreateUser(username string, password string) (model.User, error) {
	//Encode password and insert salt
	encryptedPassword, salt, err := passwords.GenerateSecurePassword(password)
	return model.User{0, username, encryptedPassword, salt}, err
}
