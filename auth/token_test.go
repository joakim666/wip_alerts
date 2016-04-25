package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAccessToken(t *testing.T) {
	assert := assert.New(t)

	ser, originalToken, err := createEncryptedTestToken()
	assert.NoError(err)

	token, err := DecryptAccessToken(ser, []byte("shared key123456"))
	assert.NoError(err)
	assert.Equal(originalToken, token)
}

func createEncryptedTestToken() (string, *Token, error) {
	scope := Scope{[]string{"role1", "role2"}, []string{"cap1", "cap2"}}

	var token Token
	token.IssueTime = time.Now().Unix()
	token.ID = "Id"
	token.AccountID = "AccountID"
	token.Type = "access_token"
	token.Scope = scope

	var sharedKey = []byte("shared key123456")

	str, err := EncryptAccessToken(&token, sharedKey)
	return str, &token, err
}
