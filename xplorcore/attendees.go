package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) attendees(accesToken string, classId *string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorAttendees, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorAttendees], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		util.AddQueryParam("class_id", classId, &queryParams)
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/attendees", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorAttendees](ctxWithTimeout, xe.client, request)
		resultChan <- result

	}()
	select {
	case res := <-resultChan:
		if res.Error == nil {
			return res.Response, res.Error
		}
		return nil, res.Error
	case <-ctxWithTimeout.Done():
		return nil, &xplorentities.ErrorResponse{
			Code:    http.StatusRequestTimeout,
			Message: "Request timeout: operation cancelled after 10 seconds",
		}
	}

}
