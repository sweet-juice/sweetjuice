# wails-mobile/biometric

This package provides access to biometric authentication (Fingerprint, Face ID, etc.) for Sweet Juice applications.

## Setup

Initialize and bind the plugin in your `engine.go`:

```go
import (
    "github.com/sweet-juice/sweetjuice/plugins/biometric"
    "github.com/sweet-juice/sweetjuice/core"
)

var biometricPlugin = biometric.NewPlugin()

// Inside StartApplication Bind list
Bind: []interface{}{
    biometricPlugin,
},

// Inside OnStart
if err := biometricPlugin.Init(app); err != nil {
    return err
}
```

## Usage (Go)

```go
// Check if biometric auth is available
res, err := biometricPlugin.CanAuthenticate()
if res.Can {
    // Trigger authentication
    biometricPlugin.Authenticate(biometric.AuthOptions{
        Title: "Login Required",
        Description: "Please authenticate to continue",
    })
}
```

## Usage (Frontend)

```js
// Check availability
const res = await SweetJuice.CallGo('BiometricPlugin.CanAuthenticate');
if (res.can_authenticate) {
    // Start auth
    await SweetJuice.CallGo('BiometricPlugin.Authenticate', {
        title: "Confirm Action",
        negative_button_text: "Cancel"
    });
}

// Listen for results
SweetJuice.on("biometric:result", (result) => {
    if (result.success) {
        console.log("Authenticated!");
    } else {
        console.error("Auth failed:", result.error);
    }
});
```

## Android Requirements

Add the biometric permission to your `AndroidManifest.xml`:

```xml
<uses-permission android:name="android.permission.USE_BIOMETRIC" />
```

Ensure you have the `androidx.biometric:biometric` dependency in your `app/build.gradle`.
