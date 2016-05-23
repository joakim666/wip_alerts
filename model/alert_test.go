package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewAlert(t *testing.T) {
	assert := assert.New(t)

	a := NewAlert("apiid")
	assert.NotNil(a)
	assert.NotNil(a.ID)
	assert.NotNil(a.APIKeyID)
	assert.NotNil(a.CreatedAt)
}

func TestAlerts(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		alerts, err := ListAlerts(db, "foo")
		assert.NoError(err)
		assert.NotNil(alerts)
		assert.True(len(*alerts) == 0)

		// add a alert
		a1 := NewAlert("APIKeyID1")
		err = a1.Save(db, "foo")
		assert.NoError(err)

		// should be one alert
		alerts, err = ListAlerts(db, "foo")
		assert.NoError(err)
		assert.NotNil(alerts)
		assert.True(len(*alerts) == 1)
		assert.Equal(*a1, (*alerts)[a1.ID])

		// add another alert
		a2 := NewAlert("APIKeyID2")
		err = a2.Save(db, "foo")
		assert.NoError(err)

		// should be two alerts
		alerts, err = ListAlerts(db, "foo")
		assert.NoError(err)
		assert.NotNil(alerts)
		assert.True(len(*alerts) == 2)
		assert.Equal(*a1, (*alerts)[a1.ID])
		assert.Equal(*a2, (*alerts)[a2.ID])
	})
}
