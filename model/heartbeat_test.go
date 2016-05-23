package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewHeartbeat(t *testing.T) {
	assert := assert.New(t)

	hb := NewHeartbeat("ApiKeyID555")
	assert.NotNil(hb)
	assert.NotNil(hb.ID)
	assert.NotNil(hb.CreatedAt)
}

func TestHeartbeats(t *testing.T) {
	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		heartbeats, err := ListHeartbeats(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 0)

		// add a heartbeat
		h1 := NewHeartbeat("APIKeyID1")
		err = h1.Save(db, "foo")
		assert.NoError(err)

		// should be one Heartbeat
		heartbeats, err = ListHeartbeats(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 1)
		assert.Equal(*h1, (*heartbeats)[h1.ID])

		// add another Heartbeat
		h2 := NewHeartbeat("APIKeyID2")
		err = h2.Save(db, "foo")
		assert.NoError(err)

		// should be two heartbeats
		heartbeats, err = ListHeartbeats(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 2)
		assert.Equal(*h1, (*heartbeats)[h1.ID])
		assert.Equal(*h2, (*heartbeats)[h2.ID])
	})
}
