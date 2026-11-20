package xplorcore

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func fastPolicy() RetryPolicy {
	return RetryPolicy{BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
}

// serverReturning answers with the given status codes in order, then 200.
func serverReturning(t *testing.T, codes ...int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(calls.Add(1))
		if n <= len(codes) {
			w.WriteHeader(codes[n-1])
			_, _ = io.WriteString(w, "error")
			return
		}
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	t.Cleanup(server.Close)
	return server, &calls
}

func doGet(t *testing.T, transport http.RoundTripper, url string, timeout time.Duration) (*http.Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	return (&http.Client{Transport: transport}).Do(req)
}

func TestRetryRecoversFromServerErrors(t *testing.T) {
	server, calls := serverReturning(t, 500, 502, 503)
	transport := &retryTransport{base: http.DefaultTransport, policy: fastPolicy()}

	resp, err := doGet(t, transport, server.URL, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || calls.Load() != 4 {
		t.Fatalf("status=%d calls=%d, want 200 after 4 calls", resp.StatusCode, calls.Load())
	}
}

func TestRetryDoesNotRetryClientErrors(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404, 409, 422, 429} {
		server, calls := serverReturning(t, code)
		transport := &retryTransport{base: http.DefaultTransport, policy: fastPolicy()}

		resp, err := doGet(t, transport, server.URL, 5*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != code || calls.Load() != 1 {
			t.Fatalf("code %d: status=%d calls=%d, want no retry", code, resp.StatusCode, calls.Load())
		}
	}
}

func TestRetryStopsBeforeContextDeadlineAndReturnsServerError(t *testing.T) {
	server, calls := serverReturning(t, 500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
		500, 500, 500, 500, 500, 500, 500, 500, 500, 500)
	policy := RetryPolicy{BaseDelay: 40 * time.Millisecond, MaxDelay: 40 * time.Millisecond}
	transport := &retryTransport{base: http.DefaultTransport, policy: policy}

	start := time.Now()
	resp, err := doGet(t, transport, server.URL, 150*time.Millisecond)
	if err != nil {
		t.Fatalf("want the last 500 response, got error %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 500 || !strings.Contains(string(body), "error") {
		t.Fatalf("status=%d body=%q", resp.StatusCode, body)
	}
	if elapsed := time.Since(start); elapsed > 150*time.Millisecond {
		t.Fatalf("took %s, exceeded the context budget", elapsed)
	}
	if calls.Load() < 2 {
		t.Fatalf("calls=%d, expected some retries within the budget", calls.Load())
	}
}

func TestRetryRespectsMaxRetries(t *testing.T) {
	server, calls := serverReturning(t, 503, 503, 503, 503)
	policy := fastPolicy()
	policy.MaxRetries = 2
	resp, err := doGet(t, &retryTransport{base: http.DefaultTransport, policy: policy}, server.URL, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 503 || calls.Load() != 3 {
		t.Fatalf("status=%d calls=%d, want 503 after 3 calls", resp.StatusCode, calls.Load())
	}

}

func TestRetryIsOptIn(t *testing.T) {
	cfg := NewConfig("example.com", "v1", "ent", "id", "secret", nil, false)
	if _, ok := newXplorHTTPClient(nil, cfg).Transport.(*authRetryTransport).base.(*retryTransport); ok {
		t.Fatal("retries must be disabled by default")
	}
	cfg.WithRetry(RetryPolicy{})
	if _, ok := newXplorHTTPClient(nil, cfg).Transport.(*authRetryTransport).base.(*retryTransport); !ok {
		t.Fatal("WithRetry must enable retries")
	}
}

func TestRetryDoesNotRetryPost(t *testing.T) {
	server, calls := serverReturning(t, 500)
	transport := &retryTransport{base: http.DefaultTransport, policy: fastPolicy()}
	req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("a=b"))
	resp, err := (&http.Client{Transport: transport}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if calls.Load() != 1 {
		t.Fatalf("calls=%d, POST must not be retried", calls.Load())
	}
}

func TestRetryRecoversFromNetworkError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close() // abrupt connection drop
			return
		}
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	resp, err := doGet(t, &retryTransport{base: http.DefaultTransport, policy: fastPolicy()}, server.URL, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 || calls.Load() != 2 {
		t.Fatalf("status=%d calls=%d", resp.StatusCode, calls.Load())
	}
}

func TestParseRetryAfter(t *testing.T) {
	if d, ok := parseRetryAfter("3"); !ok || d != 3*time.Second {
		t.Fatalf("got %v %v", d, ok)
	}
	if _, ok := parseRetryAfter("soon"); ok {
		t.Fatal("expected invalid")
	}
}
