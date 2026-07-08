# Install & Build

This document contains the installation, bootstrapping, build and uninstall steps used during development.

## Prerequisites

- Go 1.24+
- Android SDK and Android NDK installed (Android Studio recommended)
- `git`, `curl`
- **For less mess, do not alter the default location for Android SDK**

## CLI installation

You can use `Go` to install `Sweet Juice` CLI tool `wailsm`.

* Using `Go`:

```bash
go install github.com/sweet-juice/sweetjuice/cmd/wailsm@latest
```

## Bootstrapping a new project

```bash
wailsm --new my-new-app
```

This scaffolds a fresh project from the local `AppTemplate` included in the Sweet Juice module.

### Sweet Juice app structure

- `core/` — core Go runtime package for mobile bridged apps
- `plugins/` — Go-side implementations of native plugins
- `native/android/` — Android Studio project using the generated AAR
- `native/ios/` — iOS project using the generated XFramework
- `frontend/` — web UI assets

## Refresh / Sync (generate native bindings)

The CLI performs an automatic sync before `build` and `run`. Manually trigger sync when you want to force a regeneration:

```bash
wailsm --refresh android
# or
wailsm --refresh ios
```

This drops the resulting `.aar` (Android) or `.xcframework` (iOS) into the respective native project directories.

## Running App Immediately

Run on a connected device:

```bash
# For Android (via ADB)
wailsm --run android

# For iOS (via xtool)
wailsm --run ios
```

## Build release binaries:

```bash
# Android
wailsm --build android debug
wailsm --build android release
wailsm --build android bundle

# iOS
wailsm --build ios debug
wailsm --build ios release
```

## To Uninstall Sweet Juice CLI:

```shell
rm $(go env GOPATH)/bin/wailsm     
```


##  Notes & troubleshooting

- Ensure `ANDROID_HOME` or equivalent environment variables are set or Android Studio is installed in default locations.
- For iOS, ensures `xtool` and Xcode Command Line Tools are installed.
- If bindings seem stale, run `wailsm --refresh <platform>` to force a rebuild.

For additional developer notes see `docs/CONTRIBUTING.md` and the project `plugins/` directory.
