# Changelog

## [0.2.0] - 2026-09-02

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
