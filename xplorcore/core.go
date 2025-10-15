package xplorcore

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

var xplorProviderInstace *XplorProvider = nil
var syncOnce sync.Once

type XplorProvider struct {
	providers *sync.Pool
	token     *xplorentities.XPlorTokenWithTimestamp
	authMutex *sync.Mutex
}
type xplorExecutor struct {
	config         *xplorConfig
	client         *http.Client
	defaultTimeout time.Duration
	nodeId         *string
}

func Init(cfg *xplorConfig) *XplorProvider {
	syncOnce.Do(func() {
		xplorProviderInstace = &XplorProvider{
			authMutex: &sync.Mutex{},
			providers: &sync.Pool{
				New: func() any {
					return &xplorExecutor{config: cfg, client: http.DefaultClient, defaultTimeout: 30 * time.Second}
				},
			},
		}
	})
	return xplorProviderInstace
}
func checkNodeId(nodeId string) *xplorentities.ErrorResponse {
	if strings.TrimSpace(nodeId) == "" {
		return &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Node ID is required",
		}
	}
	return nil
}

func (pp XplorProvider) getExecutor(nodeId string) *xplorExecutor {
	var executor = pp.providers.Get().(*xplorExecutor)
	if strings.TrimSpace(nodeId) == "" {
		executor.nodeId = nil
	} else {

		executor.nodeId = &nodeId
	}
	return executor
}
func (pp XplorProvider) putExecutor(executor *xplorExecutor) {
	pp.providers.Put(executor)
}
func (pp XplorProvider) Close() {
	xplorProviderInstace = nil
}
func (xe *xplorExecutor) generateHeaders(accessToken string) map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + accessToken,
	}
	if xe.nodeId != nil && strings.TrimSpace(*xe.nodeId) != "" {
		headers["X-User-Network-Node-Id"] = "/" + xe.config.EnterpriseName + "/network_nodes/" + *xe.nodeId
	}
	return headers
}
func (xp XplorProvider) needsAuthentication(token *xplorentities.XPlorTokenWithTimestamp) bool {
	if token == nil {
		return true
	}

	return !token.IsValid()
}
func (xe *XplorProvider) authenticateIfNeeded(executor *xplorExecutor) *xplorentities.ErrorResponse {

	// Double-check locking pattern
	xe.authMutex.Lock()
	defer xe.authMutex.Unlock()

	// Second check with lock to prevent race conditions
	if xe.needsAuthentication(xe.token) {
		var token *xplorentities.XPlorTokenResponse
		var err *xplorentities.ErrorResponse
		if token, err = executor.authenticate(); err != nil {
			return &xplorentities.ErrorResponse{
				Code:    err.Code,
				Message: "Failed to authenticate: " + err.Message,
			}
		}
		xe.token = &xplorentities.XPlorTokenWithTimestamp{
			Token:      token,
			ObtainedAt: time.Now(),
		}
	}
	return nil

}
func (xe *XplorProvider) Families(nodeId string, params *xplorentities.XPlorFamiliesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorFamilies, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	families, err := executor.families(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get families: " + err.Message,
		}
	}

	return families, nil

}
func (xe *XplorProvider) Family(nodeId string, familyId string) (*xplorentities.XPlorFamily, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(familyId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Family ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	family, err := executor.family(xe.token.Token.AccessToken, familyId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get family: " + err.Message,
		}
	}

	return family, nil

}
func (xe *XplorProvider) Clubs(nodeId string) (*xplorentities.XPloreClubs, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}

	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}

	clubs, err := executor.clubs(xe.token.Token.AccessToken, nil)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get clubs: " + err.Message,
		}
	}

	return clubs, nil

}
func (xe *XplorProvider) Club(nodeId string, clubId string) (*xplorentities.XPlorClub, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(clubId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Club ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	club, err := executor.club(xe.token.Token.AccessToken, clubId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get club: " + err.Message,
		}
	}

	return club, nil

}
func (xe *XplorProvider) Events(nodeId string, pagination *xplorentities.XPlorPagination, timeGap *xplorentities.XPlorTimeGap) (*xplorentities.XPlorEvents, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	events, err := executor.events(xe.token.Token.AccessToken, pagination, timeGap)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get events: " + err.Message,
		}
	}

	return events, nil

}
func (xe *XplorProvider) Activities(nodeId string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorActivities, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	activities, err := executor.activities(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get activities: " + err.Message,
		}
	}

	return activities, nil
}
func (xe *XplorProvider) Activity(nodeId string, activityId string) (*xplorentities.XPlorActivity, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(activityId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Activity ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	activity, err := executor.activity(xe.token.Token.AccessToken, activityId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get activity: " + err.Message,
		}
	}

	return activity, nil

}
func (xd *XplorProvider) Studios(nodeId string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorStudios, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xd.getExecutor(nodeId)
	defer xd.putExecutor(executor)

	if err := xd.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	studios, err := executor.studios(xd.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get studios: " + err.Message,
		}
	}

	return studios, nil

}
func (xd *XplorProvider) Studio(nodeId string, studioId string) (*xplorentities.XPlorStudio, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(studioId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Studio ID is required",
		}
	}
	var executor = xd.getExecutor(nodeId)
	defer xd.putExecutor(executor)

	if err := xd.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	studio, err := executor.studio(xd.token.Token.AccessToken, studioId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get studio: " + err.Message,
		}
	}

	return studio, nil

}
func (xd *XplorProvider) Contacts(nodeId string, params *xplorentities.XPlorContactsParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorContacts, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xd.getExecutor(nodeId)
	defer xd.putExecutor(executor)

	if err := xd.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	contacts, err := executor.contacts(xd.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get contacts: " + err.Message,
		}
	}

	return contacts, nil

}
func (xd *XplorProvider) Contact(nodeId string, contactId string) (*xplorentities.XPlorContact, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(contactId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Contact ID is required",
		}
	}
	var executor = xd.getExecutor(nodeId)
	defer xd.putExecutor(executor)

	if err := xd.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	contact, err := executor.contact(xd.token.Token.AccessToken, contactId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get contact: " + err.Message,
		}
	}

	return contact, nil

}
func (xe *XplorProvider) Subscriptions(nodeId string, params *xplorentities.XPlorSubscriptionsParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorSubscriptions, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	subscriptions, err := executor.subscriptions(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get subscriptions: " + err.Message,
		}
	}

	return subscriptions, nil

}
func (xe *XplorProvider) Subscription(nodeId string, subscriptionId string) (*xplorentities.XPlorSubscription, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(subscriptionId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Subscription ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	subscription, err := executor.subscription(xe.token.Token.AccessToken, subscriptionId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get subscription: " + err.Message,
		}
	}

	return subscription, nil

}
func (xe *XplorProvider) Classes(nodeId string, params *xplorentities.XPlorClassesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorClasses, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	classes, err := executor.classes(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get classes: " + err.Message,
		}
	}

	return classes, nil

}
func (xe *XplorProvider) Class(nodeId string, classId string) (*xplorentities.XPlorClass, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(classId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Class ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	class, err := executor.class(xe.token.Token.AccessToken, classId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get class: " + err.Message,
		}
	}

	return class, nil

}
func (xe *XplorProvider) NetworkNodes(pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorNetworkNodes, *xplorentities.ErrorResponse) {
	var executor = xe.getExecutor("")
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	networkNodes, err := executor.networkNodes(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get nodes: " + err.Message,
		}
	}

	return networkNodes, nil
}
func (xe *XplorProvider) NetworkNode(nodeId string) (*xplorentities.XPlorNetworkNode, *xplorentities.ErrorResponse) {
	if strings.TrimSpace(nodeId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Node ID is required",
		}
	}
	var executor = xe.getExecutor("")
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	networkNode, err := executor.networkNode(xe.token.Token.AccessToken, nodeId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get node: " + err.Message,
		}
	}

	return networkNode, nil
}

func (xe *XplorProvider) Attendees(nodeId string, classId *string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorAttendees, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	attendees, err := executor.attendees(xe.token.Token.AccessToken, classId, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get attendees: " + err.Message,
		}
	}

	return attendees, nil

}
func (xe *XplorProvider) Coaches(nodeId string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPloreCoaches, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	coaches, err := executor.coaches(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get coaches: " + err.Message,
		}
	}

	return coaches, nil

}
func (xe *XplorProvider) Coach(nodeId string, coachId string) (*xplorentities.XPloreCoach, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(coachId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Coach ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	coach, err := executor.coach(xe.token.Token.AccessToken, coachId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get coach: " + err.Message,
		}
	}

	return coach, nil

}
func (xe *XplorProvider) Articles(nodeId string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorArticles, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	articles, err := executor.articles(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get articles: " + err.Message,
		}
	}

	return articles, nil

}
func (xe *XplorProvider) Article(nodeId string, articleId string) (*xplorentities.XPlorArticle, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(articleId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Article ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	article, err := executor.article(xe.token.Token.AccessToken, articleId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get article: " + err.Message,
		}
	}

	return article, nil

}
func (xe *XplorProvider) Recurrences(nodeId string, params *xplorentities.XPlorRecurrencesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorRecurrences, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	recurrences, err := executor.recurrences(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get recurrences: " + err.Message,
		}
	}

	return recurrences, nil

}
func (xe *XplorProvider) Recurrence(nodeId string, recurrenceId string) (*xplorentities.XPlorRecurrence, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(recurrenceId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Recurrence ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	recurrence, err := executor.frecurrence(xe.token.Token.AccessToken, recurrenceId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get recurrence: " + err.Message,
		}
	}

	return recurrence, nil

}

func (xe *XplorProvider) ClassType(nodeId string, classTypeId string) (*xplorentities.XPlorClassType, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(classTypeId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Class Type ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	classType, err := executor.classType(xe.token.Token.AccessToken, classTypeId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get class type: " + err.Message,
		}
	}

	return classType, nil
}

func (xe *XplorProvider) CounterLines(nodeId string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorCounterLines, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	counterLines, err := executor.counterLines(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get families: " + err.Message,
		}
	}

	return counterLines, nil

}
func (xe *XplorProvider) CounterLine(nodeId string, counterLineId string) (*xplorentities.XPlorCounterLine, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(counterLineId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Counter Line ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	counterLine, err := executor.counterLine(xe.token.Token.AccessToken, counterLineId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get family: " + err.Message,
		}
	}
	return counterLine, nil

}
func (xe *XplorProvider) ContactTags(nodeId string, params *xplorentities.XPlorContactTagsParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorContactTags, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	contactTags, err := executor.contactTags(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get contact tags: " + err.Message,
		}
	}

	return contactTags, nil

}
func (xe *XplorProvider) ContactTag(nodeId string, contactTagId string) (*xplorentities.XPlorContactTag, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(contactTagId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Contact Tag ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	contactTag, err := executor.contactTag(xe.token.Token.AccessToken, contactTagId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get contact tag: " + err.Message,
		}
	}

	return contactTag, nil

}
func (xe *XplorProvider) Users(pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorUsers, *xplorentities.ErrorResponse) {
	var executor = xe.getExecutor("")
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	users, err := executor.users(xe.token.Token.AccessToken, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get users: " + err.Message,
		}
	}

	return users, nil

}
func (xe *XplorProvider) User(userId string) (*xplorentities.XPlorUser, *xplorentities.ErrorResponse) {
	if strings.TrimSpace(userId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User ID is required",
		}
	}
	var executor = xe.getExecutor("")
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	user, err := executor.user(xe.token.Token.AccessToken, userId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get user: " + err.Message,
		}
	}

	return user, nil

}
func (xe *XplorProvider) Zones(nodeId string, params *xplorentities.XPlorZonesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorZones, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	zones, err := executor.zones(xe.token.Token.AccessToken, params, pagination)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get zones: " + err.Message,
		}
	}

	return zones, nil

}
func (xe *XplorProvider) Zone(nodeId string, zoneId string) (*xplorentities.XPlorZone, *xplorentities.ErrorResponse) {
	if err := checkNodeId(nodeId); err != nil {
		return nil, err
	}
	if strings.TrimSpace(zoneId) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Zone ID is required",
		}
	}
	var executor = xe.getExecutor(nodeId)
	defer xe.putExecutor(executor)

	if err := xe.authenticateIfNeeded(executor); err != nil {
		return nil, err
	}
	zone, err := executor.zone(xe.token.Token.AccessToken, zoneId)
	if err != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    err.Code,
			Message: "Failed to get zone: " + err.Message,
		}
	}

	return zone, nil

}
