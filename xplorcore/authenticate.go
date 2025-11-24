package xplorcore

import (
	"context"
	"net/http"
	"net/url"

	"github.com/angelbarreiros/XPlorGo/util"
	xplorentities "github.com/angelbarreiros/XPlorGo/xplorentities"
)

func (xe xplorExecutor) authenticate() (*xplorentities.XPlorTokenResponse, *xplorentities.ErrorResponse) {
	var ctxWithTimeout, cancel = context.WithTimeout(context.Background(), xe.defaultTimeout)
	defer cancel()
	resultChan := make(chan util.RequestResult[*xplorentities.XPlorTokenResponse], 1)

	go func() {
		var header = map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}
		formData := url.Values{}
		formData.Set("grant_type", "client_credentials")
		formData.Set("client_id", xe.config.ClientID)
		formData.Set("client_secret", xe.config.ClientSecret)

		var request = xe.config.generateRequest(http.MethodPost, "/oauth/v2/token", header, nil, formData)
		request = request.WithContext(ctxWithTimeout)
		result := util.ExecuteRequest[*xplorentities.XPlorTokenResponse](ctxWithTimeout, xe.client, request)
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
