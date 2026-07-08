# wails-mobile/wails/notification

This package allows the Go backend to post native system notifications on `**currently Android**`.

---

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
	//...
	"github.com/sweet-juice/sweetjuice/plugins/notification"
	"github.com/sweet-juice/sweetjuice/core"
)

var notifyPlugin = notification.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    notifyPlugin,
},

// Inside OnStart
if err := notifyPlugin.Init(app); err != nil {
    return err
}
```

---

## Basic Usage (Go)

```go
notifyPlugin.Post(notification.Notification{
    ID:    0, // Set to 0 to generate a unique ID automatically
    Title: "Hello World",
    Body:  "Message from Go",
    Importance: notification.ImportanceHigh,
})
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `int` | Unique notification ID. Use the same ID to update an existing notification. Set to `0` for auto-generation. |
| `Title` | `string` | The header text. |
| `Body` | `string` | The main message text. |
| `Importance` | `Importance`| `HIGH`, `DEFAULT`, `LOW`, `MIN`. |

---

## Usage from Frontend (JavaScript)

```js
await Wails.CallGo('NotificationPlugin.Post', {
    id: 0,
    title: "UI Alert",
    body: "Triggered by JS, posted by Go",
    importance: "HIGH"
});
```

---

## ⚠️ Android 13+ Requirements

On Android 13 (API 33) and higher, you must request the `POST_NOTIFICATIONS` permission before notifications will be displayed.

```go
// Using the Permission plugin
permPlugin.Request("android.permission.POST_NOTIFICATIONS")
```
