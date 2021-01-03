# Changelog

## [unreleased] - YYYY-MM-DD

### Added

- Default template supports resolved alerts
- Helpers on alert group model to group alerts by status.

## [0.3.1] - 2019-06-02

### Changed

- Update libraries.
- Use Go 1.14.

## [0.3.0] - 2019-12-20

### Added

- Forward alerts by custom chat ID based on alert labels.

### Changed

- Telegram required flags/envs are not required when using notify dry run mode.

## [0.2.1] - 2019-12-17

### Added

- Metrics to dead man's switch service.

### Changed

- Dead man's switch default interval is 15m.

## [0.2.0] - 2019-12-16

### Added

- Dead man's switch option with Alertmanager.
- Alertmanager API accepts a query string param with a custom chat ID.
- Telegram notifier can send to customized chats.

## [0.1.0] - 2019-12-13

### Added

- Custom templates.
- Docs for Kubernetes deployment and Alertmanager configuration.
- Simple health check
- Prometheus metrics.
- Telegram notifier.
- Dry-run notifier.
- Default message template.
- Alertmanager compatible webhook API.
- Models and forwarding domain service.

[unreleased]: https://github.com/slok/alertgram/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/slok/alertgram/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/slok/alertgram/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/slok/alertgram/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/slok/alertgram/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/slok/alertgram/releases/tag/v0.1.0
