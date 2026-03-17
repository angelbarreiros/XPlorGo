package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LocalTime handles API datetime values and normalizes them to UTC.
// Convention used by this SDK:
// - RFC3339 values (with timezone) keep their absolute instant.
// - Naive values (without timezone) are assumed to come in UTC+1.
type LocalTime struct {
	time.Time
}

// assumedSourceLocation defines the source timezone for naive timestamps.
var assumedSourceLocation = time.FixedZone("UTC+1", 1*60*60)

// parseTimeAsUTC parses API date/datetime strings and always returns UTC.
func parseTimeAsUTC(s string, dateTimeLayouts []string, dateLayouts []string) (time.Time, error) {
	// If payload includes timezone info, keep that absolute instant and normalize to UTC.
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}

	var lastErr error
	for _, layout := range dateTimeLayouts {
		// Naive datetime is interpreted in UTC+1 before UTC normalization.
		t, err := time.ParseInLocation(layout, s, assumedSourceLocation)
		if err == nil {
			return t.UTC(), nil
		}
		lastErr = err
	}

	for _, layout := range dateLayouts {
		// Naive date is interpreted at 00:00:00 in UTC+1, then converted to UTC.
		t, err := time.ParseInLocation(layout, s, assumedSourceLocation)
		if err == nil {
			return t.UTC(), nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

func (lt *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	// Parse and normalize to UTC to keep SDK timestamps consistent in storage.
	t, err := parseTimeAsUTC(
		s,
		[]string{"2006-01-02T15:04:05", "2006-01-02 15:04:05"},
		nil,
	)
	if err != nil {
		return err
	}
	lt.Time = t
	return nil
}

func (lt LocalTime) MarshalJSON() ([]byte, error) {
	// Keep API-compatible datetime format without timezone suffix.
	return json.Marshal(lt.Format("2006-01-02T15:04:05"))
}

// LocalDate handles API date values and normalizes them to UTC.
// Naive inputs are assumed to be in UTC+1.
type LocalDate struct {
	time.Time
}

func (ld *LocalDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := parseTimeAsUTC(
		s,
		[]string{"2006-01-02T15:04:05", "2006-01-02 15:04:05"},
		[]string{"2006-01-02"},
	)
	if err != nil {
		return fmt.Errorf("unable to parse date: %s", s)
	}
	ld.Time = t
	return nil
}

func (ld LocalDate) MarshalJSON() ([]byte, error) {
	// Keep API-compatible date format.
	return json.Marshal(ld.Format("2006-01-02"))
}

// RequestResult encapsulates the possible outcomes of an API request
type RequestResult[T any] struct {
	Response T
	Error    *ErrorResponse
}

// ExecuteRequest handles common HTTP request execution pattern including error handling and response processing
// It takes a context, http client, request, and returns a typed RequestResult
func ExecuteRequest[T any](ctx context.Context, client *http.Client, request *http.Request) RequestResult[T] {
	var zero T
	// log.Println(request.URL)
	// Log request details
	// fmt.Printf("%s %s\n", request.Method, request.URL.String())
	// for key, values := range request.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("  -H '%s: %s'\n", key, value)
	// 	}
	// }

	response, clientErr := client.Do(request)
	if clientErr != nil {
		return RequestResult[T]{
			Response: zero,
			Error: &ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to execute request: " + clientErr.Error(),
			},
		}
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	// log.Println(string(bodyBytes))
	if err != nil {
		return RequestResult[T]{
			Response: zero,
			Error: &ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to read response body: " + err.Error(),
			},
		}
	}

	// If we received a non-success status code, return an error
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return RequestResult[T]{
			Response: zero,
			Error: &ErrorResponse{
				Code:    response.StatusCode,
				Message: "Response: " + string(bodyBytes),
			},
		}
	}

	// Only try to unmarshal if we have response body
	if len(bodyBytes) > 0 {
		var target T
		err = json.Unmarshal(bodyBytes, &target)
		if err != nil {
			return RequestResult[T]{
				Response: zero,
				Error: &ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to unmarshal response: " + err.Error(),
				},
			}
		}
		return RequestResult[T]{
			Response: target,
			Error:    nil,
		}
	}

	return RequestResult[T]{
		Response: zero,
		Error:    nil,
	}
}

func AddQueryParam(name string, value *string, values *url.Values) {
	if value != nil && *value != "" {
		values.Add(name, *value)
	}

}
