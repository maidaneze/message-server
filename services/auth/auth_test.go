package auth

import (
	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/utils"
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tokens := []model.Token{}
	for i := 0; i < 10; i++ {
		token, err := GenerateToken()
		assert.True(t, tokenExpirationInRange(token))
		assert.Nil(t, err)
		assert.True(t, properTokenLength(token))
		tokens = append(tokens, token)
	}
	assert.True(t, allTokensAreDifferent(tokens))
}

func TestValidateTokenHeader(t *testing.T) {
	bearer := "Bearer "
	token1, _ := GenerateToken()
	token2, _ := GenerateToken()
	cases := []struct {
		name          string
		token         string
		expectedToken string
		expectedValid bool
	}{
		{"testValidateTokenHeaderEmptyToken", "", "", false},
		{"testValidateTokenHeaderEmptyTokenWithBearer", bearer + "", "", true},
		{"testValidateTokenHeaderTokenWithBearerAndValidToken1", bearer + token1.Uuid, token1.Uuid, true},
		{"testValidateTokenHeaderTokenWithBearerAndValidToken2", bearer + token2.Uuid, token2.Uuid, true},
		{"testValidateTokenHeaderTokenWithoutBearerAndValidToken1", token1.Uuid, "", false},
		{"testValidateTokenHeaderTokenWithoutBearerAndValidToken2", token2.Uuid, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			req := &http.Request{}
			req.Header = http.Header{}
			req.Header.Set("Authorization", c.token)
			token, valid := ValidateTokenHeader(req)
			assert.Equal(tt, c.expectedToken, token)
			assert.Equal(tt, c.expectedValid, valid)
		})
	}
}

func TestFindToken(t *testing.T) {
	emptyToken := model.Token{}
	token1 := model.Token{"123", 1}
	token2 := model.Token{"456", 2}
	token3 := model.Token{"789", 3}
	tokensSlice := []model.Token{token1, token2, token3}
	cases := []struct {
		name   string
		token  model.Token
		tokens []model.Token
		found  bool
	}{
		{"testFindTokenOnEmptyList", emptyToken, tokensSlice, false},
		{"testFindTokenReturnsTrue1", token1, tokensSlice, true},
		{"testFindTokenReturnsTrue2", token2, tokensSlice, true},
		{"testFindTokenReturnsTrue3", token3, tokensSlice, true},
		{"testFindTokenReturnsFalse1", token1, []model.Token{token2, token3}, false},
		{"testFindTokenReturnsFalse2", token2, []model.Token{token1, token3}, false},
		{"testFindTokenReturnsFalse3", token3, []model.Token{token1, token2}, false},
		{"testFindTokenReturnsFalse4", model.Token{token1.Uuid, 0}, tokensSlice, false},
		{"testFindTokenReturnsFalse5", model.Token{"invalid", token1.Expiration}, tokensSlice, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			assert.Equal(tt, c.found, FindToken(c.token, c.tokens))
		})
	}
}

func TestValidateAuthorizedUser(t *testing.T) {
	token1, _ := GenerateToken()
	token2, _ := GenerateToken()
	emptyTokens := []model.Token{}
	tokens1 := []model.Token{token1, token2}
	tokens2 := []model.Token{token2, token1}
	expiredTokens1 := []model.Token{token1, token2}
	expiredTokens1[0].Expiration = 0
	expiredTokens2 := []model.Token{token1, token2}
	expiredTokens2[1].Expiration = 0
	invalidToken1, _ := GenerateToken()
	invalidToken2, _ := GenerateToken()
	cases := []struct {
		name     string
		token    model.Token
		tokens   []model.Token
		expected bool
	}{
		{"testValidateAuthorizedUserEmptyTokens1", token1, emptyTokens, false},
		{"testValidateAuthorizedUserEmptyTokens2", token2, emptyTokens, false},
		{"testValidateAuthorizedUserEmptyTokens3", invalidToken1, emptyTokens, false},
		{"testValidateAuthorizedUserEmptyTokens4", invalidToken2, emptyTokens, false},
		{"testValidateAuthorizedUserValidTokens1", token1, tokens1, true},
		{"testValidateAuthorizedUserValidTokens2", token2, tokens1, true},
		{"testValidateAuthorizedUserValidTokens3", invalidToken1, tokens1, false},
		{"testValidateAuthorizedUserValidTokens4", invalidToken2, tokens1, false},
		{"testValidateAuthorizedUserValidTokens5", token1, tokens2, true},
		{"testValidateAuthorizedUserValidTokens6", token2, tokens2, true},
		{"testValidateAuthorizedUserValidTokens7", invalidToken1, tokens2, false},
		{"testValidateAuthorizedUserValidTokens8", invalidToken2, tokens2, false},
		{"testValidateAuthorizedUserExpiredTokens1", token1, expiredTokens1, false},
		{"testValidateAuthorizedUserExpiredTokens2", token2, expiredTokens1, true},
		{"testValidateAuthorizedUserExpiredTokens3", invalidToken1, expiredTokens1, false},
		{"testValidateAuthorizedUserExpiredTokens4", invalidToken2, expiredTokens1, false},
		{"testValidateAuthorizedUserExpiredTokens5", token1, expiredTokens2, true},
		{"testValidateAuthorizedUserExpiredTokens6", token2, expiredTokens2, false},
		{"testValidateAuthorizedUserExpiredTokens7", invalidToken1, expiredTokens2, false},
		{"testValidateAuthorizedUserExpiredTokens8", invalidToken2, expiredTokens2, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			assert.Equal(tt, c.expected, ValidateAuthorizedUser(c.token.Uuid, c.tokens))
		})
	}

}

func tokenExpirationInRange(token model.Token) bool {
	now := utils.UTCTimeMilliseconds()
	return token.Expiration <= now+tokenDurationMilliseconds && token.Expiration > now-1000
}

func properTokenLength(token model.Token) bool {
	length := len([]byte(token.Uuid))
	_ = length
	return len([]byte(token.Uuid)) == TokenLengthInBytes
}

func allTokensAreDifferent(tokens []model.Token) bool {
	for i := 0; i < len(tokens); i++ {
		for j := i + 1; j < len(tokens); j++ {
			if tokens[i].Uuid == tokens[j].Uuid {
				return false
			}
		}
	}
	return true
}
