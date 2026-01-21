package xplorentities

import (
	"errors"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/angelbarreiros/XPlorGo/util"
)

// Estructuras principales para Recurrence Collection
type XPlorRecurrences struct {
	Context     string            `json:"@context"`
	ID          string            `json:"@id"`
	Type        string            `json:"@type"`
	Recurrences []XPlorRecurrence `json:"hydra:member"`
	Pagination  HydraView         `json:"hydra:view"`
}

type XPlorRecurrence struct {
	ID                           *string         `json:"@id"`
	Type                         string          `json:"@type"`
	StartedAt                    util.LocalTime  `json:"startedAt"`
	EndedAt                      util.LocalTime  `json:"endedAt"`
	Frequency                    string          `json:"frequency"`
	Day                          string          `json:"day"`
	ExcludedDates                []string        `json:"excludedDates"`
	ExtraDates                   []string        `json:"extraDates"`
	ClassEventType               ClassEventType  `json:"classEventType"`
	ClassEvents                  []string        `json:"classEvents"`
	Processing                   bool            `json:"processing"`
	DeletedAt                    *util.LocalTime `json:"deletedAt"`
	DeletedBy                    *string         `json:"deletedBy"`
	Action                       interface{}     `json:"action"`
	CourseWaitingListRecurrences []interface{}   `json:"courseWaitingListRecurrences"`
}

type ClassEventType struct {
	ID                  *string        `json:"@id"`
	Type                string         `json:"@type"`
	StartedAt           util.LocalTime `json:"startedAt"`
	EndedAt             util.LocalTime `json:"endedAt"`
	Summary             string         `json:"summary"`
	Description         *string        `json:"description"`
	Club                *string        `json:"club"`
	Studio              *string        `json:"studio"`
	Activity            *string        `json:"activity"`
	Coach               *string        `json:"coach"`
	AttendingLimit      int            `json:"attendingLimit"`
	QueueLimit          int            `json:"queueLimit"`
	PrivateComment      interface{}    `json:"privateComment"`
	ClassLayout         interface{}    `json:"classLayout"`
	InstructionsComment interface{}    `json:"instructionsComment"`
	OnlineLimit         interface{}    `json:"onlineLimit"`
	ExternalQuota       interface{}    `json:"externalQuota"`
}

// Función genérica para extraer IDs de strings (no punteros)
func ExtractIDFromString(field string, errMsg string) (string, error) {
	if field == "" {
		return "", errors.New(errMsg)
	}
	cleanPath := strings.Split(field, "?")[0]
	base := path.Base(cleanPath)
	return base, nil
}

// Métodos para Recurrence
func (r *XPlorRecurrence) RecurrenceID() (string, error) {
	return ExtractID(r.ID, "recurrence ID field is nil")
}

func (r *XPlorRecurrence) ClassEventIDs() ([]string, error) {
	if len(r.ClassEvents) == 0 {
		return nil, errors.New("no class events available")
	}

	ids := make([]string, len(r.ClassEvents))
	for i, classEvent := range r.ClassEvents {
		id, err := ExtractIDFromString(classEvent, "class event ID field is empty")
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

// Métodos para ClassEventType
func (cet *ClassEventType) ClassEventTypeID() (string, error) {
	return ExtractID(cet.ID, "class event type ID field is nil")
}

func (cet *ClassEventType) ClubID() (string, error) {
	return ExtractID(cet.Club, "club ID field is nil")
}

func (cet *ClassEventType) StudioID() (string, error) {
	return ExtractID(cet.Studio, "studio ID field is nil")
}
func (cet *ClassEventType) CoachId() (string, error) {
	return ExtractID(cet.Coach, "coach ID field is nil")
}

func (cet *ClassEventType) ActivityID() (string, error) {
	return ExtractID(cet.Activity, "activity ID field is nil")
}

func (cet *ClassEventType) StartTime() string {
	return cet.StartedAt.Time.Format("15:04")
}

func (cet *ClassEventType) EndTime() string {
	return cet.EndedAt.Time.Format("15:04")
}

// Método para obtener el tipo de frecuencia
func (r *XPlorRecurrence) GetFrequencyType() string {
	freq := strings.ToLower(strings.TrimSpace(r.Frequency))
	switch freq {
	case "daily", "diario":
		return "daily"
	case "weekly", "semanal":
		return "weekly"
	case "monthly", "mensual":
		return "monthly"
	case "yearly", "anual":
		return "yearly"
	case "biweekly", "bisemanal":
		return "biweekly"
	case "bimonthly", "bimensual":
		return "bimonthly"

	default:
		return "unknown"
	}
}

// Método para obtener el día completo de la semana
func (r *XPlorRecurrence) GetFullDay() string {
	day := strings.ToLower(strings.TrimSpace(r.Day))
	switch day {
	// Inglés diminutivo
	case "mo":
		return "Monday"
	case "tu":
		return "Tuesday"
	case "we":
		return "Wednesday"
	case "th":
		return "Thursday"
	case "fr":
		return "Friday"
	case "sa":
		return "Saturday"
	case "su":
		return "Sunday"
	// Español diminutivo
	case "lu":
		return "Monday"
	case "ma":
		return "Tuesday"
	case "mi":
		return "Wednesday"
	case "ju":
		return "Thursday"
	case "vi":
		return "Friday"
	case "do":
		return "Sunday"
	// Si ya viene completo, devolverlo
	case "monday", "lunes":
		return "Monday"
	case "tuesday", "martes":
		return "Tuesday"
	case "wednesday", "miércoles":
		return "Wednesday"
	case "thursday", "jueves":
		return "Thursday"
	case "friday", "viernes":
		return "Friday"
	case "saturday", "sábado":
		return "Saturday"
	case "sunday", "domingo":
		return "Sunday"
	default:
		return "Unknown"
	}
}

// GetWeekdayNumber returns the weekday as an integer (0 = Sunday, 1 = Monday, ... 6 = Saturday)
func (r *XPlorRecurrence) GetWeekdayNumber() int {
	day := strings.ToLower(strings.TrimSpace(r.Day))
	switch day {
	// Inglés diminutivo
	case "mo":
		return 1
	case "tu":
		return 2
	case "we":
		return 3
	case "th":
		return 4
	case "fr":
		return 5
	case "sa":
		return 6
	case "su":
		return 0
	// Español diminutivo
	case "lu":
		return 1
	case "ma":
		return 2
	case "mi":
		return 3
	case "ju":
		return 4
	case "vi":
		return 5
	case "sá":
		return 6
	case "do":
		return 0
	// Si ya viene completo
	case "monday", "lunes":
		return 1
	case "tuesday", "martes":
		return 2
	case "wednesday", "miércoles":
		return 3
	case "thursday", "jueves":
		return 4
	case "friday", "viernes":
		return 5
	case "saturday", "sábado":
		return 6
	case "sunday", "domingo":
		return 0
	default:
		return -1
	}
}

// Métodos para la colección completa
func (rc *XPlorRecurrences) CollectionID() (string, error) {
	return ExtractIDFromString(rc.ID, "collection ID field is empty")
}

func (rc *XPlorRecurrences) ContextID() (string, error) {
	return ExtractIDFromString(rc.Context, "context ID field is empty")
}
func (rc *XPlorRecurrence) IsActive() bool {
	return rc.DeletedAt == nil && rc.EndedAt.Time.After(time.Now())
}

// Método para obtener todos los IDs de recurrencia en la colección
func (rc *XPlorRecurrences) AllRecurrenceIDs() ([]string, error) {
	if len(rc.Recurrences) == 0 {
		return nil, errors.New("no recurrences available")
	}

	ids := make([]string, len(rc.Recurrences))
	for i, recurrence := range rc.Recurrences {
		id, err := recurrence.RecurrenceID()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

// Método para obtener todos los IDs de actividades en la colección
func (rc *XPlorRecurrences) AllActivityIDs() ([]string, error) {
	if len(rc.Recurrences) == 0 {
		return nil, errors.New("no recurrences available")
	}

	activityIDs := make([]string, 0)
	for _, recurrence := range rc.Recurrences {
		activityID, err := recurrence.ClassEventType.ActivityID()
		if err == nil { // Solo agregar si no hay error
			activityIDs = append(activityIDs, activityID)
		}
	}

	if len(activityIDs) == 0 {
		return nil, errors.New("no activity IDs found")
	}

	return activityIDs, nil
}

// Método para obtener todos los IDs de estudios en la colección
func (rc *XPlorRecurrences) AllStudioIDs() ([]string, error) {
	if len(rc.Recurrences) == 0 {
		return nil, errors.New("no recurrences available")
	}

	studioIDs := make([]string, 0)
	for _, recurrence := range rc.Recurrences {
		studioID, err := recurrence.ClassEventType.StudioID()
		if err == nil { // Solo agregar si no hay error
			studioIDs = append(studioIDs, studioID)
		}
	}

	if len(studioIDs) == 0 {
		return nil, errors.New("no studio IDs found")
	}

	return studioIDs, nil
}

// Método para obtener todos los IDs de clubs en la colección
func (rc *XPlorRecurrences) AllClubIDs() ([]string, error) {
	if len(rc.Recurrences) == 0 {
		return nil, errors.New("no recurrences available")
	}

	clubIDs := make([]string, 0)
	for _, recurrence := range rc.Recurrences {
		clubID, err := recurrence.ClassEventType.ClubID()
		if err == nil { // Solo agregar si no hay error
			clubIDs = append(clubIDs, clubID)
		}
	}

	if len(clubIDs) == 0 {
		return nil, errors.New("no club IDs found")
	}

	return clubIDs, nil
}

// XPlorRecurrencesParams represents the search parameters for recurrences
type XPlorRecurrencesParams struct {
	StartedAtBefore         *time.Time
	StartedAtStrictlyBefore *time.Time
	StartedAtAfter          *time.Time
	StartedAtStrictlyAfter  *time.Time
	EndedAtBefore           *time.Time
	EndedAtStrictlyBefore   *time.Time
	EndedAtAfter            *time.Time
	EndedAtStrictlyAfter    *time.Time
	Week                    string
	Club                    string
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorRecurrencesParams) ToValues(values *url.Values) {
	// StartedAt date filters
	if p.StartedAtBefore != nil {
		values.Set("startedAt[before]", p.StartedAtBefore.Format("2006-01-02T15:04:05"))
	}
	if p.StartedAtStrictlyBefore != nil {
		values.Set("startedAt[strictly_before]", p.StartedAtStrictlyBefore.Format("2006-01-02T15:04:05"))
	}
	if p.StartedAtAfter != nil {
		values.Set("startedAt[after]", p.StartedAtAfter.Format("2006-01-02T15:04:05"))
	}
	if p.StartedAtStrictlyAfter != nil {
		values.Set("startedAt[strictly_after]", p.StartedAtStrictlyAfter.Format("2006-01-02T15:04:05"))
	}

	// EndedAt date filters
	if p.EndedAtBefore != nil {
		values.Set("endedAt[before]", p.EndedAtBefore.Format("2006-01-02T15:04:05"))
	}
	if p.EndedAtStrictlyBefore != nil {
		values.Set("endedAt[strictly_before]", p.EndedAtStrictlyBefore.Format("2006-01-02T15:04:05"))
	}
	if p.EndedAtAfter != nil {
		values.Set("endedAt[after]", p.EndedAtAfter.Format("2006-01-02T15:04:05"))
	}
	if p.EndedAtStrictlyAfter != nil {
		values.Set("endedAt[strictly_after]", p.EndedAtStrictlyAfter.Format("2006-01-02T15:04:05"))
	}
	// Club filter
	if p.Club != "" {
		values.Set("club", p.Club)
	}

	// Week filter
	if p.Week != "" {
		values.Set("week", p.Week)
	}
}
