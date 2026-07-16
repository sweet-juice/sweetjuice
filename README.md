# Sweet Juice

Sweet Juice is an early port of Go v3 implementation to support mobile devies.

⚠️ **Important Note:** iOS support is currently in beta!

<p>
	<img alt="Example UI" src="./hello-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="Notification plugin" src="./notification-example.jpg" width="200" style="border-radius:6px;"/>
	<img alt="output console" src="./screenshot-output.jpg" width="200" style="border-radius:6px;"/>
</p>

<p align="center">
	<a href="https://github.com/sweet-juice/sweetjuice/blob/main/LICENSE">
	<img src="https://img.shields.io/github/license/sweet-juice/sweetjuice" alt="license"/>
	</a>
	<a href="https://goreportcard.com/report/github.com/sweet-juice/sweetjuice">
	<img src="https://goreportcard.com/badge/github.com/sweet-juice/sweetjuice" alt="goreport"/>
	</a>
	<a href="https://pkg.go.dev/github.com/sweet-juice/sweetjuice">
	<img src="https://pkg.go.dev/badge/github.com/sweet-juice/sweetjuice.svg" alt="Go Reference"/>
	</a>
	<a href="https://github.com/sweet-juice/sweetjuice/issues">
	<img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="contrib"/>
	</a>
</p>

## Quick links

1. [Install & Build](docs/INSTALL.md)
2. [CLI reference](docs/CLI.md)
3. [Contributing](docs/CONTRIBUTING.md)
4. [Plugin docs](docs/PLUGINS_DOCS.md)
5. [Android plugin guide](docs/PLUGINS_ANDROID.md)
6. [iOS plugin guide](#)
7. [Example Sweet Juice app](github.com/sweet-juice/examples)
8. [Benchmarks](docs/BENCHMARKS.md)

## Requirements

- Go 1.24+
- Android SDK & NDK (install via Android Studio or your preferred method)
- `git`, `curl`, `unzip`

## Pre-packed plugins

| Plugin | Android | iOS |
| :--- | :---: | :---: |
| `plugins/logger` | `Yes` | `Yes` |
| `plugins/notification` | `Yes` | `Yes` |
| `plugins/permission` | `Yes` | `Yes` |
| `plugins/special-permission` | `Yes` | `No` |
| `plugins/devicestate` | `Yes` | `Yes` |
| `plugins/workmanager` | `Yes` | `Yes` |
| `plugins/osapi` | `Yes` | `Yes` |
| `plugins/biometrics` | `Yes` | `Yes` |
| `plugins/filepicker` | `Yes` | `Yes` |

See [Full List](#)

## Notes

- For contribution guidelines and plugin conventions see `docs/CONTRIBUTING.md`.
