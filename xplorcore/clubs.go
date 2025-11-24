package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) clubs(accesToken string, pagination *xplorentities.XPlorPagination) (*xplorentities.XPloreClubs, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPloreClubs], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/clubs", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPloreClubs](ctxWithTimeout, xe.client, request)
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

func (xe xplorExecutor) club(accesToken string, clubId string) (*xplorentities.XPlorClub, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorClub], 1)

	go func() {
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/clubs/"+clubId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorClub](ctxWithTimeout, xe.client, request)
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
