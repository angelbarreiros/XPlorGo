package xplorentities

// XPlorCourseWaitingListRecurrence contains a course-level waiting-list
// recurrence resource. Its payload is intentionally flexible because this
// resource is not described in the public Resamania OpenAPI specification.
type XPlorCourseWaitingListRecurrence map[string]any

// XPlorCourseWaitingListRecurrences contains course-level waiting-list
// recurrence resources returned by the Hydra collection endpoint.
type XPlorCourseWaitingListRecurrences struct {
	Context    any                                `json:"@context"`
	ID         string                             `json:"@id"`
	Type       string                             `json:"@type"`
	Items      []XPlorCourseWaitingListRecurrence `json:"hydra:member"`
	Pagination *HydraView                         `json:"hydra:view,omitempty"`
}
