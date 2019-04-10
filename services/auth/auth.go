package auth

import (
	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/utils"
	"time"

	"net/http"

	"strings"

	"github.com/google/uuid"
)

var (
	tokenDurationMilliseconds = int64(time.Hour) * 24 * 365 / int64(time.Millisecond)
	TokenLengthInBytes        = 36
)

//Generates a Token
//The uuid has 122 random bits and a length of 36 bytes
//Returns error in case of failiure

func GenerateToken() (model.Token, error) {
	var err error
	var tokenUUID uuid.UUID
	utils.Retry(func() error {
		tokenUUID, err = uuid.NewRandom()
		return err
	}, 2, 20*time.Millisecond)

	if err != nil {
		return model.Token{}, err
	}
	expiration := utils.UTCTimeMilliseconds() + tokenDurationMilliseconds
	token := model.Token{tokenUUID.String(), expiration}
	return token, nil
}

//Validates that the "Authorization" header for the requests contains a token with the format bearer
//"Bearer <token>"
//Returns true in case it is a valid token and false otherwise

func ValidateTokenHeader(r *http.Request) (string, bool) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", false
	}
	reqToken = splitToken[1]
	return reqToken, true
}

//Returns true if the token is in tokens and false otherwise

func FindToken(token model.Token, tokens []model.Token) bool {
	found := false
	for _, value := range tokens {
		if (token.Uuid == value.Uuid) && (token.Expiration == value.Expiration) {
			found = true
		}
	}
	return found
}

//Returns true if the requestToken is one of the userTokens and it hasn't expired and returns false otherwise

func ValidateAuthorizedUser(requestToken string, userTokens []model.Token) bool {
	for _, value := range userTokens {
		if (value.Uuid == requestToken) && (value.Expiration >= utils.UTCTimeMilliseconds()) {
			return true
		}
	}
	return false
}
