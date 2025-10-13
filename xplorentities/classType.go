package xplorentities

import "github.com/angelbarreiros/XPlorGo/util"

// XPlorClassType representa un tipo de evento de clase individual
type XPlorClassType struct {
	Context             string          `json:"@context"`
	ID                  *string         `json:"@id"`
	Type                string          `json:"@type"`
	InternalID          int             `json:"id"`
	StartedAt           *util.LocalTime `json:"startedAt"`
	EndedAt             *util.LocalTime `json:"endedAt"`
	Summary             string          `json:"summary"`
	Description         interface{}     `json:"description"`
	Club                *string         `json:"club"`
	Studio              *string         `json:"studio"`
	Activity            *string         `json:"activity"`
	Coach               *string         `json:"coach"`
	AttendingLimit      int             `json:"attendingLimit"`
	QueueLimit          int             `json:"queueLimit"`
	PrivateComment      interface{}     `json:"privateComment"`
	Recurrence          *string         `json:"recurrence"`
	ClassLayout         interface{}     `json:"classLayout"`
	InstructionsComment interface{}     `json:"instructionsComment"`
	OnlineLimit         interface{}     `json:"onlineLimit"`
	ExternalQuota       interface{}     `json:"externalQuota"`
	ImportRequestId     interface{}     `json:"importRequestId"`
}

type HydraSearch struct {
	Type                        string         `json:"@type"`
	HydraTemplate               string         `json:"hydra:template"`
	HydraVariableRepresentation string         `json:"hydra:variableRepresentation"`
	HydraMapping                []HydraMapping `json:"hydra:mapping"`
}

// Métodos para XPlorClassEventType
func (cet *XPlorClassType) ClassEventTypeID() (string, error) {
	return ExtractID(cet.ID, "class event type ID field is nil")
}

func (cet *XPlorClassType) ClubID() (string, error) {
	return ExtractID(cet.Club, "club ID field is nil")
}
func (cet *XPlorClassType) CoachId() (string, error) {
	return ExtractID(cet.Coach, "coach ID field is nil")
}

func (cet *XPlorClassType) StudioID() (string, error) {
	return ExtractID(cet.Studio, "studio ID field is nil")
}

func (cet *XPlorClassType) ActivityID() (string, error) {
	return ExtractID(cet.Activity, "activity ID field is nil")
}

func (cet *XPlorClassType) RecurrenceID() (string, error) {
	return ExtractID(cet.Recurrence, "recurrence ID field is nil")
}

// Método para obtener el ID interno (campo "id")
func (cet *XPlorClassType) InternalIDValue() int {
	return cet.InternalID
}
