package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

type Renewal struct {
	ID             string // uuid
	RefreshTokenID string // uuid of refresh token
	CreatedAt      time.Time
}

func (r Renewal) PersistanceID() string {
	return r.ID
}

func NewRenewal() *Renewal {
	var r Renewal
	uuid := uuid.NewV4()
	r.ID = uuid.String()
	r.CreatedAt = time.Now()
	return &r
}

func SaveRenewals(db *bolt.DB, accountUUID string, renewals *map[string]Renewal) error {
	return BoltSaveAccountObjects(db, accountUUID, "Renewals", BoltMap(renewals))
}

func ListRenewals(db *bolt.DB, accountUUID string) (*map[string]Renewal, error) {
	m, err := BoltGetAccountObjects(db, accountUUID, "Renewals", reflect.TypeOf(Renewal{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Renewal
	m2 := make(map[string]Renewal)
	for _, v := range *m {
		d := v.(*Renewal)
		m2[v.PersistanceID()] = *d
	}

	return &m2, nil
}
