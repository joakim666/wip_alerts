package main

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/model"
)

type createAPIKeyDTO struct {
	Description string `json:"description" binding:"required"`
}

type apiKeyDTO struct {
	ID          string             `json:"id"`
	Description string             `json:"description"`
	IssuedAt    time.Time          `json:"issued_at"`
	Status      model.APIKeyStatus `json:"status"`
}

func CreateAPIKeyRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("CreateAPIKeyRoute")

		accountIDInterface, exists := c.Get("accountID")
		if exists == false {
			glog.Infof("No accountID set")
			c.Status(401) // => Unauthorized
			return
		}

		accountID, ok := accountIDInterface.(string)
		if ok == false {
			glog.Infof("AccountID in context is not a string")
			c.Status(401) // => Unauthorized
			return
		}

		var json createAPIKeyDTO

		err := c.BindJSON(&json)
		if err != nil {
			glog.Infof("Binding failed: %s", err)
			c.Status(400) // => Bad Request
			return
		}

		glog.Infof("Json: %s", json)

		apiKey := model.NewAPIKey()
		apiKey.Description = json.Description

		err = apiKey.Save(db, accountID)
		if err != nil {
			glog.Errorf("Failed to save created API key: %s", err)
			c.Status(500) // => Internal Server error
			return
		}

		dto := makeDTO(apiKey)

		c.JSON(201, dto)
	}
}

func ListAPIKeyRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("ListAPIKeyRoute")

		accountIDInterface, exists := c.Get("accountID")
		if exists == false {
			glog.Infof("No accountID set")
			c.Status(401) // => Unauthorized
			return
		}

		accountID, ok := accountIDInterface.(string)
		if ok == false {
			glog.Infof("AccountID in context is not a string")
			c.Status(401) // => Unauthorized
			return
		}

		glog.Infof("Listing APIKeys for account id: %s", accountID)

		apiKeys, err := model.ListAPIKeys(db, accountID)
		if err != nil {
			glog.Infof("No account for account id=%s", accountID)
			c.Status(400) // => Bad Request
			return
		}

		dtos := make(map[string]apiKeyDTO, 0)

		for _, v := range *apiKeys {
			dto := makeDTO(&v)
			dtos[dto.ID] = dto
		}

		c.JSON(200, dtos)
	}
}

func makeDTO(apiKey *model.APIKey) apiKeyDTO {
	var dto apiKeyDTO

	dto.ID = apiKey.ID
	dto.Description = apiKey.Description
	dto.IssuedAt = apiKey.CreatedAt
	dto.Status = apiKey.Status

	return dto
}
