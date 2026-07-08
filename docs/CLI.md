# CLI Reference

Compact reference and useful examples for the `wailsm` CLI.

Synopsis

```
Sweet Juice Toolchain CLI (wailsm)
Usage:
  wailsm --new <project_name>        Create a fresh project from the template
  wailsm --refresh <platform>        Run platform sync: 'android' or 'ios'
  wailsm --build <platform> <mode>   Sync environment and compile app binary (debug/release/bundle)
  wailsm --run <platform>            Compile, install, and execute application via ADB
  wailsm --add <plugin-url>          Install a native Go/Mobile plugin
  wailsm --remove <plugin-url>       Uninstall a native Go/Mobile plugin
```

Basic examples

```bash
# Create a new project
wailsm --new my-app

# Build debug APK
wailsm --build android debug

# Build and run (the CLI performs sync automatically)
wailsm --build android debug && wailsm --run android

# Build release bundle (AAB)
wailsm --build android bundle

# Add external plugin
wailsm --add github.com/author/plugin

# Remove plugin
wailsm --remove github.com/author/plugin
```

When to use `--refresh`

- The CLI automatically runs sync during `build`/`run`. Use `wailsm --refresh android` only when you want to force a rebuild of the generated AAR/JAR artifacts (for debugging or template changes).

Environment notes

- Ensure Android SDK and NDK are installed and available to the environment (Android Studio recommended).
- If builds fail due to Gradle, open `native/android/` in Android Studio and inspect the Gradle wrapper and SDK/NDK settings.

CI / scripts

- For CI pipelines prefer explicit commands: `wailsm --refresh android && wailsm --build android release` to keep steps visible.
