# wails-mobile/datadir

This package provides access to standard Android application directories (internal/external files and cache).

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/datadir"
)

var dataDirPlugin = datadir.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    dataDirPlugin,
},

// Inside OnStart
if err := dataDirPlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
dirs, err := dataDirPlugin.GetDirs()
if err == nil {
    fmt.Println("Files Dir:", dirs.Files)
}
```

## Usage (Frontend)

```js
const dirs = await Wails.CallGo('DataDirPlugin.GetDirs');
console.log("Internal Files Dir:", dirs.files);
```
