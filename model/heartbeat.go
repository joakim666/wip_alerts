package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
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
