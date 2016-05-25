package model

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"time"
	"flag"
)

func TestNewHeartbeat(t *testing.T) {
	assert := assert.New(t)

	hb := NewHeartbeat("ApiKeyID555")
	assert.NotNil(hb)
	assert.NotNil(hb.ID)
	assert.NotNil(hb.CreatedAt)
}

func TestListHeartbeats(t *testing.T) {
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

func TestLatestHeartbeats(t *testing.T) {
	flag.Lookup("logtostderr").Value.Set("true")

	RunInTestDb(t, func(t *testing.T, db *bolt.DB) {
		assert := assert.New(t)

		// should be empty
		heartbeats, err := LatestHeartbeatPerApiKey(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 0)

		// add a heartbeat
		h1 := NewHeartbeat("APIKeyID1")
		h1.ExecutedAt = time.Now().Add(-10 * time.Minute)
		err = h1.Save(db, "foo")
		assert.NoError(err)

		// should be one Heartbeat
		heartbeats, err = LatestHeartbeatPerApiKey(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 1)
		assert.Equal(*h1, (*heartbeats)[h1.ID])

		// add another heartbeat with the same api key
		h2 := NewHeartbeat("APIKeyID1")
		h2.ExecutedAt = time.Now()
		err = h2.Save(db, "foo")
		assert.NoError(err)

		// should be still be just one heartbeats, but should be the latest one i.e. h2
		heartbeats, err = LatestHeartbeatPerApiKey(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 1)
		assert.Equal(*h2, (*heartbeats)[h2.ID])

		// add another heartbeat with the another api key
		h3 := NewHeartbeat("APIKeyID2")
		h3.ExecutedAt = time.Now()
		err = h3.Save(db, "foo")
		assert.NoError(err)

		// should be two heartbeats
		heartbeats, err = LatestHeartbeatPerApiKey(db, "foo")
		assert.NoError(err)
		assert.NotNil(heartbeats)
		assert.True(len(*heartbeats) == 2)
		assert.Equal(*h2, (*heartbeats)[h2.ID])
		assert.Equal(*h3, (*heartbeats)[h3.ID])
	})
}
