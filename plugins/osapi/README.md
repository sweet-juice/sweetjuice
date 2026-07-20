# wails-mobile/osapi

This package provides information about the Android system version and device hardware.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/osapi"
)

var osApiPlugin = osapi.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    osApiPlugin,
},

// Inside OnStart
if err := osApiPlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
info, err := osApiPlugin.GetInfo()
if err == nil {
    fmt.Printf("OS: Android %s (SDK %d)\n", info.Release, info.SdkInt)
    fmt.Printf("Device: %s %s\n", info.Manufacturer, info.Model)
}
```

## Usage (Frontend)

```js
const info = await SweetJuice.CallGo('OsApiPlugin.GetInfo');
console.log(`Running on Android ${info.release} (API ${info.sdk_int})`);
```
