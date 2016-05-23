package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

// AlertPriority shows the priority of the alert
type AlertPriority string

const (
	HighPriority AlertPriority = "high"
	NormalPriority AlertPriority = "normal"
	LowPriority AlertPriority = "low"
)

type Alert struct {
	ID               string // uuid
	APIKeyID         string // uuid of api key that sent this heartbeat
	Title            string
	ShortDescription string
	LongDescription  string
	Priority         AlertPriority
	TriggeredAt      time.Time
	CreatedAt        time.Time
}

func (a Alert) PersistanceID() string {
	return a.ID
}

// Save the alert attached to the given accountUUID
func (a Alert) Save(db *bolt.DB, accountUUID string) error {
	return BoltSaveAccountObjects(db, ParentID(accountUUID), "Alerts", BoltSingle(&a))
}

// NewAlert creates a new Alert. APIKeyID is mandatory
func NewAlert(apiKeyID string) *Alert {
	var a Alert
	uuid := uuid.NewV4()
	a.ID = uuid.String()
	a.APIKeyID = apiKeyID
	a.CreatedAt = time.Now()
	return &a
}

func ListAlerts(db *bolt.DB, accountUUID string) (*map[string]Alert, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Alerts", reflect.TypeOf(Alert{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Alert
	m2 := make(map[string]Alert)
	for _, v := range *m {
		hb := v.(*Alert)
		m2[v.PersistanceID()] = *hb
	}

	return &m2, nil
}
