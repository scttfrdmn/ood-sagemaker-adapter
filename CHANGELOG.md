# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- `launch` now waits for the SageMaker app to reach `InService` (polling `DescribeApp`,
  bounded by `--wait-timeout`, default 5m) before returning the presigned Studio URL.
  Previously the URL was generated immediately after `CreateApp`, so it often opened to a
  Studio loading spinner that timed out (#2). Errors out on a terminal non-serviceable
  status or on timeout.

### Fixed
- Repinned `github.com/scttfrdmn/substrate` 0.45.2 → 0.65.0 and regenerated go.sum. The v0.45.2 tag content was changed upstream (substrate#296), so the recorded checksum no longer matched and `go test -tags=integration` failed with a go.sum SECURITY ERROR. Integration tests now build and pass.

### Added
- Initial scaffold — OOD compute adapter for AWS SageMaker Studio, translating Open OnDemand interactive app requests to SageMaker API calls and returning presigned Studio URLs.
- CLI commands: `launch` (create a SageMaker app → presigned Studio URL), `status <app-id>` (OOD-normalized status), and `delete <app-id>` (delete a Studio app).
- Substrate integration tests for the SageMaker Studio app lifecycle.
- CI workflow with pinned action SHAs.
