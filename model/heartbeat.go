package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
	"github.com/golang/glog"
)

type Heartbeat struct {
	ID         string // uuid
	APIKeyID   string // uuid of api key that sent this heartbeat
	ExecutedAt time.Time
	CreatedAt  time.Time
}

func (h Heartbeat) PersistanceID() string {
	return h.ID
}

func (h Heartbeat) Save(db *bolt.DB, accountUUID string) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "Heartbeats", BoltSingle(&h))
}

func NewHeartbeat(apiKeyID string) *Heartbeat {
	var hb Heartbeat
	uuid := uuid.NewV4()
	hb.ID = uuid.String()
	hb.APIKeyID = apiKeyID
	hb.CreatedAt = time.Now()
	return &hb
}

func ListHeartbeats(db *bolt.DB, accountUUID string) (*map[string]Heartbeat, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Heartbeats", reflect.TypeOf(Heartbeat{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Heartbeat
	m2 := make(map[string]Heartbeat)
	for _, v := range *m {
		hb := v.(*Heartbeat)
		m2[v.PersistanceID()] = *hb
	}

	return &m2, nil
}

func LatestHeartbeatPerApiKey(db *bolt.DB, accountUUID string) (*map[string]Heartbeat, error) {
	hbs, err := ListHeartbeats(db, accountUUID)
	if err != nil {
		return nil, err
	}

	apiKeyToHeartbeat := make(map[string]Heartbeat)

	// iterate over all heartbeats
	for _, v := range *hbs {
		glog.Info("FOO v: %s", v.ID)
		h, ok := apiKeyToHeartbeat[v.APIKeyID]
		if ok == false {
			glog.Infof("FOO h is null")
			apiKeyToHeartbeat[v.APIKeyID] = v
		} else {
			glog.Info("FOO h: %s", h.ID)
			glog.Infof("%s after %s", v.ExecutedAt, h.ExecutedAt)
			if v.ExecutedAt.After(h.ExecutedAt) {
				// replace heartbeat if this heartbeat was executed later
				apiKeyToHeartbeat[v.APIKeyID] = v
			}
		}
	}

	// convert into map heartbeat id => heartbeat
	m2 := make(map[string]Heartbeat)
	for _, v := range apiKeyToHeartbeat {
		m2[v.PersistanceID()] = v
	}

	return &m2, nil
}
