package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	assert := assert.New(t)

	d := NewDevice()
	assert.NotNil(d)
	assert.NotNil(d.ID)
	assert.NotNil(d.CreatedAt)
}

func TestDevices(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		devices, err := ListDevices(db, "foo")
		assert.NoError(err)
		assert.NotNil(devices)
		assert.True(len(*devices) == 0)

		dd := make(map[string]Device)

		// add a devices
		d1 := NewDevice()
		dd[d1.ID] = *d1

		err = SaveDevices(db, "foo", &dd)
		assert.NoError(err)

		// should be one device
		devices, err = ListDevices(db, "foo")
		assert.NoError(err)
		assert.NotNil(devices)
		assert.True(len(*devices) == 1)
		assert.Equal(*d1, (*devices)[d1.ID])

		//add another device
		d2 := NewDevice()
		dd[d2.ID] = *d2

		err = SaveDevices(db, "foo", &dd)
		assert.NoError(err)

		// should be two device
		devices, err = ListDevices(db, "foo")
		assert.NoError(err)
		assert.NotNil(devices)
		assert.True(len(*devices) == 2)
		assert.Equal(*d1, (*devices)[d1.ID])
		assert.Equal(*d2, (*devices)[d2.ID])

	})
}
