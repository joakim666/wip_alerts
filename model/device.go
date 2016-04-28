package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

type Device struct {
	ID         string // uuid
	DeviceID   string
	DeviceType string
	DeviceInfo string
	CreatedAt  time.Time
}

func (d Device) PersistanceID() string {
	return d.ID
}

func NewDevice() *Device {
	var d Device
	uuid := uuid.NewV4()
	d.ID = uuid.String()
	d.CreatedAt = time.Now()
	return &d
}

func SaveDevices(db *bolt.DB, accountUUID string, devices *map[string]Device) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "Devices", BoltMap(devices))
}

func ListDevices(db *bolt.DB, accountUUID string) (*map[string]Device, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Devices", reflect.TypeOf(Device{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Device
	m2 := make(map[string]Device)
	for _, v := range *m {
		d := v.(*Device)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}
