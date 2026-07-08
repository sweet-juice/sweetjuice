# wails-mobile/filepicker

This package provides access to the native Android file picker, allowing users to select files or media from their device.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/filepicker"
)

var filePickerPlugin = filepicker.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    filePickerPlugin,
},

// Inside OnStart
if err := filePickerPlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
// Trigger file picker
filePickerPlugin.PickFile(filepicker.PickerOptions{
    MimeType: "image/*",
    Multiple: false,
})
```

## Usage (Frontend)

```js
// Trigger picker
await Wails.CallGo('FilePickerPlugin.PickFile', {
    mime_type: "*/*",
    multiple: true
});

// Listen for results
Wails.on("filepicker:result", (result) => {
    if (result.error) {
        console.error("Picker error:", result.error);
        return;
    }
    
    if (result.multiple) {
        console.log("Selected URIs:", result.uris);
    } else {
        console.log("Selected URI:", result.uri);
    }
});
```

## Notes

Android returns URIs (e.g., `content://...`). To read the actual bytes in Go, you may need a separate "Content Provider" utility or use native bridge calls to open the URI stream.
