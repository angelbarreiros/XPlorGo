package xplorcore

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) courseWaitingListRecurrences(accessToken string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorCourseWaitingListRecurrences, *xplorentities.ErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()

	queryParams := xplorentities.BuildPaginationQueryParams(pagination)
	request := xe.config.generateRequest(
		http.MethodGet,
		"/course_waiting_list_recurrences",
		xe.generateHeaders(accessToken),
		queryParams,
		url.Values{},
	).WithContext(ctx)
	result := util.ExecuteRequest[*xplorentities.XPlorCourseWaitingListRecurrences](ctx, xe.client, request, xe.config.Debug)
	return result.Response, result.Error
}

func (xe xplorExecutor) courseWaitingListRecurrence(accessToken string, resourceID string) (*xplorentities.XPlorCourseWaitingListRecurrence, *xplorentities.ErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()

	request := xe.config.generateRequest(
		http.MethodGet,
		"/course_waiting_list_recurrences/"+resourceID,
		xe.generateHeaders(accessToken),
		nil,
		url.Values{},
	).WithContext(ctx)
	result := util.ExecuteRequest[*xplorentities.XPlorCourseWaitingListRecurrence](ctx, xe.client, request, xe.config.Debug)
	return result.Response, result.Error
}

// CourseWaitingListRecurrence retrieves a course-level recurrence waiting list.
func (xe *XplorProvider) CourseWaitingListRecurrence(nodeID string, resourceID string) (*xplorentities.XPlorCourseWaitingListRecurrence, *xplorentities.ErrorResponse) {
	if strings.TrimSpace(resourceID) == "" {
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Course waiting-list recurrence ID is required",
		}
	}

	executor, apiErr := xe.getExecutorFullyInitialized(nodeID)
	if apiErr != nil {
		return nil, apiErr
	}
	defer xe.putExecutor(executor)

	resource, apiErr := executor.courseWaitingListRecurrence(xe.token.Token.AccessToken, resourceID)
	if apiErr != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    apiErr.Code,
			Message: "Failed to get course waiting-list recurrence: " + apiErr.Message,
		}
	}
	return resource, nil
}

// CourseWaitingListRecurrences retrieves course-level recurrence waiting lists.
func (xe *XplorProvider) CourseWaitingListRecurrences(nodeID string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorCourseWaitingListRecurrences, *xplorentities.ErrorResponse) {
	executor, apiErr := xe.getExecutorFullyInitialized(nodeID)
	if apiErr != nil {
		return nil, apiErr
	}
	defer xe.putExecutor(executor)

	resources, apiErr := executor.courseWaitingListRecurrences(xe.token.Token.AccessToken, pagination)
	if apiErr != nil {
		return nil, &xplorentities.ErrorResponse{
			Code:    apiErr.Code,
			Message: "Failed to get course waiting-list recurrences: " + apiErr.Message,
		}
	}
	return resources, nil
}
