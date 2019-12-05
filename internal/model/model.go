package model

import (
	"time"
)

// AlertStatus is the status of an alert.
type AlertStatus int

const (
	// AlertStatusUnknown is unknown alert status.
	AlertStatusUnknown AlertStatus = iota
	// AlertStatusFiring is when the alert is active and firing.
	AlertStatusFiring
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
	// Start is when the alert has been started.
	Start time.Time
	// End is when the alert has been ended.
	End time.Time
	// Status is the status of the alert.
	Status AlertStatus
	// Labels is data that defines the alert.
	Labels map[string]interface{}
	// Annotations is a simple map of values that can be used to
	// add more info to the alert but don't define the alert nature
	// commonly this is used to add description, titles...
	Annotations map[string]interface{}
}

// AlertGroup is a group of alerts that share some of
// the information like the state, common metadata...
// and can be grouped in order to notify at the same
// time.
type AlertGroup struct {
	// ID is the group of alerts ID.
	ID string
	// Status is the status that share the group.
	Status AlertStatus
	// Labels is like the alerts labels but shared
	// by all the alerts.
	Labels map[string]interface{}
	// Annotations is like the alerts annotations but
	// shared by all the alerts.
	Annotations map[string]interface{}
	// Alerts are the alerts in the group.
	Alerts []Alert
}
