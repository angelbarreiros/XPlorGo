package xplorentities

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

type XPlorAttendees struct {
	Context    string          `json:"@context"`
	ID         string          `json:"@id"`
	Type       string          `json:"@type"`
	Attendees  []XPlorAttendee `json:"hydra:member"`
	Pagination *HydraView      `json:"hydra:view,omitempty"`
}

type XPlorAttendee struct {
	AtID               *string          `json:"@id"`
	AtType             *string          `json:"@type"`
	ContactId          *string          `json:"contactId"`
	ContactGivenName   *string          `json:"contactGivenName"`
	ContactFamilyName  *string          `json:"contactFamilyName"`
	ContactNumber      *string          `json:"contactNumber"`
	ContactClubId      *string          `json:"contactClubId"`
	ContactDetails     *string          `json:"contactDetails"`
	ContactCreatedAt   *string          `json:"contactCreatedAt"`
	ContactTagUsed     *string          `json:"contactTagUsed"`
	ContactCounterUsed *string          `json:"contactCounterUsed"`
	ContactPictureId   *string          `json:"contactPictureId"`
	ContactChannelUsed *string          `json:"contactChannelUsed"`
	ClassEvent         *XPlorClassEvent `json:"classEvent"`

	CancelReason *string  `json:"cancelReason"`
	Warnings     []string `json:"warnings"`
	CreatedAt    *string  `json:"createdAt"`
	CreatedBy    *string  `json:"createdBy"`
	CanceledAt   *string  `json:"canceledAt"`
	CanceledBy   *string  `json:"canceledBy"`
	ValidatedAt  *string  `json:"validatedAt"`
	ValidatedBy  *string  `json:"validatedBy"`
	QueuedAt     *string  `json:"queuedAt"`
	QueuedBy     *string  `json:"queuedBy"`
	DeletedAt    *string  `json:"deletedAt"`
	DeletedBy    *string  `json:"deletedBy"`

	State           *string `json:"state"`
	CostSignIn      *string `json:"costSignIn"`
	CreditSignOut   *string `json:"creditSignOut"`
	BookedItem      *string `json:"bookedItem"`
	Showed          bool    `json:"showed"`
	AttendeeGroup   *string `json:"attendeeGroup"`
	Broker          *string `json:"broker"`
	ActivityName    *string `json:"activityName"`
	ClassEventStart *string `json:"classEventStartedAt"`
	ClassLayout     *string `json:"classLayout"`
	CancelDelayOver bool    `json:"cancelDelayOver"`
}

// Subestructura ClassEvent
type XPlorClassEvent struct {
	AtID               *string `json:"@id"`
	AtType             *string `json:"@type"`
	Club               *string `json:"club"`
	Studio             *string `json:"studio"`
	Activity           *string `json:"activity"`
	Coach              *string `json:"coach"`
	AttendingLimit     *int    `json:"attendingLimit"`
	QueueLimit         *int    `json:"queueLimit"`
	PrivateComment     *string `json:"privateComment"`
	OnlineLimit        *int    `json:"onlineLimit"`
	ExternalQuota      *int    `json:"externalQuota"`
	StartedAt          *string `json:"startedAt"`
	EndedAt            *string `json:"endedAt"`
	DefaultOnlineLimit *int    `json:"defaultOnlineLimit"`
}

// ID extracts the attendee ID from the @id field
func (a XPlorAttendee) ID() (string, error) {
	if a.AtID == nil {
		return "", errors.New("attendee ID field is nil")
	}
	base := path.Base(*a.AtID)
	return base, nil
}

// ContactID extracts the contact ID from the contactId field
func (a XPlorAttendee) ContactID() (string, error) {
	if a.ContactId == nil {
		return "", errors.New("contact ID field is nil")
	}
	base := path.Base(*a.ContactId)
	return base, nil
}
func (a XPlorAttendee) ContactTagId() (string, error) {
	return ExtractID(a.ContactTagUsed, "contactTag ID field is nil")
}

// ContactCounterId extracts the contact counter ID from the contactCounterUsed field
func (a XPlorAttendee) ContactCounterId() (string, error) {
	return ExtractID(a.ContactCounterUsed, "contactCounter ID field is nil")
}

// ContactClubID extracts the contact club ID from the contactClubId field
func (a XPlorAttendee) ContactClubID() (string, error) {
	if a.ContactClubId == nil {
		return "", errors.New("contactClub ID field is nil")
	}
	base := path.Base(*a.ContactClubId)
	return base, nil
}

// ID extracts the class event ID from the @id field
func (ce XPlorClassEvent) ID() (string, error) {
	if ce.AtID == nil {
		return "", errors.New("class event ID field is nil")
	}
	base := path.Base(*ce.AtID)
	return base, nil
}

// ClubID extracts the club ID from the club field
func (ce XPlorClassEvent) ClubID() (string, error) {
	if ce.Club == nil {
		return "", errors.New("club ID field is nil")
	}
	base := path.Base(*ce.Club)
	return base, nil
}

// StudioID extracts the studio ID from the studio field
func (ce XPlorClassEvent) StudioID() (string, error) {
	if ce.Studio == nil {
		return "", errors.New("studio ID field is nil")
	}
	base := path.Base(*ce.Studio)
	return base, nil
}

// ActivityID extracts the activity ID from the activity field
func (ce XPlorClassEvent) ActivityID() (string, error) {
	if ce.Activity == nil {
		return "", errors.New("activity ID field is nil")
	}
	base := path.Base(*ce.Activity)
	return base, nil
}

// CoachID extracts the coach ID from the coach field
func (ce XPlorClassEvent) CoachID() (string, error) {
	if ce.Coach == nil {
		return "", errors.New("coach ID field is nil")
	}
	base := path.Base(*ce.Coach)
	return base, nil
}

// StartedAtTime converts StartedAt string to time.Time
func (ce XPlorClassEvent) StartedAtTime() (*time.Time, error) {
	return parseTimeString(ce.StartedAt)
}

// EndedAtTime converts EndedAt string to time.Time
func (ce XPlorClassEvent) EndedAtTime() (*time.Time, error) {
	return parseTimeString(ce.EndedAt)
}

type XTime struct {
	time.Time
}

// UnmarshalJSON unmarshals JSON bytes into XTime, supporting RFC3339 and custom formats
func (xt *XTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}

	// probar RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		xt.Time = t
		return nil
	}

	// probar sin zona
	t, err = time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		xt.Time = t
		return nil
	}

	return fmt.Errorf("XTime: cannot parse %q", s)
}

func parseTimeString(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}

	// RFC3339 (ej: 2023-07-25T16:21:39Z o con offset)
	if t, err := time.Parse(time.RFC3339, *s); err == nil {
		return &t, nil
	}

	// Sin zona horaria (ej: 2023-07-25T16:21:39)
	if t, err := time.Parse("2006-01-02T15:04:05", *s); err == nil {
		return &t, nil
	}

	return nil, errors.New("cannot parse time string: " + *s)
}
func (a XPlorAttendee) CreatedAtTime() (*time.Time, error) {
	return parseTimeString(a.CreatedAt)
}

// CanceledAtTime converts CanceledAt string to time.Time
func (a XPlorAttendee) CanceledAtTime() (*time.Time, error) {
	return parseTimeString(a.CanceledAt)
}

// ValidatedAtTime converts ValidatedAt string to time.Time
func (a XPlorAttendee) ValidatedAtTime() (*time.Time, error) {
	return parseTimeString(a.ValidatedAt)
}

// QueuedAtTime converts QueuedAt string to time.Time
func (a XPlorAttendee) QueuedAtTime() (*time.Time, error) {
	return parseTimeString(a.QueuedAt)
}

// DeletedAtTime converts DeletedAt string to time.Time
func (a XPlorAttendee) DeletedAtTime() (*time.Time, error) {
	return parseTimeString(a.DeletedAt)
}

// ClassEventStartTime converts ClassEventStart string to time.Time
func (a XPlorAttendee) ClassEventStartTime() (*time.Time, error) {
	return parseTimeString(a.ClassEventStart)
}
