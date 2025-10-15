package xplorentities

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

type XPlorPagination struct {
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
}
type XPlorTimeGap struct {
	StartAt *time.Time `json:"startAt"`
	EndAt   *time.Time `json:"endAt"`
}

// XPlorClassesParams represents all available query parameters for filtering classes
// Based on the hydra:search template from the API
type XPlorClassesParams struct {
	// Filter by club(s)
	Club  string   `json:"club,omitempty"`
	Clubs []string `json:"club[],omitempty"`

	// Filter by coach(es)
	Coach   string   `json:"coach,omitempty"`
	Coaches []string `json:"coach[],omitempty"`

	// Filter by activity/activities
	Activity   string   `json:"activity,omitempty"`
	Activities []string `json:"activity[],omitempty"`

	// Filter by studio(s)
	Studio  string   `json:"studio,omitempty"`
	Studios []string `json:"studio[],omitempty"`

	// Filter by recurrence(s)
	Recurrence  string   `json:"recurrence,omitempty"`
	Recurrences []string `json:"recurrence[],omitempty"`

	// Filter by attendee contact ID(s)
	AttendeeContactID  string   `json:"attendees.contactId,omitempty"`
	AttendeeContactIDs []string `json:"attendees.contactId[],omitempty"`

	// Filter by attendee state(s)
	AttendeeState  string   `json:"attendees.state,omitempty"`
	AttendeeStates []string `json:"attendees.state[],omitempty"`

	// Filter by activity groups
	ActivityGroup  string   `json:"activity.activityGroups,omitempty"`
	ActivityGroups []string `json:"activity.activityGroups[],omitempty"`

	// Date/time filters for startedAt
	StartedAtBefore         *time.Time `json:"startedAt[before],omitempty"`
	StartedAtStrictlyBefore *time.Time `json:"startedAt[strictly_before],omitempty"`
	StartedAtAfter          *time.Time `json:"startedAt[after],omitempty"`
	StartedAtStrictlyAfter  *time.Time `json:"startedAt[strictly_after],omitempty"`

	// Ordering
	OrderStartedAt string `json:"order[startedAt],omitempty"` // "asc" or "desc"

	// Other filters
	Available *bool  `json:"available,omitempty"`
	Time      string `json:"time,omitempty"`
	Archived  string `json:"archived,omitempty"`
}

// ToValues converts XPlorClassesParams to url.Values following the same pattern as families
func (p XPlorClassesParams) ToValues(values *url.Values) {
	// Single club
	if strings.TrimSpace(p.Club) != "" {
		values.Set("club", p.Club)
	}
	// Multiple clubs
	for _, club := range p.Clubs {
		if strings.TrimSpace(club) != "" {
			values.Add("club[]", club)
		}
	}

	// Single coach
	if strings.TrimSpace(p.Coach) != "" {
		values.Set("coach", p.Coach)
	}
	// Multiple coaches
	for _, coach := range p.Coaches {
		if strings.TrimSpace(coach) != "" {
			values.Add("coach[]", coach)
		}
	}

	// Single activity
	if strings.TrimSpace(p.Activity) != "" {
		values.Set("activity", p.Activity)
	}
	// Multiple activities
	for _, activity := range p.Activities {
		if strings.TrimSpace(activity) != "" {
			values.Add("activity[]", activity)
		}
	}

	// Single studio
	if strings.TrimSpace(p.Studio) != "" {
		values.Set("studio", p.Studio)
	}
	// Multiple studios
	for _, studio := range p.Studios {
		if strings.TrimSpace(studio) != "" {
			values.Add("studio[]", studio)
		}
	}

	// Single recurrence
	if strings.TrimSpace(p.Recurrence) != "" {
		values.Set("recurrence", p.Recurrence)
	}
	// Multiple recurrences
	for _, recurrence := range p.Recurrences {
		if strings.TrimSpace(recurrence) != "" {
			values.Add("recurrence[]", recurrence)
		}
	}

	// Single attendee contact ID
	if strings.TrimSpace(p.AttendeeContactID) != "" {
		values.Set("attendees.contactId", p.AttendeeContactID)
	}
	// Multiple attendee contact IDs
	for _, contactID := range p.AttendeeContactIDs {
		if strings.TrimSpace(contactID) != "" {
			values.Add("attendees.contactId[]", contactID)
		}
	}

	// Single attendee state
	if strings.TrimSpace(p.AttendeeState) != "" {
		values.Set("attendees.state", p.AttendeeState)
	}
	// Multiple attendee states
	for _, state := range p.AttendeeStates {
		if strings.TrimSpace(state) != "" {
			values.Add("attendees.state[]", state)
		}
	}

	// Single activity group
	if strings.TrimSpace(p.ActivityGroup) != "" {
		values.Set("activity.activityGroups", p.ActivityGroup)
	}
	// Multiple activity groups
	for _, group := range p.ActivityGroups {
		if strings.TrimSpace(group) != "" {
			values.Add("activity.activityGroups[]", group)
		}
	}

	// Date/time filters
	if p.StartedAtBefore != nil {
		values.Set("startedAt[before]", p.StartedAtBefore.Format(time.RFC3339))
	}
	if p.StartedAtStrictlyBefore != nil {
		values.Set("startedAt[strictly_before]", p.StartedAtStrictlyBefore.Format(time.RFC3339))
	}
	if p.StartedAtAfter != nil {
		values.Set("startedAt[after]", p.StartedAtAfter.Format(time.RFC3339))
	}
	if p.StartedAtStrictlyAfter != nil {
		values.Set("startedAt[strictly_after]", p.StartedAtStrictlyAfter.Format(time.RFC3339))
	}

	// Ordering
	if strings.TrimSpace(p.OrderStartedAt) != "" {
		values.Set("order[startedAt]", p.OrderStartedAt)
	}

	// Other filters
	if p.Available != nil {
		values.Set("available", strconv.FormatBool(*p.Available))
	}
	if strings.TrimSpace(p.Time) != "" {
		values.Set("time", p.Time)
	}
	if strings.TrimSpace(p.Archived) != "" {
		values.Set("archived", p.Archived)
	}
}

// BuildPaginationQueryParams creates query parameters from pagination
// Only adds parameters if they are not empty
func BuildPaginationQueryParams(pagination *XPlorPagination) url.Values {
	queryParams := url.Values{}

	if pagination != nil {
		if pagination.ItemsPerPage > 0 {
			queryParams.Set("itemsPerPage", strconv.Itoa(pagination.ItemsPerPage))
		}

		if pagination.Page > 0 {
			queryParams.Set("page", strconv.Itoa(pagination.Page))
		}
	}

	return queryParams
}

// BuildPaginationAndTimeGapParams creates query parameters from pagination and timeGap
func BuildPaginationAndTimeGapParams(pagination *XPlorPagination, timeGap *XPlorTimeGap) url.Values {
	queryParams := url.Values{}

	// Set pagination parameters
	if pagination != nil {
		if pagination.ItemsPerPage > 0 {
			queryParams.Set("itemsPerPage", strconv.Itoa(pagination.ItemsPerPage))
		} else {
			queryParams.Set("itemsPerPage", "10") // default value
		}

		if pagination.Page > 0 {
			queryParams.Set("page", strconv.Itoa(pagination.Page))
		} else {
			queryParams.Set("page", "1") // default value
		}
	} else {
		// Default values when pagination is nil
		queryParams.Set("itemsPerPage", "10")
		queryParams.Set("page", "1")
	}

	// Set time gap parameters
	if timeGap != nil {
		if timeGap.StartAt != nil {
			queryParams.Set("startAt", timeGap.StartAt.Format("2006-01-02T15:04"))
		}
		if timeGap.EndAt != nil {
			queryParams.Set("endAt", timeGap.EndAt.Format("2006-01-02T15:04"))
		}
	}

	return queryParams
}

// BuildClassesQueryParams creates query parameters from XPlorClassesParams
func BuildClassesQueryParams(params *XPlorClassesParams) url.Values {
	queryParams := url.Values{}

	if params != nil {
		params.ToValues(&queryParams)
	}

	return queryParams
}

// BuildClassesWithPaginationQueryParams combines classes params with pagination
func BuildClassesWithPaginationQueryParams(params *XPlorClassesParams, pagination *XPlorPagination) url.Values {
	queryParams := url.Values{}

	// Add classes parameters
	if params != nil {
		params.ToValues(&queryParams)
	}

	// Add pagination parameters
	if pagination != nil {
		if pagination.ItemsPerPage > 0 {
			queryParams.Set("itemsPerPage", strconv.Itoa(pagination.ItemsPerPage))
		}
		if pagination.Page > 0 {
			queryParams.Set("page", strconv.Itoa(pagination.Page))
		}
	}

	return queryParams
}
