package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/joakim666/wip_alerts/auth"
	"github.com/joakim666/wip_alerts/model"
	"github.com/stretchr/testify/assert"
)

func TestListTokensWithNoTokens(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		a := func(token *auth.Token, ctx *gin.Context) bool {
			return true
		}
		authFunc := auth.ValidateAccessToken(a, []byte("shared key123456"))

		router.Use(authFunc)

		router.GET("/tokens", ListTokens(db))

		req, _ := http.NewRequest("GET", "/tokens", nil)
		res := httptest.NewRecorder()

		serializedToken, _, err := createEncryptedAdminTestToken()
		assert.NoError(err)

		req.Header.Add("Authorization", "Bearer "+serializedToken)
		router.ServeHTTP(res, req)

		assert.Equal(200, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		dtos := make(map[string]TokenDTO)
		err = json.Unmarshal(body, &dtos)
		assert.NoError(err)

		assert.Equal(0, len(dtos))
	})
}

func TestListTokensWithOneToken(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		scope := model.Scope{Roles: []string{"test"}, Capabilities: []string{}}

		t1 := model.NewToken()
		t1.IssueTime = time.Now()
		t1.Type = "refresh_token"
		t1.Scope = scope
		t1.RawString = "foobar"

		tokens := make(map[string]model.Token)
		tokens[t1.ID] = *t1

		err := model.SaveTokens(db, "foo", &tokens)
		assert.NoError(err)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		a := func(token *auth.Token, ctx *gin.Context) bool {
			return true
		}
		authFunc := auth.ValidateAccessToken(a, []byte("shared key123456"))

		router.Use(authFunc)

		router.GET("/tokens", ListTokens(db))

		req, _ := http.NewRequest("GET", "/tokens", nil)
		res := httptest.NewRecorder()

		serializedToken, _, err := createEncryptedAdminTestToken()
		assert.NoError(err)

		req.Header.Add("Authorization", "Bearer "+serializedToken)
		router.ServeHTTP(res, req)

		assert.Equal(200, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		dtos := make(map[string]TokenDTO)
		err = json.Unmarshal(body, &dtos)
		assert.NoError(err)

		assert.Equal(1, len(dtos))

		assert.Equal(t1.ID, dtos[t1.ID].ID)
		assert.Equal("refresh_token", dtos[t1.ID].Type)
		assert.Equal("test", dtos[t1.ID].Scope.Roles[0])
		assert.Equal("foobar", dtos[t1.ID].RawString)
	})
}

func TestListTokensWithTwoTokensSameAccount(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		scope := model.Scope{Roles: []string{"test"}, Capabilities: []string{}}

		t1 := model.NewToken()
		t1.IssueTime = time.Now()
		t1.Type = "refresh_token"
		t1.Scope = scope
		t1.RawString = "foobar"

		t2 := model.NewToken()
		t2.IssueTime = time.Now()
		t2.Type = "access_token"
		t2.Scope = scope
		t2.RawString = "foobar2"

		tokens := make(map[string]model.Token)
		tokens[t1.ID] = *t1
		tokens[t2.ID] = *t2

		err := model.SaveTokens(db, "foo", &tokens)
		assert.NoError(err)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		a := func(token *auth.Token, ctx *gin.Context) bool {
			return true
		}
		authFunc := auth.ValidateAccessToken(a, []byte("shared key123456"))

		router.Use(authFunc)

		router.GET("/tokens", ListTokens(db))

		req, _ := http.NewRequest("GET", "/tokens", nil)
		res := httptest.NewRecorder()

		serializedToken, _, err := createEncryptedAdminTestToken()
		assert.NoError(err)

		req.Header.Add("Authorization", "Bearer "+serializedToken)
		router.ServeHTTP(res, req)

		assert.Equal(200, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		dtos := make(map[string]TokenDTO)
		err = json.Unmarshal(body, &dtos)
		assert.NoError(err)

		assert.Equal(2, len(dtos))

		assert.Equal(t1.ID, dtos[t1.ID].ID)
		assert.Equal("refresh_token", dtos[t1.ID].Type)
		assert.Equal("test", dtos[t1.ID].Scope.Roles[0])
		assert.Equal("foobar", dtos[t1.ID].RawString)

		assert.Equal(t2.ID, dtos[t2.ID].ID)
		assert.Equal("access_token", dtos[t2.ID].Type)
		assert.Equal("test", dtos[t2.ID].Scope.Roles[0])
		assert.Equal("foobar2", dtos[t2.ID].RawString)
	})
}

func TestListTokensWithThreeTokensAndTwoDifferentAccount(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		scope := model.Scope{Roles: []string{"test"}, Capabilities: []string{}}

		t1 := model.NewToken()
		t1.IssueTime = time.Now()
		t1.Type = "refresh_token"
		t1.Scope = scope
		t1.RawString = "foobar"

		t2 := model.NewToken()
		t2.IssueTime = time.Now()
		t2.Type = "access_token"
		t2.Scope = scope
		t2.RawString = "foobar2"

		// Add two tokens for account 'foo'
		tokens1 := make(map[string]model.Token)
		tokens1[t1.ID] = *t1
		tokens1[t2.ID] = *t2

		err := model.SaveTokens(db, "foo", &tokens1)
		assert.NoError(err)

		t3 := model.NewToken()
		t3.IssueTime = time.Now()
		t3.Type = "access_token"
		t3.Scope = scope
		t3.RawString = "foobar3"

		// Add one token for account 'bar'
		tokens2 := make(map[string]model.Token)
		tokens2[t3.ID] = *t3

		err = model.SaveTokens(db, "bar", &tokens2)
		assert.NoError(err)

		gin.SetMode(gin.TestMode)
		router := gin.New()

		a := func(token *auth.Token, ctx *gin.Context) bool {
			return true
		}
		authFunc := auth.ValidateAccessToken(a, []byte("shared key123456"))

		router.Use(authFunc)

		router.GET("/tokens", ListTokens(db))

		req, _ := http.NewRequest("GET", "/tokens", nil)
		res := httptest.NewRecorder()

		serializedToken, _, err := createEncryptedAdminTestToken()
		assert.NoError(err)

		req.Header.Add("Authorization", "Bearer "+serializedToken)
		router.ServeHTTP(res, req)

		assert.Equal(200, res.Code)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(err)

		dtos := make(map[string]TokenDTO)
		err = json.Unmarshal(body, &dtos)
		assert.NoError(err)

		assert.Equal(3, len(dtos), "The returned map should contain all three tokens")

		assert.Equal(t1.ID, dtos[t1.ID].ID)
		assert.Equal("refresh_token", dtos[t1.ID].Type)
		assert.Equal("test", dtos[t1.ID].Scope.Roles[0])
		assert.Equal("foobar", dtos[t1.ID].RawString)

		assert.Equal(t2.ID, dtos[t2.ID].ID)
		assert.Equal("access_token", dtos[t2.ID].Type)
		assert.Equal("test", dtos[t2.ID].Scope.Roles[0])
		assert.Equal("foobar2", dtos[t2.ID].RawString)

		assert.Equal(t3.ID, dtos[t3.ID].ID)
		assert.Equal("access_token", dtos[t3.ID].Type)
		assert.Equal("test", dtos[t3.ID].Scope.Roles[0])
		assert.Equal("foobar3", dtos[t3.ID].RawString)
	})
}

func createEncryptedAdminTestToken() (string, *auth.Token, error) {
	scope := auth.Scope{Roles: []string{"admin"}, Capabilities: []string{}}

	var token auth.Token
	token.IssueTime = time.Now().Unix()
	token.ID = "Id"
	token.AccountID = "AccountID"
	token.Type = "access_token"
	token.Scope = scope

	var sharedKey = []byte("shared key123456")

	str, err := auth.EncryptAccessToken(&token, sharedKey)
	return str, &token, err
}
