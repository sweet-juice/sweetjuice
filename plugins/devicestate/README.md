# wails-mobile/devicestate

This package provides access to Android device state information from Go and forwards state change events to the SweetJuice frontend.

## What it provides

- Battery level and charging status
- Battery health, status, and temperature
- Low power mode state
- Connectivity status, network type, roaming state, and metered/unmetered flags

## Setup

Initialize the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/devicestate"
    "github.com/sweet-juice/sweetjuice/core"
)

var deviceState = devicestate.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    helloService,
    deviceState,
    // ...other plugins
},

// Inside OnStart
if err := deviceState.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
state, err := deviceState.GetState()
if err != nil {
    log.Error("Failed to read device state: %v", err)
}

fmt.Printf("Battery: %d%% charging=%v network=%s\n", state.BatteryLevel, state.IsCharging, state.Connectivity.NetworkType)

_, err = deviceState.StartMonitoring()
```

## Events

The plugin emits `devicestate:changed` whenever the native Android state changes.

```js
SweetJuice.on("devicestate:changed", (data) => {
  console.log("Device state updated", data)
})
```

## Notes

This plugin is Android-only. It uses native Android platform APIs through the Sweet Juice bridge and does not require manual JNI implementation.
