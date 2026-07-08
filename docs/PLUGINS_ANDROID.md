# Android Plugin Development for Sweet Juice

This guide covers the technical details of implementing the native side of a Sweet Juice plugin on Android.

## The `SweetJuicePlugin` Interface

All Android plugins must implement the `com.sweetjuice.plugin.SweetJuicePlugin` interface:

```java
public interface SweetJuicePlugin {
    // Unique domain for routing (e.g., "camera", "permission")
    String getDomain();
    
    // Lifecycle hook for context attachment
    void onAttach(Context context);
    
    // Primary execution entry point from Go
    String handleAction(String action, String jsonArgsPayload);
    
    // Lifecycle hooks for Android results
    default void onActivityResult(int req, int res, Intent data) {}
    default void onRequestPermissionsResult(int req, String[] perms, int[] res) {}
    default void onNewIntent(Intent intent) {}
}
```

---

## Background-Safe Design

Android apps can be woken up in the background (e.g., by `WorkManager`). In these cases, there is **no Activity**. Your plugin must be designed to handle this.

### Context Awareness
`onAttach(Context context)` is called twice:
1.  **At Startup**: Called with the `Application` context. Use this for system services (e.g `NotificationManager`, `Log`).
2.  **On UI Open**: Called with the `AppCompatActivity` context. Use this for UI tasks (e.g `ActivityCompat.requestPermissions`, `startActivityForResult`).

### Example: Context Handling
```java
public class MyPlugin implements SweetJuicePlugin {
    private Context mContext;
    private AppCompatActivity mActivity;

    @Override
    public void onAttach(Context context) {
        this.mContext = context; // Always available
        if (context instanceof AppCompatActivity) {
            this.mActivity = (AppCompatActivity) context; // Available only in foreground
        }
    }

    @Override
    public String handleAction(String action, String args) {
        if ("doUiTask".equals(action)) {
            if (mActivity == null) return "{\"error\":\"App in background\"}";
            // Perform Activity-specific task
        }
        return "{\"status\":\"ok\"}";
    }
}
```

---

## Registration

### 1. The Global scope
Register your plugin in `SweetJuiceApplication.java` to ensure it is available for background workers:

```java
// Inside registerPlugins() in SweetJuiceApplication.java
private void registerPlugins() {
    registerPlugin(new MyPlugin());
}
```

### 2. The UI Scope
Ensure your plugin is re-attached to the Activity in `SweetJuiceActivity.java`:
This is handled automatically in the Sweet Juice `AppTemplate`. 
```java
// Inside onCreate
SweetJuiceApplication app = (SweetJuiceApplication) getApplication();
for (SweetJuicePlugin plugin : app.getPlugins().values()) {
    plugin.onAttach(this);
}
```

---

## Talking back to Go

To send data from Java back to Go asynchronously:

```java
// Wrap your arguments in a JSON array []
String payload = "[{\"result\": \"success\"}]";

// handleNativeAction routes directly to Go's RegisterNativeMethod handlers
sweetjuice.Sweetjuice.handleNativeAction("myplugin:callback", payload);
```

---

## Tips for Android
*   **Permissions**: Use the existing `PermissionsPlugin` to check/request neccessary permissions before your plugin executes.
*   **Main Thread**: `handleAction` is often called on a background bridge thread. If you need to touch the UI, use `new Handler(Looper.getMainLooper()).post(...)`.
*   **JSON Parsing**: Use `org.json.JSONObject` to parse `jsonArgsPayload`.

---

## 🛠️ Background Services & Daemons

If your plugin requires a persistent background presence, you should implement a `Service`.

1.  **Foreground Service**: Required for persistent background logic on modern Android (Oreo+). You must show a notification.
2.  **Manifest Registration**: All services must be declared in `AndroidManifest.xml`.
3.  **Permissions**: Request `FOREGROUND_SERVICE` and specific types (e.g., `specialUse` for custom Go logic) in the manifest.
