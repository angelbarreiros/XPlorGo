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
	StartedAt                          util.LocalTime  `json:"startedAt"`
	EndedAt                            util.LocalTime  `json:"endedAt"`
	Summary                            string          `json:"summary"`
	Description                        *string         `json:"description"`
	CreatedAt                          util.LocalTime  `json:"createdAt"`
	CreatedBy                          string          `json:"createdBy"`
	UpdatedAt                          util.LocalTime  `json:"updatedAt"`
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

// FirstPageNumber extrae el número de página de la URL first
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

// LastPageNumber extrae el número de página de la URL last
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

// NextPageNumber extrae el número de página de la URL next
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

// extractPageNumber extrae el parámetro 'page' de una URL
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
func (a Attendee) ContactCounterId() (string, error) {
	return ExtractID(a.ContactCounterIdUsed, "contact counter ID field is nil")
}

// FullName devuelve el nombre completo del asistente
func (a Attendee) FullName() string {
	return a.ContactGivenName + " " + a.ContactFamilyName
}

// HasContact verifica si el asistente tiene información de contacto
func (a Attendee) HasContact() bool {
	return a.ContactID != nil
}

// ContactIDValue extrae el ID numérico del contacto desde el campo ContactID
func (a Attendee) ContactIDValue() (string, error) {
	return ExtractID(a.ContactID, "contact ID field is nil")
}

// FullName devuelve el nombre completo del asistente

// ---- Club ----
func (e XPlorEvent) ClubID() (string, error) {
	return ExtractID(e.Club, "club field is nil")
}

// ---- Studio ----
func (e XPlorEvent) StudioID() (string, error) {
	return ExtractID(e.Studio, "studio field is nil")
}

// ---- Activity ----
func (e XPlorEvent) ActivityID() (string, error) {
	return ExtractID(e.Activity, "activity field is nil")
}

// ---- Coach ----
func (e XPlorEvent) CoachID() (string, error) {
	return ExtractID(e.Coach, "coach field is nil")
}

// ---- Recurrence ----
func (e XPlorEvent) RecurrenceID() (string, error) {
	return ExtractID(e.Recurrence, "recurrence field is nil")
}

// ---- ClassEvent ----
func (e XPlorEvent) EventID() (string, error) {
	return ExtractID(e.ID, "event ID field is nil")
}
