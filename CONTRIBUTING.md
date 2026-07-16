# Contributing

Thanks for helping improve Sweet Juice. This file contains the project's contribution guidelines and the minimal steps to add or update plugins.

## Quick checklist

- Fork the repo and open a feature branch: `git checkout -b feature/your-change`
- Run `gofmt` / `go vet` and keep Go code compatible with Go 1.24+
- Run unit tests where available: `go test ./...`
- Keep commits small and focused; use clear messages
- Open a PR describing the change and any manual verification steps

## Writing plugins

Follow the existing conventions (see `plugins/PLUGINS.md` and `plugins/PLUGINS_ANDROID.md`):

- Plugin layout: `plugins/<name>/` with Go wrapper package and optional `android/` directory for Java/Kotlin sources.
- The Go wrapper should implement an `Init(app *wails.Application) error` method when needed and expose methods for frontend and other Go packages.
- Native Android side must implement the `SweetJuicePlugin` interface and return a stable domain string via `getDomain()`.
- Register your plugin in the sample application for manual testing by adding it to `SweetJuiceApplication.registerPlugin(...)` in `_examples/helloworld/native/android/app/src/main/java/com/example/juiceobile/SweetJuiceApplication.java`.

## Documentation

- Use `gomarkdoc` style for Go package docs (the repo already includes generated DOCS.md files). Keep README entries concise and link to deep-dive docs when required.
- If changes affect build or install instructions, update `docs/INSTALL.md` rather than the top-level README.

## Tests & Validation

- Prefer small, focused unit tests for Go packages.
- Manual validation for Android plugins: run the example app in `_examples/helloworld/` and exercise the plugin actions through the UI or via the `SweetJuicemobile` bridge.

## Code style

- Use `gofmt`/`goimports` and idiomatic Go naming.
- Java/Kotlin should follow standard Android style; keep UI work on the main thread and background work on background threads.

Thanks — your improvements make Sweet Juice better for everyone.
