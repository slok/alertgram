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
	// StartsAt is when the alert has been started.
	StartsAt time.Time
	// End is when the alert has been ended.
	EndsAt time.Time
	// Status is the status of the alert.
	Status AlertStatus
	// Labels is data that defines the alert.
	Labels map[string]string
	// Annotations is a simple map of values that can be used to
	// add more info to the alert but don't define the alert nature
	// commonly this is used to add description, titles...
	Annotations map[string]string
	// GeneratorURL is the url that generated the alert (eg. Prometheus metrics).
	GeneratorURL string
}

// IsFiring returns if the alerts is firing.
func (a Alert) IsFiring() bool { return a.Status == AlertStatusFiring }

// AlertGroup is a group of alerts that share some of
// the information like the state, common metadata...
// and can be grouped in order to notify at the same
// time.
type AlertGroup struct {
	// ID is the group of alerts ID.
	ID string
	// Labels are the labels of the group.
	Labels map[string]string
	// Alerts are all the alerts in the group (firing, resolved, unknown...).
	Alerts []Alert
}

// FiringAlerts returns the firing alerts.
func (a AlertGroup) FiringAlerts() []Alert { return a.groupByStatusAlerts(AlertStatusFiring) }

// ResolvedAlerts returns the resolved alerts.
func (a AlertGroup) ResolvedAlerts() []Alert { return a.groupByStatusAlerts(AlertStatusResolved) }

// HasFiring returns true if it has firing alerts.
func (a AlertGroup) HasFiring() bool { return a.hasAlertByStatus(AlertStatusFiring) }

// HasResolved returns true if it has resolved alerts.
func (a AlertGroup) HasResolved() bool { return a.hasAlertByStatus(AlertStatusResolved) }

func (a AlertGroup) hasAlertByStatus(status AlertStatus) bool {
	for _, al := range a.Alerts {
		if al.Status == status {
			return true
		}
	}

	return false
}

func (a AlertGroup) groupByStatusAlerts(status AlertStatus) []Alert {
	var alerts []Alert
	for _, al := range a.Alerts {
		if al.Status == status {
			alerts = append(alerts, al)
		}
	}

	return alerts
}
