# Changelog

## [0.1.44] - 2026-11-20

### Added

- Optional automatic retries for server errors (5xx) and network failures,
  enabled with `config.WithRetry(xplorcore.RetryPolicy{...})`. Retries use
  exponential backoff with jitter, honor `Retry-After`, only apply to
  idempotent methods, never retry client errors (4xx), and stay within each
  method's request timeout. Retries are disabled by default.
- `config.WithRequestTimeout(d)` to configure the total time budget of each
  SDK method call (default 30s, retries included).

### Fixed

- Timeout errors now report the configured timeout instead of always saying
  "10 seconds".
- Subscription requests no longer send `contact[]`, which the API rejects
  with 400.

### Breaking changes

- Removed `XPlorSubscriptionsParams.ContactIds`. Filtering subscriptions by
  several contacts at once is not supported by the API; use `ContactId` and
  make one request per contact.

## [0.1.43] - 2026-09-02

### Fixed

- Requests rejected with 401 because the access token expired are retried
  once after refreshing the token. Concurrent requests share a single
  refresh.

## [0.1.42] - 2026-09-02

### Added

- Added raw contact image downloads through `ContactImage` and
  `/files/contact_images/{pictureId}`.
- Added filtered article and counter-line collection requests for
  contact-scoped synchronization.
- Added course waiting-list recurrence and course waiting-list resources.
- Added isolated providers with `NewProvider` for separate API
  configurations in the same process.

### Changed

- `Articles` and `CounterLines` now accept collection filter parameters before
  pagination.
- `SharedResource` accepts the object payload returned by the API.
- Contact-scoped query parameters are serialized before pagination parameters.

### Breaking changes

- Update calls to `Articles` and `CounterLines` to pass the new params
  argument. Use `nil` when no collection filter is required.
- `ContactImage` now returns `[]byte` instead of a JSON image resource.
