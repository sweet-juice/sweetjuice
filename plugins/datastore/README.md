# wails-mobile/datastore

This package provides a simple key-value data store using Android's `SharedPreferences`.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/datastore"
)

var dataStorePlugin = datastore.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    dataStorePlugin,
},

// Inside OnStart
if err := dataStorePlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
// Save data
dataStorePlugin.Set("user_token", "abc-123")

// Retrieve data
token, err := dataStorePlugin.Get("user_token", "")
```

## Usage (Frontend)

```js
// Set value
await SweetJuice.CallGo('DataStorePlugin.Set', "theme", "dark");

// Get value
const theme = await SweetJuice.CallGo('DataStorePlugin.Get', "theme", "light");
console.log("Current theme:", theme);
```
