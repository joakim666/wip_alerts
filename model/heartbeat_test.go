package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewHeartbeat(t *testing.T) {
	assert := assert.New(t)

	hb := NewHeartbeat()
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

		hbs := make(map[string]Heartbeat)

		// add a heartbeat
		h1 := NewHeartbeat()
		h1.APIKeyID = "APIKeyID1"

		hbs[h1.ID] = *h1

		err = SaveHeartbeats(db, "foo", &hbs)
		assert.NoError(err)

		// should be one Heartbeat
		heartbeats, err = ListHeartbeats(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 1)
		assert.Equal(*h1, (*heartbeats)[h1.ID])

		//		actualH1 := (*heartbeats)[h1.ID].(Heartbeat)
		//		assert.Equal("APIKeyID1", actualH1.APIKeyID)

		// add another Heartbeat
		h2 := NewHeartbeat()
		h2.APIKeyID = "APIKeyID2"
		hbs[h2.ID] = *h2

		err = SaveHeartbeats(db, "foo", &hbs)
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
