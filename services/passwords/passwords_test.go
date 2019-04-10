package passwords

import (
	"github.com/maidaneze/message-server/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEncryptedPasswordIsDifferentFromOriginalPassword(t *testing.T) {
	pass := "pass"
	salt := "salt"

	encrypted1 := GetEncryptedPassword(pass, salt)

	assert.NotEqual(t, pass, encrypted1)
}

func TestGetEncryptedPasswordEncryptedWithDifferentSaltsAreDifferent(t *testing.T) {
	pass := "pass"
	salt1 := "salt1"
	salt2 := "salt2"

	encrypted1 := GetEncryptedPassword(pass, salt1)
	encrypted2 := GetEncryptedPassword(pass, salt2)

	assert.NotEqual(t, encrypted1, encrypted2)
}

func TestGetEncryptedPasswordDifferentPasswordsEncryptedWithTheSameSaltAreDifferent(t *testing.T) {
	pass1 := "pass1"
	pass2 := "pass2"
	salt := "salt"

	encrypted1 := GetEncryptedPassword(pass1, salt)
	encrypted2 := GetEncryptedPassword(pass2, salt)

	assert.NotEqual(t, encrypted1, encrypted2)
}

func TestGetEncryptedPasswordEncryptingTwiceGivesTheSameEncryptedResult(t *testing.T) {
	pass := "pass"
	salt := "salt"

	encrypted1 := GetEncryptedPassword(pass, salt)
	encrypted2 := GetEncryptedPassword(pass, salt)

	assert.Equal(t, encrypted1, encrypted2)
}

func TestGenerateEncryptedPasswordCanBeRecoveredByEncryptingThePasswordWithTheSalt(t *testing.T) {
	pass := "pass"

	encrypted1, salt, err := GenerateSecurePassword(pass)

	require.Nil(t, err)

	assert.NotEqual(t, encrypted1, pass)

	assert.True(t, ValidUserPassword(model.User{0, "", encrypted1, salt}, pass))

	assert.Equal(t, encrypted1, GetEncryptedPassword(pass, salt))

	assert.False(t, ValidUserPassword(model.User{0, "", encrypted1, salt}, "invalid"))

	assert.False(t, ValidUserPassword(model.User{0, "", encrypted1, "invalid"}, pass))
}

func TestGenerateEncryptedPassword2TimesGeneratesDifferentEncryptedPasswords(t *testing.T) {
	pass := "pass"

	encrypted1, salt1, err1 := GenerateSecurePassword(pass)
	encrypted2, salt2, err2 := GenerateSecurePassword(pass)

	require.Nil(t, err1)
	require.Nil(t, err2)

	assert.NotEqual(t, encrypted1, encrypted2)
	assert.NotEqual(t, salt1, salt2)
}
