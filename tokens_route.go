package main

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/model"
)

type NewTokenDTO struct {
	GrantType string `json:"grant_Type" binding:"required"`
	AccountID string `json:"account_id"`
	RenewalID string `json:"renewal_id"`
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

/*func PostTokens(db *bolt.DB, privateKey interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json NewRenewalDTO

		if c.BindJSON(&json) == nil {
			token, err := auth.DecryptRefreshToken(json.RefreshToken, privateKey)
			if err != nil {
				glog.Errorf("Failed to decrypt refresh token: %s", err)
				c.Status(500)
				return
			}

			renewal := model.NewRenewal()
			renewal.RefreshTokenID = token.ID

			// get existing renewals
			renewals, err := model.ListRenewals(db, token.AccountID)
			if err != nil {
				glog.Errorf("Failed to get renewals: %s", err)
				c.Status(500)
				return
			}

			// add new renewal
			(*renewals)[renewal.ID] = *renewal

			// save the renewals
			err = model.SaveRenewals(db, token.ID, renewals)
			if err != nil {
				glog.Errorf("Failed to save renewal in db: %s", err)
				c.Status(500)
				return
			}

			c.JSON(201, gin.H{
				"renewal_id": renewal.ID,
			})
		}
	}
}*/

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
