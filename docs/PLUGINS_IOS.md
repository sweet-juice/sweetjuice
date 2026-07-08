# iOS Plugin Development for Sweet Juice

This guide covers the technical details of implementing the native side of a Sweet Juice plugin on iOS.

## 📋 The Swift Interface

All iOS plugins must implement the `SweetJuicePlugin` protocol:

```swift
import UIKit

public protocol SweetJuicePlugin {
    func getDomain() -> String
    func onAttach(container: UIViewController)
    func handleAction(action: String, jsonArgs: String) -> String
}
```

---

## 🛠️ Implementation Example

```swift
class MyPlugin: SweetJuicePlugin {
    private var container: UIViewController?

    func getDomain() -> String {
        return "myplugin"
    }

    func onAttach(container: UIViewController) {
        self.container = container
    }

    func handleAction(action: String, jsonArgs: String) -> String {
        if action == "hello" {
            return "{\"message\":\"Hello from iOS!\"}"
        }
        return "{\"error\":\"Unknown action\"}"
    }
}
```

---

## 🏗️ Registration

Register your plugin in the `SweetJuiceManager`:

```swift
// Typically done in your App structure or a dedicated initializer
SweetJuiceManager.shared.registerPlugin(MyPlugin())
```

The `SweetJuiceWebView` will automatically call `onAttach` for all registered plugins when the UI is ready.

---

## 📡 Talking back to Go

To send data from Swift back to Go asynchronously:

```swift
let payload = "[{\"result\": \"success\"}]"
SweetjuiceHandleNativeAction("myplugin:callback", payload)
```
