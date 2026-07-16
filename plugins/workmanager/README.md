# wails-mobile/wails/workmanager

This package provides a wrapper for the Android `WorkManager` API, allowing you to run Go logic in the background even when the application is closed.

---

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
	//...
	"github.com/sweet-juice/sweetjuice/plugins/workmanager"
	"github.com/sweet-juice/sweetjuice/core"
)

var wmPlugin = workmanager.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    wmPlugin,
},

// Inside OnStart
if err := wmPlugin.Init(app); err != nil {
    return err
}
```

---

## Registering Background Tasks

You must register the Go function that should run when the Android OS triggers the background worker.

```go
wmPlugin.RegisterTask("sync_data", func() error {
    fmt.Println("Running background sync...")
    // Perform sync logic here
    return nil // Return error to trigger an OS retry
})
```

---

## Scheduling Tasks

### Periodic Work
Schedule a task to run every N minutes (Minimum 15 minutes enforced by Android).

```go
// From Go
wmPlugin.EnqueuePeriodic("sync_data", 15, nil) //15 minutes minimum
```

```js
// From JS
await SweetJuice.CallGo('WorkManagerPlugin.EnqueuePeriodic', "sync_data", 15);
```

### One-Time Work
Schedule a task to run once as soon as possible.

```go
wmPlugin.EnqueueOneTime("sync_data", nil)
```

---

## Using Constraints

You can restrict tasks to run only under certain conditions (e.g., only on Wi-Fi).

```go
constraints := &workmanager.Constraints{
    NetworkType:      workmanager.NetworkUnmetered, // Wi-Fi only
    RequiresCharging: true,                         // Plugged in only
}

wmPlugin.EnqueuePeriodic("sync_data", 60, constraints)
```

**JS Constraints Object:**
```js
const constraints = {
    network_type: "UNMETERED",
    requires_charging: true
};
await SweetJuice.CallGo('WorkManagerPlugin.EnqueuePeriodic', "sync_data", 60, constraints);
```

---

## Stopping Work

```go
// Cancel all scheduled tasks
wmPlugin.CancelAll()
```
