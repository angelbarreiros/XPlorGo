package xplorcore

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type neededHeaders struct {
	HeaderName string
	Value      string
}

type xplorConfig struct {
	Host           string
	APIVersion     string
	EnterpriseName string
	ClientID       string
	ClientSecret   string
	NeededHeaders  []neededHeaders
}

func NewConfig(host string, apiVersion string, enterpriseName, clientID, clientSecret string, headers map[string]string) *xplorConfig {
	var config = &xplorConfig{
		Host:           host,
		APIVersion:     apiVersion,
		EnterpriseName: enterpriseName,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
	}
	for headerName, value := range headers {
		config.NeededHeaders = append(config.NeededHeaders, neededHeaders{
			HeaderName: headerName,
			Value:      value,
		})
	}
	return config
}

func (xc *xplorConfig) generateRequest(method string, uri string, optionalHeaders map[string]string, queryParams url.Values, params url.Values) *http.Request {
	request := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme:   "https",
			Host:     xc.Host,
			Path:     "/" + xc.APIVersion + "/" + xc.EnterpriseName + uri,
			RawQuery: queryParams.Encode(),
		},
		Header: make(http.Header),
	}
	if params != nil {
		request.Body = io.NopCloser(bytes.NewBufferString(params.Encode()))
	}

	for _, header := range xc.NeededHeaders {
		request.Header.Add(header.HeaderName, header.Value)
	}
	for headerName, value := range optionalHeaders {
		request.Header.Add(headerName, value)
	}
	// log.Println("Request URL:", request.URL.String())

	return request

}
