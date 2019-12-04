package model

import "time"

// AlertStatus is the status of an alert.
type AlertStatus int

const (
	// AlertStatusUnknown is unknown alert status.
	AlertStatusUnknown AlertStatus = iota
	// AlertStatusTriggering is when the alert is being triggered.
	AlertStatusTriggering
	// AlertStatusResolved is when the alert was being triggered and now
	// is not being triggered anymore.
	AlertStatusResolved
)

// Alert represents an alert.
type Alert struct {
	// ID is the ID of the alert.
	ID string
	// Name is the name of the alert.
	Name string
	// TS is when the alert has been created.
	TS time.Time
	// Status is the status of the alert.
	Status AlertStatus
	// Description is a description text of the alert.
	Description string
	// Metadata is a simple map of values that can be used to
	// add more info to the alert.
	Metadata map[string]interface{}
	// Where the alert is coming from.
	Source string
}
