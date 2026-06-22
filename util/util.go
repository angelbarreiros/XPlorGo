package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LocalTime handles API datetime values as naive timestamps (without timezone).
// All timestamps are returned as-is from the API, allowing timezone handling
// to be managed at the application level.
type LocalTime struct {
	time.Time
}

// parseTimeAsNaive parses API date/datetime strings and returns naive timestamps.
// RFC3339 values with timezone info have the timezone stripped, keeping only the date/time components.
// Naive values are parsed directly without timezone assumptions.
func parseTimeAsNaive(s string, dateTimeLayouts []string, dateLayouts []string) (time.Time, error) {
	// If payload includes timezone info (RFC3339), strip it and use only date/time components
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		// Return as naive: extract the date/time components without timezone
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC), nil
	}

	var lastErr error
	for _, layout := range dateTimeLayouts {
		// Parse naive datetime directly in UTC location (no conversion)
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	for _, layout := range dateLayouts {
		// Parse naive date directly in UTC location (no conversion)
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
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
	// Parse as naive timestamp without timezone conversion
	t, err := parseTimeAsNaive(
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
	t, err := parseTimeAsNaive(
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
// It takes a context, http client, request, debug flag, and returns a typed RequestResult
func ExecuteRequest[T any](ctx context.Context, client *http.Client, request *http.Request, debug bool) RequestResult[T] {
	var zero T
	if debug {
		curlCommand, curlErr := formatCurlCommand(request)
		if curlErr != nil {
			log.Printf("Failed to format cURL command: %v", curlErr)
		} else {
			fmt.Println(curlCommand)
		}
	}

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

	if debug {
		fmt.Printf("Response status: %s\n", response.Status)
		for _, key := range sortedHeaderKeys(response.Header) {
			for _, value := range response.Header.Values(key) {
				fmt.Printf("Response header %s: %s\n", key, value)
			}
		}
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return RequestResult[T]{
			Response: zero,
			Error: &ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to read response body: " + err.Error(),
			},
		}
	}
	if debug {
		fmt.Printf("Body: %s\n", string(bodyBytes))
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

func formatCurlCommand(request *http.Request) (string, error) {
	var builder strings.Builder
	builder.WriteString("curl -X ")
	builder.WriteString(request.Method)
	builder.WriteString(" \\\n'")
	builder.WriteString(shellSingleQuote(request.URL.String()))
	builder.WriteString("'")

	for _, key := range sortedHeaderKeys(request.Header) {
		for _, value := range request.Header.Values(key) {
			builder.WriteString(" \\\n-H '")
			builder.WriteString(shellSingleQuote(key))
			builder.WriteString(": ")
			builder.WriteString(shellSingleQuote(value))
			builder.WriteString("'")
		}
	}

	if request.Body != nil {
		bodyBytes, err := io.ReadAll(request.Body)
		if err != nil {
			return "", err
		}
		request.Body.Close()
		request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		if len(bodyBytes) > 0 {
			builder.WriteString(" \\\n-d '")
			builder.WriteString(shellSingleQuote(string(bodyBytes)))
			builder.WriteString("'")
		}
	}

	return builder.String(), nil
}

func sortedHeaderKeys(header http.Header) []string {
	preferred := []string{"Authorization", "Cache-Control", "Content-Type", "Accept", "User-Agent", "Timestamp"}
	keys := make([]string, 0, len(header))
	seen := make(map[string]bool, len(header))

	for _, key := range preferred {
		if _, ok := header[key]; ok {
			keys = append(keys, key)
			seen[key] = true
		}
	}

	var remaining []string
	for key := range header {
		if !seen[key] {
			remaining = append(remaining, key)
		}
	}
	sort.Strings(remaining)

	return append(keys, remaining...)
}

func shellSingleQuote(value string) string {
	return strings.ReplaceAll(value, "'", "'\"'\"'")
}
