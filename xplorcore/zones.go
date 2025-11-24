package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) zones(accesToken string, params *xplorentities.XPlorZonesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorZones, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorZones], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		if params != nil {
			params.ToValues(&queryParams)
		}
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/zones", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorZones](ctxWithTimeout, xe.client, request)
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
func (xe xplorExecutor) zone(accesToken string, zoneId string) (*xplorentities.XPlorZone, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorZone], 1)

	go func() {
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/zones/"+zoneId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorZone](ctxWithTimeout, xe.client, request)
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
