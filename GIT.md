# Git Workflow for sb-storage

This repository contains the Slidebolt Storage service, which runs the storage server as a standalone process. It produces a standalone binary.

## Dependencies
- **Internal:**
  - `sb-contract`: Core interfaces and shared structures.
  - `sb-messenger-sdk`: Shared messaging interfaces and NATS implementation.
  - `sb-runtime`: Core execution environment and logging.
  - `sb-storage-sdk`: Shared interfaces and utilities for storage operations.
  - `sb-storage-server`: The actual storage implementation and indexing logic.
- **External:** 
  - `github.com/nats-io/nats.go`: Communication with NATS.

## Build Process
- **Type:** Go Application (Service).
- **Consumption:** Run as the primary data persistence service for Slidebolt.
- **Artifacts:** Produces a binary named `sb-storage`.
- **Command:** `go build -o sb-storage ./cmd/sb-storage`
- **Validation:** 
  - Validated through unit tests: `go test -v ./...`
  - Validated by successful compilation of the binary.

## Pre-requisites & Publishing
As the primary storage service, `sb-storage` must be updated whenever any of its internal SDKs or the core storage server implementation is changed.

**Before publishing:**
1. Determine current tag: `git tag | sort -V | tail -n 1`
2. Ensure all local tests pass: `go test -v ./...`
3. Ensure the binary builds: `go build -o sb-storage ./cmd/sb-storage`

**Publishing Order:**
1. Ensure all internal dependencies (`sb-contract`, `sb-messenger-sdk`, `sb-runtime`, `sb-storage-sdk`, `sb-storage-server`) are tagged and pushed.
2. Update `sb-storage/go.mod` to reference the latest tags.
3. Determine next semantic version for `sb-storage` (e.g., `v1.0.4`).
4. Commit and push the changes to `main`.
5. Tag the repository: `git tag v1.0.4`.
6. Push the tag: `git push origin main v1.0.4`.

## Update Workflow & Verification
1. **Modify:** Update storage service configuration or logic in `app/` or `cmd/`.
2. **Verify Local:**
   - Run `go mod tidy`.
   - Run `go test ./...`.
   - Run `go build -o sb-storage ./cmd/sb-storage`.
3. **Commit:** Ensure the commit message clearly describes the storage service change.
4. **Tag & Push:** (Follow the Publishing Order above).
