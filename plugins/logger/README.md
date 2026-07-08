# wails-mobile/wails/logger

This package routes Go logs to the native Android `Logcat` system, allowing you to see your Go application output in ADB or Android Studio's Logcat terminal.

---

## Setup

Initialize the plugin in your `engine.go`. Note that since this is typically used only from Go, you don't need to add it to the `Bind` list unless you want to log from JS too.

```go
import (
	//...
	"github.com/sweet-juice/sweetjuice/plugins/logger"
	"github.com/sweet-juice/sweetjuice/core"
)

var log = logger.NewPlugin("MyAppTag")

// Inside OnStart
if err := log.Init(app); err != nil {
    return err
}
```

---

## Usage (Go)

The logger supports standard severity levels:

```go
log.Info("App started successfully")
log.Debug("Data payload: %v", myData)
log.Warn("Connectivity is slow...")
log.Error("Failed to save state: %v", err)
```

## Example output in ADB or Logcat:

```log
2026-05-23 22:56:26.275 13811-13954 MyAppTag  com.example.wailsmobile  [I]  App started successfully
```
---

## Benefits
By using this bridge instead of `fmt.Println`, your Go logs are correctly categorized by severity (Info, Debug, Error) and tagged with your application name, making them much easier to filter and debug in a real Android environment.
