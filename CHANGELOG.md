# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.21.47] - 2026-09-22

### Added

- Initial public release of the FalaAI API Go SDK.
- Typed client for the three FalaAI API endpoints:
  - `POST /v1/audio/transcriptions` - audio to text, with speaker separation
  - `POST /v1/analyze/diagnostic` - conversation analysis
  - `POST /v1/analyze/auditoriaRisco` - compliance audit
- Generated from the FalaAI API OpenAPI contract, version 1.21.47.
- End-to-end test suite (19 tests).

[1.21.47]: https://github.com/ActionTecBr/falaai-api/releases/tag/v1.21.47
