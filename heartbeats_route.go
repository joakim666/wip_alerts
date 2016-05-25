package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"time"
	"github.com/joakim666/wip_alerts/model"
	"net/http"
)

type createHeartbeatDTO struct {
	ExecutedAt	time.Time		`json:"executed_at" binding:"required"`
}

type heartbeatDTO struct {
	ID               string                `json:"id"`
	ExecutedAt       time.Time             `json:"executed_at"`
	CreatedAt        time.Time             `json:"created_at"`
}


// CreateHeartbeatRoute creates and saves a new heartbeat
func CreateHeartbeatRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("CreateHeartbeatRoute")

		accountIDInterface, exists := c.Get("accountID")
		if exists == false {
			glog.Infof("No accountID set")
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}
		apiKeyID, exists := c.Get("apiKeyID")
		if exists == false {
			glog.Infof("No apiKeyID set")
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}

		accountID, ok := accountIDInterface.(string)
		if ok == false {
			glog.Infof("AccountID in context is not a string")
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}

		var json createHeartbeatDTO

		err := c.BindJSON(&json)
		if err != nil {
			glog.Infof("Binding failed: %s", err)
			c.Status(http.StatusBadRequest) // => Bad Request
			return
		}

		glog.Infof("Json: %s", json)

		hb := model.NewHeartbeat(apiKeyID.(string))
		hb.ExecutedAt = json.ExecutedAt

		err = hb.Save(db, accountID)
		if err != nil {
			glog.Errorf("Failed to save created heartbeat: %s", err)
			c.Status(http.StatusInternalServerError) // => Internal Server error
			return
		}

		dto := makeHeartbeatDTO(hb)

		c.JSON(http.StatusCreated, dto)
	}
}

// LatestHeartbeats returns the last heartbeat for each api key for the identified account
func LatestHeartbeatsRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("LatestHeartbeats")

		accountIDInterface, exists := c.Get("accountID")
		if exists == false {
			glog.Infof("No accountID set")
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}

		accountID, ok := accountIDInterface.(string)
		if ok == false {
			glog.Infof("AccountID in context is not a string")
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}

		glog.Infof("Listing latest heartbeat for account id: %s", accountID)

		heartbeats, err := model.LatestHeartbeatPerApiKey(db, accountID)
		if err != nil {
			glog.Infof("No account for account id=%s", accountID)
			c.Status(http.StatusBadRequest) // => Bad Request
			return
		}

		dtos := make(map[string]heartbeatDTO, 0)

		for _, v := range *heartbeats {
			dto := makeHeartbeatDTO(&v)
			dtos[dto.ID] = dto
		}

		c.JSON(http.StatusOK, dtos)
	}
}

func makeHeartbeatDTO(hb *model.Heartbeat) heartbeatDTO {
	var dto heartbeatDTO

	dto.ID = hb.ID
	dto.ExecutedAt = hb.ExecutedAt
	dto.CreatedAt = hb.CreatedAt

	return dto
}
