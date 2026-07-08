# Contributing to Sweet Juice

Short contribution guide and plugin conventions.

## Checklist

- Fork + branch: `git checkout -b feature/your-change`
- Format Go code: `gofmt -w ./`
- Run tests: `go test ./...`
- Keep PRs small and documented.

## Plugin conventions

- **Go Side**: Located in `plugins/<name>/`. Must import `github.com/sweet-juice/sweetjuice/core`.
- **Native Side**: 
    - Android source should be in a subfolder (e.g., `android/`). Must implement `SweetJuicePlugin`.
    - iOS source should be in a subfolder (e.g., `ios/`). Must implement `SweetJuicePlugin` protocol.
- **Registration**: Register new plugins in the `AppTemplate` to ensure they are available for new projects.

## Docs

- Keep README concise; place detailed install/cli/contrib docs in `docs/`.
- Update `DOCS.md` (generated via `gomarkdoc`) if core APIs change.
