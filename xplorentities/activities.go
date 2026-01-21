package xplorentities

import (
	"net/url"
	"path"
	"strconv"

	"github.com/angelbarreiros/XPlorGo/util"
)

// ----------- Hydra Collection de Activities -----------
type XPlorActivities struct {
	Context    string          `json:"@context"`
	ID         string          `json:"@id"`
	Type       string          `json:"@type"`
	Activities []XPlorActivity `json:"hydra:member"`
	Pagination HydraView       `json:"hydra:view"`
}

// ----------- XPlorActivity -----------
type XPlorActivity struct {
	ID                 *string            `json:"@id"`
	Type               string             `json:"@type"`
	Name               string             `json:"name"`
	ClubId             *string            `json:"clubId"`
	ColorHex           string             `json:"colorHex"`
	Durations          []string           `json:"durations"`
	CreatedAt          *util.LocalTime    `json:"createdAt"`
	CreatedBy          string             `json:"createdBy"`
	TemplateToken      *string            `json:"templateToken"`
	ShowcaseActivities []ShowcaseActivity `json:"showcaseActivities"`
	ActivityGroups     []ActivityGroup    `json:"activityGroups"`
	ArchivedAt         *util.LocalTime    `json:"archivedAt"`
	ArchivedBy         *string            `json:"archivedBy"`
	IsBookable         bool               `json:"isBookable"`
	IsViewable         bool               `json:"isViewable"`
	NetworkNodeName    *string            `json:"networkNodeName"`
}

// ----------- ShowcaseActivity -----------
type ShowcaseActivity struct {
	ID          string        `json:"@id"`
	Type        string        `json:"@type"`
	Description string        `json:"description"`
	IsBookable  bool          `json:"isBookable"`
	Properties  ShowcaseProps `json:"properties"`
	Context     string        `json:"context"`
}

type ShowcaseProps struct {
	ExtendedImage       *string `json:"extendedImage"`
	ShowAvailablePlaces bool    `json:"showAvailablePlaces"`
}

// ----------- ActivityGroup -----------
type ActivityGroup struct {
	ID              string   `json:"@id"`
	Type            string   `json:"@type"`
	Name            string   `json:"name"`
	Description     *string  `json:"description"`
	Type_           *string  `json:"type"`
	Activities      []string `json:"activities"`
	NetworkNodeName string   `json:"networkNodeName"`
}

// ----------- Métodos útiles -----------

func (a XPlorActivity) ActivityID() (string, error) {
	return ExtractID(a.ID, "activity ID field is nil")
}

func (a XPlorActivity) ClubID() (string, error) {
	return ExtractID(a.ClubId, "club ID field is nil")
}

// Obtener IDs de Showcase Activities asociadas
func (a XPlorActivity) ShowcaseIDs() ([]int, error) {
	var ids []int
	for _, s := range a.ShowcaseActivities {
		base := path.Base(s.ID)
		n, err := strconv.Atoi(base)
		if err != nil {
			return nil, err
		}
		ids = append(ids, n)
	}
	return ids, nil
}

// Chequear si la actividad es activa (no archivada)
func (a XPlorActivity) IsActive() bool {
	return a.ArchivedAt == nil
}

// Chequear si la actividad es de PADEL
func (a XPlorActivity) IsPadel() bool {
	return len(a.Name) >= 5 && a.Name[:5] == "PADEL"
}

// Obtener la duración en minutos (si viene en formato ISO 8601, ej: PT60M, PT90M)
func (a XPlorActivity) DurationMinutes() int {
	if len(a.Durations) == 0 {
		return 0
	}
	// asumimos formato PTxxM
	d := a.Durations[0]
	if len(d) > 2 && d[:2] == "PT" && d[len(d)-1] == 'M' {
		minStr := d[2 : len(d)-1]
		if m, err := strconv.Atoi(minStr); err == nil {
			return m
		}
	}
	return 0
}

// XPlorActivitiesParams represents the search parameters for activities
type XPlorActivitiesParams struct {
	ClubID   *string
	ClubIDs  []string
	Name     *string
	Archived *bool
}

// ToValues converts the params to url.Values for query parameters
func (p XPlorActivitiesParams) ToValues(orgName string, values *url.Values) {
	// Single value filters
	if p.ClubID != nil {
		values.Set("clubId", "/"+orgName+"/clubs/"+*p.ClubID)
	}
	if p.Name != nil {
		values.Set("name", *p.Name)
	}
	if p.Archived != nil && *p.Archived {
		values.Set("archived", "true")
	}

	// Array filters
	for _, clubID := range p.ClubIDs {
		values.Add("clubId[]", clubID)
	}
}
