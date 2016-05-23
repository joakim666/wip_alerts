package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"time"
	"github.com/joakim666/wip_alerts/model"
)

type createAlertDTO struct {
	Title            string                `json:"title" binding:"required"`
	ShortDescription string                `json:"short_description" binding:"required"`
	LongDescription  string                `json:"long_description" binding:"required"`
	Priority         model.AlertPriority   `json:"priority" binding:"required"`
	TriggeredAt      time.Time             `json:"triggered_at" binding:"required"`
}

type alertDTO struct {
	ID               string                `json:"id"`
	Title            string                `json:"title"`
	ShortDescription string                `json:"short_description"`
	LongDescription  string                `json:"long_description"`
	Priority         model.AlertPriority   `json:"priority"`
	TriggeredAt      time.Time             `json:"triggered_at"`
	CreatedAt        time.Time             `json:"created_at"`
}

func CreateAlertRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("CreateAlertRoute")

		accountIDInterface, exists := c.Get("accountID")
		if exists == false {
			glog.Infof("No accountID set")
			c.Status(401) // => Unauthorized
			return
		}
		apiKeyID, exists := c.Get("apiKeyID")
		if exists == false {
			glog.Infof("No apiKeyID set")
			c.Status(401) // => Unauthorized
			return
		}

		accountID, ok := accountIDInterface.(string)
		if ok == false {
			glog.Infof("AccountID in context is not a string")
			c.Status(401) // => Unauthorized
			return
		}

		var json createAlertDTO

		err := c.BindJSON(&json)
		if err != nil {
			glog.Infof("Binding failed: %s", err)
			c.Status(400) // => Bad Request
			return
		}

		glog.Infof("Json: %s", json)

		alert := model.NewAlert(apiKeyID.(string))
		alert.Title = json.Title
		alert.ShortDescription = json.ShortDescription
		alert.LongDescription = json.LongDescription
		alert.Priority = json.Priority
		alert.TriggeredAt = json.TriggeredAt

		err = alert.Save(db, accountID)
		if err != nil {
			glog.Errorf("Failed to save created alert: %s", err)
			c.Status(500) // => Internal Server error
			return
		}

		dto := makeAlertDTO(alert)

		c.JSON(201, dto)
	}
}

func makeAlertDTO(alert *model.Alert) alertDTO {
	var dto alertDTO

	dto.ID = alert.ID
	dto.Title = alert.Title
	dto.ShortDescription = alert.ShortDescription
	dto.LongDescription = alert.LongDescription
	dto.Priority = alert.Priority
	dto.TriggeredAt = alert.TriggeredAt
	dto.CreatedAt = alert.CreatedAt

	return dto
}
