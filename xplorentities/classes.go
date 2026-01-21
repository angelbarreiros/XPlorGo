package xplorentities

import (
	"net/url"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

// LocalTime is an alias for util.LocalTime
type LocalTime = util.LocalTime

// XPlorClasses representa la colección de eventos de clase
type XPlorClasses struct {
	Context    string       `json:"@context"`
	ID         string       `json:"@id"`
	Type       string       `json:"@type"`
	Classes    []XPlorClass `json:"hydra:member"`
	Pagination HydraView    `json:"hydra:view"`
}

// XPlorClass representa un evento de clase individual
type XPlorClass struct {
	ID                                 *string         `json:"@id"`
	Type                               string          `json:"@type"`
	Context                            interface{}     `json:"@context"`
	Club                               *string         `json:"club"`
	Studio                             *string         `json:"studio"`
	Activity                           *string         `json:"activity"`
	Coach                              *string         `json:"coach"`
	AttendingLimit                     *int            `json:"attendingLimit"`
	QueueLimit                         int             `json:"queueLimit"`
	PrivateComment                     interface{}     `json:"privateComment"`
	Recurrence                         *string         `json:"recurrence"`
	ClassLayout                        interface{}     `json:"classLayout"`
	ClassLayoutConfiguration           []interface{}   `json:"classLayoutConfiguration"`
	DisabledItems                      []string        `json:"disabledItems"`
	InstructionsComment                interface{}     `json:"instructionsComment"`
	OnlineLimit                        interface{}     `json:"onlineLimit"`
	ExternalQuota                      interface{}     `json:"externalQuota"`
	CoachAvailable                     bool            `json:"coachAvailable"`
	StartedAt                          LocalTime       `json:"startedAt"`
	EndedAt                            LocalTime       `json:"endedAt"`
	Summary                            string          `json:"summary"`
	Description                        interface{}     `json:"description"`
	CreatedAt                          *LocalTime      `json:"createdAt"`
	CreatedBy                          string          `json:"createdBy"`
	UpdatedAt                          *LocalTime      `json:"updatedAt"`
	UpdatedBy                          string          `json:"updatedBy"`
	DeletedAt                          *util.LocalTime `json:"deletedAt"`
	DeletedBy                          interface{}     `json:"deletedBy"`
	ArchivedAt                         interface{}     `json:"archivedAt"`
	ArchivedBy                         interface{}     `json:"archivedBy"`
	Processing                         bool            `json:"processing"`
	BookedAttendees                    []Attendee      `json:"bookedAttendees"`
	QueuedAttendees                    []Attendee      `json:"queuedAttendees"`
	AutoPromoteQueuedAttendeesPossible bool            `json:"autoPromoteQueuedAttendeesPossible"`
	AttendeeRemaining                  int             `json:"attendeeRemaining"`
	QueueRemaining                     int             `json:"queueRemaining"`
	DefaultOnlineLimit                 interface{}     `json:"defaultOnlineLimit"`
}

// Métodos para obtener IDs

func (c XPlorClass) ClassEventID() (string, error) {
	return ExtractID(c.ID, "class event ID field is nil")
}

func (c XPlorClass) ClubID() (string, error) {
	return ExtractID(c.Club, "club ID field is nil")
}

func (c XPlorClass) StudioID() (string, error) {
	return ExtractID(c.Studio, "studio ID field is nil")
}

func (c XPlorClass) ActivityID() (string, error) {
	return ExtractID(c.Activity, "activity ID field is nil")
}

func (c XPlorClass) CoachID() (string, error) {
	return ExtractID(c.Coach, "coach ID field is nil")
}

func (c XPlorClass) RecurrenceID() (string, error) {
	return ExtractID(c.Recurrence, "recurrence ID field is nil")
}

// Métodos de utilidad para fechas
func (c XPlorClass) GetStartedAt() time.Time {
	return c.StartedAt.Time
}

func (c XPlorClass) GetEndedAt() time.Time {
	return c.EndedAt.Time
}

func (c XPlorClass) GetCreatedAt() time.Time {
	return c.CreatedAt.Time
}

func (c XPlorClass) GetUpdatedAt() time.Time {
	return c.UpdatedAt.Time
}

// Métodos para verificar disponibilidad
func (c XPlorClass) HasAvailableSpots() bool {
	return c.AttendeeRemaining > 0
}

func (c XPlorClass) HasQueueSpots() bool {
	return c.QueueRemaining > 0
}

func (c XPlorClass) IsActive() bool {
	return c.DeletedAt == nil && c.ArchivedAt == nil
}

func (c XPlorClass) IsDeleted() bool {
	return c.DeletedAt != nil
}

// GetAllContactIDs extracts all contact IDs from both booked and queued attendees
func (c XPlorClass) GetAllContactIDs() ([]string, error) {
	var contactIDs []string

	// Extract contact IDs from all attendees
	for _, attendee := range c.BookedAttendees {
		if attendee.ContactID != nil {
			contactID, err := ExtractID(attendee.ContactID, "contact ID field is nil")
			if err != nil {
				return nil, err
			}
			contactIDs = append(contactIDs, contactID)
		}
	}

	return contactIDs, nil
}

// XPlorClassesParams represents the search parameters for classes
type XPlorClassesParams struct {
	Club                    *string
	Clubs                   []string
	Coach                   *string
	Coaches                 []string
	Activity                *string
	Activities              []string
	Studio                  *string
	Studios                 []string
	Recurrence              *string
	Recurrences             []string
	AttendeeContactID       *string
	AttendeeContactIDs      []string
	AttendeeState           *string
	AttendeeStates          []string
	ActivityGroup           *string
	ActivityGroups          []string
	StartedAtBefore         *time.Time
	StartedAtStrictlyBefore *time.Time
	StartedAtAfter          *time.Time
	StartedAtStrictlyAfter  *time.Time
	OrderStartedAt          *string // asc or desc
	Available               *string
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorClassesParams) ToValues(enterpriseName string, values *url.Values) {
	// Single value filters
	if p.Club != nil {
		values.Set("club", "/"+enterpriseName+"/clubs/"+*p.Club)
	}
	if p.Coach != nil {
		values.Set("coach", *p.Coach)
	}
	if p.Activity != nil {
		values.Set("activity", *p.Activity)
	}
	if p.Studio != nil {
		values.Set("studio", *p.Studio)
	}
	if p.Recurrence != nil {
		values.Set("recurrence", *p.Recurrence)
	}
	if p.AttendeeContactID != nil {
		values.Set("attendees.contactId", *p.AttendeeContactID)
	}
	if p.AttendeeState != nil {
		values.Set("attendees.state", *p.AttendeeState)
	}
	if p.ActivityGroup != nil {
		values.Set("activity.activityGroups", *p.ActivityGroup)
	}
	if p.OrderStartedAt != nil {
		values.Set("order[startedAt]", *p.OrderStartedAt)
	}
	if p.Available != nil {
		values.Set("available", *p.Available)
	}

	// Array filters
	for _, club := range p.Clubs {
		values.Add("club[]", club)
	}
	for _, coach := range p.Coaches {
		values.Add("coach[]", coach)
	}
	for _, activity := range p.Activities {
		values.Add("activity[]", activity)
	}
	for _, studio := range p.Studios {
		values.Add("studio[]", studio)
	}
	for _, recurrence := range p.Recurrences {
		values.Add("recurrence[]", recurrence)
	}
	for _, contactID := range p.AttendeeContactIDs {
		values.Add("attendees.contactId[]", contactID)
	}
	for _, state := range p.AttendeeStates {
		values.Add("attendees.state[]", state)
	}
	for _, group := range p.ActivityGroups {
		values.Add("activity.activityGroups[]", group)
	}

	// Date filters
	if p.StartedAtBefore != nil {
		values.Set("startedAt[before]", p.StartedAtBefore.Format("2006-01-02 15:04:05"))
	}
	if p.StartedAtStrictlyBefore != nil {
		values.Set("startedAt[strictly_before]", p.StartedAtStrictlyBefore.Format("2006-01-02 15:04:05"))
	}
	if p.StartedAtAfter != nil {
		values.Set("startedAt[after]", p.StartedAtAfter.Format("2006-01-02 15:04:05"))
	}
	if p.StartedAtStrictlyAfter != nil {
		values.Set("startedAt[strictly_after]", p.StartedAtStrictlyAfter.Format("2006-01-02 15:04:05"))
	}
}
