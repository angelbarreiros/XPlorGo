package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) studios(accesToken string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorStudios, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorStudios], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/studios", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorStudios](ctxWithTimeout, xe.client, request)
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
func (xe xplorExecutor) studio(accesToken string, familyId string) (*xplorentities.XPlorStudio, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorStudio], 1)

	go func() {
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/studios/"+familyId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorStudio](ctxWithTimeout, xe.client, request)
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
