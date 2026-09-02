package xplorcore

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/angelbarreiros/XPlorGo/xplorentities"
)

// authRetryTransport retries an API request once when XPlor rejects its
// access token. The retry is deliberately implemented at the transport layer
// so every SDK endpoint gets the same behavior, including future endpoints.
type authRetryTransport struct {
	provider *XplorProvider
	base     http.RoundTripper
}

func newXplorHTTPClient(provider *XplorProvider) *http.Client {
	base := http.DefaultClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	return &http.Client{
		Transport:     &authRetryTransport{provider: provider, base: base},
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           http.DefaultClient.Jar,
		Timeout:       http.DefaultClient.Timeout,
	}
}

func (t *authRetryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.base.RoundTrip(request)
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized {
		return response, err
	}

	staleToken, hasBearerToken := bearerToken(request.Header.Get("Authorization"))
	if isAuthenticationRequest(request) || !hasBearerToken || !requestCanBeRetried(request) || request.Context().Err() != nil {
		return response, nil
	}

	// Preserve the original response in case refreshing fails. This keeps the
	// caller's error as the API's 401 instead of turning it into a transport
	// error.
	originalBody, _ := io.ReadAll(response.Body)
	response.Body.Close()
	response.Body = io.NopCloser(bytes.NewReader(originalBody))
	response.ContentLength = int64(len(originalBody))

	accessToken, refreshed := t.provider.refreshToken(staleToken)
	if !refreshed {
		return response, nil
	}

	retryRequest, err := requestWithAccessToken(request, accessToken)
	if err != nil {
		return response, nil
	}

	// Call the base transport directly so a second 401 is returned to the SDK
	// instead of causing an unbounded retry loop.
	return t.base.RoundTrip(retryRequest)
}

func bearerToken(value string) (string, bool) {
	parts := strings.SplitN(strings.TrimSpace(value), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	return token, token != ""
}

func isAuthenticationRequest(request *http.Request) bool {
	return strings.HasSuffix(request.URL.Path, "/oauth/v2/token")
}

func requestCanBeRetried(request *http.Request) bool {
	return request.Body == nil || request.GetBody != nil
}

func requestWithAccessToken(request *http.Request, accessToken string) (*http.Request, error) {
	retryRequest := request.Clone(request.Context())
	retryRequest.Header = request.Header.Clone()
	retryRequest.Header.Set("Authorization", "Bearer "+accessToken)

	if request.Body != nil {
		body, err := request.GetBody()
		if err != nil {
			return nil, err
		}
		retryRequest.Body = body
	}

	return retryRequest, nil
}

// refreshToken refreshes only when the token that failed is still cached. If
// another request already refreshed it, all concurrent callers reuse that
// newer token instead of issuing duplicate authentication requests.
func (xp *XplorProvider) refreshToken(staleToken string) (string, bool) {
	xp.authMutex.Lock()
	defer xp.authMutex.Unlock()

	if xp.token != nil && xp.token.Token != nil &&
		xp.token.Token.AccessToken != staleToken && xp.token.IsValid() {
		return xp.token.Token.AccessToken, true
	}

	executor := xp.providers.Get().(*xplorExecutor)
	defer xp.providers.Put(executor)

	token, err := executor.authenticate()
	if err != nil || token == nil || token.AccessToken == "" {
		return "", false
	}

	xp.token = &xplorentities.XPlorTokenWithTimestamp{
		Token:      token,
		ObtainedAt: time.Now(),
	}
	return token.AccessToken, true
}
