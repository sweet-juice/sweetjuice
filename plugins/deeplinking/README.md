# wails-mobile/deeplinking

This package provides support for handling deep links in Sweet Juice applications.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/deeplinking"
)

var deepLinkPlugin = deeplinking.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    deepLinkPlugin,
},

// Inside OnStart
if err := deepLinkPlugin.Init(app); err != nil {
    return err
}

// Register a Go handler
deepLinkPlugin.OnURL(func(url string) {
    fmt.Println("Received Deep Link:", url)
})
```

## Android Configuration

To receive deep links, you must configure an intent filter in your `AndroidManifest.xml` within the `.WailsWebViewActivity` activity tag:

```xml
<intent-filter android:label="@string/app_name">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <!-- Accept URIs: myapp://example.com -->
    <data android:scheme="myapp" android:host="example.com" />
</intent-filter>
```

## Usage (Frontend)

```js
Wails.on("deeplinking:received", (url) => {
  console.log("Deep link received:", url);
  // Route to the appropriate page in your frontend
});
```
