package xplorentities

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/angelbarreiros/XPlorGo/util"
)

// Colección Hydra
type XPlorEvents struct {
	Context    string       `json:"@context"`
	ID         string       `json:"@id"`
	Type       string       `json:"@type"`
	Events     []XPlorEvent `json:"hydra:member"`
	Pagination HydraView    `json:"hydra:view"`
}

// Un evento de clase individual
type XPlorEvent struct {
	ID                                 *string         `json:"@id"`
	Type                               string          `json:"@type"`
	Club                               *string         `json:"club"`
	Studio                             *string         `json:"studio"`
	Activity                           *string         `json:"activity"`
	Coach                              *string         `json:"coach"`
	AttendingLimit                     int             `json:"attendingLimit"`
	QueueLimit                         int             `json:"queueLimit"`
	PrivateComment                     *string         `json:"privateComment"`
	Recurrence                         *string         `json:"recurrence"`
	ClassLayout                        *string         `json:"classLayout"`
	ClassLayoutConfiguration           []string        `json:"classLayoutConfiguration"`
	InstructionsComment                *string         `json:"instructionsComment"`
	OnlineLimit                        *int            `json:"onlineLimit"`
	ExternalQuota                      *string         `json:"externalQuota"`
	CoachAvailable                     bool            `json:"coachAvailable"`
	StartedAt                          *util.LocalTime `json:"startedAt"`
	EndedAt                            *util.LocalTime `json:"endedAt"`
	Summary                            string          `json:"summary"`
	Description                        *string         `json:"description"`
	CreatedAt                          *util.LocalTime `json:"createdAt"`
	CreatedBy                          string          `json:"createdBy"`
	UpdatedAt                          *util.LocalTime `json:"updatedAt"`
	UpdatedBy                          string          `json:"updatedBy"`
	DeletedAt                          *util.LocalTime `json:"deletedAt"`
	DeletedBy                          *string         `json:"deletedBy"`
	Processing                         bool            `json:"processing"`
	BookedAttendees                    []Attendee      `json:"bookedAttendees"`
	QueuedAttendees                    []Attendee      `json:"queuedAttendees"`
	AutoPromoteQueuedAttendeesPossible bool            `json:"autoPromoteQueuedAttendeesPossible"`
	AttendeeRemaining                  int             `json:"attendeeRemaining"`
	QueueRemaining                     int             `json:"queueRemaining"`
	DefaultOnlineLimit                 *int            `json:"defaultOnlineLimit"`
}

type HydraMapping struct {
	Type     string `json:"@type"`
	Variable string `json:"variable"`
	Property string `json:"property"`
	Required bool   `json:"required"`
}

// Vista Hydra (para paginación)
type HydraView struct {
	HydraFirst string `json:"hydra:first"`
	HydraLast  string `json:"hydra:last"`
	HydraNext  string `json:"hydra:next"`
}

// FirstPageNumber extracts the page number from the first URL
func (hv HydraView) FirstPageNumber() (int, error) {
	if hv.HydraFirst == "" {
		return 0, errors.New("hydra:first URL is empty")
	}
	pageStr, err := extractPageNumber(hv.HydraFirst)
	if err != nil {
		return 0, err
	}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, errors.New("page parameter is not a valid integer")
	}
	return pageInt, nil
}

// LastPageNumber extracts the page number from the last URL
func (hv HydraView) LastPageNumber() (int, error) {
	if hv.HydraLast == "" {
		return 1, nil // Return 1 if there's no last page (only one page available)
	}
	pageStr, err := extractPageNumber(hv.HydraLast)
	if err != nil {
		return 1, nil // Return 1 if we can't extract page number
	}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return 1, nil // Return 1 if page parameter is not a valid integer
	}
	return pageInt, nil
}

// NextPageNumber extracts the page number from the next URL
func (hv HydraView) NextPageNumber() (int, error) {
	if hv.HydraNext == "" {
		return 0, errors.New("hydra:next URL is empty")
	}
	pageStr, err := extractPageNumber(hv.HydraNext)
	if err != nil {
		return 0, err
	}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, errors.New("page parameter is not a valid integer")
	}
	return pageInt, nil
}

// extractPageNumber extracts the page parameter from a URL query string
func extractPageNumber(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.New("invalid URL format")
	}

	pageParam := parsedURL.Query().Get("page")
	if pageParam == "" {
		return "", errors.New("page parameter not found in URL")
	}

	return pageParam, nil
}

// Attendee representa un asistente en un evento
// Attendee representa un asistente en un evento
type Attendee struct {
	ContactID            *string         `json:"contactId"`
	ContactGivenName     string          `json:"contactGivenName"`
	ContactFamilyName    string          `json:"contactFamilyName"`
	ContactNumber        *string         `json:"contactNumber"`
	ContactClubID        *string         `json:"contactClubId"`
	ContactDetails       *string         `json:"contactDetails"`
	ContactCreatedAt     *util.LocalTime `json:"contactCreatedAt"`
	ContactTagIdUsed     *string         `json:"contactTagUsed"`
	ContactCounterIdUsed *string         `json:"contactCounterUsed"`
	ContactPictureID     *string         `json:"contactPictureId"`
	ContactChannelUsed   *string         `json:"contactChannelUsed"`
	Warnings             []string        `json:"warnings"`
	CreatedAt            util.LocalTime  `json:"createdAt"`
	CreatedBy            string          `json:"createdBy"`
	CanceledAt           *util.LocalTime `json:"canceledAt"`
	CanceledBy           *string         `json:"canceledBy"`
	ValidatedAt          *util.LocalTime `json:"validatedAt"`
	ValidatedBy          *string         `json:"validatedBy"`
	QueuedAt             *util.LocalTime `json:"queuedAt"`
	QueuedBy             *string         `json:"queuedBy"`
	DeletedAt            *util.LocalTime `json:"deletedAt"`
	DeletedBy            *string         `json:"deletedBy"`
	State                string          `json:"state"`
	CostSignIn           *float64        `json:"costSignIn"`
	BookedItem           *string         `json:"bookedItem"`
	Showed               bool            `json:"showed"`
	Broker               *string         `json:"broker"`
	FromAttendeeGroup    bool            `json:"fromAttendeeGroup"`
	CancelDelayOver      bool            `json:"cancelDelayOver"`
}

func (a Attendee) ContactTagId() (string, error) {
	return ExtractID(a.ContactTagIdUsed, "contact tag ID field is nil")
}

// ContactCounterId extracts the contact counter ID from the contactCounterUsed field
func (a Attendee) ContactCounterId() (string, error) {
	return ExtractID(a.ContactCounterIdUsed, "contact counter ID field is nil")
}

// FullName returns the full name of the attendee
func (a Attendee) FullName() string {
	return a.ContactGivenName + " " + a.ContactFamilyName
}

// HasContact checks if the attendee has contact information
func (a Attendee) HasContact() bool {
	return a.ContactID != nil
}

// ContactIDValue extracts the numeric contact ID from the ContactID field
func (a Attendee) ContactIDValue() (string, error) {
	return ExtractID(a.ContactID, "contact ID field is nil")
}

// ClubID extracts the club ID from the club field
func (e XPlorEvent) ClubID() (string, error) {
	return ExtractID(e.Club, "club field is nil")
}

// StudioID extracts the studio ID from the studio field
func (e XPlorEvent) StudioID() (string, error) {
	return ExtractID(e.Studio, "studio field is nil")
}

// ActivityID extracts the activity ID from the activity field
func (e XPlorEvent) ActivityID() (string, error) {
	return ExtractID(e.Activity, "activity field is nil")
}

// CoachID extracts the coach ID from the coach field
func (e XPlorEvent) CoachID() (string, error) {
	return ExtractID(e.Coach, "coach field is nil")
}

// RecurrenceID extracts the recurrence ID from the recurrence field
func (e XPlorEvent) RecurrenceID() (string, error) {
	return ExtractID(e.Recurrence, "recurrence field is nil")
}

// EventID extracts the event ID from the @id field
func (e XPlorEvent) EventID() (string, error) {
	return ExtractID(e.ID, "event ID field is nil")
}
