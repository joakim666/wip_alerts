package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"time"
	"github.com/joakim666/wip_alerts/model"
	"net/http"
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
	Status           model.AlertStatus     `json:"status"`
	TriggeredAt      time.Time             `json:"triggered_at"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type updateAlertDTO struct {
	Status model.AlertStatus        `json:"status" binding:"required"`
}

// CreateAlertRoute creates and saves a new alert
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

		c.JSON(http.StatusCreated, dto)
	}
}

// ListAlertsRoute lists all alerts with status not "archived"
func ListAlertsRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("ListAlertsRoute")

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

		glog.Infof("Listing alerts for account id: %s", accountID)

		alerts, err := model.ListNonArchivedAlerts(db, accountID)
		if err != nil {
			glog.Infof("No account for account id=%s", accountID)
			c.Status(400) // => Bad Request
			return
		}

		dtos := make(map[string]alertDTO, 0)

		for _, v := range *alerts {
			dto := makeAlertDTO(&v)
			dtos[dto.ID] = dto
		}

		c.JSON(http.StatusOK, dtos)
	}
}

// UpdateAlertRoute updates the status of the alert. Allowed status transitions are:
//  * NewStatus -> SeenStatus
//  * NewStatus -> ArchivedStatus
//  * SeenStatus -> ArchivedStatus
func UpdateAlertRoute(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		glog.Infof("UpdateAlertRoute")

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

		alertID := c.Param("id")

		glog.Infof("Update alert with id %s for account id: %s", alertID, accountID)

		alert, accId, err := model.GetAlert(db, alertID)
		if err != nil {
			glog.Errorf("Error search for alert with id %s", alertID, err)
			c.Status(http.StatusNotFound)
			return
		}
		if alert == nil {
			glog.Errorf("Could not find alert with id %s", alertID)
			c.Status(http.StatusNotFound)
			return
		}

		if accountID != *accId {
			glog.Errorf("Authorized with account id %s but trying to access alert %s beloging to account %s", accountID, alertID, accId)
			c.Status(http.StatusUnauthorized) // => Unauthorized
			return
		}

		var json updateAlertDTO

		err = c.BindJSON(&json)
		if err != nil {
			glog.Infof("Binding failed: %s", err)
			c.Status(400) // => Bad Request
			return
		}

		glog.Infof("Json: %s", json)

		if alert.Status == model.NewStatus && json.Status == model.SeenStatus {
			alert.Status = model.SeenStatus
		} else if alert.Status == model.NewStatus && json.Status == model.ArchivedStatus {
			alert.Status = model.ArchivedStatus
		} else if alert.Status == model.SeenStatus && json.Status == model.ArchivedStatus {
			alert.Status = model.ArchivedStatus
		} else {
			glog.Infof("Trying to do state transition from %s to %s", alert.Status, json.Status)
			c.Status(http.StatusBadRequest) // => BadRequest
			return
		}

		alert.UpdatedAt = time.Now(); // update the updated time

		err = alert.Save(db, accountID)
		if err != nil {
			glog.Errorf("Failed to save updated alert: %s", err)
			c.Status(500) // => Internal Server error
			return
		}

		dto := makeAlertDTO(alert)

		c.JSON(http.StatusOK, dto)
	}
}

func makeAlertDTO(alert *model.Alert) alertDTO {
	var dto alertDTO

	dto.ID = alert.ID
	dto.Title = alert.Title
	dto.ShortDescription = alert.ShortDescription
	dto.LongDescription = alert.LongDescription
	dto.Priority = alert.Priority
	dto.Status = alert.Status
	dto.TriggeredAt = alert.TriggeredAt
	dto.CreatedAt = alert.CreatedAt
	dto.UpdatedAt = alert.UpdatedAt

	return dto
}
