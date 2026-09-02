package xplorcore

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) courseWaitingList(accessToken string, resourceID string) (*xplorentities.XPlorCourseWaitingList, *xplorentities.ErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()

	request := xe.config.generateRequest(
		http.MethodGet,
		"/course_waiting_lists/"+resourceID,
		xe.generateHeaders(accessToken),
		nil,
		url.Values{},
	).WithContext(ctx)
	result := util.ExecuteRequest[*xplorentities.XPlorCourseWaitingList](ctx, xe.client, request, xe.config.Debug)
	return result.Response, result.Error
}

// CourseWaitingList retrieves a course-level waiting list and its entries.
func (xe *XplorProvider) CourseWaitingList(nodeID string, resourceID string) (*xplorentities.XPlorCourseWaitingList, *xplorentities.ErrorResponse) {
	if strings.TrimSpace(resourceID) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Course waiting-list ID is required",
		}
	}

	executor, apiErr := xe.getExecutorFullyInitialized(nodeID)
	if apiErr != nil {
		return nil, apiErr
	}
	defer xe.putExecutor(executor)

	resource, apiErr := executor.courseWaitingList(xe.token.Token.AccessToken, resourceID)
	if apiErr != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    apiErr.Code,
			Message: "Failed to get course waiting list: " + apiErr.Message,
		}
	}
	return resource, nil
}
