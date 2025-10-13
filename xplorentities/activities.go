package xplorentities

import (
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
	ActivityGroups     []string           `json:"activityGroups"`
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
