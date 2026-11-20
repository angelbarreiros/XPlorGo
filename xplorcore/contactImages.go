package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) contactImages(accesToken string, params *xplorentities.XPlorContactImagesParams, pagination *xplorentities.XPlorPagination) (*xplorentities.XPlorContactImages, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorContactImages], 1)

	go func() {
		var queryParams = xplorentities.BuildPaginationQueryParams(pagination)
		if params != nil {
			params.ToValues(&queryParams)
		}
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/files/contact_images", xe.generateHeaders(accesToken), queryParams, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorContactImages](ctxWithTimeout, xe.client, request, xe.config.Debug)
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
			Message: "Request timeout: operation cancelled after " + xe.defaultTimeout.String(),
		}
	}
}

func (xe xplorExecutor) contactImage(accesToken string, contactImageId string) ([]byte, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[[]byte], 1)

	go func() {
		formData := url.Values{}

		var request = xe.config.generateRequest(http.MethodGet, "/files/contact_images/"+contactImageId, xe.generateHeaders(accesToken), nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteBinaryRequest(ctxWithTimeout, xe.client, request, xe.config.Debug)
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
			Message: "Request timeout: operation cancelled after " + xe.defaultTimeout.String(),
		}
	}
}
