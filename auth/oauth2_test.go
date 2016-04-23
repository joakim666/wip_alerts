package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExtractToken(t *testing.T) {
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "/foo", nil)

	token, err := extractToken(req)
	assert.Error(err)

	serializedToken, _, err := createEncryptedTestToken()
	assert.NoError(err)

	fmt.Printf("Test token: %s\n", serializedToken)

	req.Header.Add("Authorization", "Bearer "+serializedToken)
	token, err = extractToken(req)
	assert.NoError(err)
	assert.NotNil(token)
	assert.Equal(serializedToken, token)
}

func TestValidateAccessTokenNoHeader(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	a := func(token *Token, ctx *gin.Context) bool {
		return true
	}
	authFunc := ValidateAccessToken(a, []byte("shared key123456"))

	router.Use(authFunc)

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
		c.String(200, "OK")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)

	assert.False(called)
	assert.Equal(401, w.Code)
}

func TestValidateAccessTokenHeader(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	a := func(token *Token, ctx *gin.Context) bool {
		return true
	}
	authFunc := ValidateAccessToken(a, []byte("shared key123456"))

	router.Use(authFunc)

	called := false

	router.GET("/test", func(c *gin.Context) {
		called = true
		c.String(200, "OK")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)

	serializedToken, _, err := createEncryptedTestToken()
	assert.NoError(err)

	r.Header.Add("Authorization", "Bearer "+serializedToken)
	router.ServeHTTP(w, r)

	assert.True(called)
	assert.Equal(200, w.Code)
}
