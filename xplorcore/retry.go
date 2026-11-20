package xplorcore

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

// RetryPolicy controls how the SDK retries requests that fail because of a
// server-side or network problem (5xx responses or transport errors).
// Retries are opt-in: they are only active when the config is built with
// WithRetry. A zero-value RetryPolicy enables retries with the defaults.
// Client errors (4xx: validation, not found, forbidden...) are never retried;
// 401 is handled separately by authRetryTransport.
//
// Retries are bounded by the context of each SDK method: the SDK keeps
// retrying with exponential backoff while there is time left before the
// method's deadline (see WithRequestTimeout), so a call never takes longer
// than it would without retries.
type RetryPolicy struct {
	// MaxRetries caps the number of retries. 0 means "no cap": retry until
	// the context deadline is reached.
	MaxRetries int
	// BaseDelay is the wait before the first retry. It doubles on every
	// attempt (with jitter). Defaults to 250ms.
	BaseDelay time.Duration
	// MaxDelay caps a single wait between attempts. Defaults to 5s.
	MaxDelay time.Duration
}

const (
	defaultRetryBaseDelay = 250 * time.Millisecond
	defaultRetryMaxDelay  = 5 * time.Second
)

type retryTransport struct {
	base   http.RoundTripper
	policy RetryPolicy
	debug  bool
}

func (t *retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if !isIdempotent(request.Method) || !requestCanBeRetried(request) {
		return t.base.RoundTrip(request)
	}

	ctx := request.Context()
	for attempt := 0; ; attempt++ {
		attemptRequest := request
		if attempt > 0 {
			var err error
			if attemptRequest, err = cloneRequestForRetry(request); err != nil {
				return nil, err
			}
		}

		response, err := t.base.RoundTrip(attemptRequest)
		if !shouldRetry(ctx, response, err) {
			return response, err
		}
		if t.policy.MaxRetries > 0 && attempt >= t.policy.MaxRetries {
			return response, err
		}

		delay := t.backoff(attempt, response)
		// Give up early if waiting would exceed the method's deadline: the
		// caller then gets the real server error instead of a generic timeout.
		if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) <= delay {
			return response, err
		}

		if t.debug {
			log.Printf("xplor: retrying %s %s in %s (attempt %d): %s",
				request.Method, request.URL.Path, delay, attempt+1, retryReason(response, err))
		}
		discardBody(response)

		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		}
	}
}

// shouldRetry returns true only for transient failures: network errors and
// 5xx responses. Context cancellation/deadline is never retried.
func shouldRetry(ctx context.Context, response *http.Response, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	if err != nil {
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	return response != nil && response.StatusCode >= 500 && response.StatusCode <= 599
}

func (t *retryTransport) backoff(attempt int, response *http.Response) time.Duration {
	base := t.policy.BaseDelay
	if base <= 0 {
		base = defaultRetryBaseDelay
	}
	maxDelay := t.policy.MaxDelay
	if maxDelay <= 0 {
		maxDelay = defaultRetryMaxDelay
	}

	// Honor Retry-After (seconds) when the server sends it, e.g. on 503.
	if response != nil {
		if retryAfter, ok := parseRetryAfter(response.Header.Get("Retry-After")); ok {
			return min(retryAfter, maxDelay)
		}
	}

	delay := base << min(attempt, 20)
	if delay <= 0 || delay > maxDelay {
		delay = maxDelay
	}
	// Full jitter in [delay/2, delay) so concurrent callers don't retry in sync.
	half := delay / 2
	return half + rand.N(half+1)
}

func parseRetryAfter(value string) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}
	if when, err := http.ParseTime(value); err == nil {
		return max(time.Until(when), 0), true
	}
	var seconds int
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, false
		}
		seconds = seconds*10 + int(r-'0')
		if seconds > 3600 {
			return time.Hour, true
		}
	}
	return time.Duration(seconds) * time.Second, true
}

func isIdempotent(method string) bool {
	switch method {
	case "", http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete:
		return true
	}
	return false
}

func cloneRequestForRetry(request *http.Request) (*http.Request, error) {
	clone := request.Clone(request.Context())
	if request.Body != nil {
		body, err := request.GetBody()
		if err != nil {
			return nil, err
		}
		clone.Body = body
	}
	return clone, nil
}

// discardBody drains and closes a response we are not returning, so the
// underlying connection can be reused.
func discardBody(response *http.Response) {
	if response == nil || response.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
	_ = response.Body.Close()
}

func retryReason(response *http.Response, err error) string {
	if err != nil {
		return err.Error()
	}
	return response.Status
}
