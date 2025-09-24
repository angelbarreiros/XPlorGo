package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) subscriptions(accesToken string, params *xplorentities.XPlorSubscriptionsParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorSubscriptions, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorSubscriptions], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		if params != nil {
			params.ToValues(&queryParams)
		}
		formData := url.Values{}
		formData.Set("client_id", xe.config.ClientID)
		formData.Set("client_secret", xe.config.ClientSecret)
		var request = xe.config.generateRequest(http.MethodGet, "/subscriptions", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorSubscriptions](ctxWithTimeout, xe.client, request)
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
func (xe xplorExecutor) subscription(accesToken string, subscriptionId string) (*xplorentities.XPlorSubscription, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorSubscription], 1)

	go func() {
		formData := url.Values{}
		formData.Set("client_id", xe.config.ClientID)
		formData.Set("client_secret", xe.config.ClientSecret)
		var request = xe.config.generateRequest(http.MethodGet, "/subscriptions/"+subscriptionId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorSubscription](ctxWithTimeout, xe.client, request)
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
