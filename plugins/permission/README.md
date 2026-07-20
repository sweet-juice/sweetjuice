# wails-mobile/wails/permission

This is an officially maintained `Permission` package for Android, allowing the Go backend to manage system runtime permissions.

---

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
	//...
	"github.com/sweet-juice/sweetjuice/plugins/permission"
	"github.com/sweet-juice/sweetjuice/core"
)

var permPlugin = permission.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    permPlugin,
},

// Inside OnStart
if err := permPlugin.Init(app); err != nil {
    return err
}
```

---

## Basic Usage (Go)

You can check or request permissions directly from your Go business logic.

```go
// Check if camera permission is granted
status, err := permPlugin.Check("android.permission.CAMERA")

// Request camera permission
// Results will be sent asynchronously via the "permissions:changed" event
permPlugin.Request("android.permission.CAMERA")
```

---

## Usage from Frontend (JavaScript)

The frontend should strictly use the `SweetJuice.CallGo` bridge.

```js
// Requesting permission
await SweetJuice.CallGo('PermissionPlugin.Request', "android.permission.CAMERA");

// Checking status
const result = await SweetJuice.CallGo('PermissionPlugin.Check', "android.permission.CAMERA");
const parsed = JSON.parse(result);
console.log("Status:", parsed.status);
```

### Listening for Results
When a permission is requested, the result is delivered as an event:

```js
SweetJuice.on('permissions:changed', (data) => {
    console.log(`Permission for ${data.permission} is ${data.granted ? 'GRANTED' : 'DENIED'}`);
});
```

---

## Android Manifest Requirements

Any permission you request at runtime **must** also be declared in your `AndroidManifest.xml`:

example: 

```xml
<uses-permission android:name="android.permission.CAMERA" />
<uses-permission android:name="android.permission.POST_NOTIFICATIONS" />
```
