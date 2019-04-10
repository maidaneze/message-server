package passwords

import (
	"github.com/maidaneze/message-server/utils"
	"crypto/rand"
	"time"

	"github.com/maidaneze/message-server/model"

	"github.com/mattn/go-sqlite3"
)

//Generates a secure password to safely store it on a database using a 32 bit random salt
//Returns the encrypted password and the random salt used on success
//Returns an error if it fails

func GenerateSecurePassword(password string) (encryptedPassword string, salt string, err error) {
	randbytes := make([]byte, 32)

	utils.Retry(func() error {
		_, err = rand.Read(randbytes)
		return err
	}, 2, 20*time.Millisecond)

	if err != nil {
		return "", "", err
	}

	salt = string(randbytes)

	return GetEncryptedPassword(password, salt), salt, nil
}

//Returns the encrypted password Using SHA256 and the given salt
//The salt should be a 32 bit integer

func GetEncryptedPassword(password, salt string) string {
	encrypt := sqlite3.CryptEncoderSSHA256(salt)
	bytes := encrypt([]byte(password), nil)
	encryptedPassword := string(bytes)
	return encryptedPassword
}

//Returns true if the given password matches the encrypted password and false otherwise

func ValidUserPassword(user model.User, password string) bool {
	return GetEncryptedPassword(password, user.Salt) == user.Password
}
