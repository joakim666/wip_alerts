package model

import (
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

// AlertPriority shows the priority of the alert, possible values HighPriority, NormalPriority and LowPriority
type AlertPriority string

// AlertStatus indicates the status of the alert
type AlertStatus string

const (
	HighPriority AlertPriority = "high"
	NormalPriority AlertPriority = "normal"
	LowPriority AlertPriority = "low"

	NewStatus AlertStatus = "new"
	SeenStatus AlertStatus = "seen"
	ArchivedStatus AlertStatus = "archived"
)

type Alert struct {
	ID               string // uuid
	APIKeyID         string // uuid of api key that sent this heartbeat
	Title            string
	ShortDescription string
	LongDescription  string
	Priority         AlertPriority
	Status		 AlertStatus
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
	a.Status = NewStatus
	return &a
}

// ListsAlerts returns all alerts for the given account
func ListAlerts(db *bolt.DB, accountUUID string) (*map[string]Alert, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Alerts", reflect.TypeOf(Alert{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Alert
	m2 := make(map[string]Alert)
	for _, v := range *m {
		a := v.(*Alert)
		m2[v.PersistanceID()] = *a
	}

	return &m2, nil
}

// ListNonArchivedAlerts list all alerts that do not have status "archived" for the given account
func ListNonArchivedAlerts(db *bolt.DB, accountUUID string) (*map[string]Alert, error) {
	m, err := BoltGetAccountObjects(db, ParentID(accountUUID), "Alerts", reflect.TypeOf(Alert{}))
	if err != nil {
		return nil, err
	}

	// convert to map containing Alert and filter
	m2 := make(map[string]Alert)
	for _, v := range *m {
		a := v.(*Alert)
		if ArchivedStatus != a.Status { // do not included alerts with status archived
			m2[v.PersistanceID()] = *a
		}
	}

	return &m2, nil
}
