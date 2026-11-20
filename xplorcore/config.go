package xplorcore

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
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
	Debug          bool
	// RequestTimeout is the total time budget of each SDK method, including
	// retries. Defaults to 30s.
	RequestTimeout time.Duration
	// Retry enables automatic retries on 5xx / network errors. nil (the
	// default) means no retries; set it with WithRetry.
	Retry *RetryPolicy
}

// WithRequestTimeout sets the total time budget of each SDK method call,
// retries included. Must be called before Init/NewProvider.
func (xc *xplorConfig) WithRequestTimeout(timeout time.Duration) *xplorConfig {
	xc.RequestTimeout = timeout
	return xc
}

// WithRetry enables automatic retries of 5xx / network errors with the given
// policy (RetryPolicy{} uses the defaults). Retries are disabled unless this
// is called. Must be called before Init/NewProvider.
func (xc *xplorConfig) WithRetry(policy RetryPolicy) *xplorConfig {
	xc.Retry = &policy
	return xc
}

func NewConfig(host string, apiVersion string, enterpriseName, clientID, clientSecret string, headers map[string]string, debug bool) *xplorConfig {
	var config = &xplorConfig{
		Host:           host,
		APIVersion:     apiVersion,
		EnterpriseName: enterpriseName,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Debug:          debug,
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
			RawQuery: encodeQueryParams(queryParams),
		},
		Header: make(http.Header),
	}
	if params != nil {
		body := []byte(params.Encode())
		request.Body = io.NopCloser(bytes.NewReader(body))
		request.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(body)), nil
		}
		request.ContentLength = int64(len(body))
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

// encodeQueryParams keeps contact filters first so request logs and gateway
// filters can reliably identify the attendee scope before pagination fields.
func encodeQueryParams(values url.Values) string {
	if len(values) == 0 {
		return ""
	}

	preferredKeys := []string{
		"contact",
		"contact[]",
		"contactId",
		"contactId[]",
		"sale.contactId",
		"sale.contactId[]",
	}
	preferred := make(map[string]bool, len(preferredKeys))
	parts := make([]string, 0, len(values))
	for _, key := range preferredKeys {
		preferred[key] = true
		for _, value := range values[key] {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}

	remainingKeys := make([]string, 0, len(values))
	for key := range values {
		if !preferred[key] {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Strings(remainingKeys)
	for _, key := range remainingKeys {
		for _, value := range values[key] {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}

	return strings.Join(parts, "&")
}
