package main

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/auth"
	"github.com/joakim666/wip_alerts/model"
)

type NewRenewalDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceType   string `json:"device_type" binding:"required"`
	DeviceInfo   string `json:"device_info" binding:"required"` // json as a string TODO validate that it's proper json
}

type RenewalDTO struct {
	ID             string    `json:"id"`               // uuid
	AccountID      string    `json:"account_id"`       // account uuid
	RefreshTokenID string    `json:"refresh_token_id"` // uuid of refresh token
	CreatedAt      time.Time `json:"created_at"`
}

func ListRenewals(db *bolt.DB) gin.HandlerFunc {
	glog.Infof("listRenewals")

	var renewalDTOs []RenewalDTO

	return func(c *gin.Context) {
		accounts, err := model.ListAccounts(db)
		if err != nil {
			glog.Errorf("ListRenewals failed: %s", err)
			c.Status(500)
		} else {
			for _, v := range *accounts {
				renewals, err := model.ListRenewals(db, v.ID)
				if err != nil {
					glog.Errorf("Failed to get renewals for account %s: %s", v.ID, err)
				} else {
					dtos := makeRenewalDTOs(db, v.ID, renewals)
					renewalDTOs = append(renewalDTOs, *dtos...)
				}
			}

			c.JSON(200, renewalDTOs)
		}
	}
}

func PostRenewals(db *bolt.DB, privateKey interface{}) gin.HandlerFunc {
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
}

func makeRenewalDTOs(db *bolt.DB, accountId string, renewals *map[string]model.Renewal) *[]RenewalDTO {
	glog.Infof("makeRenewalDTOs. Size=%d", len(*renewals))
	dtos := make([]RenewalDTO, 0)

	if len(*renewals) > 0 {
		for _, v := range *renewals {
			dtos = append(dtos, makeRenewalDTO(accountId, v))
		}
	} else {
	}

	return &dtos
}

func makeRenewalDTO(accountId string, renewal model.Renewal) RenewalDTO { // TODO pointers?
	glog.Infof("makeRenewalDTO")
	var dto RenewalDTO

	dto.ID = renewal.ID
	dto.AccountID = accountId
	dto.RefreshTokenID = "TODO"
	dto.CreatedAt = renewal.CreatedAt

	return dto
}
