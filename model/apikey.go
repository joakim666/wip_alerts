package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

type APIKey struct {
	ID             string // uuid
	RefreshTokenID string // uuid of refresh token
	CreatedAt      time.Time
}

func (a APIKey) PersistanceID() string {
	return a.ID
}

func NewAPIKey() *APIKey {
	var a APIKey
	uuid := uuid.NewV4()
	a.ID = uuid.String()
	a.CreatedAt = time.Now()
	return &a
}

func SaveAPIKeys(db *bolt.DB, accountUUID string, apikeys *map[string]APIKey) error {
	return BoltSaveAccountObjects(db, accountUUID, "APIKeys", BoltMap(apikeys))
}

func ListAPIKeys(db *bolt.DB, accountUUID string) (*map[string]APIKey, error) {
	m, err := BoltGetAccountObjects(db, accountUUID, "APIKeys", reflect.TypeOf(APIKey{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing APIKey
	m2 := make(map[string]APIKey)
	for _, v := range *m {
		d := v.(*APIKey)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}
