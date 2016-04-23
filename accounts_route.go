package main

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/joakim666/wip_alerts/model"
)

type NewAccountDTO struct {
	DeviceID   string `json:"device_id" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"`
	DeviceInfo string `json:"device_info" binding:"required"` // json as a string TODO validate that it's proper json
}

type AccountDTO struct {
	ID      string      `json:"id"` // uuid
	Devices []DeviceDTO `json:"devices"`
}

type DeviceDTO struct {
	ID         string `json:"id"` // uuid
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type"`
	DeviceInfo string `json:"device_info"`
	// CreatedAt
}

func postAccounts(db *bolt.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json NewAccountDTO

		if c.BindJSON(&json) == nil {

			account := model.NewAccount()

			err := account.Save(db)
			if err != nil {
				glog.Errorf("Failed to save account in db: %s", err)
				c.Status(500)
				return
			}

			device := newDeviceFromDTO(&json)

			devices := make(map[string]model.Device)
			devices[device.ID] = *device

			err = model.SaveDevices(db, account.ID, &devices)
			if err != nil {
				glog.Errorf("Failed to save devices for account in db: %s", err)
				c.Status(500)
				return
			}

			c.Status(201) // 201 Created
		}
	}
}

func newDeviceFromDTO(dto *NewAccountDTO) *model.Device {
	device := model.NewDevice()
	device.DeviceID = dto.DeviceID
	device.DeviceType = dto.DeviceType
	device.DeviceInfo = dto.DeviceInfo

	return device
}

func listAccounts(db *bolt.DB) gin.HandlerFunc {
	glog.Infof("listAccounts")

	return func(c *gin.Context) {
		accounts, err := model.ListAccounts(db)
		if err != nil {
			glog.Errorf("ListAccounts failed: %s", err)
			c.Status(500)
		} else {
			accountDTOs, err := makeAccountDTOs(db, accounts)
			if err != nil {
				glog.Errorf("Failed to transform accounts into dtos: %s", err)
				c.Status(500)
			}

			c.JSON(200, accountDTOs)
		}
	}
}

func makeAccountDTOs(db *bolt.DB, accounts *map[string]model.Account) (*[]AccountDTO, error) {
	glog.Infof("makeAccountDTOs. Size=%d", len(*accounts))
	var dtos []AccountDTO

	if len(*accounts) > 0 {
		for _, v := range *accounts {
			devices, err := model.ListDevices(db, v.ID)
			if err != nil {
				return nil, err
			}
			dtos = append(dtos, makeAccountDTO(v, devices))
		}
	} else {
		dtos = make([]AccountDTO, 0)
	}

	return &dtos, nil
}

func makeAccountDTO(account model.Account, devices *map[string]model.Device) AccountDTO { // TODO pointers?
	glog.Infof("makeAccountDTO")
	var dto AccountDTO

	dto.ID = account.ID
	dto.Devices = *makeDeviceDTOs(devices)

	return dto
}

func makeDeviceDTOs(devices *map[string]model.Device) *[]DeviceDTO {
	glog.Infof("makeDeviceDTOs. Size=%d", len(*devices))
	var dtos []DeviceDTO

	for _, v := range *devices {
		dtos = append(dtos, makeDeviceDTO(v))
	}

	return &dtos
}

func makeDeviceDTO(device model.Device) DeviceDTO { // TODO pointers?
	glog.Infof("makeDeviceDTO")
	var dto DeviceDTO

	dto.ID = device.ID
	dto.DeviceID = device.DeviceID
	dto.DeviceType = device.DeviceType
	dto.DeviceInfo = device.DeviceInfo

	return dto
}
