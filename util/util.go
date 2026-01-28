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

// LocalTime maneja fechas en formato "2006-01-02T15:04:05" sin zona horaria
type LocalTime struct {
	time.Time
}

func (lt *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	// Ajusta el layout al formato recibido por el API
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	lt.Time = t
	return nil
}

func (lt LocalTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(lt.Format("2006-01-02T15:04:05"))
}

// LocalDate maneja fechas en formato "2006-01-02" sin zona horaria ni hora
type LocalDate struct {
	time.Time
}

func (ld *LocalDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	// Intenta parsear primero formato completo con hora
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		ld.Time = t
		return nil
	}
	// Si falla, intenta parsear solo fecha
	if t, err := time.Parse("2006-01-02", s); err == nil {
		ld.Time = t
		return nil
	}
	return fmt.Errorf("unable to parse date: %s", s)
}

func (ld LocalDate) MarshalJSON() ([]byte, error) {
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
