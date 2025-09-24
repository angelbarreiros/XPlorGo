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
	EnterpriseName string
	ClientID       string
	ClientSecret   string
	NeededHeaders  []neededHeaders
}

func NewConfig(host string, enterpriseName, clientID, clientSecret string, headers map[string]string) *xplorConfig {
	var config = &xplorConfig{
		EnterpriseName: enterpriseName,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Host:           host,
	}
	for headerName, value := range headers {
		config.NeededHeaders = append(config.NeededHeaders, neededHeaders{
			HeaderName: headerName,
			Value:      value,
		})
	}
	// readEnvFile()

	return config
}

func (xc *xplorConfig) generateRequest(method string, uri string, optionalHeaders map[string]string, queryParams url.Values, params url.Values) *http.Request {

	request := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme:   "https",
			Host:     xc.Host + uri,
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

// func readEnvFile() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		panic("Error loading .env file, must have ENVIRONMENT variable set")
// 	}
// }
