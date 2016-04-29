package main

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/auth"
	"github.com/joakim666/wip_alerts/model"
)

type NewTokenDTO struct {
	GrantType string  `json:"grant_type" binding:"required"`
	AccountID *string `json:"account_id"`
	RenewalID *string `json:"renewal_id"`
}

type ScopeDTO struct {
	Roles        []string `json:"roles"`
	Capabilities []string `json:"capabilities"`
}

type TokenDTO struct {
	ID        string    `json:"id"` // uuid
	AccountID string    `json:"account_id"`
	IssueTime time.Time `json:"issue_time"`
	Type      string    `json:"type"`
	Scope     ScopeDTO  `json:"scope"`
	RawString string    `json:"raw_string"`
	CreatedAt time.Time `json:"created_at"`
}

func ListTokens(db *bolt.DB) gin.HandlerFunc {
	glog.Infof("ListTokens")

	tokenDTOs := make(map[string]TokenDTO) // key = tokenID

	return func(c *gin.Context) {
		tokens, err := model.ListAllTokens(db)
		if err != nil {
			glog.Errorf("Failed to get tokens: %s", err)
		} else {
			dtos := makeTokenDTOs(db, tokens)
			for _, v := range *dtos {
				tokenDTOs[v.ID] = v
			}
		}

		c.JSON(200, tokenDTOs)
	}
}

// PostTokens creates a new token. 'publicKey' is the public part of the private-key used to sign and encrypt the refresh tokens. 'encryptionKey' is the shared key used to sign, encrypt, validate and decrypt the access tokens.
func PostTokens(db *bolt.DB, publicKey interface{}, encryptionKey interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("PostTokens")
		var json NewTokenDTO

		err := c.BindJSON(&json)
		if err != nil {
			glog.Infof("Binding failed: %s", err)
			c.Status(400) // => Bad Request
			return
		}

		glog.Infof("Json: %s", json)

		switch json.GrantType {
		case "account":
			handleAccountRequest(c, &json, db, publicKey, encryptionKey)
			return
		case "renewal":
			handleRenewalRequest(c, &json, db, encryptionKey)
			return
		default:
			// bad request
			c.Status(400)
		}
	}
}

func handleAccountRequest(c *gin.Context, json *NewTokenDTO, db *bolt.DB, publicKey interface{}, encryptionKey interface{}) {
	glog.Infof("handleAccountRequest")
	// AccountID is mandatory
	if json.AccountID == nil {
		c.Status(400)
		return
	}

	// look up account
	account, err := model.GetAccount(db, *json.AccountID)
	if err != nil {
		glog.Errorf("Failed to find matching account for id=%s: %s", *json.AccountID, err)
		c.Status(400) // => Bad Request
		return
	}

	// get already created tokens for this account
	oldTokens, err := model.ListTokens(db, account.ID)
	if err != nil {
		glog.Errorf("Failed to find existing tokens for account with id=%s: %s", account.ID, err)
		c.Status(400) // => Bad Request
		return
	}

	// check if any of the already created tokens is a refresh token
	for _, v := range *oldTokens {
		if v.Type == "refresh_token" {
			// this account id already has a created refresh token
			glog.Errorf("Account %s already has a refresh token", account.ID)
			c.Status(400) // => Bad Request
			return
		}
	}

	if account == nil {
		glog.Errorf("Failed to find matching account for id=%s: %s", *json.AccountID, err)
		c.Status(400) // => Bad Request
		return
	}

	now := time.Now()

	// begin - create refresh token
	refreshTokenStr, err := createRefreshToken(now, *json.AccountID, db, publicKey)
	if err != nil {
		glog.Errorf("Failed to create refresh token: %s", err)
		c.Status(500)
		return
	}
	// end - create refresh token

	// begin - create access token
	accessTokenStr, err := createAccessToken(now, *json.AccountID, db, encryptionKey)
	if err != nil {
		glog.Errorf("Failed to create access token: %s", err)
		c.Status(500)
		return
	}
	// end - create access token

	c.JSON(201, gin.H{
		"refresh_token": refreshTokenStr,
		"access_token":  accessTokenStr,
	})
}

func handleRenewalRequest(c *gin.Context, json *NewTokenDTO, db *bolt.DB, encryptionKey interface{}) {
	// RenewalID is mandatory
	if json.RenewalID == nil {
		c.Status(400)
		return
	}

	// look up account from RenewalID
	renewal, accountID, err := model.GetRenewal(db, *json.RenewalID)
	if err != nil {
		glog.Errorf("Failed to find matching renewal for id=%s: %s", *json.RenewalID, err)
		c.Status(400) // => Bad Request
		return
	}

	if renewal == nil {
		glog.Errorf("Failed to find matching renewal for id=%s: %s", *json.RenewalID, err)
		c.Status(400) // => Bad Request
		return
	}

	if renewal.UsedAt != nil {
		// this renewal has already been used
		glog.Errorf("Renewal %s has already been used", *json.RenewalID)
		c.Status(400) // => Bad Request
		return
	}

	now := time.Now()

	// begin - create access token
	accessTokenStr, err := createAccessToken(now, *accountID, db, encryptionKey)
	if err != nil {
		glog.Errorf("Failed to create access token: %s", err)
		c.Status(500)
		return
	}
	// end - create access token

	// Mark renewal as used
	renewal.UsedAt = &now
	err = renewal.Save(db, *accountID)
	if err != nil {
		glog.Errorf("Failed to saved renewal: %s", err)
		c.Status(500)
		return
	}

	c.JSON(201, gin.H{
		"access_token": accessTokenStr,
	})

}

func createRefreshToken(creationTime time.Time, accountID string, db *bolt.DB, publicKey interface{}) (string, error) {
	glog.Infof("createRefreshToken")

	dbRefreshToken := model.NewToken()

	refreshToken := auth.Token{}
	refreshToken.IssueTime = creationTime.Unix()
	refreshToken.ID = dbRefreshToken.ID
	refreshToken.AccountID = accountID
	refreshToken.Type = "refresh_token" // TODO enum
	refreshToken.Scope = auth.Scope{
		Roles:        []string{"user"},
		Capabilities: []string{"refresh_token"}}

	dbRefreshToken.IssueTime = creationTime
	dbRefreshToken.Type = refreshToken.Type
	dbRefreshToken.Scope = model.Scope{
		Roles:        refreshToken.Scope.Roles,
		Capabilities: refreshToken.Scope.Capabilities,
	}

	res, err := auth.EncryptRefreshToken(&refreshToken, publicKey)
	if err != nil {
		return "", err
	}

	dbRefreshToken.RawString = res

	err = dbRefreshToken.Save(db, accountID)
	if err != nil {
		return "", err
	}

	return res, nil
}

func createAccessToken(creationTime time.Time, accountID string, db *bolt.DB, encryptionKey interface{}) (string, error) {
	dbAccessToken := model.NewToken()

	accessToken := auth.Token{}
	accessToken.IssueTime = creationTime.Unix()
	accessToken.ID = dbAccessToken.ID
	accessToken.AccountID = accountID
	accessToken.Type = "access_token" // TODO enum
	accessToken.Scope = auth.Scope{
		Roles:        []string{"user"},
		Capabilities: []string{"access_token"}}

	dbAccessToken.IssueTime = creationTime
	dbAccessToken.Type = accessToken.Type
	dbAccessToken.Scope = model.Scope{
		Roles:        accessToken.Scope.Roles,
		Capabilities: accessToken.Scope.Capabilities,
	}

	res, err := auth.EncryptAccessToken(&accessToken, encryptionKey)
	if err != nil {
		return "", err
	}

	dbAccessToken.RawString = res

	err = dbAccessToken.Save(db, accountID)
	if err != nil {
		return "", err
	}

	return res, nil
}

func makeTokenDTOs(db *bolt.DB, tokens *map[string][]model.Token) *[]TokenDTO {
	glog.Infof("makeTokenDTOs. Size=%d", len(*tokens))
	dtos := make([]TokenDTO, 0)

	if len(*tokens) > 0 {
		for k, v := range *tokens {
			for _, v := range v {
				dtos = append(dtos, makeTokenDTO(k, v))
			}
		}
	} else {
	}

	return &dtos
}

func makeTokenDTO(accountId string, token model.Token) TokenDTO { // TODO pointers?
	glog.Infof("makeTokenDTO")
	var dto TokenDTO

	dto.ID = token.ID
	dto.AccountID = accountId
	dto.IssueTime = token.IssueTime
	dto.Type = token.Type
	dto.Scope = makeScopeDTO(token.Scope)
	dto.RawString = token.RawString
	dto.CreatedAt = token.CreatedAt
	return dto
}

func makeScopeDTO(scope model.Scope) ScopeDTO {
	glog.Infof("makeScopeDTO")
	var dto ScopeDTO

	dto.Roles = scope.Roles
	dto.Capabilities = scope.Capabilities

	return dto
}
