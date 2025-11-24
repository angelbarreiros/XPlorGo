package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) families(accesToken string, params *xplorentities.XPlorFamiliesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorFamilies, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorFamilies], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		if params != nil {
			params.ToValues(&queryParams)
		}
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/families", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorFamilies](ctxWithTimeout, xe.client, request)
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
func (xe xplorExecutor) family(accesToken string, familyId string) (*xplorentities.XPlorFamily, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorFamily], 1)

	go func() {
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/families/"+familyId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorFamily](ctxWithTimeout, xe.client, request)
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
