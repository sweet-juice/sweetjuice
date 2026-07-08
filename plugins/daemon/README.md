# wails-mobile/daemon

This package provides access to Android Foreground Services, allowing your Go backend to remain active and high-priority even when the UI is closed.

On Android Oreo (API 26) and above, it starts a **Foreground Service** with a persistent notification. On older versions, it falls back to a standard background service.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/daemon"
)

var daemonPlugin = daemon.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    daemonPlugin,
},

// Inside OnStart
if err := daemonPlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
// Start the service
daemonPlugin.Start(daemon.Options{
    Title:   "Music Sync",
    Message: "Uploading your library...",
})

// Stop the service
daemonPlugin.Stop()
```

## Android Manifest Requirements

You must register the service and the foreground permission in your `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.FOREGROUND_SERVICE" />
<!-- For Android 14+ -->
<uses-permission android:name="android.permission.FOREGROUND_SERVICE_SPECIAL_USE" />

<application ...>
    <service
        android:name="com.wailspackage.daemon.WailsDaemonService"
        android:enabled="true"
        android:exported="false"
        android:foregroundServiceType="specialUse">
        <property android:name="android.app.PROPERTY_SPECIAL_USE_FGS_SUBTYPE"
                  android:value="Wails Background Logic"/>
    </service>
</application>
```

## Why use Daemon?

Android is very aggressive at killing background processes to save battery. If your Go backend is performing a critical task (like a large file upload or WebSocket connection) and the user switches apps, the OS might kill your process. 

Starting a **Daemon** promotes your app's process to the "Foreground" state, making it highly unlikely to be killed by the system.
